import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

// ==================== Mock 设置 ====================

// Mock API 调用
vi.mock("@/api/approval", () => ({
  getApprovalList: vi.fn(),
  createApproval: vi.fn(),
  deleteApproval: vi.fn(),
  findApproval: vi.fn(),
  updateApproval: vi.fn(),
}));

// Mock 用户工具
vi.mock("@/utils", () => ({
  getUser: vi.fn(),
}));

// 导入模块（在 mock 之后）
import { getApprovalList } from "@/api/approval";
import { getUser } from "@/utils";
import { createPinia, setActivePinia } from "pinia";
import { useApprovalStore } from "./approval.js";

// ==================== 测试套件 ====================

describe("Approval Store 测试", () => {
  let store;

  beforeEach(() => {
    // 创建新的 Pinia 实例
    setActivePinia(createPinia());

    // 重置所有 mocks
    vi.clearAllMocks();

    // 设置默认的用户信息
    vi.mocked(getUser).mockReturnValue({
      id: "user_123",
      name: "Test User",
      token: "test-token",
    });

    // 创建 store 实例
    store = useApprovalStore();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  // ==================== 初始状态测试 ====================

  describe("初始状态", () => {
    it("应该有正确的初始状态", () => {
      expect(store.approvalList).toEqual([]);
      expect(store.approvalTotal).toBe(0);
      expect(store.approvalLoading).toBe(false);
      expect(store.approvalError).toBeNull();
      expect(store.pendingTasks).toEqual([]);
      expect(store.pendingTasksTotal).toBe(0);
      expect(store.selectedApprovals).toEqual([]);
    });

    it("应该正确计算统计信息", () => {
      // 设置一些测试数据
      store.approvalList = [
        { id: "1", status: "pending" },
        { id: "2", status: "approved" },
        { id: "3", status: "rejected" },
      ];

      const stats = store.approvalStatusStats;
      expect(stats.pending).toBe(1);
      expect(stats.approved).toBe(1);
      expect(stats.rejected).toBe(1);
    });

    it("应该正确计算选中状态", () => {
      expect(store.hasSelectedApprovals).toBe(false);

      store.selectedApprovals = [{ id: "1" }];
      expect(store.hasSelectedApprovals).toBe(true);
    });
  });

  // ==================== 审批管理测试 ====================

  describe("审批管理", () => {
    it("应该能够获取审批列表", async () => {
      const mockApprovals = [
        { id: "1", title: "测试审批1", status: "pending" },
        { id: "2", title: "测试审批2", status: "approved" },
      ];

      vi.mocked(getApprovalList).mockResolvedValue({
        data: mockApprovals,
        total: 2,
        headers: { link: '' },
      });

      const result = await store.fetchApprovalList({ page: 1, pageSize: 10 });

      expect(vi.mocked(getApprovalList)).toHaveBeenCalledWith(
        expect.objectContaining({
          page: 1,
          pageSize: 10,
        })
      );
      expect(store.approvalList).toEqual(mockApprovals);
      expect(store.approvalTotal).toBe(2);
      expect(result.list).toEqual(mockApprovals);
      expect(result.total).toBe(2);
    });

    it("应该处理获取审批列表失败的情况", async () => {
      const errorMessage = "网络错误";
      vi.mocked(getApprovalList).mockRejectedValue(new Error(errorMessage));

      await expect(store.fetchApprovalList()).rejects.toThrow(errorMessage);
      expect(store.approvalError).toBe(errorMessage);
      expect(store.approvalLoading).toBe(false);
    });

    it("应该使用缓存", async () => {
      const mockApprovals = [{ id: "1", title: "测试审批1", status: "pending" }];

      vi.mocked(getApprovalList).mockResolvedValue({
        data: mockApprovals,
        total: 1,
        headers: { link: '' },
      });

      // 第一次调用
      await store.fetchApprovalList({ page: 1, pageSize: 10 });
      expect(vi.mocked(getApprovalList)).toHaveBeenCalledTimes(1);

      // 第二次调用应该使用缓存
      await store.fetchApprovalList({ page: 1, pageSize: 10 });
      expect(vi.mocked(getApprovalList)).toHaveBeenCalledTimes(1); // 仍然是1次
    });
  });

  // ==================== 筛选和分页测试 ====================

  describe("筛选和分页", () => {
    it("应该能够设置筛选条件", () => {
      const filters = { status: "pending", keyword: "测试" };
      store.setFilters(filters);

      expect(store.filters.status).toBe("pending");
      expect(store.filters.keyword).toBe("测试");
      expect(store.pagination.page).toBe(1); // 应该重置页码
    });

    it("应该能够设置分页", () => {
      const pagination = { page: 2, pageSize: 20 };
      store.setPagination(pagination);

      expect(store.pagination.page).toBe(2);
      expect(store.pagination.pageSize).toBe(20);
    });

    it("应该能够重置筛选条件", () => {
      store.setFilters({ status: "pending", keyword: "测试" });
      store.resetFilters();

      expect(store.filters.status).toBe("");
      expect(store.filters.keyword).toBe("");
      expect(store.pagination.page).toBe(1);
    });

    it("应该正确筛选审批列表", () => {
      store.approvalList = [
        { id: "1", title: "测试审批", status: "pending", applicant: "张三", content: "测试内容" },
        { id: "2", title: "其他审批", status: "approved", applicant: "李四", content: "其他内容" },
      ];

      // 按状态筛选
      store.setFilters({ status: "pending" });
      expect(store.filteredApprovals).toHaveLength(1);
      expect(store.filteredApprovals[0].id).toBe("1");

      // 按关键词筛选
      store.setFilters({ status: "", keyword: "测试" });
      expect(store.filteredApprovals).toHaveLength(1);
      expect(store.filteredApprovals[0].id).toBe("1");
    });
  });

  // ==================== 选择管理测试 ====================

  describe("选择管理", () => {
    beforeEach(() => {
      store.approvalList = [
        { id: "1", title: "审批1" },
        { id: "2", title: "审批2" },
        { id: "3", title: "审批3" },
      ];
    });

    it("应该能够选择审批", () => {
      const approval = store.approvalList[0];
      store.selectApproval(approval);

      expect(store.selectedApprovals).toHaveLength(1);
      expect(store.selectedApprovals[0].id).toBe("1");
    });

    it("应该能够取消选择审批", () => {
      const approval = store.approvalList[0];
      store.selectApproval(approval);
      store.unselectApproval("1");

      expect(store.selectedApprovals).toHaveLength(0);
    });

    it("应该能够切换选择状态", () => {
      const approval = store.approvalList[0];

      // 选择
      store.toggleApprovalSelection(approval);
      expect(store.selectedApprovals).toHaveLength(1);

      // 取消选择
      store.toggleApprovalSelection(approval);
      expect(store.selectedApprovals).toHaveLength(0);
    });

    it("应该能够全选/取消全选", () => {
      // 全选
      store.toggleAllApprovals();
      expect(store.selectedApprovals).toHaveLength(3);

      // 取消全选
      store.toggleAllApprovals();
      expect(store.selectedApprovals).toHaveLength(0);
    });

    it("应该能够清除所有选择", () => {
      store.selectApproval(store.approvalList[0]);
      store.clearSelections();

      expect(store.selectedApprovals).toHaveLength(0);
      expect(store.selectedTasks).toHaveLength(0);
    });
  });

  // ==================== 缓存管理测试 ====================

  describe("缓存管理", () => {
    it("应该能够清除缓存", () => {
      // 设置一些缓存
      store.cache.set("test_key", "test_value");
      store.cacheExpiry.set("test_key", Date.now() + 60000);

      store.clearCache();

      expect(store.cache.size).toBe(0);
      expect(store.cacheExpiry.size).toBe(0);
    });

    it("应该能够按模式清除缓存", () => {
      store.cache.set("approval_1", "data1");
      store.cache.set("approval_2", "data2");
      store.cache.set("task_1", "data3");

      store.clearCache("approval");

      expect(store.cache.has("approval_1")).toBe(false);
      expect(store.cache.has("approval_2")).toBe(false);
      expect(store.cache.has("task_1")).toBe(true);
    });
  });

  // ==================== 状态重置测试 ====================

  describe("状态重置", () => {
    it("应该能够重置所有状态", () => {
      // 设置一些状态
      store.approvalList = [{ id: "1" }];
      store.approvalTotal = 10;
      store.selectedApprovals = [{ id: "1" }];
      store.cache.set("test", "value");

      store.resetState();

      expect(store.approvalList).toEqual([]);
      expect(store.approvalTotal).toBe(0);
      expect(store.selectedApprovals).toEqual([]);
      expect(store.cache.size).toBe(0);
    });
  });
});
