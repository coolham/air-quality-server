package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// UnifiedSensorData 统一传感器数据模型
type UnifiedSensorData struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID   string     `json:"device_id" gorm:"type:varchar(64);not null;index:idx_device_timestamp"`
	DeviceType DeviceType `json:"device_type" gorm:"type:varchar(50);not null;index:idx_device_type"`
	SensorID   string     `json:"sensor_id" gorm:"type:varchar(64);comment:传感器ID;index:idx_sensor_id"`
	SensorType string     `json:"sensor_type" gorm:"type:varchar(50);comment:传感器类型;index:idx_sensor_type"`
	Timestamp  time.Time  `json:"timestamp" gorm:"not null;index:idx_device_timestamp;index:idx_timestamp"`

	// 核心环境指标
	PM25         *float64 `json:"pm25" gorm:"type:decimal(8,3);comment:PM2.5浓度 μg/m³;index:idx_pm25"`
	PM10         *float64 `json:"pm10" gorm:"type:decimal(8,3);comment:PM10浓度 μg/m³;index:idx_pm10"`
	CO2          *float64 `json:"co2" gorm:"type:decimal(8,3);comment:CO2浓度 ppm;index:idx_co2"`
	Formaldehyde *float64 `json:"formaldehyde" gorm:"type:decimal(8,3);comment:甲醛浓度 mg/m³;index:idx_formaldehyde"`

	// 环境参数
	Temperature *float64 `json:"temperature" gorm:"type:decimal(6,2);comment:温度 °C;index:idx_temperature"`
	Humidity    *float64 `json:"humidity" gorm:"type:decimal(6,2);comment:湿度 %;index:idx_humidity"`
	Pressure    *float64 `json:"pressure" gorm:"type:decimal(8,2);comment:气压 hPa;index:idx_pressure"`

	// 其他污染物指标（可扩展）
	O3  *float64 `json:"o3" gorm:"type:decimal(8,3);comment:臭氧浓度 μg/m³"`
	NO2 *float64 `json:"no2" gorm:"type:decimal(8,3);comment:二氧化氮浓度 μg/m³"`
	SO2 *float64 `json:"so2" gorm:"type:decimal(8,3);comment:二氧化硫浓度 μg/m³"`
	CO  *float64 `json:"co" gorm:"type:decimal(8,3);comment:一氧化碳浓度 mg/m³"`
	VOC *float64 `json:"voc" gorm:"type:decimal(8,3);comment:挥发性有机化合物 μg/m³"`

	// 设备状态信息
	Battery        *int   `json:"battery" gorm:"type:int;comment:电池电量 %"`
	SignalStrength *int   `json:"signal_strength" gorm:"type:int;comment:信号强度 dBm"`
	DataQuality    string `json:"data_quality" gorm:"type:varchar(20);default:good;comment:数据质量"`

	// 位置信息
	Latitude  *float64 `json:"latitude" gorm:"type:decimal(10,8);comment:纬度"`
	Longitude *float64 `json:"longitude" gorm:"type:decimal(11,8);comment:经度"`

	// 扩展数据（JSON格式存储非标准指标）
	ExtendedData *string `json:"extended_data" gorm:"type:json;comment:扩展数据"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (UnifiedSensorData) TableName() string {
	return "unified_sensor_data"
}

// GetMetricValue 获取指定指标的值
func (s *UnifiedSensorData) GetMetricValue(metric string) *float64 {
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
	case "o3":
		return s.O3
	case "no2":
		return s.NO2
	case "so2":
		return s.SO2
	case "co":
		return s.CO
	case "voc":
		return s.VOC
	default:
		// 尝试从扩展数据中获取
		return s.getExtendedMetricValue(metric)
	}
}

// SetMetricValue 设置指定指标的值
func (s *UnifiedSensorData) SetMetricValue(metric string, value *float64) {
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
	case "o3":
		s.O3 = value
	case "no2":
		s.NO2 = value
	case "so2":
		s.SO2 = value
	case "co":
		s.CO = value
	case "voc":
		s.VOC = value
	default:
		// 存储到扩展数据中
		s.setExtendedMetricValue(metric, value)
	}
}

// getExtendedMetricValue 从扩展数据中获取指标值
func (s *UnifiedSensorData) getExtendedMetricValue(metric string) *float64 {
	if s.ExtendedData == nil {
		return nil
	}

	var extended map[string]interface{}
	if err := json.Unmarshal([]byte(*s.ExtendedData), &extended); err != nil {
		return nil
	}

	if value, ok := extended[metric].(float64); ok {
		return &value
	}
	return nil
}

// setExtendedMetricValue 设置扩展数据中的指标值
func (s *UnifiedSensorData) setExtendedMetricValue(metric string, value *float64) {
	var extended map[string]interface{}

	if s.ExtendedData != nil {
		json.Unmarshal([]byte(*s.ExtendedData), &extended)
	} else {
		extended = make(map[string]interface{})
	}

	if value != nil {
		extended[metric] = *value
	} else {
		delete(extended, metric)
	}

	if len(extended) > 0 {
		data, _ := json.Marshal(extended)
		extendedStr := string(data)
		s.ExtendedData = &extendedStr
	} else {
		s.ExtendedData = nil
	}
}

// GetAvailableMetrics 获取当前数据中可用的指标
func (s *UnifiedSensorData) GetAvailableMetrics() []string {
	var metrics []string

	// 检查标准指标
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
	if s.O3 != nil {
		metrics = append(metrics, "o3")
	}
	if s.NO2 != nil {
		metrics = append(metrics, "no2")
	}
	if s.SO2 != nil {
		metrics = append(metrics, "so2")
	}
	if s.CO != nil {
		metrics = append(metrics, "co")
	}
	if s.VOC != nil {
		metrics = append(metrics, "voc")
	}

	// 检查扩展指标
	if s.ExtendedData != nil {
		var extended map[string]interface{}
		if err := json.Unmarshal([]byte(*s.ExtendedData), &extended); err == nil {
			for metric := range extended {
				metrics = append(metrics, metric)
			}
		}
	}

	return metrics
}

// GetDataByDeviceType 根据设备类型获取相关数据
func (s *UnifiedSensorData) GetDataByDeviceType() map[string]*float64 {
	data := make(map[string]*float64)

	// 根据设备类型返回相关指标
	switch s.DeviceType {
	case DeviceTypeFormaldehyde:
		if s.Formaldehyde != nil {
			data["formaldehyde"] = s.Formaldehyde
		}
		if s.Temperature != nil {
			data["temperature"] = s.Temperature
		}
		if s.Humidity != nil {
			data["humidity"] = s.Humidity
		}

	case DeviceTypePM25:
		if s.PM25 != nil {
			data["pm25"] = s.PM25
		}
		if s.Temperature != nil {
			data["temperature"] = s.Temperature
		}
		if s.Humidity != nil {
			data["humidity"] = s.Humidity
		}

	case DeviceTypeAirQuality:
		// 综合空气质量传感器，返回所有可用数据
		if s.PM25 != nil {
			data["pm25"] = s.PM25
		}
		if s.PM10 != nil {
			data["pm10"] = s.PM10
		}
		if s.CO2 != nil {
			data["co2"] = s.CO2
		}
		if s.Temperature != nil {
			data["temperature"] = s.Temperature
		}
		if s.Humidity != nil {
			data["humidity"] = s.Humidity
		}
		if s.Pressure != nil {
			data["pressure"] = s.Pressure
		}
	}

	return data
}

// UnifiedSensorDataStatistics 统一传感器数据统计
type UnifiedSensorDataStatistics struct {
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

// UnifiedSensorDataUpload 统一传感器数据上传请求
type UnifiedSensorDataUpload struct {
	DeviceID   string                 `json:"device_id" binding:"required"`
	DeviceType string                 `json:"device_type" binding:"required"`
	SensorID   string                 `json:"sensor_id"`
	SensorType string                 `json:"sensor_type"`
	Timestamp  int64                  `json:"timestamp" binding:"required"`
	Data       map[string]interface{} `json:"data" binding:"required"`
	Location   *LocationInfo          `json:"location,omitempty"`
	Quality    *QualityInfo           `json:"quality,omitempty"`
	Extended   map[string]interface{} `json:"extended,omitempty"`
}
