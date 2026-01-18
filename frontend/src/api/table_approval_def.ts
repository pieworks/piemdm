/**
 * Table Approval Definition API
 *
 * 表审批定义管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * TableApprovalDefinition 数据模型
 */
export interface TableApprovalDefinition {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  EntityCode: string;
  Operation: string;
  ApprovalDefCode: string;
  Description?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 创建 TableApprovalDefinition 请求参数
 */
export interface CreateTableApprovalDefRequest {
  EntityCode: string;
  Operation: string;
  ApprovalDefCode: string;
  Description?: string;
  Status?: string;
}

/**
 * 更新 TableApprovalDefinition 请求参数
 */
export interface UpdateTableApprovalDefRequest {
  id: number;
  EntityCode?: string;
  Operation?: string;
  ApprovalDefCode?: string;
  Description?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetTableApprovalDefListParams {
  entity_code: string;
  operation?: string;
}

/**
 * 获取表审批定义列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<TableApprovalDefinition[]>>>
 */
export const findTableApprovalDefList = (
  params: GetTableApprovalDefListParams
): Promise<AxiosResponse<ApiResponse<TableApprovalDefinition[]>>> => {
  return service.get('/admin/table_approval_defs', { params });
};

/**
 * 根据 ID 查询表审批定义
 *
 * @param id - 表审批定义ID
 * @returns Promise<AxiosResponse<ApiResponse<TableApprovalDefinition>>>
 */
export const getTableApprovalDef = (
  id: number
): Promise<AxiosResponse<ApiResponse<TableApprovalDefinition>>> => {
  return service.get(`/admin/table_approval_defs/${id}`);
};

/**
 * 创建表审批定义
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<TableApprovalDefinition>>>
 */
export const createTableApprovalDef = (
  data: CreateTableApprovalDefRequest
): Promise<AxiosResponse<ApiResponse<TableApprovalDefinition>>> => {
  return service.post('/admin/table_approval_defs', data);
};

/**
 * 更新表审批定义
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse<TableApprovalDefinition>>>
 */
export const updateTableApprovalDef = (
  data: UpdateTableApprovalDefRequest
): Promise<AxiosResponse<ApiResponse<TableApprovalDefinition>>> => {
  return service.put(`/admin/table_approval_defs/${data.id}`, data);
};

/**
 * 删除表审批定义
 *
 * @param id - 表审批定义ID
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteTableApprovalDef = (
  id: number
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/table_approval_defs/${id}`);
};

/**
 * 批量创建表审批定义
 *
 * @param data - 批量创建参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchCreateTableApprovalDef = (
  data: { list: TableApprovalDefinition[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/table_approval_defs/batch_create', data);
};

/**
 * 批量删除表审批定义
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteTableApprovalDef = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/table_approval_defs/batch_delete', data);
};
