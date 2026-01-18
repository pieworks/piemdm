import axios from 'axios';
import { useUserStore } from '@/pinia/modules/user';
import { updateApiConfig } from './config';

/**
 * 配置 Axios 拦截器
 */
export function setupInterceptors() {
  // 请求拦截器
  axios.interceptors.request.use(
    (config) => {
      const userStore = useUserStore();

      // 添加 token
      if (userStore.token) {
        config.headers.Authorization = `Bearer ${userStore.token}`;
        // 同时更新 API 配置
        updateApiConfig(userStore.token);
      }

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // 响应拦截器
  axios.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      // 处理 401 未授权
      if (error.response?.status === 401) {
        const userStore = useUserStore();
        userStore.logout();

        // 重定向到登录页
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
      }

      // 处理 403 禁止访问
      if (error.response?.status === 403) {
        console.error('Access denied:', error.response.data);
      }

      return Promise.reject(error);
    }
  );
}
