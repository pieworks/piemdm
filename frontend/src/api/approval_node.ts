/**
 * Approval Node API
 *
 * 审批节点管理相关 API 封装
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
 * ApprovalNode 数据模型
 */
export interface ApprovalNode {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  ApprovalDefCode: string;
  NodeCode?: string;
  NodeName: string;
  NodeType: string;
  Description?: string;
  SortOrder?: number;
  ApproverType?: string;
  ApproverConfig?: string;
  ConditionConfig?: string;
  Status?: string;
  CreatedBy?: string;
  UpdatedBy?: string;
}

/**
 * 创建 ApprovalNode 请求参数
 */
export interface CreateApprovalNodeRequest {
  ApprovalDefCode: string;
  NodeCode?: string;
  NodeName: string;
  NodeType: string;
  Description?: string;
  SortOrder?: number;
  ApproverType?: string;
  ApproverConfig?: string;
  ConditionConfig?: string;
  Status?: string;
}

/**
 * 更新 ApprovalNode 请求参数
 */
export interface UpdateApprovalNodeRequest {
  id: number;
  ApprovalDefCode?: string;
  NodeCode?: string;
  NodeName?: string;
  NodeType?: string;
  Description?: string;
  SortOrder?: number;
  ApproverType?: string;
  ApproverConfig?: string;
  ConditionConfig?: string;
  Status?: string;
}

/**
 * 查询参数
 */
export interface GetApprovalNodeListParams {
  page?: number;
  pageSize?: number;
  approvalDefCode?: string;
  nodeType?: string;
  status?: string;
}

/**
 * 查询单个 ApprovalNode 参数
 */
export interface GetApprovalNodeParams {
  id: number;
}

/**
 * 根据 ID 查询 ApprovalNode
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalNode>>>
 */
export const findApprovalNode = (
  params: GetApprovalNodeParams
): Promise<AxiosResponse<ApiResponse<ApprovalNode>>> => {
  return service.get(`/approval_nodes/${params.id}`);
};

/**
 * 创建 ApprovalNode
 *
 * @param data - 创建参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalNode>>>
 */
export const createApprovalNode = (
  data: CreateApprovalNodeRequest
): Promise<AxiosResponse<ApiResponse<ApprovalNode>>> => {
  return service.post('/approval_nodes', data);
};

/**
 * 删除 ApprovalNode
 *
 * @param data - 包含 id 的对象
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const deleteApprovalNode = (
  data: { id: number }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.delete(`/approval_nodes/${data.id}`);
};

/**
 * 批量删除 ApprovalNode
 *
 * @param data - 批量删除参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchDeleteApprovalNode = (
  data: { ids: number[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post('/approval_nodes/batch_delete', data);
};

/**
 * 更新 ApprovalNode
 *
 * @param data - 更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApprovalNode = (
  data: UpdateApprovalNodeRequest
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put(`/approval_nodes/${data.id}`, data);
};

/**
 * 批量更新 ApprovalNode 状态
 *
 * @param data - 批量更新参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const updateApprovalNodeStatus = (
  data: { ids: number[]; status: string }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.put('/approval_nodes/batch', data);
};

/**
 * 批量创建 ApprovalNode
 *
 * @param data - 批量创建参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalNode[]>>>
 */
export const batchCreateApprovalNodes = (
  data: { nodes: CreateApprovalNodeRequest[] }
): Promise<AxiosResponse<ApiResponse<ApprovalNode[]>>> => {
  return service.post('/approval_nodes/batch', data);
};

/**
 * 分页获取 ApprovalNode 列表
 *
 * @param params - 查询参数
 * @returns Promise<AxiosResponse<ApiResponse<ApprovalNode[]>>>
 */
export const getApprovalNodeList = (
  params?: GetApprovalNodeListParams
): Promise<AxiosResponse<ApiResponse<ApprovalNode[]>>> => {
  return service.get('/approval_nodes', { params });
};

/**
 * 批量同步审批节点（支持增删改）
 *
 * @param data - 同步参数
 * @returns Promise<AxiosResponse<ApiResponse>>
 */
export const batchSyncApprovalNodes = (
  data: { approvalDefCode: string; nodes: ApprovalNode[] }
): Promise<AxiosResponse<ApiResponse>> => {
  return service.post(
    `/admin/approval_nodes/sync?approvalDefCode=${data.approvalDefCode}`,
    data.nodes
  );
};
