/**
 * API Endpoints Configuration
 *
 * API 端点配置
 */

export interface ApiEndpoints {
  baseURL: string;
  userLogin: string;
  userLogout: string;
  adminLogin: string;
  adminLogout: string;
  validateToken: string;
}

const apiConfig: Record<string, ApiEndpoints> = {
  development: {
    baseURL: 'http://localhost:8787/api/v1',
    userLogin: '/auth/login',
    userLogout: '/auth/logout',
    adminLogin: '/admin/auth/login',
    adminLogout: '/admin/auth/logout',
    validateToken: '/auth/validate',
  },

  production: {
    baseURL: 'https://api.example.com/api/v1',
    userLogin: '/auth/login',
    userLogout: '/auth/logout',
    adminLogin: '/admin/auth/login',
    adminLogout: '/admin/auth/logout',
    validateToken: '/auth/validate',
  },
};

const currentEnv = (import.meta as any).env?.MODE || 'development';
export default apiConfig[currentEnv] as ApiEndpoints;
