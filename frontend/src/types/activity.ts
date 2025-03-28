// 活动类型
export enum ActivityType {
  DISCOUNT = 'discount', // 折扣活动
  COUPON = 'coupon', // 优惠券活动
  POINTS = 'points', // 积分活动
  SPECIAL = 'special', // 特价活动
}

// 活动状态
export enum ActivityStatus {
  DRAFT = 'draft', // 草稿
  PENDING = 'pending', // 待审核
  ACTIVE = 'active', // 进行中
  PAUSED = 'paused', // 已暂停
  ENDED = 'ended', // 已结束
  REJECTED = 'rejected', // 已拒绝
}

// 活动规则条件类型
export enum ActivityRuleConditionType {
  AMOUNT = 'amount', // 消费金额
  QUANTITY = 'quantity', // 购买数量
  MEMBER_LEVEL = 'member_level', // 会员等级
  MEMBER_POINTS = 'member_points', // 会员积分
}

// 活动规则奖励类型
export enum ActivityRuleRewardType {
  DISCOUNT = 'discount', // 折扣
  AMOUNT = 'amount', // 金额
  POINTS = 'points', // 积分
  COUPON = 'coupon', // 优惠券
}

// 活动规则
export interface ActivityRule {
  id: number;
  activityId: number;
  conditionType: ActivityRuleConditionType;
  conditionValue: number;
  rewardType: ActivityRuleRewardType;
  rewardValue: number;
  createdAt: string;
  updatedAt: string;
}

// 活动参与记录
export interface ActivityParticipant {
  id: number;
  activityId: number;
  memberId: number;
  member: {
    id: number;
    name: string;
    phone: string;
  };
  participateTime: string;
  amount: number;
  status: 'success' | 'failed';
  createdAt: string;
  updatedAt: string;
}

// 活动
export interface Activity {
  id: number;
  name: string;
  type: ActivityType;
  status: ActivityStatus;
  startTime: string;
  endTime: string;
  description: string;
  storeId: number;
  store: {
    id: number;
    name: string;
  };
  rules: ActivityRule[];
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

// 活动列表查询参数
export interface ActivityQueryParams {
  page?: number;
  pageSize?: number;
  type?: ActivityType;
  status?: ActivityStatus;
  storeId?: number;
  startTime?: string;
  endTime?: string;
}

// 活动列表响应
export interface ActivityListResponse {
  list: Activity[];
  total: number;
  page: number;
  pageSize: number;
}

// 活动参与记录查询参数
export interface ActivityParticipantQueryParams {
  page?: number;
  pageSize?: number;
}

// 活动参与记录列表响应
export interface ActivityParticipantListResponse {
  list: ActivityParticipant[];
  total: number;
  page: number;
  pageSize: number;
}

// 创建活动请求
export interface CreateActivityRequest {
  name: string;
  type: ActivityType;
  status: ActivityStatus;
  startTime: string;
  endTime: string;
  description: string;
  storeId: number;
  rules: Omit<ActivityRule, 'id' | 'activityId' | 'createdAt' | 'updatedAt'>[];
}

// 更新活动请求
export interface UpdateActivityRequest {
  name?: string;
  type?: ActivityType;
  status?: ActivityStatus;
  startTime?: string;
  endTime?: string;
  description?: string;
  storeId?: number;
  rules?: Omit<ActivityRule, 'id' | 'activityId' | 'createdAt' | 'updatedAt'>[];
}

// 更新活动状态请求
export interface UpdateActivityStatusRequest {
  status: ActivityStatus;
} 