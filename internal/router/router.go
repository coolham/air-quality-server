package router

import (
	"air-quality-server/api"
	"air-quality-server/internal/config"
	"air-quality-server/internal/handlers"
	"air-quality-server/internal/middleware"
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"air-quality-server/web"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器
func InitRouter(handlers *handlers.Handlers, services *services.Services, cfg *config.Config, logger utils.Logger) *gin.Engine {
	// 设置Gin模式
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 添加中间件
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// 设置API路由
	api.SetupAPIRoutes(router, handlers, services, cfg, logger)

	// 初始化Web模块
	web.InitWeb(router, services, logger)

	return router
}
