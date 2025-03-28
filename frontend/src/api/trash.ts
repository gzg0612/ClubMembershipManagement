import request from '@/utils/request';
import type {
  GetDeletedRecordsParams,
  GetDeletedRecordsResponse,
  GetDeletedRecordResponse
} from '@/types/trash';

// 获取删除记录列表
export function getDeletedRecords(params: GetDeletedRecordsParams) {
  return request<GetDeletedRecordsResponse>({
    url: '/api/v1/trash',
    method: 'get',
    params
  });
}

// 获取单个删除记录
export function getDeletedRecord(id: number) {
  return request<GetDeletedRecordResponse>({
    url: `/api/v1/trash/${id}`,
    method: 'get'
  });
}

// 恢复删除的记录
export function restoreRecord(id: number) {
  return request({
    url: `/api/v1/trash/${id}/restore`,
    method: 'post'
  });
}

// 永久删除记录
export function permanentlyDelete(id: number) {
  return request({
    url: `/api/v1/trash/${id}`,
    method: 'delete'
  });
} 