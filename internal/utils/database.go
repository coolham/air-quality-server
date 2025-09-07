package utils

import (
	"context"
	"fmt"
	"time"

	"air-quality-server/internal/config"
	"air-quality-server/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库连接管理器
type Database struct {
	DB     *gorm.DB
	config *config.DatabaseConfig
	logger Logger
}

// NewDatabase 创建新的数据库连接
func NewDatabase(cfg *config.DatabaseConfig, log Logger) (*Database, error) {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	// 设置GORM日志级别
	var gormLogger logger.Interface
	if log != nil {
		gormLogger = &GormLogger{logger: log}
	} else {
		gormLogger = logger.Default
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLife) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	database := &Database{
		DB:     db,
		config: cfg,
		logger: log,
	}

	if log != nil {
		log.Info("数据库连接成功",
			String("host", cfg.Host),
			Int("port", cfg.Port),
			String("database", cfg.Database),
		)
	}

	return database, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Ping 测试数据库连接
func (d *Database) Ping() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}
	return fmt.Errorf("数据库连接未初始化")
}

// GetStats 获取数据库连接统计信息
func (d *Database) GetStats() (map[string]interface{}, error) {
	if d.DB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}

// AutoMigrate 自动迁移数据模型
func (d *Database) AutoMigrate() error {
	if d.DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 获取所有数据模型
	models := getAllModels()

	// 执行自动迁移
	if err := d.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	if d.logger != nil {
		d.logger.Info("数据库自动迁移完成", Int("models_count", len(models)))
	}

	return nil
}

// getAllModels 获取所有数据模型
func getAllModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Device{},
		&models.UnifiedSensorData{},
		&models.DeviceRuntimeStatus{},
		&models.Alert{},
		&models.AlertRule{},
		&models.SystemConfig{},
		&models.FormaldehydeData{},
		&models.FormaldehydeDeviceStatus{},
		&models.FormaldehydeDeviceConfig{},
	}
}

// GormLogger GORM日志适配器
type GormLogger struct {
	logger Logger
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info 信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn 警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error 错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

// Trace 跟踪日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.logger.Error("SQL执行失败",
			String("sql", sql),
			Int64("rows", rows),
			Duration("elapsed", elapsed),
			ErrorField(err),
		)
	} else {
		l.logger.Debug("SQL执行成功",
			String("sql", sql),
			Int64("rows", rows),
			Duration("elapsed", elapsed),
		)
	}
}
