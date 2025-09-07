@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - 数据库状态检查
echo ========================================
echo.

echo 正在检查数据库状态...
echo.

go run cmd\migrate\main.go -action status

echo.
echo ========================================
echo 检查完成
echo ========================================
echo.

pause
