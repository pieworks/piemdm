/**
 * Entity API
 *
 * 动态实体数据管理相关 API 封装
 *
 * 注意: 实体数据结构是动态的,由 table_fields 表定义
 * 使用泛型支持类型安全,同时保持灵活性
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * 实体基础字段 (所有动态表共有)
 */
export interface EntityBase {
  id?: number | string;
  table_code: string;
  created_at?: string;
  updated_at?: string;
  created_by?: string;
  updated_by?: string;
  status?: string;
}

/**
 * 动态实体数据
 * @template T - 动态字段类型,默认为 Record<string, any>
 */
export type EntityData<T = Record<string, any>> = EntityBase & T;

/**
 * 实体查询参数
 */
export interface EntityQueryParams {
  table_code: string;
  id?: number | string;
  page?: number;
  pageSize?: number;
  [key: string]: any; // 支持动态查询条件
}

/**
 * Batch operations参数
 */
export interface EntityBatchParams {
  table_code: string;
  ids?: (number | string)[];
  status?: string;
  [key: string]: any;
}

/**
 * 实体日志查询参数
 */
export interface EntityLogParams {
  table_code: string;
  page?: number;
  pageSize?: number;
  [key: string]: any;
}

/**
 * 实体历史查询参数
 */
export interface EntityHistoryParams {
  table_code: string;
  page?: number;
  pageSize?: number;
  [key: string]: any;
}

/**
 * 导出参数
 */
export interface EntityExportParams {
  table_code: string;
  [key: string]: any;
}

/**
 * 模板下载参数
 */
export interface EntityTemplateParams {
  table_code: string;
}

/**
 * 创建实体
 *
 * @template T - 动态字段类型
 * @param data - 实体数据 (包含 table_code 和动态字段)
 * @returns Promise<AxiosResponse<ApiResponse<EntityData<T>>>>
 */
export const createEntity = <T = Record<string, any>>(
  data: EntityData<T>
): Promise<AxiosResponse<ApiResponse<EntityData<T>>>> => {
  return service.post(`/entities/${data.table_code}`, data);
};

/**
 * 删除实体
 *
 * @param data - 包含 table_code 和 id
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteEntity = (
  data: { table_code: string; id: number | string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/entities/${data.table_code}/${data.id}`);
};

/**
 * 批量删除实体
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteEntity = (
  data: EntityBatchParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post(`/entities/${data.table_code}/batch`, data);
};

/**
 * 更新实体
 *
 * @template T - 动态字段类型
 * @param data - 实体数据 (必须包含 id)
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateEntity = <T = Record<string, any>>(
  data: EntityData<T> & { id: number | string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/entities/${data.table_code}/${data.id}`, data);
};

/**
 * 批量更新实体状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateEntityStatus = (
  data: EntityBatchParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/entities/${data.table_code}/batch`, data);
};

/**
 * 根据 ID 查询实体
 *
 * @template T - 动态字段类型
 * @param params - 查询参数 (table_code 和 id)
 * @returns Promise<AxiosResponse<ApiResponse<EntityData<T>>>>
 */
export const findEntity = <T = Record<string, any>>(
  params: { table_code: string; id: number | string }
): Promise<AxiosResponse<ApiResponse<EntityData<T>>>> => {
  return service.get(`/entities/${params.table_code}/${params.id}`);
};

/**
 * 查询实体列表
 *
 * @template T - 动态字段类型
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<EntityData<T>[]>>>
 */
export const findEntityList = <T = Record<string, any>>(
  params: EntityQueryParams
): Promise<AxiosResponse<ApiResponse<EntityData<T>[]>>> => {
  return service.get(`/entities/${params.table_code}`, { params });
};

/**
 * 分页获取实体列表
 *
 * @template T - 动态字段类型
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<EntityData<T>[]>>>
 */
export const getEntityList = <T = Record<string, any>>(
  params: EntityQueryParams
): Promise<AxiosResponse<ApiResponse<EntityData<T>[]>>> => {
  return service.get(`/entities/${params.table_code}`, { params });
};

/**
 * 获取实体日志列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getEntityLogList = (
  params: EntityLogParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.get(`/entities/${params.table_code}/logs`, { params });
};

/**
 * 获取实体历史列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getEntityHistoryList = (
  params: EntityHistoryParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.get(`/entities/${params.table_code}/histories`, { params });
};

/**
 * 导出实体数据
 *
 * @param params - 导出参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getExportFile = (
  params: EntityExportParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.get(`/entities/${params.table_code}/export`, { params });
};

/**
 * 导入实体数据
 *
 * @param data - FormData 包含文件和 table_code
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const importFile = (
  data: FormData
): Promise<AxiosResponse<ApiResponse>> => {
  const tableCode = data.get('table_code');
  console.log('importFile data: ', tableCode);
  return service.post(`/entities/${tableCode}/import`, data);
};

/**
 * 获取导入模板
 *
 * @param params - 模板参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getTemplate = (
  params: EntityTemplateParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.get(`/entities/${params.table_code}/template`, { params });
};

/**
 * 获取实体统计信息
 *
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getEntityStats = (): Promise<AxiosResponse<ApiResponse>> => {
  return service.get('/entities/statistics');
};
