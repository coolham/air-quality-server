package handlers

import (
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeviceHandler 设备处理器
type DeviceHandler struct {
	deviceService services.DeviceService
	logger        utils.Logger
}

// NewDeviceHandler 创建设备处理器
func NewDeviceHandler(deviceService services.DeviceService, logger utils.Logger) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
		logger:        logger,
	}
}

// CreateDevice 创建设备
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req struct {
		SerialNumber string `json:"serial_number" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Location     string `json:"location"`
		Description  string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("创建设备请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 这里需要调用服务层创建设备
	// 由于我们还没有完整的模型定义，先返回成功响应
	h.logger.Info("创建设备请求", utils.String("serial_number", req.SerialNumber))
	c.JSON(http.StatusCreated, gin.H{
		"message": "设备创建成功",
		"data":    req,
	})
}

// GetDevice 获取设备
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	h.logger.Info("获取设备请求", utils.String("device_id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取设备成功",
		"data": gin.H{
			"id": id,
		},
	})
}

// UpdateDevice 更新设备
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("更新设备请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("更新设备请求", utils.String("device_id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "设备更新成功",
		"data": gin.H{
			"id": id,
		},
	})
}

// DeleteDevice 删除设备
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	h.logger.Info("删除设备请求", utils.String("device_id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "设备删除成功",
	})
}

// ListDevices 列出设备
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	h.logger.Info("列出设备请求", utils.Int("limit", limit), utils.Int("offset", offset))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取设备列表成功",
		"data": gin.H{
			"devices": []gin.H{},
			"total":   0,
			"limit":   limit,
			"offset":  offset,
		},
	})
}

// GetDeviceStatus 获取设备状态
func (h *DeviceHandler) GetDeviceStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("设备ID参数错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备ID参数错误"})
		return
	}

	h.logger.Info("获取设备状态请求", utils.String("device_id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取设备状态成功",
		"data": gin.H{
			"device_id": id,
			"status":    "online",
			"last_seen": "2024-01-01T00:00:00Z",
		},
	})
}
