import request from '@/utils/request'
import type {
  User,
  CreateUserParams,
  UpdateUserParams,
  UserQueryParams,
  UserListResponse,
  LoginParams,
  LoginResponse,
  ChangePasswordParams
} from '@/types/user'

// 用户管理 API
export const userApi = {
  // 用户登录
  login: (data: LoginParams) => {
    return request.post<LoginResponse>('/api/v1/auth/login', data)
  },

  // 获取当前用户信息
  getCurrentUser: () => {
    return request.get<User>('/api/v1/auth/current-user')
  },

  // 修改密码
  changePassword: (data: ChangePasswordParams) => {
    return request.put('/api/v1/auth/change-password', data)
  },

  // 创建用户
  createUser: (data: CreateUserParams) => {
    return request.post<User>('/api/v1/users', data)
  },

  // 更新用户信息
  updateUser: (id: number, data: UpdateUserParams) => {
    return request.put<User>(`/api/v1/users/${id}`, data)
  },

  // 删除用户
  deleteUser: (id: number) => {
    return request.delete(`/api/v1/users/${id}`)
  },

  // 获取用户详情
  getUser: (id: number) => {
    return request.get<User>(`/api/v1/users/${id}`)
  },

  // 获取用户列表
  getUserList: (params: UserQueryParams) => {
    return request.get<UserListResponse>('/api/v1/users', { params })
  },

  // 重置用户密码
  resetPassword: (id: number) => {
    return request.put(`/api/v1/users/${id}/reset-password`)
  },

  // 启用/禁用用户
  toggleUserStatus: (id: number, status: number) => {
    return request.put(`/api/v1/users/${id}/status`, { status })
  }
} 