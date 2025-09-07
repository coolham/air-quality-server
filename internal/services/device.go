package services

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/utils"
	"context"
)

// DeviceService 设备服务接口
type DeviceService interface {
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDevice(ctx context.Context, id string) (*models.Device, error)
	GetDeviceBySerialNumber(ctx context.Context, serialNumber string) (*models.Device, error)
	UpdateDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, id string) error
	ListDevices(ctx context.Context, limit, offset int) ([]models.Device, error)
	CountDevices(ctx context.Context) (int64, error)
	GetDeviceStatus(ctx context.Context, id string) (*models.Device, error)
	UpdateDeviceStatus(ctx context.Context, id string, status string) error
}

// deviceService 设备服务实现
type deviceService struct {
	deviceRepo repositories.DeviceRepository
	logger     utils.Logger
}

// NewDeviceService 创建设备服务
func NewDeviceService(deviceRepo repositories.DeviceRepository, logger utils.Logger) DeviceService {
	return &deviceService{
		deviceRepo: deviceRepo,
		logger:     logger,
	}
}

// CreateDevice 创建设备
func (s *deviceService) CreateDevice(ctx context.Context, device *models.Device) error {
	if err := s.deviceRepo.Create(ctx, device); err != nil {
		s.logger.Error("创建设备失败", utils.ErrorField(err))
		return err
	}
	s.logger.Info("设备创建成功", utils.String("device_id", device.ID))
	return nil
}

// GetDevice 获取设备
func (s *deviceService) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取设备失败", utils.ErrorField(err), utils.String("device_id", id))
		return nil, err
	}
	return device, nil
}

// GetDeviceBySerialNumber 根据序列号获取设备
func (s *deviceService) GetDeviceBySerialNumber(ctx context.Context, serialNumber string) (*models.Device, error) {
	device, err := s.deviceRepo.GetByDeviceID(ctx, serialNumber)
	if err != nil {
		s.logger.Error("根据序列号获取设备失败", utils.ErrorField(err), utils.String("serial_number", serialNumber))
		return nil, err
	}
	return device, nil
}

// UpdateDevice 更新设备
func (s *deviceService) UpdateDevice(ctx context.Context, device *models.Device) error {
	// 使用结构体更新，只更新非零值字段
	updateData := &models.Device{
		Name:              device.Name,
		Type:              device.Type,
		LocationLatitude:  device.LocationLatitude,
		LocationLongitude: device.LocationLongitude,
		LocationAddress:   device.LocationAddress,
		Status:            device.Status,
		Config:            device.Config,
		UpdatedAt:         device.UpdatedAt,
	}

	if err := s.deviceRepo.Update(ctx, device.ID, updateData); err != nil {
		s.logger.Error("更新设备失败", utils.ErrorField(err), utils.String("device_id", device.ID))
		return err
	}
	s.logger.Info("设备更新成功", utils.String("device_id", device.ID))
	return nil
}

// DeleteDevice 删除设备
func (s *deviceService) DeleteDevice(ctx context.Context, id string) error {
	if err := s.deviceRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除设备失败", utils.ErrorField(err), utils.String("device_id", id))
		return err
	}
	s.logger.Info("设备删除成功", utils.String("device_id", id))
	return nil
}

// ListDevices 列出设备
func (s *deviceService) ListDevices(ctx context.Context, limit, offset int) ([]models.Device, error) {
	req := &repositories.ListRequest{
		Page:     offset/limit + 1,
		PageSize: limit,
		OrderBy:  "created_at",
		Order:    "desc",
	}

	response, err := s.deviceRepo.List(ctx, req)
	if err != nil {
		s.logger.Error("列出设备失败", utils.ErrorField(err))
		return nil, err
	}
	return response.Data, nil
}

// GetDeviceStatus 获取设备状态
func (s *deviceService) GetDeviceStatus(ctx context.Context, id string) (*models.Device, error) {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取设备状态失败", utils.ErrorField(err), utils.String("device_id", id))
		return nil, err
	}

	return device, nil
}

// UpdateDeviceStatus 更新设备状态
func (s *deviceService) UpdateDeviceStatus(ctx context.Context, id string, status string) error {
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取设备失败", utils.ErrorField(err), utils.String("device_id", id))
		return err
	}

	device.Status = models.DeviceStatus(status)
	updateData := &models.Device{
		Status:    device.Status,
		UpdatedAt: device.UpdatedAt,
	}
	if err := s.deviceRepo.Update(ctx, device.ID, updateData); err != nil {
		s.logger.Error("更新设备状态失败", utils.ErrorField(err), utils.String("device_id", id))
		return err
	}

	s.logger.Info("设备状态更新成功", utils.String("device_id", id), utils.String("status", status))
	return nil
}

// CountDevices 获取设备总数
func (s *deviceService) CountDevices(ctx context.Context) (int64, error) {
	count, err := s.deviceRepo.Count(ctx, map[string]interface{}{})
	if err != nil {
		s.logger.Error("获取设备总数失败", utils.ErrorField(err))
		return 0, err
	}
	return count, nil
}
