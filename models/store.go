package models

import (
	"time"
)

// Store 店铺模型
type Store struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:100;not null"`   // 店铺名称
	Code          string    `json:"code" gorm:"size:50;uniqueIndex"` // 店铺编码
	Address       string    `json:"address" gorm:"size:255"`         // 店铺地址
	Phone         string    `json:"phone" gorm:"size:20"`            // 联系电话
	Manager       string    `json:"manager" gorm:"size:50"`          // 负责人
	Email         string    `json:"email" gorm:"size:100"`           // 电子邮箱
	BusinessHours string    `json:"business_hours" gorm:"size:100"`  // 营业时间
	Status        int       `json:"status" gorm:"default:1"`         // 状态：1营业中，0休息中
	Description   string    `json:"description" gorm:"type:text"`    // 店铺描述
	Logo          string    `json:"logo" gorm:"size:255"`            // 店铺logo
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"index"` // 软删除
}

// StoreStaff 店铺员工模型
type StoreStaff struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoreID   uint      `json:"store_id" gorm:"not null"`     // 店铺ID
	UserID    uint      `json:"user_id" gorm:"not null"`      // 用户ID
	Role      string    `json:"role" gorm:"size:50;not null"` // 角色：manager-店长，staff-店员
	Status    int       `json:"status" gorm:"default:1"`      // 状态：1在职，0离职
	StartDate time.Time `json:"start_date"`                   // 入职日期
	EndDate   time.Time `json:"end_date"`                     // 离职日期
	Remark    string    `json:"remark" gorm:"size:255"`       // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"` // 软删除
}

// StoreDevice 店铺设备模型
type StoreDevice struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoreID   uint      `json:"store_id" gorm:"not null"`        // 店铺ID
	Name      string    `json:"name" gorm:"size:100;not null"`   // 设备名称
	Type      string    `json:"type" gorm:"size:50;not null"`    // 设备类型
	Code      string    `json:"code" gorm:"size:50;uniqueIndex"` // 设备编号
	Status    int       `json:"status" gorm:"default:1"`         // 状态：1正常，0故障
	LastCheck time.Time `json:"last_check"`                      // 最后检查时间
	NextCheck time.Time `json:"next_check"`                      // 下次检查时间
	Remark    string    `json:"remark" gorm:"size:255"`          // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"` // 软删除
}

// StoreDeviceCheck 设备检查记录模型
type StoreDeviceCheck struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	DeviceID  uint      `json:"device_id" gorm:"not null"`          // 设备ID
	StoreID   uint      `json:"store_id" gorm:"not null"`           // 店铺ID
	CheckType string    `json:"check_type" gorm:"size:50;not null"` // 检查类型：daily-日常检查，weekly-周检查，monthly-月检查
	Status    int       `json:"status" gorm:"default:1"`            // 状态：1正常，0异常
	Content   string    `json:"content" gorm:"type:text"`           // 检查内容
	Result    string    `json:"result" gorm:"type:text"`            // 检查结果
	Operator  string    `json:"operator" gorm:"size:50"`            // 操作人
	CreatedAt time.Time `json:"created_at"`
}
