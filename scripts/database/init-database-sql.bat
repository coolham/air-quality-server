@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - SQL脚本数据库初始化
echo ========================================
echo.

echo 此脚本将使用SQL脚本初始化数据库
echo 请确保MySQL服务已启动，并且root用户密码正确
echo.

set /p mysql_root_password="请输入MySQL root用户密码: "

echo.
echo 正在执行SQL脚本初始化数据库...
echo.

mysql -u root -p%mysql_root_password% < scripts\init.sql

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo ✅ 数据库初始化成功！
    echo ========================================
    echo.
    echo 已创建以下内容：
    echo - 数据库: air_quality
    echo - 所有表结构和索引
    echo - 默认用户和角色
    echo - 示例设备和配置
    echo - 告警规则
    echo.
    echo 下一步操作：
    echo 1. 运行 check-database.bat 检查数据库状态
    echo 2. 运行 start-app.bat 启动应用程序
    echo.
) else (
    echo.
    echo ========================================
    echo ❌ 数据库初始化失败！
    echo ========================================
    echo.
    echo 请检查：
    echo 1. MySQL服务是否运行
    echo 2. root用户密码是否正确
    echo 3. 是否有足够的权限创建数据库
    echo 4. scripts\init.sql 文件是否存在
    echo.
)

pause
