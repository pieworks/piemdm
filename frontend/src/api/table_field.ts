/**
 * Table Field API
 *
 * 表字段管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * TableField 数据模型
 */
export interface TableField {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  TableCode: string;
  FieldName: string;
  FieldType: string;
  FieldLabel?: string;
  DefaultValue?: string;
  IsRequired?: boolean;
  IsUnique?: boolean;
  MaxLength?: number;
  MinLength?: number;
  Pattern?: string;
  Options?: string;
  Sort?: number;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 字段选项
 */
export interface FieldOption {
  code: string;
  name: string;
}

/**
 * 创建表字段
 */
export const createTableField = (
  data: Partial<TableField>
): Promise<AxiosResponse<ApiResponse<TableField>>> => {
  return service.post('/table_fields', data);
};

/**
 * 删除表字段
 */
export const deleteTableField = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/table_fields/${data.id}`);
};

/**
 * 批量删除表字段
 */
export const batchDeleteTableField = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/table_fields/batch_delete', data);
};

/**
 * 更新表字段
 */
export const updateTableField = (
  data: Partial<TableField> & { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/table_fields/${data.id}`, data);
};

/**
 * 批量更新表字段状态
 */
export const updateTableFieldStatus = (
  data: { ids: number[]; status: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/table_fields/batch', data);
};

/**
 * 根据 ID 查询表字段
 */
export const findTableField = (
  params: { id: number }
): Promise<AxiosResponse<ApiResponse<TableField>>> => {
  return service.get(`/table_fields/${params.id}`);
};

/**
 * 获取表字段列表
 */
export const findTableFieldList = (
  params?: any
): Promise<AxiosResponse<ApiResponse<TableField[]>>> => {
  return service.get('/table_fields', { params });
};

/**
 * 分页获取表字段列表
 */
export const getTableFieldList = (
  params?: any
): Promise<AxiosResponse<ApiResponse<TableField[]>>> => {
  return service.get('/table_fields', { params });
};

/**
 * 获取表字段详情
 */
export const getTableField = (
  params: { id: number }
): Promise<AxiosResponse<ApiResponse<TableField>>> => {
  return service.get(`/table_fields/${params.id}`);
};

/**
 * 发布实体表
 */
export const publicTable = (
  data: { tableCode: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/table_fields/public', data);
};

/**
 * 获取表列表
 */
export const findTableList = (
  params?: any
): Promise<AxiosResponse<ApiResponse<any[]>>> => {
  return service.get('/tables', { params });
};

/**
 * 获取表的所有字段(包括系统字段)
 * @param params - 查询参数
 * @param params.table_code - 表代码
 * @returns 返回包含用户字段和系统字段的完整列表
 */
export const getTableFields = (
  params: { table_code: string }
): Promise<AxiosResponse<ApiResponse<TableField[]>>> => {
  return service.get('/table_fields/fields', { params });
};

/**
 * 获取字段类型预设列表
 * @returns 返回所有字段类型预设配置
 */
export const getFieldTypePresets = (): Promise<
  AxiosResponse<ApiResponse<any[]>>
> => {
  return service.get('/field-type-presets');
};

/**
 * 获取字段类型分组
 * @returns 返回字段类型分组信息
 */
export const getFieldTypeGroups = (): Promise<
  AxiosResponse<ApiResponse<any[]>>
> => {
  return service.get('/field-type-groups');
};

/**
 * 获取表的选项列表（用于关联字段下拉）
 * @param tableCode - 表代码
 * @param filter - 过滤条件 (可选), 例如 {"dict_code": "DIC0045"}
 * @returns 返回 [{code, name}, ...] 格式的选项列表
 */
export const getTableOptions = (
  tableCode: string,
  filter: Record<string, any> | null = null
): Promise<AxiosResponse<ApiResponse<FieldOption[]>>> => {
  const params: any = {};
  if (filter) {
    params.filter = JSON.stringify(filter);
  }

  return service.get(`/table/${tableCode}/options`, { params });
};
