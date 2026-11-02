<#
.SYNOPSIS
WinGuard CLI - One-Command Install for Windows

.DESCRIPTION
Installs WinGuard CLI VPN with a single PowerShell command.
Supports downloading pre-built binary or building from source.

.PARAMETER Method
Installation method: "binary" (default) or "build"

.EXAMPLE
.\install.ps1 -Method binary
.\install.ps1 -Method build

.INPUTS
None

.OUTPUTS
System.Object

.NOTES
Requires PowerShell 5.1+ and Administrator privileges for TUN adapter creation.

Copyright (c) 2026 danis. MIT License.
#>

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("binary", "build")]
    [string]$Method = "binary"
)

Write-Host "╔═════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   🛡️  WinGuard CLI - One-Command Installer for Windows           ║" -ForegroundColor Cyan
Write-Host "╚═════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan

function Check-Admin {
    $principal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    return $principal.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)
}

if (!(Check-Admin)) {
    Write-Host "⚠️  Running without Administrator privileges!" -ForegroundColor Yellow
    Write-Host "   Some features may require re-running as Admin (for TUN adapter)" -ForegroundColor Yellow
}

Write-Host "▶️  Selecting installation method: $Method" -ForegroundColor Green

switch ($Method) {
    "binary" {
        Write-Host "📥  Downloading pre-built WinGuard CLI binary..." -ForegroundColor White
        
        $url = "https://github.com/danishtesting4/winguard-cli/releases/latest/download/winguard-windows-amd64.exe"
        $exePath = "$env:USERPROFILE\Desktop\winguard.exe"
        
        try {
            Write-Host "   Fetching latest release..." -ForegroundColor Gray
            Invoke-WebRequest -Uri $url -OutFile $exePath -UseBasicParsing
            
            Write-Host "✅  Installation complete!" -ForegroundColor Green
            Write-Host "   Binary saved to: $exePath" -ForegroundColor Green
            Write-Host "" -ForegroundColor White
            Write-Host "▶️  Next steps:" -ForegroundColor Cyan
            Write-Host "   1. Right-click winguard.exe → Run as administrator" -ForegroundColor Gray
            Write-Host "   2. winguard connect --obfuscation wstunnel" -ForegroundColor Gray
            Write-Host "   3. Visit: https://github.com/danishtesting4/winguard-cli#readme" -ForegroundColor Gray
        }
        catch {
            Write-Error "❌  Failed to download binary. Check your internet connection."
            Write-Error "   Alternative: Use 'build' method or visit releases manually."
        }
    }
    
    "build" {
        Write-Host "🔨  Building WinGuard CLI from source..." -ForegroundColor White
        Write-Host "   This requires Go 1.21+ installed on your system." -ForegroundColor Yellow
        
        if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
            Write-Error "❌  Go is not installed. Please install Go 1.21+ first:"
            Write-Error "   https://go.dev/dl/"
            return
        }
        
        Write-Host "   Cloning repository..." -ForegroundColor Gray
        $cloneDir = "$env:TEMP\winguard-cli"
        if (Test-Path $cloneDir) {
            Remove-Item $cloneDir -Recurse -Force
        }
        
        try {
            git clone "https://github.com/danishtesting4/winguard-cli.git" $cloneDir
            Set-Location $cloneDir
            
            Write-Host "   Building optimized binary..." -ForegroundColor Gray
            $result = & go build -ldflags="-s -w" -o winguard.exe ./src
            
            Write-Host "   Compressing with UPX..." -ForegroundColor Gray
            if (Get-Command upx.exe -ErrorAction SilentlyContinue) {
                & upx.exe --best winguard.exe
                Write-Host "   ✅  UPX compression complete" -ForegroundColor Green
            }
            else {
                Write-Host "   ⚠️  UPX not found - binary size will be larger" -ForegroundColor Yellow
            }
            
            $exePath = Join-Path $cloneDir "winguard.exe"
            Write-Host "✅  Build complete!" -ForegroundColor Green
            Write-Host "   Binary location: $exePath" -ForegroundColor Green
            Write-Host "" -ForegroundColor White
            Write-Host "▶️  Next steps:" -ForegroundColor Cyan
            Write-Host "   1. Run: `& .\winguard.exe connect --obfuscation wstunnel`" -ForegroundColor Gray
            Write-Host "   2. Requires Administrator privileges for first run" -ForegroundColor Gray
        }
        catch {
            Write-Error "❌  Build failed. Errors:"
            $_.Exception | Format-Message | ForEach-Object { Write-Error "   $_" }
        }
    }
    
    default {
        Write-Error "❌  Invalid method: $Method. Use 'binary' or 'build'."
    }
}

Write-Host "" -ForegroundColor Cyan
Write-Host "🎉  WinGuard CLI is ready to bypass country-level DPI blocks!" -ForegroundColor Green
Write-Host "📖  Documentation: https://github.com/danishtesting4/winguard-cli" -ForegroundColor Cyan