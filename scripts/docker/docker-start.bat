@echo off
chcp 65001 >nul

echo 🚀 启动空气质量监测系统...

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker未运行，请先启动Docker Desktop
    pause
    exit /b 1
)

REM 检查docker-compose是否安装
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ docker-compose未安装，请先安装docker-compose
    pause
    exit /b 1
)

REM 创建必要的目录
echo 📁 创建必要的目录...
if not exist "logs" mkdir logs
if not exist "config" mkdir config

REM 检查配置文件
if not exist "config\config.docker.yaml" (
    echo ❌ 配置文件不存在: config\config.docker.yaml
    pause
    exit /b 1
)

REM 停止现有容器
echo 🛑 停止现有容器...
docker-compose down

REM 构建并启动服务
echo 🔨 构建并启动服务...
docker-compose up --build -d

REM 等待服务启动
echo ⏳ 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
echo 🔍 检查服务状态...
docker-compose ps

REM 显示日志
echo 📋 显示服务日志...
docker-compose logs --tail=50 air-quality-server

echo ✅ 空气质量监测系统启动完成！
echo 🌐 Web界面: http://localhost:8080
echo 📊 Dashboard: http://localhost:8080/dashboard
echo 📡 MQTT Broker: localhost:1883
echo.
echo 📝 常用命令:
echo   查看日志: docker-compose logs -f air-quality-server
echo   停止服务: docker-compose down
echo   重启服务: docker-compose restart air-quality-server
echo   查看状态: docker-compose ps

pause
