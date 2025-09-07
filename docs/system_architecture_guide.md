# 系统架构设计指南

## 概述

本文档是空气质量监测服务端系统的完整架构设计指南，包含系统概述、模块设计、接口定义和部署架构的所有内容。

## 1. 系统概述

### 1.1 项目背景
企业级空气质量监测服务端系统，用于接收、存储、处理和分析来自ESP32空气质量监测设备的数据。系统需要支持大规模设备接入、实时数据处理、历史数据分析和可视化展示。

### 1.2 核心需求
- **高并发**: 支持数千台ESP32设备同时接入
- **实时性**: 数据接收和处理延迟小于100ms
- **可靠性**: 99.9%系统可用性，数据不丢失
- **可扩展**: 支持水平扩展，模块化设计
- **监控**: 完善的系统监控和告警机制

## 2. 系统架构设计

### 2.1 整体架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ESP32设备     │    │   移动端APP     │    │   Web管理端     │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTTP/WebSocket       │ HTTP/WebSocket       │ HTTP/WebSocket
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
┌───────▼────────┐    ┌─────────▼─────────┐    ┌─────────▼─────────┐
│  数据接收服务   │    │   数据处理服务     │    │   数据查询服务     │
│ (Data Ingestion)│    │ (Data Processing) │    │ (Data Query)      │
└───────┬────────┘    └─────────┬─────────┘    └─────────┬─────────┘
        │                      │                        │
        │                      │                        │
┌───────▼────────┐    ┌─────────▼─────────┐    ┌─────────▼─────────┐
│   消息队列      │    │   关系数据库      │    │   Redis缓存       │
│  (Redis Pub/Sub)│    │   (MySQL)         │    │   (Redis)         │
└────────────────┘    └───────────────────┘    └───────────────────┘
```

### 2.2 微服务模块划分

#### 2.2.1 核心服务模块
1. **设备管理服务** (Device Management)
2. **数据接收服务** (Data Ingestion)
3. **数据处理服务** (Data Processing)
4. **数据查询服务** (Data Query)
5. **告警服务** (Alert Service)
6. **用户管理服务** (User Management)
7. **配置管理服务** (Configuration Management)

#### 2.2.2 支撑服务模块
1. **认证授权服务** (Auth Service)
2. **日志服务** (Logging Service)
3. **监控服务** (Monitoring Service)
4. **通知服务** (Notification Service)

## 3. 详细模块设计

### 3.1 设备管理服务 (Device Management)
**职责**: 设备注册、状态管理、配置下发
**技术栈**: Go + gRPC + MySQL
**功能**:
- 设备注册和认证
- 设备状态监控
- 设备配置管理
- 设备分组管理
- 设备固件升级

**数据模型**:
```go
type Device struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Location    Location  `json:"location"`
    Status      string    `json:"status"`
    Config      Config    `json:"config"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Address   string  `json:"address"`
}
```

### 3.2 数据接收服务 (Data Ingestion)
**职责**: 接收设备数据、数据验证、格式转换
**技术栈**: Go + Redis Pub/Sub + MySQL
**功能**:
- HTTP/WebSocket数据接收
- 数据格式验证和清洗
- 数据压缩和批量处理
- 数据路由到消息队列

**数据模型**:
```go
type AirQualityData struct {
    DeviceID    string    `json:"device_id"`
    Timestamp   time.Time `json:"timestamp"`
    PM25        float64   `json:"pm25"`
    PM10        float64   `json:"pm10"`
    CO2         float64   `json:"co2"`
    Temperature float64   `json:"temperature"`
    Humidity    float64   `json:"humidity"`
    Pressure    float64   `json:"pressure"`
    Location    Location  `json:"location"`
}
```

### 3.3 数据处理服务 (Data Processing)
**职责**: 实时数据处理、统计分析、异常检测
**技术栈**: Go + Redis Pub/Sub + MySQL + Redis
**功能**:
- 实时数据流处理
- 数据聚合和统计
- 异常检测和告警触发
- 数据质量评估
- 数据补全和修正

### 3.4 数据查询服务 (Data Query)
**职责**: 历史数据查询、报表生成、数据导出
**技术栈**: Go + MySQL + Redis
**功能**:
- 多维度数据查询
- 实时数据API
- 历史数据报表
- 数据可视化接口
- 数据导出功能

### 3.5 告警服务 (Alert Service)
**职责**: 告警规则管理、告警触发、通知发送
**技术栈**: Go + MySQL + Redis + Redis Pub/Sub
**功能**:
- 告警规则配置
- 实时告警检测
- 告警级别管理
- 告警通知发送
- 告警历史记录

### 3.6 用户管理服务 (User Management)
**职责**: 用户认证、权限管理、组织架构
**技术栈**: Go + MySQL + Redis
**功能**:
- 用户注册和登录
- 角色权限管理
- 组织架构管理
- 用户行为审计

### 3.7 配置管理服务 (Configuration Management)
**职责**: 系统配置、设备配置、业务配置
**技术栈**: Go + MySQL + Consul
**功能**:
- 系统参数配置
- 设备配置模板
- 业务规则配置
- 配置版本管理

## 4. 模块间通信协议

### 4.1 通信方式
- **同步通信**: gRPC (服务间调用)
- **异步通信**: Redis Pub/Sub (事件驱动)
- **HTTP API**: 对外接口

### 4.2 数据格式
- **JSON**: HTTP API 数据格式
- **Protocol Buffers**: gRPC 数据格式
- **JSON**: Redis Pub/Sub 消息格式

## 5. 核心服务接口定义

### 5.1 设备管理服务 (Device Management Service)

#### 5.1.1 gRPC 接口
```protobuf
service DeviceManagementService {
    // 设备注册
    rpc RegisterDevice(RegisterDeviceRequest) returns (RegisterDeviceResponse);
    
    // 设备认证
    rpc AuthenticateDevice(AuthenticateDeviceRequest) returns (AuthenticateDeviceResponse);
    
    // 获取设备信息
    rpc GetDevice(GetDeviceRequest) returns (GetDeviceResponse);
    
    // 更新设备状态
    rpc UpdateDeviceStatus(UpdateDeviceStatusRequest) returns (UpdateDeviceStatusResponse);
    
    // 获取设备配置
    rpc GetDeviceConfig(GetDeviceConfigRequest) returns (GetDeviceConfigResponse);
    
    // 更新设备配置
    rpc UpdateDeviceConfig(UpdateDeviceConfigRequest) returns (UpdateDeviceConfigResponse);
    
    // 设备列表查询
    rpc ListDevices(ListDevicesRequest) returns (ListDevicesResponse);
}

message Device {
    string id = 1;
    string name = 2;
    string type = 3;
    Location location = 4;
    string status = 5;
    DeviceConfig config = 6;
    int64 created_at = 7;
    int64 updated_at = 8;
}

message Location {
    double latitude = 1;
    double longitude = 2;
    string address = 3;
}

message DeviceConfig {
    int32 report_interval = 1;  // 上报间隔(秒)
    map<string, string> sensors = 2;  // 传感器配置
    map<string, string> thresholds = 3;  // 阈值配置
}
```

#### 5.1.2 HTTP API 接口
```go
// 设备管理 REST API
GET    /api/v1/devices                    // 获取设备列表
POST   /api/v1/devices                    // 创建设备
GET    /api/v1/devices/{id}               // 获取设备详情
PUT    /api/v1/devices/{id}               // 更新设备信息
DELETE /api/v1/devices/{id}               // 删除设备
GET    /api/v1/devices/{id}/config        // 获取设备配置
PUT    /api/v1/devices/{id}/config        // 更新设备配置
POST   /api/v1/devices/{id}/status        // 更新设备状态
```

### 5.2 数据接收服务 (Data Ingestion Service)

#### 5.2.1 HTTP API 接口
```go
// 数据上报接口
POST   /api/v1/data/upload                // 设备数据上报
POST   /api/v1/data/batch                 // 批量数据上报
GET    /api/v1/data/status                // 数据接收状态

// WebSocket 接口
WS     /ws/data                           // 实时数据推送
```

#### 5.2.2 数据模型
```go
type DataUploadRequest struct {
    DeviceID  string                 `json:"device_id" binding:"required"`
    Timestamp int64                  `json:"timestamp" binding:"required"`
    Data      AirQualityData         `json:"data" binding:"required"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type AirQualityData struct {
    PM25        *float64 `json:"pm25,omitempty"`
    PM10        *float64 `json:"pm10,omitempty"`
    CO2         *float64 `json:"co2,omitempty"`
    Temperature *float64 `json:"temperature,omitempty"`
    Humidity    *float64 `json:"humidity,omitempty"`
    Pressure    *float64 `json:"pressure,omitempty"`
    Location    Location `json:"location,omitempty"`
}
```

### 5.3 数据处理服务 (Data Processing Service)

#### 5.3.1 Redis Pub/Sub 消息格式
```go
// 原始数据消息
type RawDataMessage struct {
    MessageID   string          `json:"message_id"`
    DeviceID    string          `json:"device_id"`
    Timestamp   int64           `json:"timestamp"`
    Data        AirQualityData  `json:"data"`
    ReceivedAt  int64           `json:"received_at"`
}

// 处理结果消息
type ProcessedDataMessage struct {
    MessageID     string                 `json:"message_id"`
    DeviceID      string                 `json:"device_id"`
    Timestamp     int64                  `json:"timestamp"`
    OriginalData  AirQualityData         `json:"original_data"`
    ProcessedData AirQualityData         `json:"processed_data"`
    Statistics    DataStatistics         `json:"statistics"`
    Alerts        []Alert                `json:"alerts,omitempty"`
    QualityScore  float64                `json:"quality_score"`
}

type DataStatistics struct {
    PM25Avg    float64 `json:"pm25_avg"`
    PM10Avg    float64 `json:"pm10_avg"`
    CO2Avg     float64 `json:"co2_avg"`
    TempAvg    float64 `json:"temp_avg"`
    HumidityAvg float64 `json:"humidity_avg"`
    PressureAvg float64 `json:"pressure_avg"`
}
```

### 5.4 数据查询服务 (Data Query Service)

#### 5.4.1 HTTP API 接口
```go
// 实时数据查询
GET    /api/v1/data/realtime/{device_id}     // 获取设备实时数据
GET    /api/v1/data/realtime                 // 获取所有设备实时数据

// 历史数据查询
GET    /api/v1/data/history/{device_id}      // 获取设备历史数据
GET    /api/v1/data/statistics/{device_id}   // 获取设备统计数据
GET    /api/v1/data/export/{device_id}       // 导出设备数据

// 聚合数据查询
GET    /api/v1/data/aggregate                // 获取聚合数据
GET    /api/v1/data/trends                   // 获取趋势数据
```

#### 5.4.2 查询参数
```go
type DataQueryParams struct {
    DeviceIDs   []string  `form:"device_ids"`
    StartTime   int64     `form:"start_time"`
    EndTime     int64     `form:"end_time"`
    Interval    string    `form:"interval"`    // 1m, 5m, 1h, 1d
    Metrics     []string  `form:"metrics"`     // pm25, pm10, co2, etc.
    Aggregation string    `form:"aggregation"` // avg, max, min, sum
    Limit       int       `form:"limit"`
    Offset      int       `form:"offset"`
}
```

### 5.5 告警服务 (Alert Service)

#### 5.5.1 gRPC 接口
```protobuf
service AlertService {
    // 创建告警规则
    rpc CreateAlertRule(CreateAlertRuleRequest) returns (CreateAlertRuleResponse);
    
    // 更新告警规则
    rpc UpdateAlertRule(UpdateAlertRuleRequest) returns (UpdateAlertRuleResponse);
    
    // 删除告警规则
    rpc DeleteAlertRule(DeleteAlertRuleRequest) returns (DeleteAlertRuleResponse);
    
    // 获取告警规则列表
    rpc ListAlertRules(ListAlertRulesRequest) returns (ListAlertRulesResponse);
    
    // 获取告警历史
    rpc GetAlertHistory(GetAlertHistoryRequest) returns (GetAlertHistoryResponse);
    
    // 处理告警
    rpc ProcessAlert(ProcessAlertRequest) returns (ProcessAlertResponse);
}

message AlertRule {
    string id = 1;
    string name = 2;
    string device_id = 3;
    string metric = 4;
    string condition = 5;  // gt, lt, eq, ne
    double threshold = 6;
    int32 duration = 7;    // 持续时间(秒)
    string severity = 8;   // critical, warning, info
    bool enabled = 9;
    repeated string notification_channels = 10;
}
```

#### 5.5.2 HTTP API 接口
```go
// 告警管理
GET    /api/v1/alerts/rules                 // 获取告警规则列表
POST   /api/v1/alerts/rules                 // 创建告警规则
PUT    /api/v1/alerts/rules/{id}            // 更新告警规则
DELETE /api/v1/alerts/rules/{id}            // 删除告警规则

// 告警历史
GET    /api/v1/alerts/history               // 获取告警历史
GET    /api/v1/alerts/active                // 获取活跃告警
POST   /api/v1/alerts/{id}/acknowledge      // 确认告警
POST   /api/v1/alerts/{id}/resolve          // 解决告警
```

### 5.6 用户管理服务 (User Management Service)

#### 5.6.1 gRPC 接口
```protobuf
service UserManagementService {
    // 用户认证
    rpc AuthenticateUser(AuthenticateUserRequest) returns (AuthenticateUserResponse);
    
    // 用户注册
    rpc RegisterUser(RegisterUserRequest) returns (RegisterUserResponse);
    
    // 获取用户信息
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    
    // 更新用户信息
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    
    // 用户列表
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    
    // 角色管理
    rpc CreateRole(CreateRoleRequest) returns (CreateRoleResponse);
    rpc AssignRole(AssignRoleRequest) returns (AssignRoleResponse);
}

message User {
    string id = 1;
    string username = 2;
    string email = 3;
    string phone = 4;
    string status = 5;
    repeated string roles = 6;
    int64 created_at = 7;
    int64 updated_at = 8;
}
```

#### 5.6.2 HTTP API 接口
```go
// 用户认证
POST   /api/v1/auth/login                   // 用户登录
POST   /api/v1/auth/logout                  // 用户登出
POST   /api/v1/auth/refresh                 // 刷新Token
POST   /api/v1/auth/register                // 用户注册

// 用户管理
GET    /api/v1/users                        // 获取用户列表
POST   /api/v1/users                        // 创建用户
GET    /api/v1/users/{id}                   // 获取用户详情
PUT    /api/v1/users/{id}                   // 更新用户信息
DELETE /api/v1/users/{id}                   // 删除用户

// 权限管理
GET    /api/v1/roles                        // 获取角色列表
POST   /api/v1/roles                        // 创建角色
PUT    /api/v1/roles/{id}                   // 更新角色
DELETE /api/v1/roles/{id}                   // 删除角色
```

## 6. 事件驱动架构

### 6.1 事件类型定义
```go
// 设备事件
type DeviceEvent struct {
    EventType string    `json:"event_type"` // registered, online, offline, config_updated
    DeviceID  string    `json:"device_id"`
    Timestamp int64     `json:"timestamp"`
    Data      interface{} `json:"data"`
}

// 数据事件
type DataEvent struct {
    EventType string    `json:"event_type"` // data_received, data_processed, alert_triggered
    DeviceID  string    `json:"device_id"`
    Timestamp int64     `json:"timestamp"`
    Data      interface{} `json:"data"`
}

// 告警事件
type AlertEvent struct {
    EventType string    `json:"event_type"` // alert_triggered, alert_resolved, alert_acknowledged
    AlertID   string    `json:"alert_id"`
    DeviceID  string    `json:"device_id"`
    Timestamp int64     `json:"timestamp"`
    Data      interface{} `json:"data"`
}
```

### 6.2 Redis Pub/Sub Channel 设计
```
# 原始数据
air-quality:raw-data

# 处理后的数据
air-quality:processed-data

# 告警事件
air-quality:alerts

# 设备事件
air-quality:device-events

# 系统事件
air-quality:system-events
```

## 7. 技术栈选择

### 7.1 后端技术栈
- **语言**: Go 1.21+
- **框架**: Gin (HTTP) + gRPC
- **数据库**: MySQL 8.0 + Redis 7.0
- **消息队列**: Redis Pub/Sub
- **服务发现**: Consul (可选)
- **配置中心**: 配置文件 + 环境变量
- **监控**: Prometheus + Grafana
- **日志**: 结构化日志 (JSON格式)

### 7.2 部署和运维
- **容器化**: Docker + Docker Compose
- **编排**: Kubernetes
- **CI/CD**: GitLab CI/CD
- **监控**: Prometheus + Grafana + Jaeger
- **日志**: ELK Stack
- **备份**: Velero + Restic

## 8. 数据存储设计

### 8.1 数据库选型
- **MySQL**: 存储所有数据（设备信息、用户信息、配置信息、空气质量数据）
- **Redis**: 缓存、会话存储、实时数据、消息队列

### 8.2 数据分片策略
- **水平分片**: 按设备ID或时间范围分片
- **读写分离**: 主从复制，读写分离
- **数据归档**: 历史数据定期归档到对象存储

## 9. 安全设计

### 9.1 认证授权
- **JWT Token**: 无状态认证
- **OAuth 2.0**: 第三方认证
- **RBAC**: 基于角色的访问控制

### 9.2 数据安全
- **传输加密**: TLS 1.3
- **存储加密**: 数据库字段加密
- **API安全**: 限流、防重放攻击

## 10. 性能优化

### 10.1 缓存策略
- **多级缓存**: 应用缓存 + Redis + CDN
- **缓存预热**: 热点数据预加载
- **缓存更新**: 异步更新策略

### 10.2 数据库优化
- **索引优化**: 合理设计索引
- **查询优化**: SQL优化、分页查询
- **连接池**: 数据库连接池管理

## 11. 监控和运维

### 11.1 系统监控
- **指标监控**: Prometheus + Grafana
- **日志监控**: ELK Stack
- **链路追踪**: Jaeger
- **健康检查**: 服务健康状态监控

### 11.2 告警机制
- **系统告警**: CPU、内存、磁盘、网络
- **业务告警**: 数据异常、服务异常
- **通知方式**: 邮件、短信、钉钉、企业微信

## 12. 扩展性设计

### 12.1 水平扩展
- **无状态服务**: 所有服务设计为无状态
- **负载均衡**: 多实例部署
- **数据分片**: 支持数据水平分片

### 12.2 模块化设计
- **微服务架构**: 服务独立部署
- **接口标准化**: RESTful API + gRPC
- **配置外部化**: 配置中心管理

## 13. 错误处理

### 13.1 错误码定义
```go
const (
    // 成功
    CodeSuccess = 0
    
    // 通用错误 (1000-1999)
    CodeInvalidRequest = 1001
    CodeUnauthorized   = 1002
    CodeForbidden      = 1003
    CodeNotFound       = 1004
    CodeInternalError  = 1005
    
    // 设备相关错误 (2000-2999)
    CodeDeviceNotFound     = 2001
    CodeDeviceOffline      = 2002
    CodeDeviceUnauthorized = 2003
    CodeDeviceConfigError  = 2004
    
    // 数据相关错误 (3000-3999)
    CodeDataInvalid        = 3001
    CodeDataNotFound       = 3002
    CodeDataProcessingError = 3003
    
    // 告警相关错误 (4000-4999)
    CodeAlertRuleNotFound  = 4001
    CodeAlertRuleInvalid   = 4002
    CodeAlertNotFound      = 4003
)
```

### 13.2 错误响应格式
```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
    RequestID string `json:"request_id"`
}
```

## 14. 接口版本管理

### 14.1 版本策略
- **URL版本**: `/api/v1/`, `/api/v2/`
- **Header版本**: `API-Version: v1`
- **向后兼容**: 保持向后兼容至少2个版本

### 14.2 版本生命周期
- **v1**: 当前稳定版本
- **v2**: 开发中版本
- **v0**: 实验性版本

## 15. 部署架构

### 15.1 开发环境
- Docker Compose 本地部署
- 单机多容器部署

### 15.2 生产环境
- Kubernetes 集群部署
- 多可用区部署
- 自动扩缩容

## 16. 开发计划

### 16.1 第一阶段 (MVP)
- 基础架构搭建
- 设备管理服务
- 数据接收服务
- 基础数据查询

### 16.2 第二阶段
- 数据处理服务
- 告警服务
- 用户管理服务
- Web管理界面

### 16.3 第三阶段
- 高级分析功能
- 移动端APP
- 第三方集成
- 性能优化

## 17. 风险评估

### 17.1 技术风险
- **数据一致性**: 分布式事务处理
- **性能瓶颈**: 高并发处理能力
- **数据丢失**: 消息队列可靠性

### 17.2 业务风险
- **设备兼容性**: 不同厂商设备适配
- **数据准确性**: 传感器数据校准
- **用户接受度**: 界面易用性

## 18. 总结

本设计方案采用微服务架构，具备高可用、高并发、可扩展的特性。通过模块化设计，系统可以灵活应对业务需求变化，支持快速迭代和功能扩展。技术栈选择成熟稳定，降低开发风险，提高系统可靠性。

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队
