/**
 * Webhook API
 *
 * Webhook 管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse, BatchDeleteRequest, BatchUpdateStatusRequest } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * Webhook 数据模型
 */
export interface Webhook {
  ID?: number;
  Name: string;
  Url: string;
  Events?: string[];
  Status?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

/**
 * 创建 Webhook 请求参数
 */
export interface CreateWebhookRequest {
  Name: string;
  Url: string;
  Events?: string[];
  Status?: string;
}

/**
 * 更新 Webhook 请求参数
 */
export interface UpdateWebhookRequest {
  ID: number;
  Name?: string;
  Url?: string;
  Events?: string[];
  Status?: string;
}

/**
 * 查询 Webhook 列表参数
 */
export interface GetWebhookListParams {
  page?: number;
  pageSize?: number;
  showType?: string;
  startDate?: string;
  endDate?: string;
  name?: string;
  url?: string;
  status?: string;
}

/**
 * 查询单个 Webhook 参数
 */
export interface GetWebhookParams {
  id: number;
}

/**
 * 创建 Webhook
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Webhook>>>
 */
export const createWebhook = (
  data: CreateWebhookRequest
): Promise<AxiosResponse<ApiResponse<Webhook>>> => {
  return service.post('/admin/webhooks', data);
};

/**
 * 删除 Webhook
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteWebhook = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/webhooks/${data.id}`);
};

/**
 * 批量删除 Webhook
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteWebhook = (
  data: BatchDeleteRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/webhooks/batch', { data });
};

/**
 * 更新 Webhook
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateWebhook = (
  data: UpdateWebhookRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/webhooks/${data.ID}`, data);
};

/**
 * 批量更新 Webhook 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateWebhookStatus = (
  data: BatchUpdateStatusRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/webhooks/batch', data);
};

/**
 * 根据 ID 查询 Webhook
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Webhook>>>
 */
export const findWebhook = (
  params: GetWebhookParams
): Promise<AxiosResponse<ApiResponse<Webhook>>> => {
  return service.get(`/admin/webhooks/${params.id}`);
};

/**
 * 分页获取 Webhook 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Webhook[]>>>
 */
export const getWebhookList = (
  params?: GetWebhookListParams
): Promise<AxiosResponse<ApiResponse<Webhook[]>>> => {
  return service.get('/admin/webhooks', { params });
};

/**
 * 根据 ID 查询 Webhook (别名)
 *
 * @deprecated 使用 findWebhook 代替
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Webhook>>>
 */
export const getWebhook = (
  params: GetWebhookParams
): Promise<AxiosResponse<ApiResponse<Webhook>>> => {
  return service.get(`/admin/webhooks/${params.id}`);
};
