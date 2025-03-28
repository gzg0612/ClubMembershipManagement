import { get, post, put } from '@/utils/request';
import type { SystemSettings, UpdateSettingsRequest, SettingsResponse } from '@/types/settings';

// 获取系统设置
export const getSettings = () => {
  return get<SettingsResponse>('/api/v1/settings');
};

// 更新系统设置
export const updateSettings = (data: UpdateSettingsRequest) => {
  return put<SettingsResponse>('/api/v1/settings', data);
};

// 重置系统设置
export const resetSettings = () => {
  return post<SettingsResponse>('/api/v1/settings/reset');
};

// 测试邮件设置
export const testEmailSettings = (email: string) => {
  return post<{ code: number; message: string }>('/api/v1/settings/test-email', { email });
};

// 测试短信设置
export const testSMSSettings = (phone: string) => {
  return post<{ code: number; message: string }>('/api/v1/settings/test-sms', { phone });
}; 