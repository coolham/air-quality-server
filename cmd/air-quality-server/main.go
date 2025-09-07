package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"air-quality-server/internal/config"
	"air-quality-server/internal/handlers"
	"air-quality-server/internal/mqtt"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/router"
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

// getProjectRoot 智能获取项目根目录
func getProjectRoot() string {
	// 方法1: 尝试从可执行文件路径推断
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)

		// 检查是否在 bin 目录下（生产环境）
		if filepath.Base(execDir) == "bin" {
			projectRoot := filepath.Join(execDir, "..")
			if isProjectRoot(projectRoot) {
				return projectRoot
			}
		}

		// 检查是否在 cmd/air-quality-server 目录下（开发环境）
		if filepath.Base(execDir) == "air-quality-server" {
			parentDir := filepath.Join(execDir, "..")
			if filepath.Base(parentDir) == "cmd" {
				projectRoot := filepath.Join(parentDir, "..")
				if isProjectRoot(projectRoot) {
					return projectRoot
				}
			}
		}
	}

	// 方法2: 从当前工作目录开始向上查找
	wd, err := os.Getwd()
	if err == nil {
		currentDir := wd
		for {
			if isProjectRoot(currentDir) {
				return currentDir
			}
			parent := filepath.Join(currentDir, "..")
			if parent == currentDir {
				break // 到达根目录
			}
			currentDir = parent
		}
	}

	// 方法3: 使用相对路径作为后备方案
	return "."
}

// isProjectRoot 检查指定目录是否为项目根目录
func isProjectRoot(dir string) bool {
	// 检查是否存在项目特征文件/目录
	requiredPaths := []string{
		"go.mod",
		"web/templates",
		"web/static",
		"internal",
		"cmd",
	}

	for _, path := range requiredPaths {
		fullPath := filepath.Join(dir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := utils.InitGlobalLogger(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
		cfg.Log.Compress,
	); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	logger := utils.GetLogger()
	logger.Info("启动空气质量监测服务",
		utils.String("version", cfg.Service.Version),
		utils.String("environment", cfg.Service.Environment),
	)

	// 初始化数据库
	db, err := utils.NewDatabase(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("初始化数据库失败", utils.ErrorField(err))
	}
	defer db.Close()

	// 初始化Redis
	redis, err := utils.NewRedis(&cfg.Redis, logger)
	if err != nil {
		logger.Warn("初始化Redis失败，程序将在无Redis模式下运行", utils.ErrorField(err))
		redis = nil // 设置为nil，表示Redis不可用
	} else {
		defer redis.Close()
	}

	// 初始化仓储层
	repos := initRepositories(db.DB, logger)

	// 初始化服务层
	svcs := initServices(repos, redis, logger)

	// 初始化MQTT服务器
	mqttServer := initMQTTServer(cfg, logger, repos, svcs)
	if mqttServer != nil {
		defer mqttServer.Stop()
	}

	// 初始化处理器
	handlers := initHandlers(svcs, logger)

	// 初始化路由
	router := router.InitRouter(handlers, svcs, cfg, logger)

	// 启动服务器
	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("服务器启动", utils.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", utils.ErrorField(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", utils.ErrorField(err))
	}

	logger.Info("服务器已关闭")
}

// loadConfig 加载配置
func loadConfig() (*config.Config, error) {
	// 优先级1: 环境变量 AIR_QUALITY_CONFIG
	if configPath := os.Getenv("AIR_QUALITY_CONFIG"); configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			return config.Load(configPath)
		}
		// 如果环境变量指定的文件不存在，记录警告但继续尝试其他方式
		fmt.Printf("警告: 环境变量 AIR_QUALITY_CONFIG 指定的配置文件不存在: %s\n", configPath)
	}

	// 尝试自动查找配置文件
	projectRoot := getProjectRoot()
	possiblePaths := []string{
		filepath.Join(projectRoot, "config", "config.yaml"),
		filepath.Join(projectRoot, "config.yaml"),
		"config/config.yaml",
		"config.yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return config.Load(path)
		}
	}

	// 如果找不到配置文件，使用环境变量配置
	return config.LoadFromEnv()
}

// initRepositories 初始化仓储层
func initRepositories(db *gorm.DB, logger utils.Logger) *repositories.Repositories {
	return &repositories.Repositories{
		Device:            repositories.NewDeviceRepository(db, logger),
		AirQuality:        repositories.NewAirQualityRepository(db, logger),
		UnifiedSensorData: repositories.NewUnifiedSensorDataRepository(db, logger),
		User:              repositories.NewUserRepository(db, logger),
		Alert:             repositories.NewAlertRepository(db, logger),
		Config:            repositories.NewConfigRepository(db, logger),
	}
}

// initServices 初始化服务层
func initServices(repos *repositories.Repositories, redis *utils.Redis, logger utils.Logger) *services.Services {
	return &services.Services{
		Device:            services.NewDeviceService(repos.Device, logger),
		AirQuality:        services.NewAirQualityService(repos.AirQuality, repos.Device, logger),
		UnifiedSensorData: services.NewUnifiedSensorDataService(repos.UnifiedSensorData, repos.Device, services.NewAlertService(repos.Alert, logger), logger),
		User:              services.NewUserService(repos.User, logger),
		Alert:             services.NewAlertService(repos.Alert, logger),
		Config:            services.NewConfigService(repos.Config, logger),
	}
}

// initMQTTServer 初始化MQTT服务器
func initMQTTServer(cfg *config.Config, logger utils.Logger, repos *repositories.Repositories, svcs *services.Services) *mqtt.Server {
	// 检查MQTT配置
	if cfg.MQTT.Broker == "" {
		logger.Warn("MQTT配置为空，跳过MQTT服务器启动")
		return nil
	}

	// 创建传感器数据处理器
	sensorDataHandler := mqtt.NewSensorDataHandler(
		repos.UnifiedSensorData,
		repos.Device,
		svcs.Alert,
		logger,
	)

	// 创建MQTT服务器
	mqttServer := mqtt.NewServer(&cfg.MQTT, logger, sensorDataHandler)

	// 启动MQTT服务器
	if err := mqttServer.Start(); err != nil {
		logger.Error("启动MQTT服务器失败", utils.ErrorField(err))
		return nil
	}

	logger.Info("MQTT服务器启动成功",
		utils.String("broker", cfg.MQTT.Broker),
		utils.String("client_id", cfg.MQTT.ClientID))

	return mqttServer
}

// initHandlers 初始化处理器
func initHandlers(svcs *services.Services, logger utils.Logger) *handlers.Handlers {
	return &handlers.Handlers{
		Device:     handlers.NewDeviceHandler(svcs.Device, logger),
		AirQuality: handlers.NewAirQualityHandler(svcs.AirQuality, logger),
		User:       handlers.NewUserHandler(svcs.User, logger),
		Alert:      handlers.NewAlertHandler(svcs.Alert, logger),
		Config:     handlers.NewConfigHandler(svcs.Config, logger),
	}
}
