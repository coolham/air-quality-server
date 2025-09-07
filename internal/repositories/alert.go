package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

// AlertRepository 告警仓储接口
type AlertRepository interface {
	BaseRepository[models.Alert]
	GetByDeviceID(deviceID string) ([]models.Alert, error)
	GetByStatus(status string) ([]models.Alert, error)
	GetByType(alertType string) ([]models.Alert, error)
	GetUnresolved() ([]models.Alert, error)
	MarkAsResolved(alertID uint) error
	GetByTimeRange(startTime, endTime int64) ([]models.Alert, error)
}

// alertRepository 告警仓储实现
type alertRepository struct {
	*baseRepository[models.Alert]
	db     *gorm.DB
	logger utils.Logger
}

// NewAlertRepository 创建告警仓储
func NewAlertRepository(db *gorm.DB, logger utils.Logger) AlertRepository {
	return &alertRepository{
		baseRepository: NewBaseRepository[models.Alert](db, logger).(*baseRepository[models.Alert]),
		db:             db,
		logger:         logger,
	}
}

// GetByDeviceID 根据设备ID获取告警
func (r *alertRepository) GetByDeviceID(deviceID string) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("device_id = ?", deviceID).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		r.logger.Error("根据设备ID获取告警失败", utils.ErrorField(err), utils.String("device_id", deviceID))
		return nil, err
	}
	return alerts, nil
}

// GetByStatus 根据状态获取告警
func (r *alertRepository) GetByStatus(status string) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("status = ?", status).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		r.logger.Error("根据状态获取告警失败", utils.ErrorField(err), utils.String("status", status))
		return nil, err
	}
	return alerts, nil
}

// GetByType 根据类型获取告警
func (r *alertRepository) GetByType(alertType string) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("type = ?", alertType).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		r.logger.Error("根据类型获取告警失败", utils.ErrorField(err), utils.String("type", alertType))
		return nil, err
	}
	return alerts, nil
}

// GetUnresolved 获取未解决的告警
func (r *alertRepository) GetUnresolved() ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("status = ?", "active").Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		r.logger.Error("获取未解决告警失败", utils.ErrorField(err))
		return nil, err
	}
	return alerts, nil
}

// MarkAsResolved 标记告警为已解决
func (r *alertRepository) MarkAsResolved(alertID uint) error {
	err := r.db.Model(&models.Alert{}).Where("id = ?", alertID).Updates(map[string]interface{}{
		"status":      "resolved",
		"resolved_at": "NOW()",
	}).Error
	if err != nil {
		r.logger.Error("标记告警为已解决失败", utils.ErrorField(err), utils.Int("alert_id", int(alertID)))
		return err
	}
	return nil
}

// GetByTimeRange 根据时间范围获取告警
func (r *alertRepository) GetByTimeRange(startTime, endTime int64) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("created_at BETWEEN ? AND ?", startTime, endTime).Order("created_at DESC").Find(&alerts).Error
	if err != nil {
		r.logger.Error("根据时间范围获取告警失败", utils.ErrorField(err))
		return nil, err
	}
	return alerts, nil
}
