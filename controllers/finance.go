package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"huiyuanguanli/models"
	"huiyuanguanli/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// Transaction 交易记录
type Transaction struct {
	gorm.Model
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	Status        string  `json:"status"`
	Remark        string  `json:"remark"`
	OrderID       *uint   `json:"orderId,omitempty"`
}

// FinanceController 财务控制器
type FinanceController struct {
	db *gorm.DB
}

// NewFinanceController 创建财务控制器
func NewFinanceController(db *gorm.DB) *FinanceController {
	return &FinanceController{db: db}
}

// CreateTransaction 创建交易记录
func (c *FinanceController) CreateTransaction(ctx *gin.Context) {
	var transaction models.Transaction
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证交易数据
	if err := transaction.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "交易数据验证失败", "error": err.Error()})
		return
	}

	// 开启事务
	tx := c.db.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建交易记录失败"})
		return
	}

	// 创建交易记录
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建交易记录失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建交易记录失败"})
		return
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("transaction:%d", transaction.ID)
	utils.DeleteCache(cacheKey)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": transaction})
}

// GetTransaction 获取交易记录详情
func (c *FinanceController) GetTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	// 尝试从缓存获取
	var transaction models.Transaction
	cacheKey := fmt.Sprintf("transaction:%s", id)
	if err := utils.GetCache(cacheKey, &transaction); err == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": transaction})
		return
	}

	// 从数据库获取
	if err := c.db.Preload("Order").First(&transaction, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "交易记录不存在"})
		return
	}

	// 缓存交易记录数据
	utils.Cache(cacheKey, transaction, time.Hour)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": transaction})
}

// UpdateTransactionStatus 更新交易状态
func (c *FinanceController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		Status models.TransactionStatus `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 开启事务
	tx := c.db.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新交易状态失败"})
		return
	}

	// 更新交易状态
	if err := tx.Model(&models.Transaction{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新交易状态失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新交易状态失败"})
		return
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("transaction:%s", id)
	utils.DeleteCache(cacheKey)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

// GetTransactionList 获取交易记录列表
func (c *FinanceController) GetTransactionList(ctx *gin.Context) {
	var transactions []models.Transaction
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	transactionType := ctx.Query("type")
	status := ctx.Query("status")
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	// 转换分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码格式错误"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "每页数量格式错误"})
		return
	}

	// 构建查询条件
	query := c.db.Model(&models.Transaction{}).Preload("Order")

	// 添加查询条件
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	// 使用索引优化排序
	query = query.Order("created_at DESC")

	// 分页查询
	var total int64
	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取交易记录列表失败"})
		return
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取交易记录列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     transactions,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetFinancialSummary 获取财务汇总
func (c *FinanceController) GetFinancialSummary(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month") // 默认按月统计
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	// 构建查询条件
	query := c.db.Model(&models.Transaction{})

	// 添加时间范围条件
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	// 计算收入
	var totalIncome float64
	if err := query.Where("type = ? AND status = ?", models.TransactionTypeIncome, models.TransactionStatusCompleted).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalIncome); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取财务汇总失败"})
		return
	}

	// 计算支出
	var totalExpense float64
	if err := query.Where("type = ? AND status = ?", models.TransactionTypeExpense, models.TransactionStatusCompleted).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalExpense); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取财务汇总失败"})
		return
	}

	summary := models.FinancialSummary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetIncome:    totalIncome - totalExpense,
		Period:       period,
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": summary})
}

// GenerateFinancialReport 生成财务报表
func (c *FinanceController) GenerateFinancialReport(ctx *gin.Context) {
	var req struct {
		StartDate time.Time `json:"startDate" binding:"required"`
		EndDate   time.Time `json:"endDate" binding:"required"`
		Type      string    `json:"type" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 构建查询条件
	query := c.db.Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at <= ?", req.StartDate, req.EndDate)

	// 根据报表类型生成数据
	var reportData interface{}
	switch req.Type {
	case "daily":
		// 按日统计
		var dailyStats []struct {
			Date      string  `json:"date"`
			Income    float64 `json:"income"`
			Expense   float64 `json:"expense"`
			NetIncome float64 `json:"netIncome"`
		}
		if err := query.Select("DATE(created_at) as date, "+
			"SUM(CASE WHEN type = ? AND status = ? THEN amount ELSE 0 END) as income, "+
			"SUM(CASE WHEN type = ? AND status = ? THEN amount ELSE 0 END) as expense",
			models.TransactionTypeIncome, models.TransactionStatusCompleted,
			models.TransactionTypeExpense, models.TransactionStatusCompleted).
			Group("DATE(created_at)").
			Order("date").
			Find(&dailyStats).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成财务报表失败"})
			return
		}
		reportData = dailyStats

	case "monthly":
		// 按月统计
		var monthlyStats []struct {
			Month     string  `json:"month"`
			Income    float64 `json:"income"`
			Expense   float64 `json:"expense"`
			NetIncome float64 `json:"netIncome"`
		}
		if err := query.Select("DATE_FORMAT(created_at, '%Y-%m') as month, "+
			"SUM(CASE WHEN type = ? AND status = ? THEN amount ELSE 0 END) as income, "+
			"SUM(CASE WHEN type = ? AND status = ? THEN amount ELSE 0 END) as expense",
			models.TransactionTypeIncome, models.TransactionStatusCompleted,
			models.TransactionTypeExpense, models.TransactionStatusCompleted).
			Group("DATE_FORMAT(created_at, '%Y-%m')").
			Order("month").
			Find(&monthlyStats).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成财务报表失败"})
			return
		}
		reportData = monthlyStats

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的报表类型"})
		return
	}

	// 将报表数据转换为JSON
	data, err := json.Marshal(reportData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成财务报表失败"})
		return
	}

	// 创建报表记录
	report := models.FinancialReport{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Type:      req.Type,
		Data:      string(data),
	}

	if err := c.db.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存财务报表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": report})
}

// ImportPreviewResponse 导入预览响应
type ImportPreviewResponse struct {
	FileID  string                   `json:"fileId"`
	Headers []string                 `json:"headers"`
	Preview []map[string]interface{} `json:"preview"`
	Mapping map[string]string        `json:"mapping"`
	Errors  []ImportError            `json:"errors,omitempty"`
}

// ImportError 导入错误信息
type ImportError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// ImportResult 导入结果
type ImportResult struct {
	Success        int           `json:"success"`
	Failed         int           `json:"failed"`
	Total          int           `json:"total"`
	Errors         []ImportError `json:"errors"`
	SuccessRecords []Transaction `json:"successRecords"`
}

// GetImportTemplate 获取导入模板
func (c *FinanceController) GetImportTemplate(ctx *gin.Context) {
	templateType := ctx.Query("type")
	if templateType != "basic" && templateType != "detailed" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板类型"})
		return
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置表头
	headers := []string{"交易类型", "金额", "支付方式"}
	if templateType == "detailed" {
		headers = append(headers, "状态", "备注", "订单号")
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// 设置示例数据
	if templateType == "basic" {
		f.SetCellValue("Sheet1", "A2", "income")
		f.SetCellValue("Sheet1", "B2", "100.00")
		f.SetCellValue("Sheet1", "C2", "支付宝")
	} else {
		f.SetCellValue("Sheet1", "A2", "income")
		f.SetCellValue("Sheet1", "B2", "100.00")
		f.SetCellValue("Sheet1", "C2", "支付宝")
		f.SetCellValue("Sheet1", "D2", "completed")
		f.SetCellValue("Sheet1", "E2", "会员充值")
		f.SetCellValue("Sheet1", "F2", "ORDER123456")
	}

	// 设置响应头
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=transaction_template_%s.xlsx", templateType))

	// 写入响应
	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成模板失败"})
		return
	}
}

// PreviewImport 预览导入数据
func (c *FinanceController) PreviewImport(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
		return
	}
	defer src.Close()

	// 解析Excel文件
	f, err := excelize.OpenReader(src)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "解析Excel文件失败"})
		return
	}
	defer f.Close()

	// 读取第一个sheet
	rows, err := f.GetRows("Sheet1")
	if err != nil || len(rows) < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的Excel文件格式"})
		return
	}

	// 获取表头
	headers := rows[0]

	// 预览数据（最多10条）
	previewCount := 10
	if len(rows) < previewCount {
		previewCount = len(rows)
	}

	preview := make([]map[string]interface{}, 0, previewCount-1)
	errors := make([]ImportError, 0)

	// 自动映射字段
	mapping := map[string]string{}
	for i, header := range headers {
		switch strings.ToLower(header) {
		case "交易类型", "type":
			mapping["type"] = headers[i]
		case "金额", "amount":
			mapping["amount"] = headers[i]
		case "支付方式", "payment_method":
			mapping["paymentMethod"] = headers[i]
		case "状态", "status":
			mapping["status"] = headers[i]
		case "备注", "remark":
			mapping["remark"] = headers[i]
		}
	}

	// 生成预览数据
	for i := 1; i < previewCount; i++ {
		row := rows[i]
		data := make(map[string]interface{})

		for j, cell := range row {
			if j < len(headers) {
				data[headers[j]] = cell
			}
		}

		// 验证数据
		if err := c.validateImportRow(data, mapping); err != nil {
			errors = append(errors, ImportError{
				Row:     i + 1,
				Message: err.Error(),
			})
		}

		preview = append(preview, data)
	}

	// 生成临时文件ID
	fileID := uuid.New().String()

	// 保存文件到临时目录
	tempPath := filepath.Join(os.TempDir(), fileID+".xlsx")
	if err := f.SaveAs(tempPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存临时文件失败"})
		return
	}

	ctx.JSON(http.StatusOK, ImportPreviewResponse{
		FileID:  fileID,
		Headers: headers,
		Preview: preview,
		Mapping: mapping,
		Errors:  errors,
	})
}

// ConfirmImport 确认导入数据
func (c *FinanceController) ConfirmImport(ctx *gin.Context) {
	var params struct {
		FileID     string            `json:"fileId" binding:"required"`
		Mapping    map[string]string `json:"mapping" binding:"required"`
		SkipErrors bool              `json:"skipErrors"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 读取临时文件
	tempPath := filepath.Join(os.TempDir(), params.FileID+".xlsx")
	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "临时文件不存在或已过期"})
		return
	}
	defer f.Close()
	defer os.Remove(tempPath)

	// 读取所有数据
	rows, err := f.GetRows("Sheet1")
	if err != nil || len(rows) < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的Excel文件格式"})
		return
	}

	// 处理导入
	result := ImportResult{
		Total:          len(rows) - 1,
		SuccessRecords: make([]Transaction, 0),
		Errors:         make([]ImportError, 0),
	}

	// 开启事务
	tx := c.db.Begin()

	// 批量导入数据
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		data := make(map[string]interface{})

		// 根据映射关系处理数据
		for field, header := range params.Mapping {
			for j, h := range rows[0] {
				if h == header && j < len(row) {
					data[field] = row[j]
				}
			}
		}

		// 验证数据
		if err := c.validateImportRow(data, params.Mapping); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, ImportError{
				Row:     i + 1,
				Message: err.Error(),
			})
			if !params.SkipErrors {
				tx.Rollback()
				ctx.JSON(http.StatusBadRequest, result)
				return
			}
			continue
		}

		// 创建交易记录
		transaction := Transaction{
			Type:          data["type"].(string),
			Amount:        c.parseAmount(data["amount"].(string)),
			PaymentMethod: data["paymentMethod"].(string),
			Status:        c.getDefaultStatus(data),
			Remark:        c.getString(data, "remark"),
		}

		if err := tx.Create(&transaction).Error; err != nil {
			result.Failed++
			result.Errors = append(result.Errors, ImportError{
				Row:     i + 1,
				Message: "保存记录失败: " + err.Error(),
			})
			if !params.SkipErrors {
				tx.Rollback()
				ctx.JSON(http.StatusBadRequest, result)
				return
			}
			continue
		}

		result.Success++
		result.SuccessRecords = append(result.SuccessRecords, transaction)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// 辅助方法

func (c *FinanceController) validateImportRow(data map[string]interface{}, mapping map[string]string) error {
	// 验证必填字段
	requiredFields := []string{"type", "amount", "paymentMethod"}
	for _, field := range requiredFields {
		header, exists := mapping[field]
		if !exists {
			return fmt.Errorf("缺少必填字段映射: %s", field)
		}
		value, exists := data[header]
		if !exists || value == "" {
			return fmt.Errorf("缺少必填字段: %s", header)
		}
	}

	// 验证交易类型
	if typeHeader := mapping["type"]; typeHeader != "" {
		if value, ok := data[typeHeader].(string); ok {
			if value != "income" && value != "expense" {
				return fmt.Errorf("无效的交易类型: %s", value)
			}
		}
	}

	// 验证金额
	if amountHeader := mapping["amount"]; amountHeader != "" {
		if value, ok := data[amountHeader].(string); ok {
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				return fmt.Errorf("无效的金额格式: %s", value)
			}
		}
	}

	return nil
}

func (c *FinanceController) parseAmount(s string) float64 {
	amount, _ := strconv.ParseFloat(s, 64)
	return amount
}

func (c *FinanceController) getString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func (c *FinanceController) getDefaultStatus(data map[string]interface{}) string {
	if status, ok := data["status"].(string); ok && status != "" {
		return status
	}
	return "completed"
}
