@echo off
echo 配置Windows防火墙规则以允许MQTT服务器...
echo.

REM 检查是否以管理员权限运行
net session >nul 2>&1
if %errorLevel% == 0 (
    echo 检测到管理员权限，继续配置...
) else (
    echo 错误：需要管理员权限来配置防火墙规则
    echo 请右键点击此脚本，选择"以管理员身份运行"
    pause
    exit /b 1
)

echo.
echo 正在添加防火墙规则...

REM 添加入站规则 - 允许MQTT服务器端口1883
netsh advfirewall firewall add rule name="MQTT Server - Air Quality (Inbound)" dir=in action=allow protocol=TCP localport=1883 profile=any
if %errorLevel% == 0 (
    echo ✅ 入站规则添加成功
) else (
    echo ❌ 入站规则添加失败
)

REM 添加出站规则 - 允许MQTT服务器端口1883
netsh advfirewall firewall add rule name="MQTT Server - Air Quality (Outbound)" dir=out action=allow protocol=TCP localport=1883 profile=any
if %errorLevel% == 0 (
    echo ✅ 出站规则添加成功
) else (
    echo ❌ 出站规则添加失败
)

REM 添加Web服务器端口8080的规则
netsh advfirewall firewall add rule name="Web Server - Air Quality (Inbound)" dir=in action=allow protocol=TCP localport=8080 profile=any
if %errorLevel% == 0 (
    echo ✅ Web服务器入站规则添加成功
) else (
    echo ❌ Web服务器入站规则添加失败
)

echo.
echo 🎉 防火墙配置完成！
echo.
echo 已添加的规则：
echo - MQTT Server - Air Quality (Inbound) - 端口1883
echo - MQTT Server - Air Quality (Outbound) - 端口1883  
echo - Web Server - Air Quality (Inbound) - 端口8080
echo.
echo 现在运行程序时不会再弹出防火墙确认对话框了。
echo.
pause
