import request from '@/utils/request'

// 获取会员列表
export function getMemberList(params) {
  return request({
    url: '/api/v1/members',
    method: 'get',
    params
  })
}

// 创建会员
export function createMember(data) {
  return request({
    url: '/api/v1/members',
    method: 'post',
    data
  })
}

// 更新会员
export function updateMember(id, data) {
  return request({
    url: `/api/v1/members/${id}`,
    method: 'put',
    data
  })
}

// 删除会员
export function deleteMember(id) {
  return request({
    url: `/api/v1/members/${id}`,
    method: 'delete'
  })
}

// 会员充值
export function rechargeMember(id, data) {
  return request({
    url: `/api/v1/members/${id}/recharge`,
    method: 'post',
    data
  })
}

// 会员消费
export function consumeMember(id, data) {
  return request({
    url: `/api/v1/members/${id}/consume`,
    method: 'post',
    data
  })
}

// 获取交易记录
export function getTransactionList(params) {
  return request({
    url: '/api/v1/transactions',
    method: 'get',
    params
  })
}

// 导出会员数据
export function exportMembers(params) {
  return request({
    url: '/api/v1/members/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
} 