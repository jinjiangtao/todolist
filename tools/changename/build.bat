@echo off
echo ========================================
echo  Smart Renamer - 构建脚本
echo ========================================
echo.

echo [1/1] 编译 Go 程序...
go mod tidy
if errorlevel 1 (
    echo 依赖安装失败！
    pause
    exit /b 1
)

go build -ldflags "-s -w" -o build\smart-renamer.exe
if errorlevel 1 (
    echo 编译失败！
    pause
    exit /b 1
)

echo.
echo ========================================
echo  构建成功！
echo  可执行文件位置: build\smart-renamer.exe
echo ========================================
echo.
pause
