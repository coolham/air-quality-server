package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

// UnifiedSensorDataRepository 统一传感器数据仓库接口
type UnifiedSensorDataRepository interface {
	BaseRepository[models.UnifiedSensorData]

	// 获取设备最新数据
	GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.UnifiedSensorData, error)

	// 获取设备历史数据
	GetHistoryByDeviceID(ctx context.Context, deviceID string, limit, offset int) ([]models.UnifiedSensorData, error)

	// 获取设备指定时间范围的数据
	GetByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.UnifiedSensorData, error)

	// 获取设备统计数据
	GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error)

	// 批量插入数据
	BatchInsert(ctx context.Context, data []models.UnifiedSensorData) error

	// 根据设备类型获取数据
	GetByDeviceType(ctx context.Context, deviceType models.DeviceType, limit, offset int) ([]models.UnifiedSensorData, error)

	// 获取指定指标的数据
	GetMetricData(ctx context.Context, deviceID, metric string, startTime, endTime int64) ([]models.UnifiedSensorData, error)

	// 获取多设备数据
	GetMultiDeviceData(ctx context.Context, deviceIDs []string, startTime, endTime int64) ([]models.UnifiedSensorData, error)

	// 获取设备类型统计
	GetDeviceTypeStatistics(ctx context.Context, deviceType models.DeviceType, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error)

	// 数据迁移相关方法
	MigrateFromAirQualityData(ctx context.Context, data []models.AirQualityData) error
	MigrateFromFormaldehydeData(ctx context.Context, data []models.FormaldehydeData) error

	// 获取传感器ID列表
	GetSensorIDs(ctx context.Context, deviceID string) ([]string, error)

	// 获取所有设备数据
	GetAllData(ctx context.Context, limit, offset int) ([]models.UnifiedSensorData, error)
}

// unifiedSensorDataRepository 统一传感器数据仓库实现
type unifiedSensorDataRepository struct {
	*baseRepository[models.UnifiedSensorData]
}

// NewUnifiedSensorDataRepository 创建统一传感器数据仓库
func NewUnifiedSensorDataRepository(db *gorm.DB, logger utils.Logger) UnifiedSensorDataRepository {
	return &unifiedSensorDataRepository{
		baseRepository: NewBaseRepository[models.UnifiedSensorData](db, logger).(*baseRepository[models.UnifiedSensorData]),
	}
}

// GetLatestByDeviceID 获取设备最新数据
func (r *unifiedSensorDataRepository) GetLatestByDeviceID(ctx context.Context, deviceID string) (*models.UnifiedSensorData, error) {
	var data models.UnifiedSensorData
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
func (r *unifiedSensorDataRepository) GetHistoryByDeviceID(ctx context.Context, deviceID string, limit, offset int) ([]models.UnifiedSensorData, error) {
	var data []models.UnifiedSensorData
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}

// GetByTimeRange 获取设备指定时间范围的数据
func (r *unifiedSensorDataRepository) GetByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.UnifiedSensorData, error) {
	var data []models.UnifiedSensorData
	query := r.db.WithContext(ctx).Where("timestamp BETWEEN ? AND ?",
		time.Unix(startTime, 0),
		time.Unix(endTime, 0))

	// 如果指定了设备ID，则添加设备筛选条件
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	err := query.Order("timestamp ASC").Find(&data).Error
	return data, err
}

// GetStatistics 获取设备统计数据
func (r *unifiedSensorDataRepository) GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error) {
	var stats models.UnifiedSensorDataStatistics

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&models.UnifiedSensorData{}).
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
func (r *unifiedSensorDataRepository) BatchInsert(ctx context.Context, data []models.UnifiedSensorData) error {
	if len(data) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).CreateInBatches(data, 100).Error
}

// GetByDeviceType 根据设备类型获取数据
func (r *unifiedSensorDataRepository) GetByDeviceType(ctx context.Context, deviceType models.DeviceType, limit, offset int) ([]models.UnifiedSensorData, error) {
	var data []models.UnifiedSensorData
	err := r.db.WithContext(ctx).
		Where("device_type = ?", deviceType).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}

// GetMetricData 获取指定指标的数据
func (r *unifiedSensorDataRepository) GetMetricData(ctx context.Context, deviceID, metric string, startTime, endTime int64) ([]models.UnifiedSensorData, error) {
	var data []models.UnifiedSensorData
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
	case "o3":
		query = query.Where("o3 IS NOT NULL")
	case "no2":
		query = query.Where("no2 IS NOT NULL")
	case "so2":
		query = query.Where("so2 IS NOT NULL")
	case "co":
		query = query.Where("co IS NOT NULL")
	case "voc":
		query = query.Where("voc IS NOT NULL")
	}

	err := query.Order("timestamp ASC").Find(&data).Error
	return data, err
}

// GetMultiDeviceData 获取多设备数据
func (r *unifiedSensorDataRepository) GetMultiDeviceData(ctx context.Context, deviceIDs []string, startTime, endTime int64) ([]models.UnifiedSensorData, error) {
	var data []models.UnifiedSensorData
	err := r.db.WithContext(ctx).
		Where("device_id IN ? AND timestamp BETWEEN ? AND ?",
			deviceIDs,
			time.Unix(startTime, 0),
			time.Unix(endTime, 0)).
		Order("device_id, timestamp ASC").
		Find(&data).Error
	return data, err
}

// GetDeviceTypeStatistics 获取设备类型统计
func (r *unifiedSensorDataRepository) GetDeviceTypeStatistics(ctx context.Context, deviceType models.DeviceType, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error) {
	var stats models.UnifiedSensorDataStatistics

	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&models.UnifiedSensorData{}).
		Where("device_type = ? AND timestamp BETWEEN ? AND ?",
			deviceType,
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

// MigrateFromAirQualityData 从AirQualityData迁移数据
func (r *unifiedSensorDataRepository) MigrateFromAirQualityData(ctx context.Context, data []models.AirQualityData) error {
	if len(data) == 0 {
		return nil
	}

	var unifiedData []models.UnifiedSensorData
	for _, item := range data {
		unified := models.UnifiedSensorData{
			DeviceID:    item.DeviceID,
			DeviceType:  models.DeviceTypeAirQuality, // 默认为综合空气质量传感器
			Timestamp:   item.Timestamp,
			PM25:        item.PM25,
			PM10:        item.PM10,
			CO2:         item.CO2,
			Temperature: item.Temperature,
			Humidity:    item.Humidity,
			Pressure:    item.Pressure,
			DataQuality: "good",
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
		unifiedData = append(unifiedData, unified)
	}

	return r.BatchInsert(ctx, unifiedData)
}

// MigrateFromFormaldehydeData 从FormaldehydeData迁移数据
func (r *unifiedSensorDataRepository) MigrateFromFormaldehydeData(ctx context.Context, data []models.FormaldehydeData) error {
	if len(data) == 0 {
		return nil
	}

	var unifiedData []models.UnifiedSensorData
	for _, item := range data {
		unified := models.UnifiedSensorData{
			DeviceID:       item.DeviceID,
			DeviceType:     models.DeviceTypeFormaldehyde,
			Timestamp:      item.Timestamp,
			Formaldehyde:   item.Formaldehyde,
			Temperature:    item.Temperature,
			Humidity:       item.Humidity,
			Battery:        item.Battery,
			SignalStrength: item.SignalStrength,
			DataQuality:    item.DataQuality,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
		unifiedData = append(unifiedData, unified)
	}

	return r.BatchInsert(ctx, unifiedData)
}

// GetSensorIDs 获取传感器ID列表
func (r *unifiedSensorDataRepository) GetSensorIDs(ctx context.Context, deviceID string) ([]string, error) {
	var sensorIDs []string

	query := r.db.WithContext(ctx).Model(&models.UnifiedSensorData{}).
		Select("DISTINCT sensor_id").
		Where("sensor_id IS NOT NULL AND sensor_id != ''")

	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	err := query.Pluck("sensor_id", &sensorIDs).Error
	if err != nil {
		return nil, err
	}

	return sensorIDs, nil
}

// GetAllData 获取所有设备数据
func (r *unifiedSensorDataRepository) GetAllData(ctx context.Context, limit, offset int) ([]models.UnifiedSensorData, error) {
	var data []models.UnifiedSensorData
	err := r.db.WithContext(ctx).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&data).Error
	return data, err
}
