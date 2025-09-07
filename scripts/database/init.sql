-- 空气质量监测系统数据库初始化脚本
-- 版本: v2.0 (新系统完整初始化)
-- 创建日期: 2024-09-07
-- 说明: 此脚本用于新系统的完整数据库初始化，包含所有表结构和初始数据
-- 注意: 此脚本不兼容旧系统，仅适用于全新部署

-- 创建数据库
CREATE DATABASE IF NOT EXISTS air_quality CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE air_quality;

-- 设备表
CREATE TABLE IF NOT EXISTS devices (
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

-- 统一传感器数据表
CREATE TABLE IF NOT EXISTS unified_sensor_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '数据ID',
    device_id VARCHAR(64) NOT NULL COMMENT '设备ID',
    device_type VARCHAR(50) NOT NULL COMMENT '设备类型',
    sensor_id VARCHAR(64) COMMENT '传感器ID',
    sensor_type VARCHAR(50) COMMENT '传感器类型',
    timestamp TIMESTAMP NOT NULL COMMENT '数据时间戳',
    pm25 DECIMAL(8, 2) COMMENT 'PM2.5浓度',
    pm10 DECIMAL(8, 2) COMMENT 'PM10浓度',
    co2 DECIMAL(8, 2) COMMENT 'CO2浓度',
    formaldehyde DECIMAL(8, 2) COMMENT '甲醛浓度',
    temperature DECIMAL(6, 2) COMMENT '温度',
    humidity DECIMAL(6, 2) COMMENT '湿度',
    pressure DECIMAL(8, 2) COMMENT '气压',
    battery DECIMAL(5, 2) COMMENT '电池电量',
    data_quality VARCHAR(20) DEFAULT 'good' COMMENT '数据质量',
    location_latitude DECIMAL(10, 8) COMMENT '纬度',
    location_longitude DECIMAL(11, 8) COMMENT '经度',
    location_address VARCHAR(200) COMMENT '地址',
    quality_score DECIMAL(4, 2) COMMENT '数据质量评分',
    extended_data JSON COMMENT '扩展数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_device_timestamp (device_id, timestamp),
    INDEX idx_device_type (device_type),
    INDEX idx_sensor_id (sensor_id),
    INDEX idx_sensor_type (sensor_type),
    INDEX idx_timestamp (timestamp),
    INDEX idx_device_id (device_id),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统一传感器数据表';

-- 用户表
CREATE TABLE IF NOT EXISTS users (
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

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
    name VARCHAR(50) UNIQUE NOT NULL COMMENT '角色名称',
    description VARCHAR(200) COMMENT '角色描述',
    permissions JSON COMMENT '权限列表',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 告警规则表
CREATE TABLE IF NOT EXISTS alert_rules (
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

-- 告警记录表
CREATE TABLE IF NOT EXISTS alerts (
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

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
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

-- 设备运行时状态表
CREATE TABLE IF NOT EXISTS device_runtime_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '状态ID',
    device_id VARCHAR(64) NOT NULL COMMENT '设备ID',
    status ENUM('online', 'offline', 'maintenance', 'error') DEFAULT 'offline' COMMENT '设备状态',
    last_seen TIMESTAMP NULL COMMENT '最后在线时间',
    uptime_seconds BIGINT DEFAULT 0 COMMENT '运行时间(秒)',
    memory_usage DECIMAL(5, 2) COMMENT '内存使用率(%)',
    cpu_usage DECIMAL(5, 2) COMMENT 'CPU使用率(%)',
    disk_usage DECIMAL(5, 2) COMMENT '磁盘使用率(%)',
    network_status ENUM('connected', 'disconnected', 'unstable') DEFAULT 'disconnected' COMMENT '网络状态',
    firmware_version VARCHAR(50) COMMENT '固件版本',
    config_version VARCHAR(50) COMMENT '配置版本',
    error_count INT DEFAULT 0 COMMENT '错误计数',
    last_error TEXT COMMENT '最后错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_device_id (device_id),
    INDEX idx_status (status),
    INDEX idx_last_seen (last_seen),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备运行时状态表';


-- 插入默认角色
INSERT IGNORE INTO roles (name, description, permissions) VALUES 
('admin', '系统管理员', '["*"]'),
('operator', '操作员', '["device:read", "device:write", "data:read", "alert:read", "alert:write"]'),
('viewer', '查看者', '["device:read", "data:read", "alert:read"]');

-- 插入默认用户 (密码: admin123)
INSERT IGNORE INTO users (username, email, password_hash, status) VALUES 
('admin', 'admin@air-quality.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active');

-- 为用户分配管理员角色
INSERT IGNORE INTO user_roles (user_id, role_id) 
SELECT u.id, r.id FROM users u, roles r WHERE u.username = 'admin' AND r.name = 'admin';

-- 插入默认系统配置
INSERT IGNORE INTO system_configs (key_name, value, description, category) VALUES 
('data_retention_days', '365', '数据保留天数', 'data'),
('alert_check_interval', '60', '告警检查间隔(秒)', 'alert'),
('max_devices_per_user', '100', '每用户最大设备数', 'device'),
('api_rate_limit', '1000', 'API请求限制(每小时)', 'api'),
('data_quality_threshold', '0.8', '数据质量阈值', 'data');

-- 插入示例设备
INSERT IGNORE INTO devices (id, name, type, location_latitude, location_longitude, location_address, status, config) VALUES 
('ESP32_001', '测试设备001', 'ESP32', 39.9042, 116.4074, '北京市朝阳区', 'online', '{"report_interval": 60, "sensors": {"pm25": true, "pm10": true, "co2": true, "temperature": true, "humidity": true, "pressure": true}}'),
('ESP32_002', '测试设备002', 'ESP32', 31.2304, 121.4737, '上海市黄浦区', 'online', '{"report_interval": 60, "sensors": {"pm25": true, "pm10": true, "co2": true, "temperature": true, "humidity": true, "pressure": true}}'),
('HCHO_001', '甲醛监测设备001', 'hcho', 39.9042, 116.4074, '北京市朝阳区', 'online', '{"report_interval": 60, "sensors": {"formaldehyde": true, "temperature": true, "humidity": true, "battery": true}}'),
('HCHO_002', '甲醛监测设备002', 'hcho', 31.2304, 121.4737, '上海市黄浦区', 'online', '{"report_interval": 60, "sensors": {"formaldehyde": true, "temperature": true, "humidity": true, "battery": true}}');

-- 插入示例告警规则
INSERT IGNORE INTO alert_rules (name, device_id, metric, condition_type, threshold_value, duration_seconds, severity, enabled, notification_channels) VALUES 
('PM2.5超标告警', NULL, 'pm25', 'gt', 75.0, 300, 'warning', TRUE, '["email", "sms"]'),
('PM10超标告警', NULL, 'pm10', 'gt', 150.0, 300, 'warning', TRUE, '["email", "sms"]'),
('CO2浓度告警', NULL, 'co2', 'gt', 1000.0, 600, 'critical', TRUE, '["email", "sms", "webhook"]'),
('甲醛浓度告警', NULL, 'formaldehyde', 'gt', 0.08, 300, 'critical', TRUE, '["email", "sms", "webhook"]'),
('温度异常告警', NULL, 'temperature', 'gt', 40.0, 180, 'warning', TRUE, '["email"]'),
('湿度异常告警', NULL, 'humidity', 'lt', 20.0, 300, 'info', TRUE, '["email"]'),
('电池电量低告警', NULL, 'battery', 'lt', 20.0, 600, 'warning', TRUE, '["email"]');

-- 创建数据迁移版本表
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入初始版本
INSERT IGNORE INTO schema_migrations (version) VALUES ('20240907000001');


-- 插入示例设备运行时状态
INSERT IGNORE INTO device_runtime_status (device_id, status, last_seen, uptime_seconds, memory_usage, cpu_usage, network_status, firmware_version) VALUES 
('ESP32_001', 'online', NOW(), 86400, 45.2, 12.8, 'connected', 'v1.2.3'),
('ESP32_002', 'online', NOW(), 172800, 38.7, 8.5, 'connected', 'v1.2.3');

-- 创建视图：设备实时状态
CREATE OR REPLACE VIEW device_realtime_status AS
SELECT 
    d.id,
    d.name,
    d.type,
    d.location_latitude,
    d.location_longitude,
    d.location_address,
    d.status,
    d.config,
    usd.timestamp as last_data_time,
    usd.pm25,
    usd.pm10,
    usd.co2,
    usd.formaldehyde,
    usd.temperature,
    usd.humidity,
    usd.pressure,
    usd.battery,
    usd.quality_score,
    CASE 
        WHEN TIMESTAMPDIFF(MINUTE, usd.timestamp, NOW()) > 10 THEN 'offline'
        ELSE d.status
    END as realtime_status
FROM devices d
LEFT JOIN (
    SELECT 
        device_id,
        timestamp,
        pm25,
        pm10,
        co2,
        formaldehyde,
        temperature,
        humidity,
        pressure,
        battery,
        quality_score,
        ROW_NUMBER() OVER (PARTITION BY device_id ORDER BY timestamp DESC) as rn
    FROM unified_sensor_data
) usd ON d.id = usd.device_id AND usd.rn = 1;

-- 创建视图：告警统计
CREATE OR REPLACE VIEW alert_statistics AS
SELECT 
    DATE(triggered_at) as alert_date,
    severity,
    status,
    COUNT(*) as alert_count
FROM alerts
WHERE triggered_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY DATE(triggered_at), severity, status
ORDER BY alert_date DESC, severity, status;
