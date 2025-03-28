package models

import (
	"time"

	"gorm.io/gorm"
)

// ActivityType 活动类型
type ActivityType string

const (
	ActivityTypeRecharge ActivityType = "recharge" // 储值活动
	ActivityTypeDiscount ActivityType = "discount" // 折扣活动
	ActivityTypeGift     ActivityType = "gift"     // 满赠活动
)

// ActivityStatus 活动状态
type ActivityStatus string

const (
	ActivityStatusDraft     ActivityStatus = "draft"     // 草稿
	ActivityStatusPending   ActivityStatus = "pending"   // 待审核
	ActivityStatusActive    ActivityStatus = "active"    // 进行中
	ActivityStatusPaused    ActivityStatus = "paused"    // 已暂停
	ActivityStatusEnded     ActivityStatus = "ended"     // 已结束
	ActivityStatusRejected  ActivityStatus = "rejected"  // 已拒绝
	ActivityStatusCancelled ActivityStatus = "cancelled" // 已取消
)

// Activity 活动
type Activity struct {
	gorm.Model
	Name        string         `json:"name" gorm:"size:100;not null"`      // 活动名称
	Type        ActivityType   `json:"type" gorm:"size:20;not null"`       // 活动类型
	Status      ActivityStatus `json:"status" gorm:"size:20;not null"`     // 活动状态
	StartTime   time.Time      `json:"startTime" gorm:"not null"`          // 开始时间
	EndTime     time.Time      `json:"endTime" gorm:"not null"`            // 结束时间
	Description string         `json:"description" gorm:"type:text"`       // 活动描述
	StoreID     uint           `json:"storeId" gorm:"not null"`            // 所属店铺ID
	Store       Store          `json:"store" gorm:"foreignKey:StoreID"`    // 关联店铺
	Rules       []ActivityRule `json:"rules" gorm:"foreignKey:ActivityID"` // 活动规则
	CreatedBy   uint           `json:"createdBy" gorm:"not null"`          // 创建人ID
	UpdatedBy   uint           `json:"updatedBy" gorm:"not null"`          // 更新人ID
}

// ActivityRule 活动规则
type ActivityRule struct {
	gorm.Model
	ActivityID     uint     `json:"activityId" gorm:"not null"`            // 活动ID
	Activity       Activity `json:"activity" gorm:"foreignKey:ActivityID"` // 关联活动
	ConditionType  string   `json:"conditionType" gorm:"size:50"`          // 条件类型
	ConditionValue float64  `json:"conditionValue"`                        // 条件值
	RewardType     string   `json:"rewardType" gorm:"size:50"`             // 奖励类型
	RewardValue    float64  `json:"rewardValue"`                           // 奖励值
	Description    string   `json:"description" gorm:"type:text"`          // 规则描述
}

// ActivityParticipant 活动参与记录
type ActivityParticipant struct {
	gorm.Model
	ActivityID      uint      `json:"activityId" gorm:"not null"`            // 活动ID
	Activity        Activity  `json:"activity" gorm:"foreignKey:ActivityID"` // 关联活动
	MemberID        uint      `json:"memberId" gorm:"not null"`              // 会员ID
	Member          Member    `json:"member" gorm:"foreignKey:MemberID"`     // 关联会员
	ParticipateTime time.Time `json:"participateTime" gorm:"not null"`       // 参与时间
	Amount          float64   `json:"amount"`                                // 参与金额
	RewardAmount    float64   `json:"rewardAmount"`                          // 奖励金额
	Status          string    `json:"status" gorm:"size:20"`                 // 参与状态
	Remark          string    `json:"remark" gorm:"type:text"`               // 备注
}

// Validate 验证活动数据
func (a *Activity) Validate() error {
	if a.Name == "" {
		return ErrInvalidName
	}
	if a.Type == "" {
		return ErrInvalidType
	}
	if a.Status == "" {
		return ErrInvalidStatus
	}
	if a.StartTime.IsZero() {
		return ErrInvalidStartTime
	}
	if a.EndTime.IsZero() {
		return ErrInvalidEndTime
	}
	if a.EndTime.Before(a.StartTime) {
		return ErrInvalidEndTime
	}
	if a.StoreID == 0 {
		return ErrInvalidStoreID
	}
	return nil
}

// IsActive 检查活动是否进行中
func (a *Activity) IsActive() bool {
	now := time.Now()
	return a.Status == ActivityStatusActive &&
		now.After(a.StartTime) &&
		now.Before(a.EndTime)
}

// CanParticipate 检查会员是否可以参与活动
func (a *Activity) CanParticipate(member *Member) bool {
	if !a.IsActive() {
		return false
	}
	// TODO: 根据活动规则判断会员是否符合参与条件
	return true
}

// CalculateReward 计算活动奖励
func (a *Activity) CalculateReward(amount float64) float64 {
	var reward float64
	for _, rule := range a.Rules {
		switch rule.ConditionType {
		case "amount":
			if amount >= rule.ConditionValue {
				switch rule.RewardType {
				case "discount":
					reward = amount * (1 - rule.RewardValue)
				case "gift":
					reward = rule.RewardValue
				}
			}
		}
	}
	return reward
}
