package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TransactionType 交易类型
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"  // 收入
	TransactionTypeExpense TransactionType = "expense" // 支出
)

// TransactionStatus 交易状态
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"   // 待处理
	TransactionStatusCompleted TransactionStatus = "completed" // 已完成
	TransactionStatusFailed    TransactionStatus = "failed"    // 失败
)

// Transaction 交易记录模型
type Transaction struct {
	ID            uint              `gorm:"primarykey" json:"id"`
	OrderID       *uint             `gorm:"index:idx_order_id" json:"orderId"`                          // 关联订单ID
	Type          TransactionType   `gorm:"type:varchar(20);not null;index:idx_type" json:"type"`       // 交易类型
	Amount        float64           `gorm:"type:decimal(10,2);not null;index:idx_amount" json:"amount"` // 交易金额
	Status        TransactionStatus `gorm:"type:varchar(20);not null;index:idx_status" json:"status"`   // 交易状态
	PaymentMethod string            `gorm:"type:varchar(50)" json:"paymentMethod"`                      // 支付方式
	Remark        string            `gorm:"type:varchar(255)" json:"remark"`                            // 备注
	CreatedAt     time.Time         `gorm:"index:idx_created_at" json:"createdAt"`                      // 创建时间
	UpdatedAt     time.Time         `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
	Order         *Order            `gorm:"foreignKey:OrderID" json:"order"` // 关联订单
}

// Validate 验证交易数据
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return fmt.Errorf("交易金额必须大于0")
	}
	if t.Type != TransactionTypeIncome && t.Type != TransactionTypeExpense {
		return fmt.Errorf("无效的交易类型")
	}
	if t.Status != TransactionStatusPending && t.Status != TransactionStatusCompleted && t.Status != TransactionStatusFailed {
		return fmt.Errorf("无效的交易状态")
	}
	return nil
}

// BeforeCreate 创建交易记录前的处理
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	// 设置默认状态
	if t.Status == "" {
		t.Status = TransactionStatusPending
	}
	return nil
}

// FinancialSummary 财务汇总模型
type FinancialSummary struct {
	TotalIncome  float64 `json:"totalIncome"`  // 总收入
	TotalExpense float64 `json:"totalExpense"` // 总支出
	NetIncome    float64 `json:"netIncome"`    // 净收入
	Period       string  `json:"period"`       // 统计周期
}

// FinancialReport 财务报表模型
type FinancialReport struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	StartDate time.Time      `gorm:"not null;index:idx_date" json:"startDate"` // 开始日期
	EndDate   time.Time      `gorm:"not null;index:idx_date" json:"endDate"`   // 结束日期
	Type      string         `gorm:"type:varchar(20);not null" json:"type"`    // 报表类型
	Data      string         `gorm:"type:text" json:"data"`                    // 报表数据(JSON)
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
