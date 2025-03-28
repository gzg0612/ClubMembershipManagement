package services

import (
	"huiyuanguanli/models"

	"gorm.io/gorm"
)

// WechatMediaService 微信媒体服务
type WechatMediaService struct {
	db *gorm.DB
}

// NewWechatMediaService 创建微信媒体服务实例
func NewWechatMediaService(db *gorm.DB) *WechatMediaService {
	return &WechatMediaService{db: db}
}

// UploadMedia 上传媒体文件
func (s *WechatMediaService) UploadMedia(fileType string, filePath string) (*models.WechatMedia, error) {
	media := &models.WechatMedia{
		Type:     fileType,
		FilePath: filePath,
		Status:   "pending",
	}

	if err := s.db.Create(media).Error; err != nil {
		return nil, err
	}

	// TODO: 调用微信API上传媒体文件
	// 这里需要实现具体的微信API调用逻辑

	return media, nil
}

// GetMedia 获取媒体文件信息
func (s *WechatMediaService) GetMedia(id uint) (*models.WechatMedia, error) {
	var media models.WechatMedia
	if err := s.db.First(&media, id).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

// DeleteMedia 删除媒体文件
func (s *WechatMediaService) DeleteMedia(id uint) error {
	media, err := s.GetMedia(id)
	if err != nil {
		return err
	}

	// TODO: 调用微信API删除媒体文件
	// 这里需要实现具体的微信API调用逻辑

	return s.db.Delete(media).Error
}

// UpdateMediaStatus 更新媒体文件状态
func (s *WechatMediaService) UpdateMediaStatus(id uint, status string) error {
	return s.db.Model(&models.WechatMedia{}).Where("id = ?", id).Update("status", status).Error
}
