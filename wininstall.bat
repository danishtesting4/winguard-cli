@echo off
title WinGuard CLI Installer
color 0A

echo.
echo ============================================
echo  🛡️  WinGuard CLI - One-Click Installer
echo ============================================
echo.

:: Check for Administrator
net session >nul 2>&1
if %errorLevel% == 0 (
    echo ✅  Administrator privileges detected
) else (
    echo ⚠️  Not running as Administrator
    echo   Some features may limited
    pause
    exit /b
)

:: Download directory
set INSTALL_DIR=%USERPROFILE%\WinGuard
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo ▶️  Downloading WinGuard CLI binary...
curl -L -o "%INSTALL_DIR%\winguard.exe" "https://github.com/danishtesting4/winguard-cli/releases/latest/download/winguard-windows-amd64.exe"

echo ▶️  Installing...
echo.
echo   Binary placed at: %INSTALL_DIR%\winguard.exe
echo.
echo ✅  WinGuard CLI installed successfully!
echo.

:: Add to PATH (optional)
set PATH=%INSTALL_DIR%;%PATH%
setx PATH "%PATH%" /M

echo ▶️  Quick test...
echo.
echo   Running: winguard status
winguard status

echo.
echo ============================================
echo  Installation complete!
echo  Run: winguard connect --obfuscation wstunnel
echo ============================================
pause