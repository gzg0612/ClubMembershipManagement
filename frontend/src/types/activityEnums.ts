export enum ActivityType {
  DISCOUNT = 'DISCOUNT',
  COUPON = 'COUPON',
  POINTS = 'POINTS',
  SPECIAL = 'SPECIAL',
}

export enum ActivityStatus {
  DRAFT = 'DRAFT',
  PENDING = 'PENDING',
  ACTIVE = 'ACTIVE',
  PAUSED = 'PAUSED',
  ENDED = 'ENDED',
  REJECTED = 'REJECTED',
}

export enum ActivityRuleConditionType {
  AMOUNT = 'AMOUNT',
  QUANTITY = 'QUANTITY',
  MEMBER_LEVEL = 'MEMBER_LEVEL',
  MEMBER_POINTS = 'MEMBER_POINTS',
}

export enum ActivityRuleRewardType {
  DISCOUNT = 'DISCOUNT',
  AMOUNT = 'AMOUNT',
  POINTS = 'POINTS',
  COUPON = 'COUPON',
} 