/**
 * Application API
 *
 * 应用管理相关 API 封装
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
 * Application 数据模型
 */
export interface Application {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  AppId?: string;
  AppSecret?: string;
  Name: string;
  IP?: string;
  Description?: string;
  Status?: 'Normal' | 'Frozen';
}

/**
 * 创建 Application 请求参数
 */
export interface CreateApplicationRequest {
  Name: string;
  IP?: string;
  Description?: string;
  Status?: 'Normal' | 'Frozen';
}

/**
 * 更新 Application 请求参数
 */
export interface UpdateApplicationRequest {
  Id: number;
  Name: string;
  IP?: string;
  Description?: string;
  Status?: 'Normal' | 'Frozen';
}

/**
 * 批量更新状态请求参数
 */
export interface BatchUpdateApplicationStatusRequest {
  ids: number[];
  status: 'Normal' | 'Frozen';
}

/**
 * 批量删除请求参数
 */
export interface BatchDeleteApplicationRequest {
  ids: number[];
}

/**
 * 查询参数
 */
export interface GetApplicationListParams {
  page?: number;
  pageSize?: number;
  showType?: string;
  startDate?: string;
  endDate?: string;
  appId?: string;
}

/**
 * 查询单个 Application 参数
 */
export interface GetApplicationParams {
  id: number;
}

/**
 * 创建 Application
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Application>>>
 *
 * @example
 * ```typescript
 * const response = await createApplication({
 *   Name: 'My App',
 *   IP: '192.168.1.1',
 *   Description: 'Test application',
 *   Status: 'Normal'
 * });
 * ```
 */
export const createApplication = (
  data: CreateApplicationRequest
): Promise<AxiosResponse<ApiResponse<Application>>> => {
  return service.post('/admin/applications', data);
};

/**
 * 删除 Application
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteApplication = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/applications/${data.id}`);
};

/**
 * 批量删除 Application
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteApplication = (
  data: BatchDeleteApplicationRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/applications/batch', { data });
};

/**
 * 更新 Application
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApplication = (
  data: UpdateApplicationRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/applications/${data.Id}`, data);
};

/**
 * 批量更新 Application 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApplicationStatus = (
  data: BatchUpdateApplicationStatusRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/applications/batch', data);
};

/**
 * 根据 ID 查询 Application
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Application>>>
 */
export const findApplication = (
  params: GetApplicationParams
): Promise<AxiosResponse<ApiResponse<Application>>> => {
  return service.get(`/admin/applications/${params.id}`);
};

/**
 * 分页获取 Application 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Application[]>>>
 */
export const getApplicationList = (
  params?: GetApplicationListParams
): Promise<AxiosResponse<ApiResponse<Application[]>>> => {
  return service.get('/admin/applications', { params });
};
