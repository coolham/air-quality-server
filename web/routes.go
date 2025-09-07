package web

import (
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"air-quality-server/web/handlers"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InitWeb 初始化Web模块
func InitWeb(router *gin.Engine, services *services.Services, logger utils.Logger) {
	// 获取Web资源路径
	webTemplatesPath, webStaticPath, webAssetsPath := GetWebPaths()

	logger.Info("Web路径配置",
		zap.String("templates_path", webTemplatesPath),
		zap.String("static_path", webStaticPath),
		zap.String("assets_path", webAssetsPath))

	// 设置模板函数
	router.SetFuncMap(TemplateFuncs)

	// 加载HTML模板
	router.LoadHTMLGlob(filepath.Join(webTemplatesPath, "*"))

	// 设置静态文件服务
	router.Static("/static", webStaticPath)
	router.Static("/assets", webAssetsPath)

	// 设置Web路由
	SetupRoutes(router, services, logger)
}

// SetupRoutes 设置Web路由
func SetupRoutes(router *gin.Engine, services *services.Services, logger utils.Logger) {
	// 创建Web处理器
	webHandlers := handlers.NewWebHandlers(services, logger)

	// 创建错误处理器
	errorHandler := handlers.NewErrorHandler(logger)

	// 添加错误处理中间件
	router.Use(errorHandler.RequestIDMiddleware())
	router.Use(errorHandler.ErrorLoggerMiddleware())
	router.Use(errorHandler.RecoveryMiddleware())

	// 创建Web路由组
	webGroup := router.Group("/")
	{
		// 首页重定向到仪表板
		webGroup.GET("/", func(c *gin.Context) {
			c.Redirect(302, "/dashboard")
		})

		// 仪表板
		webGroup.GET("/dashboard", webHandlers.Dashboard)

		// 设备管理
		webGroup.GET("/devices", webHandlers.DeviceList)
		webGroup.GET("/devices/:id", webHandlers.DeviceDetail)

		// 数据查看
		webGroup.GET("/sensor-data", webHandlers.DataView)

		// 图表分析
		webGroup.GET("/charts", webHandlers.Charts)

		// 告警管理
		webGroup.GET("/alerts", webHandlers.Alerts)

		// 数据导出
		webGroup.GET("/export", webHandlers.DataExportAPI)
	}

	// Web API路由组
	webAPI := router.Group("/web/api")
	{
		// 设备统计API
		webAPI.GET("/device-stats", webHandlers.API)

		// 最新数据API
		webAPI.GET("/latest-data", webHandlers.API)

		// 图表数据API
		webAPI.GET("/chart-data", webHandlers.API)

		// 传感器列表API
		webAPI.GET("/sensors", webHandlers.SensorsAPI)

		// 数据查询API（放在最后，避免路由冲突）
		webAPI.GET("/data", webHandlers.DataAPI)
		webAPI.GET("/data/export", webHandlers.DataExportAPI)
	}

	// 设置404和405处理器
	router.NoRoute(errorHandler.NotFoundHandler)
	router.NoMethod(errorHandler.MethodNotAllowedHandler)
}
