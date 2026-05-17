# htop PowerShell 自动加载脚本
# 将此文件内容添加到您的 PowerShell Profile 中

# 添加 htop 到 PATH
$htopPath = "D:\trae\todolist\tools\htop"
if ($env:Path -notlike "*$htopPath*") {
    $env:Path = $env:Path + ";" + $htopPath
    Write-Host "htop added to PATH: $htopPath" -ForegroundColor Green
}

# 设置 htop 别名（可选）
if (-not (Get-Alias -Name htop -ErrorAction SilentlyContinue)) {
    Set-Alias -Name htop -Value "$htopPath\htop.exe" -Scope Global
    Write-Host "htop alias created" -ForegroundColor Green
}
