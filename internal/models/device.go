package models

import (
	"time"

	"gorm.io/gorm"
)

// Device 通用设备模型
type Device struct {
	ID                string         `json:"id" gorm:"primaryKey;type:varchar(64)"`
	Name              string         `json:"name" gorm:"type:varchar(100);not null"`
	Type              DeviceType     `json:"type" gorm:"type:varchar(50);not null"`
	LocationLatitude  *float64       `json:"location_latitude" gorm:"type:decimal(10,8)"`
	LocationLongitude *float64       `json:"location_longitude" gorm:"type:decimal(11,8)"`
	LocationAddress   *string        `json:"location_address" gorm:"type:varchar(200)"`
	Status            DeviceStatus   `json:"status" gorm:"type:varchar(20);default:'offline'"`
	Config            *string        `json:"config" gorm:"type:json"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}

// DeviceType 设备类型
type DeviceType string

const (
	DeviceTypeFormaldehyde DeviceType = "hcho"        // 甲醛传感器
	DeviceTypePM25         DeviceType = "pm25"        // PM2.5传感器
	DeviceTypePM10         DeviceType = "pm10"        // PM10传感器
	DeviceTypeCO2          DeviceType = "co2"         // CO2传感器
	DeviceTypeAirQuality   DeviceType = "air_quality" // 综合空气质量传感器
)

// IsValid 验证设备类型
func (t DeviceType) IsValid() bool {
	switch t {
	case DeviceTypeFormaldehyde, DeviceTypePM25, DeviceTypePM10, DeviceTypeCO2, DeviceTypeAirQuality:
		return true
	default:
		return false
	}
}

// GetSupportedMetrics 获取设备类型支持的指标
func (t DeviceType) GetSupportedMetrics() []string {
	switch t {
	case DeviceTypeFormaldehyde:
		return []string{"formaldehyde", "temperature", "humidity"}
	case DeviceTypePM25:
		return []string{"pm25", "temperature", "humidity"}
	case DeviceTypePM10:
		return []string{"pm10", "temperature", "humidity"}
	case DeviceTypeCO2:
		return []string{"co2", "temperature", "humidity"}
	case DeviceTypeAirQuality:
		return []string{"pm25", "pm10", "co2", "temperature", "humidity", "pressure"}
	default:
		return []string{}
	}
}

// DeviceConfig 通用设备配置
type DeviceConfig struct {
	ReportInterval int                    `json:"report_interval"`
	Sensors        map[string]bool        `json:"sensors"`
	Thresholds     map[string]interface{} `json:"thresholds"`
}

// Location 位置信息
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

// DeviceStatus 设备状态
type DeviceStatus string

const (
	DeviceStatusOnline      DeviceStatus = "online"
	DeviceStatusOffline     DeviceStatus = "offline"
	DeviceStatusMaintenance DeviceStatus = "maintenance"
	DeviceStatusError       DeviceStatus = "error"
)

// IsValid 验证设备状态
func (s DeviceStatus) IsValid() bool {
	switch s {
	case DeviceStatusOnline, DeviceStatusOffline, DeviceStatusMaintenance, DeviceStatusError:
		return true
	default:
		return false
	}
}

// DeviceCreateRequest 创建设备请求
type DeviceCreateRequest struct {
	ID       string       `json:"id" binding:"required"`
	Name     string       `json:"name" binding:"required"`
	Type     string       `json:"type" binding:"required"`
	Location Location     `json:"location"`
	Config   DeviceConfig `json:"config"`
}

// DeviceUpdateRequest 更新设备请求
type DeviceUpdateRequest struct {
	Name     *string       `json:"name,omitempty"`
	Type     *string       `json:"type,omitempty"`
	Location *Location     `json:"location,omitempty"`
	Config   *DeviceConfig `json:"config,omitempty"`
}

// DeviceListRequest 设备列表请求
type DeviceListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Status   string `form:"status"`
	Type     string `form:"type"`
	Keyword  string `form:"keyword"`
}

// DeviceListResponse 设备列表响应
type DeviceListResponse struct {
	Devices  []Device `json:"devices"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

// DeviceRuntimeStatus 设备运行时状态
type DeviceRuntimeStatus struct {
	ID              uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID        string     `json:"device_id" gorm:"type:varchar(64);not null;uniqueIndex"`
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
func (DeviceRuntimeStatus) TableName() string {
	return "device_runtime_status"
}

// DeviceRealtimeStatus 设备实时状态（用于API响应）
type DeviceRealtimeStatus struct {
	Device
	LastDataTime   *time.Time `json:"last_data_time"`
	PM25           *float64   `json:"pm25"`
	PM10           *float64   `json:"pm10"`
	CO2            *float64   `json:"co2"`
	Formaldehyde   *float64   `json:"formaldehyde"`
	Temperature    *float64   `json:"temperature"`
	Humidity       *float64   `json:"humidity"`
	Pressure       *float64   `json:"pressure"`
	QualityScore   *float64   `json:"quality_score"`
	RealtimeStatus string     `json:"realtime_status"`
}

// DeviceStatistics 设备统计信息
type DeviceStatistics struct {
	DeviceID        string     `json:"device_id"`
	TotalDataCount  int64      `json:"total_data_count"`
	OnlineTime      int64      `json:"online_time"`  // 秒
	OfflineTime     int64      `json:"offline_time"` // 秒
	LastOnlineTime  *time.Time `json:"last_online_time"`
	LastOfflineTime *time.Time `json:"last_offline_time"`
	AveragePM25     *float64   `json:"average_pm25"`
	AveragePM10     *float64   `json:"average_pm10"`
	AverageCO2      *float64   `json:"average_co2"`
	AverageTemp     *float64   `json:"average_temp"`
	AverageHumidity *float64   `json:"average_humidity"`
	AveragePressure *float64   `json:"average_pressure"`
}
