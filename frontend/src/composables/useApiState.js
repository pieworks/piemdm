import { computed, onUnmounted, reactive, ref } from 'vue';

/**
 * API状态管理组合式函数
 * 用于统一管理API调用的加载状态、错误处理和数据缓存
 */
export function useApiState() {
  // 全局加载状态
  const globalLoading = ref(false);
  const loadingCount = ref(0);

  // API调用状态映射
  const apiStates = reactive(new Map());

  // 错误状态
  const errors = reactive(new Map());

  // 缓存状态
  const cache = reactive(new Map());
  const cacheExpiry = reactive(new Map());

  // 请求取消控制器
  const abortControllers = reactive(new Map());

  // 计算属性
  const hasErrors = computed(() => errors.size > 0);
  const errorList = computed(() => Array.from(errors.values()));

  /**
   * 创建API状态
   * @param {string} key - API标识键
   * @param {Object} options - 配置选项
   */
  const createApiState = (key, options = {}) => {
    const state = reactive({
      loading: false,
      error: null,
      data: null,
      lastFetch: null,
      retryCount: 0,
      ...options.initialState,
    });

    apiStates.set(key, state);
    return state;
  };

  /**
   * 获取API状态
   * @param {string} key - API标识键
   */
  const getApiState = key => {
    return apiStates.get(key) || createApiState(key);
  };

  /**
   * 设置加载状态
   * @param {string} key - API标识键
   * @param {boolean} loading - 加载状态
   */
  const setLoading = (key, loading) => {
    const state = getApiState(key);

    if (loading && !state.loading) {
      loadingCount.value++;
      state.loading = true;
    } else if (!loading && state.loading) {
      loadingCount.value--;
      state.loading = false;
    }

    // 更新全局加载状态
    globalLoading.value = loadingCount.value > 0;
  };

  /**
   * 设置错误状态
   * @param {string} key - API标识键
   * @param {Error|string|null} error - 错误信息
   */
  const setError = (key, error) => {
    const state = getApiState(key);

    if (error) {
      const errorInfo = {
        key,
        message: typeof error === 'string' ? error : error.message,
        code: error.code || 'UNKNOWN_ERROR',
        timestamp: new Date(),
        stack: error.stack,
      };

      state.error = errorInfo;
      errors.set(key, errorInfo);
    } else {
      state.error = null;
      errors.delete(key);
    }
  };

  /**
   * 设置数据
   * @param {string} key - API标识键
   * @param {any} data - 数据
   */
  const setData = (key, data) => {
    const state = getApiState(key);
    state.data = data;
    state.lastFetch = new Date();
  };

  /**
   * 执行API调用
   * @param {string} key - API标识键
   * @param {Function} apiCall - API调用函数
   * @param {Object} options - 配置选项
   */
  const executeApi = async (key, apiCall, options = {}) => {
    const {
      useCache = false,
      cacheTTL = 5 * 60 * 1000, // 5分钟
      retryCount = 0,
      retryDelay = 1000,
      timeout = 30000,
      onSuccess,
      onError,
      onFinally,
      ...requestConfig
    } = options;

    // 检查缓存
    if (useCache) {
      const cached = getCache(key);
      if (cached) {
        setData(key, cached);
        return cached;
      }
    }

    // 取消之前的请求
    cancelRequest(key);

    // 创建取消控制器
    const abortController = new AbortController();
    abortControllers.set(key, abortController);

    try {
      setLoading(key, true);
      setError(key, null);

      // 设置超时
      const timeoutId = setTimeout(() => {
        abortController.abort();
      }, timeout);

      // 执行API调用
      const result = await apiCall({
        signal: abortController.signal,
        ...requestConfig,
      });

      clearTimeout(timeoutId);

      // 设置数据
      setData(key, result);

      // 缓存结果
      if (useCache) {
        setCache(key, result, cacheTTL);
      }

      // 重置重试计数
      const state = getApiState(key);
      state.retryCount = 0;

      // 成功回调
      if (onSuccess) {
        onSuccess(result);
      }

      return result;
    } catch (error) {
      // 如果是取消请求，不处理错误
      if (error.name === 'AbortError') {
        return null;
      }

      const state = getApiState(key);
      state.retryCount++;

      // 设置错误
      setError(key, error);

      // 重试逻辑
      if (state.retryCount <= retryCount) {
        console.log(`[API] 重试第 ${state.retryCount} 次: ${key}`);

        await new Promise(resolve => setTimeout(resolve, retryDelay * state.retryCount));

        return executeApi(key, apiCall, options);
      }

      // 错误回调
      if (onError) {
        onError(error);
      }

      throw error;
    } finally {
      setLoading(key, false);
      abortControllers.delete(key);

      // 完成回调
      if (onFinally) {
        onFinally();
      }
    }
  };

  /**
   * 取消请求
   * @param {string} key - API标识键
   */
  const cancelRequest = key => {
    const controller = abortControllers.get(key);
    if (controller) {
      controller.abort();
      abortControllers.delete(key);
    }
  };

  /**
   * 取消所有请求
   */
  const cancelAllRequests = () => {
    abortControllers.forEach((controller, key) => {
      controller.abort();
    });
    abortControllers.clear();
  };

  /**
   * 设置缓存
   * @param {string} key - 缓存键
   * @param {any} data - 数据
   * @param {number} ttl - 过期时间（毫秒）
   */
  const setCache = (key, data, ttl = 5 * 60 * 1000) => {
    cache.set(key, data);
    cacheExpiry.set(key, Date.now() + ttl);
  };

  /**
   * 获取缓存
   * @param {string} key - 缓存键
   */
  const getCache = key => {
    const expiry = cacheExpiry.get(key);
    if (expiry && Date.now() < expiry) {
      return cache.get(key);
    }

    // 清除过期缓存
    cache.delete(key);
    cacheExpiry.delete(key);
    return null;
  };

  /**
   * 清除缓存
   * @param {string|RegExp} pattern - 缓存键模式
   */
  const clearCache = pattern => {
    if (!pattern) {
      cache.clear();
      cacheExpiry.clear();
      return;
    }

    const keys = Array.from(cache.keys());

    if (typeof pattern === 'string') {
      // 字符串匹配
      keys.forEach(key => {
        if (key.includes(pattern)) {
          cache.delete(key);
          cacheExpiry.delete(key);
        }
      });
    } else if (pattern instanceof RegExp) {
      // 正则表达式匹配
      keys.forEach(key => {
        if (pattern.test(key)) {
          cache.delete(key);
          cacheExpiry.delete(key);
        }
      });
    }
  };

  /**
   * 清除错误
   * @param {string} key - API标识键
   */
  const clearError = key => {
    if (key) {
      setError(key, null);
    } else {
      errors.clear();
      apiStates.forEach(state => {
        state.error = null;
      });
    }
  };

  /**
   * 重置API状态
   * @param {string} key - API标识键
   */
  const resetApiState = key => {
    if (key) {
      const state = getApiState(key);
      state.loading = false;
      state.error = null;
      state.data = null;
      state.lastFetch = null;
      state.retryCount = 0;

      errors.delete(key);
    } else {
      // 重置所有状态
      apiStates.clear();
      errors.clear();
      loadingCount.value = 0;
      globalLoading.value = false;
    }
  };

  /**
   * 批量执行API
   * @param {Array} apiCalls - API调用配置数组
   * @param {Object} options - 配置选项
   */
  const executeBatch = async (apiCalls, options = {}) => {
    const { concurrent = true, stopOnError = false } = options;

    if (concurrent) {
      // 并行执行
      const promises = apiCalls.map(({ key, apiCall, options }) =>
        executeApi(key, apiCall, options).catch(error => ({ error, key }))
      );

      const results = await Promise.all(promises);

      // 检查是否有错误
      const errors = results.filter(result => result && result.error);
      if (errors.length > 0 && stopOnError) {
        throw new Error(`批量API调用失败: ${errors.map(e => e.key).join(', ')}`);
      }

      return results;
    } else {
      // 串行执行
      const results = [];

      for (const { key, apiCall, options } of apiCalls) {
        try {
          const result = await executeApi(key, apiCall, options);
          results.push(result);
        } catch (error) {
          results.push({ error, key });

          if (stopOnError) {
            throw error;
          }
        }
      }

      return results;
    }
  };

  /**
   * 获取API状态摘要
   */
  const getStateSummary = () => {
    return {
      totalApis: apiStates.size,
      loadingApis: Array.from(apiStates.entries())
        .filter(([key, state]) => state.loading)
        .map(([key]) => key),
      errorApis: Array.from(errors.keys()),
      globalLoading: globalLoading.value,
      cacheSize: cache.size,
    };
  };

  // 组件卸载时清理
  onUnmounted(() => {
    cancelAllRequests();
  });

  return {
    // 状态
    globalLoading,
    hasErrors,
    errorList,

    // API状态管理
    createApiState,
    getApiState,
    setLoading,
    setError,
    setData,
    executeApi,
    cancelRequest,
    cancelAllRequests,

    // 缓存管理
    setCache,
    getCache,
    clearCache,

    // 错误管理
    clearError,

    // 状态重置
    resetApiState,

    // Batch operations
    executeBatch,

    // 工具方法
    getStateSummary,
  };
}

/**
 * 创建特定API的状态管理
 * @param {string} apiName - API名称
 * @param {Object} options - 配置选项
 */
export function useSpecificApiState(apiName, options = {}) {
  const apiState = useApiState();
  const state = apiState.getApiState(apiName);

  const execute = (apiCall, executeOptions = {}) => {
    return apiState.executeApi(apiName, apiCall, { ...options, ...executeOptions });
  };

  const cancel = () => {
    apiState.cancelRequest(apiName);
  };

  const reset = () => {
    apiState.resetApiState(apiName);
  };

  const clearError = () => {
    apiState.clearError(apiName);
  };

  return {
    // 状态
    loading: computed(() => state.loading),
    error: computed(() => state.error),
    data: computed(() => state.data),
    lastFetch: computed(() => state.lastFetch),
    retryCount: computed(() => state.retryCount),

    // 方法
    execute,
    cancel,
    reset,
    clearError,
  };
}
