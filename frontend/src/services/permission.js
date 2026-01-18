// services/permission.js
import { useUserStore } from '@/pinia/modules/user';

export const usePermission = () => {
  const userStore = useUserStore();

  return {
    // 检查是否可以访问路由
    canAccessRoute(route) {
      // 公共路由
      if (!route.meta || !route.meta.requiresAuth) {
        return true;
      }

      // 管理员路由
      if (route.meta.userType === "admin") {
        return userStore.isAdminUser;
      }

      // 普通用户路由
      if (route.meta.userType === "regular") {
        return userStore.isRegularUser;
      }

      // 需要特定权限的路由
      if (route.meta.permission) {
        return userStore.hasPermission(route.meta.permission);
      }

      return true;
    },

    // 检查是否可以执行操作
    canPerformAction(action) {
      // 管理员可以执行所有操作
      if (userStore.isAdminUser) {
        return true;
      }

      // 普通用户需要检查特定权限
      return userStore.hasPermission(action);
    },
  };
};
