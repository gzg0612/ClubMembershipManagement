import request from '@/utils/request'
import {
  Store,
  StoreStaff,
  StoreDevice,
  StoreDeviceCheck,
  PaginationParams,
  PaginatedResponse,
  StoreStatus,
  StoreStaffRole,
  StoreStaffStatus,
  DeviceType,
  DeviceStatus,
  CheckType,
  CheckStatus
} from '@/types/store'

// 店铺管理API
export const storeApi = {
  // 创建店铺
  createStore: (data: Partial<Store>) => {
    return request.post<Store>('/api/v1/stores', data)
  },

  // 更新店铺
  updateStore: (id: number, data: Partial<Store>) => {
    return request.put<Store>(`/api/v1/stores/${id}`, data)
  },

  // 删除店铺
  deleteStore: (id: number) => {
    return request.delete(`/api/v1/stores/${id}`)
  },

  // 获取店铺详情
  getStore: (id: number) => {
    return request.get<Store>(`/api/v1/stores/${id}`)
  },

  // 获取店铺列表
  getStoreList: (params: PaginationParams) => {
    return request.get<PaginatedResponse<Store>>('/api/v1/stores', { params })
  },

  // 店铺员工管理
  staff: {
    // 添加店铺员工
    addStaff: (storeId: number, data: Partial<StoreStaff>) => {
      return request.post<StoreStaff>(`/api/v1/stores/${storeId}/staff`, data)
    },

    // 更新员工信息
    updateStaff: (storeId: number, staffId: number, data: Partial<StoreStaff>) => {
      return request.put<StoreStaff>(`/api/v1/stores/${storeId}/staff/${staffId}`, data)
    },

    // 移除员工
    removeStaff: (storeId: number, staffId: number) => {
      return request.delete(`/api/v1/stores/${storeId}/staff/${staffId}`)
    },

    // 获取员工列表
    getStaffList: (storeId: number, params: PaginationParams) => {
      return request.get<PaginatedResponse<StoreStaff>>(`/api/v1/stores/${storeId}/staff`, { params })
    }
  },

  // 店铺设备管理
  device: {
    // 添加店铺设备
    addDevice: (storeId: number, data: Partial<StoreDevice>) => {
      return request.post<StoreDevice>(`/api/v1/stores/${storeId}/devices`, data)
    },

    // 更新设备信息
    updateDevice: (storeId: number, deviceId: number, data: Partial<StoreDevice>) => {
      return request.put<StoreDevice>(`/api/v1/stores/${storeId}/devices/${deviceId}`, data)
    },

    // 删除设备
    deleteDevice: (storeId: number, deviceId: number) => {
      return request.delete(`/api/v1/stores/${storeId}/devices/${deviceId}`)
    },

    // 获取设备列表
    getDeviceList: (storeId: number, params: PaginationParams) => {
      return request.get<PaginatedResponse<StoreDevice>>(`/api/v1/stores/${storeId}/devices`, { params })
    },

    // 获取设备详情
    getDevice: (storeId: number, deviceId: number) => {
      return request.get<StoreDevice>(`/api/v1/stores/${storeId}/devices/${deviceId}`)
    },

    // 设备检查记录管理
    check: {
      // 添加检查记录
      addCheck: (storeId: number, deviceId: number, data: Partial<StoreDeviceCheck>) => {
        return request.post<StoreDeviceCheck>(`/api/v1/stores/${storeId}/devices/${deviceId}/checks`, data)
      },

      // 获取检查记录列表
      getCheckList: (storeId: number, deviceId: number, params: PaginationParams) => {
        return request.get<PaginatedResponse<StoreDeviceCheck>>(`/api/v1/stores/${storeId}/devices/${deviceId}/checks`, { params })
      }
    }
  }
} 