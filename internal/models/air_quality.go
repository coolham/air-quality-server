package models

import (
	"time"

	"gorm.io/gorm"
)

// AirQualityData 空气质量数据模型
type AirQualityData struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID     string         `json:"device_id" gorm:"type:varchar(64);not null;index:idx_device_timestamp"`
	Timestamp    time.Time      `json:"timestamp" gorm:"not null;index:idx_device_timestamp;index:idx_timestamp"`
	PM25         *float64       `json:"pm25" gorm:"type:decimal(8,2)"`
	PM10         *float64       `json:"pm10" gorm:"type:decimal(8,2)"`
	CO2          *float64       `json:"co2" gorm:"type:decimal(8,2)"`
	Temperature  *float64       `json:"temperature" gorm:"type:decimal(6,2)"`
	Humidity     *float64       `json:"humidity" gorm:"type:decimal(6,2)"`
	Pressure     *float64       `json:"pressure" gorm:"type:decimal(8,2)"`
	QualityScore *float64       `json:"quality_score" gorm:"type:decimal(4,2)"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (AirQualityData) TableName() string {
	return "air_quality_data"
}

// AirQualityDataUpload 数据上传请求
type AirQualityDataUpload struct {
	DeviceID  string                 `json:"device_id" binding:"required"`
	Timestamp int64                  `json:"timestamp" binding:"required"`
	Data      AirQualityDataPayload  `json:"data" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AirQualityDataPayload 数据载荷
type AirQualityDataPayload struct {
	PM25        *float64  `json:"pm25,omitempty"`
	PM10        *float64  `json:"pm10,omitempty"`
	CO2         *float64  `json:"co2,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
	Humidity    *float64  `json:"humidity,omitempty"`
	Pressure    *float64  `json:"pressure,omitempty"`
	Location    *Location `json:"location,omitempty"`
}

// AirQualityDataBatch 批量数据上传请求
type AirQualityDataBatch struct {
	DeviceID string                 `json:"device_id" binding:"required"`
	Data     []AirQualityDataUpload `json:"data" binding:"required"`
}

// AirQualityDataQuery 数据查询请求
type AirQualityDataQuery struct {
	DeviceIDs   []string `form:"device_ids"`
	StartTime   int64    `form:"start_time"`
	EndTime     int64    `form:"end_time"`
	Interval    string   `form:"interval"`    // 1m, 5m, 1h, 1d
	Metrics     []string `form:"metrics"`     // pm25, pm10, co2, etc.
	Aggregation string   `form:"aggregation"` // avg, max, min, sum
	Limit       int      `form:"limit" binding:"min=1,max=1000"`
	Offset      int      `form:"offset" binding:"min=0"`
}

// AirQualityDataResponse 数据查询响应
type AirQualityDataResponse struct {
	Data     []AirQualityData `json:"data"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// AirQualityStatistics 空气质量统计
type AirQualityStatistics struct {
	DeviceID     string   `json:"device_id"`
	TimeRange    string   `json:"time_range"`
	PM25Avg      *float64 `json:"pm25_avg"`
	PM25Max      *float64 `json:"pm25_max"`
	PM25Min      *float64 `json:"pm25_min"`
	PM10Avg      *float64 `json:"pm10_avg"`
	PM10Max      *float64 `json:"pm10_max"`
	PM10Min      *float64 `json:"pm10_min"`
	CO2Avg       *float64 `json:"co2_avg"`
	CO2Max       *float64 `json:"co2_max"`
	CO2Min       *float64 `json:"co2_min"`
	TempAvg      *float64 `json:"temp_avg"`
	TempMax      *float64 `json:"temp_max"`
	TempMin      *float64 `json:"temp_min"`
	HumidityAvg  *float64 `json:"humidity_avg"`
	HumidityMax  *float64 `json:"humidity_max"`
	HumidityMin  *float64 `json:"humidity_min"`
	PressureAvg  *float64 `json:"pressure_avg"`
	PressureMax  *float64 `json:"pressure_max"`
	PressureMin  *float64 `json:"pressure_min"`
	DataCount    int64    `json:"data_count"`
	QualityScore *float64 `json:"quality_score"`
}

// AirQualityTrend 空气质量趋势
type AirQualityTrend struct {
	Time         time.Time `json:"time"`
	PM25         *float64  `json:"pm25"`
	PM10         *float64  `json:"pm10"`
	CO2          *float64  `json:"co2"`
	Temperature  *float64  `json:"temperature"`
	Humidity     *float64  `json:"humidity"`
	Pressure     *float64  `json:"pressure"`
	QualityScore *float64  `json:"quality_score"`
}

// AirQualityAlert 空气质量告警
type AirQualityAlert struct {
	DeviceID  string    `json:"device_id"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// AirQualityLevel 空气质量等级
type AirQualityLevel string

const (
	AirQualityLevelExcellent AirQualityLevel = "excellent" // 优
	AirQualityLevelGood      AirQualityLevel = "good"      // 良
	AirQualityLevelLight     AirQualityLevel = "light"     // 轻度污染
	AirQualityLevelModerate  AirQualityLevel = "moderate"  // 中度污染
	AirQualityLevelHeavy     AirQualityLevel = "heavy"     // 重度污染
	AirQualityLevelSevere    AirQualityLevel = "severe"    // 严重污染
)

// GetAirQualityLevel 根据PM2.5值获取空气质量等级
func GetAirQualityLevel(pm25 float64) AirQualityLevel {
	switch {
	case pm25 <= 35:
		return AirQualityLevelExcellent
	case pm25 <= 75:
		return AirQualityLevelGood
	case pm25 <= 115:
		return AirQualityLevelLight
	case pm25 <= 150:
		return AirQualityLevelModerate
	case pm25 <= 250:
		return AirQualityLevelHeavy
	default:
		return AirQualityLevelSevere
	}
}

// GetAirQualityLevelDescription 获取空气质量等级描述
func GetAirQualityLevelDescription(level AirQualityLevel) string {
	switch level {
	case AirQualityLevelExcellent:
		return "优"
	case AirQualityLevelGood:
		return "良"
	case AirQualityLevelLight:
		return "轻度污染"
	case AirQualityLevelModerate:
		return "中度污染"
	case AirQualityLevelHeavy:
		return "重度污染"
	case AirQualityLevelSevere:
		return "严重污染"
	default:
		return "未知"
	}
}

// CalculateQualityScore 计算数据质量评分
func CalculateQualityScore(data *AirQualityDataPayload) float64 {
	score := 0.0
	count := 0

	if data.PM25 != nil {
		score += 1.0
		count++
	}
	if data.PM10 != nil {
		score += 1.0
		count++
	}
	if data.CO2 != nil {
		score += 1.0
		count++
	}
	if data.Temperature != nil {
		score += 1.0
		count++
	}
	if data.Humidity != nil {
		score += 1.0
		count++
	}
	if data.Pressure != nil {
		score += 1.0
		count++
	}

	if count == 0 {
		return 0.0
	}

	return score / float64(count)
}
