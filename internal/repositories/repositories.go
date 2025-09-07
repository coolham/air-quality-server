package repositories

// Repositories 仓储层集合
type Repositories struct {
	Device              DeviceRepository
	AirQuality          AirQualityRepository
	UnifiedSensorData   UnifiedSensorDataRepository
	DeviceRuntimeStatus DeviceRuntimeStatusRepository
	User                UserRepository
	Alert               AlertRepository
	Config              ConfigRepository
}
