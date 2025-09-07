package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

// AirQualityRepository 空气质量数据仓储接口
type AirQualityRepository interface {
	BaseRepository[models.AirQualityData]
	GetByDeviceID(deviceID string, limit int) ([]models.AirQualityData, error)
	GetByTimeRange(deviceID string, startTime, endTime int64) ([]models.AirQualityData, error)
	GetLatestByDeviceID(deviceID string) (*models.AirQualityData, error)
	GetStatistics(deviceID string, startTime, endTime int64) (*models.AirQualityStatistics, error)
}

// airQualityRepository 空气质量数据仓储实现
type airQualityRepository struct {
	*baseRepository[models.AirQualityData]
	db     *gorm.DB
	logger utils.Logger
}

// NewAirQualityRepository 创建空气质量数据仓储
func NewAirQualityRepository(db *gorm.DB, logger utils.Logger) AirQualityRepository {
	return &airQualityRepository{
		baseRepository: NewBaseRepository[models.AirQualityData](db, logger).(*baseRepository[models.AirQualityData]),
		db:             db,
		logger:         logger,
	}
}

// GetByDeviceID 根据设备ID获取数据
func (r *airQualityRepository) GetByDeviceID(deviceID string, limit int) ([]models.AirQualityData, error) {
	var data []models.AirQualityData
	err := r.db.Where(&models.AirQualityData{DeviceID: deviceID}).
		Order("created_at DESC").
		Limit(limit).
		Find(&data).Error
	if err != nil {
		r.logger.Error("获取设备空气质量数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return data, nil
}

// GetByTimeRange 根据时间范围获取数据
func (r *airQualityRepository) GetByTimeRange(deviceID string, startTime, endTime int64) ([]models.AirQualityData, error) {
	var data []models.AirQualityData
	err := r.db.Where("device_id = ? AND created_at BETWEEN ? AND ?", deviceID, startTime, endTime).
		Order("created_at ASC").
		Find(&data).Error
	if err != nil {
		r.logger.Error("获取时间范围内空气质量数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return data, nil
}

// GetLatestByDeviceID 获取设备最新数据
func (r *airQualityRepository) GetLatestByDeviceID(deviceID string) (*models.AirQualityData, error) {
	var data models.AirQualityData
	err := r.db.Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("获取设备最新空气质量数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return &data, nil
}

// GetStatistics 获取统计数据
func (r *airQualityRepository) GetStatistics(deviceID string, startTime, endTime int64) (*models.AirQualityStatistics, error) {
	var stats models.AirQualityStatistics

	// 获取平均值
	err := r.db.Model(&models.AirQualityData{}).
		Where("device_id = ? AND created_at BETWEEN ? AND ?", deviceID, startTime, endTime).
		Select("AVG(pm25) as pm25_avg, AVG(pm10) as pm10_avg, AVG(co2) as co2_avg, AVG(temperature) as temperature_avg, AVG(humidity) as humidity_avg").
		Scan(&stats).Error
	if err != nil {
		r.logger.Error("获取空气质量统计数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}

	// 获取最大值
	var maxData models.AirQualityData
	err = r.db.Where("device_id = ? AND created_at BETWEEN ? AND ?", deviceID, startTime, endTime).
		Order("pm25 DESC").
		First(&maxData).Error
	if err == nil {
		if maxData.PM25 != nil {
			stats.PM25Max = maxData.PM25
		}
		if maxData.PM10 != nil {
			stats.PM10Max = maxData.PM10
		}
		if maxData.CO2 != nil {
			stats.CO2Max = maxData.CO2
		}
		if maxData.Temperature != nil {
			stats.TempMax = maxData.Temperature
		}
		if maxData.Humidity != nil {
			stats.HumidityMax = maxData.Humidity
		}
	}

	// 获取最小值
	var minData models.AirQualityData
	err = r.db.Where("device_id = ? AND created_at BETWEEN ? AND ?", deviceID, startTime, endTime).
		Order("pm25 ASC").
		First(&minData).Error
	if err == nil {
		if minData.PM25 != nil {
			stats.PM25Min = minData.PM25
		}
		if minData.PM10 != nil {
			stats.PM10Min = minData.PM10
		}
		if minData.CO2 != nil {
			stats.CO2Min = minData.CO2
		}
		if minData.Temperature != nil {
			stats.TempMin = minData.Temperature
		}
		if minData.Humidity != nil {
			stats.HumidityMin = minData.Humidity
		}
	}

	// 获取数据点数量
	var count int64
	err = r.db.Model(&models.AirQualityData{}).
		Where("device_id = ? AND created_at BETWEEN ? AND ?", deviceID, startTime, endTime).
		Count(&count).Error
	if err != nil {
		r.logger.Error("获取数据点数量失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	stats.DataCount = count

	return &stats, nil
}
