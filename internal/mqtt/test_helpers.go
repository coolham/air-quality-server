package mqtt

import (
	"air-quality-server/internal/models"
	"encoding/json"
	"time"
)

// TestDataGenerator 测试数据生成器
type TestDataGenerator struct{}

// NewTestDataGenerator 创建测试数据生成器
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{}
}

// GenerateValidSensorData 生成有效的传感器数据
func (g *TestDataGenerator) GenerateValidSensorData(deviceID string) map[string]interface{} {
	return map[string]interface{}{
		"device_id":   deviceID,
		"device_type": "hcho",
		"sensor_id":   "sensor_" + deviceID + "_01",
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
			"address":   "北京市朝阳区",
		},
		"quality": map[string]interface{}{
			"signal_strength": -65,
			"data_quality":    "good",
		},
	}
}

// GenerateAlertTriggeringData 生成触发告警的数据
func (g *TestDataGenerator) GenerateAlertTriggeringData(deviceID string, formaldehydeLevel float64) map[string]interface{} {
	data := g.GenerateValidSensorData(deviceID)
	data["data"].(map[string]interface{})["formaldehyde"] = formaldehydeLevel
	return data
}

// GenerateInvalidData 生成无效数据
func (g *TestDataGenerator) GenerateInvalidData() map[string]interface{} {
	return map[string]interface{}{
		"device_type": "hcho",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"formaldehyde": 0.05,
		},
	}
}

// GenerateDataWithMissingFields 生成缺少字段的数据
func (g *TestDataGenerator) GenerateDataWithMissingFields(deviceID string) map[string]interface{} {
	return map[string]interface{}{
		"device_id":   deviceID,
		"device_type": "hcho",
		"timestamp":   time.Now().Unix(),
		"data":        map[string]interface{}{},
	}
}

// GenerateDataWithInvalidTimestamp 生成无效时间戳的数据
func (g *TestDataGenerator) GenerateDataWithInvalidTimestamp(deviceID string) map[string]interface{} {
	return map[string]interface{}{
		"device_id":   deviceID,
		"device_type": "hcho",
		"timestamp":   "invalid_timestamp",
		"data": map[string]interface{}{
			"formaldehyde": 0.05,
		},
	}
}

// GenerateMultipleDeviceData 生成多设备数据
func (g *TestDataGenerator) GenerateMultipleDeviceData(deviceIDs []string) []map[string]interface{} {
	var dataList []map[string]interface{}

	for i, deviceID := range deviceIDs {
		data := g.GenerateValidSensorData(deviceID)
		// 为每个设备生成不同的数据
		data["data"].(map[string]interface{})["formaldehyde"] = 0.03 + float64(i)*0.01
		data["data"].(map[string]interface{})["temperature"] = 20.0 + float64(i)*2.0
		data["data"].(map[string]interface{})["humidity"] = 40.0 + float64(i)*5.0
		data["data"].(map[string]interface{})["battery"] = 90 - i*10
		dataList = append(dataList, data)
	}

	return dataList
}

// ToJSON 将数据转换为JSON字节数组
func (g *TestDataGenerator) ToJSON(data map[string]interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}

// ToJSONList 将数据列表转换为JSON字节数组列表
func (g *TestDataGenerator) ToJSONList(dataList []map[string]interface{}) [][]byte {
	var jsonList [][]byte
	for _, data := range dataList {
		jsonList = append(jsonList, g.ToJSON(data))
	}
	return jsonList
}

// TestTopicGenerator 测试主题生成器
type TestTopicGenerator struct{}

// NewTestTopicGenerator 创建测试主题生成器
func NewTestTopicGenerator() *TestTopicGenerator {
	return &TestTopicGenerator{}
}

// GenerateValidTopic 生成有效主题
func (g *TestTopicGenerator) GenerateValidTopic(deviceType, deviceID string) string {
	return "air-quality/" + deviceType + "/" + deviceID + "/data"
}

// GenerateInvalidTopic 生成无效主题
func (g *TestTopicGenerator) GenerateInvalidTopic() string {
	return "invalid/topic/format"
}

// GenerateTopicsForMultipleDevices 为多个设备生成主题
func (g *TestTopicGenerator) GenerateTopicsForMultipleDevices(deviceType string, deviceIDs []string) []string {
	var topics []string
	for _, deviceID := range deviceIDs {
		topics = append(topics, g.GenerateValidTopic(deviceType, deviceID))
	}
	return topics
}

// TestAssertions 测试断言辅助函数
type TestAssertions struct{}

// NewTestAssertions 创建测试断言辅助函数
func NewTestAssertions() *TestAssertions {
	return &TestAssertions{}
}

// AssertSensorDataFields 断言传感器数据字段
func (a *TestAssertions) AssertSensorDataFields(t interface{}, data *models.UnifiedSensorData, expected map[string]interface{}) {
	// 这里可以使用testify的断言，但为了保持简单，我们只提供结构
	// 实际使用时可以调用具体的断言函数
}

// AssertDeviceStatusFields 断言设备状态字段
func (a *TestAssertions) AssertDeviceStatusFields(t interface{}, status *models.DeviceRuntimeStatus, expected map[string]interface{}) {
	// 这里可以使用testify的断言，但为了保持简单，我们只提供结构
	// 实际使用时可以调用具体的断言函数
}

// AssertAlertFields 断言告警字段
func (a *TestAssertions) AssertAlertFields(t interface{}, alert *models.Alert, expected map[string]interface{}) {
	// 这里可以使用testify的断言，但为了保持简单，我们只提供结构
	// 实际使用时可以调用具体的断言函数
}

// TestConstants 测试常量
const (
	// 测试设备ID
	TestDeviceID1 = "hcho_001"
	TestDeviceID2 = "hcho_002"
	TestDeviceID3 = "hcho_003"

	// 测试传感器ID
	TestSensorID1 = "sensor_hcho_001_01"
	TestSensorID2 = "sensor_hcho_002_01"
	TestSensorID3 = "sensor_hcho_003_01"

	// 测试设备类型
	TestDeviceTypeHCHO   = "hcho"
	TestDeviceTypeESP32  = "esp32"
	TestDeviceTypeSensor = "sensor"

	// 测试数据值
	TestFormaldehydeNormal   = 0.05
	TestFormaldehydeWarning  = 0.08
	TestFormaldehydeCritical = 0.12

	TestTemperature = 22.5
	TestHumidity    = 45.0
	TestBattery     = 85

	// 测试位置
	TestLatitude  = 39.9042
	TestLongitude = 116.4074
	TestAddress   = "北京市朝阳区"

	// 测试信号强度
	TestSignalStrength = -65

	// 测试数据质量
	TestDataQuality = "good"
)

// TestMQTTMessage 测试MQTT消息结构
type TestMQTTMessage struct {
	DeviceID   string                 `json:"device_id"`
	DeviceType string                 `json:"device_type"`
	SensorID   string                 `json:"sensor_id"`
	SensorType string                 `json:"sensor_type"`
	Timestamp  int64                  `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	Location   *TestLocationInfo      `json:"location,omitempty"`
	Quality    *TestQualityInfo       `json:"quality,omitempty"`
}

// TestLocationInfo 测试位置信息
type TestLocationInfo struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Address   string   `json:"address,omitempty"`
}

// TestQualityInfo 测试质量信息
type TestQualityInfo struct {
	SignalStrength *int   `json:"signal_strength,omitempty"`
	DataQuality    string `json:"data_quality,omitempty"`
}

// CreateTestMQTTMessage 创建测试MQTT消息
func CreateTestMQTTMessage(deviceID string) *TestMQTTMessage {
	latitude := TestLatitude
	longitude := TestLongitude
	signalStrength := TestSignalStrength

	return &TestMQTTMessage{
		DeviceID:   deviceID,
		DeviceType: TestDeviceTypeHCHO,
		SensorID:   "sensor_" + deviceID + "_01",
		SensorType: TestDeviceTypeHCHO,
		Timestamp:  time.Now().Unix(),
		Data: map[string]interface{}{
			"formaldehyde": TestFormaldehydeNormal,
			"temperature":  TestTemperature,
			"humidity":     TestHumidity,
			"battery":      TestBattery,
		},
		Location: &TestLocationInfo{
			Latitude:  &latitude,
			Longitude: &longitude,
			Address:   TestAddress,
		},
		Quality: &TestQualityInfo{
			SignalStrength: &signalStrength,
			DataQuality:    TestDataQuality,
		},
	}
}

// ToJSON 将测试MQTT消息转换为JSON
func (m *TestMQTTMessage) ToJSON() []byte {
	jsonData, _ := json.Marshal(m)
	return jsonData
}

// SetFormaldehydeLevel 设置甲醛浓度
func (m *TestMQTTMessage) SetFormaldehydeLevel(level float64) {
	m.Data["formaldehyde"] = level
}

// SetBatteryLevel 设置电池电量
func (m *TestMQTTMessage) SetBatteryLevel(level float64) {
	m.Data["battery"] = level
}

// SetTemperature 设置温度
func (m *TestMQTTMessage) SetTemperature(temp float64) {
	m.Data["temperature"] = temp
}

// SetHumidity 设置湿度
func (m *TestMQTTMessage) SetHumidity(humidity float64) {
	m.Data["humidity"] = humidity
}
