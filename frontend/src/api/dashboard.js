import request from '@/utils/request'

/**
 * 获取控制台统计数据
 * @param {String} period - 时间周期：today/week/month/year
 * @returns {Promise}
 */
export function getStatistics(period = 'month') {
  return request({
    url: '/api/v1/dashboard/statistics',
    method: 'get',
    params: { period }
  })
}

/**
 * 获取销售趋势数据
 * @param {String} period - 时间周期：today/week/month/year
 * @param {String} type - 数据类型：amount/order
 * @returns {Promise}
 */
export function getSalesTrend(period = 'month', type = 'amount') {
  return request({
    url: '/api/v1/dashboard/sales-trend',
    method: 'get',
    params: { period, type }
  })
}

/**
 * 获取会员分布数据
 * @returns {Promise}
 */
export function getMemberDistribution() {
  return request({
    url: '/api/v1/dashboard/member-distribution',
    method: 'get'
  })
}

/**
 * 获取热销商品数据
 * @param {Number} limit - 获取数量限制
 * @returns {Promise}
 */
export function getHotProducts(limit = 5) {
  return request({
    url: '/api/v1/dashboard/hot-products',
    method: 'get',
    params: { limit }
  })
}

/**
 * 获取最近添加的会员数据
 * @param {Number} limit - 获取数量限制
 * @returns {Promise}
 */
export function getRecentMembers(limit = 5) {
  return request({
    url: '/api/v1/dashboard/recent-members',
    method: 'get',
    params: { limit }
  })
}

/**
 * 获取最近活动数据
 * @param {Number} limit - 获取数量限制
 * @returns {Promise}
 */
export function getRecentActivities(limit = 5) {
  return request({
    url: '/api/v1/dashboard/recent-activities',
    method: 'get',
    params: { limit }
  })
} 