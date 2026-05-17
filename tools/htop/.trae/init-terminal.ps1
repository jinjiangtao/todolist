# Trae IDE terminal initialization script

# Ensure project directory is in PATH for running htop directly
$projectPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$env:PATH = "$projectPath;$env:PATH"

# Display welcome message
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Windows HTOP Development Environment Ready" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Available commands:" -ForegroundColor Gray
Write-Host "  htop                  - Run process monitor" -ForegroundColor White
Write-Host "  go test ./...         - Run all tests" -ForegroundColor White
Write-Host "  go build -o htop.exe . - Build program" -ForegroundColor White
Write-Host ""
Write-Host "Shortcuts:" -ForegroundColor Gray
Write-Host "  ↑ / ↓    - Scroll process list" -ForegroundColor White
Write-Host "  P        - Sort by CPU" -ForegroundColor White
Write-Host "  M        - Sort by Memory" -ForegroundColor White
Write-Host "  L/F1     - Show help" -ForegroundColor White
Write-Host "  Q        - Quit" -ForegroundColor White
Write-Host ""

# Change to project directory
Set-Location $projectPath