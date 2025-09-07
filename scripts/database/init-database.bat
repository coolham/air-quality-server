@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - 数据库初始化工具
echo ========================================
echo.

echo 正在初始化数据库...
echo.

go run cmd\migrate\main.go -action init

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo ✅ 数据库初始化成功！
    echo ========================================
    echo.
    echo 下一步操作：
    echo 1. 运行 check-database.bat 检查数据库状态
    echo 2. 运行 start-app.bat 启动应用程序
    echo.
) else (
    echo.
    echo ========================================
    echo ❌ 数据库初始化失败！
    echo ========================================
    echo.
    echo 请检查：
    echo 1. MySQL服务是否运行
    echo 2. 数据库连接配置是否正确
    echo 3. 用户权限是否足够
    echo.
)

pause
