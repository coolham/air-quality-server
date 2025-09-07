# 部署指南

## 1. 部署架构

### 1.1 简化部署架构
```
┌─────────────────────────────────────────────────────────────┐
│                    负载均衡器 (Nginx)                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
┌───────▼──────┐ ┌───▼──────┐ ┌────▼──────┐
│ 数据接收服务  │ │ 数据查询  │ │ 设备管理  │
│ (端口:8081)  │ │ (端口:8082)│ │ (端口:8083)│
└───────┬──────┘ └──────────┘ └───────────┘
        │
┌───────▼──────┐
│ 数据处理服务  │
│ (端口:8084)  │
└───────┬──────┘
        │
┌───────▼──────┐
│ 告警服务     │
│ (端口:8085)  │
└──────────────┘
        │
┌───────▼──────┐
│ MySQL + Redis│
│ (端口:3306)  │
│ (端口:6379)  │
└──────────────┘
```

### 1.2 服务端口分配
- **数据接收服务**: 8081
- **数据查询服务**: 8082
- **设备管理服务**: 8083
- **数据处理服务**: 8084
- **告警服务**: 8085
- **用户管理服务**: 8086
- **配置管理服务**: 8087

## 2. 环境要求

### 2.1 硬件要求
- **CPU**: 4核心以上
- **内存**: 8GB以上
- **存储**: 100GB以上SSD
- **网络**: 100Mbps以上

### 2.2 软件要求
- **操作系统**: Ubuntu 20.04+ / CentOS 8+ / Windows Server 2019+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Go**: 1.21+

## 3. 快速部署

### 3.1 使用Docker Compose部署

#### 3.1.1 克隆项目
```bash
git clone <repository-url>
cd air-quality-server
```

#### 3.1.2 配置环境变量
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
vim .env
```

#### 3.1.3 启动服务
```bash
# 启动基础服务 (MySQL, Redis)
docker-compose up -d mysql redis

# 等待数据库启动完成
sleep 30

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 3.2 手动部署

#### 3.2.1 安装依赖
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y mysql-server redis-server nginx

# CentOS/RHEL
sudo yum install -y mysql-server redis nginx
```

#### 3.2.2 配置数据库
```bash
# 启动MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# 创建数据库和用户
mysql -u root -p < scripts/init.sql

# 启动Redis
sudo systemctl start redis
sudo systemctl enable redis
```

#### 3.2.3 构建和启动服务
```bash
# 安装Go依赖
go mod download

# 构建所有服务
make build

# 启动所有服务
make start-all
```

## 4. 配置说明

### 4.1 环境变量配置
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_NAME=air_quality
DB_USER=air_quality
DB_PASSWORD=air_quality123

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# 服务配置
SERVICE_PORT=8081
LOG_LEVEL=info
ENVIRONMENT=production

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
```

### 4.2 Nginx配置
```nginx
upstream air_quality_backend {
    server 127.0.0.1:8081;  # 数据接收服务
    server 127.0.0.1:8082;  # 数据查询服务
    server 127.0.0.1:8083;  # 设备管理服务
    server 127.0.0.1:8084;  # 数据处理服务
    server 127.0.0.1:8085;  # 告警服务
    server 127.0.0.1:8086;  # 用户管理服务
    server 127.0.0.1:8087;  # 配置管理服务
}

server {
    listen 80;
    server_name your-domain.com;

    # API路由
    location /api/v1/data/ {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/query/ {
        proxy_pass http://127.0.0.1:8082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/devices/ {
        proxy_pass http://127.0.0.1:8083;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/alerts/ {
        proxy_pass http://127.0.0.1:8085;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/users/ {
        proxy_pass http://127.0.0.1:8086;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket支持
    location /ws/ {
        proxy_pass http://127.0.0.1:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 5. 监控和日志

### 5.1 服务监控
```bash
# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]

# 查看资源使用情况
docker stats
```

### 5.2 日志管理
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f data-ingestion

# 查看错误日志
docker-compose logs | grep ERROR
```

### 5.3 健康检查
```bash
# 检查服务健康状态
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

## 6. 数据备份

### 6.1 数据库备份
```bash
# 创建备份脚本
cat > backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backup/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# 备份MySQL
mysqldump -u root -p air_quality > $BACKUP_DIR/air_quality.sql

# 备份Redis
redis-cli --rdb $BACKUP_DIR/redis.rdb

echo "Backup completed: $BACKUP_DIR"
EOF

chmod +x backup.sh

# 设置定时备份
crontab -e
# 添加: 0 2 * * * /path/to/backup.sh
```

### 6.2 数据恢复
```bash
# 恢复MySQL数据
mysql -u root -p air_quality < /backup/20240101/air_quality.sql

# 恢复Redis数据
redis-cli --rdb /backup/20240101/redis.rdb
```

## 7. 性能优化

### 7.1 数据库优化
```sql
-- 优化MySQL配置
SET GLOBAL innodb_buffer_pool_size = 2G;
SET GLOBAL max_connections = 1000;
SET GLOBAL query_cache_size = 256M;

-- 创建索引
CREATE INDEX idx_air_quality_data_device_timestamp ON air_quality_data(device_id, timestamp);
CREATE INDEX idx_air_quality_data_timestamp ON air_quality_data(timestamp);
```

### 7.2 Redis优化
```bash
# 优化Redis配置
echo "maxmemory 2gb" >> /etc/redis/redis.conf
echo "maxmemory-policy allkeys-lru" >> /etc/redis/redis.conf
systemctl restart redis
```

### 7.3 应用优化
```bash
# 设置Go环境变量
export GOGC=100
export GOMAXPROCS=4

# 启动服务时设置参数
./bin/data-ingestion -workers=10 -batch-size=100
```

## 8. 故障排除

### 8.1 常见问题

#### 8.1.1 服务启动失败
```bash
# 检查端口占用
netstat -tlnp | grep :8081

# 检查日志
docker-compose logs data-ingestion

# 检查配置文件
cat config/config.yaml
```

#### 8.1.2 数据库连接失败
```bash
# 检查MySQL状态
systemctl status mysql

# 检查连接
mysql -u air_quality -p -h localhost air_quality

# 检查防火墙
ufw status
```

#### 8.1.3 Redis连接失败
```bash
# 检查Redis状态
systemctl status redis

# 检查连接
redis-cli ping

# 检查配置
cat /etc/redis/redis.conf
```

### 8.2 性能问题

#### 8.2.1 高CPU使用率
```bash
# 查看进程
top -p $(pgrep data-ingestion)

# 分析性能
go tool pprof http://localhost:8081/debug/pprof/profile
```

#### 8.2.2 高内存使用率
```bash
# 查看内存使用
free -h

# 查看进程内存
ps aux | grep data-ingestion

# 分析内存
go tool pprof http://localhost:8081/debug/pprof/heap
```

## 9. 安全配置

### 9.1 防火墙配置
```bash
# Ubuntu/Debian
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable

# CentOS/RHEL
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload
```

### 9.2 SSL证书配置
```bash
# 使用Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 9.3 数据库安全
```sql
-- 创建只读用户
CREATE USER 'readonly'@'%' IDENTIFIED BY 'readonly_password';
GRANT SELECT ON air_quality.* TO 'readonly'@'%';

-- 限制用户权限
REVOKE ALL PRIVILEGES ON *.* FROM 'air_quality'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE ON air_quality.* TO 'air_quality'@'%';
```

## 10. 扩展部署

### 10.1 水平扩展
```bash
# 启动多个实例
./bin/data-ingestion -port=8081 &
./bin/data-ingestion -port=8088 &
./bin/data-ingestion -port=8089 &

# 使用负载均衡器分发请求
```

### 10.2 数据库主从复制
```sql
-- 主库配置
[mysqld]
server-id = 1
log-bin = mysql-bin
binlog-format = ROW

-- 从库配置
[mysqld]
server-id = 2
relay-log = mysql-relay-bin
read-only = 1
```

### 10.3 Redis集群
```bash
# 启动Redis集群
redis-server --port 7000 --cluster-enabled yes --cluster-config-file nodes-7000.conf
redis-server --port 7001 --cluster-enabled yes --cluster-config-file nodes-7001.conf
redis-server --port 7002 --cluster-enabled yes --cluster-config-file nodes-7002.conf

# 创建集群
redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002
```

## 11. 维护操作

### 11.1 服务更新
```bash
# 停止服务
docker-compose down

# 更新代码
git pull

# 重新构建
docker-compose build

# 启动服务
docker-compose up -d
```

### 11.2 数据清理
```sql
-- 清理旧数据
DELETE FROM air_quality_data WHERE timestamp < DATE_SUB(NOW(), INTERVAL 1 YEAR);

-- 优化表
OPTIMIZE TABLE air_quality_data;
```

### 11.3 配置更新
```bash
# 更新配置
vim config/config.yaml

# 重启服务
docker-compose restart [service-name]
```
