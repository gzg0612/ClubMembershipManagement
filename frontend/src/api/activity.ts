import request from '@/utils/request';
import type {
  Activity,
  ActivityQueryParams,
  ActivityListResponse,
  ActivityParticipantQueryParams,
  ActivityParticipantListResponse,
  CreateActivityRequest,
  UpdateActivityRequest,
  UpdateActivityStatusRequest,
} from '@/types/activity';

// 获取活动列表
export function getActivityList(params: ActivityQueryParams) {
  return request<ActivityListResponse>({
    url: '/api/v1/activities',
    method: 'get',
    params,
  });
}

// 获取活动详情
export function getActivity(id: number) {
  return request<Activity>({
    url: `/api/v1/activities/${id}`,
    method: 'get',
  });
}

// 创建活动
export function createActivity(data: CreateActivityRequest) {
  return request<Activity>({
    url: '/api/v1/activities',
    method: 'post',
    data,
  });
}

// 更新活动
export function updateActivity(id: number, data: UpdateActivityRequest) {
  return request<Activity>({
    url: `/api/v1/activities/${id}`,
    method: 'put',
    data,
  });
}

// 删除活动
export function deleteActivity(id: number) {
  return request({
    url: `/api/v1/activities/${id}`,
    method: 'delete',
  });
}

// 更新活动状态
export function updateActivityStatus(id: number, data: UpdateActivityStatusRequest) {
  return request({
    url: `/api/v1/activities/${id}/status`,
    method: 'patch',
    data,
  });
}

// 获取活动参与记录
export function getActivityParticipants(id: number, params: ActivityParticipantQueryParams) {
  return request<ActivityParticipantListResponse>({
    url: `/api/v1/activities/${id}/participants`,
    method: 'get',
    params,
  });
} 