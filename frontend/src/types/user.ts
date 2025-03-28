// 用户角色
export const UserRole = {
  SuperAdmin: 1,    // 超级管理员
  Admin: 2,         // 管理员
  Manager: 3,       // 店长
  Staff: 4,         // 店员
  Technician: 5     // 技师
} as const

// 用户状态
export const UserStatus = {
  Active: 1,        // 正常
  Inactive: 2,      // 禁用
  Pending: 3        // 待审核
} as const

// 类型定义
export type UserRole = typeof UserRole[keyof typeof UserRole]
export type UserStatus = typeof UserStatus[keyof typeof UserStatus]

// 用户信息
export interface User {
  id: number
  username: string
  name: string
  phone: string
  email: string
  role: typeof UserRole[keyof typeof UserRole]
  status: typeof UserStatus[keyof typeof UserStatus]
  storeId?: number
  storeName?: string
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
}

// 用户创建参数
export interface CreateUserParams {
  username: string
  password: string
  name: string
  phone: string
  email: string
  role: typeof UserRole[keyof typeof UserRole]
  storeId?: number
}

// 用户更新参数
export interface UpdateUserParams {
  name?: string
  phone?: string
  email?: string
  role?: typeof UserRole[keyof typeof UserRole]
  status?: typeof UserStatus[keyof typeof UserStatus]
  storeId?: number
  password?: string
}

// 用户查询参数
export interface UserQueryParams {
  page: number
  pageSize: number
  keyword?: string
  role?: typeof UserRole[keyof typeof UserRole]
  status?: typeof UserStatus[keyof typeof UserStatus]
  storeId?: number
  startDate?: string
  endDate?: string
  dateRange?: [string, string]
}

// 用户列表响应
export interface UserListResponse {
  total: number
  page: number
  pageSize: number
  items: User[]
}

// 用户登录参数
export interface LoginParams {
  username: string
  password: string
}

// 用户登录响应
export interface LoginResponse {
  token: string
  user: User
}

// 修改密码参数
export interface ChangePasswordParams {
  oldPassword: string
  newPassword: string
} 