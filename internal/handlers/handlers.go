package handlers

// Handlers 处理器集合
type Handlers struct {
	Device     *DeviceHandler
	AirQuality *AirQualityHandler
	User       *UserHandler
	Alert      *AlertHandler
	Config     *ConfigHandler
}
