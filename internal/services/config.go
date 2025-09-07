package services

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/utils"
	"context"
	"time"
)

// ConfigService 配置服务接口
type ConfigService interface {
	GetConfig(ctx context.Context, key string) (*models.SystemConfig, error)
	SetConfig(ctx context.Context, key, value, category, description string) error
	GetConfigsByCategory(ctx context.Context, category string) ([]models.SystemConfig, error)
	GetAllConfigs(ctx context.Context) (map[string]string, error)
	UpdateConfig(ctx context.Context, config *models.SystemConfig) error
	DeleteConfig(ctx context.Context, key string) error
	GetSystemSettings(ctx context.Context) (map[string]interface{}, error)
	UpdateSystemSettings(ctx context.Context, settings map[string]interface{}) error
}

// configService 配置服务实现
type configService struct {
	configRepo repositories.ConfigRepository
	logger     utils.Logger
}

// NewConfigService 创建配置服务
func NewConfigService(configRepo repositories.ConfigRepository, logger utils.Logger) ConfigService {
	return &configService{
		configRepo: configRepo,
		logger:     logger,
	}
}

// GetConfig 获取配置
func (s *configService) GetConfig(ctx context.Context, key string) (*models.SystemConfig, error) {
	config, err := s.configRepo.GetByKey(key)
	if err != nil {
		s.logger.Error("获取配置失败", utils.ErrorField(err), utils.String("key", key))
		return nil, err
	}
	return config, nil
}

// SetConfig 设置配置
func (s *configService) SetConfig(ctx context.Context, key, value, category, description string) error {
	config := &models.SystemConfig{
		KeyName:     key,
		Value:       &value,
		Category:    &category,
		Description: &description,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	// 检查配置是否已存在
	existingConfig, err := s.configRepo.GetByKey(key)
	if err != nil {
		s.logger.Error("检查配置是否存在失败", utils.ErrorField(err), utils.String("key", key))
		return err
	}

	if existingConfig != nil {
		// 更新现有配置
		existingConfig.Value = &value
		existingConfig.Category = &category
		existingConfig.Description = &description
		existingConfig.UpdatedAt = now

		updateData := &models.SystemConfig{
			Value:       &value,
			Category:    &category,
			Description: &description,
			UpdatedAt:   now,
		}
		if err := s.configRepo.Update(ctx, existingConfig.ID, updateData); err != nil {
			s.logger.Error("更新配置失败", utils.ErrorField(err), utils.String("key", key))
			return err
		}
	} else {
		// 创建新配置
		if err := s.configRepo.Create(ctx, config); err != nil {
			s.logger.Error("创建配置失败", utils.ErrorField(err), utils.String("key", key))
			return err
		}
	}

	s.logger.Info("配置设置成功", utils.String("key", key))
	return nil
}

// GetConfigsByCategory 根据分类获取配置
func (s *configService) GetConfigsByCategory(ctx context.Context, category string) ([]models.SystemConfig, error) {
	configs, err := s.configRepo.GetByCategory(category)
	if err != nil {
		s.logger.Error("根据分类获取配置失败", utils.ErrorField(err), utils.String("category", category))
		return nil, err
	}
	return configs, nil
}

// GetAllConfigs 获取所有配置
func (s *configService) GetAllConfigs(ctx context.Context) (map[string]string, error) {
	configs, err := s.configRepo.GetAllAsMap()
	if err != nil {
		s.logger.Error("获取所有配置失败", utils.ErrorField(err))
		return nil, err
	}
	return configs, nil
}

// UpdateConfig 更新配置
func (s *configService) UpdateConfig(ctx context.Context, config *models.SystemConfig) error {
	now := time.Now()
	config.UpdatedAt = now

	updateData := &models.SystemConfig{
		Value:       config.Value,
		Category:    config.Category,
		Description: config.Description,
		UpdatedAt:   config.UpdatedAt,
	}
	if err := s.configRepo.Update(ctx, config.ID, updateData); err != nil {
		s.logger.Error("更新配置失败", utils.ErrorField(err), utils.String("key", config.KeyName))
		return err
	}

	s.logger.Info("配置更新成功", utils.String("key", config.KeyName))
	return nil
}

// DeleteConfig 删除配置
func (s *configService) DeleteConfig(ctx context.Context, key string) error {
	config, err := s.configRepo.GetByKey(key)
	if err != nil {
		s.logger.Error("获取配置失败", utils.ErrorField(err), utils.String("key", key))
		return err
	}

	if config == nil {
		return nil // 配置不存在，无需删除
	}

	if err := s.configRepo.Delete(ctx, config.ID); err != nil {
		s.logger.Error("删除配置失败", utils.ErrorField(err), utils.String("key", key))
		return err
	}

	s.logger.Info("配置删除成功", utils.String("key", key))
	return nil
}

// GetSystemSettings 获取系统设置
func (s *configService) GetSystemSettings(ctx context.Context) (map[string]interface{}, error) {
	configs, err := s.configRepo.GetByCategory("system")
	if err != nil {
		s.logger.Error("获取系统设置失败", utils.ErrorField(err))
		return nil, err
	}

	settings := make(map[string]interface{})

	for _, config := range configs {
		if config.Value == nil {
			continue
		}
		switch config.KeyName {
		case "data_retention_days":
			// 解析数据保留天数
			if days, err := utils.StringToInt(*config.Value); err == nil {
				settings["data_retention_days"] = days
			}
		case "alert_check_interval":
			// 解析告警检查间隔
			if interval, err := utils.StringToInt(*config.Value); err == nil {
				settings["alert_check_interval"] = interval
			}
		case "max_devices":
			// 解析最大设备数
			if max, err := utils.StringToInt(*config.Value); err == nil {
				settings["max_devices"] = max
			}
		case "enable_notifications":
			// 解析是否启用通知
			settings["enable_notifications"] = *config.Value == "true"
		case "notification_email":
			settings["notification_email"] = *config.Value
		}
	}

	return settings, nil
}

// UpdateSystemSettings 更新系统设置
func (s *configService) UpdateSystemSettings(ctx context.Context, settings map[string]interface{}) error {
	// 更新数据保留天数
	if days, ok := settings["data_retention_days"].(int); ok {
		if err := s.SetConfig(ctx, "data_retention_days", utils.IntToString(days), "system", "数据保留天数"); err != nil {
			return err
		}
	}

	// 更新告警检查间隔
	if interval, ok := settings["alert_check_interval"].(int); ok {
		if err := s.SetConfig(ctx, "alert_check_interval", utils.IntToString(interval), "system", "告警检查间隔(秒)"); err != nil {
			return err
		}
	}

	// 更新最大设备数
	if max, ok := settings["max_devices"].(int); ok {
		if err := s.SetConfig(ctx, "max_devices", utils.IntToString(max), "system", "最大设备数"); err != nil {
			return err
		}
	}

	// 更新通知设置
	if enable, ok := settings["enable_notifications"].(bool); ok {
		if err := s.SetConfig(ctx, "enable_notifications", utils.BoolToString(enable), "system", "是否启用通知"); err != nil {
			return err
		}
	}

	// 更新通知邮箱
	if email, ok := settings["notification_email"].(string); ok {
		if err := s.SetConfig(ctx, "notification_email", email, "system", "通知邮箱"); err != nil {
			return err
		}
	}

	s.logger.Info("系统设置更新成功")
	return nil
}
