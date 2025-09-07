package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	KeyName     string         `json:"key_name" gorm:"type:varchar(100);uniqueIndex;not null"`
	Value       *string        `json:"value" gorm:"type:text"`
	Description *string        `json:"description" gorm:"type:varchar(200)"`
	Category    *string        `json:"category" gorm:"type:varchar(50)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// SystemConfigCreateRequest 创建系统配置请求
type SystemConfigCreateRequest struct {
	KeyName     string `json:"key_name" binding:"required"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// SystemConfigUpdateRequest 更新系统配置请求
type SystemConfigUpdateRequest struct {
	Value       *string `json:"value,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

// SystemConfigListRequest 系统配置列表请求
type SystemConfigListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
}

// SystemConfigListResponse 系统配置列表响应
type SystemConfigListResponse struct {
	Configs  []SystemConfig `json:"configs"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}
