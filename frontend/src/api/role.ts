/**
 * Role API
 *
 * 角色管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * Role 数据模型
 */
export interface Role {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  Code: string;
  Name: string;
  Description?: string;
  DataScope?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
  Permissions?: Permission[];
}

/**
 * Permission 数据模型(简化版)
 */
export interface Permission {
  ID?: number;
  Code: string;
  Name: string;
  Resource?: string;
  Action?: string;
  ParentID?: number;
  Description?: string;
}

/**
 * 创建 Role 请求参数
 */
export interface CreateRoleRequest {
  Code: string;
  Name: string;
  Description?: string;
  DataScope?: string;
  Status?: string;
}

/**
 * 更新 Role 请求参数
 */
export interface UpdateRoleRequest {
  id: number;
  Code?: string;
  Name?: string;
  Description?: string;
  DataScope?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetRoleListParams {
  page?: number;
  pageSize?: number;
  code?: string;
  name?: string;
  status?: string;
}

/**
 * 查询单个 Role 参数
 */
export interface GetRoleParams {
  id: number;
}

/**
 * 分配权限请求参数
 */
export interface AssignPermissionsRequest {
  permission_ids: number[];
}

/**
 * 创建 Role
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Role>>>
 */
export const createRole = (
  data: CreateRoleRequest
): Promise<AxiosResponse<ApiResponse<Role>>> => {
  return service.post('/admin/roles', data);
};

/**
 * 删除 Role
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteRole = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/roles/${data.id}`);
};

/**
 * 批量删除 Role
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteRole = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/roles/batch', { data });
};

/**
 * 更新 Role
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateRole = (
  data: UpdateRoleRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/roles/${data.id}`, data);
};

/**
 * 批量更新 Role 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateRoleStatus = (
  data: { ids: number[]; status: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/roles/batch', data);
};

/**
 * 根据 ID 查询 Role
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Role>>>
 */
export const findRole = (
  params: GetRoleParams
): Promise<AxiosResponse<ApiResponse<Role>>> => {
  return service.get(`/admin/roles/${params.id}`);
};

/**
 * 分页获取 Role 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Role[]>>>
 */
export const getRoleList = (
  params?: GetRoleListParams
): Promise<AxiosResponse<ApiResponse<Role[]>>> => {
  return service.get('/admin/roles', { params });
};

/**
 * 获取角色的权限列表
 *
 * @param id - 角色ID
 * @returns Promise<AxiosResponse<ApiResponse<Permission[]>>>
 */
export const getRolePermissions = (
  id: number
): Promise<AxiosResponse<ApiResponse<Permission[]>>> => {
  return service.get(`/admin/roles/${id}/permissions`);
};

/**
 * 为角色分配权限
 *
 * @param id - 角色ID
 * @param data - 权限ID列表
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const assignPermissions = (
  id: number,
  data: AssignPermissionsRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post(`/admin/roles/${id}/permissions`, data);
};
