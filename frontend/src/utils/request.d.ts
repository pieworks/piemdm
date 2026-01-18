/**
 * Type declarations for request.js
 *
 * 为 request.js 提供 TypeScript 类型声明
 */

import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

/**
 * 扩展的 Axios 请求配置
 *
 * 添加了自定义的全局 loading 和错误处理控制
 */
export interface ExtendedAxiosRequestConfig extends AxiosRequestConfig {
  /**
   * 是否显示全局 loading 动画
   * @default true
   */
  globalLoading?: boolean;

  /**
   * 是否显示全局错误提示
   * @default true
   */
  globalError?: boolean;
}

/**
 * 扩展的 Axios 实例
 *
 * 包含自定义的拦截器和配置
 */
export interface ExtendedAxiosInstance extends AxiosInstance {
  request<T = any, R = AxiosResponse<T>>(config: ExtendedAxiosRequestConfig): Promise<R>;
  get<T = any, R = AxiosResponse<T>>(url: string, config?: ExtendedAxiosRequestConfig): Promise<R>;
  delete<T = any, R = AxiosResponse<T>>(url: string, config?: ExtendedAxiosRequestConfig): Promise<R>;
  head<T = any, R = AxiosResponse<T>>(url: string, config?: ExtendedAxiosRequestConfig): Promise<R>;
  options<T = any, R = AxiosResponse<T>>(url: string, config?: ExtendedAxiosRequestConfig): Promise<R>;
  post<T = any, R = AxiosResponse<T>>(url: string, data?: any, config?: ExtendedAxiosRequestConfig): Promise<R>;
  put<T = any, R = AxiosResponse<T>>(url: string, data?: any, config?: ExtendedAxiosRequestConfig): Promise<R>;
  patch<T = any, R = AxiosResponse<T>>(url: string, data?: any, config?: ExtendedAxiosRequestConfig): Promise<R>;
}

/**
 * 默认导出的 service 实例
 *
 * 配置了以下拦截器功能:
 * - 自动添加 Authorization token
 * - 自动添加 X-User-Id 请求头
 * - 全局 loading 控制
 * - 全局错误处理
 * - 401/403 自动跳转登录
 * - new-token 自动更新
 */
declare const service: ExtendedAxiosInstance;

export default service;

/**
 * 登录 API
 *
 * @param credentials - 登录凭证
 * @param endpoint - API 端点
 * @param config - 请求配置
 */
export function loginApi(
  credentials: any,
  endpoint: string,
  config?: ExtendedAxiosRequestConfig
): Promise<AxiosResponse>;

/**
 * 验证 Token API
 *
 * @param config - 请求配置
 */
export function validateTokenApi(
  config?: ExtendedAxiosRequestConfig
): Promise<AxiosResponse>;

/**
 * 登出 API
 *
 * @param config - 请求配置
 */
export function logoutApi(
  config?: ExtendedAxiosRequestConfig
): Promise<AxiosResponse>;
