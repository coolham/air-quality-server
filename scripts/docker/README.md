# Docker管理脚本

本目录包含Docker容器的启动、停止和管理脚本。

## 脚本说明

- **docker-start.bat** - 启动生产环境Docker服务
- **docker-start.sh** - Linux环境启动脚本
- **docker-dev-start.bat** - 启动开发环境数据库和缓存服务
- **docker-dev-start.sh** - Linux开发环境启动脚本
- **docker-dev-full.bat** - 启动完整开发环境（包含应用服务）
- **docker-stop.sh** - 停止所有Docker容器
- **docker-build-china.bat** - 中国大陆网络环境Docker构建脚本（Windows）
- **docker-build-china.sh** - 中国大陆网络环境Docker构建脚本（Linux/macOS）

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
# 只启动数据库和缓存服务（推荐）
docker-dev-start.bat

# 启动完整开发环境（包含应用服务）
docker-dev-full.bat

# Linux开发环境
./docker-dev-start.sh
```

### 停止服务
```bash
./docker-stop.sh
```

### 中国大陆网络环境构建
如果在中国大陆遇到Go模块下载超时问题，请使用专门的构建脚本：

```cmd
# Windows环境
REM 构建生产环境镜像
scripts\docker\docker-build-china.bat

REM 构建开发环境镜像
scripts\docker\docker-build-china.bat dev
```

```bash
# Linux/macOS环境
# 构建生产环境镜像
./scripts/docker/docker-build-china.sh

# 构建开发环境镜像
./scripts/docker/docker-build-china.sh dev
```

详细说明请参考：[Docker中国大陆网络问题解决方案](../docs/docker_china_guide.md)

## 环境说明

### 生产环境 (docker-compose.yml)
- MySQL: localhost:3306
- Redis: localhost:6379
- Web应用: localhost:8080
- MQTT: localhost:1883

### 开发环境 (docker-compose.dev.yml)
- MySQL: localhost:3307
- Redis: localhost:6380
- MQTT: localhost:1884
- Web应用: localhost:8081 (仅完整开发环境)
