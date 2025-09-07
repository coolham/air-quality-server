@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - 一键安装脚本
echo ========================================
echo.

echo 此脚本将帮助您完成以下操作：
echo 1. 检查Go环境
echo 2. 下载项目依赖
echo 3. 初始化数据库
echo 4. 构建应用程序
echo.

set /p confirm="是否继续？(Y/N): "
if /i not "%confirm%"=="Y" (
    echo 操作已取消
    pause
    exit /b 0
)

echo.
echo ========================================
echo 步骤1: 检查Go环境
echo ========================================
go version
if %errorlevel% neq 0 (
    echo ❌ Go环境未正确安装！
    echo 请先安装Go 1.21或更高版本
    pause
    exit /b 1
)
echo ✅ Go环境检查通过

echo.
echo ========================================
echo 步骤2: 下载项目依赖
echo ========================================
go mod download
if %errorlevel% neq 0 (
    echo ❌ 依赖下载失败！
    pause
    exit /b 1
)
echo ✅ 依赖下载完成

echo.
echo ========================================
echo 步骤3: 初始化数据库
echo ========================================
echo 请选择数据库初始化方式：
echo 1. 使用SQL脚本（推荐，更快速）
echo 2. 使用Go迁移工具
echo 3. 跳过数据库初始化
echo.
set /p db_choice="请选择 (1/2/3): "

if "%db_choice%"=="1" (
    echo 使用SQL脚本初始化数据库...
    call init-database-sql.bat
    if %errorlevel% neq 0 (
        echo ❌ SQL脚本数据库初始化失败！
        goto :build
    )
) else if "%db_choice%"=="2" (
    echo 使用Go迁移工具初始化数据库...
    go run cmd\migrate\main.go -action init
    if %errorlevel% neq 0 (
        echo ❌ Go迁移工具数据库初始化失败！
        echo 请检查数据库配置和连接
        goto :build
    )
) else if "%db_choice%"=="3" (
    echo 跳过数据库初始化
    goto :build
) else (
    echo 无效选择，跳过数据库初始化
    goto :build
)
echo ✅ 数据库初始化完成

:build
echo.
echo ========================================
echo 步骤4: 构建应用程序
echo ========================================
if not exist bin mkdir bin
go build -o bin\air-quality-server.exe .\cmd\air-quality-server
if %errorlevel% neq 0 (
    echo ❌ 应用程序构建失败！
    pause
    exit /b 1
)
echo ✅ 应用程序构建完成

echo.
echo ========================================
echo 🎉 安装完成！
echo ========================================
echo.
echo 下一步操作：
echo 1. 如果数据库未初始化，请运行: init-database-sql.bat 或 init-database.bat
echo 2. 启动应用程序: run-app.bat
echo 3. 或开发模式运行: start-app.bat
echo.
echo 访问地址：
echo - 健康检查: http://localhost:8080/health
echo - 仪表板: http://localhost:8080/dashboard
echo.

pause
