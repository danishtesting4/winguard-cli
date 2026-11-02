# WinGuard CLI - Free Lightweight VPN for Windows

## 📦 Quick Download (Recommended)
Pre-built binaries are available in the [Releases](https://github.com/danishtesting4/winguard-cli/releases) section. Just download the `.exe` and run it - no installation required!

## 🚀 One-Command Install Methods
Choose your preferred installation method:

### ☁️ Option 1: PowerShell One-Liner (Easiest)
```powershell
iwr https://bit.ly/winguard-install -UseBasicParsing | iex
# OR run the included script:
.\install.ps1 -Method binary   # Download pre-built binary (recommended)
.\install.ps1 -Method build    # Build from source (requires Go)
```

### 📦 Option 2: Chocolatey Package Manager
```powershell
choco install winguard-cli
```
*Uses the `winguard-cli.nuspec` in this repo. Requires Chocolatey installed.*

### 📦 Option 3: Winget Package Manager
```powershell
winget install danishtesting4.WinGuardCLI
```
*Will be available once package is published to Microsoft Store/Repository.*

### 🔨 Option 4: Batch File Installer (Offline)
```powershell
.\wininstall.bat
```
*Full installer with GUI prompts. Downloads binary and optionally adds to PATH.*

### 📁 Option 5: Manual Binary Download (No Package Manager)
```powershell
# 1. Download single .exe from Releases
# 2. Run as Administrator:
Right-Click winguard.exe → Run as administrator

# 3. Connect immediately:
winguard connect --obfuscation wstunnel
```

### 🛠️ Option 6: Build from Source (Full Control)
```powershell
# Requires Go 1.21+
go get git.zx2c4.com/wireguard-go
go build -ldflags="-s -w" -o winguard.exe ./src
# Optional: & .\upx.exe --best winguard.exe  # Compress to ~3-5MB

# Result: Standalone .exe ready to use
.\winguard connect --obfuscation wstunnel
```

## 📋 Recommended Method
**PowerShell one-liner** (`.\install.ps1 -Method binary`) - Most users, 30 seconds, no dependencies.

## 📦 All Methods Summary

| Method | Time | Dependencies | Binary Size | Best For |
|--------|------|-------------|-------------|----------|
| **PowerShell script** | 30 sec | None | 3-8MB | Most users |
| **Chocolatey** | 45 sec | Chocolatey | 3-8MB | choco users |
| **Winget** | 30 sec | Winget | 3-8MB | winget users |
| **Batch installer** | 2 min | None | 3-8MB | Offline/GUI preference |
| **Manual binary** | 5 min | None | 3-8MB | Full control |
| **Build from source** | 2-5 min | Go 1.21+ | 3-8MB | Developers, customization |

## 💡 Quick Start - Choose One

**Fastest (PowerShell):**
```powershell
iwr https://github.com/danishtesting4/winguard-cli/raw/main/install.ps1 | iex
# Then run:
.\winguard connect --obfuscation wstunnel
```

**Easiest (Binary):**
```powershell
# 1. Go to: https://github.com/danishtesting4/winguard-cli/releases
# 2. Download winguard-windows-amd64.exe
# 3. Run as Administrator
# 4. winguard connect --obfuscation wstunnel
```

**Most Customizable (Build):**
```powershell
go get git.zx2c4.com/wireguard-go
go build -ldflags="-s -w" -o winguard.exe ./src
& .\upx.exe --best winguard.exe
.\winguard connect --obfuscation wstunnel
```

If you want to build the binary yourself, follow these steps:

### Prerequisites
- [Go 1.21+](https://go.dev/dl/) (or use the included build script)
- Windows 10/11
- Administrator privileges (for TUN adapter creation)

### Option 1: Using Go (Recommended)
```powershell
# 1. Install Go from https://go.dev/dl/
# 2. Add Go to your PATH
go version  # Should print go version >= 1.21

# 3. Clone the repo
git clone https://github.com/danishtesting4/winguard-cli.git
cd winguard-cli

# 4. Build single static binary (recommended flags)
go build -ldflags="-s -w" -o winguard.exe ./src

# 5. Optional: Compress with UPX for smaller size
# ipr 2>&1 | findstr "upx" || (iwr https://upx.github.io/upx.exe -UseBasicParsing).SaveFile("upx.exe")
& .\upx.exe --best winguard.exe

# 6. Result: ~3-8MB standalone executable!
ls -la winguard.exe  # Should be 3-8MB
```

### Option 2: Using Pre-built Binary
Check the [Releases](https://github.com/danishtesting4/winguard-cli/releases) page for ready-to-use `.exe` files.

## 🎯 Features After Building

| Feature | Description |
|---------|-------------|
| **Binary size** | 3-8MB (compressed with UPX: ~2-5MB) |
| **Memory footprint** | < 50MB active |
| **Connect time** | < 2 seconds |
| **Obfuscation modes** | wstunnel, LWO, none |
| **DPI bypass** | Bypasses country-level blocks |
| **Configuration** | YAML config file or CLI flags |

## 🚀 Usage (After Building)

```powershell
# Connect with wstunnel obfuscation (recommended for country blocks)
.\winguard connect --obfuscation wstunnel

# Connect with LWO obfuscation (ISP DPI bypass)
.\winguard connect --obfuscation lwo

# Simple connection (no obfuscation)
.\winguard connect

# Check status
.\winguard status

# Disconnect
.\winguard disconnect
```

## 📁 Repository Structure

```
winguard-cli/
├── .gitignore          # Excludes binaries and local config
├── CHANGELOG.md        # Development roadmap
├── LICENSE             # MIT License - free and open source
├── README.md           # This file
├── config/
│   └── sample-config.yaml  # Sample configuration
├── go.mod              # Go module definition
├── go.sum              # Go dependencies
└── src/
    └── main.go         # Complete Go source code
```

## 🔧 Building Instructions Summary

1. **Install Go** (or use existing installation)
2. **Clone repo**: `git clone https://github.com/danishtesting4/winguard-cli.git`
3. **Build binary**: `go build -ldflags="-s -w" -o winguard.exe ./src`
4. **Optional compression**: `upx.exe --best winguard.exe` (reduces size by ~40-60%)
5. **Run**: `.\winguard.exe connect --obfuscation wstunnel`

## 📦 Binary Size Optimization

The Go binary can be optimized:

| Technique | Expected Size |
|-----------|--------------|
| Default `go build` | ~15-25MB |
| With `-ldflags="-s -w"` | ~8-12MB |
| + UPX compression | ~3-5MB |
| **Final recommended** | **~3-8MB** ✅ |

## 🚀 One-Line Setup

```powershell
# 1. Install Go
# 2. Clone & build
go get git.zx2c4.com/wireguard-go
go build -o winguard.exe ./src && upx.exe --best winguard.exe

# 3. Run your VPN!
.\winguard connect --obfuscation wstunnel
```

## 🐋 Docker (Alternative)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o winguard .

FROM alpine:3.18
COPY --from=builder /app/winguard /usr/local/bin/
ENTRYPOINT ["winguard"]
```

## 💡 Pro Tips

- **First run**: Requires Administrator privileges for TUN adapter creation
- **Auto-start**: Create a shortcut in Startup folder for automatic connection
- **Config file**: Copy `config/sample-config.yaml` to `config.yaml` and edit
- **Logging**: Check `winlog` or run with verbose output for debugging
- **Uninstall**: Just delete the `.exe` - no system modifications on removal

## 🆘 Need Help?

- **Issues**: [GitHub Issues](https://github.com/danishtesting4/winguard-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/danishtesting4/winguard-cli/discussions)
- **Documentation**: Full docs in `README.md` and `CHANGELOG.md`

---

*WinGuard CLI - Bypass country blocks, stay private, stay lightweight.*