/**
 * Cron API
 *
 * 定时任务管理相关 API 封装
 *
 * 注意: 由于后端未添加 Swagger 注释,这里使用手动封装的 TypeScript 版本
 * TODO: 后端添加 Swagger 注释后,可以使用自动生成的 API 客户端
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * Cron 数据模型
 */
export interface Cron {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  Code: string;
  Expression: string;
  Name: string;
  EntityCode?: string;
  System?: string;
  Url?: string;
  Protocol?: string;
  Method?: string;
  AppId?: string;
  AppKey?: string;
  SignType?: string;
  Description?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 创建 Cron 请求参数
 */
export interface CreateCronRequest {
  Code: string;
  Expression: string;
  Name: string;
  EntityCode?: string;
  System?: string;
  Url?: string;
  Protocol?: string;
  Method?: string;
  AppId?: string;
  AppKey?: string;
  SignType?: string;
  Description?: string;
  Status?: string;
}

/**
 * 更新 Cron 请求参数
 */
export interface UpdateCronRequest {
  ID: number;
  Code?: string;
  Expression?: string;
  Name?: string;
  EntityCode?: string;
  System?: string;
  Url?: string;
  Protocol?: string;
  Method?: string;
  AppId?: string;
  AppKey?: string;
  SignType?: string;
  Description?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetCronListParams {
  page?: number;
  pageSize?: number;
  code?: string;
  name?: string;
  entityCode?: string;
  status?: string;
}

/**
 * 查询单个 Cron 参数
 */
export interface GetCronParams {
  id: number;
}

/**
 * 创建 Cron
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Cron>>>
 */
export const createCron = (
  data: CreateCronRequest
): Promise<AxiosResponse<ApiResponse<Cron>>> => {
  return service.post('/admin/crons', data);
};

/**
 * 删除 Cron
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteCron = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/crons/${data.id}`, { data });
};

/**
 * 批量删除 Cron
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteCron = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/crons/batch', { data });
};

/**
 * 更新 Cron
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateCron = (
  data: UpdateCronRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/crons/${data.ID}`, data);
};

/**
 * 批量更新 Cron 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateCronStatus = (
  data: { ids: number[]; status: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/crons/batch', data);
};

/**
 * 根据 ID 查询 Cron
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Cron>>>
 */
export const findCron = (
  params: GetCronParams
): Promise<AxiosResponse<ApiResponse<Cron>>> => {
  return service.get(`/admin/crons/${params.id}`);
};

/**
 * 分页获取 Cron 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Cron[]>>>
 */
export const getCronList = (
  params?: GetCronListParams
): Promise<AxiosResponse<ApiResponse<Cron[]>>> => {
  return service.get('/admin/crons', { params });
};
