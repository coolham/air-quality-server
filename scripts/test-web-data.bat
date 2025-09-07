@echo off
chcp 65001 >nul
echo ========================================
echo 空气质量监测系统 - Web数据查看功能测试
echo ========================================
echo.

echo 正在启动应用程序...
echo.

cd /d "%~dp0.."

echo 启动服务器...
start "Air Quality Server" bin\air-quality-server.exe

echo 等待服务器启动...
timeout /t 5 /nobreak >nul

echo.
echo ========================================
echo Web数据查看功能测试
echo ========================================
echo.

echo 1. 打开浏览器访问数据查看页面:
echo    http://127.0.0.1:8080/data
echo.

echo 2. 测试API接口:
echo    - 数据查询API: http://127.0.0.1:8080/web/api/data
echo    - 数据导出API: http://127.0.0.1:8080/web/api/data/export?format=csv
echo.

echo 3. 功能特性:
echo    ✓ 设备筛选
echo    ✓ 设备类型筛选
echo    ✓ 传感器ID筛选
echo    ✓ 时间范围筛选
echo    ✓ 分页显示
echo    ✓ 数据导出(CSV/JSON)
echo    ✓ 实时数据展示
echo.

echo 4. 测试步骤:
echo    a) 在浏览器中打开 http://127.0.0.1:8080/data
echo    b) 选择设备ID进行筛选
echo    c) 设置时间范围查询历史数据
echo    d) 测试数据导出功能
echo    e) 验证分页功能
echo.

echo 按任意键打开浏览器...
pause >nul

echo 正在打开浏览器...
start http://127.0.0.1:8080/data

echo.
echo ========================================
echo 测试完成
echo ========================================
echo.
echo 提示: 服务器正在后台运行，按Ctrl+C停止服务器
echo.

pause
