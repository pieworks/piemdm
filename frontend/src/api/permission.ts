/**
 * Permission API
 *
 * 权限管理相关 API 封装
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
 * Permission 数据模型
 */
export interface Permission {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  Code: string;
  Name: string;
  Resource?: string;
  Action?: string;
  ParentID?: number;
  Description?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
  Children?: Permission[];
}

/**
 * 创建 Permission 请求参数
 */
export interface CreatePermissionRequest {
  Code: string;
  Name: string;
  Resource?: string;
  Action?: string;
  ParentID?: number;
  Description?: string;
  Status?: string;
}

/**
 * 更新 Permission 请求参数
 */
export interface UpdatePermissionRequest {
  Code?: string;
  Name?: string;
  Resource?: string;
  Action?: string;
  ParentID?: number;
  Description?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetPermissionListParams {
  page?: number;
  pageSize?: number;
  name?: string;
  code?: string;
  resource?: string;
  status?: string;
}

/**
 * 查询单个 Permission 参数
 */
export interface GetPermissionParams {
  id: number;
}


/**
 * 创建 Permission
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Permission>>>
 */
export const createPermission = (
  data: CreatePermissionRequest
): Promise<AxiosResponse<ApiResponse<Permission>>> => {
  return service.post('/admin/permissions', data);
};

/**
 * 删除 Permission
 *
 * @param id - 权限ID
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deletePermission = (
  id: number
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/permissions/${id}`);
};

/**
 * 批量删除 Permission
 *
 * @param params - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeletePermission = (
  params: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/permissions/batch_delete', { data: params });
};

/**
 * 更新 Permission
 *
 * @param id - 权限ID
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updatePermission = (
  id: number,
  data: UpdatePermissionRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/permissions/${id}`, data);
};

/**
 * 根据 ID 查询 Permission
 *
 * @param id - 权限ID
 * @returns Promise<AxiosResponse<ApiResponse<Permission>>>
 */
export const getPermission = (
  id: number
): Promise<AxiosResponse<ApiResponse<Permission>>> => {
  return service.get(`/admin/permissions/${id}`);
};

/**
 * 分页获取 Permission 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Permission[]>>>
 */
export const getPermissionList = (
  params?: GetPermissionListParams
): Promise<AxiosResponse<ApiResponse<Permission[]>>> => {
  return service.get('/admin/permissions', { params });
};

/**
 * 获取权限树
 *
 * @returns Promise<AxiosResponse<ApiResponse<Permission[]>>>
 */
export const getPermissionTree = (): Promise<
  AxiosResponse<ApiResponse<Permission[]>>
> => {
  return service.get('/admin/permissions/tree');
};
