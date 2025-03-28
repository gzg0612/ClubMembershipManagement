package models

import (
	"time"
)

// Product 商品模型
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`        // 商品名称
	Code        string    `json:"code" gorm:"size:50;uniqueIndex"`      // 商品编码
	CategoryID  uint      `json:"category_id" gorm:"not null"`          // 商品分类ID
	Price       float64   `json:"price" gorm:"type:decimal(10,2)"`      // 售价
	CostPrice   float64   `json:"cost_price" gorm:"type:decimal(10,2)"` // 成本价
	Stock       int       `json:"stock" gorm:"default:0"`               // 库存数量
	Unit        string    `json:"unit" gorm:"size:20"`                  // 单位
	Description string    `json:"description" gorm:"type:text"`         // 商品描述
	Image       string    `json:"image" gorm:"size:255"`                // 商品图片
	Status      int       `json:"status" gorm:"default:1"`              // 状态：1上架，0下架
	StoreID     uint      `json:"store_id" gorm:"not null"`             // 所属店铺ID
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"` // 软删除
}

// ProductCategory 商品分类模型
type ProductCategory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null"` // 分类名称
	ParentID  uint      `json:"parent_id" gorm:"default:0"`   // 父分类ID
	Level     int       `json:"level" gorm:"default:1"`       // 分类层级
	Sort      int       `json:"sort" gorm:"default:0"`        // 排序
	Status    int       `json:"status" gorm:"default:1"`      // 状态：1启用，0禁用
	StoreID   uint      `json:"store_id" gorm:"not null"`     // 所属店铺ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"` // 软删除
}

// ProductStockLog 商品库存变动记录
type ProductStockLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null"` // 商品ID
	Type      int       `json:"type" gorm:"not null"`       // 变动类型：1入库，2出库，3盘点
	Quantity  int       `json:"quantity" gorm:"not null"`   // 变动数量
	Before    int       `json:"before" gorm:"not null"`     // 变动前库存
	After     int       `json:"after" gorm:"not null"`      // 变动后库存
	Remark    string    `json:"remark" gorm:"size:255"`     // 备注
	Operator  string    `json:"operator" gorm:"size:50"`    // 操作人
	StoreID   uint      `json:"store_id" gorm:"not null"`   // 所属店铺ID
	CreatedAt time.Time `json:"created_at"`
}

// ProductPriceLog 商品价格变动记录
type ProductPriceLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null"`       // 商品ID
	Type      int       `json:"type" gorm:"not null"`             // 变动类型：1售价，2成本价
	Before    float64   `json:"before" gorm:"type:decimal(10,2)"` // 变动前价格
	After     float64   `json:"after" gorm:"type:decimal(10,2)"`  // 变动后价格
	Remark    string    `json:"remark" gorm:"size:255"`           // 备注
	Operator  string    `json:"operator" gorm:"size:50"`          // 操作人
	StoreID   uint      `json:"store_id" gorm:"not null"`         // 所属店铺ID
	CreatedAt time.Time `json:"created_at"`
}
