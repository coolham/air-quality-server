@echo off
chcp 65001 >nul
echo ========================================
echo   启动空气质量监测服务
echo   (包含内置MQTT服务器)
echo ========================================
echo.

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 未找到Go，请先安装Go
    pause
    exit /b 1
)

echo ✅ Go已安装
echo.

REM 设置环境变量
set AIR_QUALITY_CONFIG=config\config.yaml
set SERVER_HOST=0.0.0.0
set SERVER_PORT=8080
set LOG_LEVEL=info

echo 📋 启动配置:
echo   - HTTP服务器: http://localhost:8080
echo   - MQTT服务器: tcp://localhost:1883
echo   - 用户名: admin
echo   - 密码: password
echo.

echo 🚀 正在启动服务...
echo 按 Ctrl+C 停止服务
echo ----------------------------------------

REM 启动Go应用程序
cd /d "%~dp0.."
go run cmd/air-quality-server/main.go

echo.
echo ✅ 服务已停止
pause
