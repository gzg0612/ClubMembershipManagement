package models

import (
	"time"

	"gorm.io/gorm"
)

// RecordType 记录类型
type RecordType string

const (
	RecordTypeMember   RecordType = "member"   // 会员
	RecordTypeProduct  RecordType = "product"  // 商品
	RecordTypeActivity RecordType = "activity" // 活动
)

// DeletedRecord 已删除记录
type DeletedRecord struct {
	gorm.Model
	RecordType   RecordType `json:"recordType" gorm:"size:20;not null"`   // 记录类型
	RecordID     uint       `json:"recordId" gorm:"not null"`             // 记录ID
	RecordData   string     `json:"recordData" gorm:"type:text;not null"` // 删除前数据的JSON备份
	DeleteReason string     `json:"deleteReason" gorm:"type:text"`        // 删除原因
	OperatorID   uint       `json:"operatorId" gorm:"not null"`           // 删除操作人ID
	OperatorName string     `json:"operatorName" gorm:"size:50;not null"` // 删除操作人姓名
	DeletedAt    time.Time  `json:"deletedAt" gorm:"not null"`            // 删除时间
}

// Validate 验证删除记录数据
func (d *DeletedRecord) Validate() error {
	if d.RecordType == "" {
		return ErrInvalidRecordType
	}
	if d.RecordID == 0 {
		return ErrInvalidRecordID
	}
	if d.RecordData == "" {
		return ErrInvalidRecordData
	}
	if d.OperatorID == 0 {
		return ErrInvalidOperatorID
	}
	if d.OperatorName == "" {
		return ErrInvalidOperatorName
	}
	return nil
}
