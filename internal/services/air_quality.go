package services

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/utils"
	"context"
)

// AirQualityService 空气质量服务接口
type AirQualityService interface {
	CreateData(ctx context.Context, data *models.AirQualityData) error
	CreateBatchData(ctx context.Context, data []models.AirQualityData) error
	GetDataByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.AirQualityData, error)
	GetDataByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.AirQualityData, error)
	GetLatestData(ctx context.Context, deviceID string) (*models.AirQualityData, error)
	GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.AirQualityStatistics, error)
	AnalyzeData(ctx context.Context, data *models.AirQualityData) (map[string]interface{}, error)
}

// airQualityService 空气质量服务实现
type airQualityService struct {
	airQualityRepo repositories.AirQualityRepository
	deviceRepo     repositories.DeviceRepository
	logger         utils.Logger
}

// NewAirQualityService 创建空气质量服务
func NewAirQualityService(airQualityRepo repositories.AirQualityRepository, deviceRepo repositories.DeviceRepository, logger utils.Logger) AirQualityService {
	return &airQualityService{
		airQualityRepo: airQualityRepo,
		deviceRepo:     deviceRepo,
		logger:         logger,
	}
}

// CreateData 创建空气质量数据
func (s *airQualityService) CreateData(ctx context.Context, data *models.AirQualityData) error {
	if err := s.airQualityRepo.Create(ctx, data); err != nil {
		s.logger.Error("创建空气质量数据失败", utils.ErrorField(err))
		return err
	}

	// 分析数据并记录
	if analysis, err := s.AnalyzeData(ctx, data); err == nil {
		s.logger.Info("空气质量数据分析完成",
			utils.String("device_id", data.DeviceID),
			utils.String("quality", analysis["quality"].(string)),
			utils.Float64("aqi", analysis["aqi"].(float64)))
	}

	return nil
}

// CreateBatchData 批量创建空气质量数据
func (s *airQualityService) CreateBatchData(ctx context.Context, data []models.AirQualityData) error {
	for _, item := range data {
		if err := s.airQualityRepo.Create(ctx, &item); err != nil {
			s.logger.Error("批量创建空气质量数据失败", utils.ErrorField(err))
			return err
		}
	}

	s.logger.Info("批量创建空气质量数据成功", utils.Int("count", len(data)))
	return nil
}

// GetDataByDeviceID 根据设备ID获取数据
func (s *airQualityService) GetDataByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.AirQualityData, error) {
	data, err := s.airQualityRepo.GetByDeviceID(deviceID, limit)
	if err != nil {
		s.logger.Error("根据设备ID获取空气质量数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return data, nil
}

// GetDataByTimeRange 根据时间范围获取数据
func (s *airQualityService) GetDataByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.AirQualityData, error) {
	data, err := s.airQualityRepo.GetByTimeRange(deviceID, startTime, endTime)
	if err != nil {
		s.logger.Error("根据时间范围获取空气质量数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return data, nil
}

// GetLatestData 获取最新数据
func (s *airQualityService) GetLatestData(ctx context.Context, deviceID string) (*models.AirQualityData, error) {
	data, err := s.airQualityRepo.GetLatestByDeviceID(deviceID)
	if err != nil {
		s.logger.Error("获取最新空气质量数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return data, nil
}

// GetStatistics 获取统计数据
func (s *airQualityService) GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.AirQualityStatistics, error) {
	stats, err := s.airQualityRepo.GetStatistics(deviceID, startTime, endTime)
	if err != nil {
		s.logger.Error("获取空气质量统计数据失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return stats, nil
}

// AnalyzeData 分析空气质量数据
func (s *airQualityService) AnalyzeData(ctx context.Context, data *models.AirQualityData) (map[string]interface{}, error) {
	analysis := make(map[string]interface{})
	analysis["data_id"] = data.ID
	analysis["device_id"] = data.DeviceID
	analysis["timestamp"] = data.CreatedAt

	// 计算AQI（简化版本）
	var pm25, pm10 float64
	if data.PM25 != nil {
		pm25 = *data.PM25
	}
	if data.PM10 != nil {
		pm10 = *data.PM10
	}
	aqi := s.calculateAQI(pm25, pm10)
	analysis["aqi"] = aqi
	analysis["quality"] = s.getQualityLevel(aqi)

	// 健康建议
	analysis["health_advice"] = s.getHealthAdvice(analysis["quality"].(string))

	// 趋势分析（简化版本）
	analysis["trend"] = s.analyzeTrend(data)

	return analysis, nil
}

// calculateAQI 计算AQI
func (s *airQualityService) calculateAQI(pm2_5, pm10 float64) float64 {
	// 简化的AQI计算，实际应该根据国家标准计算
	aqi := (pm2_5*2.5 + pm10*1.5) / 4
	if aqi > 500 {
		return 500
	}
	return aqi
}

// getQualityLevel 获取空气质量等级
func (s *airQualityService) getQualityLevel(aqi float64) string {
	if aqi <= 50 {
		return "优"
	} else if aqi <= 100 {
		return "良"
	} else if aqi <= 150 {
		return "轻度污染"
	} else if aqi <= 200 {
		return "中度污染"
	} else if aqi <= 300 {
		return "重度污染"
	} else {
		return "严重污染"
	}
}

// getHealthAdvice 获取健康建议
func (s *airQualityService) getHealthAdvice(qualityLevel string) string {
	switch qualityLevel {
	case "优":
		return "空气质量令人满意，基本无空气污染，各类人群可正常活动。"
	case "良":
		return "空气质量可以接受，但某些污染物可能对极少数异常敏感人群健康有较弱影响。"
	case "轻度污染":
		return "易感人群症状有轻度加剧，健康人群出现刺激症状。建议儿童、老年人及心脏病、呼吸系统疾病患者应减少长时间、高强度的户外锻炼。"
	case "中度污染":
		return "进一步加剧易感人群症状，可能对健康人群心脏、呼吸系统有影响。建议儿童、老年人及心脏病、呼吸系统疾病患者避免长时间、高强度的户外锻炼。"
	case "重度污染":
		return "心脏病和肺病患者症状显著加剧，运动耐受力降低，健康人群普遍出现症状。建议儿童、老年人和心脏病、肺病患者应停留在室内，停止户外运动。"
	case "严重污染":
		return "健康人群运动耐受力降低，有明显强烈症状，提前出现某些疾病。建议儿童、老年人和病人应当停留在室内，避免体力消耗，一般人群应避免户外活动。"
	default:
		return "建议关注空气质量变化，做好防护措施。"
	}
}

// analyzeTrend 分析趋势
func (s *airQualityService) analyzeTrend(data *models.AirQualityData) string {
	// 简化的趋势分析，实际应该对比历史数据
	var pm25, pm10 float64
	if data.PM25 != nil {
		pm25 = *data.PM25
	}
	if data.PM10 != nil {
		pm10 = *data.PM10
	}

	if pm25 > 75 || pm10 > 150 {
		return "上升"
	} else if pm25 < 25 && pm10 < 50 {
		return "下降"
	} else {
		return "稳定"
	}
}
