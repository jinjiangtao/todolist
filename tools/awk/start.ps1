# GoAWK 快捷启动脚本

# 添加项目路径到 PATH
$projectPath = $PSScriptRoot
if (-not ($env:PATH -split ';' -contains $projectPath)) {
    $env:PATH = "$projectPath;$env:PATH"
    Write-Host "已将 $projectPath 添加到 PATH" -ForegroundColor Green
} else {
    Write-Host "$projectPath 已在 PATH 中" -ForegroundColor Gray
}

# 显示快速帮助
Write-Host ""
Write-Host "GoAWK 快捷命令:" -ForegroundColor Cyan
Write-Host "  awk '$1' test.txt      - 打印第一列"
Write-Host "  go test ./...         - 运行所有测试"
Write-Host "  go build -o awk.exe ./cmd/awk  - 编译"
Write-Host ""

# 保持会话打开
if ($args.Count -eq 0) {
    Write-Host "保持会话打开，按 Ctrl+C 退出" -ForegroundColor Gray
    # 进入交互模式
    $host.EnterNestedPrompt()
}
