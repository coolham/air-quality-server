@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - 构建应用程序
echo ========================================
echo.

echo 正在构建应用程序...
echo.

if not exist bin mkdir bin

go build -o bin\air-quality-server.exe .\cmd\air-quality-server

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo ✅ 构建成功！
    echo ========================================
    echo.
    echo 可执行文件位置: bin\air-quality-server.exe
    echo.
    echo 运行方式：
    echo 1. 直接运行: bin\air-quality-server.exe
    echo 2. 或使用: run-app.bat
    echo.
) else (
    echo.
    echo ========================================
    echo ❌ 构建失败！
    echo ========================================
    echo.
    echo 请检查：
    echo 1. Go环境是否正确安装
    echo 2. 项目依赖是否完整
    echo 3. 代码是否有编译错误
    echo.
)

pause
