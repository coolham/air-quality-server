package services

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// UnifiedSensorDataService 统一传感器数据服务接口
type UnifiedSensorDataService interface {
	// 数据创建
	CreateData(ctx context.Context, data *models.UnifiedSensorData) error
	CreateBatchData(ctx context.Context, data []models.UnifiedSensorData) error
	CreateFromUpload(ctx context.Context, upload *models.UnifiedSensorDataUpload) error

	// 数据查询
	GetDataByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.UnifiedSensorData, error)
	GetDataByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.UnifiedSensorData, error)
	GetLatestData(ctx context.Context, deviceID string) (*models.UnifiedSensorData, error)
	GetDataByDeviceType(ctx context.Context, deviceType models.DeviceType, limit int) ([]models.UnifiedSensorData, error)
	GetMultiDeviceData(ctx context.Context, deviceIDs []string, startTime, endTime int64) ([]models.UnifiedSensorData, error)

	// 统计分析
	GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error)
	GetDeviceTypeStatistics(ctx context.Context, deviceType models.DeviceType, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error)
	GetMetricData(ctx context.Context, deviceID, metric string, startTime, endTime int64) ([]models.UnifiedSensorData, error)

	// 数据分析
	AnalyzeData(ctx context.Context, data *models.UnifiedSensorData) (map[string]interface{}, error)
	GetDataQualityScore(ctx context.Context, data *models.UnifiedSensorData) (float64, error)

	// 告警检查
	CheckAlerts(ctx context.Context, data *models.UnifiedSensorData) ([]models.Alert, error)
}

// unifiedSensorDataService 统一传感器数据服务实现
type unifiedSensorDataService struct {
	dataRepo   repositories.UnifiedSensorDataRepository
	deviceRepo repositories.DeviceRepository
	alertSvc   AlertService
	logger     utils.Logger
}

// NewUnifiedSensorDataService 创建统一传感器数据服务
func NewUnifiedSensorDataService(
	dataRepo repositories.UnifiedSensorDataRepository,
	deviceRepo repositories.DeviceRepository,
	alertSvc AlertService,
	logger utils.Logger,
) UnifiedSensorDataService {
	return &unifiedSensorDataService{
		dataRepo:   dataRepo,
		deviceRepo: deviceRepo,
		alertSvc:   alertSvc,
		logger:     logger,
	}
}

// CreateData 创建传感器数据
func (s *unifiedSensorDataService) CreateData(ctx context.Context, data *models.UnifiedSensorData) error {
	// 验证设备是否存在
	device, err := s.deviceRepo.GetByID(ctx, data.DeviceID)
	if err != nil {
		s.logger.Error("获取设备失败", utils.ErrorField(err), utils.String("device_id", data.DeviceID))
		return fmt.Errorf("获取设备失败: %w", err)
	}

	// 验证设备类型是否匹配
	if device.Type != data.DeviceType {
		s.logger.Warn("设备类型不匹配",
			utils.String("device_id", data.DeviceID),
			utils.String("expected_type", string(device.Type)),
			utils.String("actual_type", string(data.DeviceType)))
	}

	// 计算数据质量分数
	qualityScore, err := s.GetDataQualityScore(ctx, data)
	if err != nil {
		s.logger.Warn("计算数据质量分数失败", utils.ErrorField(err))
	} else {
		// 可以根据质量分数调整数据质量标记
		if qualityScore < 0.7 {
			data.DataQuality = "poor"
		} else if qualityScore < 0.9 {
			data.DataQuality = "fair"
		} else {
			data.DataQuality = "good"
		}
	}

	// 保存数据
	if err := s.dataRepo.Create(ctx, data); err != nil {
		s.logger.Error("创建传感器数据失败", utils.ErrorField(err), utils.String("device_id", data.DeviceID))
		return fmt.Errorf("创建传感器数据失败: %w", err)
	}

	s.logger.Info("创建传感器数据成功",
		utils.String("device_id", data.DeviceID),
		utils.String("device_type", string(data.DeviceType)),
		utils.String("data_quality", data.DataQuality))

	return nil
}

// CreateBatchData 批量创建传感器数据
func (s *unifiedSensorDataService) CreateBatchData(ctx context.Context, data []models.UnifiedSensorData) error {
	if len(data) == 0 {
		return nil
	}

	// 批量保存数据
	if err := s.dataRepo.BatchInsert(ctx, data); err != nil {
		s.logger.Error("批量创建传感器数据失败", utils.ErrorField(err))
		return fmt.Errorf("批量创建传感器数据失败: %w", err)
	}

	s.logger.Info("批量创建传感器数据成功", utils.Int("count", len(data)))
	return nil
}

// CreateFromUpload 从上传请求创建数据
func (s *unifiedSensorDataService) CreateFromUpload(ctx context.Context, upload *models.UnifiedSensorDataUpload) error {
	// 转换时间戳
	timestamp := time.Unix(upload.Timestamp, 0)

	// 确定设备类型
	deviceType := models.DeviceType(upload.DeviceType)
	if !deviceType.IsValid() {
		return fmt.Errorf("无效的设备类型: %s", upload.DeviceType)
	}

	// 创建传感器数据
	sensorData := &models.UnifiedSensorData{
		DeviceID:    upload.DeviceID,
		DeviceType:  deviceType,
		Timestamp:   timestamp,
		DataQuality: "good",
	}

	// 解析数据字段
	for metric, value := range upload.Data {
		if floatValue, ok := value.(float64); ok {
			sensorData.SetMetricValue(metric, &floatValue)
		}
	}

	// 解析位置信息
	if upload.Location != nil {
		sensorData.Latitude = upload.Location.Latitude
		sensorData.Longitude = upload.Location.Longitude
	}

	// 解析质量信息
	if upload.Quality != nil {
		sensorData.SignalStrength = upload.Quality.SignalStrength
		if upload.Quality.DataQuality != "" {
			sensorData.DataQuality = upload.Quality.DataQuality
		}
	}

	// 解析扩展数据
	if upload.Extended != nil && len(upload.Extended) > 0 {
		extendedJSON, err := json.Marshal(upload.Extended)
		if err == nil {
			extendedStr := string(extendedJSON)
			sensorData.ExtendedData = &extendedStr
		}
	}

	return s.CreateData(ctx, sensorData)
}

// GetDataByDeviceID 根据设备ID获取数据
func (s *unifiedSensorDataService) GetDataByDeviceID(ctx context.Context, deviceID string, limit int) ([]models.UnifiedSensorData, error) {
	return s.dataRepo.GetHistoryByDeviceID(ctx, deviceID, limit, 0)
}

// GetDataByTimeRange 根据时间范围获取数据
func (s *unifiedSensorDataService) GetDataByTimeRange(ctx context.Context, deviceID string, startTime, endTime int64) ([]models.UnifiedSensorData, error) {
	return s.dataRepo.GetByTimeRange(ctx, deviceID, startTime, endTime)
}

// GetLatestData 获取最新数据
func (s *unifiedSensorDataService) GetLatestData(ctx context.Context, deviceID string) (*models.UnifiedSensorData, error) {
	return s.dataRepo.GetLatestByDeviceID(ctx, deviceID)
}

// GetDataByDeviceType 根据设备类型获取数据
func (s *unifiedSensorDataService) GetDataByDeviceType(ctx context.Context, deviceType models.DeviceType, limit int) ([]models.UnifiedSensorData, error) {
	return s.dataRepo.GetByDeviceType(ctx, deviceType, limit, 0)
}

// GetMultiDeviceData 获取多设备数据
func (s *unifiedSensorDataService) GetMultiDeviceData(ctx context.Context, deviceIDs []string, startTime, endTime int64) ([]models.UnifiedSensorData, error) {
	return s.dataRepo.GetMultiDeviceData(ctx, deviceIDs, startTime, endTime)
}

// GetStatistics 获取统计数据
func (s *unifiedSensorDataService) GetStatistics(ctx context.Context, deviceID string, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error) {
	return s.dataRepo.GetStatistics(ctx, deviceID, startTime, endTime)
}

// GetDeviceTypeStatistics 获取设备类型统计
func (s *unifiedSensorDataService) GetDeviceTypeStatistics(ctx context.Context, deviceType models.DeviceType, startTime, endTime int64) (*models.UnifiedSensorDataStatistics, error) {
	return s.dataRepo.GetDeviceTypeStatistics(ctx, deviceType, startTime, endTime)
}

// GetMetricData 获取指定指标数据
func (s *unifiedSensorDataService) GetMetricData(ctx context.Context, deviceID, metric string, startTime, endTime int64) ([]models.UnifiedSensorData, error) {
	return s.dataRepo.GetMetricData(ctx, deviceID, metric, startTime, endTime)
}

// AnalyzeData 分析数据
func (s *unifiedSensorDataService) AnalyzeData(ctx context.Context, data *models.UnifiedSensorData) (map[string]interface{}, error) {
	analysis := make(map[string]interface{})

	// 获取可用指标
	availableMetrics := data.GetAvailableMetrics()
	analysis["available_metrics"] = availableMetrics

	// 根据设备类型进行特定分析
	switch data.DeviceType {
	case models.DeviceTypeFormaldehyde:
		if data.Formaldehyde != nil {
			analysis["formaldehyde_level"] = *data.Formaldehyde
			analysis["formaldehyde_status"] = s.getFormaldehydeStatus(*data.Formaldehyde)
		}

	case models.DeviceTypePM25:
		if data.PM25 != nil {
			analysis["pm25_level"] = *data.PM25
			analysis["pm25_status"] = s.getPM25Status(*data.PM25)
		}

	case models.DeviceTypeAirQuality:
		// 综合空气质量分析
		analysis["air_quality_index"] = s.calculateAirQualityIndex(data)
	}

	// 环境参数分析
	if data.Temperature != nil {
		analysis["temperature_status"] = s.getTemperatureStatus(*data.Temperature)
	}
	if data.Humidity != nil {
		analysis["humidity_status"] = s.getHumidityStatus(*data.Humidity)
	}

	return analysis, nil
}

// GetDataQualityScore 获取数据质量分数
func (s *unifiedSensorDataService) GetDataQualityScore(ctx context.Context, data *models.UnifiedSensorData) (float64, error) {
	score := 1.0

	// 检查数据完整性
	availableMetrics := data.GetAvailableMetrics()
	expectedMetrics := data.DeviceType.GetSupportedMetrics()

	completeness := float64(len(availableMetrics)) / float64(len(expectedMetrics))
	score *= completeness

	// 检查数据合理性
	if data.Temperature != nil && (*data.Temperature < -50 || *data.Temperature > 60) {
		score *= 0.5 // 温度异常
	}
	if data.Humidity != nil && (*data.Humidity < 0 || *data.Humidity > 100) {
		score *= 0.5 // 湿度异常
	}

	// 检查设备状态
	if data.Battery != nil && *data.Battery < 10 {
		score *= 0.8 // 电池电量低
	}
	if data.SignalStrength != nil && *data.SignalStrength < -100 {
		score *= 0.9 // 信号强度弱
	}

	return score, nil
}

// CheckAlerts 检查告警
func (s *unifiedSensorDataService) CheckAlerts(ctx context.Context, data *models.UnifiedSensorData) ([]models.Alert, error) {
	var alerts []models.Alert

	// 根据设备类型检查相应的告警
	switch data.DeviceType {
	case models.DeviceTypeFormaldehyde:
		if data.Formaldehyde != nil {
			alert := s.checkFormaldehydeAlert(ctx, data)
			if alert != nil {
				alerts = append(alerts, *alert)
			}
		}

	case models.DeviceTypePM25:
		if data.PM25 != nil {
			alert := s.checkPM25Alert(ctx, data)
			if alert != nil {
				alerts = append(alerts, *alert)
			}
		}
	}

	// 检查设备状态告警
	if data.Battery != nil && *data.Battery < 20 {
		message := fmt.Sprintf("设备电池电量低: %d%%", *data.Battery)
		alert := &models.Alert{
			RuleID:         0, // 临时规则ID
			DeviceID:       data.DeviceID,
			Metric:         "battery",
			CurrentValue:   float64(*data.Battery),
			ThresholdValue: 20,
			Severity:       "warning",
			Status:         "active",
			TriggeredAt:    time.Now(),
			Message:        &message,
			CreatedAt:      time.Now(),
		}
		alerts = append(alerts, *alert)
	}

	return alerts, nil
}

// 辅助方法
func (s *unifiedSensorDataService) getFormaldehydeStatus(value float64) string {
	if value >= 0.1 {
		return "严重超标"
	} else if value >= 0.08 {
		return "超标"
	} else {
		return "正常"
	}
}

func (s *unifiedSensorDataService) getPM25Status(value float64) string {
	if value >= 75 {
		return "严重污染"
	} else if value >= 35 {
		return "轻度污染"
	} else {
		return "良好"
	}
}

func (s *unifiedSensorDataService) getTemperatureStatus(value float64) string {
	if value < 0 || value > 40 {
		return "异常"
	} else if value < 10 || value > 30 {
		return "偏高/偏低"
	} else {
		return "正常"
	}
}

func (s *unifiedSensorDataService) getHumidityStatus(value float64) string {
	if value < 20 || value > 80 {
		return "异常"
	} else if value < 30 || value > 70 {
		return "偏高/偏低"
	} else {
		return "正常"
	}
}

func (s *unifiedSensorDataService) calculateAirQualityIndex(data *models.UnifiedSensorData) float64 {
	// 简化的空气质量指数计算
	index := 0.0
	count := 0

	if data.PM25 != nil {
		index += *data.PM25
		count++
	}
	if data.PM10 != nil {
		index += *data.PM10 * 0.5
		count++
	}
	if data.CO2 != nil {
		index += *data.CO2 * 0.1
		count++
	}

	if count > 0 {
		return index / float64(count)
	}
	return 0
}

func (s *unifiedSensorDataService) checkFormaldehydeAlert(ctx context.Context, data *models.UnifiedSensorData) *models.Alert {
	formaldehyde := *data.Formaldehyde
	var alertLevel string
	var message string
	var threshold float64

	if formaldehyde >= 0.1 {
		alertLevel = "critical"
		message = fmt.Sprintf("甲醛浓度严重超标: %.3f mg/m³", formaldehyde)
		threshold = 0.1
	} else if formaldehyde >= 0.08 {
		alertLevel = "warning"
		message = fmt.Sprintf("甲醛浓度超标: %.3f mg/m³", formaldehyde)
		threshold = 0.08
	} else {
		return nil
	}

	return &models.Alert{
		RuleID:         0, // 临时规则ID
		DeviceID:       data.DeviceID,
		Metric:         "formaldehyde",
		CurrentValue:   formaldehyde,
		ThresholdValue: threshold,
		Severity:       alertLevel,
		Status:         "active",
		TriggeredAt:    time.Now(),
		Message:        &message,
		CreatedAt:      time.Now(),
	}
}

func (s *unifiedSensorDataService) checkPM25Alert(ctx context.Context, data *models.UnifiedSensorData) *models.Alert {
	pm25 := *data.PM25
	var alertLevel string
	var message string
	var threshold float64

	if pm25 >= 75 {
		alertLevel = "critical"
		message = fmt.Sprintf("PM2.5浓度严重超标: %.1f μg/m³", pm25)
		threshold = 75
	} else if pm25 >= 35 {
		alertLevel = "warning"
		message = fmt.Sprintf("PM2.5浓度超标: %.1f μg/m³", pm25)
		threshold = 35
	} else {
		return nil
	}

	return &models.Alert{
		RuleID:         0, // 临时规则ID
		DeviceID:       data.DeviceID,
		Metric:         "pm25",
		CurrentValue:   pm25,
		ThresholdValue: threshold,
		Severity:       alertLevel,
		Status:         "active",
		TriggeredAt:    time.Now(),
		Message:        &message,
		CreatedAt:      time.Now(),
	}
}
