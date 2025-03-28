package services

import (
	"errors"
	"time"

	"huiyuanguanli/models"

	"gorm.io/gorm"
)

// ProductService 商品服务
type ProductService struct {
	db *gorm.DB
}

// NewProductService 创建商品服务实例
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(product *models.Product) error {
	// 检查商品编码是否已存在
	var count int64
	s.db.Model(&models.Product{}).Where("code = ?", product.Code).Count(&count)
	if count > 0 {
		return errors.New("商品编码已存在")
	}

	// 创建商品
	return s.db.Create(product).Error
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(product *models.Product) error {
	// 检查商品是否存在
	var existingProduct models.Product
	if err := s.db.First(&existingProduct, product.ID).Error; err != nil {
		return err
	}

	// 如果修改了商品编码，检查新编码是否已存在
	if existingProduct.Code != product.Code {
		var count int64
		s.db.Model(&models.Product{}).Where("code = ?", product.Code).Count(&count)
		if count > 0 {
			return errors.New("商品编码已存在")
		}
	}

	// 更新商品
	return s.db.Save(product).Error
}

// DeleteProduct 删除商品
func (s *ProductService) DeleteProduct(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}

// GetProduct 获取商品详情
func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	var product models.Product
	err := s.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductList 获取商品列表
func (s *ProductService) GetProductList(page, pageSize int, keyword string, categoryID uint, status int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{})

	// 添加查询条件
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateProductStock 更新商品库存
func (s *ProductService) UpdateProductStock(productID uint, quantity int, stockType int, operator string, remark string) error {
	// 开启事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 获取商品信息
	var product models.Product
	if err := tx.First(&product, productID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 计算新库存
	newStock := product.Stock
	if stockType == 1 { // 入库
		newStock += quantity
	} else if stockType == 2 { // 出库
		if newStock < quantity {
			tx.Rollback()
			return errors.New("库存不足")
		}
		newStock -= quantity
	} else if stockType == 3 { // 盘点
		newStock = quantity
	}

	// 更新商品库存
	product.Stock = newStock
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建库存变动记录
	stockLog := models.ProductStockLog{
		ProductID: productID,
		Type:      stockType,
		Quantity:  quantity,
		Before:    product.Stock,
		After:     newStock,
		Remark:    remark,
		Operator:  operator,
		StoreID:   product.StoreID,
		CreatedAt: time.Now(),
	}
	if err := tx.Create(&stockLog).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateProductPrice 更新商品价格
func (s *ProductService) UpdateProductPrice(productID uint, price float64, priceType int, operator string, remark string) error {
	// 开启事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 获取商品信息
	var product models.Product
	if err := tx.First(&product, productID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新价格
	var oldPrice float64
	if priceType == 1 { // 售价
		oldPrice = product.Price
		product.Price = price
	} else if priceType == 2 { // 成本价
		oldPrice = product.CostPrice
		product.CostPrice = price
	}

	// 保存商品
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建价格变动记录
	priceLog := models.ProductPriceLog{
		ProductID: productID,
		Type:      priceType,
		Before:    oldPrice,
		After:     price,
		Remark:    remark,
		Operator:  operator,
		StoreID:   product.StoreID,
		CreatedAt: time.Now(),
	}
	if err := tx.Create(&priceLog).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateProductCategory 创建商品分类
func (s *ProductService) CreateProductCategory(category *models.ProductCategory) error {
	// 如果有父分类，设置层级
	if category.ParentID > 0 {
		var parentCategory models.ProductCategory
		if err := s.db.First(&parentCategory, category.ParentID).Error; err != nil {
			return err
		}
		category.Level = parentCategory.Level + 1
	} else {
		category.Level = 1
	}

	return s.db.Create(category).Error
}

// UpdateProductCategory 更新商品分类
func (s *ProductService) UpdateProductCategory(category *models.ProductCategory) error {
	return s.db.Save(category).Error
}

// DeleteProductCategory 删除商品分类
func (s *ProductService) DeleteProductCategory(id uint) error {
	// 检查是否有子分类
	var count int64
	s.db.Model(&models.ProductCategory{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下有子分类，无法删除")
	}

	// 检查是否有关联商品
	s.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下有商品，无法删除")
	}

	return s.db.Delete(&models.ProductCategory{}, id).Error
}

// GetProductCategoryList 获取商品分类列表
func (s *ProductService) GetProductCategoryList(storeID uint) ([]models.ProductCategory, error) {
	var categories []models.ProductCategory
	err := s.db.Where("store_id = ?", storeID).Order("sort ASC").Find(&categories).Error
	return categories, err
}

// GetProductStockLogs 获取商品库存变动记录
func (s *ProductService) GetProductStockLogs(productID uint, page, pageSize int) ([]models.ProductStockLog, int64, error) {
	var logs []models.ProductStockLog
	var total int64

	query := s.db.Model(&models.ProductStockLog{}).Where("product_id = ?", productID)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetProductPriceLogs 获取商品价格变动记录
func (s *ProductService) GetProductPriceLogs(productID uint, page, pageSize int) ([]models.ProductPriceLog, int64, error) {
	var logs []models.ProductPriceLog
	var total int64

	query := s.db.Model(&models.ProductPriceLog{}).Where("product_id = ?", productID)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
