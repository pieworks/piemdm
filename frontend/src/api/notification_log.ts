/**
 * NotificationLog API
 *
 * 通知日志管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

const service = requestService as AxiosInstance;

/**
 * NotificationLog 数据模型
 */
export interface NotificationLog {
  ID: number;
  ApprovalID: string;
  TaskID: string;
  RecipientID: string;
  RecipientType: string;
  NotificationType: string;
  TemplateID: string;
  TemplateCode: string;
  Title: string;
  Content: string;
  Status: string;
  SendTime?: string;
  ErrorMessage?: string;
  RetryCount: number;
  MaxRetryCount: number;
  NextRetryTime?: string;
  ExtraData?: string;
  CreatedAt: string;
  UpdatedAt: string;
}

/**
 * 查询列表参数
 */
export interface GetNotificationLogListParams {
  page?: number;
  pageSize?: number;
  approvalId?: string;
  taskId?: string;
  recipientId?: string;
  notificationType?: string;
  status?: string;
  startTime?: string;
  endTime?: string;
}

/**
 * 根据 ID 查询通知日志
 */
export const findNotificationLog = (
  id: number
): Promise<AxiosResponse<ApiResponse<NotificationLog>>> => {
  return service.get(`/admin/notification_logs/${id}`);
};

/**
 * 分页获取通知日志列表
 */
export const getNotificationLogList = (
  params?: GetNotificationLogListParams
): Promise<AxiosResponse<ApiResponse<NotificationLog[]>>> => {
  return service.get('/admin/notification_logs', { params });
};
