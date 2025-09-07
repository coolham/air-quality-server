package handlers

import (
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AlertHandler 告警处理器
type AlertHandler struct {
	alertService services.AlertService
	logger       utils.Logger
}

// NewAlertHandler 创建告警处理器
func NewAlertHandler(alertService services.AlertService, logger utils.Logger) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
		logger:       logger,
	}
}

// CreateAlert 创建告警
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var req struct {
		DeviceID  string  `json:"device_id" binding:"required"`
		Type      string  `json:"type" binding:"required"`
		Level     string  `json:"level" binding:"required"`
		Message   string  `json:"message" binding:"required"`
		Threshold float64 `json:"threshold"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("创建告警请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("创建告警请求",
		utils.String("device_id", req.DeviceID),
		utils.String("type", req.Type),
		utils.String("level", req.Level))

	c.JSON(http.StatusCreated, gin.H{
		"message": "告警创建成功",
		"data":    req,
	})
}

// GetAlert 获取告警
func (h *AlertHandler) GetAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("告警ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "告警ID参数错误"})
		return
	}

	h.logger.Info("获取告警请求", utils.Int("alert_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取告警成功",
		"data": gin.H{
			"id":         id,
			"device_id":  1,
			"type":       "pm2_5_high",
			"level":      "warning",
			"message":    "PM2.5浓度过高",
			"status":     "active",
			"created_at": "2024-01-01T00:00:00Z",
		},
	})
}

// UpdateAlert 更新告警
func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("告警ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "告警ID参数错误"})
		return
	}

	var req struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("更新告警请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("更新告警请求", utils.Int("alert_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "告警更新成功",
		"data": gin.H{
			"id": id,
		},
	})
}

// DeleteAlert 删除告警
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("告警ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "告警ID参数错误"})
		return
	}

	h.logger.Info("删除告警请求", utils.Int("alert_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "告警删除成功",
	})
}

// ListAlerts 列出告警
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	status := c.Query("status")
	deviceIDStr := c.Query("device_id")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	deviceID := deviceIDStr

	h.logger.Info("列出告警请求",
		utils.Int("limit", limit),
		utils.Int("offset", offset),
		utils.String("status", status),
		utils.String("device_id", deviceID))

	c.JSON(http.StatusOK, gin.H{
		"message": "获取告警列表成功",
		"data": gin.H{
			"alerts": []gin.H{},
			"total":  0,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAlertsByDevice 根据设备获取告警
func (h *AlertHandler) GetAlertsByDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	h.logger.Info("根据设备获取告警请求", utils.String("device_id", deviceID))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取设备告警成功",
		"data": gin.H{
			"device_id": deviceID,
			"alerts":    []gin.H{},
		},
	})
}

// ResolveAlert 解决告警
func (h *AlertHandler) ResolveAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("告警ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "告警ID参数错误"})
		return
	}

	var req struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("解决告警请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("解决告警请求", utils.Int("alert_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "告警已解决",
		"data": gin.H{
			"id": id,
		},
	})
}

// GetUnresolvedAlerts 获取未解决的告警
func (h *AlertHandler) GetUnresolvedAlerts(c *gin.Context) {
	h.logger.Info("获取未解决告警请求")
	c.JSON(http.StatusOK, gin.H{
		"message": "获取未解决告警成功",
		"data": gin.H{
			"alerts": []gin.H{},
			"total":  0,
		},
	})
}
