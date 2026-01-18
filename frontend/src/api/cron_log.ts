/**
 * Cron Log API
 *
 * 定时任务日志管理相关 API 封装
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
 * CronLog 数据模型
 */
export interface CronLog {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  Method: string;
  Param?: string;
  ErrMsg?: string;
  StartTime?: string;
  EndTime?: string;
  ExecTime?: number;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetCronLogListParams {
  page?: number;
  pageSize?: number;
  method?: string;
  status?: string;
  startDate?: string;
  endDate?: string;
}

/**
 * 查询单个 CronLog 参数
 */
export interface GetCronLogParams {
  id: number;
}

/**
 * 根据 ID 查询 CronLog
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<CronLog>>>
 */
export const findCronLog = (
  params: GetCronLogParams
): Promise<AxiosResponse<ApiResponse<CronLog>>> => {
  return service.get(`/admin/cron_logs/${params.id}`);
};

/**
 * 分页获取 CronLog 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<CronLog[]>>>
 */
export const getCronLogList = (
  params?: GetCronLogListParams
): Promise<AxiosResponse<ApiResponse<CronLog[]>>> => {
  return service.get('/admin/cron_logs', { params });
};

/**
 * 删除 CronLog
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteCronLog = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/admin/cron_logs/${data.id}`);
};

/**
 * 批量删除 CronLog
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteCronLog = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete('/admin/cron_logs/batch', { data });
};
