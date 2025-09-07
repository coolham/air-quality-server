package models

import (
	"time"

	"gorm.io/gorm"
)

// FormaldehydeData 甲醛传感器数据模型
type FormaldehydeData struct {
	ID             uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID       string         `json:"device_id" gorm:"type:varchar(64);not null;index:idx_device_timestamp"`
	Timestamp      time.Time      `json:"timestamp" gorm:"not null;index:idx_device_timestamp;index:idx_timestamp"`
	Formaldehyde   *float64       `json:"formaldehyde" gorm:"type:decimal(8,3);comment:甲醛浓度 mg/m³"`
	Temperature    *float64       `json:"temperature" gorm:"type:decimal(6,2);comment:温度 °C"`
	Humidity       *float64       `json:"humidity" gorm:"type:decimal(6,2);comment:湿度 %"`
	Battery        *int           `json:"battery" gorm:"type:int;comment:电池电量 %"`
	SignalStrength *int           `json:"signal_strength" gorm:"type:int;comment:信号强度 dBm"`
	DataQuality    string         `json:"data_quality" gorm:"type:varchar(20);default:good;comment:数据质量"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (FormaldehydeData) TableName() string {
	return "formaldehyde_data"
}

// FormaldehydeDeviceStatus 甲醛设备状态模型
type FormaldehydeDeviceStatus struct {
	ID              uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID        string     `json:"device_id" gorm:"type:varchar(64);not null;uniqueIndex"`
	DeviceType      string     `json:"device_type" gorm:"type:varchar(50);not null;default:hcho"`
	Online          bool       `json:"online" gorm:"default:false"`
	BatteryLevel    *int       `json:"battery_level" gorm:"type:int;comment:电池电量 %"`
	SignalStrength  *int       `json:"signal_strength" gorm:"type:int;comment:信号强度 dBm"`
	LastDataTime    *time.Time `json:"last_data_time" gorm:"comment:最后数据时间"`
	ErrorCode       int        `json:"error_code" gorm:"default:0;comment:错误代码"`
	ErrorMessage    string     `json:"error_message" gorm:"type:text;comment:错误信息"`
	FirmwareVersion string     `json:"firmware_version" gorm:"type:varchar(50);comment:固件版本"`
	LastHeartbeat   time.Time  `json:"last_heartbeat" gorm:"autoUpdateTime;comment:最后心跳时间"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (FormaldehydeDeviceStatus) TableName() string {
	return "formaldehyde_device_status"
}

// FormaldehydeDeviceConfig 甲醛设备配置模型
type FormaldehydeDeviceConfig struct {
	ID                   uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID             string    `json:"device_id" gorm:"type:varchar(64);not null;uniqueIndex"`
	ReportInterval       int       `json:"report_interval" gorm:"default:300;comment:上报间隔(秒)"`
	FormaldehydeWarning  *float64  `json:"formaldehyde_warning" gorm:"type:decimal(8,3);default:0.08;comment:甲醛警告阈值 mg/m³"`
	FormaldehydeCritical *float64  `json:"formaldehyde_critical" gorm:"type:decimal(8,3);default:0.1;comment:甲醛严重阈值 mg/m³"`
	CalibrationEnabled   bool      `json:"calibration_enabled" gorm:"default:true;comment:校准功能启用"`
	CalibrationInterval  int       `json:"calibration_interval" gorm:"default:86400;comment:校准间隔(秒)"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (FormaldehydeDeviceConfig) TableName() string {
	return "formaldehyde_device_config"
}

// MQTTMessage MQTT消息结构
type MQTTMessage struct {
	DeviceID   string                 `json:"device_id"`
	DeviceType string                 `json:"device_type"`
	SensorID   string                 `json:"sensor_id"`
	SensorType string                 `json:"sensor_type"`
	Timestamp  int64                  `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	Location   *LocationInfo          `json:"location,omitempty"`
	Quality    *QualityInfo           `json:"quality,omitempty"`
}

// LocationInfo 位置信息
type LocationInfo struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Address   string   `json:"address,omitempty"`
}

// QualityInfo 数据质量信息
type QualityInfo struct {
	SignalStrength *int   `json:"signal_strength,omitempty"`
	DataQuality    string `json:"data_quality,omitempty"`
}

// DeviceStatusMessage 设备状态消息
type DeviceStatusMessage struct {
	DeviceID   string           `json:"device_id"`
	DeviceType string           `json:"device_type"`
	Timestamp  int64            `json:"timestamp"`
	Status     DeviceStatusData `json:"status"`
	Firmware   *FirmwareInfo    `json:"firmware,omitempty"`
}

// DeviceStatusData 设备状态数据
type DeviceStatusData struct {
	Online         bool   `json:"online"`
	BatteryLevel   *int   `json:"battery_level,omitempty"`
	SignalStrength *int   `json:"signal_strength,omitempty"`
	LastDataTime   *int64 `json:"last_data_time,omitempty"`
	ErrorCode      int    `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
}

// FirmwareInfo 固件信息
type FirmwareInfo struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
}

// ConfigMessage 配置下发消息
type ConfigMessage struct {
	DeviceID  string     `json:"device_id"`
	Timestamp int64      `json:"timestamp"`
	Config    ConfigData `json:"config"`
}

// ConfigData 配置数据
type ConfigData struct {
	ReportInterval int               `json:"report_interval"`
	Thresholds     ThresholdConfig   `json:"thresholds"`
	Calibration    CalibrationConfig `json:"calibration"`
}

// ThresholdConfig 阈值配置
type ThresholdConfig struct {
	FormaldehydeWarning  *float64 `json:"formaldehyde_warning"`
	FormaldehydeCritical *float64 `json:"formaldehyde_critical"`
}

// CalibrationConfig 校准配置
type CalibrationConfig struct {
	Enabled  bool `json:"enabled"`
	Interval int  `json:"interval"`
}

// CommandMessage 控制命令消息
type CommandMessage struct {
	DeviceID  string      `json:"device_id"`
	Timestamp int64       `json:"timestamp"`
	Command   CommandData `json:"command"`
}

// CommandData 命令数据
type CommandData struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// FormaldehydeStatistics 甲醛数据统计
type FormaldehydeStatistics struct {
	DataCount       int64   `json:"data_count"`
	FormaldehydeAvg float64 `json:"formaldehyde_avg"`
	FormaldehydeMin float64 `json:"formaldehyde_min"`
	FormaldehydeMax float64 `json:"formaldehyde_max"`
	TemperatureAvg  float64 `json:"temperature_avg"`
	TemperatureMin  float64 `json:"temperature_min"`
	TemperatureMax  float64 `json:"temperature_max"`
	HumidityAvg     float64 `json:"humidity_avg"`
	HumidityMin     float64 `json:"humidity_min"`
	HumidityMax     float64 `json:"humidity_max"`
}
