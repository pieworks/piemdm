/**
 * Table API
 *
 * 表管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse, BatchDeleteRequest, BatchUpdateStatusRequest } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * Table 数据模型
 */
export interface Table {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  Code: string;
  Name: string;
  DisplayMode?: string;
  TableType?: string;
  ParentTable?: string;
  ParentField?: string;
  SelfField?: string;
  Sort?: number;
  Description?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 创建 Table 请求参数
 */
export interface CreateTableRequest {
  Code: string;
  Name: string;
  DisplayMode?: string;
  TableType?: string;
  ParentTable?: string;
  ParentField?: string;
  SelfField?: string;
  TreeTable?: string;
  Sort?: number;
  Description?: string;
  Status?: string;
}

/**
 * 更新 Table 请求参数
 */
export interface UpdateTableRequest {
  id: number;
  Code?: string;
  Name?: string;
  DisplayMode?: string;
  TableType?: string;
  ParentTable?: string;
  ParentField?: string;
  SelfField?: string;
  TreeTable?: string;
  Sort?: number;
  Description?: string;
  Status?: string;
}

/**
 * 查询 Table 列表参数
 */
export interface GetTableListParams {
  page?: number;
  pageSize?: number;
  showType?: string;
  startDate?: string;
  endDate?: string;
  code?: string;
  name?: string;
  status?: string;
  relation?: string; // 关联关系: Entity | Item (已废弃，使用 table_type)
  table_type?: string; // 表类型: Entity | Item
}

/**
 * 查询单个 Table 参数
 */
export interface GetTableParams {
  id: number;
}

/**
 * 创建 Table
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Table>>>
 *
 * @example
 * ```typescript
 * const response = await createTable({
 *   Code: 'user_table',
 *   Name: '用户表',
 *   Description: '用户信息表',
 *   Status: 'Active'
 * });
 * ```
 */
export const createTable = (
  data: CreateTableRequest
): Promise<AxiosResponse<ApiResponse<Table>>> => {
  return service.post('/tables', data);
};

/**
 * 删除 Table
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteTable = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/tables/${data.id}`);
};

/**
 * 批量删除 Table
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteTable = (
  data: BatchDeleteRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/tables/batch_delete', data);
};

/**
 * 更新 Table
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateTable = (
  data: UpdateTableRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/tables/${data.id}`, data);
};

/**
 * 批量更新 Table 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateTableStatus = (
  data: BatchUpdateStatusRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/tables/batch', data);
};

/**
 * 根据 ID 查询 Table
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Table>>>
 */
export const getTable = (
  params: GetTableParams
): Promise<AxiosResponse<ApiResponse<Table>>> => {
  return service.get(`/tables/${params.id}`);
};

/**
 * 分页获取 Table 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Table[]>>>
 */
export const getTableList = (
  params?: GetTableListParams
): Promise<AxiosResponse<ApiResponse<Table[]>>> => {
  return service.get('/tables', { params });
};

/**
 * 查询 Table 列表 (别名)
 *
 * @deprecated 使用 getTableList 代替
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Table[]>>>
 */
export const findTableList = (
  params?: GetTableListParams
): Promise<AxiosResponse<ApiResponse<Table[]>>> => {
  return service.get('/tables', { params });
};
