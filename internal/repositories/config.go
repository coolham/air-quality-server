package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

// ConfigRepository 配置仓储接口
type ConfigRepository interface {
	BaseRepository[models.SystemConfig]
	GetByKey(key string) (*models.SystemConfig, error)
	GetByCategory(category string) ([]models.SystemConfig, error)
	UpdateByKey(key string, value string) error
	GetAllAsMap() (map[string]string, error)
}

// configRepository 配置仓储实现
type configRepository struct {
	*baseRepository[models.SystemConfig]
	db     *gorm.DB
	logger utils.Logger
}

// NewConfigRepository 创建配置仓储
func NewConfigRepository(db *gorm.DB, logger utils.Logger) ConfigRepository {
	return &configRepository{
		baseRepository: NewBaseRepository[models.SystemConfig](db, logger).(*baseRepository[models.SystemConfig]),
		db:             db,
		logger:         logger,
	}
}

// GetByKey 根据键获取配置
func (r *configRepository) GetByKey(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.Where(&models.SystemConfig{KeyName: key}).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("根据键获取配置失败", utils.ErrorField(err), utils.String("key", key))
		return nil, err
	}
	return &config, nil
}

// GetByCategory 根据分类获取配置
func (r *configRepository) GetByCategory(category string) ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	err := r.db.Where(&models.SystemConfig{Category: &category}).Find(&configs).Error
	if err != nil {
		r.logger.Error("根据分类获取配置失败", utils.ErrorField(err), utils.String("category", category))
		return nil, err
	}
	return configs, nil
}

// UpdateByKey 根据键更新配置值
func (r *configRepository) UpdateByKey(key string, value string) error {
	updateData := &models.SystemConfig{Value: &value}
	err := r.db.Model(&models.SystemConfig{}).Where("key_name = ?", key).Updates(updateData).Error
	if err != nil {
		r.logger.Error("根据键更新配置失败", utils.ErrorField(err), utils.String("key", key))
		return err
	}
	return nil
}

// GetAllAsMap 获取所有配置为map格式
func (r *configRepository) GetAllAsMap() (map[string]string, error) {
	var configs []models.SystemConfig
	err := r.db.Find(&configs).Error
	if err != nil {
		r.logger.Error("获取所有配置失败", utils.ErrorField(err))
		return nil, err
	}

	result := make(map[string]string)
	for _, config := range configs {
		if config.Value != nil {
			result[config.KeyName] = *config.Value
		}
	}

	return result, nil
}
