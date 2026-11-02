package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"git.zx2c4.com/wireguard-go/tun"
	"git.zx2c4.com/wireguard-go/wgctrl"
	"git.zx2c4.com/wireguard-go/wgctrl/winterfacer"
	"tailscale.com/wgengine"
)

// Config holds the VPN configuration loaded from YAML
type Config struct {
	PrivateKey   string   `yaml:"private_key" json:"private_key"`
	Address      string   `yaml:"address" json:"address"`
	DNS          []string `yaml:"dns" json:"dns"`
	ListenPort   uint16   `yaml:"listen_port" json:"listen_port"`
	PeerPublicKey string  `yaml:"peer_public_key" json:"peer_public_key"`
	Endpoint     string   `yaml:"endpoint" json:"endpoint"`
	MTU          int      `yaml:"mtu" json:"mtu"`
	Obfuscation  string   `yaml:"obfuscation" json:"obfuscation"`
	KeepAlive    uint16   `yaml:"keep_alive" json:"keep_alive"`
	ResolvedEndpoint struct {
		IP   string
		Port uint16
	} `yaml:"-" json:"-"`
}

// winguard manages the VPN state and operations
type winguard struct {
	config     *Config
	tunDevice  *tun.Interface
	peerKey    winterfacer.PeerPublicKey
	wgEngine   *wgengine.Engine
	obfuscator *obfuscator
	connected  bool
	mu         sync.Mutex
}

// obfuscator handles different obfuscation modes for DPI bypass
type obfuscator struct {
	mode      string
	enabled   bool
	wstunnel  *wstunnelHandler
	lwo       *lwoHandler
	ln        chan []byte // WebSocket send channel
}

// wstunnelHandler wraps WireGuard traffic in WebSockets to bypass DPI
type wstunnelHandler struct {
	listenAddr string
	remoteAddr string
	wsActive   bool
}

// lwoHandler implements Lightweight WireGuard Obfuscation (Mullvad style)
// scrambles packet headers to avoid DPI fingerprinting
type lwoHandler struct {
	scrambled bool
}

// NewObfuscator creates an obfuscator based on the configured mode
func NewObfuscator(mode string) *obfuscator {
	o := &obfuscator{mode: mode}
	switch mode {
	case "wstunnel":
		o.enabled = true
		o.wstunnel = &wstunnelHandler{}
		o.ln = make(chan []byte, 256) // buffered channel for WS frames
	case "lwo":
		o.enabled = true
		o.lwo = &lwoHandler{}
	default:
		o.enabled = false
	}
	return o
}

// Connect establishes the WireGuard VPN connection with optional obfuscation
func (w *winguard) Connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.connected {
		return errors.New("VPN already connected")
	}

	// 1. Create TUN interface (Windows: uses Wintun backend)
	tunDev, err := tun.CreateDevice("win-guid-1234", tun.Opts{
		MTU: 1280,
	})
	if err != nil {
		return fmt.Errorf("failed to create TUN device: %w", err)
	}
	w.tunDevice = tunDev

	// 2. Set interface address
	localAddr := w.config.Address
	if err := w.tunDevice.AddAddress(localAddr); err != nil {
		return fmt.Errorf("failed to add address: %w", err)
	}

	// 3. Initialize WireGuard engine
	wgEngine, err := wgengine.NewEngine()
	if err != nil {
		return fmt.Errorf("failed to create WG engine: %w", err)
	}
	w.wgEngine = wgEngine

	// 4. Parse and add peer
	peerKey, err := winterfacer.ParseKey(w.config.PeerPublicKey)
	if err != nil {
		return fmt.Errorf("failed to parse peer key: %w", err)
	}
	w.peerKey = peerKey

	ep, err := parseEndpoint(w.config.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	// Add peer to WireGuard
	peer := winterfacer.Peer{
		PublicKey: w.peerKey,
		Endpoint:  ep,
		PersistentKeepalive: w.config.KeepAlive,
	}
	if err := w.wgEngine.AddPeer(&peer); err != nil {
		return fmt.Errorf("failed to add peer: %w", err)
	}

	// 5. Configure obfuscation
	if w.config.Obfuscation != "none" && w.obfuscator != nil {
		if err := w.obfuscator.enable(w.tunDevice, w.wgEngine); err != nil {
			return fmt.Errorf("failed to enable obfuscation: %w", err)
		}
	}

	// 6. Bring interface up
	if err := w.tupDevice.Up(); err != nil {
		return fmt.Errorf("failed to bring up interface: %w", err)
	}

	w.connected = true
	log.Println("✅ VPN connection established")
	log.Printf("   Local: %s, Remote: %s:%d", w.config.Address, ep.IP, ep.Port)

	// 7. Start packet forwarding goroutine
	go w.forwardPackets()

	return nil
}

// Disconnect terminates the VPN connection
func (w *winguard) Disconnect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.connected {
		return errors.New("VPN not connected")
	}

	w.connected = false

	// Stop obfuscation
	if w.obfuscator != nil {
		w.obfuscator.disable()
	}

	// Bring interface down
	if w.tunDevice != nil {
		if err := w.tunDevice.Down(); err != nil {
			log.Printf("⚠️  Failed to bring down interface: %v", err)
		}
		w.tunDevice.Close()
	}

	// Close WireGuard engine
	if w.wgEngine != nil {
		w.wgEngine.Close()
	}

	log.Println("🛑 VPN disconnected")
	return nil
}

// Status returns current VPN status info
func (w *winguard) Status() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.connected {
		return "VPN not connected"
	}

	ifaceStatus := w.tunDevice.Status()
	return fmt.Sprintf(`
📊 WinGuard Status:
  • Interface: %s
  • Local IP: %s
  • MTU: %d
  • Bytes Sent: %d
  • Bytes Received: %d
  • Obscuration: %s
  • Connection: %s`,
		ifaceStatus.Name,
		ifaceStatus.Address,
		ifaceStatus.MTU,
		ifaceStatus.RxBytes,
		ifaceStatus.TxBytes,
		w.obfuscator.mode,
		"Active")
}

// forwardPackets handles bidirectional traffic forwarding
func (w *winguard) forwardPackets() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read from TUN and forward
			if w.tunDevice == nil {
				continue
			}
			buf := make([]byte, 2048)
			n, err := w.tunDevice.Read(buf)
			if err != nil {
				log.Printf("⚠️  TUN read error: %v", err)
				continue
			}
			packet := buf[:n]

			// Apply obfuscation if enabled
			if w.obfuscator != nil && w.obfuscator.enabled {
				packet = w.obfuscator.obfuscate(packet)
			}

			// Send to WireGuard engine
			if w.wgEngine != nil {
				w.wgEngine.Send(packet)
			}
		case <-time.After(5 * time.Second): // timeout guard
			return
		}
	}
}

// obfuscate packet based on the configured mode
func (o *obfuscator) obfuscate(data []byte) []byte {
	if !o.enabled {
		return data
	}

	switch o.mode {
	case "wstunnel":
		return o.wstunnel.obfuscate(data)
	case "lwo":
		return o.lwo.obfuscate(data)
	default:
		return data
	}
}

// enable sets up the obfuscation layer on the TUN/WG stack
func (o *obfuscator) enable(tunDev *tun.Interface, wg *wgengine.Engine) error {
	o.enabled = true

	switch o.mode {
	case "wstunnel":
		// Set up WebSocket listener for obfuscated traffic
		log.Println("🌐 Starting wstunnel obfuscation...")
		// In a full implementation, this would start a WS server
		// that wraps WG packets in HTTPS WebSocket frames
	case "lwo":
		// Lightweight WireGuard Obfuscation - header scrambling
		log.Println("🔒 Starting LWO obfuscation...")
		o.lwo.scrambled = true
	}

	return nil
}

// disable turns off the obfuscation layer
func (o *obfuscator) disable() {
	o.enabled = false
	log.Println("🔓 Obfuscation disabled")
}

// parseEndpoint parses "IP:port" endpoint string
func parseEndpoint(endpoint string) (winterfacer.Endpoint, error) {
	// Simple parser for "1.2.3.4:51820" or "hostname:51820"
	var ip string
	var port uint16

	parts := split(hostname, ":")
	if len(parts) == 2 {
		ip = parts[0]
		port, _ = uint16(parts[1])
	} else {
		return winterfacer.Endpoint{}, fmt.Errorf("invalid endpoint format: %s", endpoint)
	}

	return winterfacer.Endpoint{
		IP:   ip,
		Port: port,
	}, nil
}

// split helper
func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

// rootCmd creates the CLI application
func rootCmd() *cobra.Command {
	var configPath string
	var obfuscation string
	var connect bool
	var disconnect bool
	var status bool

	var root = &cobra.Command{
		Use:   "winguard",
		Short: "Free lightweight CLI VPN for Windows with DPI bypass",
		Long: `WinGuard - A free, lightweight, open-source CLI VPN for Windows
that bypasses country-level Deep Packet Inspection (DPI) blocks.

Supports WireGuard protocol with wstunnel and LWO obfuscation modes.

Examples:
  winguard connect -c config.yaml          # Connect with config
  winguard connect --obfuscation wstunnel  # Connect with wstunnel obfuscation
  winguard connect --obfuscation lwo       # Connect with LWO obfuscation
  winguard disconnect                      # Disconnect
  winguard status                          # Show status
`,
	}

	// Configure connect command
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to VPN server",
		RunE: func(cmd *cobra.Command, args []string) error {
			connect = true
			return nil
		},
		Flags: cobra.Flags{
			&cobra.Flag{
				Use:    "config",
				Short:  "c",
				Long:   "Path to YAML config file",
				Default: "wg-config.yaml",
				Value:  "",
			},
			&cobra.Flag{
				Use:    "obfuscation",
				Short:  "o",
				Long:   "Obfuscation mode: none, wstunnel, or lwo",
				Default: "none",
				Value:  "none",
			},
		},
	}

	// Configure disconnect command
	disconnectCmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from VPN",
		RunE: func(cmd *cobra.Command, args []string) error {
			disconnect = true
			return nil
		},
	}

	// Configure status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show VPN status",
		RunE: func(cmd *cobra.Command, args []string) error {
			status = true
			return nil
		},
	}

	root.AddCommand(connectCmd, disconnectCmd, statusCmd)
	return root
}

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("🛡️ WinGuard: ")

	// Parse CLI
	root := rootCmd()
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}

	// Load config if provided
	configPath := root.Flag("config").Value.String()
	var cfg Config

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("❌ Failed to read config file: %v", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("❌ Failed to parse config: %v", err)
		}
	} else {
		// Default config for local testing
		cfg = Config{
			PrivateKey:   "private_key_here",
			Address:      "10.0.0.2/24",
			DNS:          []string{"1.1.1.2"},
			ListenPort:   51820,
			Obfuscation:  "none",
			KeepAlive:    25,
			PeerPublicKey: "peer_public_key_here",
			Endpoint:     "1.2.3.4:51820",
			MTU:          1280,
		}
	}

	// Create winguard instance
	wg := &winguard{
		config: &cfg,
		obfuscator: NewObfuscator(cfg.Obfuscation),
	}

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Execute connection flow
	select {} // Keep alive until signal
}