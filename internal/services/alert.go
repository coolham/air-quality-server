package services

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/utils"
	"context"
	"time"
)

// AlertService 告警服务接口
type AlertService interface {
	CreateAlert(ctx context.Context, alert *models.Alert) error
	GetAlert(ctx context.Context, id uint) (*models.Alert, error)
	UpdateAlert(ctx context.Context, alert *models.Alert) error
	DeleteAlert(ctx context.Context, id uint) error
	ListAlerts(ctx context.Context, limit, offset int) ([]models.Alert, error)
	CountAlerts(ctx context.Context) (int64, error)
	GetAlertsByDeviceID(ctx context.Context, deviceID string) ([]models.Alert, error)
	GetAlertsByStatus(ctx context.Context, status string) ([]models.Alert, error)
	GetUnresolvedAlerts(ctx context.Context) ([]models.Alert, error)
	ResolveAlert(ctx context.Context, alertID uint) error
	GetAlertsByTimeRange(ctx context.Context, startTime, endTime int64) ([]models.Alert, error)
	CheckAirQualityAlerts(ctx context.Context, data *models.AirQualityData) error
}

// alertService 告警服务实现
type alertService struct {
	alertRepo repositories.AlertRepository
	logger    utils.Logger
}

// NewAlertService 创建告警服务
func NewAlertService(alertRepo repositories.AlertRepository, logger utils.Logger) AlertService {
	return &alertService{
		alertRepo: alertRepo,
		logger:    logger,
	}
}

// CreateAlert 创建告警
func (s *alertService) CreateAlert(ctx context.Context, alert *models.Alert) error {
	now := time.Now()
	alert.CreatedAt = now
	alert.UpdatedAt = now

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		s.logger.Error("创建告警失败", utils.ErrorField(err))
		return err
	}

	s.logger.Info("告警创建成功", utils.Int("alert_id", int(alert.ID)), utils.String("metric", alert.Metric))
	return nil
}

// GetAlert 获取告警
func (s *alertService) GetAlert(ctx context.Context, id uint) (*models.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取告警失败", utils.ErrorField(err), utils.Int("alert_id", int(id)))
		return nil, err
	}
	return alert, nil
}

// UpdateAlert 更新告警
func (s *alertService) UpdateAlert(ctx context.Context, alert *models.Alert) error {
	now := time.Now()
	alert.UpdatedAt = now

	updateData := &models.Alert{
		Status:    alert.Status,
		Message:   alert.Message,
		UpdatedAt: alert.UpdatedAt,
	}
	if err := s.alertRepo.Update(ctx, alert.ID, updateData); err != nil {
		s.logger.Error("更新告警失败", utils.ErrorField(err), utils.Int("alert_id", int(alert.ID)))
		return err
	}

	s.logger.Info("告警更新成功", utils.Int("alert_id", int(alert.ID)))
	return nil
}

// DeleteAlert 删除告警
func (s *alertService) DeleteAlert(ctx context.Context, id uint) error {
	if err := s.alertRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除告警失败", utils.ErrorField(err), utils.Int("alert_id", int(id)))
		return err
	}

	s.logger.Info("告警删除成功", utils.Int("alert_id", int(id)))
	return nil
}

// ListAlerts 列出告警
func (s *alertService) ListAlerts(ctx context.Context, limit, offset int) ([]models.Alert, error) {
	req := &repositories.ListRequest{
		Page:     offset/limit + 1,
		PageSize: limit,
		OrderBy:  "created_at",
		Order:    "desc",
	}

	response, err := s.alertRepo.List(ctx, req)
	if err != nil {
		s.logger.Error("列出告警失败", utils.ErrorField(err))
		return nil, err
	}
	return response.Data, nil
}

// GetAlertsByDeviceID 根据设备ID获取告警
func (s *alertService) GetAlertsByDeviceID(ctx context.Context, deviceID string) ([]models.Alert, error) {
	alerts, err := s.alertRepo.GetByDeviceID(deviceID)
	if err != nil {
		s.logger.Error("根据设备ID获取告警失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return alerts, nil
}

// GetAlertsByStatus 根据状态获取告警
func (s *alertService) GetAlertsByStatus(ctx context.Context, status string) ([]models.Alert, error) {
	alerts, err := s.alertRepo.GetByStatus(status)
	if err != nil {
		s.logger.Error("根据状态获取告警失败", utils.ErrorField(err), utils.String("status", status))
		return nil, err
	}
	return alerts, nil
}

// GetUnresolvedAlerts 获取未解决的告警
func (s *alertService) GetUnresolvedAlerts(ctx context.Context) ([]models.Alert, error) {
	alerts, err := s.alertRepo.GetUnresolved()
	if err != nil {
		s.logger.Error("获取未解决告警失败", utils.ErrorField(err))
		return nil, err
	}
	return alerts, nil
}

// ResolveAlert 解决告警
func (s *alertService) ResolveAlert(ctx context.Context, alertID uint) error {
	if err := s.alertRepo.MarkAsResolved(alertID); err != nil {
		s.logger.Error("解决告警失败", utils.ErrorField(err), utils.Int("alert_id", int(alertID)))
		return err
	}

	s.logger.Info("告警已解决", utils.Int("alert_id", int(alertID)))
	return nil
}

// GetAlertsByTimeRange 根据时间范围获取告警
func (s *alertService) GetAlertsByTimeRange(ctx context.Context, startTime, endTime int64) ([]models.Alert, error) {
	alerts, err := s.alertRepo.GetByTimeRange(startTime, endTime)
	if err != nil {
		s.logger.Error("根据时间范围获取告警失败", utils.ErrorField(err))
		return nil, err
	}
	return alerts, nil
}

// CheckAirQualityAlerts 检查空气质量告警
func (s *alertService) CheckAirQualityAlerts(ctx context.Context, data *models.AirQualityData) error {
	var alerts []models.Alert

	// PM2.5告警检查
	if data.PM25 != nil && *data.PM25 > 75 {
		message := s.getPM2_5AlertMessage(*data.PM25)
		alert := &models.Alert{
			DeviceID:  data.DeviceID,
			Metric:    "pm25",
			Severity:  s.getAlertLevel(*data.PM25, 75, 150, 250),
			Message:   &message,
			Status:    "active",
			CreatedAt: data.CreatedAt,
		}
		alerts = append(alerts, *alert)
	}

	// PM10告警检查
	if data.PM10 != nil && *data.PM10 > 150 {
		message := s.getPM10AlertMessage(*data.PM10)
		alert := &models.Alert{
			DeviceID:  data.DeviceID,
			Metric:    "pm10",
			Severity:  s.getAlertLevel(*data.PM10, 150, 300, 500),
			Message:   &message,
			Status:    "active",
			CreatedAt: data.CreatedAt,
		}
		alerts = append(alerts, *alert)
	}

	// CO2告警检查
	if data.CO2 != nil && *data.CO2 > 1000 {
		message := s.getCO2AlertMessage(*data.CO2)
		alert := &models.Alert{
			DeviceID:  data.DeviceID,
			Metric:    "co2",
			Severity:  s.getAlertLevel(*data.CO2, 1000, 2000, 5000),
			Message:   &message,
			Status:    "active",
			CreatedAt: data.CreatedAt,
		}
		alerts = append(alerts, *alert)
	}

	// 温度告警检查
	if data.Temperature != nil && (*data.Temperature > 35 || *data.Temperature < -10) {
		message := s.getTemperatureAlertMessage(*data.Temperature)
		alert := &models.Alert{
			DeviceID:  data.DeviceID,
			Metric:    "temperature",
			Severity:  "warning",
			Message:   &message,
			Status:    "active",
			CreatedAt: data.CreatedAt,
		}
		alerts = append(alerts, *alert)
	}

	// 湿度告警检查
	if data.Humidity != nil && (*data.Humidity > 80 || *data.Humidity < 20) {
		message := s.getHumidityAlertMessage(*data.Humidity)
		alert := &models.Alert{
			DeviceID:  data.DeviceID,
			Metric:    "humidity",
			Severity:  "warning",
			Message:   &message,
			Status:    "active",
			CreatedAt: data.CreatedAt,
		}
		alerts = append(alerts, *alert)
	}

	// 创建告警
	for _, alert := range alerts {
		if err := s.CreateAlert(ctx, &alert); err != nil {
			s.logger.Error("创建空气质量告警失败", utils.ErrorField(err))
		}
	}

	return nil
}

// getAlertLevel 获取告警级别
func (s *alertService) getAlertLevel(value, warning, critical, severe float64) string {
	if value >= severe {
		return "severe"
	} else if value >= critical {
		return "critical"
	} else if value >= warning {
		return "warning"
	}
	return "info"
}

// getPM2_5AlertMessage 获取PM2.5告警消息
func (s *alertService) getPM2_5AlertMessage(value float64) string {
	message := "PM2.5浓度过高"
	return message
}

// getPM10AlertMessage 获取PM10告警消息
func (s *alertService) getPM10AlertMessage(value float64) string {
	message := "PM10浓度过高"
	return message
}

// getCO2AlertMessage 获取CO2告警消息
func (s *alertService) getCO2AlertMessage(value float64) string {
	message := "CO2浓度过高"
	return message
}

// getTemperatureAlertMessage 获取温度告警消息
func (s *alertService) getTemperatureAlertMessage(value float64) string {
	message := "温度异常"
	return message
}

// getHumidityAlertMessage 获取湿度告警消息
func (s *alertService) getHumidityAlertMessage(value float64) string {
	message := "湿度异常"
	return message
}

// CountAlerts 获取告警总数
func (s *alertService) CountAlerts(ctx context.Context) (int64, error) {
	count, err := s.alertRepo.Count(ctx, map[string]interface{}{})
	if err != nil {
		s.logger.Error("获取告警总数失败", utils.ErrorField(err))
		return 0, err
	}
	return count, nil
}
