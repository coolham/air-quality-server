package models

import (
	"time"

	"gorm.io/gorm"
)

// SensorData 通用传感器数据模型
type SensorData struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID   string     `json:"device_id" gorm:"type:varchar(64);not null;index:idx_device_timestamp"`
	DeviceType DeviceType `json:"device_type" gorm:"type:varchar(50);not null;index:idx_device_type"`
	Timestamp  time.Time  `json:"timestamp" gorm:"not null;index:idx_device_timestamp;index:idx_timestamp"`

	// 通用传感器数据
	PM25         *float64 `json:"pm25" gorm:"type:decimal(8,3);comment:PM2.5浓度 μg/m³"`
	PM10         *float64 `json:"pm10" gorm:"type:decimal(8,3);comment:PM10浓度 μg/m³"`
	CO2          *float64 `json:"co2" gorm:"type:decimal(8,3);comment:CO2浓度 ppm"`
	Formaldehyde *float64 `json:"formaldehyde" gorm:"type:decimal(8,3);comment:甲醛浓度 mg/m³"`
	Temperature  *float64 `json:"temperature" gorm:"type:decimal(6,2);comment:温度 °C"`
	Humidity     *float64 `json:"humidity" gorm:"type:decimal(6,2);comment:湿度 %"`
	Pressure     *float64 `json:"pressure" gorm:"type:decimal(8,2);comment:气压 hPa"`

	// 设备状态信息
	Battery        *int   `json:"battery" gorm:"type:int;comment:电池电量 %"`
	SignalStrength *int   `json:"signal_strength" gorm:"type:int;comment:信号强度 dBm"`
	DataQuality    string `json:"data_quality" gorm:"type:varchar(20);default:good;comment:数据质量"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (SensorData) TableName() string {
	return "sensor_data"
}

// GetMetricValue 获取指定指标的值
func (s *SensorData) GetMetricValue(metric string) *float64 {
	switch metric {
	case "pm25":
		return s.PM25
	case "pm10":
		return s.PM10
	case "co2":
		return s.CO2
	case "formaldehyde":
		return s.Formaldehyde
	case "temperature":
		return s.Temperature
	case "humidity":
		return s.Humidity
	case "pressure":
		return s.Pressure
	default:
		return nil
	}
}

// SetMetricValue 设置指定指标的值
func (s *SensorData) SetMetricValue(metric string, value *float64) {
	switch metric {
	case "pm25":
		s.PM25 = value
	case "pm10":
		s.PM10 = value
	case "co2":
		s.CO2 = value
	case "formaldehyde":
		s.Formaldehyde = value
	case "temperature":
		s.Temperature = value
	case "humidity":
		s.Humidity = value
	case "pressure":
		s.Pressure = value
	}
}

// GetAvailableMetrics 获取当前数据中可用的指标
func (s *SensorData) GetAvailableMetrics() []string {
	var metrics []string
	if s.PM25 != nil {
		metrics = append(metrics, "pm25")
	}
	if s.PM10 != nil {
		metrics = append(metrics, "pm10")
	}
	if s.CO2 != nil {
		metrics = append(metrics, "co2")
	}
	if s.Formaldehyde != nil {
		metrics = append(metrics, "formaldehyde")
	}
	if s.Temperature != nil {
		metrics = append(metrics, "temperature")
	}
	if s.Humidity != nil {
		metrics = append(metrics, "humidity")
	}
	if s.Pressure != nil {
		metrics = append(metrics, "pressure")
	}
	return metrics
}

// SensorDataStatistics 传感器数据统计
type SensorDataStatistics struct {
	DataCount       int64   `json:"data_count"`
	PM25Avg         float64 `json:"pm25_avg"`
	PM25Min         float64 `json:"pm25_min"`
	PM25Max         float64 `json:"pm25_max"`
	PM10Avg         float64 `json:"pm10_avg"`
	PM10Min         float64 `json:"pm10_min"`
	PM10Max         float64 `json:"pm10_max"`
	CO2Avg          float64 `json:"co2_avg"`
	CO2Min          float64 `json:"co2_min"`
	CO2Max          float64 `json:"co2_max"`
	FormaldehydeAvg float64 `json:"formaldehyde_avg"`
	FormaldehydeMin float64 `json:"formaldehyde_min"`
	FormaldehydeMax float64 `json:"formaldehyde_max"`
	TemperatureAvg  float64 `json:"temperature_avg"`
	TemperatureMin  float64 `json:"temperature_min"`
	TemperatureMax  float64 `json:"temperature_max"`
	HumidityAvg     float64 `json:"humidity_avg"`
	HumidityMin     float64 `json:"humidity_min"`
	HumidityMax     float64 `json:"humidity_max"`
	PressureAvg     float64 `json:"pressure_avg"`
	PressureMin     float64 `json:"pressure_min"`
	PressureMax     float64 `json:"pressure_max"`
}
