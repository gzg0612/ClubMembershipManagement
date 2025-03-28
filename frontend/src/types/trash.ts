// 记录类型
export enum RecordType {
  MEMBER = 'member',   // 会员
  PRODUCT = 'product', // 商品
  ACTIVITY = 'activity' // 活动
}

// 删除记录
export interface DeletedRecord {
  id: number;
  recordType: RecordType;
  recordId: number;
  recordData: string;
  deleteReason: string;
  operatorId: number;
  operatorName: string;
  deletedAt: string;
  createdAt: string;
  updatedAt: string;
}

// 获取删除记录列表请求参数
export interface GetDeletedRecordsParams {
  type?: RecordType;
  page?: number;
  pageSize?: number;
}

// 获取删除记录列表响应
export interface GetDeletedRecordsResponse {
  list: DeletedRecord[];
  total: number;
  page: number;
  pageSize: number;
}

// 获取单个删除记录响应
export interface GetDeletedRecordResponse {
  data: DeletedRecord;
} 