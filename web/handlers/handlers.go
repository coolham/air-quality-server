package handlers

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// WebHandlers Web处理器集合
type WebHandlers struct {
	services *services.Services
	logger   utils.Logger
}

// NewWebHandlers 创建Web处理器
func NewWebHandlers(services *services.Services, logger utils.Logger) *WebHandlers {
	return &WebHandlers{
		services: services,
		logger:   logger,
	}
}

// Dashboard 仪表板页面
func (h *WebHandlers) Dashboard(c *gin.Context) {
	ctx := context.Background()

	// 获取设备统计信息
	deviceStats, err := h.getDeviceStats(ctx)
	if err != nil {
		h.logger.Error("获取设备统计失败", utils.ErrorField(err))
		deviceStats = &DeviceStats{}
	}

	// 获取最新数据
	latestData, err := h.getLatestData(ctx)
	if err != nil {
		h.logger.Error("获取最新数据失败", utils.ErrorField(err))
		latestData = []AirQualityDataSummary{}
	}

	// 获取告警统计
	alertStats, err := h.getAlertStats(ctx)
	if err != nil {
		h.logger.Error("获取告警统计失败", utils.ErrorField(err))
		alertStats = &AlertStats{}
	}

	data := gin.H{
		"Title":       "空气质量监测系统 - 仪表板",
		"CurrentPage": "dashboard",
		"DeviceStats": deviceStats,
		"LatestData":  latestData,
		"AlertStats":  alertStats,
		"CurrentTime": time.Now().Format("2006-01-02 15:04:05"),
	}

	c.HTML(http.StatusOK, "base.html", data)
}

// DeviceList 设备列表页面
func (h *WebHandlers) DeviceList(c *gin.Context) {
	h.logger.Info("🔍 DeviceList方法被调用", utils.String("path", c.Request.URL.Path), utils.String("method", c.Request.Method))
	ctx := context.Background()

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 获取设备列表
	devices, err := h.services.Device.ListDevices(ctx, pageSize, offset)
	if err != nil {
		h.logger.Error("获取设备列表失败", utils.ErrorField(err))
		devices = []models.Device{}
	}

	// 获取设备总数
	total, err := h.services.Device.CountDevices(ctx)
	if err != nil {
		h.logger.Error("获取设备总数失败", utils.ErrorField(err))
		total = 0
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	data := gin.H{
		"Title":       "设备管理",
		"CurrentPage": "devices",
		"Devices":     devices,
		"Pagination": Pagination{
			CurrentPage: page,
			TotalPages:  int(totalPages),
			TotalItems:  int(total),
			PageSize:    pageSize,
		},
	}

	c.HTML(http.StatusOK, "base.html", data)
}

// DeviceDetail 设备详情页面
func (h *WebHandlers) DeviceDetail(c *gin.Context) {
	deviceID := c.Param("id")
	ctx := context.Background()

	// 获取设备信息
	device, err := h.services.Device.GetDevice(ctx, deviceID)
	if err != nil {
		h.logger.Error("获取设备信息失败", utils.ErrorField(err))
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Title": "设备不存在",
			"Error": "设备不存在或已被删除",
		})
		return
	}

	// 获取设备最新数据
	latestData, err := h.services.AirQuality.GetLatestData(ctx, deviceID)
	if err != nil {
		h.logger.Error("获取设备最新数据失败", utils.ErrorField(err))
		latestData = nil
	}

	// 获取设备统计数据
	startTime := time.Now().Add(-24 * time.Hour).Unix()
	endTime := time.Now().Unix()
	stats, err := h.services.AirQuality.GetStatistics(ctx, deviceID, startTime, endTime)
	if err != nil {
		h.logger.Error("获取设备统计数据失败", utils.ErrorField(err))
		stats = nil
	}

	data := gin.H{
		"Title":      fmt.Sprintf("设备详情 - %s", device.Name),
		"Device":     device,
		"LatestData": latestData,
		"Stats":      stats,
	}

	c.HTML(http.StatusOK, "base.html", data)
}

// DataView 数据查看页面
func (h *WebHandlers) DataView(c *gin.Context) {
	h.logger.Info("🔍 DataView方法被调用", utils.String("path", c.Request.URL.Path), utils.String("method", c.Request.Method))
	ctx := context.Background()

	// 获取设备列表用于筛选
	devices, err := h.services.Device.ListDevices(ctx, 100, 0) // 获取前100个设备
	if err != nil {
		h.logger.Error("获取设备列表失败", utils.ErrorField(err))
		devices = []models.Device{}
	}

	// 获取查询参数
	deviceID := c.Query("device_id")
	deviceType := c.Query("device_type")
	sensorID := c.Query("sensor_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 构建筛选条件
	filters := map[string]interface{}{}
	if deviceID != "" {
		filters["device_id"] = deviceID
	}
	if deviceType != "" {
		filters["device_type"] = deviceType
	}
	if sensorID != "" {
		filters["sensor_id"] = sensorID
	}

	var data gin.H
	var historyData []models.UnifiedSensorData
	var total int64

	// 如果有时间范围参数，则查询数据
	if startTime != "" && endTime != "" {
		// 解析时间
		start, err1 := time.Parse("2006-01-02T15:04", startTime)
		end, err2 := time.Parse("2006-01-02T15:04", endTime)

		if err1 == nil && err2 == nil {
			// 获取指定时间范围的数据
			historyData, err = h.services.UnifiedSensorData.GetDataByTimeRange(ctx, deviceID, start.Unix(), end.Unix())
			if err != nil {
				h.logger.Error("获取历史数据失败", utils.ErrorField(err))
				historyData = []models.UnifiedSensorData{}
			}

			// 获取总数
			total = int64(len(historyData))
		} else {
			data = gin.H{
				"Title":   "数据查看",
				"Devices": devices,
				"Error":   "时间格式错误",
			}
			c.HTML(http.StatusOK, "base.html", data)
			return
		}
	} else {
		// 获取最新数据（分页）
		if deviceID != "" {
			// 获取指定设备的数据
			historyData, err = h.services.UnifiedSensorData.GetDataByDeviceID(ctx, deviceID, pageSize)
			if err != nil {
				h.logger.Error("获取设备数据失败", utils.ErrorField(err))
				historyData = []models.UnifiedSensorData{}
			}
			total = int64(len(historyData))
		} else {
			// 获取所有设备的最新数据（简化处理）
			historyData = []models.UnifiedSensorData{}
			total = 0
		}
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	data = gin.H{
		"Title":          "数据查看",
		"CurrentPage":    "sensor-data",
		"Devices":        devices,
		"SelectedDevice": deviceID,
		"DeviceType":     deviceType,
		"SensorID":       sensorID,
		"StartTime":      startTime,
		"EndTime":        endTime,
		"HistoryData":    historyData,
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

// Charts 图表页面
func (h *WebHandlers) Charts(c *gin.Context) {
	ctx := context.Background()

	// 获取设备列表
	devices, err := h.services.Device.ListDevices(ctx, 100, 0)
	if err != nil {
		h.logger.Error("获取设备列表失败", utils.ErrorField(err))
		devices = []models.Device{}
	}

	// 获取查询参数
	deviceID := c.Query("device_id")
	timeRange := c.DefaultQuery("time_range", "24") // 默认24小时

	var chartData *ChartData
	if deviceID != "" {
		// 根据时间范围获取数据
		hours, _ := strconv.Atoi(timeRange)
		endTime := time.Now()
		startTime := endTime.Add(-time.Duration(hours) * time.Hour)

		historyData, err := h.services.AirQuality.GetDataByTimeRange(ctx, deviceID, startTime.Unix(), endTime.Unix())
		if err != nil {
			h.logger.Error("获取图表数据失败", utils.ErrorField(err))
			historyData = []models.AirQualityData{}
		}

		// 转换为图表数据格式
		chartData = h.convertToChartData(historyData)
	}

	data := gin.H{
		"Title":          "数据图表",
		"CurrentPage":    "charts",
		"Devices":        devices,
		"SelectedDevice": deviceID,
		"TimeRange":      timeRange,
		"ChartData":      chartData,
	}

	c.HTML(http.StatusOK, "base.html", data)
}

// Alerts 告警管理页面
func (h *WebHandlers) Alerts(c *gin.Context) {
	ctx := context.Background()

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 获取告警列表
	alerts, err := h.services.Alert.ListAlerts(ctx, pageSize, offset)
	if err != nil {
		h.logger.Error("获取告警列表失败", utils.ErrorField(err))
		alerts = []models.Alert{}
	}

	// 获取告警总数
	total, err := h.services.Alert.CountAlerts(ctx)
	if err != nil {
		h.logger.Error("获取告警总数失败", utils.ErrorField(err))
		total = 0
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	data := gin.H{
		"Title":       "告警管理",
		"CurrentPage": "alerts",
		"Alerts":      alerts,
		"Pagination": Pagination{
			CurrentPage: page,
			TotalPages:  int(totalPages),
			TotalItems:  int(total),
			PageSize:    pageSize,
		},
	}

	c.HTML(http.StatusOK, "base.html", data)
}

// API 数据API接口
func (h *WebHandlers) API(c *gin.Context) {
	apiType := c.Param("type")
	ctx := context.Background()

	switch apiType {
	case "device-stats":
		stats, err := h.getDeviceStats(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)

	case "latest-data":
		deviceID := c.Query("device_id")
		if deviceID == "" {
			data, err := h.getLatestData(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, data)
		} else {
			data, err := h.services.AirQuality.GetLatestData(ctx, deviceID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, data)
		}

	case "chart-data":
		deviceID := c.Query("device_id")
		timeRange := c.DefaultQuery("time_range", "24")

		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
			return
		}

		hours, _ := strconv.Atoi(timeRange)
		endTime := time.Now()
		startTime := endTime.Add(-time.Duration(hours) * time.Hour)

		historyData, err := h.services.AirQuality.GetDataByTimeRange(ctx, deviceID, startTime.Unix(), endTime.Unix())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		chartData := h.convertToChartData(historyData)
		c.JSON(http.StatusOK, chartData)

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
	}
}

// DataAPI 数据查询API
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

	var data []models.UnifiedSensorData
	var total int64
	var err error

	// 如果有时间范围参数，则查询数据
	if startTime != "" && endTime != "" {
		// 解析时间
		start, err1 := time.Parse("2006-01-02T15:04", startTime)
		end, err2 := time.Parse("2006-01-02T15:04", endTime)

		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误"})
			return
		}

		// 获取指定时间范围的数据
		data, err = h.services.UnifiedSensorData.GetDataByTimeRange(ctx, deviceID, start.Unix(), end.Unix())
		if err != nil {
			h.logger.Error("获取历史数据失败", utils.ErrorField(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
			return
		}

		total = int64(len(data))
	} else {
		// 获取最新数据（分页）
		if deviceID != "" {
			// 获取指定设备的数据
			data, err = h.services.UnifiedSensorData.GetDataByDeviceID(ctx, deviceID, pageSize)
			if err != nil {
				h.logger.Error("获取设备数据失败", utils.ErrorField(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
				return
			}
			total = int64(len(data))
		} else {
			// 获取所有设备的最新数据（简化处理）
			data = []models.UnifiedSensorData{}
			total = 0
		}
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	response := gin.H{
		"data": data,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
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

// DataExportAPI 数据导出API
func (h *WebHandlers) DataExportAPI(c *gin.Context) {
	ctx := context.Background()

	// 获取查询参数
	deviceID := c.Query("device_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	format := c.DefaultQuery("format", "csv") // 支持csv, json

	var data []models.UnifiedSensorData
	var err error

	// 如果有时间范围参数，则查询数据
	if startTime != "" && endTime != "" {
		// 解析时间
		start, err1 := time.Parse("2006-01-02T15:04", startTime)
		end, err2 := time.Parse("2006-01-02T15:04", endTime)

		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误"})
			return
		}

		// 获取指定时间范围的数据
		data, err = h.services.UnifiedSensorData.GetDataByTimeRange(ctx, deviceID, start.Unix(), end.Unix())
		if err != nil {
			h.logger.Error("获取历史数据失败", utils.ErrorField(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
			return
		}
	} else {
		// 获取最新数据（限制1000条）
		if deviceID != "" {
			data, err = h.services.UnifiedSensorData.GetDataByDeviceID(ctx, deviceID, 1000)
			if err != nil {
				h.logger.Error("获取设备数据失败", utils.ErrorField(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
				return
			}
		} else {
			data = []models.UnifiedSensorData{}
		}
	}

	// 根据格式返回数据
	switch format {
	case "csv":
		csvData := h.convertToCSV(data)
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=sensor_data.csv")
		c.String(http.StatusOK, csvData)
	case "json":
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=sensor_data.json")
		c.JSON(http.StatusOK, data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的导出格式"})
	}
}

// getFloatValue 安全获取浮点数值
func getFloatValue(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}

// getDeviceStats 获取设备统计信息
func (h *WebHandlers) getDeviceStats(ctx context.Context) (*DeviceStats, error) {
	// 获取设备总数
	total, err := h.services.Device.CountDevices(ctx)
	if err != nil {
		return nil, err
	}

	// 获取在线设备数（这里简化处理，实际应该根据设备状态统计）
	onlineDevices := 0 // TODO: 实现在线设备统计
	offlineDevices := int(total) - onlineDevices
	activeDevices := onlineDevices // 假设在线设备都是活跃的

	return &DeviceStats{
		TotalDevices:   int(total),
		OnlineDevices:  onlineDevices,
		OfflineDevices: offlineDevices,
		ActiveDevices:  activeDevices,
	}, nil
}

// getLatestData 获取最新数据
func (h *WebHandlers) getLatestData(ctx context.Context) ([]AirQualityDataSummary, error) {
	// 获取所有设备
	devices, err := h.services.Device.ListDevices(ctx, 10, 0) // 获取前10个设备
	if err != nil {
		return nil, err
	}

	var summaries []AirQualityDataSummary
	for _, device := range devices {
		// 获取设备最新数据
		latestData, err := h.services.AirQuality.GetLatestData(ctx, device.ID)
		if err != nil {
			h.logger.Warn("获取设备最新数据失败", utils.String("device_id", device.ID), utils.ErrorField(err))
			continue
		}

		if latestData != nil {
			summary := AirQualityDataSummary{
				DeviceID:   device.ID,
				DeviceName: device.Name,
				PM25:       getFloatValue(latestData.PM25),
				PM10:       getFloatValue(latestData.PM10),
				Temp:       getFloatValue(latestData.Temperature),
				Humidity:   getFloatValue(latestData.Humidity),
				CreatedAt:  latestData.CreatedAt,
				Status:     string(device.Status),
			}
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

// getAlertStats 获取告警统计信息
func (h *WebHandlers) getAlertStats(ctx context.Context) (*AlertStats, error) {
	// 获取告警总数
	total, err := h.services.Alert.CountAlerts(ctx)
	if err != nil {
		return nil, err
	}

	// 获取未解决告警数
	unresolvedAlerts := 0 // TODO: 实现未解决告警统计
	criticalAlerts := 0   // TODO: 实现严重告警统计
	warningAlerts := 0    // TODO: 实现警告告警统计

	return &AlertStats{
		TotalAlerts:      int(total),
		UnresolvedAlerts: unresolvedAlerts,
		CriticalAlerts:   criticalAlerts,
		WarningAlerts:    warningAlerts,
	}, nil
}

// convertToChartData 将历史数据转换为图表数据格式
func (h *WebHandlers) convertToChartData(historyData []models.AirQualityData) *ChartData {
	var labels []string
	var pm25Data, pm10Data, tempData, humidityData []float64

	for _, data := range historyData {

		// 格式化时间标签
		label := data.CreatedAt.Format("15:04")
		labels = append(labels, label)

		// 添加数据点
		pm25Data = append(pm25Data, getFloatValue(data.PM25))
		pm10Data = append(pm10Data, getFloatValue(data.PM10))
		tempData = append(tempData, getFloatValue(data.Temperature))
		humidityData = append(humidityData, getFloatValue(data.Humidity))
	}

	return &ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:           "PM2.5",
				Data:            pm25Data,
				BorderColor:     "rgb(255, 99, 132)",
				BackgroundColor: "rgba(255, 99, 132, 0.2)",
				Fill:            false,
			},
			{
				Label:           "PM10",
				Data:            pm10Data,
				BorderColor:     "rgb(54, 162, 235)",
				BackgroundColor: "rgba(54, 162, 235, 0.2)",
				Fill:            false,
			},
			{
				Label:           "温度",
				Data:            tempData,
				BorderColor:     "rgb(255, 205, 86)",
				BackgroundColor: "rgba(255, 205, 86, 0.2)",
				Fill:            false,
			},
			{
				Label:           "湿度",
				Data:            humidityData,
				BorderColor:     "rgb(75, 192, 192)",
				BackgroundColor: "rgba(75, 192, 192, 0.2)",
				Fill:            false,
			},
		},
	}
}

// convertToCSV 将统一传感器数据转换为CSV格式
func (h *WebHandlers) convertToCSV(data []models.UnifiedSensorData) string {
	if len(data) == 0 {
		return ""
	}

	// CSV头部
	csv := "ID,设备ID,设备类型,传感器ID,传感器类型,时间戳,PM2.5,PM10,CO2,甲醛,温度,湿度,气压,电池,数据质量,纬度,经度,地址,质量评分,信号强度,创建时间\n"

	// 数据行
	for _, item := range data {
		csv += fmt.Sprintf("%d,%s,%s,%s,%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%s,%.8f,%.8f,%s,%.2f,%d,%s\n",
			item.ID,
			item.DeviceID,
			item.DeviceType,
			item.SensorID,
			item.SensorType,
			item.Timestamp.Format("2006-01-02 15:04:05"),
			getFloatValue(item.PM25),
			getFloatValue(item.PM10),
			getFloatValue(item.CO2),
			getFloatValue(item.Formaldehyde),
			getFloatValue(item.Temperature),
			getFloatValue(item.Humidity),
			getFloatValue(item.Pressure),
			float64(getIntValue(item.Battery)),
			item.DataQuality,
			getFloatValue(item.Latitude),
			getFloatValue(item.Longitude),
			"",  // Address字段不存在
			0.0, // QualityScore字段不存在
			getIntValue(item.SignalStrength),
			item.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	return csv
}

// getIntValue 安全获取整数值
func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
