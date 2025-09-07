# Docker管理脚本

本目录包含Docker容器的启动、停止和管理脚本。

## 脚本说明

- **docker-start.bat** - Windows生产环境启动脚本
- **docker-start.sh** - Linux/macOS生产环境启动脚本
- **docker-dev-start.bat** - Windows开发环境启动脚本
- **docker-dev-start.sh** - Linux/macOS开发环境启动脚本
- **docker-stop.sh** - 停止所有Docker容器
- **start-services.sh** - 智能启动脚本（自动选择Docker Compose版本）
- **diagnose-containers.bat** - Windows容器诊断脚本
- **diagnose-containers.sh** - Linux/macOS容器诊断脚本
- **fix-container-restart.bat** - Windows容器重启修复脚本
- **fix-container-restart.sh** - Linux/macOS容器重启修复脚本
- **fix-config-and-logs.bat** - Windows配置和日志修复脚本
- **rebuild-with-clean-names.bat** - Windows镜像重新构建脚本
- **rebuild-with-clean-names.sh** - Linux/macOS镜像重新构建脚本
- **cleanup-and-restart.bat** - Windows清理旧容器并重启脚本
- **cleanup-and-restart.sh** - Linux/macOS清理旧容器并重启脚本
- **test-database-connection.bat** - Windows数据库连接测试脚本

## 使用流程

### 生产环境
```cmd
# Windows生产环境
docker-start.bat

# Linux环境
./docker-start.sh
```

### 开发环境
```cmd
# Windows开发环境
docker-dev-start.bat

# Linux/macOS开发环境
./docker-dev-start.sh
```

### 停止服务
```bash
./docker-stop.sh
```

### 智能启动（推荐）
自动选择可用的Docker Compose版本：

```bash
# 给脚本添加执行权限
chmod +x start-services.sh

# 运行智能启动脚本
./start-services.sh
```

### 容器诊断
当容器出现问题时，使用诊断脚本：

```cmd
# Windows
diagnose-containers.bat

# Linux/macOS
chmod +x diagnose-containers.sh
./diagnose-containers.sh
```

### 容器重启修复
当容器一直重启时，使用修复脚本：

```cmd
# Windows
fix-container-restart.bat

# Linux/macOS
chmod +x fix-container-restart.sh
./fix-container-restart.sh
```

## 环境说明

### 生产环境 (docker-compose.yml)
- MySQL: localhost:3308
- Redis: localhost:6381
- Web应用: localhost:8082
- MQTT: localhost:1883

### 开发环境 (docker-compose.dev.yml)
- MySQL: localhost:3307
- Redis: localhost:6380
- Web应用: localhost:8083
- MQTT: localhost:1884
