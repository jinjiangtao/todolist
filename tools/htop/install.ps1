$htopPath = "D:\trae\todolist\tools\htop"
$envPath = [Environment]::GetEnvironmentVariable("Path", "User")

Write-Host "Checking htop installation..." -ForegroundColor Cyan
Write-Host "htop path: $htopPath" -ForegroundColor Yellow

if ($envPath -like "*$htopPath*") {
    Write-Host "htop directory already exists in PATH" -ForegroundColor Green
} else {
    Write-Host "Adding htop to PATH environment variable..." -ForegroundColor Cyan
    
    if ($envPath -eq $null -or $envPath -eq "") {
        $newPath = $htopPath
    } else {
        $newPath = $envPath + ";" + $htopPath
    }
    
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "htop added to PATH environment variable" -ForegroundColor Green
    Write-Host "New PATH: $newPath" -ForegroundColor Yellow
}

$currentSessionPath = $env:Path
if ($currentSessionPath -like "*$htopPath*") {
    Write-Host "Current session already contains htop path" -ForegroundColor Green
} else {
    Write-Host "Adding to current session..." -ForegroundColor Cyan
    $env:Path = $currentSessionPath + ";" + $htopPath
    Write-Host "Added to current session" -ForegroundColor Green
}

Write-Host "`nVerifying installation:" -ForegroundColor Cyan
Write-Host "Note: If htop is not recognized in current terminal, please close and reopen terminal" -ForegroundColor Yellow
Write-Host "`nYou can try these commands:" -ForegroundColor Green
Write-Host "  1. htop" -ForegroundColor White
Write-Host "  2. D:\trae\todolist\tools\htop\htop.exe" -ForegroundColor White
Write-Host "`nInstallation complete!" -ForegroundColor Green
