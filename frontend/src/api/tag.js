import request from '@/utils/request'

// 获取标签列表
export function getTagList(params) {
  return request({
    url: '/api/v1/tags',
    method: 'get',
    params
  })
}

// 创建标签
export function createTag(data) {
  return request({
    url: '/api/v1/tags',
    method: 'post',
    data
  })
}

// 更新标签
export function updateTag(id, data) {
  return request({
    url: `/api/v1/tags/${id}`,
    method: 'put',
    data
  })
}

// 删除标签
export function deleteTag(id) {
  return request({
    url: `/api/v1/tags/${id}`,
    method: 'delete'
  })
}

// 获取标签详情
export function getTagDetail(id) {
  return request({
    url: `/api/v1/tags/${id}`,
    method: 'get'
  })
}

// 获取标签下的会员列表
export function getTagMembers(id, params) {
  return request({
    url: `/api/v1/tags/${id}/members`,
    method: 'get',
    params
  })
} 