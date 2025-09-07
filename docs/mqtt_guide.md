# MQTT功能设计与实现指南

## 概述

本文档是空气质量监测系统MQTT功能的完整指南，包含MQTT协议设计、数据存储、消息处理和系统实现的所有内容。

## 1. 功能目标与技术栈

### 1.1 功能目标
- 通过MQTT协议接收传感器设备上传的甲醛数据
- 支持多个传感器设备同时连接
- 实时处理和存储传感器数据
- 提供数据验证和错误处理机制
- 支持设备状态监控和告警

### 1.2 技术栈
- **MQTT Broker**: 嵌入式Mochi MQTT服务器
- **消息格式**: JSON
- **数据存储**: MySQL (通过GORM ORM)
- **编程语言**: Go

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────┐    MQTT     ┌─────────────────┐    HTTP/WebSocket    ┌─────────────────┐
│   ESP32设备     │ ──────────► │   MQTT Broker   │ ◄─────────────────► │   空气质量服务   │
│  (甲醛传感器)    │             │  (Mochi MQTT)   │                     │   (Go服务)      │
└─────────────────┘             └─────────────────┘                     └─────────────────┘
                                        │                                        │
                                        │                                        │
                                ┌─────────────────┐                     ┌─────────────────┐
                                │   Web管理界面   │                     │   MySQL数据库   │
                                │   (实时监控)    │                     │   (数据存储)    │
                                └─────────────────┘                     └─────────────────┘
```

### 2.2 组件说明

1. **ESP32设备**: 甲醛传感器，通过MQTT发布数据
2. **MQTT Broker**: 嵌入式Mochi MQTT服务器，负责消息路由
3. **空气质量服务**: Go服务，作为MQTT服务器接收数据
4. **MySQL数据库**: 存储传感器数据和设备信息
5. **Web界面**: 实时显示数据和设备状态

## 3. MQTT协议设计

### 3.1 连接配置

```yaml
mqtt:
  broker: "tcp://localhost:1883"  # MQTT Broker地址
  client_id: "air-quality-server" # 客户端ID
  keep_alive: 60                  # 心跳间隔(秒)
  clean_session: true             # 清理会话
  qos: 1                          # 服务质量等级
  auto_reconnect: true            # 自动重连
  connect_timeout: 30             # 连接超时(秒)
```

### 3.2 主题设计

#### 3.2.1 主题命名规范
```
air-quality/{device_type}/{device_id}/{data_type}
```

#### 3.2.2 具体主题列表

| 主题 | 描述 | 方向 | QoS |
|------|------|------|-----|
| `air-quality/hcho/{device_id}/data` | 甲醛传感器数据 | 设备→服务 | 1 |
| `air-quality/hcho/{device_id}/status` | 设备状态信息 | 设备→服务 | 1 |
| `air-quality/hcho/{device_id}/config` | 设备配置下发 | 服务→设备 | 1 |
| `air-quality/hcho/{device_id}/command` | 设备控制命令 | 服务→设备 | 1 |
| `air-quality/hcho/{device_id}/response` | 设备响应 | 设备→服务 | 1 |

### 3.3 消息格式

#### 3.3.1 传感器数据消息
```json
{
  "device_id": "hcho_001",
  "device_type": "hcho",
  "sensor_id": "sensor_01",
  "sensor_type": "hcho",
  "timestamp": 1694000000,
  "data": {
    "formaldehyde": 0.08,
    "temperature": 25.5,
    "humidity": 60.2,
    "battery": 85
  },
  "location": {
    "latitude": 39.9042,
    "longitude": 116.4074,
    "address": "北京市朝阳区"
  },
  "quality": {
    "signal_strength": -65,
    "data_quality": "good"
  }
}
```

#### 3.3.2 设备状态消息
```json
{
  "device_id": "hcho_001",
  "device_type": "hcho",
  "timestamp": 1694000000,
  "status": {
    "online": true,
    "battery_level": 85,
    "signal_strength": -65,
    "last_data_time": 1693999990,
    "error_code": 0,
    "error_message": ""
  },
  "firmware": {
    "version": "1.2.3",
    "build_date": "2024-01-15"
  }
}
```

#### 3.3.3 配置下发消息
```json
{
  "device_id": "hcho_001",
  "timestamp": 1694000000,
  "config": {
    "report_interval": 300,
    "thresholds": {
      "formaldehyde_warning": 0.08,
      "formaldehyde_critical": 0.1
    },
    "calibration": {
      "enabled": true,
      "interval": 86400
    }
  }
}
```

#### 3.3.4 控制命令消息
```json
{
  "device_id": "hcho_001",
  "timestamp": 1694000000,
  "command": {
    "action": "calibrate",
    "parameters": {
      "duration": 300
    }
  }
}
```

## 4. 技术实现

### 4.1 MQTT服务器架构

```go
type Server struct {
    config            *config.MQTTConfig
    logger            utils.Logger
    server            *mqtt.Server
    sensorDataHandler *SensorDataHandler
}
```

### 4.2 消息处理钩子

```go
type MessageHandlerHook struct {
    logger            utils.Logger
    sensorDataHandler *SensorDataHandler
}
```

### 4.3 数据处理器

```go
type SensorDataHandler struct {
    dataRepo   repositories.UnifiedSensorDataRepository
    deviceRepo repositories.DeviceRepository
    statusRepo repositories.DeviceRuntimeStatusRepository
    alertSvc   services.AlertService
    logger     utils.Logger
}
```

### 4.4 功能特性

#### 4.4.1 嵌入式MQTT服务器
- 使用Mochi MQTT作为嵌入式MQTT服务器
- 支持TCP连接（端口1883）
- 自动处理客户端连接和断开
- 支持QoS 0、1、2消息传递

#### 4.4.2 消息处理机制
- 自动识别传感器数据主题
- 支持主题格式：`air-quality/{device_type}/{device_id}/data`
- 支持设备类型：`hcho`、`esp32`、`sensor`
- 实时处理接收到的消息

#### 4.4.3 数据存储功能
- 自动解析JSON格式的传感器数据
- 存储到统一传感器数据表（`unified_sensor_data`）
- 支持多种传感器数据类型
- 自动更新设备运行时状态

#### 4.4.4 告警功能
- 实时检查传感器数据阈值
- 自动生成告警记录
- 支持甲醛浓度告警（阈值：0.08 mg/m³）
- 支持电池电量告警（阈值：20%）

## 5. 支持的数据格式

### 5.1 环境指标
- `formaldehyde`: 甲醛浓度 (mg/m³)
- `pm25`: PM2.5浓度 (μg/m³)
- `pm10`: PM10浓度 (μg/m³)
- `co2`: CO2浓度 (ppm)
- `temperature`: 温度 (°C)
- `humidity`: 湿度 (%)
- `pressure`: 气压 (hPa)

### 5.2 设备状态
- `battery`: 电池电量 (%)
- `signal_strength`: 信号强度 (dBm)
- `data_quality`: 数据质量

### 5.3 位置信息
- `latitude`: 纬度
- `longitude`: 经度
- `address`: 地址

## 6. 告警规则设计

### 6.1 甲醛浓度告警

| 等级 | 甲醛浓度范围 | 告警类型 | 处理方式 |
|------|-------------|----------|----------|
| 正常 | < 0.08 mg/m³ | 无 | 正常记录 |
| 警告 | 0.08 - 0.1 mg/m³ | 警告 | 发送通知 |
| 严重 | > 0.1 mg/m³ | 严重 | 立即告警 |

### 6.2 设备状态告警

| 告警类型 | 触发条件 | 处理方式 |
|----------|----------|----------|
| 设备离线 | 超过5分钟未收到心跳 | 发送离线通知 |
| 电池低电量 | 电池电量 < 20% | 发送低电量警告 |
| 信号弱 | 信号强度 < -80 dBm | 发送信号弱警告 |
| 数据异常 | 数据质量 = "poor" | 记录异常日志 |

### 6.3 其他告警规则
- PM2.5超标告警（阈值：75.0 μg/m³）
- PM10超标告警（阈值：150.0 μg/m³）
- CO2浓度告警（阈值：1000.0 ppm）
- 温度异常告警（阈值：40.0°C）
- 湿度异常告警（阈值：20%）

## 7. 功能模块设计

### 7.1 MQTT客户端模块

```go
type MQTTClient struct {
    client   mqtt.Client
    config   *MQTTConfig
    logger   utils.Logger
    handlers map[string]MessageHandler
}

type MessageHandler interface {
    HandleMessage(topic string, payload []byte) error
}
```

**主要功能**:
- 连接MQTT Broker
- 订阅主题
- 处理接收到的消息
- 发布消息到设备
- 自动重连机制
- 连接状态监控

### 7.2 消息处理器模块

```go
type FormaldehydeDataHandler struct {
    dataRepo   repositories.FormaldehydeDataRepository
    deviceRepo repositories.DeviceRepository
    alertSvc   services.AlertService
    logger     utils.Logger
}

type DeviceStatusHandler struct {
    statusRepo repositories.DeviceStatusRepository
    logger     utils.Logger
}
```

**主要功能**:
- 解析JSON消息
- 数据验证和清洗
- 存储到数据库
- 触发告警检查
- 更新设备状态

### 7.3 设备管理模块

```go
type DeviceManager struct {
    deviceRepo   repositories.DeviceRepository
    statusRepo   repositories.DeviceStatusRepository
    configRepo   repositories.DeviceConfigRepository
    mqttClient   *MQTTClient
    logger       utils.Logger
}
```

**主要功能**:
- 设备注册和认证
- 设备状态监控
- 配置下发
- 远程控制
- 设备离线检测

## 8. API接口设计

### 8.1 设备管理接口

```http
GET    /api/v1/devices/hcho           # 获取甲醛设备列表
GET    /api/v1/devices/hcho/{id}      # 获取设备详情
POST   /api/v1/devices/hcho           # 注册新设备
PUT    /api/v1/devices/hcho/{id}      # 更新设备信息
DELETE /api/v1/devices/hcho/{id}      # 删除设备
```

### 8.2 数据查询接口

```http
GET    /api/v1/data/hcho/{device_id}           # 获取设备数据
GET    /api/v1/data/hcho/{device_id}/realtime  # 获取实时数据
GET    /api/v1/data/hcho/{device_id}/history   # 获取历史数据
GET    /api/v1/data/hcho/{device_id}/export    # 导出数据
```

### 8.3 设备控制接口

```http
POST   /api/v1/devices/hcho/{id}/config    # 下发配置
POST   /api/v1/devices/hcho/{id}/command   # 发送命令
GET    /api/v1/devices/hcho/{id}/status    # 获取设备状态
```

## 9. 测试工具

### 9.1 快速测试脚本
```bash
# 运行快速测试
python tools/mqtt/quick_test.py

# 或使用批处理脚本
test-mqtt-storage.bat
```

### 9.2 完整测试脚本
```bash
# 运行完整测试
python tools/mqtt/test_data_storage.py
```

### 9.3 测试功能
- 基本传感器数据发布
- 多设备数据发布
- 告警条件测试
- 数据存储验证

## 10. 日志记录

### 10.1 连接日志
- 客户端连接/断开
- 认证状态
- 会话建立

### 10.2 消息日志
- 消息接收
- 消息内容解析
- 数据处理结果

### 10.3 存储日志
- 数据保存状态
- 设备状态更新
- 告警生成

## 11. 配置说明

### 11.1 MQTT配置
```yaml
mqtt:
  broker: "localhost:1883"
  client_id: "air-quality-server"
  keep_alive: 60
  clean_session: true
  qos: 1
  auto_reconnect: true
  connect_timeout: 30
```

### 11.2 主题配置
```yaml
topics:
  data: "air-quality/{device_type}/{device_id}/data"
  status: "air-quality/{device_type}/{device_id}/status"
  config: "air-quality/{device_type}/{device_id}/config"
```

## 12. 部署配置

### 12.1 Docker Compose配置

```yaml
version: '3.8'
services:
  mosquitto:
    image: eclipse-mosquitto:2.0
    container_name: mosquitto
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto/config:/mosquitto/config
      - ./mosquitto/data:/mosquitto/data
      - ./mosquitto/log:/mosquitto/log
    networks:
      - air-quality-network

  air-quality-server:
    build: .
    container_name: air-quality-server
    ports:
      - "8080:8080"
    environment:
      - MQTT_BROKER=mosquitto:1883
      - MQTT_USERNAME=admin
      - MQTT_PASSWORD=password
    depends_on:
      - mysql
      - redis
      - mosquitto
    networks:
      - air-quality-network

networks:
  air-quality-network:
    driver: bridge
```

### 12.2 Mosquitto配置

```conf
# mosquitto.conf
listener 1883
allow_anonymous false
password_file /mosquitto/config/passwd

listener 9001
protocol websockets
allow_anonymous false
password_file /mosquitto/config/passwd

log_dest file /mosquitto/log/mosquitto.log
log_type error
log_type warning
log_type notice
log_type information
```

## 13. 故障排除

### 13.1 常见问题

#### MQTT服务器无法启动
- 检查端口1883是否被占用
- 确认防火墙设置
- 检查配置文件

#### 消息无法接收
- 确认主题格式正确
- 检查客户端连接状态
- 查看服务器日志

#### 数据存储失败
- 检查数据库连接
- 确认表结构正确
- 查看错误日志

### 13.2 调试方法

#### 启用详细日志
```yaml
log:
  level: "debug"
  format: "json"
```

#### 检查数据库
```sql
-- 查看传感器数据
SELECT * FROM unified_sensor_data ORDER BY created_at DESC LIMIT 10;

-- 查看设备状态
SELECT * FROM device_runtime_status;

-- 查看告警记录
SELECT * FROM alerts ORDER BY triggered_at DESC LIMIT 10;
```

## 14. 性能优化

### 14.1 数据库优化
- 添加适当的索引
- 定期清理历史数据
- 使用连接池

### 14.2 MQTT优化
- 合理设置QoS级别
- 控制消息频率
- 使用持久会话

### 14.3 内存优化
- 及时释放资源
- 控制日志级别
- 优化数据结构

## 15. 安全考虑

### 15.1 认证机制
- 当前使用AllowHook（允许所有连接）
- 生产环境应实现自定义认证
- 支持用户名/密码认证

### 15.2 数据验证
- 验证消息格式
- 检查数据范围
- 防止SQL注入

### 15.3 网络安全
- 使用TLS加密
- 限制客户端IP
- 监控异常连接

## 16. 监控和运维

### 16.1 监控指标
- MQTT连接状态
- 消息吞吐量
- 设备在线率
- 数据延迟
- 错误率

### 16.2 日志管理
- 结构化日志
- 日志轮转
- 集中收集
- 告警通知

### 16.3 性能优化
- 连接池管理
- 消息批处理
- 数据库索引优化
- 缓存策略

## 17. 扩展功能

### 17.1 支持更多设备类型
- 添加新的设备类型支持
- 扩展数据字段
- 自定义处理逻辑

### 17.2 实时数据推送
- WebSocket支持
- 实时数据流
- 客户端订阅

### 17.3 数据分析
- 数据聚合
- 趋势分析
- 预测模型

## 18. 实施计划

### 阶段1: 基础功能 (1-2周)
- [x] MQTT服务器实现
- [x] 基础消息处理
- [x] 数据模型创建
- [x] 数据库迁移

### 阶段2: 核心功能 (2-3周)
- [x] 设备管理功能
- [x] 告警系统集成
- [x] API接口实现
- [x] Web界面集成

### 阶段3: 高级功能 (1-2周)
- [ ] 设备配置下发
- [ ] 远程控制功能
- [ ] 监控和运维
- [ ] 性能优化

### 阶段4: 测试和部署 (1周)
- [x] 全面测试
- [x] 文档完善
- [x] 生产部署
- [ ] 用户培训

## 19. 风险评估

### 19.1 技术风险
- MQTT Broker稳定性
- 网络连接可靠性
- 数据一致性
- 性能瓶颈

### 19.2 业务风险
- 设备兼容性
- 数据准确性
- 用户体验
- 运维复杂度

### 19.3 缓解措施
- 高可用部署
- 数据备份策略
- 监控告警机制
- 文档和培训

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队
