package dto

// CreateMemberRequest 创建会员请求
type CreateMemberRequest struct {
	Name    string  `json:"name" binding:"required,max=50"`    // 姓名
	Gender  int     `json:"gender" binding:"oneof=0 1 2"`      // 性别：0未知，1男，2女
	Phone   string  `json:"phone" binding:"required,len=11"`   // 手机号
	CardID  string  `json:"card_id" binding:"required,max=50"` // 会员卡号
	StoreID uint    `json:"store_id" binding:"required"`       // 所属店铺ID
	Balance float64 `json:"balance" binding:"gte=0"`           // 余额
	Remark  string  `json:"remark" binding:"max=255"`          // 备注
	TagIDs  []uint  `json:"tag_ids"`                           // 标签ID列表
}

// UpdateMemberRequest 更新会员请求
type UpdateMemberRequest struct {
	Name    string `json:"name" binding:"required,max=50"`    // 姓名
	Gender  int    `json:"gender" binding:"oneof=0 1 2"`      // 性别：0未知，1男，2女
	Phone   string `json:"phone" binding:"required,len=11"`   // 手机号
	CardID  string `json:"card_id" binding:"required,max=50"` // 会员卡号
	StoreID uint   `json:"store_id" binding:"required"`       // 所属店铺ID
	Remark  string `json:"remark" binding:"max=255"`          // 备注
	TagIDs  []uint `json:"tag_ids"`                           // 标签ID列表
}

// MemberListRequest 会员列表请求
type MemberListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`      // 页码
	PageSize int    `form:"page_size" binding:"required,min=1"` // 每页数量
	Name     string `form:"name"`                               // 姓名
	Phone    string `form:"phone"`                              // 手机号
	CardID   string `form:"card_id"`                            // 会员卡号
	StoreID  uint   `form:"store_id"`                           // 店铺ID
	TagID    uint   `form:"tag_id"`                             // 标签ID
	Status   int    `form:"status" binding:"oneof=0 1"`         // 状态：1正常，0删除
}

// MemberListResponse 会员列表响应
type MemberListResponse struct {
	Total int64          `json:"total"` // 总数
	Items []MemberDetail `json:"items"` // 会员列表
}

// MemberDetail 会员详情
type MemberDetail struct {
	ID              uint    `json:"id"`                // 会员ID
	Name            string  `json:"name"`              // 姓名
	Gender          int     `json:"gender"`            // 性别
	Phone           string  `json:"phone"`             // 手机号
	CardID          string  `json:"card_id"`           // 会员卡号
	StoreID         uint    `json:"store_id"`          // 所属店铺ID
	StoreName       string  `json:"store_name"`        // 店铺名称
	Balance         float64 `json:"balance"`           // 余额
	GiftBalance     float64 `json:"gift_balance"`      // 赠送余额
	TotalRecharge   float64 `json:"total_recharge"`    // 总充值金额
	TotalConsume    float64 `json:"total_consume"`     // 总消费金额
	FirstActiveTime string  `json:"first_active_time"` // 首次活跃时间
	LastActiveTime  string  `json:"last_active_time"`  // 最后活跃时间
	Status          int     `json:"status"`            // 状态
	Remark          string  `json:"remark"`            // 备注
	Tags            []Tag   `json:"tags"`              // 标签列表
	CreatedAt       string  `json:"created_at"`        // 创建时间
	UpdatedAt       string  `json:"updated_at"`        // 更新时间
}

// Tag 标签信息
type Tag struct {
	ID   uint   `json:"id"`   // 标签ID
	Name string `json:"name"` // 标签名称
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"` // 充值金额
	GiftAmount  float64 `json:"gift_amount" binding:"gte=0"`    // 赠送金额
	Description string  `json:"description" binding:"max=255"`  // 充值说明
}

// ConsumeRequest 消费请求
type ConsumeRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`      // 消费金额
	ItemType    string  `json:"item_type" binding:"required,max=50"` // 消费项目类型
	ItemID      uint    `json:"item_id" binding:"required"`          // 消费项目ID
	Description string  `json:"description" binding:"max=255"`       // 消费说明
}

// TransactionListRequest 交易记录列表请求
type TransactionListRequest struct {
	Page      int    `form:"page" binding:"required,min=1"`      // 页码
	PageSize  int    `form:"page_size" binding:"required,min=1"` // 每页数量
	MemberID  uint   `form:"member_id" binding:"required"`       // 会员ID
	Type      string `form:"type"`                               // 交易类型
	StartTime string `form:"start_time"`                         // 开始时间
	EndTime   string `form:"end_time"`                           // 结束时间
}

// TransactionListResponse 交易记录列表响应
type TransactionListResponse struct {
	Total int64               `json:"total"` // 总数
	Items []TransactionDetail `json:"items"` // 交易记录列表
}

// TransactionDetail 交易记录详情
type TransactionDetail struct {
	ID           uint    `json:"id"`            // 交易ID
	Type         string  `json:"type"`          // 交易类型
	Amount       float64 `json:"amount"`        // 交易金额
	Balance      float64 `json:"balance"`       // 交易后余额
	GiftBalance  float64 `json:"gift_balance"`  // 交易后赠送余额
	ItemType     string  `json:"item_type"`     // 消费项目类型
	ItemID       uint    `json:"item_id"`       // 消费项目ID
	Description  string  `json:"description"`   // 交易描述
	StoreID      uint    `json:"store_id"`      // 消费店铺ID
	StoreName    string  `json:"store_name"`    // 店铺名称
	OperatorID   uint    `json:"operator_id"`   // 操作人ID
	OperatorName string  `json:"operator_name"` // 操作人姓名
	CreatedAt    string  `json:"created_at"`    // 创建时间
}
