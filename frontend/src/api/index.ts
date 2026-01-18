import { apiConfig } from './config';
import { DefaultApi } from './generated/api';

/**
 * 导出 API 实例
 *
 * 使用方式：
 * ```typescript
 * import { api } from '@/api';
 *
 * // 调用 API
 * const response = await api.apiV1UsersGet();
 * console.log(response.data);
 * ```
 */
export const api = new DefaultApi(apiConfig);

/**
 * 导出所有生成的类型
 */
export * from './generated/models';

/**
 * 导出配置和拦截器
 */
export { apiConfig, updateApiConfig } from './config';
export { setupInterceptors } from './interceptors';
