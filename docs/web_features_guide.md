# Web功能开发指南

## 概述

本文档是空气质量监测系统Web功能的完整指南，包含数据查看、传感器字段管理、Web界面开发和功能实现的所有内容。

## 1. Web数据查看功能

### 1.1 功能特性

#### 1.1.1 数据查询
- **设备筛选**: 按设备ID筛选数据
- **设备类型筛选**: 按设备类型（hcho、esp32、sensor）筛选
- **传感器ID筛选**: 按具体传感器ID筛选
- **时间范围筛选**: 设置开始和结束时间查询历史数据
- **分页显示**: 支持分页浏览大量数据

#### 1.1.2 数据展示
- **表格展示**: 清晰的数据表格显示
- **状态标识**: 不同颜色标识数据状态
- **实时更新**: 支持自动刷新功能
- **统计信息**: 显示数据统计和分页信息

#### 1.1.3 数据导出
- **CSV格式**: 导出为CSV文件，便于Excel分析
- **JSON格式**: 导出为JSON文件，便于程序处理
- **筛选导出**: 支持按筛选条件导出数据

### 1.2 访问地址

- **数据查看页面**: http://127.0.0.1:8080/sensor-data
- **数据查询API**: http://127.0.0.1:8080/web/api/data
- **数据导出API**: http://127.0.0.1:8080/web/api/data/export

### 1.3 API接口

#### 1.3.1 数据查询API

**请求**: `GET /web/api/data`

**参数**:
- `device_id`: 设备ID（可选）
- `device_type`: 设备类型（可选）
- `sensor_id`: 传感器ID（可选）
- `start_time`: 开始时间，格式：2006-01-02T15:04（可选）
- `end_time`: 结束时间，格式：2006-01-02T15:04（可选）
- `page`: 页码，默认1（可选）
- `page_size`: 每页大小，默认20（可选）

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "device_id": "hcho_001",
      "device_type": "hcho",
      "sensor_id": "sensor_hcho_001_01",
      "sensor_type": "hcho",
      "timestamp": "2024-09-07T14:00:00Z",
      "pm25": 35.5,
      "pm10": 50.2,
      "co2": 400.0,
      "formaldehyde": 0.05,
      "temperature": 22.5,
      "humidity": 45.0,
      "pressure": 1013.25,
      "battery": 85,
      "data_quality": "good",
      "latitude": 39.9042,
      "longitude": 116.4074,
      "signal_strength": -65,
      "created_at": "2024-09-07T14:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 200,
    "page_size": 20
  },
  "filters": {
    "device_id": "hcho_001",
    "device_type": "hcho",
    "sensor_id": "",
    "start_time": "",
    "end_time": ""
  }
}
```

#### 1.3.2 数据导出API

**请求**: `GET /web/api/data/export`

**参数**:
- `device_id`: 设备ID（可选）
- `device_type`: 设备类型（可选）
- `sensor_id`: 传感器ID（可选）
- `start_time`: 开始时间（可选）
- `end_time`: 结束时间（可选）
- `format`: 导出格式，支持csv、json（默认csv）

**响应**: 直接下载文件

### 1.4 使用示例

#### 1.4.1 查询指定设备的数据
```
GET /web/api/data?device_id=hcho_001&page=1&page_size=10
```

#### 1.4.2 查询时间范围内的数据
```
GET /web/api/data?device_id=hcho_001&start_time=2024-09-07T00:00&end_time=2024-09-07T23:59
```

#### 1.4.3 导出CSV数据
```
GET /web/api/data/export?device_id=hcho_001&format=csv
```

#### 1.4.4 导出JSON数据
```
GET /web/api/data/export?device_id=hcho_001&format=json
```

## 2. 传感器字段管理

### 2.1 字段更新说明

本次更新在"甲醛数据消息"中增加了`sensor_id`和`sensor_type`字段，以支持更精细的传感器管理和数据追踪。

### 2.2 数据模型更新

#### 2.2.1 UnifiedSensorData模型
- 添加了`sensor_id`字段：`string`类型，用于标识具体的传感器
- 添加了`sensor_type`字段：`string`类型，用于标识传感器类型
- 为两个字段添加了数据库索引以提高查询性能

#### 2.2.2 MQTTMessage结构体
- 在MQTT消息结构中添加了`sensor_id`和`sensor_type`字段
- 保持了向后兼容性，字段为可选

#### 2.2.3 UnifiedSensorDataUpload结构体
- 在数据上传请求中添加了`sensor_id`和`sensor_type`字段

### 2.3 消息处理更新

#### 2.3.1 MQTT消息处理器
- 更新了`SensorDataHandler.HandleMessage`方法，解析并处理新的字段
- 在日志输出中包含了`sensor_id`和`sensor_type`信息
- 确保数据正确存储到数据库

### 2.4 消息格式示例

#### 2.4.1 更新后的甲醛数据消息格式

```json
{
  "device_id": "hcho_001",
  "device_type": "hcho",
  "sensor_id": "sensor_hcho_001_01",
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

## 3. 数据字段说明

### 3.1 完整字段列表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint64 | 数据ID |
| device_id | string | 设备ID |
| device_type | string | 设备类型 |
| sensor_id | string | 传感器ID |
| sensor_type | string | 传感器类型 |
| timestamp | time | 数据时间戳 |
| pm25 | float64 | PM2.5浓度 (μg/m³) |
| pm10 | float64 | PM10浓度 (μg/m³) |
| co2 | float64 | CO2浓度 (ppm) |
| formaldehyde | float64 | 甲醛浓度 (mg/m³) |
| temperature | float64 | 温度 (°C) |
| humidity | float64 | 湿度 (%) |
| pressure | float64 | 气压 (hPa) |
| battery | int | 电池电量 (%) |
| data_quality | string | 数据质量 |
| latitude | float64 | 纬度 |
| longitude | float64 | 经度 |
| signal_strength | int | 信号强度 (dBm) |
| created_at | time | 创建时间 |

### 3.2 状态标识

#### 3.2.1 数据质量
- **good**: 绿色 - 数据质量良好
- **fair**: 黄色 - 数据质量一般
- **poor**: 红色 - 数据质量较差

#### 3.2.2 甲醛浓度
- **正常**: 绿色 - < 0.08 mg/m³
- **超标**: 红色 - ≥ 0.08 mg/m³

#### 3.2.3 电池电量
- **正常**: 绿色 - ≥ 50%
- **警告**: 黄色 - 20-50%
- **危险**: 红色 - < 20%

#### 3.2.4 信号强度
- **良好**: 绿色 - > -70 dBm
- **一般**: 黄色 - -80 到 -70 dBm
- **较差**: 红色 - < -80 dBm

## 4. Web界面开发

### 4.1 页面结构

#### 4.1.1 导航栏
- **仪表板**: 系统概览和关键指标
- **设备管理**: 设备列表和状态管理
- **数据查看**: 历史数据查询和导出
- **图表分析**: 数据可视化展示
- **告警管理**: 告警规则和历史记录

#### 4.1.2 页面模板
- **base.html**: 基础模板，包含导航栏和页脚
- **dashboard.html**: 仪表板页面
- **devices.html**: 设备管理页面
- **data_view.html**: 数据查看页面
- **charts.html**: 图表分析页面
- **alerts.html**: 告警管理页面

### 4.2 模板系统

#### 4.2.1 模板继承
```html
<!-- base.html -->
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <title>{{template "title" .}} - 空气质量监测系统</title>
</head>
<body>
    <nav class="navbar">
        <!-- 导航链接 -->
    </nav>
    
    <main class="container-fluid mt-4">
        {{if eq .CurrentPage "dashboard"}}
            {{template "dashboard_content" .}}
        {{else if eq .CurrentPage "sensor-data"}}
            {{template "data_view_content" .}}
        {{end}}
    </main>
</body>
</html>
```

#### 4.2.2 页面模板
```html
<!-- data_view.html -->
{{define "title"}}数据查看{{end}}

{{define "data_view_content"}}
<div class="row">
    <div class="col-12">
        <h1 class="h3 mb-4">
            <i class="fas fa-database"></i> 数据查看
        </h1>
    </div>
</div>

<!-- 筛选表单 -->
<div class="row mb-4">
    <div class="col-12">
        <div class="card">
            <div class="card-body">
                <form id="filterForm" method="GET" action="/sensor-data">
                    <!-- 筛选字段 -->
                </form>
            </div>
        </div>
    </div>
</div>

<!-- 数据表格 -->
<div class="row">
    <div class="col-12">
        <div class="card">
            <div class="card-body">
                <div class="table-responsive">
                    <table class="table table-bordered" id="dataTable">
                        <!-- 表格内容 -->
                    </table>
                </div>
            </div>
        </div>
    </div>
</div>
{{end}}
```

### 4.3 模板函数

#### 4.3.1 自定义函数
```go
var TemplateFuncs = template.FuncMap{
    "buildQuery": buildQuery,
    "add":        add,
    "sub":        sub,
    "seq":        seq,
    "contains":   contains,
    "join":       strings.Join,
}

func buildQuery(params map[string]string) string {
    values := url.Values{}
    for k, v := range params {
        values.Set(k, v)
    }
    return values.Encode()
}

func add(a, b int) int {
    return a + b
}

func sub(a, b int) int {
    return a - b
}

func seq(start, end int) []int {
    result := make([]int, end-start+1)
    for i := range result {
        result[i] = start + i
    }
    return result
}
```

#### 4.3.2 分页功能
```html
<!-- 分页导航 -->
<nav aria-label="数据分页">
    <ul class="pagination justify-content-center">
        {{if gt .Pagination.CurrentPage 1}}
            <li class="page-item">
                <a class="page-link" href="?{{buildQuery (dict "page" (sub .Pagination.CurrentPage 1))}}">上一页</a>
            </li>
        {{end}}
        
        {{range $page := seq 1 .Pagination.TotalPages}}
            {{if eq $page $.Pagination.CurrentPage}}
                <li class="page-item active">
                    <span class="page-link">{{$page}}</span>
                </li>
            {{else}}
                <li class="page-item">
                    <a class="page-link" href="?{{buildQuery (dict "page" $page)}}">{{$page}}</a>
                </li>
            {{end}}
        {{end}}
        
        {{if lt .Pagination.CurrentPage .Pagination.TotalPages}}
            <li class="page-item">
                <a class="page-link" href="?{{buildQuery (dict "page" (add .Pagination.CurrentPage 1))}}">下一页</a>
            </li>
        {{end}}
    </ul>
</nav>
```

## 5. 路由和处理器

### 5.1 Web路由配置

```go
// web/routes.go
func SetupRoutes(router *gin.Engine, services *services.Services, logger utils.Logger) {
    webHandlers := handlers.NewWebHandlers(services, logger)

    // 页面路由
    router.GET("/", webHandlers.Dashboard)
    router.GET("/dashboard", webHandlers.Dashboard)
    router.GET("/devices", webHandlers.DeviceList)
    router.GET("/devices/:id", webHandlers.DeviceDetail)
    router.GET("/sensor-data", webHandlers.DataView)
    router.GET("/charts", webHandlers.Charts)
    router.GET("/alerts", webHandlers.Alerts)
    router.GET("/export", webHandlers.DataExportAPI)

    // Web API路由
    webAPI := router.Group("/web/api")
    {
        webAPI.GET("/device-stats", webHandlers.API)
        webAPI.GET("/latest-data", webHandlers.API)
        webAPI.GET("/chart-data", webHandlers.API)
        webAPI.GET("/data", webHandlers.DataAPI)
        webAPI.GET("/data/export", webHandlers.DataExportAPI)
    }
}
```

### 5.2 处理器实现

#### 5.2.1 数据查看处理器
```go
func (h *WebHandlers) DataView(c *gin.Context) {
    ctx := context.Background()

    // 获取查询参数
    deviceID := c.Query("device_id")
    deviceType := c.Query("device_type")
    sensorID := c.Query("sensor_id")
    startTime := c.Query("start_time")
    endTime := c.Query("end_time")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    // 查询数据
    var historyData []models.UnifiedSensorData
    var total int64
    var err error

    if startTime != "" && endTime != "" {
        // 时间范围查询
        start, err1 := time.Parse("2006-01-02T15:04", startTime)
        end, err2 := time.Parse("2006-01-02T15:04", endTime)
        
        if err1 != nil || err2 != nil {
            data := gin.H{
                "Title":   "数据查看",
                "Devices": devices,
                "Error":   "时间格式错误",
            }
            c.HTML(http.StatusOK, "base.html", data)
            return
        }
        
        historyData, total, err = h.services.UnifiedSensorData.GetDataByTimeRange(
            ctx, deviceID, start, end, pageSize, (page-1)*pageSize)
    } else {
        // 获取最新数据
        if deviceID != "" {
            historyData, err = h.services.UnifiedSensorData.GetDataByDeviceID(ctx, deviceID, pageSize)
        } else {
            historyData, err = h.services.UnifiedSensorData.GetLatestData(ctx, pageSize)
        }
        total = int64(len(historyData))
    }

    if err != nil {
        h.logger.Error("获取数据失败", utils.ErrorField(err))
        historyData = []models.UnifiedSensorData{}
    }

    // 计算分页信息
    totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

    data := gin.H{
        "Title":       "数据查看",
        "CurrentPage": "sensor-data",
        "Devices":     devices,
        "HistoryData": historyData,
        "Pagination": Pagination{
            CurrentPage: page,
            TotalPages:  int(totalPages),
            TotalItems:  int(total),
            PageSize:    pageSize,
        },
        "CurrentTime": time.Now().Format("2006-01-02 15:04:05"),
    }

    c.HTML(http.StatusOK, "base.html", data)
}
```

#### 5.2.2 数据API处理器
```go
func (h *WebHandlers) DataAPI(c *gin.Context) {
    ctx := context.Background()

    // 获取查询参数
    deviceID := c.Query("device_id")
    deviceType := c.Query("device_type")
    sensorID := c.Query("sensor_id")
    startTime := c.Query("start_time")
    endTime := c.Query("end_time")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    // 查询数据
    var historyData []models.UnifiedSensorData
    var total int64
    var err error

    // 实现查询逻辑...

    // 返回JSON响应
    response := gin.H{
        "data": historyData,
        "pagination": gin.H{
            "current_page": page,
            "total_pages":  int(totalPages),
            "total_items":  int(total),
            "page_size":    pageSize,
        },
        "filters": gin.H{
            "device_id":   deviceID,
            "device_type": deviceType,
            "sensor_id":   sensorID,
            "start_time":  startTime,
            "end_time":    endTime,
        },
    }

    c.JSON(http.StatusOK, response)
}
```

## 6. 测试和验证

### 6.1 测试脚本

使用 `scripts/test-web-data.bat` 脚本可以快速测试Web数据查看功能：

```bash
# Windows
scripts\test-web-data.bat
```

该脚本会：
1. 启动应用程序
2. 打开浏览器访问数据查看页面
3. 提供测试指导

### 6.2 功能验证

#### 6.2.1 数据查询验证
- 验证设备筛选功能
- 验证时间范围查询
- 验证分页功能
- 验证数据导出功能

#### 6.2.2 界面验证
- 验证页面加载
- 验证导航功能
- 验证响应式设计
- 验证状态标识

## 7. 部署和配置

### 7.1 数据库迁移

执行以下命令更新数据库结构：
```bash
mysql -u username -p database_name < scripts/add_sensor_fields_migration.sql
```

### 7.2 应用配置

#### 7.2.1 Web配置
```yaml
web:
  templates_path: "web/templates"
  static_path: "web/static"
  assets_path: "web/assets"
```

#### 7.2.2 路由配置
```yaml
routes:
  web:
    - path: "/"
      handler: "Dashboard"
    - path: "/sensor-data"
      handler: "DataView"
    - path: "/devices"
      handler: "DeviceList"
```

## 8. 注意事项

### 8.1 开发注意事项
1. **时间格式**: 时间参数使用 `2006-01-02T15:04` 格式
2. **分页限制**: 建议每页不超过100条记录
3. **导出限制**: 单次导出最多1000条记录
4. **数据权限**: 所有数据都是只读的，不支持修改
5. **性能优化**: 大量数据查询时建议使用时间范围筛选

### 8.2 向后兼容性
- 新的`sensor_id`和`sensor_type`字段为可选字段
- 如果消息中不包含这些字段，系统会使用默认值
- 现有的MQTT客户端无需立即更新，但建议逐步迁移到新格式

### 8.3 性能考虑
1. **数据库索引**: 新字段已添加索引，可能影响写入性能，但会显著提升查询性能
2. **存储空间**: 每个记录增加约100字节的存储空间
3. **日志输出**: 日志中会包含新的传感器信息，便于调试和监控

## 9. 后续计划

1. 更新Web界面以显示传感器信息
2. 添加基于传感器类型的告警规则
3. 实现传感器级别的数据统计和分析
4. 支持多传感器设备的独立管理
5. 添加数据可视化图表
6. 实现实时数据推送功能

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队
