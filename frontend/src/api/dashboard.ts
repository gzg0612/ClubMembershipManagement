import { request } from '@/utils/request'

// 统计数据接口
export interface DashboardStats {
  memberCount: number
  storeCount: number
  todayRevenue: number
  todayNewMembers: number
}

// 图表数据类型
export interface ChartData {
  date: string
  value: number
}

// 图表数据接口
export interface ChartResponse {
  data: ChartData[]
}

// 仪表盘 API
export const dashboardApi = {
  // 获取统计数据
  getStats: () => {
    return request<DashboardStats>({
      url: '/api/v1/dashboard/stats',
      method: 'get'
    })
  },

  // 获取会员增长趋势
  getMemberTrend: (type: 'week' | 'month' | 'year') => {
    return request<ChartResponse>({
      url: '/api/v1/dashboard/member-trend',
      method: 'get',
      params: { type }
    })
  },

  // 获取营业额统计
  getRevenueTrend: (type: 'week' | 'month' | 'year') => {
    return request<ChartResponse>({
      url: '/api/v1/dashboard/revenue-trend',
      method: 'get',
      params: { type }
    })
  }
} 