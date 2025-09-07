package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

// SensorDataRepository 传感器数据仓库接口
type SensorDataRepository interface {
	BaseRepository[models.SensorData]

	// 获取设备最新数据
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.SensorData, error)

	// 获取设备历史数据
	GetHistoryByDeviceID(ctx context.Context, deviceID string, limit, offset int) ([]models.SensorData, error)

	// 获取设备指定时间范围的数据
	GetByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.SensorData, error)

	// 获取设备统计数据
	GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.SensorDataStatistics, error)

	// 批量插入数据
	BatchInsert(ctx context.Context, data []models.SensorData) error

	// 根据设备类型获取数据
	GetByDeviceType(ctx context.Context, deviceType models.DeviceType, limit, offset int) ([]models.SensorData, error)

	// 获取指定指标的数据
	GetMetricData(ctx context.Context, deviceID, metric string, startTime, endTime int64) ([]models.SensorData, error)
}

// DeviceRuntimeStatusRepository 设备运行时状态仓库接口
type DeviceRuntimeStatusRepository interface {
	BaseRepository[models.DeviceRuntimeStatus]

	// 根据设备ID获取状态
	GetByDeviceID(ctx context.Context, deviceID string) (*models.DeviceRuntimeStatus, error)

	// 获取在线设备列表
	GetOnlineDevices(ctx context.Context) ([]models.DeviceRuntimeStatus, error)

	// 获取离线设备列表
	GetOfflineDevices(ctx context.Context) ([]models.DeviceRuntimeStatus, error)

	// 更新设备在线状态
	UpdateOnlineStatus(ctx context.Context, deviceID string, online bool) error

	// 批量更新设备状态
	BatchUpdateStatus(ctx context.Context, updates map[string]bool) error

	// 根据设备类型获取状态
	GetByDeviceType(ctx context.Context, deviceType models.DeviceType) ([]models.DeviceRuntimeStatus, error)
}

// sensorDataRepository 传感器数据仓库实现
type sensorDataRepository struct {
	*baseRepository[models.SensorData]
}

// NewSensorDataRepository 创建传感器数据仓库
func NewSensorDataRepository(db *gorm.DB, logger utils.Logger) SensorDataRepository {
	return &sensorDataRepository{
		baseRepository: NewBaseRepository[models.SensorData](db, logger).(*baseRepository[models.SensorData]),
	}
}

// GetLatestByDeviceID 获取设备最新数据
func (r *sensorDataRepository) GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.SensorData, error) {
	var data models.SensorData
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
func (r *sensorDataRepository) GetHistoryByDeviceID(ctx context.Context, deviceID string, limit, offset int) ([]models.SensorData, error) {
	var data []models.SensorData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}

// GetByTimeRange 获取设备指定时间范围的数据
func (r *sensorDataRepository) GetByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.SensorData, error) {
	var data []models.SensorData
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
func (r *sensorDataRepository) GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.SensorDataStatistics, error) {
	var stats models.SensorDataStatistics

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&models.SensorData{}).
		Where("device_id = ? AND timestamp BETWEEN ? AND ?",
			deviceID,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0))

	// 计算统计数据
	err := query.Select(`
		COUNT(*) as data_count,
		AVG(pm25) as pm25_avg,
		MIN(pm25) as pm25_min,
		MAX(pm25) as pm25_max,
		AVG(pm10) as pm10_avg,
		MIN(pm10) as pm10_min,
		MAX(pm10) as pm10_max,
		AVG(co2) as co2_avg,
		MIN(co2) as co2_min,
		MAX(co2) as co2_max,
		AVG(formaldehyde) as formaldehyde_avg,
		MIN(formaldehyde) as formaldehyde_min,
		MAX(formaldehyde) as formaldehyde_max,
		AVG(temperature) as temperature_avg,
		MIN(temperature) as temperature_min,
		MAX(temperature) as temperature_max,
		AVG(humidity) as humidity_avg,
		MIN(humidity) as humidity_min,
		MAX(humidity) as humidity_max,
		AVG(pressure) as pressure_avg,
		MIN(pressure) as pressure_min,
		MAX(pressure) as pressure_max
	`).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// BatchInsert 批量插入数据
func (r *sensorDataRepository) BatchInsert(ctx context.Context, data []models.SensorData) error {
	if len(data) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).CreateInBatches(data, 100).Error
}

// GetByDeviceType 根据设备类型获取数据
func (r *sensorDataRepository) GetByDeviceType(ctx context.Context, deviceType models.DeviceType, limit, offset int) ([]models.SensorData, error) {
	var data []models.SensorData
	err := r.db.WithContext(ctx).
		Where("device_type = ?", deviceType).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}

// GetMetricData 获取指定指标的数据
func (r *sensorDataRepository) GetMetricData(ctx context.Context, deviceID, metric string, startTime, endTime int64) ([]models.SensorData, error) {
	var data []models.SensorData
	query := r.db.WithContext(ctx).
		Where("device_id = ? AND timestamp BETWEEN ? AND ?",
			deviceID,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0))

	// 只选择包含指定指标的记录
	switch metric {
	case "pm25":
		query = query.Where("pm25 IS NOT NULL")
	case "pm10":
		query = query.Where("pm10 IS NOT NULL")
	case "co2":
		query = query.Where("co2 IS NOT NULL")
	case "formaldehyde":
		query = query.Where("formaldehyde IS NOT NULL")
	case "temperature":
		query = query.Where("temperature IS NOT NULL")
	case "humidity":
		query = query.Where("humidity IS NOT NULL")
	case "pressure":
		query = query.Where("pressure IS NOT NULL")
	}

	err := query.Order("timestamp ASC").Find(&data).Error
	return data, err
}

// deviceRuntimeStatusRepository 设备运行时状态仓库实现
type deviceRuntimeStatusRepository struct {
	*baseRepository[models.DeviceRuntimeStatus]
}

// NewDeviceRuntimeStatusRepository 创建设备运行时状态仓库
func NewDeviceRuntimeStatusRepository(db *gorm.DB, logger utils.Logger) DeviceRuntimeStatusRepository {
	return &deviceRuntimeStatusRepository{
		baseRepository: NewBaseRepository[models.DeviceRuntimeStatus](db, logger).(*baseRepository[models.DeviceRuntimeStatus]),
	}
}

// GetByDeviceID 根据设备ID获取状态
func (r *deviceRuntimeStatusRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.DeviceRuntimeStatus, error) {
	var status models.DeviceRuntimeStatus
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetOnlineDevices 获取在线设备列表
func (r *deviceRuntimeStatusRepository) GetOnlineDevices(ctx context.Context) ([]models.DeviceRuntimeStatus, error) {
	var devices []models.DeviceRuntimeStatus
	err := r.db.WithContext(ctx).
		Where("online = ?", true).
		Find(&devices).Error
	return devices, err
}

// GetOfflineDevices 获取离线设备列表
func (r *deviceRuntimeStatusRepository) GetOfflineDevices(ctx context.Context) ([]models.DeviceRuntimeStatus, error) {
	var devices []models.DeviceRuntimeStatus
	err := r.db.WithContext(ctx).
		Where("online = ?", false).
		Find(&devices).Error
	return devices, err
}

// UpdateOnlineStatus 更新设备在线状态
func (r *deviceRuntimeStatusRepository) UpdateOnlineStatus(ctx context.Context, deviceID string, online bool) error {
	return r.db.WithContext(ctx).
		Model(&models.DeviceRuntimeStatus{}).
		Where("device_id = ?", deviceID).
		Update("online", online).Error
}

// BatchUpdateStatus 批量更新设备状态
func (r *deviceRuntimeStatusRepository) BatchUpdateStatus(ctx context.Context, updates map[string]bool) error {
	if len(updates) == 0 {
		return nil
	}

	// 使用事务批量更新
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for deviceID, online := range updates {
			if err := tx.Model(&models.DeviceRuntimeStatus{}).
				Where("device_id = ?", deviceID).
				Update("online", online).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByDeviceType 根据设备类型获取状态
func (r *deviceRuntimeStatusRepository) GetByDeviceType(ctx context.Context, deviceType models.DeviceType) ([]models.DeviceRuntimeStatus, error) {
	var devices []models.DeviceRuntimeStatus
	err := r.db.WithContext(ctx).
		Joins("JOIN devices ON device_runtime_status.device_id = devices.id").
		Where("devices.type = ?", deviceType).
		Find(&devices).Error
	return devices, err
}
