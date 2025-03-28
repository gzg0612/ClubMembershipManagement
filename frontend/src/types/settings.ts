// 系统设置类型
export interface SystemSettings {
  // 基本设置
  siteName: string;           // 系统名称
  siteLogo: string;          // 系统Logo
  siteDescription: string;   // 系统描述
  contactPhone: string;      // 联系电话
  contactEmail: string;      // 联系邮箱
  address: string;           // 地址

  // 会员设置
  defaultMemberLevel: number;    // 默认会员等级
  enableAutoUpgrade: boolean;    // 是否启用自动升级
  pointsPerYuan: number;        // 每消费1元获得积分
  enablePoints: boolean;        // 是否启用积分系统

  // 通知设置
  enableSMS: boolean;           // 是否启用短信通知
  enableEmail: boolean;         // 是否启用邮件通知
  enableWechat: boolean;        // 是否启用微信通知

  // 安全设置
  passwordMinLength: number;    // 密码最小长度
  requireSpecialChar: boolean;  // 是否要求特殊字符
  requireNumber: boolean;       // 是否要求数字
  requireUppercase: boolean;    // 是否要求大写字母
  loginAttempts: number;        // 最大登录尝试次数
  lockoutDuration: number;      // 账号锁定时间(分钟)

  // 备份设置
  autoBackup: boolean;          // 是否自动备份
  backupFrequency: string;      // 备份频率
  backupTime: string;          // 备份时间
  backupRetention: number;     // 备份保留天数

  // 其他设置
  timezone: string;            // 时区
  dateFormat: string;          // 日期格式
  timeFormat: string;          // 时间格式
  currency: string;            // 货币单位
  language: string;            // 系统语言
}

// 系统设置更新请求
export interface UpdateSettingsRequest {
  settings: Partial<SystemSettings>;
}

// 系统设置响应
export interface SettingsResponse {
  code: number;
  message: string;
  data: SystemSettings;
} 