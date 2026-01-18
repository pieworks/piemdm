/**
 * API 通用类型定义
 *
 * 此文件包含所有 API 调用共享的类型定义,避免在各个 API 文件中重复定义
 */

/**
 * API 响应包装
 *
 * 后端统一的响应格式
 *
 * @template T - 响应数据的类型
 */
export interface ApiResponse<T = any> {
  /** 响应数据 */
  data: T;
  /** 响应代码 */
  code?: number;
  /** 响应消息 */
  message?: string;
}

/**
 * 分页查询参数
 */
export interface PaginationParams {
  /** 页码,从 1 开始 */
  page?: number;
  /** 每页数量 */
  pageSize?: number;
}

/**
 * 时间范围查询参数
 */
export interface DateRangeParams {
  /** 开始日期 */
  startDate?: string;
  /** 结束日期 */
  endDate?: string;
}

/**
 * 批量删除请求参数
 */
export interface BatchDeleteRequest {
  /** 要删除的 ID 列表 */
  ids: number[];
}

/**
 * 批量更新状态请求参数
 */
export interface BatchUpdateStatusRequest {
  /** 要更新的 ID 列表 */
  ids: number[];
  /** 目标状态 */
  status: string;
}

/**
 * 通用查询参数
 *
 * 包含分页和时间范围
 */
export interface CommonQueryParams extends PaginationParams, DateRangeParams {
  /** 显示类型 */
  showType?: string;
}
