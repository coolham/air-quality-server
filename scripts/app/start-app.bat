@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - 启动应用程序
echo ========================================
echo.

echo 正在启动应用程序...
echo.
echo 应用程序将在以下地址启动：
echo - HTTP服务: http://localhost:8080
echo - 健康检查: http://localhost:8080/health
echo - 仪表板: http://localhost:8080/dashboard
echo.
echo 按 Ctrl+C 停止应用程序
echo.

go run cmd\air-quality-server\main.go

echo.
echo 应用程序已停止
pause
