package mqtt

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDatabase 设置测试数据库
func setupTestDatabase(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Device{},
		&models.UnifiedSensorData{},
		&models.DeviceRuntimeStatus{},
		&models.Alert{},
		&models.AlertRule{},
		&models.SystemConfig{},
	)
	require.NoError(t, err)

	return db
}

// TestMQTTIntegration_CompleteFlow 测试完整的MQTT数据流
func TestMQTTIntegration_CompleteFlow(t *testing.T) {
	// 设置测试数据库
	db := setupTestDatabase(t)

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	require.NoError(t, err)

	// 创建仓库
	repos := &repositories.Repositories{
		Device:            repositories.NewDeviceRepository(db, logger),
		UnifiedSensorData: repositories.NewUnifiedSensorDataRepository(db, logger),
		Alert:             repositories.NewAlertRepository(db, logger),
	}

	// 创建服务
	svcs := &services.Services{
		Alert: services.NewAlertService(repos.Alert, logger),
	}

	// 创建数据处理器
	handler := NewSensorDataHandler(
		repos.UnifiedSensorData,
		repos.Device,
		svcs.Alert,
		logger,
	)

	// 准备测试数据
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

	// 转换为JSON
	payload, err := json.Marshal(testData)
	require.NoError(t, err)

	// 执行处理
	err = handler.HandleMessage("air-quality/hcho/hcho_001/data", payload)
	require.NoError(t, err)

	// 验证数据存储
	ctx := context.Background()

	// 检查传感器数据
	sensorData, err := repos.UnifiedSensorData.GetHistoryByDeviceID(ctx, "hcho_001", 1, 0)
	require.NoError(t, err)
	require.Len(t, sensorData, 1)

	assert.Equal(t, "hcho_001", sensorData[0].DeviceID)
	assert.Equal(t, models.DeviceTypeFormaldehyde, sensorData[0].DeviceType)
	assert.Equal(t, "sensor_hcho_001_01", sensorData[0].SensorID)
	assert.Equal(t, "hcho", sensorData[0].SensorType)
	assert.NotNil(t, sensorData[0].Formaldehyde)
	assert.Equal(t, 0.05, *sensorData[0].Formaldehyde)
	assert.NotNil(t, sensorData[0].Temperature)
	assert.Equal(t, 22.5, *sensorData[0].Temperature)
	assert.NotNil(t, sensorData[0].Humidity)
	assert.Equal(t, 45.0, *sensorData[0].Humidity)
	assert.NotNil(t, sensorData[0].Battery)
	assert.Equal(t, 85, *sensorData[0].Battery)
	assert.NotNil(t, sensorData[0].Latitude)
	assert.Equal(t, 39.9042, *sensorData[0].Latitude)
	assert.NotNil(t, sensorData[0].Longitude)
	assert.Equal(t, 116.4074, *sensorData[0].Longitude)
	assert.NotNil(t, sensorData[0].SignalStrength)
	assert.Equal(t, -65, *sensorData[0].SignalStrength)
	assert.Equal(t, "good", sensorData[0].DataQuality)

	// 设备状态检查已移除，因为DeviceRuntimeStatusRepository已被删除
}

// TestMQTTIntegration_AlertGeneration 测试告警生成
func TestMQTTIntegration_AlertGeneration(t *testing.T) {
	// 设置测试数据库
	db := setupTestDatabase(t)

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	require.NoError(t, err)

	// 创建仓库
	repos := &repositories.Repositories{
		Device:            repositories.NewDeviceRepository(db, logger),
		UnifiedSensorData: repositories.NewUnifiedSensorDataRepository(db, logger),
		Alert:             repositories.NewAlertRepository(db, logger),
	}

	// 创建服务
	svcs := &services.Services{
		Alert: services.NewAlertService(repos.Alert, logger),
	}

	// 创建数据处理器
	handler := NewSensorDataHandler(
		repos.UnifiedSensorData,
		repos.Device,
		svcs.Alert,
		logger,
	)

	// 准备触发告警的测试数据（甲醛浓度超过阈值）
	testData := map[string]interface{}{
		"device_id":   "hcho_001",
		"device_type": "hcho",
		"sensor_id":   "sensor_hcho_001_01",
		"sensor_type": "hcho",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"formaldehyde": 0.12, // 超过0.08的警告阈值
			"temperature":  22.5,
			"humidity":     45.0,
			"battery":      85.0,
		},
	}

	// 转换为JSON
	payload, err := json.Marshal(testData)
	require.NoError(t, err)

	// 执行处理
	err = handler.HandleMessage("air-quality/hcho/hcho_001/data", payload)
	require.NoError(t, err)

	// 验证告警生成
	alerts, err := repos.Alert.GetUnresolved()
	require.NoError(t, err)
	require.Len(t, alerts, 1)

	alert := alerts[0]
	assert.Equal(t, "hcho_001", alert.DeviceID)
	assert.Equal(t, "formaldehyde", alert.Metric)
	assert.Equal(t, 0.12, alert.CurrentValue)
	assert.Equal(t, 0.08, alert.ThresholdValue)
	assert.Equal(t, "critical", alert.Severity)
	assert.Equal(t, "active", alert.Status)
	assert.NotNil(t, alert.Message)
	assert.Contains(t, *alert.Message, "甲醛浓度严重超标")
}

// TestMQTTIntegration_MultipleDevices 测试多设备数据处理
func TestMQTTIntegration_MultipleDevices(t *testing.T) {
	// 设置测试数据库
	db := setupTestDatabase(t)

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	require.NoError(t, err)

	// 创建仓库
	repos := &repositories.Repositories{
		Device:            repositories.NewDeviceRepository(db, logger),
		UnifiedSensorData: repositories.NewUnifiedSensorDataRepository(db, logger),
		Alert:             repositories.NewAlertRepository(db, logger),
	}

	// 创建服务
	svcs := &services.Services{
		Alert: services.NewAlertService(repos.Alert, logger),
	}

	// 创建数据处理器
	handler := NewSensorDataHandler(
		repos.UnifiedSensorData,
		repos.Device,
		svcs.Alert,
		logger,
	)

	// 准备多个设备的测试数据
	devices := []string{"hcho_001", "hcho_002", "hcho_003"}

	for i, deviceID := range devices {
		testData := map[string]interface{}{
			"device_id":   deviceID,
			"device_type": "hcho",
			"sensor_id":   "sensor_" + deviceID + "_01",
			"sensor_type": "hcho",
			"timestamp":   time.Now().Unix(),
			"data": map[string]interface{}{
				"formaldehyde": 0.03 + float64(i)*0.01, // 不同的甲醛浓度
				"temperature":  20.0 + float64(i)*2.0,  // 不同的温度
				"humidity":     40.0 + float64(i)*5.0,  // 不同的湿度
				"battery":      90 - i*10,              // 不同的电池电量
			},
		}

		// 转换为JSON
		payload, err := json.Marshal(testData)
		require.NoError(t, err)

		// 执行处理
		err = handler.HandleMessage("air-quality/hcho/"+deviceID+"/data", payload)
		require.NoError(t, err)
	}

	// 验证所有设备的数据都正确存储
	ctx := context.Background()

	for _, deviceID := range devices {
		// 检查传感器数据
		sensorData, err := repos.UnifiedSensorData.GetHistoryByDeviceID(ctx, deviceID, 1, 0)
		require.NoError(t, err)
		require.Len(t, sensorData, 1)
		assert.Equal(t, deviceID, sensorData[0].DeviceID)

		// 设备状态检查已移除，因为DeviceRuntimeStatusRepository已被删除
	}
}

// TestMQTTIntegration_DataValidation 测试数据验证
func TestMQTTIntegration_DataValidation(t *testing.T) {
	// 设置测试数据库
	db := setupTestDatabase(t)

	// 创建测试日志器
	logger, err := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)
	require.NoError(t, err)

	// 创建仓库
	repos := &repositories.Repositories{
		Device:            repositories.NewDeviceRepository(db, logger),
		UnifiedSensorData: repositories.NewUnifiedSensorDataRepository(db, logger),
		Alert:             repositories.NewAlertRepository(db, logger),
	}

	// 创建服务
	svcs := &services.Services{
		Alert: services.NewAlertService(repos.Alert, logger),
	}

	// 创建数据处理器
	handler := NewSensorDataHandler(
		repos.UnifiedSensorData,
		repos.Device,
		svcs.Alert,
		logger,
	)

	// 测试各种无效数据
	testCases := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
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
		},
		{
			name: "无效的时间戳",
			data: map[string]interface{}{
				"device_id":   "hcho_001",
				"device_type": "hcho",
				"timestamp":   "invalid",
				"data": map[string]interface{}{
					"formaldehyde": 0.05,
				},
			},
			wantErr: true, // 时间戳解析失败会报错
		},
		{
			name: "缺少数据字段",
			data: map[string]interface{}{
				"device_id":   "hcho_001",
				"device_type": "hcho",
				"timestamp":   time.Now().Unix(),
				"data":        map[string]interface{}{},
			},
			wantErr: false, // 缺少数据字段不会报错，只是不会存储
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 转换为JSON
			payload, err := json.Marshal(tc.data)
			require.NoError(t, err)

			// 执行处理
			err = handler.HandleMessage("air-quality/hcho/hcho_001/data", payload)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkMQTTIntegration_HandleMessage 性能测试
func BenchmarkMQTTIntegration_HandleMessage(b *testing.B) {
	// 设置测试数据库
	db := setupTestDatabase(&testing.T{})

	// 创建测试日志器
	logger, _ := utils.NewLogger("info", "console", "stdout", 100, 3, 28, true)

	// 创建仓库
	repos := &repositories.Repositories{
		Device:            repositories.NewDeviceRepository(db, logger),
		UnifiedSensorData: repositories.NewUnifiedSensorDataRepository(db, logger),
		Alert:             repositories.NewAlertRepository(db, logger),
	}

	// 创建服务
	svcs := &services.Services{
		Alert: services.NewAlertService(repos.Alert, logger),
	}

	// 创建数据处理器
	handler := NewSensorDataHandler(
		repos.UnifiedSensorData,
		repos.Device,
		svcs.Alert,
		logger,
	)

	// 准备测试数据
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
	}

	payload, _ := json.Marshal(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.HandleMessage("air-quality/hcho/hcho_001/data", payload)
	}
}
