package models

import "errors"

var (
	// 通用错误
	ErrInvalidInput = errors.New("无效的输入参数")
	ErrNotFound     = errors.New("记录不存在")
	ErrDuplicate    = errors.New("记录已存在")

	// 用户相关错误
	ErrInvalidUsername = errors.New("无效的用户名")
	ErrInvalidPassword = errors.New("无效的密码")
	ErrUserNotFound    = errors.New("用户不存在")
	ErrUserExists      = errors.New("用户已存在")

	// 媒体相关错误
	ErrInvalidMediaType = errors.New("无效的媒体类型")
	ErrInvalidFilePath  = errors.New("无效的文件路径")
	ErrInvalidStatus    = errors.New("无效的状态")

	// 交易相关错误
	ErrInvalidAmount      = errors.New("无效的金额")
	ErrInsufficientFunds  = errors.New("余额不足")
	ErrInvalidOrderStatus = errors.New("无效的订单状态")

	// 活动相关错误
	ErrInvalidActivityStatus = errors.New("无效的活动状态")
	ErrActivityNotFound      = errors.New("活动不存在")
	ErrActivityExpired       = errors.New("活动已过期")

	ErrInvalidName      = errors.New("活动名称不能为空")
	ErrInvalidType      = errors.New("活动类型不能为空")
	ErrInvalidStartTime = errors.New("开始时间不能为空")
	ErrInvalidEndTime   = errors.New("结束时间不能为空")
	ErrInvalidStoreID   = errors.New("店铺ID不能为空")

	// 数据回收站相关错误
	ErrInvalidRecordType   = errors.New("记录类型不能为空")
	ErrInvalidRecordID     = errors.New("记录ID不能为空")
	ErrInvalidRecordData   = errors.New("记录数据不能为空")
	ErrInvalidOperatorID   = errors.New("操作人ID不能为空")
	ErrInvalidOperatorName = errors.New("操作人姓名不能为空")
)
