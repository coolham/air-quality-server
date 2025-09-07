package handlers

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// API 数据API接口
func (h *WebHandlers) API(c *gin.Context) {
	// 从URL路径中提取API类型
	path := c.Request.URL.Path
	var apiType string

	// 添加调试日志
	h.logger.Info("API请求",
		utils.String("path", path),
		utils.String("method", c.Request.Method))

	// 根据路径确定API类型
	switch {
	case path == "/web/api/device-stats":
		apiType = "device-stats"
	case path == "/web/api/latest-data":
		apiType = "latest-data"
	case path == "/web/api/chart-data":
		apiType = "chart-data"
	case path == "/web/api/sensors":
		apiType = "sensors"
	default:
		apiType = c.Param("type")
		h.logger.Warn("未知API类型",
			utils.String("path", path),
			utils.String("apiType", apiType))
	}

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
		sensorID := c.Query("sensor_id")
		timeRange := c.DefaultQuery("time_range", "24")
		metric := c.DefaultQuery("metric", "all")

		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
			return
		}

		hours, _ := strconv.Atoi(timeRange)
		endTime := time.Now()
		startTime := endTime.Add(-time.Duration(hours) * time.Hour)

		historyData, err := h.services.UnifiedSensorData.GetDataByTimeRange(ctx, deviceID, startTime.Unix(), endTime.Unix())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		chartData := h.convertToChartDataFromUnified(historyData, metric, sensorID)
		c.JSON(http.StatusOK, chartData)

	case "sensors":
		deviceID := c.Query("device_id")
		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
			return
		}

		sensorIDs, err := h.services.UnifiedSensorData.GetSensorIDsByDeviceID(ctx, deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sensors": sensorIDs})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
	}
}

// SensorsAPI 传感器列表API
func (h *WebHandlers) SensorsAPI(c *gin.Context) {
	ctx := context.Background()

	deviceID := c.Query("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	sensorIDs, err := h.services.UnifiedSensorData.GetSensorIDsByDeviceID(ctx, deviceID)
	if err != nil {
		h.logger.Error("获取传感器列表失败", utils.ErrorField(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sensors": sensorIDs})
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
