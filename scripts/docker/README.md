# Docker管理脚本

本目录包含Docker容器的启动、停止和管理脚本。

## 脚本说明

- **docker-start.bat** - 启动生产环境Docker服务
- **docker-start.sh** - Linux环境启动脚本
- **docker-dev-start.bat** - 启动开发环境Docker服务
- **docker-dev-start.sh** - Linux开发环境启动脚本
- **docker-stop.sh** - 停止所有Docker容器

## 使用流程

### 启动服务
```cmd
# Windows生产环境
docker-start.bat

# Windows开发环境
docker-dev-start.bat

# Linux环境
./docker-start.sh
```

### 停止服务
```bash
./docker-stop.sh
```
