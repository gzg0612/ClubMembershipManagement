package controllers

import (
	"net/http"
	"strconv"

	"huiyuanguanli/models"
	"huiyuanguanli/services"

	"github.com/gin-gonic/gin"
)

// StoreController 店铺控制器
type StoreController struct {
	storeService *services.StoreService
}

// NewStoreController 创建店铺控制器实例
func NewStoreController(storeService *services.StoreService) *StoreController {
	return &StoreController{
		storeService: storeService,
	}
}

// CreateStore 创建店铺
func (c *StoreController) CreateStore(ctx *gin.Context) {
	var store models.Store
	if err := ctx.ShouldBindJSON(&store); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	if err := c.storeService.CreateStore(&store); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建店铺失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建店铺成功", "data": store})
}

// UpdateStore 更新店铺
func (c *StoreController) UpdateStore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	var store models.Store
	if err := ctx.ShouldBindJSON(&store); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	store.ID = uint(id)

	if err := c.storeService.UpdateStore(&store); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新店铺失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新店铺成功", "data": store})
}

// DeleteStore 删除店铺
func (c *StoreController) DeleteStore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	if err := c.storeService.DeleteStore(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除店铺失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除店铺成功"})
}

// GetStore 获取店铺详情
func (c *StoreController) GetStore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	store, err := c.storeService.GetStore(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "店铺不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取店铺成功", "data": store})
}

// GetStoreList 获取店铺列表
func (c *StoreController) GetStoreList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	keyword := ctx.Query("keyword")
	status, _ := strconv.Atoi(ctx.Query("status"))

	stores, total, err := c.storeService.GetStoreList(page, pageSize, keyword, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取店铺列表失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取店铺列表成功",
		"data": gin.H{
			"items":     stores,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AddStoreStaff 添加店铺员工
func (c *StoreController) AddStoreStaff(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	var staff models.StoreStaff
	if err := ctx.ShouldBindJSON(&staff); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	staff.StoreID = uint(storeID)

	if err := c.storeService.AddStoreStaff(&staff); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加店铺员工失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加店铺员工成功", "data": staff})
}

// UpdateStoreStaff 更新店铺员工信息
func (c *StoreController) UpdateStoreStaff(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	staffID, err := strconv.ParseUint(ctx.Param("staff_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "员工ID格式错误"})
		return
	}

	var staff models.StoreStaff
	if err := ctx.ShouldBindJSON(&staff); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	staff.ID = uint(staffID)
	staff.StoreID = uint(storeID)

	if err := c.storeService.UpdateStoreStaff(&staff); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新店铺员工失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新店铺员工成功", "data": staff})
}

// RemoveStoreStaff 移除店铺员工
func (c *StoreController) RemoveStoreStaff(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	staffID, err := strconv.ParseUint(ctx.Param("staff_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "员工ID格式错误"})
		return
	}

	if err := c.storeService.RemoveStoreStaff(uint(staffID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "移除店铺员工失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "移除店铺员工成功"})
}

// GetStoreStaffList 获取店铺员工列表
func (c *StoreController) GetStoreStaffList(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	staffs, total, err := c.storeService.GetStoreStaffList(uint(storeID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取店铺员工列表失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取店铺员工列表成功",
		"data": gin.H{
			"items":     staffs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AddStoreDevice 添加店铺设备
func (c *StoreController) AddStoreDevice(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	var device models.StoreDevice
	if err := ctx.ShouldBindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	device.StoreID = uint(storeID)

	if err := c.storeService.AddStoreDevice(&device); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加店铺设备失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加店铺设备成功", "data": device})
}

// UpdateStoreDevice 更新店铺设备信息
func (c *StoreController) UpdateStoreDevice(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	deviceID, err := strconv.ParseUint(ctx.Param("device_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID格式错误"})
		return
	}

	var device models.StoreDevice
	if err := ctx.ShouldBindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	device.ID = uint(deviceID)
	device.StoreID = uint(storeID)

	if err := c.storeService.UpdateStoreDevice(&device); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新店铺设备失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新店铺设备成功", "data": device})
}

// DeleteStoreDevice 删除店铺设备
func (c *StoreController) DeleteStoreDevice(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	deviceID, err := strconv.ParseUint(ctx.Param("device_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID格式错误"})
		return
	}

	if err := c.storeService.DeleteStoreDevice(uint(deviceID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除店铺设备失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除店铺设备成功"})
}

// GetStoreDeviceList 获取店铺设备列表
func (c *StoreController) GetStoreDeviceList(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	devices, total, err := c.storeService.GetStoreDeviceList(uint(storeID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取店铺设备列表失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取店铺设备列表成功",
		"data": gin.H{
			"items":     devices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AddDeviceCheck 添加设备检查记录
func (c *StoreController) AddDeviceCheck(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	deviceID, err := strconv.ParseUint(ctx.Param("device_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID格式错误"})
		return
	}

	var check models.StoreDeviceCheck
	if err := ctx.ShouldBindJSON(&check); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "error": err.Error()})
		return
	}

	check.DeviceID = uint(deviceID)
	check.StoreID = uint(storeID)

	// 获取当前用户ID作为操作人
	userID, _ := ctx.Get("user_id")
	check.Operator = strconv.FormatUint(uint64(userID.(uint)), 10)

	if err := c.storeService.AddDeviceCheck(&check); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加设备检查记录失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加设备检查记录成功", "data": check})
}

// GetDeviceCheckList 获取设备检查记录列表
func (c *StoreController) GetDeviceCheckList(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "店铺ID格式错误"})
		return
	}

	deviceID, err := strconv.ParseUint(ctx.Param("device_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID格式错误"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	checks, total, err := c.storeService.GetDeviceCheckList(uint(deviceID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取设备检查记录列表失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取设备检查记录列表成功",
		"data": gin.H{
			"items":     checks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
