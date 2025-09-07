# 防火墙管理脚本

本目录包含Windows防火墙规则的配置和管理脚本。

## 脚本说明

- **configure_firewall.bat** - 配置防火墙规则，允许应用程序和MQTT端口通过
- **remove_firewall_rules.bat** - 移除之前添加的防火墙规则
- **check_firewall_rules.ps1** - 检查现有的防火墙规则
- **start_with_firewall.bat** - 自动配置防火墙并启动应用程序

## 使用流程

### 配置防火墙
```cmd
# 需要管理员权限
configure_firewall.bat
```

### 检查防火墙规则
```cmd
powershell -ExecutionPolicy Bypass -File check_firewall_rules.ps1
```

### 移除防火墙规则
```cmd
# 需要管理员权限
remove_firewall_rules.bat
```

### 带防火墙启动
```cmd
# 需要管理员权限
start_with_firewall.bat
```
