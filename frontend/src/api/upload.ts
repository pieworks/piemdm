/**
 * Upload API
 *
 * 文件上传相关 API 封装
 */

import requestService from '@/utils/request';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { ApiResponse } from '@/api/types';

// 类型断言: request.js 导出的 service 是一个 Axios 实例
const service = requestService as AxiosInstance;

/**
 * 上传响应数据
 */
export interface UploadResponse {
  url: string;
  filename: string;
  size?: number;
}

/**
 * 上传文件
 *
 * @param formData - 包含文件的 FormData 对象
 * @returns Promise<AxiosResponse<ApiResponse<UploadResponse>>>
 */
export const uploadFile = (
  formData: FormData
): Promise<AxiosResponse<ApiResponse<UploadResponse>>> => {
  return service.post('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};
