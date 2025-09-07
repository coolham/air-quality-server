@echo off
chcp 65001 >nul

echo 🚀 启动空气质量监测系统开发环境...

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

REM 停止现有容器
echo 🛑 停止现有开发环境容器...
docker-compose -f docker-compose.dev.yml down

REM 启动开发环境服务
echo 🔨 启动开发环境服务...
docker-compose -f docker-compose.dev.yml up -d

REM 等待服务启动
echo ⏳ 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
echo 🔍 检查开发环境服务状态...
docker-compose -f docker-compose.dev.yml ps

echo ✅ 开发环境启动完成！
echo 🗄️  MySQL: localhost:3307 (用户名: root, 密码: admin)
echo 🔴  Redis: localhost:6380
echo 📡  MQTT: localhost:1884 (可选)
echo.
echo 📝 常用命令:
echo   查看MySQL日志: docker-compose -f docker-compose.dev.yml logs -f mysql-dev
echo   查看Redis日志: docker-compose -f docker-compose.dev.yml logs -f redis-dev
echo   停止开发环境: docker-compose -f docker-compose.dev.yml down
echo   查看状态: docker-compose -f docker-compose.dev.yml ps
echo.
echo 💡 提示: 开发环境只启动数据库和缓存服务，应用服务请在本地运行

pause
