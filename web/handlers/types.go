package handlers

import "time"

// DeviceStats 设备统计信息
type DeviceStats struct {
	TotalDevices   int `json:"total_devices"`
	OnlineDevices  int `json:"online_devices"`
	OfflineDevices int `json:"offline_devices"`
	ActiveDevices  int `json:"active_devices"`
}

// AlertStats 告警统计信息
type AlertStats struct {
	TotalAlerts      int `json:"total_alerts"`
	UnresolvedAlerts int `json:"unresolved_alerts"`
	CriticalAlerts   int `json:"critical_alerts"`
	WarningAlerts    int `json:"warning_alerts"`
}

// AirQualityDataSummary 空气质量数据摘要
type AirQualityDataSummary struct {
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	DeviceType   string    `json:"device_type"`
	SensorID     string    `json:"sensor_id"`
	SensorType   string    `json:"sensor_type"`
	PM25         float64   `json:"pm25"`
	PM10         float64   `json:"pm10"`
	CO2          float64   `json:"co2"`
	Formaldehyde float64   `json:"formaldehyde"`
	Temp         float64   `json:"temp"`
	Humidity     float64   `json:"humidity"`
	Pressure     float64   `json:"pressure"`
	Battery      int       `json:"battery"`
	DataQuality  string    `json:"data_quality"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

// Pagination 分页信息
type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
	PageSize    int `json:"page_size"`
}

// ChartData 图表数据
type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

// Dataset 数据集
type Dataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BorderColor     string    `json:"borderColor"`
	BackgroundColor string    `json:"backgroundColor"`
	Fill            bool      `json:"fill"`
}

// WebConfig Web配置
type WebConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}
