/**
 * Approval Task API
 *
 * 审批任务管理相关 API 封装
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
 * ApprovalTask 数据模型
 */
export interface ApprovalTask {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  TaskCode?: string;
  ApprovalCode: string;
  NodeCode?: string;
  NodeName?: string;
  NodeType?: string;
  ApproverType?: string;
  ApproverConfig?: string;
  Urgency?: string;
  Comment?: string;
  RemindCount?: number;
  Attachments?: string;
  AssigneeID?: string;
  AssigneeName?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 创建 ApprovalTask 请求参数
 */
export interface CreateApprovalTaskRequest {
  ApprovalCode: string;
  NodeCode?: string;
  NodeName?: string;
  NodeType?: string;
  ApproverType?: string;
  ApproverConfig?: string;
  Urgency?: string;
  Comment?: string;
  Attachments?: string;
  AssigneeID?: string;
  AssigneeName?: string;
  Status?: string;
}

/**
 * 更新 ApprovalTask 请求参数
 */
export interface UpdateApprovalTaskRequest {
  id: number;
  ApprovalCode?: string;
  NodeCode?: string;
  NodeName?: string;
  NodeType?: string;
  ApproverType?: string;
  ApproverConfig?: string;
  Urgency?: string;
  Comment?: string;
  RemindCount?: number;
  Attachments?: string;
  AssigneeID?: string;
  AssigneeName?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetApprovalTaskListParams {
  page?: number;
  pageSize?: number;
  approvalCode?: string;
  assigneeName?: string;
  status?: string;
  urgency?: string;
}

/**
 * 查询单个 ApprovalTask 参数
 */
export interface GetApprovalTaskParams {
  id: number;
}

/**
 * 创建 ApprovalTask
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalTask>>>
 */
export const createApprovalTask = (
  data: CreateApprovalTaskRequest
): Promise<AxiosResponse<ApiResponse<ApprovalTask>>> => {
  return service.post('/approval_tasks', data);
};

/**
 * 删除 ApprovalTask
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteApprovalTask = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/approval_tasks/${data.id}`);
};

/**
 * 批量删除 ApprovalTask
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteApprovalTask = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/approval_tasks/batch_delete', data);
};

/**
 * 更新 ApprovalTask
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApprovalTask = (
  data: UpdateApprovalTaskRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/approval_tasks/${data.id}`, data);
};

/**
 * 批量更新 ApprovalTask 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApprovalTaskStatus = (
  data: { ids: number[]; status: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/approval_tasks/batch', data);
};

/**
 * 根据 ID 查询 ApprovalTask
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalTask>>>
 */
export const findApprovalTask = (
  params: GetApprovalTaskParams
): Promise<AxiosResponse<ApiResponse<ApprovalTask>>> => {
  return service.get(`/approval_tasks/${params.id}`);
};

/**
 * 分页获取 ApprovalTask 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalTask[]>>>
 */
export const getApprovalTaskList = (
  params?: GetApprovalTaskListParams
): Promise<AxiosResponse<ApiResponse<ApprovalTask[]>>> => {
  return service.get('/approval_tasks', { params });
};
