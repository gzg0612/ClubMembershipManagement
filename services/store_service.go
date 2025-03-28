package services

import (
	"errors"
	"time"

	"huiyuanguanli/models"

	"gorm.io/gorm"
)

// StoreService 店铺服务
type StoreService struct {
	db *gorm.DB
}

// NewStoreService 创建店铺服务实例
func NewStoreService(db *gorm.DB) *StoreService {
	return &StoreService{db: db}
}

// CreateStore 创建店铺
func (s *StoreService) CreateStore(store *models.Store) error {
	// 检查店铺编码是否已存在
	var count int64
	s.db.Model(&models.Store{}).Where("code = ?", store.Code).Count(&count)
	if count > 0 {
		return errors.New("店铺编码已存在")
	}

	// 创建店铺
	return s.db.Create(store).Error
}

// UpdateStore 更新店铺
func (s *StoreService) UpdateStore(store *models.Store) error {
	// 检查店铺是否存在
	var existingStore models.Store
	if err := s.db.First(&existingStore, store.ID).Error; err != nil {
		return err
	}

	// 如果修改了店铺编码，检查新编码是否已存在
	if existingStore.Code != store.Code {
		var count int64
		s.db.Model(&models.Store{}).Where("code = ?", store.Code).Count(&count)
		if count > 0 {
			return errors.New("店铺编码已存在")
		}
	}

	// 更新店铺
	return s.db.Save(store).Error
}

// DeleteStore 删除店铺
func (s *StoreService) DeleteStore(id uint) error {
	// 检查是否有员工
	var staffCount int64
	s.db.Model(&models.StoreStaff{}).Where("store_id = ?", id).Count(&staffCount)
	if staffCount > 0 {
		return errors.New("该店铺还有员工，无法删除")
	}

	// 检查是否有设备
	var deviceCount int64
	s.db.Model(&models.StoreDevice{}).Where("store_id = ?", id).Count(&deviceCount)
	if deviceCount > 0 {
		return errors.New("该店铺还有设备，无法删除")
	}

	return s.db.Delete(&models.Store{}, id).Error
}

// GetStore 获取店铺详情
func (s *StoreService) GetStore(id uint) (*models.Store, error) {
	var store models.Store
	err := s.db.First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

// GetStoreList 获取店铺列表
func (s *StoreService) GetStoreList(page, pageSize int, keyword string, status int) ([]models.Store, int64, error) {
	var stores []models.Store
	var total int64

	query := s.db.Model(&models.Store{})

	// 添加查询条件
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&stores).Error
	if err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

// AddStoreStaff 添加店铺员工
func (s *StoreService) AddStoreStaff(staff *models.StoreStaff) error {
	// 检查店铺是否存在
	var store models.Store
	if err := s.db.First(&store, staff.StoreID).Error; err != nil {
		return err
	}

	// 检查用户是否已在其他店铺任职
	var existingStaff models.StoreStaff
	err := s.db.Where("user_id = ? AND status = 1", staff.UserID).First(&existingStaff).Error
	if err == nil {
		return errors.New("该用户已在其他店铺任职")
	}

	// 设置入职日期
	staff.StartDate = time.Now()

	// 创建员工记录
	return s.db.Create(staff).Error
}

// UpdateStoreStaff 更新店铺员工信息
func (s *StoreService) UpdateStoreStaff(staff *models.StoreStaff) error {
	return s.db.Save(staff).Error
}

// RemoveStoreStaff 移除店铺员工
func (s *StoreService) RemoveStoreStaff(id uint) error {
	// 更新员工状态为离职
	return s.db.Model(&models.StoreStaff{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   0,
		"end_date": time.Now(),
	}).Error
}

// GetStoreStaffList 获取店铺员工列表
func (s *StoreService) GetStoreStaffList(storeID uint, page, pageSize int) ([]models.StoreStaff, int64, error) {
	var staffs []models.StoreStaff
	var total int64

	query := s.db.Model(&models.StoreStaff{}).Where("store_id = ?", storeID)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&staffs).Error
	if err != nil {
		return nil, 0, err
	}

	return staffs, total, nil
}

// AddStoreDevice 添加店铺设备
func (s *StoreService) AddStoreDevice(device *models.StoreDevice) error {
	// 检查店铺是否存在
	var store models.Store
	if err := s.db.First(&store, device.StoreID).Error; err != nil {
		return err
	}

	// 检查设备编号是否已存在
	var count int64
	s.db.Model(&models.StoreDevice{}).Where("code = ?", device.Code).Count(&count)
	if count > 0 {
		return errors.New("设备编号已存在")
	}

	// 设置检查时间
	now := time.Now()
	device.LastCheck = now
	device.NextCheck = now.AddDate(0, 1, 0) // 默认一个月后检查

	// 创建设备记录
	return s.db.Create(device).Error
}

// UpdateStoreDevice 更新店铺设备信息
func (s *StoreService) UpdateStoreDevice(device *models.StoreDevice) error {
	return s.db.Save(device).Error
}

// DeleteStoreDevice 删除店铺设备
func (s *StoreService) DeleteStoreDevice(id uint) error {
	return s.db.Delete(&models.StoreDevice{}, id).Error
}

// GetStoreDeviceList 获取店铺设备列表
func (s *StoreService) GetStoreDeviceList(storeID uint, page, pageSize int) ([]models.StoreDevice, int64, error) {
	var devices []models.StoreDevice
	var total int64

	query := s.db.Model(&models.StoreDevice{}).Where("store_id = ?", storeID)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&devices).Error
	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

// AddDeviceCheck 添加设备检查记录
func (s *StoreService) AddDeviceCheck(check *models.StoreDeviceCheck) error {
	// 开启事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 创建检查记录
	if err := tx.Create(check).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新设备检查时间
	if err := tx.Model(&models.StoreDevice{}).Where("id = ?", check.DeviceID).Updates(map[string]interface{}{
		"last_check": check.CreatedAt,
		"next_check": check.CreatedAt.AddDate(0, 1, 0), // 默认一个月后检查
		"status":     check.Status,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetDeviceCheckList 获取设备检查记录列表
func (s *StoreService) GetDeviceCheckList(deviceID uint, page, pageSize int) ([]models.StoreDeviceCheck, int64, error) {
	var checks []models.StoreDeviceCheck
	var total int64

	query := s.db.Model(&models.StoreDeviceCheck{}).Where("device_id = ?", deviceID)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&checks).Error
	if err != nil {
		return nil, 0, err
	}

	return checks, total, nil
}
