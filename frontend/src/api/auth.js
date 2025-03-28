import request from '@/utils/request'

// 登录
export function login(data) {
  return request({
    url: '/api/v1/auth/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/api/v1/auth/user-info',
    method: 'get'
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/api/v1/auth/change-password',
    method: 'post',
    data
  })
}

// 退出登录
export function logout() {
  return request({
    url: '/api/v1/auth/logout',
    method: 'post'
  })
} 