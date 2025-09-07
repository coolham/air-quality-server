package handlers

import (
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AirQualityHandler 空气质量处理器
type AirQualityHandler struct {
	airQualityService services.AirQualityService
	logger            utils.Logger
}

// NewAirQualityHandler 创建空气质量处理器
func NewAirQualityHandler(airQualityService services.AirQualityService, logger utils.Logger) *AirQualityHandler {
	return &AirQualityHandler{
		airQualityService: airQualityService,
		logger:            logger,
	}
}

// UploadData 上传空气质量数据
func (h *AirQualityHandler) UploadData(c *gin.Context) {
	var req struct {
		DeviceID    string  `json:"device_id" binding:"required"`
		PM2_5       float64 `json:"pm2_5" binding:"required"`
		PM10        float64 `json:"pm10" binding:"required"`
		CO2         float64 `json:"co2" binding:"required"`
		Temperature float64 `json:"temperature" binding:"required"`
		Humidity    float64 `json:"humidity" binding:"required"`
		Timestamp   int64   `json:"timestamp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("上传空气质量数据请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("上传空气质量数据请求", utils.String("device_id", req.DeviceID))
	c.JSON(http.StatusCreated, gin.H{
		"message": "数据上传成功",
		"data":    req,
	})
}

// GetRealtimeData 获取实时数据
func (h *AirQualityHandler) GetRealtimeData(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	h.logger.Info("获取实时数据请求", utils.String("device_id", deviceID))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取实时数据成功",
		"data": gin.H{
			"device_id":   deviceID,
			"pm2_5":       25.5,
			"pm10":        45.2,
			"co2":         420.0,
			"temperature": 22.5,
			"humidity":    65.0,
			"timestamp":   1640995200,
			"quality":     "良",
		},
	})
}

// GetHistoryData 获取历史数据
func (h *AirQualityHandler) GetHistoryData(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	h.logger.Info("获取历史数据请求",
		utils.String("device_id", deviceID),
		utils.Int("limit", limit),
		utils.String("start_time", startTimeStr),
		utils.String("end_time", endTimeStr))

	c.JSON(http.StatusOK, gin.H{
		"message": "获取历史数据成功",
		"data": gin.H{
			"device_id": deviceID,
			"data":      []gin.H{},
			"total":     0,
			"limit":     limit,
		},
	})
}

// GetStatistics 获取统计数据
func (h *AirQualityHandler) GetStatistics(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	h.logger.Info("获取统计数据请求",
		utils.String("device_id", deviceID),
		utils.String("start_time", startTimeStr),
		utils.String("end_time", endTimeStr))

	c.JSON(http.StatusOK, gin.H{
		"message": "获取统计数据成功",
		"data": gin.H{
			"device_id": deviceID,
			"statistics": gin.H{
				"pm2_5_avg":    25.5,
				"pm10_avg":     45.2,
				"co2_avg":      420.0,
				"temp_avg":     22.5,
				"humidity_avg": 65.0,
				"data_points":  100,
			},
		},
	})
}

// ExportData 导出数据
func (h *AirQualityHandler) ExportData(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	format := c.DefaultQuery("format", "csv")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	h.logger.Info("导出数据请求",
		utils.String("device_id", deviceID),
		utils.String("format", format),
		utils.String("start_time", startTimeStr),
		utils.String("end_time", endTimeStr))

	// 设置响应头
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=air_quality_data.csv")

	// 返回CSV数据
	csvData := "timestamp,pm2_5,pm10,co2,temperature,humidity\n"
	c.String(http.StatusOK, csvData)
}
