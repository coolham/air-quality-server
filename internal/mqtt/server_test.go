package mqtt

import (
	"air-quality-server/internal/config"
	"air-quality-server/internal/utils"
	"testing"
	"time"
)

// TestNewServer 测试创建MQTT服务器
func TestNewServer(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:         "localhost:1884",
		ClientID:       "test-client",
		KeepAlive:      60,
		CleanSession:   true,
		QoS:            1,
		AutoReconnect:  true,
		ConnectTimeout: 30,
	}

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	// 创建数据处理器（使用nil，因为这是单元测试）
	var sensorDataHandler *SensorDataHandler = nil

	// 创建服务器
	server := NewServer(cfg, logger, sensorDataHandler)

	// 验证服务器创建
	if server == nil {
		t.Fatal("服务器创建失败")
	}

	if server.config != cfg {
		t.Error("配置未正确设置")
	}

	if server.logger != logger {
		t.Error("日志器未正确设置")
	}

	if server.sensorDataHandler != sensorDataHandler {
		t.Error("数据处理器未正确设置")
	}

	if server.running {
		t.Error("服务器不应在创建时运行")
	}
}

// TestServerStartStop 测试服务器启动和停止
func TestServerStartStop(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:         "localhost:1884", // 使用不同端口避免冲突
		ClientID:       "test-client",
		KeepAlive:      60,
		CleanSession:   true,
		QoS:            1,
		AutoReconnect:  true,
		ConnectTimeout: 30,
	}

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	// 创建服务器
	server := NewServer(cfg, logger, nil)

	// 测试启动
	err = server.Start()
	if err != nil {
		t.Fatalf("服务器启动失败: %v", err)
	}

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	// 验证服务器状态
	if !server.IsRunning() {
		t.Error("服务器应该正在运行")
	}

	// 测试停止
	server.Stop()

	// 等待服务器停止
	time.Sleep(100 * time.Millisecond)

	// 验证服务器状态
	if server.IsRunning() {
		t.Error("服务器应该已停止")
	}
}

// TestServerPublish 测试服务器发布消息
func TestServerPublish(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:         "localhost:1884",
		ClientID:       "test-client",
		KeepAlive:      60,
		CleanSession:   true,
		QoS:            1,
		AutoReconnect:  true,
		ConnectTimeout: 30,
	}

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	// 创建服务器
	server := NewServer(cfg, logger, nil)

	// 测试未启动状态下的发布
	err = server.Publish("test/topic", "test message")
	if err == nil {
		t.Error("未启动的服务器不应能发布消息")
	}

	// 启动服务器
	err = server.Start()
	if err != nil {
		t.Fatalf("服务器启动失败: %v", err)
	}

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	// 测试发布字符串消息
	err = server.Publish("test/topic", "test message")
	if err != nil {
		t.Errorf("发布字符串消息失败: %v", err)
	}

	// 测试发布字节数组消息
	err = server.Publish("test/topic", []byte("test message"))
	if err != nil {
		t.Errorf("发布字节数组消息失败: %v", err)
	}

	// 测试发布JSON消息
	testData := map[string]interface{}{
		"device_id": "test_device",
		"data":      "test_data",
	}
	err = server.Publish("test/topic", testData)
	if err != nil {
		t.Errorf("发布JSON消息失败: %v", err)
	}

	// 停止服务器
	server.Stop()
}

// TestServerGetStatus 测试获取服务器状态
func TestServerGetStatus(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:         "localhost:1884",
		ClientID:       "test-client",
		KeepAlive:      60,
		CleanSession:   true,
		QoS:            1,
		AutoReconnect:  true,
		ConnectTimeout: 30,
	}

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	// 创建服务器
	server := NewServer(cfg, logger, nil)

	// 测试未启动状态
	status := server.GetStatus()
	if status["running"] != false {
		t.Error("未启动的服务器状态应为false")
	}

	// 启动服务器
	err = server.Start()
	if err != nil {
		t.Fatalf("服务器启动失败: %v", err)
	}

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	// 测试启动状态
	status = server.GetStatus()
	if status["running"] != true {
		t.Error("启动的服务器状态应为true")
	}

	if status["address"] != ":1883" {
		t.Error("服务器地址不正确")
	}

	if status["client_id"] != cfg.ClientID {
		t.Error("客户端ID不正确")
	}

	// 停止服务器
	server.Stop()
}

// TestIsSensorDataTopic 测试传感器数据主题识别
func TestIsSensorDataTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected bool
	}{
		{
			name:     "有效的甲醛传感器主题",
			topic:    "air-quality/hcho/hcho_001/data",
			expected: true,
		},
		{
			name:     "有效的ESP32传感器主题",
			topic:    "air-quality/esp32/esp32_001/data",
			expected: true,
		},
		{
			name:     "有效的通用传感器主题",
			topic:    "air-quality/sensor/sensor_001/data",
			expected: true,
		},
		{
			name:     "无效的主题格式-部分不足",
			topic:    "air-quality/hcho/hcho_001",
			expected: false,
		},
		{
			name:     "无效的主题格式-前缀错误",
			topic:    "wrong-prefix/hcho/hcho_001/data",
			expected: false,
		},
		{
			name:     "无效的主题格式-后缀错误",
			topic:    "air-quality/hcho/hcho_001/status",
			expected: false,
		},
		{
			name:     "无效的主题格式-设备类型错误",
			topic:    "air-quality/invalid/hcho_001/data",
			expected: false,
		},
		{
			name:     "无效的主题格式-设备ID为空",
			topic:    "air-quality/hcho//data",
			expected: false,
		},
		{
			name:     "空主题",
			topic:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSensorDataTopic(tt.topic)
			if result != tt.expected {
				t.Errorf("isSensorDataTopic(%s) = %v, expected %v", tt.topic, result, tt.expected)
			}
		})
	}
}

// TestMessageHandlerHookID 测试消息处理钩子ID
func TestMessageHandlerHookID(t *testing.T) {
	hook := &MessageHandlerHook{}
	id := hook.ID()
	if id != "message-handler" {
		t.Errorf("钩子ID应为'message-handler'，实际为'%s'", id)
	}
}

// TestMessageHandlerHookProvides 测试消息处理钩子提供的事件
func TestMessageHandlerHookProvides(t *testing.T) {
	hook := &MessageHandlerHook{}

	// 测试支持的事件
	supportedEvents := []byte{
		0x10, // OnConnect
		0x20, // OnDisconnect
		0x30, // OnPublish
		// 添加其他支持的事件...
	}

	for _, event := range supportedEvents {
		if !hook.Provides(event) {
			t.Errorf("钩子应该支持事件 0x%02X", event)
		}
	}

	// 测试不支持的事件
	unsupportedEvent := byte(0xFF)
	if hook.Provides(unsupportedEvent) {
		t.Errorf("钩子不应该支持事件 0x%02X", unsupportedEvent)
	}
}

// TestMessageHandlerHookInit 测试消息处理钩子初始化
func TestMessageHandlerHookInit(t *testing.T) {
	hook := &MessageHandlerHook{}

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	// 测试配置
	config := map[string]interface{}{
		"logger":            logger,
		"sensorDataHandler": (*SensorDataHandler)(nil),
	}

	// 测试初始化
	err = hook.Init(config)
	if err != nil {
		t.Errorf("钩子初始化失败: %v", err)
	}

	// 验证初始化结果
	if hook.logger != logger {
		t.Error("日志器未正确设置")
	}

	// 测试无效配置
	invalidConfig := map[string]interface{}{
		"invalid": "config",
	}

	err = hook.Init(invalidConfig)
	if err == nil {
		t.Error("无效配置应该导致初始化失败")
	}
}

// BenchmarkIsSensorDataTopic 性能测试
func BenchmarkIsSensorDataTopic(b *testing.B) {
	topic := "air-quality/hcho/hcho_001/data"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isSensorDataTopic(topic)
	}
}

// BenchmarkServerStartStop 性能测试
func BenchmarkServerStartStop(b *testing.B) {
	cfg := &config.MQTTConfig{
		Broker:         "localhost:1884",
		ClientID:       "test-client",
		KeepAlive:      60,
		CleanSession:   true,
		QoS:            1,
		AutoReconnect:  true,
		ConnectTimeout: 30,
	}

	logger, _ := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := NewServer(cfg, logger, nil)
		server.Start()
		time.Sleep(10 * time.Millisecond)
		server.Stop()
	}
}
