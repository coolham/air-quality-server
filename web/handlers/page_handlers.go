package handlers

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

	// 获取查询参数
	deviceID := c.Query("device_id")
	deviceType := c.Query("device_type")
	sensorID := c.Query("sensor_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 获取设备列表用于筛选
	devices, err := h.services.Device.ListDevices(ctx, 100, 0) // 获取前100个设备
	if err != nil {
		h.logger.Error("获取设备列表失败", utils.ErrorField(err))
		devices = []models.Device{}
	}

	// 获取传感器ID列表
	var sensorIDs []string
	if deviceID != "" {
		// 如果选择了特定设备，获取该设备的传感器列表
		sensorIDs, err = h.services.UnifiedSensorData.GetSensorIDs(ctx, deviceID)
		if err != nil {
			h.logger.Error("获取设备传感器列表失败", utils.ErrorField(err))
			sensorIDs = []string{}
		}
	} else {
		// 如果没有选择设备，获取所有传感器
		sensorIDs, err = h.services.UnifiedSensorData.GetSensorIDs(ctx, "")
		if err != nil {
			h.logger.Error("获取传感器列表失败", utils.ErrorField(err))
			sensorIDs = []string{}
		}
	}

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
		// 解析时间，使用本地时区
		loc, _ := time.LoadLocation("Asia/Shanghai")
		start, err1 := time.ParseInLocation("2006-01-02T15:04", startTime, loc)
		end, err2 := time.ParseInLocation("2006-01-02T15:04", endTime, loc)

		if err1 == nil && err2 == nil {
			// 获取指定时间范围的数据
			historyData, err = h.services.UnifiedSensorData.GetDataByTimeRange(ctx, deviceID, start.Unix(), end.Unix())
			if err != nil {
				h.logger.Error("获取历史数据失败", utils.ErrorField(err))
				historyData = []models.UnifiedSensorData{}
			}

			// 如果指定了传感器ID，进行筛选
			if sensorID != "" {
				h.logger.Info("传感器筛选", utils.String("sensorID", sensorID), utils.Int("筛选前数据量", len(historyData)))
				var filteredData []models.UnifiedSensorData
				for _, data := range historyData {
					if data.SensorID == sensorID {
						filteredData = append(filteredData, data)
					}
				}
				historyData = filteredData
				h.logger.Info("传感器筛选结果", utils.String("sensorID", sensorID), utils.Int("筛选后数据量", len(historyData)))
			}

			// 获取总数
			total = int64(len(historyData))

			// 应用分页
			offset := (page - 1) * pageSize
			if offset < len(historyData) {
				end := offset + pageSize
				if end > len(historyData) {
					end = len(historyData)
				}
				historyData = historyData[offset:end]
			} else {
				historyData = []models.UnifiedSensorData{}
			}
		} else {
			data = gin.H{
				"Title":     "数据查看",
				"Devices":   devices,
				"SensorIDs": sensorIDs,
				"Error":     "时间格式错误",
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
			// 获取所有设备的最新数据
			offset := (page - 1) * pageSize
			historyData, err = h.services.UnifiedSensorData.GetAllData(ctx, pageSize, offset)
			if err != nil {
				h.logger.Error("获取所有设备数据失败", utils.ErrorField(err))
				historyData = []models.UnifiedSensorData{}
			}

			// 如果指定了传感器ID，进行筛选
			if sensorID != "" {
				var filteredData []models.UnifiedSensorData
				for _, data := range historyData {
					if data.SensorID == sensorID {
						filteredData = append(filteredData, data)
					}
				}
				historyData = filteredData
			}

			// 这里简化处理，实际应该有一个CountAllData方法
			total = int64(len(historyData))
		}
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	data = gin.H{
		"Title":          "数据查看",
		"CurrentPage":    "sensor-data",
		"Devices":        devices,
		"SensorIDs":      sensorIDs,
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
