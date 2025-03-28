package validators

import (
	"fmt"
	"huiyuanguanli/backend/models"

	"gorm.io/gorm"
)

// OrderValidator 订单验证器
type OrderValidator struct {
	db *gorm.DB
}

// NewOrderValidator 创建订单验证器
func NewOrderValidator(db *gorm.DB) *OrderValidator {
	return &OrderValidator{db: db}
}

// ValidateCreate 验证创建订单
func (v *OrderValidator) ValidateCreate(order *models.Order) error {
	// 验证会员是否存在
	var member models.Member
	if err := v.db.First(&member, order.MemberID).Error; err != nil {
		return fmt.Errorf("会员不存在")
	}

	// 验证店铺是否存在
	var store models.Store
	if err := v.db.First(&store, order.StoreID).Error; err != nil {
		return fmt.Errorf("店铺不存在")
	}

	// 验证订单项
	if len(order.Items) == 0 {
		return fmt.Errorf("订单项不能为空")
	}

	// 验证每个订单项
	for _, item := range order.Items {
		// 验证商品是否存在
		var product models.Product
		if err := v.db.First(&product, item.ProductID).Error; err != nil {
			return fmt.Errorf("商品不存在")
		}

		// 验证商品库存
		if product.Stock < item.Quantity {
			return fmt.Errorf("商品库存不足")
		}

		// 验证商品价格
		if product.Price != item.Price {
			return fmt.Errorf("商品价格已变更")
		}

		// 验证小计金额
		if item.Amount != float64(item.Quantity)*item.Price {
			return fmt.Errorf("订单项金额计算错误")
		}
	}

	// 验证订单总金额
	var total float64
	for _, item := range order.Items {
		total += item.Amount
	}
	if total != order.TotalAmount {
		return fmt.Errorf("订单总金额计算错误")
	}

	return nil
}

// ValidateStatusUpdate 验证订单状态更新
func (v *OrderValidator) ValidateStatusUpdate(orderID uint, newStatus models.OrderStatus) error {
	var order models.Order
	if err := v.db.First(&order, orderID).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}

	// 验证状态转换是否合法
	switch order.Status {
	case models.OrderStatusPending:
		if newStatus != models.OrderStatusPaid && newStatus != models.OrderStatusCancelled {
			return fmt.Errorf("待支付订单只能更新为已支付或已取消")
		}
	case models.OrderStatusPaid:
		if newStatus != models.OrderStatusCompleted {
			return fmt.Errorf("已支付订单只能更新为已完成")
		}
	case models.OrderStatusCompleted:
		return fmt.Errorf("已完成订单不能更新状态")
	case models.OrderStatusCancelled:
		return fmt.Errorf("已取消订单不能更新状态")
	default:
		return fmt.Errorf("无效的订单状态")
	}

	return nil
}
