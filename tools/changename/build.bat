@echo off
echo ========================================
echo  Smart Renamer - 构建脚本
echo ========================================
echo.

echo [1/3] 安装前端依赖...
cd frontend
call npm install
if errorlevel 1 (
    echo 前端依赖安装失败！
    pause
    exit /b 1
)
echo 前端依赖安装成功！
echo.

echo [2/3] 构建前端...
call npm run build
if errorlevel 1 (
    echo 前端构建失败！
    pause
    exit /b 1
)
echo 前端构建成功！
cd ..
echo.

echo [3/3] 编译 Go 程序...
go build -o build/smart-renamer.exe .
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
