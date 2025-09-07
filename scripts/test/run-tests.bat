@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - 运行单元测试
echo ========================================
echo.

echo 正在运行Go单元测试...
echo.

cd /d "%~dp0.."

echo 运行MQTT模块测试...
go test -v ./internal/mqtt/...

echo.
echo 运行所有测试...
go test -v ./...

echo.
echo ========================================
echo 测试完成
echo ========================================
echo.

pause
