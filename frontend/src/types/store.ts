// 店铺状态
export const StoreStatus = {
  Active: 1,
  Inactive: 2,
  Pending: 3
} as const

// 员工角色
export const StoreStaffRole = {
  Manager: 1,
  Staff: 2,
  Technician: 3
} as const

// 员工状态
export const StoreStaffStatus = {
  Active: 1,
  Inactive: 2,
  OnLeave: 3
} as const

// 设备类型
export const DeviceType = {
  Equipment: 1,
  Tool: 2,
  Furniture: 3
} as const

// 设备状态
export const DeviceStatus = {
  Normal: 1,
  Repair: 2,
  Scrapped: 3
} as const

// 检查类型
export const CheckType = {
  Daily: 1,
  Weekly: 2,
  Monthly: 3,
  Quarterly: 4
} as const

// 检查状态
export const CheckStatus = {
  Pass: 1,
  Fail: 2,
  Pending: 3
} as const

// 类型定义
export type StoreStatus = typeof StoreStatus[keyof typeof StoreStatus]
export type StoreStaffRole = typeof StoreStaffRole[keyof typeof StoreStaffRole]
export type StoreStaffStatus = typeof StoreStaffStatus[keyof typeof StoreStaffStatus]
export type DeviceType = typeof DeviceType[keyof typeof DeviceType]
export type DeviceStatus = typeof DeviceStatus[keyof typeof DeviceStatus]
export type CheckType = typeof CheckType[keyof typeof CheckType]
export type CheckStatus = typeof CheckStatus[keyof typeof CheckStatus]

// 店铺接口
export interface Store {
  id: number
  name: string
  code: string
  address: string
  phone: string
  manager: string
  email: string
  businessHours: string
  status: StoreStatus
  description: string
  logo: string
  createdAt: string
  updatedAt: string
}

// 店铺员工接口
export interface StoreStaff {
  id: number
  storeId: number
  userId: number
  role: StoreStaffRole
  status: StoreStaffStatus
  startDate: string
  endDate?: string
  remark: string
  createdAt: string
  updatedAt: string
}

// 店铺设备接口
export interface StoreDevice {
  id: number
  storeId: number
  name: string
  type: DeviceType
  code: string
  status: DeviceStatus
  lastCheck: string
  nextCheck: string
  remark: string
  createdAt: string
  updatedAt: string
}

// 设备检查记录接口
export interface StoreDeviceCheck {
  id: number
  deviceId: number
  storeId: number
  checkType: CheckType
  status: CheckStatus
  content: string
  result: string
  operator: string
  createdAt: string
}

// 分页请求参数接口
export interface PaginationParams {
  page: number
  pageSize: number
  keyword?: string
  status?: number
  startDate?: string
  endDate?: string
}

// 分页响应数据接口
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
} 