package models

import (
	"time"
)

// SystemSettings 系统设置模型
type SystemSettings struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 基本设置
	SiteName        string `gorm:"size:100" json:"site_name"`        // 系统名称
	SiteLogo        string `gorm:"size:255" json:"site_logo"`        // 系统Logo
	SiteDescription string `gorm:"size:500" json:"site_description"` // 系统描述
	ContactPhone    string `gorm:"size:20" json:"contact_phone"`     // 联系电话
	ContactEmail    string `gorm:"size:100" json:"contact_email"`    // 联系邮箱
	Address         string `gorm:"size:255" json:"address"`          // 地址

	// 会员设置
	DefaultMemberLevel int  `gorm:"default:1" json:"default_member_level"`    // 默认会员等级
	EnableAutoUpgrade  bool `gorm:"default:false" json:"enable_auto_upgrade"` // 是否启用自动升级
	PointsPerYuan      int  `gorm:"default:1" json:"points_per_yuan"`         // 每消费1元获得积分
	EnablePoints       bool `gorm:"default:true" json:"enable_points"`        // 是否启用积分系统

	// 通知设置
	EnableSMS    bool `gorm:"default:false" json:"enable_sms"`    // 是否启用短信通知
	EnableEmail  bool `gorm:"default:false" json:"enable_email"`  // 是否启用邮件通知
	EnableWechat bool `gorm:"default:false" json:"enable_wechat"` // 是否启用微信通知

	// 安全设置
	PasswordMinLength  int  `gorm:"default:8" json:"password_min_length"`     // 密码最小长度
	RequireSpecialChar bool `gorm:"default:true" json:"require_special_char"` // 是否要求特殊字符
	RequireNumber      bool `gorm:"default:true" json:"require_number"`       // 是否要求数字
	RequireUppercase   bool `gorm:"default:true" json:"require_uppercase"`    // 是否要求大写字母
	LoginAttempts      int  `gorm:"default:5" json:"login_attempts"`          // 最大登录尝试次数
	LockoutDuration    int  `gorm:"default:30" json:"lockout_duration"`       // 账号锁定时间(分钟)

	// 备份设置
	AutoBackup      bool   `gorm:"default:false" json:"auto_backup"`                // 是否自动备份
	BackupFrequency string `gorm:"size:20;default:'daily'" json:"backup_frequency"` // 备份频率
	BackupTime      string `gorm:"size:5;default:'00:00'" json:"backup_time"`       // 备份时间
	BackupRetention int    `gorm:"default:30" json:"backup_retention"`              // 备份保留天数

	// 其他设置
	Timezone   string `gorm:"size:50;default:'Asia/Shanghai'" json:"timezone"` // 时区
	DateFormat string `gorm:"size:20;default:'YYYY-MM-DD'" json:"date_format"` // 日期格式
	TimeFormat string `gorm:"size:2;default:'24'" json:"time_format"`          // 时间格式
	Currency   string `gorm:"size:3;default:'CNY'" json:"currency"`            // 货币单位
	Language   string `gorm:"size:10;default:'zh-CN'" json:"language"`         // 系统语言
}

// TableName 指定表名
func (SystemSettings) TableName() string {
	return "system_settings"
}
