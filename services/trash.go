package services

import (
	"encoding/json"
	"huiyuanguanli/models"
	"time"

	"gorm.io/gorm"
)

// TrashService 数据回收站服务
type TrashService struct {
	db *gorm.DB
}

// NewTrashService 创建数据回收站服务实例
func NewTrashService(db *gorm.DB) *TrashService {
	return &TrashService{db: db}
}

// CreateDeletedRecord 创建删除记录
func (s *TrashService) CreateDeletedRecord(recordType models.RecordType, recordID uint, recordData interface{}, deleteReason string, operatorID uint, operatorName string) error {
	jsonData, err := json.Marshal(recordData)
	if err != nil {
		return err
	}

	record := &models.DeletedRecord{
		RecordType:   recordType,
		RecordID:     recordID,
		RecordData:   string(jsonData),
		DeleteReason: deleteReason,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		DeletedAt:    time.Now(),
	}

	if err := record.Validate(); err != nil {
		return err
	}

	return s.db.Create(record).Error
}

// GetDeletedRecords 获取删除记录列表
func (s *TrashService) GetDeletedRecords(recordType models.RecordType, page, pageSize int) ([]models.DeletedRecord, int64, error) {
	var records []models.DeletedRecord
	var total int64

	query := s.db.Model(&models.DeletedRecord{})
	if recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("deleted_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

// GetDeletedRecord 获取单个删除记录
func (s *TrashService) GetDeletedRecord(id uint) (*models.DeletedRecord, error) {
	var record models.DeletedRecord
	err := s.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// RestoreRecord 恢复删除的记录
func (s *TrashService) RestoreRecord(id uint) error {
	record, err := s.GetDeletedRecord(id)
	if err != nil {
		return err
	}

	// 根据记录类型恢复数据
	switch record.RecordType {
	case models.RecordTypeMember:
		var member models.Member
		if err := json.Unmarshal([]byte(record.RecordData), &member); err != nil {
			return err
		}
		return s.db.Create(&member).Error

	case models.RecordTypeProduct:
		var product models.Product
		if err := json.Unmarshal([]byte(record.RecordData), &product); err != nil {
			return err
		}
		return s.db.Create(&product).Error

	case models.RecordTypeActivity:
		var activity models.Activity
		if err := json.Unmarshal([]byte(record.RecordData), &activity); err != nil {
			return err
		}
		return s.db.Create(&activity).Error

	default:
		return models.ErrInvalidRecordType
	}
}

// PermanentlyDelete 永久删除记录
func (s *TrashService) PermanentlyDelete(id uint) error {
	return s.db.Unscoped().Delete(&models.DeletedRecord{}, id).Error
}
