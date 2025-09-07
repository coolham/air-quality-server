package utils

import (
	"fmt"

	"air-quality-server/internal/models"

	"gorm.io/gorm"
)

// InitDatabase 初始化数据库（迁移 + 初始数据）
func InitDatabase(db *gorm.DB, logger Logger, force bool) error {
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 检查数据库是否已初始化
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err == nil && count > 0 && !force {
		if logger != nil {
			logger.Warn("数据库已包含数据，跳过初始化")
		}
		return nil
	}

	// 执行自动迁移
	if err := db.AutoMigrate(getAllModels()...); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	if logger != nil {
		logger.Info("数据库自动迁移完成")
	}

	// 插入初始数据
	if err := insertInitialData(db, logger); err != nil {
		return fmt.Errorf("插入初始数据失败: %w", err)
	}

	if logger != nil {
		logger.Info("数据库初始化完成")
	}

	return nil
}

// insertInitialData 插入初始数据
func insertInitialData(db *gorm.DB, logger Logger) error {
	if logger != nil {
		logger.Info("插入初始数据...")
	}

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
			if logger != nil {
				logger.Error("创建角色失败", String("role", role.Name), ErrorField(err))
			}
		} else if logger != nil {
			logger.Info("角色已创建", String("role", role.Name))
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
		if logger != nil {
			logger.Error("创建管理员用户失败", ErrorField(err))
		}
	} else if logger != nil {
		logger.Info("管理员用户已创建", String("username", adminUser.Username))
	}

	// 为用户分配管理员角色
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == nil {
		userRole := models.UserRole{
			UserID: adminUser.ID,
			RoleID: adminRole.ID,
		}
		if err := db.FirstOrCreate(&userRole, models.UserRole{UserID: userRole.UserID, RoleID: userRole.RoleID}).Error; err != nil {
			if logger != nil {
				logger.Error("分配角色失败", ErrorField(err))
			}
		} else if logger != nil {
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
		{
			KeyName:     "mqtt_broker_url",
			Value:       stringPtr("tcp://localhost:1883"),
			Description: stringPtr("MQTT Broker地址"),
			Category:    stringPtr("mqtt"),
		},
		{
			KeyName:     "mqtt_username",
			Value:       stringPtr("admin"),
			Description: stringPtr("MQTT用户名"),
			Category:    stringPtr("mqtt"),
		},
		{
			KeyName:     "mqtt_password",
			Value:       stringPtr("password"),
			Description: stringPtr("MQTT密码"),
			Category:    stringPtr("mqtt"),
		},
	}

	for _, cfg := range configs {
		if err := db.FirstOrCreate(&cfg, models.SystemConfig{KeyName: cfg.KeyName}).Error; err != nil {
			if logger != nil {
				logger.Error("创建系统配置失败", String("key", cfg.KeyName), ErrorField(err))
			}
		} else if logger != nil {
			logger.Info("系统配置已创建", String("key", cfg.KeyName))
		}
	}

	// 插入示例设备
	devices := []models.Device{
		{
			ID:                "hcho_001",
			Name:              "甲醛传感器001",
			Type:              models.DeviceTypeFormaldehyde,
			LocationLatitude:  float64Ptr(39.9042),
			LocationLongitude: float64Ptr(116.4074),
			LocationAddress:   stringPtr("北京市朝阳区测试位置"),
			Status:            "online",
			Config:            stringPtr(`{"report_interval": 60, "sensors": {"formaldehyde": true, "temperature": true, "humidity": true}}`),
		},
		{
			ID:                "hcho_002",
			Name:              "甲醛传感器002",
			Type:              models.DeviceTypeFormaldehyde,
			LocationLatitude:  float64Ptr(39.9142),
			LocationLongitude: float64Ptr(116.4174),
			LocationAddress:   stringPtr("北京市海淀区测试位置"),
			Status:            "online",
			Config:            stringPtr(`{"report_interval": 60, "sensors": {"formaldehyde": true, "temperature": true, "humidity": true}}`),
		},
	}

	for _, device := range devices {
		if err := db.FirstOrCreate(&device, models.Device{ID: device.ID}).Error; err != nil {
			if logger != nil {
				logger.Error("创建设备失败", String("device_id", device.ID), ErrorField(err))
			}
		} else if logger != nil {
			logger.Info("设备已创建", String("device_id", device.ID), String("name", device.Name))
		}
	}

	// 插入默认告警规则
	alertRules := []models.AlertRule{
		{
			Name:                 "甲醛浓度警告",
			DeviceID:             stringPtr(""),
			Metric:               "formaldehyde",
			ConditionType:        "gt",
			ThresholdValue:       0.08,
			DurationSeconds:      300,
			Severity:             "warning",
			Enabled:              true,
			NotificationChannels: stringPtr(`["email", "sms"]`),
		},
		{
			Name:                 "甲醛浓度严重",
			DeviceID:             stringPtr(""),
			Metric:               "formaldehyde",
			ConditionType:        "gt",
			ThresholdValue:       0.1,
			DurationSeconds:      180,
			Severity:             "critical",
			Enabled:              true,
			NotificationChannels: stringPtr(`["email", "sms", "webhook"]`),
		},
		{
			Name:                 "设备离线告警",
			DeviceID:             stringPtr(""),
			Metric:               "online_status",
			ConditionType:        "eq",
			ThresholdValue:       0,
			DurationSeconds:      300,
			Severity:             "warning",
			Enabled:              true,
			NotificationChannels: stringPtr(`["email"]`),
		},
	}

	for _, rule := range alertRules {
		if err := db.FirstOrCreate(&rule, models.AlertRule{Name: rule.Name}).Error; err != nil {
			if logger != nil {
				logger.Error("创建告警规则失败", String("rule", rule.Name), ErrorField(err))
			}
		} else if logger != nil {
			logger.Info("告警规则已创建", String("rule", rule.Name))
		}
	}

	if logger != nil {
		logger.Info("初始数据插入完成")
	}

	return nil
}

// stringPtr 创建字符串指针
func stringPtr(s string) *string {
	return &s
}

// float64Ptr 创建float64指针
func float64Ptr(f float64) *float64 {
	return &f
}
