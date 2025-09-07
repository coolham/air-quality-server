# Docker部署指南

## 概述

空气质量监测系统支持Docker容器化部署，提供生产环境和开发环境两种配置。

## 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少2GB可用内存
- 至少5GB可用磁盘空间

## 快速开始

### 生产环境部署

1. **启动完整服务**
   ```bash
   # Linux/macOS
   ./scripts/docker-start.sh
   
   # Windows
   scripts\docker-start.bat
   ```

2. **手动启动**
   ```bash
   docker-compose up --build -d
   ```

3. **访问服务**
   - Web界面: http://localhost:8082
   - Dashboard: http://localhost:8082/dashboard
   - MQTT Broker: localhost:1883

### 开发环境部署

1. **启动开发环境**
   ```bash
   # Linux/macOS
   ./scripts/docker-dev-start.sh
   
   # Windows
   scripts\docker-dev-start.bat
   ```

2. **手动启动**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

3. **服务地址**
   - MySQL: localhost:3307
   - Redis: localhost:6380
   - MQTT: localhost:1884 (可选)

## 服务配置

### 生产环境 (docker-compose.yml)

| 服务 | 端口 | 描述 |
|------|------|------|
| air-quality-server | 8082, 1883 | 主应用服务 |
| mysql | 3308 | MySQL数据库 |
| redis | 6381 | Redis缓存 |

### 开发环境 (docker-compose.dev.yml)

| 服务 | 端口 | 描述 |
|------|------|------|
| mysql-dev | 3307 | 开发数据库 |
| redis-dev | 6380 | 开发Redis |
| mosquitto-dev | 1884, 9001 | 开发MQTT Broker |

## 配置文件

### 生产环境配置
- `config/config.docker.yaml` - Docker环境专用配置
- 数据库连接: `mysql:3306` (容器内部端口)
- Redis连接: `redis:6379` (容器内部端口)

### 开发环境配置
- `config/config.dev.yaml` - 开发环境配置
- 数据库连接: `192.168.3.24:3306`
- Redis连接: `192.168.3.24:6379`

## 常用命令

### 服务管理
```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f air-quality-server
```

### 开发环境
```bash
# 启动开发环境
docker-compose -f docker-compose.dev.yml up -d

# 停止开发环境
docker-compose -f docker-compose.dev.yml down

# 查看开发环境日志
docker-compose -f docker-compose.dev.yml logs -f mysql-dev
```

### 数据管理
```bash
# 备份数据库
docker exec air-quality-mysql mysqldump -u root -padmin air_quality > backup.sql

# 恢复数据库
docker exec -i air-quality-mysql mysql -u root -padmin air_quality < backup.sql

# 清理数据卷
docker-compose down -v
```

## 健康检查

系统内置健康检查机制：

- **应用服务**: HTTP GET /health
- **MySQL**: mysqladmin ping
- **Redis**: redis-cli ping

检查间隔: 30秒
超时时间: 10秒
重试次数: 3次

## 监控和日志

### 查看日志
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs -f air-quality-server

# 查看最近100行日志
docker-compose logs --tail=100 air-quality-server
```

### 资源监控
```bash
# 查看容器资源使用
docker stats

# 查看容器详细信息
docker inspect air-quality-server
```

## 故障排除

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tulpn | grep :8080
   
   # 修改docker-compose.yml中的端口映射
   ports:
     - "8081:8080"  # 改为8081端口
   ```

2. **数据库连接失败**
   ```bash
   # 检查MySQL容器状态
   docker-compose logs mysql
   
   # 重启MySQL服务
   docker-compose restart mysql
   ```

3. **应用启动失败**
   ```bash
   # 查看应用日志
   docker-compose logs air-quality-server
   
   # 检查配置文件
   docker exec air-quality-server cat /app/config/config.docker.yaml
   ```

### 性能优化

1. **内存限制**
   ```yaml
   services:
     air-quality-server:
       deploy:
         resources:
           limits:
             memory: 512M
           reservations:
             memory: 256M
   ```

2. **CPU限制**
   ```yaml
   services:
     air-quality-server:
       deploy:
         resources:
           limits:
             cpus: '0.5'
   ```

## 安全建议

1. **修改默认密码**
   - 修改MySQL root密码
   - 修改Redis密码（如需要）
   - 修改JWT密钥

2. **网络安全**
   - 使用Docker网络隔离
   - 限制端口暴露
   - 配置防火墙规则

3. **数据安全**
   - 定期备份数据
   - 使用数据卷加密
   - 设置文件权限

## 更新和维护

### 应用更新
```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up --build -d
```

### 系统维护
```bash
# 清理未使用的镜像
docker system prune -a

# 清理数据卷
docker volume prune

# 更新基础镜像
docker-compose pull
docker-compose up -d
```

## 备份和恢复

### 数据备份
```bash
# 创建备份脚本
cat > backup.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec air-quality-mysql mysqldump -u root -padmin air_quality > "backup_${DATE}.sql"
docker exec air-quality-redis redis-cli BGSAVE
EOF

chmod +x backup.sh
```

### 数据恢复
```bash
# 恢复数据库
docker exec -i air-quality-mysql mysql -u root -padmin air_quality < backup.sql

# 恢复Redis
docker exec air-quality-redis redis-cli FLUSHALL
# 然后复制RDB文件到容器
```

## 扩展部署

### 多实例部署
```yaml
services:
  air-quality-server-1:
    build: .
    ports:
      - "8080:8080"
  
  air-quality-server-2:
    build: .
    ports:
      - "8081:8080"
```

### 负载均衡
```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - air-quality-server-1
      - air-quality-server-2
```

## 端口配置与冲突解决

### 端口映射表

#### 生产环境 (docker-compose.yml)

| 服务 | 宿主机端口 | 容器端口 | 说明 |
|------|------------|----------|------|
| air-quality-server | **8082** | 8080 | Web应用服务 |
| air-quality-server | 1883 | 1883 | MQTT Broker |
| mysql | **3308** | 3306 | MySQL数据库 |
| redis | **6381** | 6379 | Redis缓存 |

#### 开发环境 (docker-compose.dev.yml)

| 服务 | 宿主机端口 | 容器端口 | 说明 |
|------|------------|----------|------|
| mysql-dev | **3307** | 3306 | MySQL数据库 |
| redis-dev | **6380** | 6379 | Redis缓存 |
| mosquitto-dev | **1884** | 1883 | MQTT Broker |
| mosquitto-dev | 9001 | 9001 | MQTT WebSocket |
| air-quality-server-dev | **8083** | 8080 | Web应用服务 |
| air-quality-server-dev | **1885** | 1883 | MQTT Broker |

### 端口冲突解决方案

#### 问题描述
如果宿主机已经占用了以下端口：
- 3306 (MySQL默认端口)
- 6379 (Redis默认端口)
- 8080 (Web应用默认端口)

#### 解决方案
1. **生产环境**：使用3308、6381、8082端口
2. **开发环境**：使用3307、6380、8083端口

#### 连接方式

##### 从宿主机连接数据库
```bash
# 生产环境
mysql -h localhost -P 3308 -u root -p
redis-cli -h localhost -p 6381

# 开发环境  
mysql -h localhost -P 3307 -u root -p
redis-cli -h localhost -p 6380
```

##### 从应用程序连接（容器内部）
```yaml
# 配置文件中的连接方式不变
database:
  host: "mysql"  # Docker服务名
  port: 3306     # 容器内部端口

redis:
  host: "redis"  # Docker服务名
  port: 6379     # 容器内部端口
```

### 端口检查

#### 检查端口占用
```bash
# Windows
netstat -an | findstr :3306
netstat -an | findstr :6379
netstat -an | findstr :8080

# Linux/macOS
lsof -i :3306
lsof -i :6379
lsof -i :8080
```

#### 检查Docker容器端口
```bash
# 查看所有容器端口映射
docker ps

# 查看特定容器端口
docker port air-quality-mysql
docker port air-quality-redis
```

### 自定义端口配置

如果需要使用其他端口，可以修改docker-compose文件：

```yaml
# docker-compose.yml
services:
  mysql:
    ports:
      - "3309:3306"  # 修改宿主机端口为3309
      
  redis:
    ports:
      - "6382:6379"  # 修改宿主机端口为6382
      
  air-quality-server:
    ports:
      - "8084:8080"  # 修改宿主机端口为8084
```

### 常见问题

#### Q: 为什么容器内部端口不变？
A: 容器内部端口保持标准端口（3306、6379、8080），这样应用程序配置不需要修改。只有宿主机端口映射改变。

#### Q: 如何从外部访问数据库？
A: 使用宿主机端口访问：
- 生产环境：localhost:3308
- 开发环境：localhost:3307

#### Q: 应用程序如何连接数据库？
A: 在Docker容器内部，应用程序使用服务名和标准端口连接：
- 生产环境：mysql:3306
- 开发环境：mysql-dev:3306

## 联系支持

如遇到问题，请查看：
1. 系统日志: `docker-compose logs`
2. 健康检查: `docker-compose ps`
3. 资源使用: `docker stats`

更多帮助请参考项目文档或提交Issue。
