package models

import (
	"time"

	"gorm.io/gorm"
)

// Member 会员模型
type Member struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	Name            string    `gorm:"size:50;not null" json:"name"`           // 姓名
	Gender          int       `gorm:"default:0" json:"gender"`                // 性别：0未知，1男，2女
	Phone           string    `gorm:"size:20;not null;unique" json:"phone"`   // 手机号
	CardID          string    `gorm:"size:50;not null;unique" json:"card_id"` // 会员卡号
	StoreID         uint      `gorm:"not null" json:"store_id"`               // 所属店铺ID
	Store           Store     `gorm:"foreignKey:StoreID" json:"store"`        // 关联店铺
	Balance         float64   `gorm:"default:0" json:"balance"`               // 余额
	GiftBalance     float64   `gorm:"default:0" json:"gift_balance"`          // 赠送余额
	TotalRecharge   float64   `gorm:"default:0" json:"total_recharge"`        // 总充值金额
	TotalConsume    float64   `gorm:"default:0" json:"total_consume"`         // 总消费金额
	FirstActiveTime time.Time `json:"first_active_time"`                      // 首次活跃时间
	LastActiveTime  time.Time `json:"last_active_time"`                       // 最后活跃时间
	Status          int       `gorm:"default:1" json:"status"`                // 状态：1正常，0删除
	Remark          string    `gorm:"size:255" json:"remark"`                 // 备注
	CreatedAt       time.Time `json:"created_at"`                             // 创建时间
	UpdatedAt       time.Time `json:"updated_at"`                             // 更新时间
	Tags            []Tag     `gorm:"many2many:member_tags;" json:"tags"`     // 会员标签
}

// Tag 会员标签模型
type Tag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"` // 标签名称
	CreatedAt time.Time `json:"created_at"`                   // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                   // 更新时间
}

// Transaction 交易记录模型
type Transaction struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	MemberID    uint      `gorm:"not null" json:"member_id"`             // 会员ID
	Member      Member    `gorm:"foreignKey:MemberID" json:"member"`     // 关联会员
	StoreID     uint      `gorm:"not null" json:"store_id"`              // 消费店铺ID
	Store       Store     `gorm:"foreignKey:StoreID" json:"store"`       // 关联店铺
	Type        string    `gorm:"size:20;not null" json:"type"`          // 交易类型：RECHARGE充值，CONSUME消费，GIFT赠送
	Amount      float64   `gorm:"not null" json:"amount"`                // 交易金额
	Balance     float64   `gorm:"not null" json:"balance"`               // 交易后余额
	GiftBalance float64   `gorm:"not null" json:"gift_balance"`          // 交易后赠送余额
	ItemType    string    `gorm:"size:50" json:"item_type"`              // 消费项目类型
	ItemID      uint      `json:"item_id"`                               // 消费项目ID
	Description string    `gorm:"size:255" json:"description"`           // 交易描述
	OperatorID  uint      `json:"operator_id"`                           // 操作人ID
	Operator    User      `gorm:"foreignKey:OperatorID" json:"operator"` // 关联操作人
	CreatedAt   time.Time `json:"created_at"`                            // 创建时间
}

// Store 店铺模型
type Store struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"` // 店铺名称
	Address   string    `gorm:"size:255" json:"address"`       // 店铺地址
	Phone     string    `gorm:"size:20" json:"phone"`          // 联系电话
	Status    int       `gorm:"default:1" json:"status"`       // 状态：1正常，0关闭
	CreatedAt time.Time `json:"created_at"`                    // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                    // 更新时间
}

// BeforeCreate 创建会员前的钩子
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	// 设置首次活跃时间和最后活跃时间为当前时间
	now := time.Now()
	m.FirstActiveTime = now
	m.LastActiveTime = now
	return nil
}

// BeforeUpdate 更新会员前的钩子
func (m *Member) BeforeUpdate(tx *gorm.DB) error {
	// 更新最后活跃时间
	m.LastActiveTime = time.Now()
	return nil
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0"`
	GiftAmount  float64 `json:"gift_amount" binding:"min=0"`
	Description string  `json:"description"`
}

// ConsumeRequest 消费请求
type ConsumeRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0"`
	Description string  `json:"description"`
}

// MemberResponse 会员响应
type MemberResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Balance     float64   `json:"balance"`
	GiftBalance float64   `json:"gift_balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TransactionResponse 交易记录响应
type TransactionResponse struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Balance     float64   `json:"balance"`
	GiftBalance float64   `json:"gift_balance"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
