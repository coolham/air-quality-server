package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Service  ServiceConfig  `mapstructure:"service"`
	MQTT     MQTTConfig     `mapstructure:"mqtt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
	MaxLife  int    `mapstructure:"max_life"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
	Issuer      string `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// MQTTConfig MQTT配置
type MQTTConfig struct {
	Broker               string        `mapstructure:"broker"`
	ClientID             string        `mapstructure:"client_id"`
	Username             string        `mapstructure:"username"`
	Password             string        `mapstructure:"password"`
	KeepAlive            int           `mapstructure:"keep_alive"`
	CleanSession         bool          `mapstructure:"clean_session"`
	QoS                  int           `mapstructure:"qos"`
	AutoReconnect        bool          `mapstructure:"auto_reconnect"`
	MaxReconnectInterval int           `mapstructure:"max_reconnect_interval"`
	ReconnectDelay       int           `mapstructure:"reconnect_delay"`
	ConnectTimeout       int           `mapstructure:"connect_timeout"`
	WriteTimeout         int           `mapstructure:"write_timeout"`
	ReadTimeout          int           `mapstructure:"read_timeout"`
	Topics               TopicConfig   `mapstructure:"topics"`
	PublishPrefix        string        `mapstructure:"publish_prefix"`
	Message              MessageConfig `mapstructure:"message"`
	Device               DeviceConfig  `mapstructure:"device"`
	Alert                AlertConfig   `mapstructure:"alert"`
}

// TopicConfig 主题配置
type TopicConfig struct {
	DeviceStatus   string `mapstructure:"device_status"`
	DeviceResponse string `mapstructure:"device_response"`
}

// MessageConfig 消息配置
type MessageConfig struct {
	MaxSize    int `mapstructure:"max_size"`
	BufferSize int `mapstructure:"buffer_size"`
	BatchSize  int `mapstructure:"batch_size"`
}

// DeviceConfig 设备配置
type DeviceConfig struct {
	OfflineTimeout    int `mapstructure:"offline_timeout"`
	HeartbeatInterval int `mapstructure:"heartbeat_interval"`
	ReportInterval    int `mapstructure:"report_interval"`
}

// AlertConfig 告警配置
type AlertConfig struct {
	FormaldehydeWarning  float64 `mapstructure:"formaldehyde_warning"`
	FormaldehydeCritical float64 `mapstructure:"formaldehyde_critical"`
	BatteryLow           int     `mapstructure:"battery_low"`
	SignalWeak           int     `mapstructure:"signal_weak"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 120)

	// 数据库默认配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "air_quality")
	viper.SetDefault("database.password", "air_quality123")
	viper.SetDefault("database.database", "air_quality")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle", 10)
	viper.SetDefault("database.max_open", 100)
	viper.SetDefault("database.max_life", 3600)

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT默认配置
	viper.SetDefault("jwt.secret", "air-quality-secret-key")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.issuer", "air-quality-server")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 28)
	viper.SetDefault("log.compress", true)

	// 服务默认配置
	viper.SetDefault("service.name", "air-quality-server")
	viper.SetDefault("service.version", "1.0.0")
	viper.SetDefault("service.environment", "development")
	viper.SetDefault("service.debug", false)

	// MQTT默认配置
	viper.SetDefault("mqtt.broker", "tcp://localhost:1883")
	viper.SetDefault("mqtt.client_id", "air-quality-server")
	viper.SetDefault("mqtt.username", "admin")
	viper.SetDefault("mqtt.password", "password")
	viper.SetDefault("mqtt.keep_alive", 60)
	viper.SetDefault("mqtt.clean_session", true)
	viper.SetDefault("mqtt.qos", 1)
	viper.SetDefault("mqtt.auto_reconnect", true)
	viper.SetDefault("mqtt.max_reconnect_interval", 300)
	viper.SetDefault("mqtt.reconnect_delay", 5)
	viper.SetDefault("mqtt.connect_timeout", 30)
	viper.SetDefault("mqtt.write_timeout", 10)
	viper.SetDefault("mqtt.read_timeout", 10)
	viper.SetDefault("mqtt.topics.formaldehyde_data", "air-quality/hcho/+/data")
	viper.SetDefault("mqtt.topics.device_status", "air-quality/hcho/+/status")
	viper.SetDefault("mqtt.topics.device_response", "air-quality/hcho/+/response")
	viper.SetDefault("mqtt.publish_prefix", "air-quality/hcho")
	viper.SetDefault("mqtt.message.max_size", 1048576)
	viper.SetDefault("mqtt.message.buffer_size", 1000)
	viper.SetDefault("mqtt.message.batch_size", 100)
	viper.SetDefault("mqtt.device.offline_timeout", 300)
	viper.SetDefault("mqtt.device.heartbeat_interval", 60)
	viper.SetDefault("mqtt.device.report_interval", 300)
	viper.SetDefault("mqtt.alert.formaldehyde_warning", 0.08)
	viper.SetDefault("mqtt.alert.formaldehyde_critical", 0.1)
	viper.SetDefault("mqtt.alert.battery_low", 20)
	viper.SetDefault("mqtt.alert.signal_weak", -80)
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535之间")
	}

	if config.Database.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}

	if config.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}

	if config.Database.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if config.Redis.Host == "" {
		return fmt.Errorf("Redis主机地址不能为空")
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetServerAddr 获取服务器地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Service.Environment == "development"
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Service.Environment == "production"
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			Host:         getEnvString("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 30),
			IdleTimeout:  getEnvInt("SERVER_IDLE_TIMEOUT", 120),
		},
		Database: DatabaseConfig{
			Host:     getEnvString("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 3306),
			Username: getEnvString("DB_USERNAME", "air_quality"),
			Password: getEnvString("DB_PASSWORD", "air_quality123"),
			Database: getEnvString("DB_NAME", "air_quality"),
			Charset:  getEnvString("DB_CHARSET", "utf8mb4"),
			MaxIdle:  getEnvInt("DB_MAX_IDLE", 10),
			MaxOpen:  getEnvInt("DB_MAX_OPEN", 100),
			MaxLife:  getEnvInt("DB_MAX_LIFE", 3600),
		},
		Redis: RedisConfig{
			Host:     getEnvString("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnvString("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 10),
		},
		JWT: JWTConfig{
			Secret:      getEnvString("JWT_SECRET", "air-quality-secret-key"),
			ExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
			Issuer:      getEnvString("JWT_ISSUER", "air-quality-server"),
		},
		Log: LogConfig{
			Level:      getEnvString("LOG_LEVEL", "info"),
			Format:     getEnvString("LOG_FORMAT", "json"),
			Output:     getEnvString("LOG_OUTPUT", "stdout"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 3),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 28),
			Compress:   getEnvBool("LOG_COMPRESS", true),
		},
		Service: ServiceConfig{
			Name:        getEnvString("SERVICE_NAME", "air-quality-server"),
			Version:     getEnvString("SERVICE_VERSION", "1.0.0"),
			Environment: getEnvString("ENVIRONMENT", "development"),
			Debug:       getEnvBool("DEBUG", false),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// 辅助函数
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
