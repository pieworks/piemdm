/**
 * Approval API
 *
 * 审批实例管理相关 API 封装
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
 * Approval 数据模型
 */
export interface Approval {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  Code?: string;
  Title: string;
  ApprovalDefCode: string;
  EntityCode?: string;
  SerialNumber?: string;
  CurrentTaskID?: string;
  CurrentTaskName?: string;
  FormData?: string;
  FormSchema?: string;
  Priority?: number;
  Urgency?: string;
  Description?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 创建 Approval 请求参数
 */
export interface CreateApprovalRequest {
  Title: string;
  ApprovalDefCode: string;
  EntityCode?: string;
  SerialNumber?: string;
  FormData?: string;
  FormSchema?: string;
  Priority?: number;
  Urgency?: string;
  Description?: string;
  Status?: string;
}

/**
 * 更新 Approval 请求参数
 */
export interface UpdateApprovalRequest {
  id: number;
  Title?: string;
  ApprovalDefCode?: string;
  EntityCode?: string;
  SerialNumber?: string;
  CurrentTaskID?: string;
  CurrentTaskName?: string;
  FormData?: string;
  FormSchema?: string;
  Priority?: number;
  Urgency?: string;
  Description?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetApprovalListParams {
  page?: number;
  pageSize?: number;
  approvalDefCode?: string;
  entityCode?: string;
  status?: string;
  urgency?: string;
  createdBy?: string;
}

/**
 * 查询单个 Approval 参数
 */
export interface GetApprovalParams {
  id: number;
}

/**
 * 审批任务操作请求参数
 */
export interface ApproveTaskRequest {
  comment?: string;
  attachments?: string;
}

/**
 * 拒绝任务请求参数
 */
export interface RejectTaskRequest {
  comment: string;
  attachments?: string;
}

/**
 * 获取审批历史参数
 */
export interface GetApprovalHistoryParams {
  approval_id: number;
}

/**
 * 创建 Approval
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<Approval>>>
 */
export const createApproval = (
  data: CreateApprovalRequest
): Promise<AxiosResponse<ApiResponse<Approval>>> => {
  return service.post('/approvals', data);
};

/**
 * 删除 Approval
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteApproval = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/approvals/${data.id}`);
};

/**
 * 批量删除 Approval
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteApproval = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/approvals/batch_delete', data);
};

/**
 * 更新 Approval
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApproval = (
  data: UpdateApprovalRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/approvals/${data.id}`, data);
};

/**
 * 根据 ID 查询 Approval
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Approval>>>
 */
export const findApproval = (
  params: GetApprovalParams
): Promise<AxiosResponse<ApiResponse<Approval>>> => {
  return service.get(`/approvals/${params.id}`);
};

/**
 * 分页获取 Approval 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<Approval[]>>>
 */
export const getApprovalList = (
  params?: GetApprovalListParams
): Promise<AxiosResponse<ApiResponse<Approval[]>>> => {
  return service.get('/approvals', { params });
};

/**
 * 审批通过任务
 *
 * @param taskId - 任务ID
 * @param data - 审批数据
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const approveTask = (
  taskId: number | string,
  data: ApproveTaskRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post(`/approvals/task/${taskId}/approve`, data);
};

/**
 * 拒绝任务
 *
 * @param taskId - 任务ID
 * @param data - 拒绝数据
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const rejectTask = (
  taskId: number | string,
  data: RejectTaskRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post(`/approvals/task/${taskId}/reject`, data);
};

/**
 * 获取审批历史
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getApprovalHistory = (
  params: GetApprovalHistoryParams
): Promise<AxiosResponse<ApiResponse>> => {
  return service.get(`/approvals/${params.approval_id}/history`, { params });
};

/**
 * 获取审批统计信息
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const getApprovalStats = (
  params?: any
): Promise<AxiosResponse<ApiResponse>> => {
  return service.get('/approvals/statistics', { params });
};
