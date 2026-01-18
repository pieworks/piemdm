/**
 * User API
 *
 * 用户管理相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse, BatchDeleteRequest, BatchUpdateStatusRequest } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * User 数据模型
 */
export interface User {
  ID?: number;
  Username: string;
  Password?: string;
  Email?: string;
  Phone?: string;
  Status?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

/**
 * 创建 User 请求参数
 */
export interface CreateUserRequest {
  Username: string;
  Password: string;
  Email?: string;
  Phone?: string;
  Status?: string;
}

/**
 * 更新 User 请求参数
 */
export interface UpdateUserRequest {
  ID: number;
  Username?: string;
  Email?: string;
  Phone?: string;
  Status?: string;
}

/**
 * 查询 User 列表参数
 */
export interface GetUserListParams {
  page?: number;
  pageSize?: number;
  showType?: string;
  startDate?: string;
  endDate?: string;
  username?: string;
  email?: string;
  status?: string;
}

/**
 * 查询单个 User 参数
 */
export interface GetUserParams {
  id: number;
}

/**
 * 登录请求参数
 */
export interface LoginRequest {
  username: string;
  password: string;
}

/**
 * 修改密码请求参数
 */
export interface ChangePasswordRequest {
  username: string;
  password: string;
  newPassword: string;
}

/**
 * 创建 User
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<User>>>
 */
export const createUser = (
  data: CreateUserRequest
): Promise<AxiosResponse<ApiResponse<User>>> => {
  return service.post('/admin/users', data);
};

/**
 * 删除 User
 *
 * @param data - 包含 ID 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteUser = (
  data: { ID: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/users/${data.ID}`);
};

/**
 * 批量删除 User
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteUser = (
  data: BatchDeleteRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/users/batch', { data });
};

/**
 * 更新 User
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateUser = (
  data: UpdateUserRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/admin/users/${data.ID}`, data);
};

/**
 * 批量更新 User 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateUserStatus = (
  data: BatchUpdateStatusRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/users/batch', data);
};

/**
 * 根据 ID 查询 User
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<User>>>
 */
export const findUser = (
  params: GetUserParams
): Promise<AxiosResponse<ApiResponse<User>>> => {
  return service.get(`/admin/users/${params.id}`, { params });
};

/**
 * 用户登录
 *
 * @param data - 登录参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const login = (
  data: LoginRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/auth/login', data);
};

/**
 * 获取验证码
 *
 * @param data - 验证码请求参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const captcha = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/base/captcha', data);
};

/**
 * 用户注册
 *
 * @param data - 注册参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const register = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/auth/register', data);
};

/**
 * 修改密码
 *
 * @param data - 修改密码参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const changePassword = (
  data: ChangePasswordRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/users/changePassword', data);
};

/**
 * 分页获取 User 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<User[]>>>
 */
export const getUserList = (
  params?: GetUserListParams
): Promise<AxiosResponse<ApiResponse<User[]>>> => {
  return service.get('/admin/users', { params });
};

/**
 * 根据 ID 查询 User
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<User>>>
 */
export const getUser = (
  params: GetUserParams
): Promise<AxiosResponse<ApiResponse<User>>> => {
  return service.get(`/admin/users/${params.id}`);
};

/**
 * 设置用户权限
 *
 * @param data - 权限设置参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const setUserAuthority = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/users/setAuthority', data);
};

/**
 * 设置用户信息
 *
 * @param data - 用户信息参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const setUserInfo = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/users/setInfo', data);
};

/**
 * 设置当前用户信息
 *
 * @param data - 用户信息参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const setSelfInfo = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/admin/users/setSelfInfo', data);
};

/**
 * 设置用户多个权限
 *
 * @param data - 权限设置参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const setUserAuthorities = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/users/setUserAuthorities', data);
};

/**
 * 获取当前用户信息
 *
 * @returns Promise<AxiosResponse<ApiResponse<User>>>
 */
export const getUserInfo = (): Promise<AxiosResponse<ApiResponse<User>>> => {
  return service.get('/admin/users');
};

/**
 * 重置密码
 *
 * @param data - 重置密码参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const resetPassword = (
  data: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/admin/users/resetPassword', data);
};
