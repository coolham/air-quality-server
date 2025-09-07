package models

import (
	"time"

	"gorm.io/gorm"
)

// AlertRule 告警规则模型
type AlertRule struct {
	ID                   uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                 string         `json:"name" gorm:"type:varchar(100);not null"`
	DeviceID             *string        `json:"device_id" gorm:"type:varchar(64);index"`
	Metric               string         `json:"metric" gorm:"type:varchar(50);not null"`
	ConditionType        string         `json:"condition_type" gorm:"type:enum('gt','lt','eq','ne','gte','lte');not null"`
	ThresholdValue       float64        `json:"threshold_value" gorm:"type:decimal(10,2);not null"`
	DurationSeconds      int            `json:"duration_seconds" gorm:"default:0"`
	Severity             string         `json:"severity" gorm:"type:enum('critical','warning','info');default:'warning'"`
	Enabled              bool           `json:"enabled" gorm:"default:true"`
	NotificationChannels *string        `json:"notification_channels" gorm:"type:json"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (AlertRule) TableName() string {
	return "alert_rules"
}

// Alert 告警记录模型
type Alert struct {
	ID             uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	RuleID         uint64         `json:"rule_id" gorm:"not null;index"`
	DeviceID       string         `json:"device_id" gorm:"type:varchar(64);not null;index"`
	Metric         string         `json:"metric" gorm:"type:varchar(50);not null"`
	CurrentValue   float64        `json:"current_value" gorm:"type:decimal(10,2);not null"`
	ThresholdValue float64        `json:"threshold_value" gorm:"type:decimal(10,2);not null"`
	Severity       string         `json:"severity" gorm:"type:enum('critical','warning','info');not null"`
	Status         string         `json:"status" gorm:"type:enum('active','acknowledged','resolved');default:'active'"`
	TriggeredAt    time.Time      `json:"triggered_at" gorm:"not null;index"`
	AcknowledgedAt *time.Time     `json:"acknowledged_at"`
	ResolvedAt     *time.Time     `json:"resolved_at"`
	AcknowledgedBy *uint64        `json:"acknowledged_by"`
	ResolvedBy     *uint64        `json:"resolved_by"`
	Message        *string        `json:"message" gorm:"type:text"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (Alert) TableName() string {
	return "alerts"
}

// AlertConditionType 告警条件类型
type AlertConditionType string

const (
	AlertConditionGT  AlertConditionType = "gt"  // 大于
	AlertConditionLT  AlertConditionType = "lt"  // 小于
	AlertConditionEQ  AlertConditionType = "eq"  // 等于
	AlertConditionNE  AlertConditionType = "ne"  // 不等于
	AlertConditionGTE AlertConditionType = "gte" // 大于等于
	AlertConditionLTE AlertConditionType = "lte" // 小于等于
)

// IsValid 验证告警条件类型
func (t AlertConditionType) IsValid() bool {
	switch t {
	case AlertConditionGT, AlertConditionLT, AlertConditionEQ, AlertConditionNE, AlertConditionGTE, AlertConditionLTE:
		return true
	default:
		return false
	}
}

// AlertSeverity 告警严重程度
type AlertSeverity string

const (
	AlertSeverityCritical AlertSeverity = "critical"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityInfo     AlertSeverity = "info"
)

// IsValid 验证告警严重程度
func (s AlertSeverity) IsValid() bool {
	switch s {
	case AlertSeverityCritical, AlertSeverityWarning, AlertSeverityInfo:
		return true
	default:
		return false
	}
}

// AlertStatus 告警状态
type AlertStatus string

const (
	AlertStatusActive       AlertStatus = "active"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
)

// IsValid 验证告警状态
func (s AlertStatus) IsValid() bool {
	switch s {
	case AlertStatusActive, AlertStatusAcknowledged, AlertStatusResolved:
		return true
	default:
		return false
	}
}

// AlertRuleCreateRequest 创建告警规则请求
type AlertRuleCreateRequest struct {
	Name                 string   `json:"name" binding:"required"`
	DeviceID             *string  `json:"device_id,omitempty"`
	Metric               string   `json:"metric" binding:"required"`
	ConditionType        string   `json:"condition_type" binding:"required"`
	ThresholdValue       float64  `json:"threshold_value" binding:"required"`
	DurationSeconds      int      `json:"duration_seconds"`
	Severity             string   `json:"severity"`
	Enabled              bool     `json:"enabled"`
	NotificationChannels []string `json:"notification_channels,omitempty"`
}

// AlertRuleUpdateRequest 更新告警规则请求
type AlertRuleUpdateRequest struct {
	Name                 *string  `json:"name,omitempty"`
	DeviceID             *string  `json:"device_id,omitempty"`
	Metric               *string  `json:"metric,omitempty"`
	ConditionType        *string  `json:"condition_type,omitempty"`
	ThresholdValue       *float64 `json:"threshold_value,omitempty"`
	DurationSeconds      *int     `json:"duration_seconds,omitempty"`
	Severity             *string  `json:"severity,omitempty"`
	Enabled              *bool    `json:"enabled,omitempty"`
	NotificationChannels []string `json:"notification_channels,omitempty"`
}

// AlertRuleListRequest 告警规则列表请求
type AlertRuleListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	DeviceID string `form:"device_id"`
	Metric   string `form:"metric"`
	Severity string `form:"severity"`
	Enabled  *bool  `form:"enabled"`
}

// AlertRuleListResponse 告警规则列表响应
type AlertRuleListResponse struct {
	Rules    []AlertRule `json:"rules"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// AlertListRequest 告警列表请求
type AlertListRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	PageSize  int    `form:"page_size" binding:"min=1,max=100"`
	DeviceID  string `form:"device_id"`
	Metric    string `form:"metric"`
	Severity  string `form:"severity"`
	Status    string `form:"status"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
}

// AlertListResponse 告警列表响应
type AlertListResponse struct {
	Alerts   []Alert `json:"alerts"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

// AlertAcknowledgeRequest 确认告警请求
type AlertAcknowledgeRequest struct {
	Message string `json:"message,omitempty"`
}

// AlertResolveRequest 解决告警请求
type AlertResolveRequest struct {
	Message string `json:"message,omitempty"`
}

// AlertStatistics 告警统计
type AlertStatistics struct {
	TotalAlerts        int64 `json:"total_alerts"`
	ActiveAlerts       int64 `json:"active_alerts"`
	AcknowledgedAlerts int64 `json:"acknowledged_alerts"`
	ResolvedAlerts     int64 `json:"resolved_alerts"`
	CriticalAlerts     int64 `json:"critical_alerts"`
	WarningAlerts      int64 `json:"warning_alerts"`
	InfoAlerts         int64 `json:"info_alerts"`
}

// AlertTrend 告警趋势
type AlertTrend struct {
	Date          time.Time `json:"date"`
	AlertCount    int64     `json:"alert_count"`
	CriticalCount int64     `json:"critical_count"`
	WarningCount  int64     `json:"warning_count"`
	InfoCount     int64     `json:"info_count"`
}

// NotificationChannel 通知渠道
type NotificationChannel string

const (
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelSMS      NotificationChannel = "sms"
	NotificationChannelWebhook  NotificationChannel = "webhook"
	NotificationChannelDingTalk NotificationChannel = "dingtalk"
	NotificationChannelWeChat   NotificationChannel = "wechat"
)

// IsValid 验证通知渠道
func (c NotificationChannel) IsValid() bool {
	switch c {
	case NotificationChannelEmail, NotificationChannelSMS, NotificationChannelWebhook, NotificationChannelDingTalk, NotificationChannelWeChat:
		return true
	default:
		return false
	}
}

// DefaultAlertRules 默认告警规则
var DefaultAlertRules = []AlertRule{
	{
		Name:                 "PM2.5超标告警",
		Metric:               "pm25",
		ConditionType:        string(AlertConditionGT),
		ThresholdValue:       75.0,
		DurationSeconds:      300,
		Severity:             string(AlertSeverityWarning),
		Enabled:              true,
		NotificationChannels: &[]string{`["email","sms"]`}[0],
	},
	{
		Name:                 "PM10超标告警",
		Metric:               "pm10",
		ConditionType:        string(AlertConditionGT),
		ThresholdValue:       150.0,
		DurationSeconds:      300,
		Severity:             string(AlertSeverityWarning),
		Enabled:              true,
		NotificationChannels: &[]string{`["email","sms"]`}[0],
	},
	{
		Name:                 "CO2浓度告警",
		Metric:               "co2",
		ConditionType:        string(AlertConditionGT),
		ThresholdValue:       1000.0,
		DurationSeconds:      600,
		Severity:             string(AlertSeverityCritical),
		Enabled:              true,
		NotificationChannels: &[]string{`["email","sms","webhook"]`}[0],
	},
	{
		Name:                 "温度异常告警",
		Metric:               "temperature",
		ConditionType:        string(AlertConditionGT),
		ThresholdValue:       40.0,
		DurationSeconds:      180,
		Severity:             string(AlertSeverityWarning),
		Enabled:              true,
		NotificationChannels: &[]string{`["email"]`}[0],
	},
	{
		Name:                 "湿度异常告警",
		Metric:               "humidity",
		ConditionType:        string(AlertConditionLT),
		ThresholdValue:       20.0,
		DurationSeconds:      300,
		Severity:             string(AlertSeverityInfo),
		Enabled:              true,
		NotificationChannels: &[]string{`["email"]`}[0],
	},
}
