package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"air-quality-server/internal/config"
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

func main() {
	var (
		configPath = flag.String("config", "", "配置文件路径")
		action     = flag.String("action", "init", "操作类型: init, status")
		help       = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 加载配置
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	if err != nil {
		fmt.Printf("❌ 初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	logger.Info("数据库初始化工具启动")

	// 连接数据库
	db, err := utils.NewDatabase(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("连接数据库失败", utils.ErrorField(err))
	}
	defer db.Close()

	// 执行操作
	switch *action {
	case "init":
		if err := runInit(db.DB, logger); err != nil {
			logger.Fatal("数据库初始化失败", utils.ErrorField(err))
		}
	case "status":
		if err := showStatus(db.DB, logger); err != nil {
			logger.Fatal("获取状态失败", utils.ErrorField(err))
		}
	default:
		fmt.Printf("❌ 未知操作: %s\n", *action)
		showHelp()
		os.Exit(1)
	}

	logger.Info("操作完成")
}

// loadConfig 加载配置
func loadConfig(configPath string) (*config.Config, error) {
	if configPath != "" {
		return config.Load(configPath)
	}

	// 自动查找配置文件
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

	// 使用环境变量配置
	return config.LoadFromEnv()
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		if filepath.Base(execDir) == "migrate" {
			parentDir := filepath.Join(execDir, "..")
			if filepath.Base(parentDir) == "cmd" {
				return filepath.Join(parentDir, "..")
			}
		}
	}

	wd, err := os.Getwd()
	if err == nil {
		currentDir := wd
		for {
			if isProjectRoot(currentDir) {
				return currentDir
			}
			parent := filepath.Join(currentDir, "..")
			if parent == currentDir {
				break
			}
			currentDir = parent
		}
	}

	return "."
}

// isProjectRoot 检查是否为项目根目录
func isProjectRoot(dir string) bool {
	requiredPaths := []string{
		"go.mod",
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

// runInit 执行数据库初始化
func runInit(db *gorm.DB, logger utils.Logger) error {
	logger.Info("开始执行数据库初始化...")

	// 获取所有数据模型
	models := getAllModels()

	// 执行自动迁移
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	logger.Info("数据库表结构创建完成", utils.Int("models_count", len(models)))

	// 插入初始数据
	if err := insertInitialData(db, logger); err != nil {
		return fmt.Errorf("插入初始数据失败: %w", err)
	}

	logger.Info("数据库初始化完成")
	return nil
}

// showStatus 显示数据库状态
func showStatus(db *gorm.DB, logger utils.Logger) error {
	logger.Info("数据库状态信息:")

	// 检查表是否存在
	tables := []string{
		"users", "roles", "user_roles",
		"devices", "unified_sensor_data", "device_runtime_status",
		"alerts", "alert_rules", "system_configs",
	}

	for _, table := range tables {
		var exists bool
		if err := db.Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table).Scan(&exists).Error; err != nil {
			logger.Error("检查表失败", utils.String("table", table), utils.ErrorField(err))
			continue
		}

		status := "❌ 不存在"
		if exists {
			status = "✅ 存在"
		}
		logger.Info("表状态", utils.String("table", table), utils.String("status", status))
	}

	// 检查数据量
	models := getAllModels()
	for _, model := range models {
		var count int64
		if err := db.Model(model).Count(&count).Error; err == nil {
			tableName := getTableName(model)
			logger.Info("数据统计", utils.String("table", tableName), utils.Int64("count", count))
		}
	}

	return nil
}

// getAllModels 获取所有数据模型
func getAllModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Device{},
		&models.UnifiedSensorData{},
		&models.DeviceRuntimeStatus{},
		&models.Alert{},
		&models.AlertRule{},
		&models.SystemConfig{},
	}
}

// getTableName 获取表名
func getTableName(model interface{}) string {
	if tableName, ok := model.(interface{ TableName() string }); ok {
		return tableName.TableName()
	}
	return "unknown"
}

// insertInitialData 插入初始数据
func insertInitialData(db *gorm.DB, logger utils.Logger) error {
	logger.Info("插入初始数据...")

	// 插入默认角色
	roles := []models.Role{
		{
			Name:        "admin",
			Description: stringPtr("系统管理员"),
			Permissions: stringPtr(`["*"]`),
		},
		{
			Name:        "operator",
			Description: stringPtr("操作员"),
			Permissions: stringPtr(`["device:read", "device:write", "data:read", "alert:read", "alert:write"]`),
		},
		{
			Name:        "viewer",
			Description: stringPtr("查看者"),
			Permissions: stringPtr(`["device:read", "data:read", "alert:read"]`),
		},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			logger.Error("创建角色失败", utils.String("role", role.Name), utils.ErrorField(err))
		} else {
			logger.Info("角色已创建", utils.String("role", role.Name))
		}
	}

	// 插入默认用户
	adminUser := models.User{
		Username:     "admin",
		Email:        "admin@air-quality.com",
		PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // admin123
		Status:       "active",
	}

	if err := db.FirstOrCreate(&adminUser, models.User{Username: adminUser.Username}).Error; err != nil {
		logger.Error("创建管理员用户失败", utils.ErrorField(err))
	} else {
		logger.Info("管理员用户已创建", utils.String("username", adminUser.Username))
	}

	// 为用户分配管理员角色
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == nil {
		userRole := models.UserRole{
			UserID: adminUser.ID,
			RoleID: adminRole.ID,
		}
		if err := db.FirstOrCreate(&userRole, models.UserRole{UserID: userRole.UserID, RoleID: userRole.RoleID}).Error; err != nil {
			logger.Error("分配角色失败", utils.ErrorField(err))
		} else {
			logger.Info("管理员角色已分配")
		}
	}

	// 插入默认系统配置
	configs := []models.SystemConfig{
		{
			KeyName:     "data_retention_days",
			Value:       stringPtr("365"),
			Description: stringPtr("数据保留天数"),
			Category:    stringPtr("data"),
		},
		{
			KeyName:     "alert_check_interval",
			Value:       stringPtr("60"),
			Description: stringPtr("告警检查间隔(秒)"),
			Category:    stringPtr("alert"),
		},
		{
			KeyName:     "max_devices_per_user",
			Value:       stringPtr("100"),
			Description: stringPtr("每用户最大设备数"),
			Category:    stringPtr("device"),
		},
		{
			KeyName:     "api_rate_limit",
			Value:       stringPtr("1000"),
			Description: stringPtr("API请求限制(每小时)"),
			Category:    stringPtr("api"),
		},
		{
			KeyName:     "data_quality_threshold",
			Value:       stringPtr("0.8"),
			Description: stringPtr("数据质量阈值"),
			Category:    stringPtr("data"),
		},
	}

	for _, cfg := range configs {
		if err := db.FirstOrCreate(&cfg, models.SystemConfig{KeyName: cfg.KeyName}).Error; err != nil {
			logger.Error("创建系统配置失败", utils.String("key", cfg.KeyName), utils.ErrorField(err))
		} else {
			logger.Info("系统配置已创建", utils.String("key", cfg.KeyName))
		}
	}

	logger.Info("初始数据插入完成")
	return nil
}

// stringPtr 创建字符串指针
func stringPtr(s string) *string {
	return &s
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("数据库初始化工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  go run cmd/migrate/main.go [选项]")
	fmt.Println("")
	fmt.Println("选项:")
	fmt.Println("  -config string")
	fmt.Println("        配置文件路径 (默认自动查找)")
	fmt.Println("  -action string")
	fmt.Println("        操作类型: init, status (默认: init)")
	fmt.Println("  -help")
	fmt.Println("        显示帮助信息")
	fmt.Println("")
	fmt.Println("操作说明:")
	fmt.Println("  init     - 初始化数据库（创建表结构 + 插入初始数据）")
	fmt.Println("  status   - 显示数据库状态信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  # 初始化数据库（推荐首次部署时使用）")
	fmt.Println("  go run cmd/migrate/main.go -action init")
	fmt.Println("")
	fmt.Println("  # 查看数据库状态")
	fmt.Println("  go run cmd/migrate/main.go -action status")
	fmt.Println("")
	fmt.Println("  # 使用指定配置文件")
	fmt.Println("  go run cmd/migrate/main.go -config config/config.yaml -action init")
	fmt.Println("")
	fmt.Println("注意:")
	fmt.Println("  - 此工具用于全新系统初始化，会创建所有必要的表和数据")
	fmt.Println("  - 如果数据库已存在数据，请先备份或清空数据库")
	fmt.Println("  - 初始化完成后即可启动主程序")
}
