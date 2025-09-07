@echo off
REM 空气质量监测系统 - 中国大陆Docker构建脚本 (Windows版本)
REM 解决Go模块下载超时问题

echo 🚀 为中国大陆网络环境构建Docker镜像...

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker未运行，请先启动Docker
    exit /b 1
)

REM 检查端口占用情况
echo 🔍 检查端口占用情况...
netstat -an | findstr :3308 >nul
if not errorlevel 1 (
    echo ⚠️  端口3308已被占用，可能影响MySQL服务
)

netstat -an | findstr :6381 >nul
if not errorlevel 1 (
    echo ⚠️  端口6381已被占用，可能影响Redis服务
)

netstat -an | findstr :8082 >nul
if not errorlevel 1 (
    echo ⚠️  端口8082已被占用，可能影响Web服务
)

REM 设置Go代理环境变量
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=sum.golang.google.cn

echo 🔧 设置Go代理环境变量:
echo    GOPROXY=%GOPROXY%
echo    GOSUMDB=%GOSUMDB%

REM 检查构建类型
set BUILD_TYPE=%1
if "%BUILD_TYPE%"=="" set BUILD_TYPE=production

if "%BUILD_TYPE%"=="dev" (
    echo 🔨 构建开发环境镜像...
    docker build --build-arg GOPROXY=%GOPROXY% --build-arg GOSUMDB=%GOSUMDB% -f Dockerfile.dev -t air-quality-server:dev .
) else (
    echo 🔨 构建生产环境镜像...
    docker build --build-arg GOPROXY=%GOPROXY% --build-arg GOSUMDB=%GOSUMDB% -f Dockerfile -t air-quality-server:latest .
)

if errorlevel 1 (
    echo ❌ Docker镜像构建失败
    exit /b 1
)

echo ✅ Docker镜像构建完成！

REM 显示镜像信息
echo 📋 构建的镜像:
docker images | findstr air-quality-server

echo.
echo 📝 使用说明:
if "%BUILD_TYPE%"=="dev" (
    echo   启动开发环境: docker-compose -f docker-compose.dev.yml up -d
) else (
    echo   启动生产环境: docker-compose up -d
)
echo   查看镜像: docker images ^| findstr air-quality-server
echo   删除镜像: docker rmi air-quality-server:latest
