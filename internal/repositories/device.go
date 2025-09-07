package repositories

import (
	"context"
	"fmt"
	"time"

	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

// DeviceRepository 设备仓储接口
type DeviceRepository interface {
	BaseRepository[models.Device]
	GetByDeviceID(ctx context.Context, deviceID string) (*models.Device, error)
	UpdateStatus(ctx context.Context, deviceID string, status string) error
	GetRealtimeStatus(ctx context.Context, deviceID string) (*models.DeviceRealtimeStatus, error)
	GetRealtimeStatusList(ctx context.Context, req *models.DeviceListRequest) (*ListResponse[models.DeviceRealtimeStatus], error)
	GetStatistics(ctx context.Context, deviceID string, startTime, endTime time.Time) (*models.DeviceStatistics, error)
	GetOnlineDevices(ctx context.Context) ([]models.Device, error)
	GetOfflineDevices(ctx context.Context, duration time.Duration) ([]models.Device, error)
}

// deviceRepository 设备仓储实现
type deviceRepository struct {
	BaseRepository[models.Device]
	db     *gorm.DB
	logger utils.Logger
}

// NewDeviceRepository 创建设备仓储
func NewDeviceRepository(db *gorm.DB, logger utils.Logger) DeviceRepository {
	return &deviceRepository{
		BaseRepository: NewBaseRepository[models.Device](db, logger),
		db:             db,
		logger:         logger,
	}
}

// GetByDeviceID 根据设备ID获取设备
func (r *deviceRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.Device, error) {
	var device models.Device
	if err := r.db.WithContext(ctx).Where("id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("设备不存在")
		}
		r.logger.Error("根据设备ID获取设备失败", utils.String("device_id", deviceID), utils.ErrorField(err))
		return nil, fmt.Errorf("获取设备失败: %w", err)
	}
	return &device, nil
}

// UpdateStatus 更新设备状态
func (r *deviceRepository) UpdateStatus(ctx context.Context, deviceID string, status string) error {
	updateData := &models.Device{Status: models.DeviceStatus(status)}
	if err := r.db.WithContext(ctx).Model(&models.Device{}).Where("id = ?", deviceID).Updates(updateData).Error; err != nil {
		r.logger.Error("更新设备状态失败", utils.String("device_id", deviceID), utils.String("status", status), utils.ErrorField(err))
		return fmt.Errorf("更新设备状态失败: %w", err)
	}
	return nil
}

// GetRealtimeStatus 获取设备实时状态
func (r *deviceRepository) GetRealtimeStatus(ctx context.Context, deviceID string) (*models.DeviceRealtimeStatus, error) {
	var status models.DeviceRealtimeStatus

	// 查询设备信息和最新数据
	query := `
		SELECT 
			d.id, d.name, d.type, d.location_latitude, d.location_longitude, d.location_address,
			d.status, d.config, d.created_at, d.updated_at,
			aqd.timestamp as last_data_time,
			aqd.pm25, aqd.pm10, aqd.co2, aqd.temperature, aqd.humidity, aqd.pressure, aqd.quality_score,
			CASE 
				WHEN TIMESTAMPDIFF(MINUTE, aqd.timestamp, NOW()) > 10 THEN 'offline'
				ELSE d.status
			END as realtime_status
		FROM devices d
		LEFT JOIN (
			SELECT 
				device_id, timestamp, pm25, pm10, co2, temperature, humidity, pressure, quality_score,
				ROW_NUMBER() OVER (PARTITION BY device_id ORDER BY timestamp DESC) as rn
			FROM air_quality_data
		) aqd ON d.id = aqd.device_id AND aqd.rn = 1
		WHERE d.id = ?
	`

	if err := r.db.WithContext(ctx).Raw(query, deviceID).Scan(&status).Error; err != nil {
		r.logger.Error("获取设备实时状态失败", utils.String("device_id", deviceID), utils.ErrorField(err))
		return nil, fmt.Errorf("获取设备实时状态失败: %w", err)
	}

	return &status, nil
}

// GetRealtimeStatusList 获取设备实时状态列表
func (r *deviceRepository) GetRealtimeStatusList(ctx context.Context, req *models.DeviceListRequest) (*ListResponse[models.DeviceRealtimeStatus], error) {
	var statuses []models.DeviceRealtimeStatus
	var total int64

	// 构建查询
	query := `
		SELECT 
			d.id, d.name, d.type, d.location_latitude, d.location_longitude, d.location_address,
			d.status, d.config, d.created_at, d.updated_at,
			aqd.timestamp as last_data_time,
			aqd.pm25, aqd.pm10, aqd.co2, aqd.temperature, aqd.humidity, aqd.pressure, aqd.quality_score,
			CASE 
				WHEN TIMESTAMPDIFF(MINUTE, aqd.timestamp, NOW()) > 10 THEN 'offline'
				ELSE d.status
			END as realtime_status
		FROM devices d
		LEFT JOIN (
			SELECT 
				device_id, timestamp, pm25, pm10, co2, temperature, humidity, pressure, quality_score,
				ROW_NUMBER() OVER (PARTITION BY device_id ORDER BY timestamp DESC) as rn
			FROM air_quality_data
		) aqd ON d.id = aqd.device_id AND aqd.rn = 1
		WHERE 1=1
	`

	// 添加条件
	args := []interface{}{}
	if req.Status != "" {
		query += " AND d.status = ?"
		args = append(args, req.Status)
	}
	if req.Type != "" {
		query += " AND d.type = ?"
		args = append(args, req.Type)
	}
	if req.Keyword != "" {
		query += " AND (d.name LIKE ? OR d.id LIKE ?)"
		keyword := "%" + req.Keyword + "%"
		args = append(args, keyword, keyword)
	}

	// 计算总数
	countQuery := "SELECT COUNT(*) FROM devices d WHERE 1=1"
	countArgs := []interface{}{}
	if req.Status != "" {
		countQuery += " AND d.status = ?"
		countArgs = append(countArgs, req.Status)
	}
	if req.Type != "" {
		countQuery += " AND d.type = ?"
		countArgs = append(countArgs, req.Type)
	}
	if req.Keyword != "" {
		countQuery += " AND (d.name LIKE ? OR d.id LIKE ?)"
		keyword := "%" + req.Keyword + "%"
		countArgs = append(countArgs, keyword, keyword)
	}

	if err := r.db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
		r.logger.Error("计算设备总数失败", utils.ErrorField(err))
		return nil, fmt.Errorf("计算设备总数失败: %w", err)
	}

	// 设置排序
	query += " ORDER BY d.created_at DESC"

	// 设置分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query += " LIMIT ? OFFSET ?"
		args = append(args, req.PageSize, offset)
	}

	// 执行查询
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&statuses).Error; err != nil {
		r.logger.Error("查询设备实时状态列表失败", utils.ErrorField(err))
		return nil, fmt.Errorf("查询设备实时状态列表失败: %w", err)
	}

	return &ListResponse[models.DeviceRealtimeStatus]{
		Data:     statuses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetStatistics 获取设备统计信息
func (r *deviceRepository) GetStatistics(ctx context.Context, deviceID string, startTime, endTime time.Time) (*models.DeviceStatistics, error) {
	var stats models.DeviceStatistics

	// 查询统计数据
	query := `
		SELECT 
			device_id,
			COUNT(*) as total_data_count,
			AVG(pm25) as average_pm25,
			AVG(pm10) as average_pm10,
			AVG(co2) as average_co2,
			AVG(temperature) as average_temp,
			AVG(humidity) as average_humidity,
			AVG(pressure) as average_pressure
		FROM air_quality_data
		WHERE device_id = ? AND timestamp BETWEEN ? AND ?
		GROUP BY device_id
	`

	if err := r.db.WithContext(ctx).Raw(query, deviceID, startTime, endTime).Scan(&stats).Error; err != nil {
		r.logger.Error("获取设备统计信息失败", utils.String("device_id", deviceID), utils.ErrorField(err))
		return nil, fmt.Errorf("获取设备统计信息失败: %w", err)
	}

	stats.DeviceID = deviceID

	return &stats, nil
}

// GetOnlineDevices 获取在线设备
func (r *deviceRepository) GetOnlineDevices(ctx context.Context) ([]models.Device, error) {
	var devices []models.Device
	if err := r.db.WithContext(ctx).Where("status = ?", models.DeviceStatusOnline).Find(&devices).Error; err != nil {
		r.logger.Error("获取在线设备失败", utils.ErrorField(err))
		return nil, fmt.Errorf("获取在线设备失败: %w", err)
	}
	return devices, nil
}

// GetOfflineDevices 获取离线设备
func (r *deviceRepository) GetOfflineDevices(ctx context.Context, duration time.Duration) ([]models.Device, error) {
	var devices []models.Device

	// 查询超过指定时间没有数据的设备
	query := `
		SELECT d.* FROM devices d
		LEFT JOIN (
			SELECT device_id, MAX(timestamp) as last_timestamp
			FROM air_quality_data
			GROUP BY device_id
		) aqd ON d.id = aqd.device_id
		WHERE d.status = 'online' AND (
			aqd.last_timestamp IS NULL OR 
			aqd.last_timestamp < ?
		)
	`

	cutoffTime := time.Now().Add(-duration)
	if err := r.db.WithContext(ctx).Raw(query, cutoffTime).Scan(&devices).Error; err != nil {
		r.logger.Error("获取离线设备失败", utils.ErrorField(err))
		return nil, fmt.Errorf("获取离线设备失败: %w", err)
	}

	return devices, nil
}
