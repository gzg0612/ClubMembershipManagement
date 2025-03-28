package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Username         string         `json:"username" gorm:"type:varchar(50);not null;uniqueIndex"`
	Password         string         `json:"-" gorm:"type:varchar(255);not null"`
	RealName         string         `json:"real_name" gorm:"type:varchar(50)"`
	Phone            string         `json:"phone" gorm:"type:varchar(20)"`
	Email            string         `json:"email" gorm:"type:varchar(100)"`
	Avatar           string         `json:"avatar" gorm:"type:varchar(255)"`
	Gender           int8           `json:"gender" gorm:"type:tinyint;default:0"` // 0未知,1男,2女
	Status           int8           `json:"status" gorm:"type:tinyint;default:1"` // 1正常,0禁用,2锁定
	DepartmentID     *uint          `json:"department_id" gorm:""`
	Position         string         `json:"position" gorm:"type:varchar(50)"`
	LastLoginTime    *time.Time     `json:"last_login_time"`
	LastLoginIP      string         `json:"last_login_ip" gorm:"type:varchar(50)"`
	LoginFailedCount int            `json:"login_failed_count" gorm:"default:0"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(50);not null"`
	Code        string         `json:"code" gorm:"type:varchar(50);not null;uniqueIndex"`
	Description string         `json:"description" gorm:"type:varchar(200)"`
	Status      int8           `json:"status" gorm:"type:tinyint;default:1"` // 1启用,0禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
}

// Permission 权限模型
type Permission struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Code      string         `json:"code" gorm:"type:varchar(50);not null;uniqueIndex"`
	Type      string         `json:"type" gorm:"type:varchar(20);not null"` // menu菜单,button按钮,api接口
	ParentID  *uint          `json:"parent_id"`
	Path      string         `json:"path" gorm:"type:varchar(200)"`
	Component string         `json:"component" gorm:"type:varchar(200)"`
	Icon      string         `json:"icon" gorm:"type:varchar(50)"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int8           `json:"status" gorm:"type:tinyint;default:1"` // 1启用,0禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Roles []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	RoleID    uint      `json:"role_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoleID       uint      `json:"role_id" gorm:"not null"`
	PermissionID uint      `json:"permission_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      *uint     `json:"user_id"`
	Username    string    `json:"username" gorm:"type:varchar(50);not null"`
	LoginType   string    `json:"login_type" gorm:"type:varchar(20)"`
	LoginStatus string    `json:"login_status" gorm:"type:varchar(20);not null"`
	IP          string    `json:"ip" gorm:"type:varchar(50);not null"`
	Location    string    `json:"location" gorm:"type:varchar(100)"`
	Device      string    `json:"device" gorm:"type:varchar(200)"`
	Browser     string    `json:"browser" gorm:"type:varchar(200)"`
	OS          string    `json:"os" gorm:"type:varchar(200)"`
	UserAgent   string    `json:"user_agent" gorm:"type:text"`
	LoginTime   time.Time `json:"login_time" gorm:"default:CURRENT_TIMESTAMP"`
}
