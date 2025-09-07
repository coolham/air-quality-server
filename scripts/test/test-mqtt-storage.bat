@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - MQTT数据存储测试
echo ========================================
echo.

echo 此脚本将测试MQTT服务器的数据存储功能
echo 请确保：
echo 1. 应用程序正在运行
echo 2. MQTT服务器已启动
echo 3. 数据库已初始化
echo.

set /p confirm="是否继续测试？(Y/N): "
if /i not "%confirm%"=="Y" (
    echo 测试已取消
    pause
    exit /b 0
)

echo.
echo 正在运行MQTT数据存储测试...
echo.

cd /d "%~dp0.."
python tools\mqtt\quick_test.py

echo.
echo ========================================
echo 测试完成
echo ========================================
echo.

pause
