import request from '@/utils/request'

// 获取店铺列表
export function getStoreList(params) {
  return request({
    url: '/api/v1/stores',
    method: 'get',
    params
  })
}

// 创建店铺
export function createStore(data) {
  return request({
    url: '/api/v1/stores',
    method: 'post',
    data
  })
}

// 更新店铺
export function updateStore(id, data) {
  return request({
    url: `/api/v1/stores/${id}`,
    method: 'put',
    data
  })
}

// 删除店铺
export function deleteStore(id) {
  return request({
    url: `/api/v1/stores/${id}`,
    method: 'delete'
  })
}

// 获取店铺详情
export function getStoreDetail(id) {
  return request({
    url: `/api/v1/stores/${id}`,
    method: 'get'
  })
}

// 获取店铺统计数据
export function getStoreStatistics(id, params) {
  return request({
    url: `/api/v1/stores/${id}/statistics`,
    method: 'get',
    params
  })
} 