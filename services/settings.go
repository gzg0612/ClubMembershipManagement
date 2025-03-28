package services

import (
	"errors"
	"huiyuanguanli/models"
	"huiyuanguanli/utils"
)

// SettingsService 系统设置服务
type SettingsService struct {
	db *utils.DB
}

// NewSettingsService 创建系统设置服务实例
func NewSettingsService(db *utils.DB) *SettingsService {
	return &SettingsService{db: db}
}

// GetSettings 获取系统设置
func (s *SettingsService) GetSettings() (*models.SystemSettings, error) {
	var settings models.SystemSettings
	result := s.db.First(&settings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &settings, nil
}

// UpdateSettings 更新系统设置
func (s *SettingsService) UpdateSettings(settings *models.SystemSettings) error {
	if settings == nil {
		return errors.New("设置不能为空")
	}

	// 验证必填字段
	if settings.SiteName == "" {
		return errors.New("系统名称不能为空")
	}

	// 验证数值范围
	if settings.DefaultMemberLevel < 1 || settings.DefaultMemberLevel > 10 {
		return errors.New("默认会员等级必须在1-10之间")
	}
	if settings.PointsPerYuan < 0 || settings.PointsPerYuan > 100 {
		return errors.New("积分比例必须在0-100之间")
	}
	if settings.PasswordMinLength < 6 || settings.PasswordMinLength > 20 {
		return errors.New("密码最小长度必须在6-20之间")
	}
	if settings.LoginAttempts < 1 || settings.LoginAttempts > 10 {
		return errors.New("登录尝试次数必须在1-10之间")
	}
	if settings.LockoutDuration < 1 || settings.LockoutDuration > 1440 {
		return errors.New("锁定时间必须在1-1440分钟之间")
	}
	if settings.BackupRetention < 1 || settings.BackupRetention > 365 {
		return errors.New("备份保留天数必须在1-365天之间")
	}

	// 更新设置
	result := s.db.Save(settings)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ResetSettings 重置系统设置为默认值
func (s *SettingsService) ResetSettings() error {
	defaultSettings := &models.SystemSettings{
		SiteName:           "会员管理系统",
		SiteDescription:    "专业的会员管理系统",
		DefaultMemberLevel: 1,
		PointsPerYuan:      1,
		EnablePoints:       true,
		PasswordMinLength:  8,
		RequireSpecialChar: true,
		RequireNumber:      true,
		RequireUppercase:   true,
		LoginAttempts:      5,
		LockoutDuration:    30,
		BackupFrequency:    "daily",
		BackupTime:         "00:00",
		BackupRetention:    30,
		Timezone:           "Asia/Shanghai",
		DateFormat:         "YYYY-MM-DD",
		TimeFormat:         "24",
		Currency:           "CNY",
		Language:           "zh-CN",
	}

	return s.UpdateSettings(defaultSettings)
}

// TestEmailSettings 测试邮件设置
func (s *SettingsService) TestEmailSettings(email string) error {
	if email == "" {
		return errors.New("邮箱地址不能为空")
	}

	// TODO: 实现邮件发送测试
	return nil
}

// TestSMSSettings 测试短信设置
func (s *SettingsService) TestSMSSettings(phone string) error {
	if phone == "" {
		return errors.New("手机号不能为空")
	}

	// TODO: 实现短信发送测试
	return nil
}
