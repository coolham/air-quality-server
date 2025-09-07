package services

// Services 服务层集合
type Services struct {
	Device            DeviceService
	AirQuality        AirQualityService
	UnifiedSensorData UnifiedSensorDataService
	User              UserService
	Alert             AlertService
	Config            ConfigService
}
