@echo off
echo 删除Windows防火墙规则...
echo.

REM 检查是否以管理员权限运行
net session >nul 2>&1
if %errorLevel% == 0 (
    echo 检测到管理员权限，继续删除规则...
) else (
    echo 错误：需要管理员权限来删除防火墙规则
    echo 请右键点击此脚本，选择"以管理员身份运行"
    pause
    exit /b 1
)

echo.
echo 正在删除防火墙规则...

REM 删除MQTT服务器入站规则
netsh advfirewall firewall delete rule name="MQTT Server - Air Quality (Inbound)"
if %errorLevel% == 0 (
    echo ✅ MQTT入站规则删除成功
) else (
    echo ⚠️ MQTT入站规则可能不存在或已删除
)

REM 删除MQTT服务器出站规则
netsh advfirewall firewall delete rule name="MQTT Server - Air Quality (Outbound)"
if %errorLevel% == 0 (
    echo ✅ MQTT出站规则删除成功
) else (
    echo ⚠️ MQTT出站规则可能不存在或已删除
)

REM 删除Web服务器规则
netsh advfirewall firewall delete rule name="Web Server - Air Quality (Inbound)"
if %errorLevel% == 0 (
    echo ✅ Web服务器规则删除成功
) else (
    echo ⚠️ Web服务器规则可能不存在或已删除
)

echo.
echo 🎉 防火墙规则清理完成！
echo.
pause
