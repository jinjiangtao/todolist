Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  htop PowerShell Profile Config Tool" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$profilePath = $PROFILE
$profileDir = Split-Path -Parent $profilePath

Write-Host "Checking PowerShell Profile..." -ForegroundColor Yellow
Write-Host "Profile path: $profilePath" -ForegroundColor Gray

if (-not (Test-Path $profileDir)) {
    Write-Host "Creating Profile directory..." -ForegroundColor Cyan
    New-Item -ItemType Directory -Force -Path $profileDir | Out-Null
}

$alreadyConfigured = $false
if (Test-Path $profilePath) {
    $profileContent = Get-Content $profilePath -Raw -ErrorAction SilentlyContinue
    if ($profileContent -like "*htop*") {
        Write-Host "htop configuration already exists in Profile" -ForegroundColor Green
        $alreadyConfigured = $true
    }
}

if (-not $alreadyConfigured) {
    Write-Host "Adding htop configuration to Profile..." -ForegroundColor Cyan
    
    $htopConfig = @"

# htop auto-load configuration
`$htopPath = "D:\trae\todolist\tools\htop"
if (`$env:Path -notlike "*`$htopPath*") {
    `$env:Path = `$env:Path + ";" + `$htopPath
}
"@
    
    Add-Content -Path $profilePath -Value $htopConfig
    
    Write-Host "htop added to PowerShell Profile" -ForegroundColor Green
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Configuration Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Close current PowerShell window" -ForegroundColor White
Write-Host "2. Reopen PowerShell" -ForegroundColor White
Write-Host "3. Type 'htop' anywhere to run the program" -ForegroundColor White
Write-Host ""
Write-Host "Note: Reopened PowerShell will automatically load htop" -ForegroundColor Cyan
