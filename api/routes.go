package api

import (
	"air-quality-server/internal/config"
	"air-quality-server/internal/handlers"
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes 设置API路由
func SetupAPIRoutes(router *gin.Engine, handlers *handlers.Handlers, services *services.Services, cfg *config.Config, logger utils.Logger) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   cfg.Service.Name,
			"version":   cfg.Service.Version,
		})
	})

	// API路由组
	api := router.Group("/api/v1")
	{
		// 设备管理
		devices := api.Group("/devices")
		{
			devices.GET("", handlers.Device.ListDevices)
			devices.POST("", handlers.Device.CreateDevice)
			devices.GET("/:id", handlers.Device.GetDevice)
			devices.PUT("/:id", handlers.Device.UpdateDevice)
			devices.DELETE("/:id", handlers.Device.DeleteDevice)
			devices.GET("/:id/status", handlers.Device.GetDeviceStatus)
			// devices.PUT("/:id/status", handlers.Device.UpdateDeviceStatus) // 方法未实现
			// devices.GET("/:id/statistics", handlers.Device.GetDeviceStatistics) // 方法未实现
		}

		// 数据管理
		data := api.Group("/data")
		{
			data.POST("/upload", handlers.AirQuality.UploadData)
			data.GET("/realtime/:device_id", handlers.AirQuality.GetRealtimeData)
			data.GET("/history/:device_id", handlers.AirQuality.GetHistoryData)
			data.GET("/statistics/:device_id", handlers.AirQuality.GetStatistics)
			data.GET("/export/:device_id", handlers.AirQuality.ExportData)
		}

		// 用户管理
		users := api.Group("/users")
		{
			users.GET("", handlers.User.ListUsers)
			users.POST("", handlers.User.CreateUser)
			users.GET("/:id", handlers.User.GetUser)
			users.PUT("/:id", handlers.User.UpdateUser)
			users.DELETE("/:id", handlers.User.DeleteUser)
		}

		// 认证
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.User.Login)
			auth.POST("/logout", handlers.User.Logout)
			auth.POST("/change-password", handlers.User.ChangePassword)
		}

		// 告警管理
		alerts := api.Group("/alerts")
		{
			alerts.GET("", handlers.Alert.ListAlerts)
			alerts.POST("", handlers.Alert.CreateAlert)
			alerts.GET("/:id", handlers.Alert.GetAlert)
			alerts.PUT("/:id", handlers.Alert.UpdateAlert)
			alerts.DELETE("/:id", handlers.Alert.DeleteAlert)
			alerts.GET("/device/:device_id", handlers.Alert.GetAlertsByDevice)
			alerts.POST("/:id/resolve", handlers.Alert.ResolveAlert)
			alerts.GET("/unresolved", handlers.Alert.GetUnresolvedAlerts)
		}

		// 配置管理
		configs := api.Group("/configs")
		{
			configs.GET("", handlers.Config.GetAllConfigs)
			configs.GET("/:key", handlers.Config.GetConfig)
			configs.POST("/:key", handlers.Config.SetConfig)
			configs.PUT("/:key", handlers.Config.UpdateConfig)
			configs.DELETE("/:key", handlers.Config.DeleteConfig)
			configs.GET("/category/:category", handlers.Config.GetConfigsByCategory)
			configs.GET("/system/settings", handlers.Config.GetSystemSettings)
			configs.PUT("/system/settings", handlers.Config.UpdateSystemSettings)
		}
	}

	// WebSocket支持
	// router.GET("/ws/data", handlers.AirQuality.WebSocket) // 暂时注释掉，因为WebSocket方法还未实现
}
