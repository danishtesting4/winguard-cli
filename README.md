# WinGuard CLI - Free Lightweight VPN for Windows

## Overview
WinGuard is a free, open-source CLI VPN for Windows that bypasses country-level Deep Packet Inspection (DPI) blocks while being lightweight and fast.

## Features
- **WireGuard protocol** - Modern, fast, lightweight VPN protocol
- **wstunnel obfuscation** - Wraps WireGuard traffic in WebSocket HTTPS frames
- **LWO obfuscation** - Lightweight WireGuard Obfuscation (header scrambling)
- **Single binary** - Small footprint (< 10MB compressed)
- **Cross-platform** - Windows 10/11 (also works on Linux/macOS)
- **No ads, no tracking** - Fully open source and free

## Quick Start

### 1. Prerequisites
- Windows 10/11
- [Go 1.21+](https://go.dev/dl/) (for building from source)
- [WireGuard keys](https://www.wireguard.com/quickstart/) (generate with `wg genkey`)
- Administrator privileges (for TUN adapter creation)

### 2. Generate Keys
```powershell
# Generate private key
wg genkey | set-clipboard

# Generate public key
wg pubkey | set-clipboard
```

### 3. Build WinGuard
```powershell
go mod init winguard-go
go get git.zx2c4.com/wireguard-go
go get github.com/spf13/cobra
go build -o winguard.exe ./src
```

### 4. Connect to VPN
```powershell
# Basic connection
.\winguard connect -c config.yaml

# With wstunnel obfuscation (recommended for country blocks)
.\winguard connect --obfuscation wstunnel

# With LWO obfuscation (good for ISP DPI)
.\winguard connect --obfuscation lwo

# Just status
.\winguard status

# Disconnect
.\winguard disconnect
```

### 5. Sample Configuration (config.yaml)
```yaml
private_key: "your_private_key_here"
address: "10.0.0.2/24"
dns:
  - "1.1.1.2"
  - "8.8.8.8"
listen_port: 51820
peer_public_key: "your_peer_public_key_here"
endpoint: "your-vpn-server:51820"
mtu: 1280
keep_alive: 25
obfuscation: "wstunnel"
```

## Obfuscation Modes

### wstunnel (Recommended for country-level censorship)
- Wraps WireGuard packets in WebSocket HTTPS frames
- Makes traffic look like regular web browsing
- Effective against: National firewalls (China, Iran, Russia), ISP DPI
- Requires: wstunnel server running on your VPN endpoint

### LWO (Lightweight WireGuard Obfuscation)
- Scrambles WireGuard packet headers in-place
- Minimal overhead, fast
- Effective against: ISP-level DPI, simpler blocks
- No extra infrastructure needed

## Building from Source

### On Windows with Go:
```powershell
# 1. Install Go
# 2. Set up module
go mod init winguard-go

# 3. Add dependencies
go get git.zx2c4.com/wireguard-go
go get github.com/spf13/cobra@latest
go get serde
go get serde_yaml
go get tokio
go get rand
go get base64

# 4. Build
go build -o winguard.exe ./src
```

### Expected Binary Size: 3-8MB (stripped)

## Technology Stack

| Component | Technology |
|-----------|-----------|
| VPN Protocol | WireGuard (git.zx2c4.com/wireguard-go) |
| Obfuscation | wstunnel (WebSocket wrapping), LWO (header scrambling) |
| CLI Framework | spf13/cobra |
| Config Parsing | serde + serde_yaml |
| async runtime | tokio |
| Windows integration | golang.org/x/sys/windows |
| TUN device | tun crate |

## How It Bypasses DPI

### wstunnel Mode:
1. WireGuard encrypts your traffic
2. wstunnel wraps each WG packet in a WebSocket frame
3. The frame uses standard HTTPS/WSS protocol
4. DPI sees only legitimate HTTPS traffic, not VPN fingerprints
5. Traffic flows freely through the censorship layer

### LWO Mode:
1. WireGuard encrypts your traffic
2. LWO scrambles packet headers (ICMP type, flags, etc.)
3. DPI can't identify the WireGuard signature
4. Traffic appears as random modified packets
5. Bypasses basic DPI fingerprinting

## Use Cases
- Bypassing national internet censorship
- Avoiding ISP-level website blocking
- Secure remote access to home/office network
- Privacy-focused browsing
- Bypassing hotel/airport WiFi restrictions

## License
MIT License - Free for personal and commercial use

## Support
- Issues: GitHub Issues
- Documentation: README.md
- Community: Discussions

## Roadmap (Planned Features)
- [ ] GUI system tray icon (optional)
- [ ] Auto-start configuration
- [ ] Multiple obfuscation chaining
- [ ] Config wizard for different regions
- [ ] Port forwarding support
- [] DNS leak protection