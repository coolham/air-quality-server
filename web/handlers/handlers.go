package handlers

import (
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
)

// WebHandlers Web处理器集合
type WebHandlers struct {
	services *services.Services
	logger   utils.Logger
}

// NewWebHandlers 创建Web处理器
func NewWebHandlers(services *services.Services, logger utils.Logger) *WebHandlers {
	return &WebHandlers{
		services: services,
		logger:   logger,
	}
}
