# Trae IDE 终端初始化脚本

# 确保项目目录在 PATH 中，以便可以直接运行 awk
$projectPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$env:PATH = "$projectPath;$env:PATH"

# 显示欢迎信息
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  GoAWK 开发环境已就绪" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "可以直接使用以下命令：" -ForegroundColor Gray
Write-Host "  awk '$1' test.txt      - 打印第一列" -ForegroundColor White
Write-Host "  go test ./...         - 运行所有测试" -ForegroundColor White
Write-Host "  go build -o awk.exe ./cmd/awk  - 编译程序" -ForegroundColor White
Write-Host ""

# 切换到项目目录
Set-Location $projectPath
