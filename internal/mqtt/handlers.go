package mqtt

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SensorDataHandler 传感器数据处理器
type SensorDataHandler struct {
	dataRepo   repositories.UnifiedSensorDataRepository
	deviceRepo repositories.DeviceRepository
	alertSvc   services.AlertService
	logger     utils.Logger
}

// NewSensorDataHandler 创建传感器数据处理器
func NewSensorDataHandler(
	dataRepo repositories.UnifiedSensorDataRepository,
	deviceRepo repositories.DeviceRepository,
	alertSvc services.AlertService,
	logger utils.Logger,
) *SensorDataHandler {
	return &SensorDataHandler{
		dataRepo:   dataRepo,
		deviceRepo: deviceRepo,
		alertSvc:   alertSvc,
		logger:     logger,
	}
}

// HandleMessage 处理传感器数据消息
func (h *SensorDataHandler) HandleMessage(topic string, payload []byte) error {
	var msg models.MQTTMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.logger.Error("解析甲醛数据消息失败", utils.ErrorField(err))
		return err
	}

	// 验证必要字段
	if msg.DeviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}

	// 转换时间戳
	timestamp := time.Unix(msg.Timestamp, 0)

	// 确定设备类型
	deviceType := models.DeviceTypeFormaldehyde // 默认甲醛传感器
	if msg.DeviceType != "" {
		deviceType = models.DeviceType(msg.DeviceType)
	}

	// 提取数据
	sensorData := &models.UnifiedSensorData{
		DeviceID:    msg.DeviceID,
		DeviceType:  deviceType,
		SensorID:    msg.SensorID,
		SensorType:  msg.SensorType,
		Timestamp:   timestamp,
		DataQuality: "good",
	}

	// 解析数据字段
	if data, ok := msg.Data["formaldehyde"].(float64); ok {
		sensorData.Formaldehyde = &data
	}
	if data, ok := msg.Data["pm25"].(float64); ok {
		sensorData.PM25 = &data
	}
	if data, ok := msg.Data["pm10"].(float64); ok {
		sensorData.PM10 = &data
	}
	if data, ok := msg.Data["co2"].(float64); ok {
		sensorData.CO2 = &data
	}
	if data, ok := msg.Data["temperature"].(float64); ok {
		sensorData.Temperature = &data
	}
	if data, ok := msg.Data["humidity"].(float64); ok {
		sensorData.Humidity = &data
	}
	if data, ok := msg.Data["pressure"].(float64); ok {
		sensorData.Pressure = &data
	}
	if data, ok := msg.Data["battery"].(float64); ok {
		battery := int(data)
		sensorData.Battery = &battery
	}

	// 解析质量信息
	if msg.Quality != nil {
		if msg.Quality.SignalStrength != nil {
			sensorData.SignalStrength = msg.Quality.SignalStrength
		}
		if msg.Quality.DataQuality != "" {
			sensorData.DataQuality = msg.Quality.DataQuality
		}
		// Battery 字段在 QualityInfo 中不存在，已在上面处理
	}

	// 解析位置信息
	if msg.Location != nil {
		if msg.Location.Latitude != nil {
			sensorData.Latitude = msg.Location.Latitude
		}
		if msg.Location.Longitude != nil {
			sensorData.Longitude = msg.Location.Longitude
		}
	}

	// 保存数据
	ctx := context.Background()
	if err := h.dataRepo.Create(ctx, sensorData); err != nil {
		h.logger.Error("保存传感器数据失败",
			utils.String("device_id", msg.DeviceID),
			utils.ErrorField(err))
		return err
	}

	// 检查告警
	if err := h.checkAlerts(ctx, sensorData); err != nil {
		h.logger.Error("检查告警失败",
			utils.String("device_id", msg.DeviceID),
			utils.ErrorField(err))
	}

	h.logger.Info("处理传感器数据成功",
		utils.String("device_id", msg.DeviceID),
		utils.String("device_type", string(deviceType)),
		utils.String("sensor_id", msg.SensorID),
		utils.String("sensor_type", msg.SensorType),
		utils.Float64("formaldehyde", getFloatValue(sensorData.Formaldehyde)))

	return nil
}

// checkAlerts 检查告警
func (h *SensorDataHandler) checkAlerts(ctx context.Context, data *models.UnifiedSensorData) error {
	if data.Formaldehyde == nil {
		return nil
	}

	formaldehyde := *data.Formaldehyde
	var alertLevel string
	var message string

	// 检查甲醛浓度告警
	if formaldehyde >= 0.1 {
		alertLevel = "critical"
		message = fmt.Sprintf("甲醛浓度严重超标: %.3f mg/m³", formaldehyde)
	} else if formaldehyde >= 0.08 {
		alertLevel = "warning"
		message = fmt.Sprintf("甲醛浓度超标: %.3f mg/m³", formaldehyde)
	} else {
		return nil // 正常范围，无需告警
	}

	// 创建告警 (使用默认规则ID 0，表示系统自动生成的告警)
	alert := &models.Alert{
		RuleID:         0, // 系统自动告警，无对应规则
		DeviceID:       data.DeviceID,
		Metric:         "formaldehyde",
		CurrentValue:   formaldehyde,
		ThresholdValue: 0.08, // 默认阈值
		Severity:       alertLevel,
		Status:         "active",
		TriggeredAt:    time.Now(),
		Message:        &message,
	}

	return h.alertSvc.CreateAlert(ctx, alert)
}

// getFloatValue 安全获取浮点数值
func getFloatValue(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}
