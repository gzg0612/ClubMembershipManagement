package tests

import (
	"huiyuanguanli/backend/models"
	"huiyuanguanli/backend/services"
	"huiyuanguanli/backend/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingsService(t *testing.T) {
	// 初始化数据库连接
	db := utils.InitTestDB()

	// 创建服务实例
	settingsService := services.NewSettingsService(db)

	// 测试获取系统设置
	t.Run("GetSettings", func(t *testing.T) {
		settings, err := settingsService.GetSettings()
		assert.NoError(t, err)
		assert.NotNil(t, settings)
	})

	// 测试更新系统设置
	t.Run("UpdateSettings", func(t *testing.T) {
		settings := &models.SystemSettings{
			SiteName:           "测试系统",
			SiteDescription:    "测试描述",
			DefaultMemberLevel: 2,
			PointsPerYuan:      2,
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

		err := settingsService.UpdateSettings(settings)
		assert.NoError(t, err)

		// 验证更新是否成功
		updated, err := settingsService.GetSettings()
		assert.NoError(t, err)
		assert.Equal(t, settings.SiteName, updated.SiteName)
		assert.Equal(t, settings.SiteDescription, updated.SiteDescription)
		assert.Equal(t, settings.DefaultMemberLevel, updated.DefaultMemberLevel)
		assert.Equal(t, settings.PointsPerYuan, updated.PointsPerYuan)
	})

	// 测试重置系统设置
	t.Run("ResetSettings", func(t *testing.T) {
		err := settingsService.ResetSettings()
		assert.NoError(t, err)

		// 验证重置是否成功
		settings, err := settingsService.GetSettings()
		assert.NoError(t, err)
		assert.Equal(t, "会员管理系统", settings.SiteName)
		assert.Equal(t, "专业的会员管理系统", settings.SiteDescription)
		assert.Equal(t, 1, settings.DefaultMemberLevel)
		assert.Equal(t, 1, settings.PointsPerYuan)
	})

	// 测试验证规则
	t.Run("ValidationRules", func(t *testing.T) {
		// 测试空设置
		err := settingsService.UpdateSettings(nil)
		assert.Error(t, err)
		assert.Equal(t, "设置不能为空", err.Error())

		// 测试空系统名称
		settings := &models.SystemSettings{}
		err = settingsService.UpdateSettings(settings)
		assert.Error(t, err)
		assert.Equal(t, "系统名称不能为空", err.Error())

		// 测试无效的默认会员等级
		settings.SiteName = "测试系统"
		settings.DefaultMemberLevel = 11
		err = settingsService.UpdateSettings(settings)
		assert.Error(t, err)
		assert.Equal(t, "默认会员等级必须在1-10之间", err.Error())

		// 测试无效的积分比例
		settings.DefaultMemberLevel = 1
		settings.PointsPerYuan = 101
		err = settingsService.UpdateSettings(settings)
		assert.Error(t, err)
		assert.Equal(t, "积分比例必须在0-100之间", err.Error())

		// 测试无效的密码最小长度
		settings.PointsPerYuan = 1
		settings.PasswordMinLength = 5
		err = settingsService.UpdateSettings(settings)
		assert.Error(t, err)
		assert.Equal(t, "密码最小长度必须在6-20之间", err.Error())
	})

	// 测试邮件和短信设置
	t.Run("NotificationSettings", func(t *testing.T) {
		// 测试空邮箱
		err := settingsService.TestEmailSettings("")
		assert.Error(t, err)
		assert.Equal(t, "邮箱地址不能为空", err.Error())

		// 测试空手机号
		err = settingsService.TestSMSSettings("")
		assert.Error(t, err)
		assert.Equal(t, "手机号不能为空", err.Error())
	})
}
