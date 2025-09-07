package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

// FormaldehydeDataRepository 甲醛数据仓库接口
type FormaldehydeDataRepository interface {
	BaseRepository[models.FormaldehydeData]

	// 获取设备最新数据
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.FormaldehydeData, error)

	// 获取设备历史数据
	GetHistoryByDeviceID(ctx context.Context, deviceID string, limit, offset int) ([]models.FormaldehydeData, error)

	// 获取设备指定时间范围的数据
	GetByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.FormaldehydeData, error)

	// 获取设备统计数据
	GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.FormaldehydeStatistics, error)

	// 批量插入数据
	BatchInsert(ctx context.Context, data []models.FormaldehydeData) error
}

// DeviceStatusRepository 设备状态仓库接口
type DeviceStatusRepository interface {
	BaseRepository[models.FormaldehydeDeviceStatus]

	// 根据设备ID获取状态
	GetByDeviceID(ctx context.Context, deviceID string) (*models.FormaldehydeDeviceStatus, error)

	// 获取在线设备列表
	GetOnlineDevices(ctx context.Context) ([]models.FormaldehydeDeviceStatus, error)

	// 获取离线设备列表
	GetOfflineDevices(ctx context.Context) ([]models.FormaldehydeDeviceStatus, error)

	// 更新设备在线状态
	UpdateOnlineStatus(ctx context.Context, deviceID string, online bool) error

	// 批量更新设备状态
	BatchUpdateStatus(ctx context.Context, updates map[string]bool) error
}

// DeviceConfigRepository 设备配置仓库接口
type DeviceConfigRepository interface {
	BaseRepository[models.FormaldehydeDeviceConfig]

	// 根据设备ID获取配置
	GetByDeviceID(ctx context.Context, deviceID string) (*models.FormaldehydeDeviceConfig, error)

	// 更新设备配置
	UpdateConfig(ctx context.Context, deviceID string, config *models.FormaldehydeDeviceConfig) error
}

// formaldehydeDataRepository 甲醛数据仓库实现
type formaldehydeDataRepository struct {
	*baseRepository[models.FormaldehydeData]
}

// NewFormaldehydeDataRepository 创建甲醛数据仓库
func NewFormaldehydeDataRepository(db *gorm.DB, logger utils.Logger) FormaldehydeDataRepository {
	return &formaldehydeDataRepository{
		baseRepository: NewBaseRepository[models.FormaldehydeData](db, logger).(*baseRepository[models.FormaldehydeData]),
	}
}

// GetLatestByDeviceID 获取设备最新数据
func (r *formaldehydeDataRepository) GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.FormaldehydeData, error) {
	var data models.FormaldehydeData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetHistoryByDeviceID 获取设备历史数据
func (r *formaldehydeDataRepository) GetHistoryByDeviceID(ctx context.Context, deviceID string, limit, offset int) ([]models.FormaldehydeData, error) {
	var data []models.FormaldehydeData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}

// GetByTimeRange 获取设备指定时间范围的数据
func (r *formaldehydeDataRepository) GetByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.FormaldehydeData, error) {
	var data []models.FormaldehydeData
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp BETWEEN ? AND ?",
			deviceID,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0)).
		Order("timestamp ASC").
		Find(&data).Error
	return data, err
}

// GetStatistics 获取设备统计数据
func (r *formaldehydeDataRepository) GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.FormaldehydeStatistics, error) {
	var stats models.FormaldehydeStatistics

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&models.FormaldehydeData{}).
		Where("device_id = ? AND timestamp BETWEEN ? AND ?",
			deviceID,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0))

	// 计算统计数据
	err := query.Select(`
		COUNT(*) as data_count,
		AVG(formaldehyde) as formaldehyde_avg,
		MIN(formaldehyde) as formaldehyde_min,
		MAX(formaldehyde) as formaldehyde_max,
		AVG(temperature) as temperature_avg,
		MIN(temperature) as temperature_min,
		MAX(temperature) as temperature_max,
		AVG(humidity) as humidity_avg,
		MIN(humidity) as humidity_min,
		MAX(humidity) as humidity_max
	`).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// BatchInsert 批量插入数据
func (r *formaldehydeDataRepository) BatchInsert(ctx context.Context, data []models.FormaldehydeData) error {
	if len(data) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).CreateInBatches(data, 100).Error
}

// deviceStatusRepository 设备状态仓库实现
type deviceStatusRepository struct {
	*baseRepository[models.FormaldehydeDeviceStatus]
}

// NewDeviceStatusRepository 创建设备状态仓库
func NewDeviceStatusRepository(db *gorm.DB, logger utils.Logger) DeviceStatusRepository {
	return &deviceStatusRepository{
		baseRepository: NewBaseRepository[models.FormaldehydeDeviceStatus](db, logger).(*baseRepository[models.FormaldehydeDeviceStatus]),
	}
}

// GetByDeviceID 根据设备ID获取状态
func (r *deviceStatusRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.FormaldehydeDeviceStatus, error) {
	var status models.FormaldehydeDeviceStatus
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetOnlineDevices 获取在线设备列表
func (r *deviceStatusRepository) GetOnlineDevices(ctx context.Context) ([]models.FormaldehydeDeviceStatus, error) {
	var devices []models.FormaldehydeDeviceStatus
	err := r.db.WithContext(ctx).
		Where("online = ?", true).
		Find(&devices).Error
	return devices, err
}

// GetOfflineDevices 获取离线设备列表
func (r *deviceStatusRepository) GetOfflineDevices(ctx context.Context) ([]models.FormaldehydeDeviceStatus, error) {
	var devices []models.FormaldehydeDeviceStatus
	err := r.db.WithContext(ctx).
		Where("online = ?", false).
		Find(&devices).Error
	return devices, err
}

// UpdateOnlineStatus 更新设备在线状态
func (r *deviceStatusRepository) UpdateOnlineStatus(ctx context.Context, deviceID string, online bool) error {
	return r.db.WithContext(ctx).
		Model(&models.FormaldehydeDeviceStatus{}).
		Where("device_id = ?", deviceID).
		Update("online", online).Error
}

// BatchUpdateStatus 批量更新设备状态
func (r *deviceStatusRepository) BatchUpdateStatus(ctx context.Context, updates map[string]bool) error {
	if len(updates) == 0 {
		return nil
	}

	// 使用事务批量更新
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for deviceID, online := range updates {
			if err := tx.Model(&models.FormaldehydeDeviceStatus{}).
				Where("device_id = ?", deviceID).
				Update("online", online).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// deviceConfigRepository 设备配置仓库实现
type deviceConfigRepository struct {
	*baseRepository[models.FormaldehydeDeviceConfig]
}

// NewDeviceConfigRepository 创建设备配置仓库
func NewDeviceConfigRepository(db *gorm.DB, logger utils.Logger) DeviceConfigRepository {
	return &deviceConfigRepository{
		baseRepository: NewBaseRepository[models.FormaldehydeDeviceConfig](db, logger).(*baseRepository[models.FormaldehydeDeviceConfig]),
	}
}

// GetByDeviceID 根据设备ID获取配置
func (r *deviceConfigRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.FormaldehydeDeviceConfig, error) {
	var config models.FormaldehydeDeviceConfig
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateConfig 更新设备配置
func (r *deviceConfigRepository) UpdateConfig(ctx context.Context, deviceID string, config *models.FormaldehydeDeviceConfig) error {
	return r.db.WithContext(ctx).
		Model(&models.FormaldehydeDeviceConfig{}).
		Where("device_id = ?", deviceID).
		Updates(config).Error
}
