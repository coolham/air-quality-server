package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string         `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Phone        *string        `json:"phone" gorm:"type:varchar(20)"`
	PasswordHash string         `json:"-" gorm:"type:varchar(255);not null"`
	Status       string         `json:"status" gorm:"type:varchar(20);default:'active'"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色模型
type Role struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	Description *string        `json:"description" gorm:"type:varchar(200)"`
	Permissions *string        `json:"permissions" gorm:"type:json"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `json:"user_id" gorm:"not null;uniqueIndex:uk_user_role"`
	RoleID    uint64    `json:"role_id" gorm:"not null;uniqueIndex:uk_user_role"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// IsValid 验证用户状态
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	default:
		return false
	}
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Phone    string   `json:"phone,omitempty"`
	Password string   `json:"password" binding:"required,min=6"`
	RoleIDs  []uint64 `json:"role_ids,omitempty"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Username *string  `json:"username,omitempty"`
	Email    *string  `json:"email,omitempty"`
	Phone    *string  `json:"phone,omitempty"`
	Status   *string  `json:"status,omitempty"`
	RoleIDs  []uint64 `json:"role_ids,omitempty"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token     string    `json:"token"`
	User      UserInfo  `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint64  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Status   string  `json:"status"`
	Roles    []Role  `json:"roles"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users    []UserInfo `json:"users"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=50"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Roles    []Role `json:"roles"`
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// Permission 权限常量
const (
	// 设备权限
	PermissionDeviceRead   = "device:read"
	PermissionDeviceWrite  = "device:write"
	PermissionDeviceDelete = "device:delete"

	// 数据权限
	PermissionDataRead   = "data:read"
	PermissionDataWrite  = "data:write"
	PermissionDataExport = "data:export"

	// 告警权限
	PermissionAlertRead    = "alert:read"
	PermissionAlertWrite   = "alert:write"
	PermissionAlertAck     = "alert:ack"
	PermissionAlertResolve = "alert:resolve"

	// 用户权限
	PermissionUserRead   = "user:read"
	PermissionUserWrite  = "user:write"
	PermissionUserDelete = "user:delete"

	// 角色权限
	PermissionRoleRead   = "role:read"
	PermissionRoleWrite  = "role:write"
	PermissionRoleDelete = "role:delete"

	// 系统权限
	PermissionSystemRead   = "system:read"
	PermissionSystemWrite  = "system:write"
	PermissionSystemConfig = "system:config"

	// 管理员权限
	PermissionAdmin = "*"
)

// DefaultRoles 默认角色
var DefaultRoles = []Role{
	{
		Name:        "admin",
		Description: &[]string{"系统管理员"}[0],
		Permissions: &[]string{`["*"]`}[0],
	},
	{
		Name:        "operator",
		Description: &[]string{"操作员"}[0],
		Permissions: stringPtr(`["device:read","device:write","data:read","alert:read","alert:write"]`),
	},
	{
		Name:        "viewer",
		Description: &[]string{"查看者"}[0],
		Permissions: stringPtr(`["device:read","data:read","alert:read"]`),
	},
}

// stringPtr 字符串指针辅助函数
func stringPtr(s string) *string {
	return &s
}
