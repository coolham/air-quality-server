# 数据库设计与管理指南

## 概述

本文档是空气质量监测系统数据库的完整指南，包含数据库设计、设置、数据模型和管理的所有内容。

## 1. 数据库选型与架构

### 1.1 主要数据库
- **MySQL 8.0**: 存储所有结构化数据
- **Redis 7.0**: 缓存、会话存储、消息队列

### 1.2 设计原则
- 简化架构，减少组件依赖
- 使用MySQL存储时序数据，通过索引优化查询性能
- 使用Redis作为消息队列，简化消息传递

### 1.3 数据库结构
```sql
-- 创建数据库
CREATE DATABASE air_quality CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE air_quality;
```

## 2. 数据模型设计

### 2.1 核心模型分类

1. **用户管理模型** - 用户、角色、权限管理
2. **设备管理模型** - 设备信息、状态、配置
3. **数据存储模型** - 传感器数据、统计数据
4. **告警管理模型** - 告警规则、告警记录
5. **系统配置模型** - 系统参数、业务配置

### 2.2 用户管理模型

#### 2.2.1 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    email VARCHAR(100) UNIQUE NOT NULL COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active' COMMENT '用户状态',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

#### 2.2.2 角色表 (roles)
```sql
CREATE TABLE roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
    name VARCHAR(50) UNIQUE NOT NULL COMMENT '角色名称',
    description VARCHAR(200) COMMENT '角色描述',
    permissions JSON COMMENT '权限列表',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';
```

#### 2.2.3 用户角色关联表 (user_roles)
```sql
CREATE TABLE user_roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';
```

### 2.3 设备管理模型

#### 2.3.1 设备表 (devices)
```sql
CREATE TABLE devices (
    id VARCHAR(64) PRIMARY KEY COMMENT '设备ID',
    name VARCHAR(100) NOT NULL COMMENT '设备名称',
    type VARCHAR(50) NOT NULL COMMENT '设备类型',
    location_latitude DECIMAL(10, 8) COMMENT '纬度',
    location_longitude DECIMAL(11, 8) COMMENT '经度',
    location_address VARCHAR(200) COMMENT '地址',
    status ENUM('online', 'offline', 'maintenance') DEFAULT 'offline' COMMENT '设备状态',
    config JSON COMMENT '设备配置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status (status),
    INDEX idx_location (location_latitude, location_longitude),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备信息表';
```

#### 2.3.2 设备运行时状态表 (device_runtime_status)
```sql
CREATE TABLE device_runtime_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '状态ID',
    device_id VARCHAR(64) NOT NULL COMMENT '设备ID',
    online BOOLEAN DEFAULT FALSE COMMENT '在线状态',
    battery_level INT COMMENT '电池电量百分比',
    signal_strength INT COMMENT '信号强度(dBm)',
    last_data_time TIMESTAMP NULL COMMENT '最后数据时间',
    last_heartbeat TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后心跳时间',
    error_code INT DEFAULT 0 COMMENT '错误代码',
    error_message TEXT COMMENT '错误信息',
    firmware_version VARCHAR(50) COMMENT '固件版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_device_id (device_id),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备运行时状态表';
```

### 2.4 数据存储模型

#### 2.4.1 统一传感器数据表 (unified_sensor_data)
```sql
CREATE TABLE unified_sensor_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '数据ID',
    device_id VARCHAR(64) NOT NULL COMMENT '设备ID',
    device_type VARCHAR(50) NOT NULL COMMENT '设备类型',
    sensor_id VARCHAR(64) COMMENT '传感器ID',
    sensor_type VARCHAR(50) COMMENT '传感器类型',
    timestamp TIMESTAMP NOT NULL COMMENT '数据时间戳',
    
    -- 核心环境指标
    pm25 DECIMAL(8, 3) COMMENT 'PM2.5浓度 μg/m³',
    pm10 DECIMAL(8, 3) COMMENT 'PM10浓度 μg/m³',
    co2 DECIMAL(8, 3) COMMENT 'CO2浓度 ppm',
    formaldehyde DECIMAL(8, 3) COMMENT '甲醛浓度 mg/m³',
    
    -- 环境参数
    temperature DECIMAL(6, 2) COMMENT '温度 °C',
    humidity DECIMAL(6, 2) COMMENT '湿度 %',
    pressure DECIMAL(8, 2) COMMENT '气压 hPa',
    
    -- 其他污染物指标
    o3 DECIMAL(8, 3) COMMENT '臭氧浓度 μg/m³',
    no2 DECIMAL(8, 3) COMMENT '二氧化氮浓度 μg/m³',
    so2 DECIMAL(8, 3) COMMENT '二氧化硫浓度 μg/m³',
    co DECIMAL(8, 3) COMMENT '一氧化碳浓度 mg/m³',
    voc DECIMAL(8, 3) COMMENT '挥发性有机化合物 μg/m³',
    
    -- 设备状态信息
    battery INT COMMENT '电池电量 %',
    signal_strength INT COMMENT '信号强度 dBm',
    data_quality VARCHAR(20) DEFAULT 'good' COMMENT '数据质量',
    
    -- 位置信息
    latitude DECIMAL(10, 8) COMMENT '纬度',
    longitude DECIMAL(11, 8) COMMENT '经度',
    
    -- 扩展数据
    extended_data JSON COMMENT '扩展数据',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_device_timestamp (device_id, timestamp),
    INDEX idx_device_type (device_type),
    INDEX idx_sensor_id (sensor_id),
    INDEX idx_sensor_type (sensor_type),
    INDEX idx_timestamp (timestamp),
    INDEX idx_formaldehyde (formaldehyde),
    INDEX idx_temperature (temperature),
    INDEX idx_humidity (humidity),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统一传感器数据表';
```

#### 2.4.2 甲醛专用数据表 (formaldehyde_data)
```sql
CREATE TABLE formaldehyde_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '数据ID',
    device_id VARCHAR(64) NOT NULL COMMENT '设备ID',
    timestamp TIMESTAMP NOT NULL COMMENT '数据时间戳',
    formaldehyde DECIMAL(8, 3) COMMENT '甲醛浓度 mg/m³',
    temperature DECIMAL(6, 2) COMMENT '温度 °C',
    humidity DECIMAL(6, 2) COMMENT '湿度 %',
    battery INT COMMENT '电池电量 %',
    signal_strength INT COMMENT '信号强度 dBm',
    data_quality VARCHAR(20) DEFAULT 'good' COMMENT '数据质量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_device_timestamp (device_id, timestamp),
    INDEX idx_timestamp (timestamp),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='甲醛数据表';
```

### 2.5 告警管理模型

#### 2.5.1 告警规则表 (alert_rules)
```sql
CREATE TABLE alert_rules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '规则ID',
    name VARCHAR(100) NOT NULL COMMENT '规则名称',
    device_id VARCHAR(64) COMMENT '设备ID，NULL表示所有设备',
    metric VARCHAR(50) NOT NULL COMMENT '监控指标',
    condition_type ENUM('gt', 'lt', 'eq', 'ne', 'gte', 'lte') NOT NULL COMMENT '条件类型',
    threshold_value DECIMAL(10, 2) NOT NULL COMMENT '阈值',
    duration_seconds INT DEFAULT 0 COMMENT '持续时间(秒)',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning' COMMENT '严重程度',
    enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    notification_channels JSON COMMENT '通知渠道',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_device_id (device_id),
    INDEX idx_enabled (enabled),
    INDEX idx_severity (severity),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';
```

#### 2.5.2 告警记录表 (alerts)
```sql
CREATE TABLE alerts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '告警ID',
    rule_id BIGINT NOT NULL COMMENT '规则ID',
    device_id VARCHAR(64) NOT NULL COMMENT '设备ID',
    metric VARCHAR(50) NOT NULL COMMENT '监控指标',
    current_value DECIMAL(10, 2) NOT NULL COMMENT '当前值',
    threshold_value DECIMAL(10, 2) NOT NULL COMMENT '阈值',
    severity ENUM('critical', 'warning', 'info') NOT NULL COMMENT '严重程度',
    status ENUM('active', 'acknowledged', 'resolved') DEFAULT 'active' COMMENT '告警状态',
    triggered_at TIMESTAMP NOT NULL COMMENT '触发时间',
    acknowledged_at TIMESTAMP NULL COMMENT '确认时间',
    resolved_at TIMESTAMP NULL COMMENT '解决时间',
    acknowledged_by BIGINT NULL COMMENT '确认人',
    resolved_by BIGINT NULL COMMENT '解决人',
    message TEXT COMMENT '告警消息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_rule_id (rule_id),
    INDEX idx_device_id (device_id),
    INDEX idx_status (status),
    INDEX idx_triggered_at (triggered_at),
    INDEX idx_severity (severity),
    FOREIGN KEY (rule_id) REFERENCES alert_rules(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (acknowledged_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警记录表';
```

### 2.6 系统配置模型

#### 2.6.1 系统配置表 (system_configs)
```sql
CREATE TABLE system_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
    key_name VARCHAR(100) UNIQUE NOT NULL COMMENT '配置键',
    value TEXT COMMENT '配置值',
    description VARCHAR(200) COMMENT '配置描述',
    category VARCHAR(50) COMMENT '配置分类',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_key_name (key_name),
    INDEX idx_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
```

## 3. 数据库设置与初始化

### 3.1 前置条件

1. **MySQL 8.0+** 已安装并运行
2. **Go 1.21+** 开发环境
3. **配置文件** 已正确配置数据库连接信息

### 3.2 配置数据库连接

编辑 `config/config.yaml` 文件，确保数据库配置正确：

```yaml
database:
  host: localhost
  port: 3306
  username: air_quality
  password: air_quality123
  database: air_quality
  charset: utf8mb4
  max_idle: 10
  max_open: 100
  max_life: 3600
```

### 3.3 创建数据库

```sql
-- 连接到MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE air_quality CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选）
CREATE USER 'air_quality'@'localhost' IDENTIFIED BY 'air_quality123';
GRANT ALL PRIVILEGES ON air_quality.* TO 'air_quality'@'localhost';
FLUSH PRIVILEGES;
```

### 3.4 初始化数据库

使用以下命令之一初始化数据库：

#### 方法1：使用Makefile（推荐）
```bash
# 初始化数据库
make migrate

# 查看数据库状态
make migrate-status
```

#### 方法2：直接运行迁移工具
```bash
# 初始化数据库
go run cmd/migrate/main.go -action init

# 查看数据库状态
go run cmd/migrate/main.go -action status

# 查看帮助信息
go run cmd/migrate/main.go -help
```

#### 方法3：使用指定配置文件
```bash
# 使用指定配置文件初始化
go run cmd/migrate/main.go -config config/config.yaml -action init
```

### 3.5 初始化过程

数据库初始化工具会执行以下操作：

1. **创建表结构** - 自动创建所有数据表
2. **创建索引** - 为所有表创建必要的索引以优化查询性能
3. **插入初始数据** - 自动插入默认数据

### 3.6 默认数据

#### 默认角色
- `admin` - 系统管理员（全部权限）
- `operator` - 操作员（设备、数据、告警管理）
- `viewer` - 查看者（只读权限）

#### 默认用户
- 用户名：`admin`
- 密码：`admin123`
- 角色：管理员

#### 默认系统配置
- 数据保留天数：365天
- 告警检查间隔：60秒
- 每用户最大设备数：100
- API请求限制：1000/小时
- 数据质量阈值：0.8
- MQTT配置参数

#### 默认告警规则
- 甲醛浓度警告：> 0.08 mg/m³
- 甲醛浓度严重：> 0.1 mg/m³
- 设备离线告警

#### 示例设备
- `hcho_001` - 甲醛传感器001
- `hcho_002` - 甲醛传感器002

## 4. 性能优化

### 4.1 数据分区策略

#### 4.1.1 空气质量数据表分区
```sql
-- 按月分区，提高查询性能
ALTER TABLE unified_sensor_data 
PARTITION BY RANGE (UNIX_TIMESTAMP(timestamp)) (
    PARTITION p202401 VALUES LESS THAN (UNIX_TIMESTAMP('2024-02-01')),
    PARTITION p202402 VALUES LESS THAN (UNIX_TIMESTAMP('2024-03-01')),
    PARTITION p202403 VALUES LESS THAN (UNIX_TIMESTAMP('2024-04-01')),
    PARTITION p202404 VALUES LESS THAN (UNIX_TIMESTAMP('2024-05-01')),
    PARTITION p202405 VALUES LESS THAN (UNIX_TIMESTAMP('2024-06-01')),
    PARTITION p202406 VALUES LESS THAN (UNIX_TIMESTAMP('2024-07-01')),
    PARTITION p202407 VALUES LESS THAN (UNIX_TIMESTAMP('2024-08-01')),
    PARTITION p202408 VALUES LESS THAN (UNIX_TIMESTAMP('2024-09-01')),
    PARTITION p202409 VALUES LESS THAN (UNIX_TIMESTAMP('2024-10-01')),
    PARTITION p202410 VALUES LESS THAN (UNIX_TIMESTAMP('2024-11-01')),
    PARTITION p202411 VALUES LESS THAN (UNIX_TIMESTAMP('2024-12-01')),
    PARTITION p202412 VALUES LESS THAN (UNIX_TIMESTAMP('2025-01-01')),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

### 4.2 索引优化

#### 4.2.1 复合索引
```sql
-- 设备数据查询优化
CREATE INDEX idx_device_timestamp_metric ON unified_sensor_data (device_id, timestamp, formaldehyde);
CREATE INDEX idx_device_timestamp_temp ON unified_sensor_data (device_id, timestamp, temperature);
CREATE INDEX idx_device_timestamp_humidity ON unified_sensor_data (device_id, timestamp, humidity);

-- 告警查询优化
CREATE INDEX idx_alert_device_status ON alerts (device_id, status, triggered_at);
CREATE INDEX idx_alert_severity_status ON alerts (severity, status, triggered_at);
```

### 4.3 查询优化
- 使用合适的索引
- 避免全表扫描
- 使用分区表
- 优化JOIN查询

### 4.4 写入优化
- 批量插入数据
- 使用事务
- 异步写入
- 数据压缩

## 5. Redis 数据结构设计

### 5.1 缓存键设计

#### 5.1.1 设备信息缓存
```
# 设备基本信息
device:info:{device_id} -> JSON

# 设备状态
device:status:{device_id} -> "online" | "offline" | "maintenance"

# 设备配置
device:config:{device_id} -> JSON
```

#### 5.1.2 实时数据缓存
```
# 设备最新数据
device:latest:{device_id} -> JSON

# 设备统计数据（1小时）
device:stats:1h:{device_id} -> JSON

# 设备统计数据（1天）
device:stats:1d:{device_id} -> JSON
```

#### 5.1.3 用户会话缓存
```
# 用户会话
session:{session_id} -> JSON

# 用户权限
user:permissions:{user_id} -> JSON

# 用户Token
token:{token} -> JSON
```

### 5.2 消息队列设计

#### 5.2.1 Redis Pub/Sub Channels
```
# 原始数据通道
air-quality:raw-data

# 处理后的数据通道
air-quality:processed-data

# 告警事件通道
air-quality:alerts

# 设备事件通道
air-quality:device-events
```

#### 5.2.2 Redis Streams (可选)
```
# 数据流
air-quality:data-stream

# 告警流
air-quality:alert-stream
```

## 6. 数据迁移和备份

### 6.1 数据迁移脚本
```sql
-- 创建迁移版本表
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入初始版本
INSERT INTO schema_migrations (version) VALUES ('20240101000000');
```

### 6.2 数据备份策略
- **全量备份**: 每日凌晨进行全量备份
- **增量备份**: 每小时进行增量备份
- **日志备份**: 实时备份binlog
- **数据归档**: 超过1年的数据归档到对象存储

## 7. 监控和维护

### 7.1 数据库监控
- 连接数监控
- 查询性能监控
- 磁盘空间监控
- 慢查询监控

### 7.2 数据质量
- 数据完整性检查
- 数据一致性验证
- 异常数据检测
- 数据清理任务

### 7.3 缓存策略
- 热点数据缓存
- 查询结果缓存
- 分布式缓存
- 缓存预热

## 8. 安全设计

### 8.1 访问控制
- 数据库用户权限管理
- 网络访问控制
- 加密传输
- 审计日志

### 8.2 数据保护
- 敏感数据加密
- 数据脱敏
- 备份加密
- 访问日志记录

## 9. 验证安装

### 9.1 检查数据库状态
```bash
make migrate-status
```

输出应该显示所有表都已创建，并且包含初始数据。

### 9.2 连接数据库验证
```sql
-- 连接到数据库
mysql -u air_quality -p air_quality

-- 查看所有表
SHOW TABLES;

-- 检查用户数据
SELECT * FROM users;

-- 检查角色数据
SELECT * FROM roles;

-- 检查系统配置
SELECT * FROM system_configs;
```

### 9.3 启动应用程序
```bash
# 构建并启动应用
make start

# 或者直接运行
go run cmd/air-quality-server/main.go
```

访问 `http://localhost:8080/health` 检查应用是否正常启动。

## 10. 故障排除

### 10.1 常见问题

#### 数据库连接失败
**错误信息**：`连接数据库失败`

**解决方案**：
- 检查MySQL服务是否运行
- 验证数据库连接配置
- 确认用户权限

#### 表已存在错误
**错误信息**：`table already exists`

**解决方案**：
- 如果数据库已初始化，这是正常现象
- 如需重新初始化，请先清空数据库

#### 权限不足
**错误信息**：`Access denied`

**解决方案**：
- 检查数据库用户权限
- 确保用户有CREATE、INSERT、SELECT权限

### 10.2 重新初始化

如果需要重新初始化数据库：

```sql
-- 删除数据库（谨慎操作）
DROP DATABASE air_quality;

-- 重新创建数据库
CREATE DATABASE air_quality CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

然后重新运行初始化命令。

## 11. 生产环境注意事项

1. **数据备份**：在生产环境中，执行任何数据库操作前请先备份数据
2. **用户权限**：使用最小权限原则，避免使用root用户
3. **网络安全**：配置防火墙，限制数据库访问
4. **监控告警**：设置数据库性能监控和告警
5. **定期维护**：定期检查和优化数据库性能

## 12. 下一步

数据库初始化完成后，您可以：

1. 启动应用程序：`make start`
2. 配置MQTT服务器
3. 添加实际设备
4. 设置告警规则
5. 配置用户权限

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队
