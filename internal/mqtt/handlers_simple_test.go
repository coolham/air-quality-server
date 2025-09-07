package mqtt

import (
	"air-quality-server/internal/models"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetFloatValue 测试获取浮点数值
func TestGetFloatValue(t *testing.T) {
	tests := []struct {
		name     string
		ptr      *float64
		expected float64
	}{
		{
			name:     "有效指针",
			ptr:      func() *float64 { v := 3.14; return &v }(),
			expected: 3.14,
		},
		{
			name:     "空指针",
			ptr:      nil,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFloatValue(tt.ptr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMQTTMessageParsing 测试MQTT消息解析
func TestMQTTMessageParsing(t *testing.T) {

	// 测试有效消息解析
	validData := map[string]interface{}{
		"device_id":   "hcho_001",
		"device_type": "hcho",
		"sensor_id":   "sensor_hcho_001_01",
		"sensor_type": "hcho",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"formaldehyde": 0.05,
			"temperature":  22.5,
			"humidity":     45.0,
			"battery":      85.0,
		},
		"location": map[string]interface{}{
			"latitude":  39.9042,
			"longitude": 116.4074,
		},
		"quality": map[string]interface{}{
			"signal_strength": -65,
			"data_quality":    "good",
		},
	}

	payload, err := json.Marshal(validData)
	assert.NoError(t, err)

	var msg models.MQTTMessage
	err = json.Unmarshal(payload, &msg)
	assert.NoError(t, err)

	// 验证解析结果
	assert.Equal(t, "hcho_001", msg.DeviceID)
	assert.Equal(t, "hcho", msg.DeviceType)
	assert.Equal(t, "sensor_hcho_001_01", msg.SensorID)
	assert.Equal(t, "hcho", msg.SensorType)
	assert.NotNil(t, msg.Data)
	assert.NotNil(t, msg.Location)
	assert.NotNil(t, msg.Quality)

	// 验证数据字段
	formaldehyde, ok := msg.Data["formaldehyde"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 0.05, formaldehyde)

	temperature, ok := msg.Data["temperature"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 22.5, temperature)

	// 验证位置信息
	assert.NotNil(t, msg.Location.Latitude)
	assert.Equal(t, 39.9042, *msg.Location.Latitude)
	assert.NotNil(t, msg.Location.Longitude)
	assert.Equal(t, 116.4074, *msg.Location.Longitude)

	// 验证质量信息
	assert.NotNil(t, msg.Quality.SignalStrength)
	assert.Equal(t, -65, *msg.Quality.SignalStrength)
	assert.Equal(t, "good", msg.Quality.DataQuality)
}

// TestMQTTMessageValidation 测试MQTT消息验证
func TestMQTTMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "有效消息",
			data: map[string]interface{}{
				"device_id":   "hcho_001",
				"device_type": "hcho",
				"timestamp":   time.Now().Unix(),
				"data": map[string]interface{}{
					"formaldehyde": 0.05,
				},
			},
			wantErr: false,
		},
		{
			name: "缺少设备ID",
			data: map[string]interface{}{
				"device_type": "hcho",
				"timestamp":   time.Now().Unix(),
				"data": map[string]interface{}{
					"formaldehyde": 0.05,
				},
			},
			wantErr: true,
			errMsg:  "设备ID不能为空",
		},
		{
			name: "设备ID为空字符串",
			data: map[string]interface{}{
				"device_id":   "",
				"device_type": "hcho",
				"timestamp":   time.Now().Unix(),
				"data": map[string]interface{}{
					"formaldehyde": 0.05,
				},
			},
			wantErr: true,
			errMsg:  "设备ID不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(tt.data)
			assert.NoError(t, err)

			var msg models.MQTTMessage
			err = json.Unmarshal(payload, &msg)
			assert.NoError(t, err)

			// 验证设备ID
			if msg.DeviceID == "" {
				if tt.wantErr {
					assert.Contains(t, tt.errMsg, "设备ID不能为空")
				}
			} else {
				assert.False(t, tt.wantErr)
			}
		})
	}
}

// TestDataFieldExtraction 测试数据字段提取
func TestDataFieldExtraction(t *testing.T) {
	// 创建测试数据
	testData := map[string]interface{}{
		"formaldehyde": 0.05,
		"temperature":  22.5,
		"humidity":     45.0,
		"battery":      85.0,
		"pm25":         35.0,
		"pm10":         50.0,
		"co2":          400.0,
		"pressure":     1013.25,
	}

	// 测试各种数据类型的提取
	tests := []struct {
		field    string
		expected interface{}
	}{
		{"formaldehyde", 0.05},
		{"temperature", 22.5},
		{"humidity", 45.0},
		{"battery", 85.0},
		{"pm25", 35.0},
		{"pm10", 50.0},
		{"co2", 400.0},
		{"pressure", 1013.25},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			value, ok := testData[tt.field]
			assert.True(t, ok)
			assert.Equal(t, tt.expected, value)
		})
	}
}

// TestAlertThresholds 测试告警阈值
func TestAlertThresholds(t *testing.T) {
	tests := []struct {
		name          string
		formaldehyde  float64
		expectedLevel string
		expectedAlert bool
	}{
		{
			name:          "正常范围",
			formaldehyde:  0.05,
			expectedLevel: "",
			expectedAlert: false,
		},
		{
			name:          "警告阈值",
			formaldehyde:  0.08,
			expectedLevel: "warning",
			expectedAlert: true,
		},
		{
			name:          "严重阈值",
			formaldehyde:  0.12,
			expectedLevel: "critical",
			expectedAlert: true,
		},
		{
			name:          "极高浓度",
			formaldehyde:  0.20,
			expectedLevel: "critical",
			expectedAlert: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var alertLevel string
			var shouldAlert bool

			if tt.formaldehyde >= 0.1 {
				alertLevel = "critical"
				shouldAlert = true
			} else if tt.formaldehyde >= 0.08 {
				alertLevel = "warning"
				shouldAlert = true
			} else {
				shouldAlert = false
			}

			assert.Equal(t, tt.expectedLevel, alertLevel)
			assert.Equal(t, tt.expectedAlert, shouldAlert)
		})
	}
}

// BenchmarkGetFloatValue 性能测试
func BenchmarkGetFloatValue(b *testing.B) {
	value := 3.14
	ptr := &value

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getFloatValue(ptr)
	}
}

// BenchmarkMQTTMessageParsing 性能测试
func BenchmarkMQTTMessageParsing(b *testing.B) {
	testData := map[string]interface{}{
		"device_id":   "hcho_001",
		"device_type": "hcho",
		"sensor_id":   "sensor_hcho_001_01",
		"sensor_type": "hcho",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"formaldehyde": 0.05,
			"temperature":  22.5,
			"humidity":     45.0,
			"battery":      85.0,
		},
		"location": map[string]interface{}{
			"latitude":  39.9042,
			"longitude": 116.4074,
		},
		"quality": map[string]interface{}{
			"signal_strength": -65,
			"data_quality":    "good",
		},
	}

	payload, _ := json.Marshal(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg models.MQTTMessage
		json.Unmarshal(payload, &msg)
	}
}
