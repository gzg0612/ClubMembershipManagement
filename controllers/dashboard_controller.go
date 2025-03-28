package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardController 控制台控制器
type DashboardController struct {
	DB *gorm.DB
}

// NewDashboardController 创建控制台控制器
func NewDashboardController(db *gorm.DB) *DashboardController {
	return &DashboardController{DB: db}
}

// GetStatistics 获取统计数据
func (c *DashboardController) GetStatistics(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")

	// 这里实际应该根据period参数从数据库查询统计数据
	// 为了演示，我们返回一些模拟数据

	var (
		memberCount         = 2815
		salesAmount         = 158620
		orderCount          = 1253
		newMemberCount      = 56
		memberGrowthRate    = 12
		salesGrowthRate     = 8
		orderGrowthRate     = 5
		newMemberGrowthRate = -2
	)

	if period == "today" {
		memberCount = 2815
		salesAmount = 5620
		orderCount = 42
		newMemberCount = 3
		memberGrowthRate = 5
		salesGrowthRate = 12
		orderGrowthRate = 8
		newMemberGrowthRate = -5
	} else if period == "week" {
		memberCount = 2815
		salesAmount = 38620
		orderCount = 283
		newMemberCount = 23
		memberGrowthRate = 8
		salesGrowthRate = 10
		orderGrowthRate = 6
		newMemberGrowthRate = -3
	} else if period == "year" {
		memberCount = 2815
		salesAmount = 1958620
		orderCount = 15253
		newMemberCount = 892
		memberGrowthRate = 25
		salesGrowthRate = 15
		orderGrowthRate = 18
		newMemberGrowthRate = 8
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取统计数据成功",
		"data": gin.H{
			"memberCount":         memberCount,
			"salesAmount":         salesAmount,
			"orderCount":          orderCount,
			"newMemberCount":      newMemberCount,
			"memberGrowthRate":    memberGrowthRate,
			"salesGrowthRate":     salesGrowthRate,
			"orderGrowthRate":     orderGrowthRate,
			"newMemberGrowthRate": newMemberGrowthRate,
		},
	})
}

// GetSalesTrend 获取销售趋势数据
func (c *DashboardController) GetSalesTrend(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")
	dataType := ctx.DefaultQuery("type", "amount")

	// 这里实际应该根据period和type参数从数据库查询销售趋势数据
	// 为了演示，我们返回一些模拟数据

	// 模拟数据
	var labels []string
	var data []float64

	now := time.Now()

	if period == "today" {
		// 按小时统计
		labels = []string{"00:00", "02:00", "04:00", "06:00", "08:00", "10:00", "12:00", "14:00", "16:00", "18:00", "20:00", "22:00"}
		if dataType == "amount" {
			data = []float64{320, 280, 150, 320, 580, 920, 1200, 950, 850, 920, 680, 450}
		} else {
			data = []float64{6, 5, 3, 6, 12, 18, 25, 20, 16, 19, 15, 10}
		}
	} else if period == "week" {
		// 按天统计
		labels = []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
		if dataType == "amount" {
			data = []float64{4500, 5200, 4800, 6200, 7800, 8500, 6800}
		} else {
			data = []float64{38, 45, 42, 56, 68, 75, 59}
		}
	} else if period == "year" {
		// 按月统计
		labels = []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
		if dataType == "amount" {
			data = []float64{152000, 148000, 165000, 178000, 185000, 190000, 182000, 175000, 168000, 158000, 162000, 180000}
		} else {
			data = []float64{1250, 1180, 1320, 1420, 1480, 1520, 1460, 1400, 1350, 1260, 1300, 1450}
		}
	} else {
		// 默认按天统计（本月）
		daysInMonth := 30 // 简化处理
		labels = make([]string, daysInMonth)
		data = make([]float64, daysInMonth)

		for i := 0; i < daysInMonth; i++ {
			date := now.AddDate(0, 0, i-daysInMonth+1)
			labels[i] = date.Format("01-02")

			if dataType == "amount" {
				// 销售额在4000-8000之间波动
				data[i] = 4000 + float64(i%4+1)*1000
			} else {
				// 订单量在30-70之间波动
				data[i] = 30 + float64(i%4+1)*10
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取销售趋势数据成功",
		"data": gin.H{
			"labels": labels,
			"data":   data,
			"type":   dataType,
		},
	})
}

// GetMemberDistribution 获取会员分布数据
func (c *DashboardController) GetMemberDistribution(ctx *gin.Context) {
	// 这里实际应该从数据库查询会员分布数据
	// 为了演示，我们返回一些模拟数据

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取会员分布数据成功",
		"data": gin.H{
			"labels": []string{"普通会员", "白银会员", "黄金会员", "铂金会员", "钻石会员"},
			"data":   []int{1250, 850, 420, 220, 75},
		},
	})
}

// GetHotProducts 获取热销商品数据
func (c *DashboardController) GetHotProducts(ctx *gin.Context) {
	_ = ctx.DefaultQuery("limit", "5") // 标记为已使用，实际项目中应该用于限制返回数量

	// 这里实际应该从数据库查询热销商品数据
	// 为了演示，我们返回一些模拟数据

	hotProducts := []gin.H{
		{"name": "高级会员卡", "icon": "el-icon-vip", "sales": 256, "amount": "¥76,800", "stockStatus": "充足"},
		{"name": "护肤礼包", "icon": "el-icon-present", "sales": 128, "amount": "¥38,400", "stockStatus": "偏低"},
		{"name": "精品SPA套餐", "icon": "el-icon-moon", "sales": 96, "amount": "¥28,800", "stockStatus": "紧缺"},
		{"name": "月度按摩卡", "icon": "el-icon-bangzhu", "sales": 72, "amount": "¥21,600", "stockStatus": "充足"},
		{"name": "健身季卡", "icon": "el-icon-medal", "sales": 64, "amount": "¥19,200", "stockStatus": "充足"},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取热销商品数据成功",
		"data":    hotProducts,
	})
}

// GetRecentMembers 获取最近添加的会员数据
func (c *DashboardController) GetRecentMembers(ctx *gin.Context) {
	_ = ctx.DefaultQuery("limit", "5") // 标记为已使用，实际项目中应该用于限制返回数量

	// 这里实际应该从数据库查询最近添加的会员数据
	// 为了演示，我们返回一些模拟数据

	recentMembers := []gin.H{
		{"id": 1001, "name": "张三", "phone": "13800138001", "level": "黄金会员", "joinTime": "2025-03-15"},
		{"id": 1002, "name": "李四", "phone": "13800138002", "level": "白银会员", "joinTime": "2025-03-16"},
		{"id": 1003, "name": "王五", "phone": "13800138003", "level": "铂金会员", "joinTime": "2025-03-17"},
		{"id": 1004, "name": "赵六", "phone": "13800138004", "level": "普通会员", "joinTime": "2025-03-18"},
		{"id": 1005, "name": "钱七", "phone": "13800138005", "level": "黄金会员", "joinTime": "2025-03-19"},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取最近添加的会员数据成功",
		"data":    recentMembers,
	})
}

// GetRecentActivities 获取最近活动数据
func (c *DashboardController) GetRecentActivities(ctx *gin.Context) {
	_ = ctx.DefaultQuery("limit", "5") // 标记为已使用，实际项目中应该用于限制返回数量

	// 这里实际应该从数据库查询最近活动数据
	// 为了演示，我们返回一些模拟数据

	recentActivities := []gin.H{
		{"name": "618购物节特惠", "icon": "el-icon-shopping-bag-1", "startTime": "2025-06-01", "endTime": "2025-06-18", "participants": 1253, "status": "进行中"},
		{"name": "新会员专享礼", "icon": "el-icon-present", "startTime": "2025-05-15", "endTime": "2025-06-15", "participants": 865, "status": "进行中"},
		{"name": "五一假期折扣", "icon": "el-icon-discount", "startTime": "2025-05-01", "endTime": "2025-05-05", "participants": 1056, "status": "已结束"},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取最近活动数据成功",
		"data":    recentActivities,
	})
}
