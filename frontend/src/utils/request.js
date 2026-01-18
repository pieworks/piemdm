import router from '@/router';
import apiConfig from '@/api/endpoints';
import { load } from '@/components/loading.js';
import { AppModal } from '@/components/Modal/modal';
import { AppToast } from '@/components/toast.js';
import { useUserStore } from '@/pinia/modules/user';
import axios from 'axios'; // 引入axios

const service = axios.create({
  // baseURL: "http://127.0.0.1:8787/api/v1", // url = base url + request url
  baseURL: apiConfig.baseURL, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000,
});

// http request 拦截器
service.interceptors.request.use(
  config => {
    const cacheControl = import.meta.env.VITE_CACHE_CONTROL || 'no-cache';
    const { globalLoading = true } = config;

    // ===== Batch operations IDs 校验 =====
    // 检查是否是Batch operations (URL 包含 /batch)
    const isBatchOperation = config.url && config.url.includes('/batch');
    if (isBatchOperation && config.data) {
      // 支持 ids/IDs，以及 DELETE 请求中嵌套的 data.ids
      const data = config.data;
      const ids = data.ids || data.IDs || data.data?.ids || data.data?.IDs;
      if (!ids || (Array.isArray(ids) && ids.length === 0)) {
        AppModal.alert({
          title: '提示',
          content: '请先选择要操作的记录',
        });
        // 返回一个被取消的 Promise，阻止请求发送
        return Promise.reject();
      }
    }
    // ===== Batch operations校验结束 =====

    if (globalLoading) {
      load.show();
    }

    // ===== 动态路径转换: 根据当前路由自动添加 /admin 前缀 =====
    // 检查当前路由是否在 /admin 下
    const isAdminRoute = window.location.pathname.startsWith('/admin');

    // 只处理相对路径,跳过完整 URL
    const isFullUrl = config.url && (config.url.startsWith('http://') || config.url.startsWith('https://'));

    // 排除不需要添加 /admin 前缀的路径
    const excludePaths = [
      '/auth/',           // 认证相关 (公共端点)
      '/admin/',          // 已经有 /admin 前缀的
    ];

    const shouldAddAdminPrefix =
      isAdminRoute &&
      config.url &&
      !isFullUrl &&  // 跳过完整 URL
      !excludePaths.some(path => config.url.startsWith(path));

    if (shouldAddAdminPrefix) {
      config.url = `/admin${config.url}`;
    }
    // ===== 路径转换结束 =====

    const userStore = useUserStore();
    config.headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + userStore.token,
      'X-User-Id': userStore.userInfo?.ID || '',
      'Cache-Control': cacheControl,
      ...config.headers,
    };
    return config;
  },
  error => {
    const { globalLoading = true, globalError = true } = error.config || {};
    let message = error.response?.data?.message || decodeURI(error.response?.headers?.message || '') || error.message;

    if (globalError) {
      AppToast.show((message));
      AppModal.alert({
        title: '服务器错误',
        bodyHtml: true,
        bodyContent: message,
      });
    }

    if (globalLoading) {
      load.hide();
    }
    return Promise.reject(error);
  }
);

// http response 拦截器
service.interceptors.response.use(
  response => {
    const userStore = useUserStore();
    const { globalLoading = true, globalError = true } = response.config;

    if (globalLoading) {
      load.hide();
    }

    if (response.headers['new-token']) {
      userStore.setToken(response.headers['new-token']);
    }

    if (response.status === 200 || response.statusText === 'OK') {
      return response;
    } else {
      let message = response.data.message || decodeURI(response.headers.message);

      if (globalError) {
        AppModal.alert({
          title: '服务器错误',
          bodyHtml: true,
          bodyContent: message,
        });
      }

      // 如果响应中包含 reload 标志，则执行以下操作:
      // 1. 清除用户 token
      // 2. 清除本地存储
      // 3. 获取当前路由路径
      // 4. 根据路径判断重定向到管理员登录页还是普通登录页
      if (response.data && response.data.reload) {
        userStore.token = '';
        localStorage.clear();
        const currentPath = router.currentRoute.value.path;
        if (currentPath.startsWith('/admin')) {
          router.push({ name: 'AdminLogin', replace: true });
        } else {
          router.push({ name: 'Login', replace: true });
        }
      }
      return response.data.message ? response.data : response;
    }
  },
  error => {
    const { globalLoading = true, globalError = true } = error.config || {};

    if (globalLoading) {
      load.hide();
    }

    if (globalError) {
      // 根据状态码进行错误提示
      switch (error.response?.status) {
        case 500:
          AppModal.alert({
            title: '接口报错',
            bodyHtml: true,
            bodyContent: `
              <p>检测到接口错误${error}</p>
              <p>错误信息: <span style="color:red">${error.response?.data?.message}</span></p>
              <p>错误码<span style="color:red"> ${error.response?.status} </span>：此类错误内容常见于后台panic，请先查看后台日志，如果影响您正常使用可强制登出清理缓存</p>
              `,
          });
          break;
        case 404:
        case 400:
          AppModal.alert({
            title: '接口报错',
            bodyHtml: true,
            bodyContent: `
                <p>检测到接口错误: ${error}</p>
                <p>错误信息: <span style="color:red">${error.response?.data?.message}</span></p>
                <p>错误码<span style="color:red"> ${error.response?.status} </span>：此类错误多为接口未注册（或未重启）或者请求路径（方法）与api路径（方法）不符--如果为自动化代码请检查是否存在空格</p>
                `,
          });
          break;
        case 403:
          const errorMsg = error.response?.data?.message || '您没有权限执行此操作';
          if (errorMsg === 'admin access required') {
            router.push('/approval/index');
          } else {
            AppModal.alert({
              title: '权限不足',
              bodyHtml: true,
              bodyContent: `
                <p>检测到权限错误</p>
                <p>错误信息: <span style="color:red">${errorMsg}</span></p>
                <p>请联系管理员为您分配相应的权限</p>
                `,
            });
          }
          break;
        case 401:
          const userStore = useUserStore();
          userStore.token = '';
          localStorage.clear();
          const currentPath = router.currentRoute.value.path;
          if (currentPath.startsWith('/admin')) {
            router.push({ name: 'AdminLogin', replace: true });
          } else {
            router.push({ name: 'Login', replace: true });
          }
          break;
        default:
          AppModal.alert({
            title: '接口报错',
            bodyHtml: true,
            bodyContent: `
              <p>检测到请求错误</p>
              <p>${error}</p>
              `,
          });
      }
    }

    // 抛出异常而不是返回 error,确保调用方能正确捕获错误
    return Promise.reject(error);
  }
);

export const loginApi = (credentials, endpoint, config = {}) => {
  return service.post(apiConfig.baseURL + endpoint, credentials, config);
};

export const validateTokenApi = (config = {}) => {
  return service.post(
    apiConfig.validateToken,
    {},
    {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
      ...config,
    }
  );
};

export const logoutApi = (config = {}) => {
  return service.post(
    apiConfig.logout,
    {},
    {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
      ...config,
    }
  );
};

export default service;
