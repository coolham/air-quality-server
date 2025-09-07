@echo off
echo 空气质量监测系统 - 开发环境启动脚本
echo =====================================
echo.

REM 检查是否以管理员权限运行
net session >nul 2>&1
if %errorLevel% == 0 (
    echo ✅ 检测到管理员权限
    goto :configure_firewall
) else (
    echo ⚠️ 未检测到管理员权限，跳过防火墙配置
    echo 注意：可能会弹出防火墙确认对话框
    goto :start_server
)

:configure_firewall
echo.
echo 🔧 配置防火墙规则...
netsh advfirewall firewall add rule name="MQTT Server - Air Quality (Inbound)" dir=in action=allow protocol=TCP localport=1883 profile=any >nul 2>&1
netsh advfirewall firewall add rule name="Web Server - Air Quality (Inbound)" dir=in action=allow protocol=TCP localport=8080 profile=any >nul 2>&1
echo ✅ 防火墙规则配置完成

:start_server
echo.
echo 🚀 启动空气质量监测服务器...
echo.

REM 切换到项目根目录
cd /d "%~dp0.."

REM 启动服务器
go run cmd/air-quality-server/main.go

echo.
echo 服务器已停止
pause
