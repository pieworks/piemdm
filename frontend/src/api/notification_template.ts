/**
 * NotificationTemplate API
 *
 * 通知模板管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

const service = requestService as AxiosInstance;

/**
 * NotificationTemplate 数据模型
 */
export interface NotificationTemplate {
  ID?: string;
  TemplateCode: string;
  TemplateName: string;
  TemplateType: string;
  NotificationType: string;
  TitleTemplate: string;
  ContentTemplate: string;
  Variables?: string; // JSON string

  Status?: string;
  Description?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

/**
 * 创建请求参数
 */
export interface CreateNotificationTemplateRequest {
  TemplateCode: string;
  TemplateName: string;
  TemplateType: string;
  NotificationType: string;
  TitleTemplate: string;
  ContentTemplate: string;
  Variables?: string;
  Status?: string;
  Description?: string;
}

/**
 * 更新请求参数
 */
export interface UpdateNotificationTemplateRequest {
  ID: string;
  TemplateCode?: string;
  TemplateName?: string;
  TemplateType?: string;
  NotificationType?: string;
  TitleTemplate?: string;
  ContentTemplate?: string;
  Variables?: string;
  Status?: string;
  Description?: string;
}

/**
 * 查询列表参数
 */
export interface GetNotificationTemplateListParams {
  page?: number;
  pageSize?: number;
  templateCode?: string;
  templateName?: string;
  templateType?: string;
  notificationType?: string;

}

/**
 * 验证模板请求参数
 */
export interface ValidateTemplateRequest {
  TemplateCode: string;
  TitleTemplate: string;
  ContentTemplate: string;
  Variables: Record<string, any>;
}

/**
 * 渲染模板请求参数
 */
export interface RenderTemplateRequest {
  TemplateID: string;
  Variables: Record<string, any>;
}

/**
 * 创建通知模板
 */
export const createNotificationTemplate = (
  data: CreateNotificationTemplateRequest
): Promise<AxiosResponse<ApiResponse<NotificationTemplate>>> => {
  return service.post('/admin/notification_templates', data);
};

/**
 * 删除通知模板
 */
export const deleteNotificationTemplate = (
  id: string
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/notification_templates/${id}`);
};

/**
 * 更新通知模板
 */
export const updateNotificationTemplate = (
  data: UpdateNotificationTemplateRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/notification_templates/${data.ID}`, data);
};

/**
 * 批量更新通知模板状态
 */
export const updateNotificationTemplateStatus = (
  data: { ids: string[]; status: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/notification_templates/batch', data);
};

/**
 * 根据 ID 查询通知模板
 */
export const findNotificationTemplate = (
  id: string
): Promise<AxiosResponse<ApiResponse<NotificationTemplate>>> => {
  return service.get(`/admin/notification_templates/${id}`);
};

/**
 * 分页获取通知模板列表
 */
export const getNotificationTemplateList = (
  params?: GetNotificationTemplateListParams
): Promise<AxiosResponse<ApiResponse<NotificationTemplate[]>>> => {
  return service.get('/admin/notification_templates', { params });
};

/**
 * 渲染模板预览
 */
export const renderTemplate = (
  id: string,
  data: RenderTemplateRequest
): Promise<AxiosResponse<ApiResponse<any>>> => {
  return service.post(`/admin/notification_templates/${id}/render`, data);
};

/**
 * 验证模板格式
 */
export const validateTemplate = (
  data: ValidateTemplateRequest
): Promise<AxiosResponse<ApiResponse<any>>> => {
  return service.post('/admin/notification_templates/validate', data);
};
