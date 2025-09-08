# Docker 数据库管理

## 概述

本项目使用 Docker 容器化部署，数据库通过 MySQL 官方镜像的初始化脚本自动创建。

## 解决方案

**问题**：`Error 1054 (42S22): Unknown column 'o3' in 'field list'`

**解决**：更新了 `scripts/database/init.sql` 文件，包含所有必要的字段：
- `o3` (臭氧浓度)
- `no2` (二氧化氮浓度)  
- `so2` (二氧化硫浓度)
- `co` (一氧化碳浓度)
- `voc` (挥发性有机化合物)
- `signal_strength` (信号强度)

## 使用方法

### 快速启动

**Windows:**
```bash
cd scripts/docker
quick-start.bat
```

**Linux/macOS:**
```bash
cd scripts/docker
chmod +x quick-start.sh
./quick-start.sh
```

### 基本操作

| 操作 | Windows | Linux/macOS |
|------|---------|-------------|
| 启动服务 | `quick-start.bat` | `./quick-start.sh` |
| 停止服务 | `stop.bat` | `./stop.sh` |
| 清理数据 | `clean.bat` | `./clean.sh` |

### 直接使用 docker-compose

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 清理所有数据
docker-compose down -v
```

## 服务地址

- **Web 界面**: http://localhost:8082
- **MQTT 端口**: 1883
- **MySQL 端口**: 3308

## 数据库配置

- **主机**: localhost (容器内: mysql)
- **端口**: 3308 (容器内: 3306)
- **用户**: root
- **密码**: admin
- **数据库**: air_quality

## 注意事项

1. **首次启动**：数据库会自动通过 `init.sql` 脚本初始化
2. **数据持久化**：数据保存在 Docker 卷中，重启容器不会丢失数据
3. **清理数据**：使用 `clean.bat/clean.sh` 会删除所有数据，请谨慎使用

## 故障排除

### 端口冲突
如果端口被占用，可以修改 `docker-compose.yml` 中的端口映射：
```yaml
ports:
  - "8083:8080"  # 改为其他端口
```

### 数据库连接失败
1. 检查 MySQL 容器是否正常运行：`docker-compose ps`
2. 查看容器日志：`docker-compose logs mysql`
3. 等待数据库完全启动（首次启动需要时间）

### 清理重新开始
```bash
# 停止并删除所有数据
docker-compose down -v

# 重新启动
docker-compose up -d
```