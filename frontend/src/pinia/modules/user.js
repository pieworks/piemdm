import router from '@/router/index';
import { loginApi, logoutApi, validateTokenApi } from '@/utils/request';
import { defineStore } from 'pinia';


// 用户状态管理优化
export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    userInfo: {},
    roles: [],
    permissions: [], // User permissions
    isAdmin: false, // 明确标记是否管理员
    loginTime: 0,
  }),

  actions: {
    // 用户登录（普通用户）
    async userLogin(credentials) {
      this.isAdmin = false;
      const result = await this.commonLogin(credentials, '/auth/login');
      if (result.success) {
        router.push('/approval/index');
      }
      return result;
    },

    // 管理员登录
    async adminLogin(credentials) {
      this.isAdmin = true;
      const result = await this.commonLogin(credentials, '/auth/login');
      if (result.success) {
        router.push('/admin/dashboard');
      }
      return result;
    },

    // 切换到管理员模式
    enableAdminMode() {
      this.isAdmin = true;
      localStorage.setItem('isAdmin', 'true');
    },

    // 通用登录逻辑
    async commonLogin(credentials, endpoint) {
      try {
        const { data } = await loginApi(credentials, endpoint);
        data.roles = [];
        this.setUserData(data);
        return { success: true };
      } catch (error) {
        return { success: false, message: error.message };
      }
    },

    // 设置用户数据（提取共同逻辑）
    setUserData(data) {
      this.token = data.token;
      this.userInfo = data.userInfo;
      this.roles = data.roles || [];
      // Support direct permissions list from backend
      this.permissions = data.permissions || [];
      this.loginTime = new Date().getTime();

      // 持久化存储
      localStorage.setItem('token', data.token);
      localStorage.setItem('userInfo', JSON.stringify(data.userInfo));
      localStorage.setItem('roles', JSON.stringify(this.roles));
      localStorage.setItem('permissions', JSON.stringify(this.permissions));
      localStorage.setItem('isAdmin', this.isAdmin.toString());
      localStorage.setItem('loginTime', this.loginTime.toString());
    },

    // 登出处理
    logout() {
      // 获取当前用户类型（保存以便决定重定向位置）
      const isAdmin = this.isAdmin;

      // 清除状态和存储
      this.clearUserData();

      // JWT 认证不需要后端登出接口，只需清除前端 token
      // 使用 window.location 强制跳转，避免路由守卫问题
      window.location.href = isAdmin ? '/admin/site/login' : '/site/login';
    },

    // 清除用户数据
    clearUserData() {
      // 清除状态
      this.token = '';
      this.userInfo = null;
      this.roles = [];
      this.permissions = [];
      this.isAdmin = false;
      this.loginTime = 0;

      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      localStorage.removeItem('roles');
      localStorage.removeItem('permissions');
      localStorage.removeItem('isAdmin');
      localStorage.removeItem('loginTime');
    },

    // 会话恢复
    restoreSession() {
      const token = localStorage.getItem('token');
      if (!token) return false;

      try {
        // 恢复基本数据
        this.token = token;
        this.userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
        this.roles = JSON.parse(localStorage.getItem('roles') || '[]');
        this.permissions = JSON.parse(localStorage.getItem('permissions') || '[]');
        this.isAdmin = localStorage.getItem('isAdmin') === 'true';
        this.loginTime = parseInt(localStorage.getItem('loginTime') || '0');

        // 验证令牌有效性（可选）
        this.validateToken();

        return true;
      } catch (error) {
        console.error('会话恢复失败:', error);
        this.clearUserData();
        return false;
      }
    },

    // 验证令牌有效性
    async validateToken() {
      try {
        await validateTokenApi();
        return true;
      } catch (error) {
        this.clearUserData();
        return false;
      }
    },
  },

  getters: {
    // 明确的角色检查
    isAdminUser: state => state.isAdmin,
    isRegularUser: state => !state.isAdmin && !!state.token,

    // 其他通用 getters
    isLoggedIn: state => !!state.token,
    username: state => state.userInfo?.username || '',

    // Check if user can access admin interface
    canAccessAdmin: state => state.isAdmin || state.userInfo?.username === 'admin' || state.roles.some(r => r.code === 'admin'),


    // 权限检查
    hasPermission: state => rule => {
      if (state.isAdmin) return true; // 管理员拥有所有权限
      if (!rule) return true;

      // Check direct permissions first
      if (state.permissions && state.permissions.includes(rule)) {
        return true;
      }

      // Fallback to role permissions (legacy support)
      return state.roles.some(role => role.permissions && role.permissions.includes(rule));
    },
  },
});
