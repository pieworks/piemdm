import { config } from "@vue/test-utils";
import { vi } from "vitest";

// ==================== 全局测试配置 ====================

// 配置 Vue Test Utils
config.global.stubs = {
  // 存根化路由组件
  "router-link": true,
  "router-view": true,

  // 存根化第三方组件
  "vue-select": true,
  "datepicker": true,
};

// ==================== Mock 全局对象 ====================

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
vi.stubGlobal("localStorage", localStorageMock);

// Mock sessionStorage
const sessionStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
vi.stubGlobal("sessionStorage", sessionStorageMock);

// Mock window.location
const locationMock = {
  href: "http://localhost:3000",
  origin: "http://localhost:3000",
  protocol: "http:",
  host: "localhost:3000",
  hostname: "localhost",
  port: "3000",
  pathname: "/",
  search: "",
  hash: "",
  assign: vi.fn(),
  replace: vi.fn(),
  reload: vi.fn(),
};
vi.stubGlobal("location", locationMock);

// Mock window.navigator
const navigatorMock = {
  userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
  language: "zh-CN",
  languages: ["zh-CN", "zh", "en"],
  onLine: true,
};
vi.stubGlobal("navigator", navigatorMock);

// Mock console methods for cleaner test output
global.console = {
  ...console,
  log: vi.fn(),
  debug: vi.fn(),
  info: vi.fn(),
  warn: vi.fn(),
  error: vi.fn(),
};

// ==================== Mock Bootstrap 组件 ====================

// Mock Bootstrap Modal
const mockModal = {
  show: vi.fn(),
  hide: vi.fn(),
  toggle: vi.fn(),
  dispose: vi.fn(),
};

// Mock Bootstrap Toast
const mockToast = {
  show: vi.fn(),
  hide: vi.fn(),
  dispose: vi.fn(),
};

// Mock Bootstrap Tooltip
const mockTooltip = {
  show: vi.fn(),
  hide: vi.fn(),
  toggle: vi.fn(),
  dispose: vi.fn(),
};

// Mock Bootstrap 全局对象
vi.stubGlobal("bootstrap", {
  Modal: vi.fn().mockImplementation(() => mockModal),
  Toast: vi.fn().mockImplementation(() => mockToast),
  Tooltip: vi.fn().mockImplementation(() => mockTooltip),
});

// ==================== Mock Vue Router ====================

vi.mock("vue-router", () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    go: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
  }),
  useRoute: () => ({
    path: "/",
    name: "home",
    params: {},
    query: {},
    meta: {},
    fullPath: "/",
  }),
  createRouter: vi.fn(),
  createWebHistory: vi.fn(),
}));

// ==================== Mock Axios ====================

vi.mock("axios", () => ({
  default: {
    create: vi.fn(() => ({
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
      patch: vi.fn(),
      interceptors: {
        request: {
          use: vi.fn(),
        },
        response: {
          use: vi.fn(),
        },
      },
    })),
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    patch: vi.fn(),
    CancelToken: {
      source: vi.fn(() => ({
        token: "mock-token",
        cancel: vi.fn(),
      })),
    },
  },
}));

// ==================== Mock WebSocket ====================

class MockWebSocket {
  constructor(url) {
    this.url = url;
    this.readyState = WebSocket.CONNECTING;
    this.onopen = null;
    this.onclose = null;
    this.onmessage = null;
    this.onerror = null;

    // 模拟异步连接
    setTimeout(() => {
      this.readyState = WebSocket.OPEN;
      if (this.onopen) {
        this.onopen({ type: "open" });
      }
    }, 100);
  }

  send(data) {
    // 模拟发送消息
  }

  close(code, reason) {
    this.readyState = WebSocket.CLOSED;
    if (this.onclose) {
      this.onclose({ type: "close", code, reason });
    }
  }
}

// WebSocket 常量
MockWebSocket.CONNECTING = 0;
MockWebSocket.OPEN = 1;
MockWebSocket.CLOSING = 2;
MockWebSocket.CLOSED = 3;

vi.stubGlobal("WebSocket", MockWebSocket);

// ==================== Mock Notification API ====================

const mockNotification = vi.fn().mockImplementation((title, options) => ({
  title,
  ...options,
  close: vi.fn(),
}));

mockNotification.permission = "granted";
mockNotification.requestPermission = vi.fn().mockResolvedValue("granted");

vi.stubGlobal("Notification", mockNotification);

// ==================== Mock Audio ====================

class MockAudio {
  constructor(src) {
    this.src = src;
    this.volume = 1;
    this.currentTime = 0;
    this.duration = 0;
    this.paused = true;
  }

  play() {
    this.paused = false;
    return Promise.resolve();
  }

  pause() {
    this.paused = true;
  }

  load() {
    // Mock load
  }
}

vi.stubGlobal("Audio", MockAudio);

// ==================== 测试工具函数 ====================

// 创建测试用的 Pinia store
export function createTestPinia() {
  const { createPinia } = require("pinia");
  return createPinia();
}

// 创建测试用的路由器
export function createTestRouter() {
  const { createRouter, createMemoryHistory } = require("vue-router");
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", component: { template: "<div>Home</div>" } },
      { path: "/approval", component: { template: "<div>Approval</div>" } },
      { path: "/task", component: { template: "<div>Task</div>" } },
    ],
  });
}

// 等待 Vue 的下一个 tick
export function nextTick() {
  return new Promise((resolve) => setTimeout(resolve, 0));
}

// 等待指定时间
export function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// 模拟用户交互
export function mockUserInteraction(wrapper, selector, event = "click") {
  const element = wrapper.find(selector);
  if (element.exists()) {
    element.trigger(event);
  }
  return nextTick();
}

// 模拟表单输入
export function mockFormInput(wrapper, selector, value) {
  const input = wrapper.find(selector);
  if (input.exists()) {
    input.setValue(value);
    input.trigger("input");
  }
  return nextTick();
}

// ==================== 清理函数 ====================

// 在每个测试后清理 mocks
afterEach(() => {
  vi.clearAllMocks();
  localStorageMock.getItem.mockClear();
  localStorageMock.setItem.mockClear();
  localStorageMock.removeItem.mockClear();
  localStorageMock.clear.mockClear();

  sessionStorageMock.getItem.mockClear();
  sessionStorageMock.setItem.mockClear();
  sessionStorageMock.removeItem.mockClear();
  sessionStorageMock.clear.mockClear();
});
