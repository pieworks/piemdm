/**
 * Webhook Delivery API
 *
 * Webhook 投递记录管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse, BatchDeleteRequest, BatchUpdateStatusRequest } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * WebhookDelivery 数据模型
 */
export interface WebhookDelivery {
  ID?: number;
  WebhookID: number;
  Event: string;
  Payload?: string;
  Status?: string;
  ResponseCode?: number;
  ResponseBody?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

/**
 * 创建 WebhookDelivery 请求参数
 */
export interface CreateWebhookDeliveryRequest {
  WebhookID: number;
  Event: string;
  Payload?: string;
  Status?: string;
}

/**
 * 更新 WebhookDelivery 请求参数
 */
export interface UpdateWebhookDeliveryRequest {
  id: number;
  WebhookID?: number;
  Event?: string;
  Payload?: string;
  Status?: string;
  ResponseCode?: number;
  ResponseBody?: string;
}

/**
 * 查询 WebhookDelivery 列表参数
 */
export interface GetWebhookDeliveryListParams {
  page?: number;
  pageSize?: number;
  showType?: string;
  startDate?: string;
  endDate?: string;
  webhookId?: number;
  event?: string;
  status?: string;
}

/**
 * 查询单个 WebhookDelivery 参数
 */
export interface GetWebhookDeliveryParams {
  id: number;
}

/**
 * 创建 WebhookDelivery
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<WebhookDelivery>>>
 */
export const createWebhookDelivery = (
  data: CreateWebhookDeliveryRequest
): Promise<AxiosResponse<ApiResponse<WebhookDelivery>>> => {
  return service.post('/webhook_deliveries', data);
};

/**
 * 删除 WebhookDelivery
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteWebhookDelivery = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/webhook_deliveries/${data.id}`, { data });
};

/**
 * 批量删除 WebhookDelivery
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteWebhookDelivery = (
  data: BatchDeleteRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/webhook_deliveries/batch', { data });
};

/**
 * 更新 WebhookDelivery
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateWebhookDelivery = (
  data: UpdateWebhookDeliveryRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/webhook_deliveries/${data.id}`, data);
};

/**
 * 批量更新 WebhookDelivery 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateWebhookDeliveryStatus = (
  data: BatchUpdateStatusRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/webhook_deliveries/batch', data);
};

/**
 * 根据 ID 查询 WebhookDelivery
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<WebhookDelivery>>>
 */
export const findWebhookDelivery = (
  params: GetWebhookDeliveryParams
): Promise<AxiosResponse<ApiResponse<WebhookDelivery>>> => {
  return service.get(`/webhook_deliveries/${params.id}`);
};

/**
 * 分页获取 WebhookDelivery 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<WebhookDelivery[]>>>
 */
export const getWebhookDeliveryList = (
  params?: GetWebhookDeliveryListParams
): Promise<AxiosResponse<ApiResponse<WebhookDelivery[]>>> => {
  return service.get('/webhook_deliveries', { params });
};

/**
 * 根据 ID 查询 WebhookDelivery (别名)
 *
 * @deprecated 使用 findWebhookDelivery 代替
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<WebhookDelivery>>>
 */
export const getWebhookDelivery = (
  params: GetWebhookDeliveryParams
): Promise<AxiosResponse<ApiResponse<WebhookDelivery>>> => {
  return service.get(`/webhook_deliveries/${params.id}`);
};
