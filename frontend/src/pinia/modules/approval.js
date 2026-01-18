import { getApprovalList } from '@/api/approval';
import httpLinkHeader from 'http-link-header';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useApprovalStore = defineStore('approval', () => {
  // ==================== 状态定义 ====================

  // 审批列表相关状态
  const approvalList = ref([]);
  const approvalTotal = ref(0);
  const approvalLoading = ref(false);
  const approvalError = ref(null);

  // 当前审批详情
  const currentApproval = ref(null);
  const approvalHistory = ref([]);
  const approvalStatistics = ref({});

  // 待办任务相关状态
  const pendingTasks = ref([]);
  const pendingTasksTotal = ref(0);
  const pendingTasksLoading = ref(false);

  // 已办任务相关状态
  const completedTasks = ref([]);
  const completedTasksTotal = ref(0);

  // 筛选和分页状态
  const filters = ref({
    status: '',
    keyword: '',
    startDate: '',
    endDate: '',
    applicant: '',
    approver: '',
  });

  const pagination = ref({
    page: 1,
    pageSize: 20,
  });

  // 选中项状态
  const selectedApprovals = ref([]);
  const selectedTasks = ref([]);

  // 缓存状态
  const cache = ref(new Map());
  const cacheExpiry = ref(new Map());

  // ==================== 计算属性 ====================

  // 审批状态统计
  const approvalStatusStats = computed(() => {
    const stats = {
      pending: 0,
      approved: 0,
      rejected: 0,
      cancelled: 0,
      withdrawn: 0,
    };

    approvalList.value.forEach(approval => {
      if (stats.hasOwnProperty(approval.status)) {
        stats[approval.status]++;
      }
    });

    return stats;
  });

  // 是否有选中项
  const hasSelectedApprovals = computed(() => selectedApprovals.value.length > 0);
  const hasSelectedTasks = computed(() => selectedTasks.value.length > 0);

  // 筛选后的审批列表
  const filteredApprovals = computed(() => {
    let result = approvalList.value;

    if (filters.value.status) {
      result = result.filter(item => item.status === filters.value.status);
    }

    if (filters.value.keyword) {
      const keyword = filters.value.keyword.toLowerCase();
      result = result.filter(
        item =>
          item.title.toLowerCase().includes(keyword) ||
          item.content.toLowerCase().includes(keyword) ||
          item.applicant.toLowerCase().includes(keyword)
      );
    }

    if (filters.value.startDate) {
      result = result.filter(item => new Date(item.createdAt) >= new Date(filters.value.startDate));
    }

    if (filters.value.endDate) {
      result = result.filter(item => new Date(item.createdAt) <= new Date(filters.value.endDate));
    }

    return result;
  });

  // ==================== 缓存管理 ====================

  const getCacheKey = (type, params = {}) => {
    return `${type}_${JSON.stringify(params)}`;
  };

  const setCache = (key, data, ttl = 5 * 60 * 1000) => {
    cache.value.set(key, data);
    cacheExpiry.value.set(key, Date.now() + ttl);
  };

  const getCache = key => {
    const expiry = cacheExpiry.value.get(key);
    if (expiry && Date.now() < expiry) {
      return cache.value.get(key);
    }

    // 清除过期缓存
    cache.value.delete(key);
    cacheExpiry.value.delete(key);
    return null;
  };

  const clearCache = pattern => {
    if (pattern) {
      for (const key of cache.value.keys()) {
        if (key.includes(pattern)) {
          cache.value.delete(key);
          cacheExpiry.value.delete(key);
        }
      }
    } else {
      cache.value.clear();
      cacheExpiry.value.clear();
    }
  };

  // ==================== Actions ====================

  // 获取审批列表
  const fetchApprovalList = async (params = {}) => {
    const cacheKey = getCacheKey('approvals', { ...filters.value, ...pagination.value, ...params });
    const cached = getCache(cacheKey);

    console.log('fetchApprovalList cached: ', cached);

    if (cached) {
      approvalList.value = cached.list;
      approvalTotal.value = cached.total;
      return cached;
    }

    try {
      approvalLoading.value = true;
      approvalError.value = null;

      console.log("filters.value: ", filters.value)
      console.log("pagination.value: ", pagination.value)
      console.log("params: ", params)

      const response = await getApprovalList({
        ...filters.value,
        ...pagination.value,
        ...params,
      });

      console.log('fetchApprovalList response: ', response);

      approvalList.value = response.data || [];

      if (response.total !== undefined) {
        approvalTotal.value = response.total;
      } else if (response.headers && response.headers.link) {
        const links = httpLinkHeader.parse(response.headers.link).refs;
        links.forEach(link => {
          if (['last'].includes(link.rel)) {
            const url = new URL(link.uri);
            approvalTotal.value = parseInt(url.searchParams.get('page')) || 1;
          }
        });
      }

      // 缓存结果
      const result = {
        list: approvalList.value,
        total: approvalTotal.value,
      };
      setCache(cacheKey, result);

      return result;
    } catch (error) {
      approvalError.value = error.message;
      console.error('获取审批列表失败:', error);
      throw error;
    } finally {
      approvalLoading.value = false;
    }
  };

  // ==================== 筛选和分页 ====================

  const setFilters = newFilters => {
    filters.value = { ...filters.value, ...newFilters };
    pagination.value.page = 1; // 重置页码
  };

  const setPagination = newPagination => {
    pagination.value = { ...pagination.value, ...newPagination };
  };

  const resetFilters = () => {
    filters.value = {
      status: '',
      keyword: '',
      startDate: '',
      endDate: '',
      applicant: '',
      approver: '',
    };
    pagination.value.page = 1;
  };

  // ==================== 选择管理 ====================

  const selectApproval = approval => {
    if (!selectedApprovals.value.find(item => item.id === approval.id)) {
      selectedApprovals.value.push(approval);
    }
  };

  const unselectApproval = approvalId => {
    const index = selectedApprovals.value.findIndex(item => item.id === approvalId);
    if (index > -1) {
      selectedApprovals.value.splice(index, 1);
    }
  };

  const toggleApprovalSelection = approval => {
    const isSelected = selectedApprovals.value.find(item => item.id === approval.id);
    if (isSelected) {
      unselectApproval(approval.id);
    } else {
      selectApproval(approval);
    }
  };

  const toggleAllApprovals = () => {
    if (selectedApprovals.value.length === approvalList.value.length) {
      selectedApprovals.value = [];
    } else {
      selectedApprovals.value = [...approvalList.value];
    }
  };

  const clearSelections = () => {
    selectedApprovals.value = [];
    selectedTasks.value = [];
  };

  // ==================== 状态重置 ====================

  const resetState = () => {
    approvalList.value = [];
    approvalTotal.value = 0;
    currentApproval.value = null;
    approvalHistory.value = [];
    pendingTasks.value = [];
    pendingTasksTotal.value = 0;
    completedTasks.value = [];
    completedTasksTotal.value = 0;
    selectedApprovals.value = [];
    selectedTasks.value = [];
    clearCache();
  };

  // ==================== 返回 Store 接口 ====================

  return {
    // 状态
    approvalList,
    approvalTotal,
    approvalLoading,
    approvalError,
    currentApproval,
    approvalHistory,
    approvalStatistics,
    pendingTasks,
    pendingTasksTotal,
    pendingTasksLoading,
    completedTasks,
    completedTasksTotal,
    filters,
    pagination,
    selectedApprovals,
    selectedTasks,

    // 计算属性
    approvalStatusStats,
    hasSelectedApprovals,
    hasSelectedTasks,
    filteredApprovals,

    // Actions
    fetchApprovalList,

    // 筛选和分页
    setFilters,
    setPagination,
    resetFilters,

    // 选择管理
    selectApproval,
    unselectApproval,
    toggleApprovalSelection,
    toggleAllApprovals,
    clearSelections,

    // 缓存管理
    clearCache,
    cache,
    cacheExpiry,

    // 生命周期
    resetState,
  };
});
