# WinGuard CLI - Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-11-02
### Added
- Initial release of WinGuard CLI
- WireGuard VPN protocol support
- wstunnel obfuscation mode (WebSocket HTTPS wrapping)
- LWO obfuscation mode (Lightweight header scrambling)
- CLI interface with cobra (connect, disconnect, status)
- YAML configuration parsing
- Windows TUN device management
- MIT License

## Planned
- GUI system tray icon (optional)
- Auto-start configuration
- Multiple obfuscation chaining
- Config wizard for different censorship levels
- Port forwarding support
- DNS leak protection