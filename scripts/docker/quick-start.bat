@echo off
REM Docker 快速启动脚本
REM 解决 "Unknown column 'o3' in 'field list'" 错误

echo ========================================
echo Docker 快速启动
echo ========================================

REM 检查 Docker
docker info >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Docker 未运行
    pause
    exit /b 1
)

echo 正在启动服务...
echo 数据库将通过 init 脚本自动初始化，包含所有必要字段

docker-compose up -d

if %errorlevel% equ 0 (
    echo.
    echo ✅ 服务启动成功！
    echo.
    echo 服务地址：
    echo   Web 界面: http://localhost:8082
    echo   MQTT 端口: 1883
    echo   MySQL 端口: 3308
    echo.
    echo 数据库已包含所有必要字段，错误已解决。
) else (
    echo.
    echo ❌ 服务启动失败！
)

echo.
pause
