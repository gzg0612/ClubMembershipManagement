package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusCompleted OrderStatus = "completed" // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
)

// Order 订单模型
type Order struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	OrderNo     string         `gorm:"type:varchar(32);uniqueIndex:idx_order_no" json:"orderNo"`              // 订单编号
	MemberID    uint           `gorm:"not null;index:idx_member_id" json:"memberId"`                          // 会员ID
	StoreID     uint           `gorm:"not null;index:idx_store_id" json:"storeId"`                            // 店铺ID
	TotalAmount float64        `gorm:"type:decimal(10,2);not null;index:idx_total_amount" json:"totalAmount"` // 订单总金额
	Status      OrderStatus    `gorm:"type:varchar(20);not null;index:idx_status" json:"status"`              // 订单状态
	Remark      string         `gorm:"type:varchar(255)" json:"remark"`                                       // 备注
	CreatedAt   time.Time      `gorm:"index:idx_created_at" json:"createdAt"`                                 // 创建时间
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Member      Member         `gorm:"foreignKey:MemberID" json:"member"` // 关联会员
	Store       Store          `gorm:"foreignKey:StoreID" json:"store"`   // 关联店铺
	Items       []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`   // 订单项
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	OrderID   uint           `gorm:"not null;index:idx_order_id" json:"orderId"`     // 订单ID
	ProductID uint           `gorm:"not null;index:idx_product_id" json:"productId"` // 商品ID
	Quantity  int            `gorm:"not null" json:"quantity"`                       // 数量
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`       // 单价
	Amount    float64        `gorm:"type:decimal(10,2);not null" json:"amount"`      // 小计金额
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Order     Order          `gorm:"foreignKey:OrderID" json:"order"`     // 关联订单
	Product   Product        `gorm:"foreignKey:ProductID" json:"product"` // 关联商品
}

// BeforeCreate 创建订单前生成订单编号
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	// 生成订单编号：年月日时分秒 + 4位随机数
	o.OrderNo = time.Now().Format("20060102150405") + RandomString(4)
	return nil
}

// BeforeSave 保存订单前计算总金额
func (o *Order) BeforeSave(tx *gorm.DB) error {
	// 计算订单总金额
	var total float64
	for _, item := range o.Items {
		total += item.Amount
	}
	o.TotalAmount = total
	return nil
}

// Validate 验证订单数据
func (o *Order) Validate() error {
	if o.MemberID == 0 {
		return fmt.Errorf("会员ID不能为空")
	}
	if o.StoreID == 0 {
		return fmt.Errorf("店铺ID不能为空")
	}
	if len(o.Items) == 0 {
		return fmt.Errorf("订单项不能为空")
	}
	if o.TotalAmount <= 0 {
		return fmt.Errorf("订单总金额必须大于0")
	}
	return nil
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	const letters = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
