import { useUserStore } from '@/pinia/modules/user';
import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  // 公共路由
  {
    path: '/',
    children: [
      { path: 'site/login', component: () => import('views/auth/Login.vue') },
      {
        path: 'admin/site/login',
        component: () => import('@/views/admin/auth/Login.vue'),
      },
    ],
  },

  // 用户路由（需要普通用户权限）
  {
    path: '/',
    name: 'Layout',
    component: () => import('views/layout/Index.vue'),
    redirect: '/site/login',
    meta: { requiresAuth: true, userType: 'regular' },
    children: [
      {
        path: '/site/logout',
        component: () => import('views/auth/Logout.vue'),
      },
      // {
      //   path: '/dashboard/index',
      //   component: () => import('views/dashboard/Index.vue'),
      // },
      {
        path: 'approval',
        children: [
          {
            path: 'index',
            component: () => import('views/approval/Index.vue'),
          },
          {
            path: 'task',
            component: () => import('views/approval/Task.vue'),
          },
        ],
      },
      {
        path: 'entity',
        children: [
          {
            path: 'index',
            component: () => import('views/entity/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/entity/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/entity/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/entity/Create.vue'),
          },
        ],
      },
    ],
  },

  // 管理员路由（需要管理员权限）
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('views/admin/layout/Index.vue'),
    redirect: '/admin/dashboard/index',
    meta: { requiresAuth: true, userType: 'admin' },
    children: [
      {
        path: 'dashboard/index',
        component: () => import('views/admin/dashboard/Index.vue'),
      },
      {
        path: 'site/logout',
        component: () => import('views/admin/auth/Logout.vue'),
      },
      {
        path: 'approval',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/approval/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/approval/View.vue'),
          },
        ],
      },
      {
        path: 'approval_def',
        children: [
          {
            path: '',
            component: () => import('views/admin/approval_def/Index.vue'),
          },
          {
            path: 'index',
            component: () => import('views/admin/approval_def/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/approval_def/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/admin/approval_def/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/admin/approval_def/Create.vue'),
          },
          {
            path: 'designer/:id?',
            component: () => import('views/admin/approval_def/Designer.vue'),
          },
        ],
      },
      {
        path: 'table',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/table/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/table/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/admin/table/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/admin/table/Create.vue'),
          },
        ],
      },
      {
        path: 'table_field',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/table_field/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/table_field/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/admin/table_field/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/admin/table_field/Create.vue'),
          },
        ],
      },
      {
        path: 'application',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/application/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/application/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/admin/application/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/admin/application/Create.vue'),
          },
        ],
      },
      {
        path: 'webhook',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/webhook/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/webhook/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/admin/webhook/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/admin/webhook/Create.vue'),
          },
        ],
      },
      {
        path: 'webhook_delivery',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/webhook_delivery/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/webhook_delivery/View.vue'),
          },
        ],
      },
      {
        path: 'notification_template',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/notification_template/Index.vue'),
            meta: { permission: 'notification_template:list' }
          },
          {
            path: 'view',
            component: () => import('views/admin/notification_template/View.vue'),
            meta: { permission: 'notification_template:list' }
          },
          {
            path: 'update',
            component: () => import('views/admin/notification_template/Update.vue'),
            meta: { permission: 'notification_template:update' }
          },
          {
            path: 'create',
            component: () => import('views/admin/notification_template/Create.vue'),
            meta: { permission: 'notification_template:create' }
          },
        ],
      },
      {
        path: 'notification_log',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/notification_log/Index.vue'),
            meta: { permission: 'notification_log:list' }
          },
          {
            path: 'view',
            component: () => import('views/admin/notification_log/View.vue'),
            meta: { permission: 'notification_log:list' }
          },
        ],
      },
      {
        path: 'cron',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/cron/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/cron/View.vue'),
          },
          {
            path: 'update',
            component: () => import('views/admin/cron/Update.vue'),
          },
          {
            path: 'create',
            component: () => import('views/admin/cron/Create.vue'),
          },
        ],
      },
      {
        path: 'cron_log',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/cron_log/Index.vue'),
          },
          {
            path: 'view',
            component: () => import('views/admin/cron_log/View.vue'),
          },
        ],
      },
      {
        path: 'permission',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/permission/Index.vue'),
            meta: { permission: 'permission:list' }
          },
          {
            path: 'view',
            component: () => import('views/admin/permission/View.vue'),
            meta: { permission: 'permission:list' }
          },
          {
            path: 'update',
            component: () => import('views/admin/permission/Update.vue'),
            meta: { permission: 'permission:update' }
          },
          {
            path: 'create',
            component: () => import('views/admin/permission/Create.vue'),
            meta: { permission: 'permission:create' }
          },
        ],
      },
      {
        path: 'role',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/role/Index.vue'),
            meta: { permission: 'role:list' }
          },
          {
            path: 'view',
            component: () => import('views/admin/role/View.vue'),
            meta: { permission: 'role:list' }
          },
          {
            path: 'update',
            component: () => import('views/admin/role/Update.vue'),
            meta: { permission: 'role:update' }
          },
          {
            path: 'create',
            component: () => import('views/admin/role/Create.vue'),
            meta: { permission: 'role:create' }
          },
        ],
      },
      {
        path: 'user',
        children: [
          {
            path: 'index',
            component: () => import('views/admin/user/Index.vue'),
            meta: { permission: 'user:list' }
          },
          {
            path: 'view',
            component: () => import('views/admin/user/View.vue'),
            meta: { permission: 'user:list' }
          },
          {
            path: 'update',
            component: () => import('views/admin/user/Update.vue'),
            meta: { permission: 'user:update' }
          },
          {
            path: 'create',
            component: () => import('views/admin/user/Create.vue'),
            meta: { permission: 'user:create' }
          },
        ],
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_BASE),
  routes,
});

// router.beforeEach((to, from, next) => {
//   let isAuthenticated = false;
//   // 检查用户登录状态的两种方式
//   let user = getUser();
//   const token = window.localStorage.getItem("token");
//   if ((user && user.name) || token) {
//     isAuthenticated = true;
//   }

//   console.log("路由守卫:", {
//     to: to.path,
//     isAuthenticated,
//     hasUser: !!user,
//     hasToken: !!token,
//   });

//   if (
//     !isAuthenticated &&
//     to.name !== "Login" &&
//     to.name !== "AdminLogin" &&
//     to.name !== "OAuth"
//   ) {
//     console.log("未登录，重定向到登录页");
//     if (to.path.startsWith("/admin")) {
//       next({ name: "AdminLogin" });
//     } else {
//       next({ name: "Login" });
//     }
//   }
//   // else if(isAuthenticated && (to.name == 'Login' || to.name == 'OAuth' )) next({ name: 'Index'})
//   else next();
// });

// 导航守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  // 需要认证的路由
  if (to.meta.requiresAuth) {
    if (!userStore.token) {
      // 未登录时重定向到相应登录页
      next(to.meta.userType === 'admin' ? '/admin/site/login' : '/site/login');
      return;
    }

    // 检查用户类型与路由要求是否匹配
    if (to.meta.userType === 'admin') {
      if (!userStore.isAdminUser) {
        if (userStore.canAccessAdmin) {
          // Auto-switch to admin mode if allowed
          userStore.enableAdminMode();
        } else {
          next('/approval/index'); // 普通用户尝试访问管理页面 -> 跳转到审批首页
          return;
        }
      }
    }

    // if (to.meta.userType === 'regular' && userStore.isAdminUser) {
    //   next('/admin/site/login'); // allow admin to access regular pages
    //   return;
    // }

    // Check granular permission
    if (to.meta.permission) {
      if (!userStore.hasPermission(to.meta.permission)) {
        // Access denied
        // console.warn("Access denied: missing permission", to.meta.permission);
        // next(false); // Cancel navigation or redirect to 403
        // For better UX, might redirect to a forbidden page or just prevent access
        // Since 403 page might not be set up, redirect to dashboard or login?
        // Let's just cancel navigation for now
        return next(false);
      }
    }
  }

  next();
});

export default router;
