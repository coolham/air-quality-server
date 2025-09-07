# 检查Windows防火墙规则
Write-Host "检查空气质量监测系统的防火墙规则..." -ForegroundColor Green
Write-Host ""

# 检查MQTT服务器规则
Write-Host "🔍 检查MQTT服务器规则 (端口1883):" -ForegroundColor Yellow
$mqttInbound = Get-NetFirewallRule -DisplayName "*MQTT Server - Air Quality (Inbound)*" -ErrorAction SilentlyContinue
$mqttOutbound = Get-NetFirewallRule -DisplayName "*MQTT Server - Air Quality (Outbound)*" -ErrorAction SilentlyContinue

if ($mqttInbound) {
    Write-Host "✅ MQTT入站规则存在" -ForegroundColor Green
    Write-Host "   规则名称: $($mqttInbound.DisplayName)" -ForegroundColor Gray
    Write-Host "   状态: $($mqttInbound.Enabled)" -ForegroundColor Gray
} else {
    Write-Host "❌ MQTT入站规则不存在" -ForegroundColor Red
}

if ($mqttOutbound) {
    Write-Host "✅ MQTT出站规则存在" -ForegroundColor Green
    Write-Host "   规则名称: $($mqttOutbound.DisplayName)" -ForegroundColor Gray
    Write-Host "   状态: $($mqttOutbound.Enabled)" -ForegroundColor Gray
} else {
    Write-Host "❌ MQTT出站规则不存在" -ForegroundColor Red
}

Write-Host ""

# 检查Web服务器规则
Write-Host "🔍 检查Web服务器规则 (端口8080):" -ForegroundColor Yellow
$webInbound = Get-NetFirewallRule -DisplayName "*Web Server - Air Quality (Inbound)*" -ErrorAction SilentlyContinue

if ($webInbound) {
    Write-Host "✅ Web服务器入站规则存在" -ForegroundColor Green
    Write-Host "   规则名称: $($webInbound.DisplayName)" -ForegroundColor Gray
    Write-Host "   状态: $($webInbound.Enabled)" -ForegroundColor Gray
} else {
    Write-Host "❌ Web服务器入站规则不存在" -ForegroundColor Red
}

Write-Host ""

# 检查端口监听状态
Write-Host "🔍 检查端口监听状态:" -ForegroundColor Yellow
$port1883 = Get-NetTCPConnection -LocalPort 1883 -ErrorAction SilentlyContinue
$port8080 = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue

if ($port1883) {
    Write-Host "✅ 端口1883正在监听" -ForegroundColor Green
    Write-Host "   状态: $($port1883.State)" -ForegroundColor Gray
} else {
    Write-Host "❌ 端口1883未监听" -ForegroundColor Red
}

if ($port8080) {
    Write-Host "✅ 端口8080正在监听" -ForegroundColor Green
    Write-Host "   状态: $($port8080.State)" -ForegroundColor Gray
} else {
    Write-Host "❌ 端口8080未监听" -ForegroundColor Red
}

Write-Host ""
Write-Host "检查完成！" -ForegroundColor Green
