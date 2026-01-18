import { useApprovalStore } from '@/pinia/modules/approval';
import { useUserStore } from '@/pinia/modules/user';
import { reactive, ref } from 'vue';

/**
 * WebSocket连接管理器
 * 用于实时更新审批状态和通知
 */
class WebSocketManager {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectInterval = 5000;
    this.heartbeatInterval = 30000;
    this.heartbeatTimer = null;
    this.reconnectTimer = null;

    // 连接状态
    this.isConnected = ref(false);
    this.isConnecting = ref(false);
    this.lastError = ref(null);

    // 消息统计
    this.stats = reactive({
      messagesReceived: 0,
      messagesSent: 0,
      lastMessageTime: null,
      connectionTime: null,
    });

    // 事件监听器
    this.eventListeners = new Map();

    // 消息队列（离线时缓存）
    this.messageQueue = [];

    this.init();
  }

  /**
   * 初始化WebSocket连接
   */
  init() {
    const userStore = useUserStore();

    // 监听用户登录状态变化
    userStore.$subscribe((mutation, state) => {
      if (state.token && !this.isConnected.value) {
        this.connect();
      } else if (!state.token && this.isConnected.value) {
        this.disconnect();
      }
    });

    // 如果用户已登录，立即连接
    if (userStore.token) {
      this.connect();
    }
  }

  /**
   * 建立WebSocket连接
   */
  connect() {
    if (this.isConnecting.value || this.isConnected.value) {
      return;
    }

    try {
      this.isConnecting.value = true;
      this.lastError.value = null;

      const userStore = useUserStore();
      const wsUrl = this.getWebSocketUrl();

      console.log('[WebSocket] 正在连接...', wsUrl);

      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = this.handleOpen.bind(this);
      this.ws.onmessage = this.handleMessage.bind(this);
      this.ws.onclose = this.handleClose.bind(this);
      this.ws.onerror = this.handleError.bind(this);
    } catch (error) {
      console.error('[WebSocket] 连接失败:', error);
      this.lastError.value = error.message;
      this.isConnecting.value = false;
      this.scheduleReconnect();
    }
  }

  /**
   * 断开WebSocket连接
   */
  disconnect() {
    console.log('[WebSocket] 主动断开连接');

    this.clearTimers();
    this.reconnectAttempts = 0;

    if (this.ws) {
      this.ws.close(1000, '用户主动断开');
      this.ws = null;
    }

    this.isConnected.value = false;
    this.isConnecting.value = false;
  }

  /**
   * 获取WebSocket连接URL
   */
  getWebSocketUrl() {
    const userStore = useUserStore();
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_WS_HOST || window.location.host;
    const token = userStore.token;

    return `${protocol}//${host}/ws/approval?token=${encodeURIComponent(token)}`;
  }

  /**
   * 处理连接打开事件
   */
  handleOpen(event) {
    console.log('[WebSocket] 连接已建立');

    this.isConnected.value = true;
    this.isConnecting.value = false;
    this.reconnectAttempts = 0;
    this.stats.connectionTime = new Date();

    // 发送认证消息
    this.sendAuth();

    // 启动心跳
    this.startHeartbeat();

    // 处理离线期间的消息队列
    this.processMessageQueue();

    // 触发连接成功事件
    this.emit('connected', event);
  }

  /**
   * 处理消息接收事件
   */
  handleMessage(event) {
    try {
      const message = JSON.parse(event.data);
      this.stats.messagesReceived++;
      this.stats.lastMessageTime = new Date();

      console.log('[WebSocket] 收到消息:', message);

      // 处理不同类型的消息
      this.handleMessageByType(message);
    } catch (error) {
      console.error('[WebSocket] 消息解析失败:', error, event.data);
    }
  }

  /**
   * 处理连接关闭事件
   */
  handleClose(event) {
    console.log('[WebSocket] 连接已关闭:', event.code, event.reason);

    this.isConnected.value = false;
    this.isConnecting.value = false;
    this.clearTimers();

    // 触发断开连接事件
    this.emit('disconnected', event);

    // 如果不是正常关闭，尝试重连
    if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
      this.scheduleReconnect();
    }
  }

  /**
   * 处理连接错误事件
   */
  handleError(event) {
    console.error('[WebSocket] 连接错误:', event);

    this.lastError.value = '连接错误';
    this.isConnecting.value = false;

    // 触发错误事件
    this.emit('error', event);
  }

  /**
   * 根据消息类型处理消息
   */
  handleMessageByType(message) {
    const { type, data } = message;

    switch (type) {
      case 'approval_status_changed':
        this.handleApprovalStatusChanged(data);
        break;

      case 'new_task_assigned':
        this.handleNewTaskAssigned(data);
        break;

      case 'task_completed':
        this.handleTaskCompleted(data);
        break;

      case 'approval_withdrawn':
        this.handleApprovalWithdrawn(data);
        break;

      case 'approval_urged':
        this.handleApprovalUrged(data);
        break;

      case 'notification':
        this.handleNotification(data);
        break;

      case 'heartbeat':
        this.handleHeartbeat(data);
        break;

      default:
        console.warn('[WebSocket] 未知消息类型:', type);
        this.emit('message', message);
    }
  }

  /**
   * 处理审批状态变更
   */
  handleApprovalStatusChanged(data) {
    const approvalStore = useApprovalStore();

    // 清除相关缓存
    approvalStore.clearCache('approvals');

    // 刷新审批列表
    approvalStore.fetchApprovalList();

    // 如果是当前查看的审批，更新详情
    if (approvalStore.currentApproval?.id === data.approvalId) {
      approvalStore.fetchApprovalDetail(data.approvalId);
    }

    this.emit('approval_status_changed', data);

    // 显示通知
    this.showNotification({
      title: '审批状态更新',
      message: `审批"${data.title}"状态已更新为${data.status}`,
      type: 'info',
    });
  }

  /**
   * 处理新任务分配
   */
  handleNewTaskAssigned(data) {
    const approvalStore = useApprovalStore();

    // 刷新待办任务列表
    approvalStore.fetchPendingTasks();

    this.emit('new_task_assigned', data);

    // 显示通知
    this.showNotification({
      title: '新的待办任务',
      message: `您有新的审批任务："${data.title}"`,
      type: 'warning',
    });
  }

  /**
   * 处理任务完成
   */
  handleTaskCompleted(data) {
    const approvalStore = useApprovalStore();

    // 刷新任务列表
    approvalStore.fetchPendingTasks();
    approvalStore.fetchCompletedTasks();

    this.emit('task_completed', data);
  }

  /**
   * 处理审批撤回
   */
  handleApprovalWithdrawn(data) {
    const approvalStore = useApprovalStore();

    // 清除相关缓存并刷新
    approvalStore.clearCache('approvals');
    approvalStore.fetchApprovalList();

    this.emit('approval_withdrawn', data);

    // 显示通知
    this.showNotification({
      title: '审批已撤回',
      message: `审批"${data.title}"已被撤回`,
      type: 'info',
    });
  }

  /**
   * 处理审批催办
   */
  handleApprovalUrged(data) {
    this.emit('approval_urged', data);

    // 显示通知
    this.showNotification({
      title: '审批催办',
      message: `审批"${data.title}"被催办：${data.message}`,
      type: 'warning',
    });
  }

  /**
   * 处理通知消息
   */
  handleNotification(data) {
    this.emit('notification', data);
    this.showNotification(data);
  }

  /**
   * 处理心跳消息
   */
  handleHeartbeat(data) {
    // 回复心跳
    this.send({
      type: 'heartbeat_response',
      timestamp: Date.now(),
    });
  }

  /**
   * 发送认证消息
   */
  sendAuth() {
    const userStore = useUserStore();

    this.send({
      type: 'auth',
      token: userStore.token,
      userId: userStore.userInfo?.id,
      timestamp: Date.now(),
    });
  }

  /**
   * 启动心跳
   */
  startHeartbeat() {
    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected.value) {
        this.send({
          type: 'heartbeat',
          timestamp: Date.now(),
        });
      }
    }, this.heartbeatInterval);
  }

  /**
   * 发送消息
   */
  send(message) {
    if (!this.isConnected.value) {
      console.warn('[WebSocket] 连接未建立，消息已加入队列:', message);
      this.messageQueue.push(message);
      return false;
    }

    try {
      const messageStr = JSON.stringify(message);
      this.ws.send(messageStr);
      this.stats.messagesSent++;

      console.log('[WebSocket] 发送消息:', message);
      return true;
    } catch (error) {
      console.error('[WebSocket] 发送消息失败:', error);
      return false;
    }
  }

  /**
   * 处理消息队列
   */
  processMessageQueue() {
    while (this.messageQueue.length > 0 && this.isConnected.value) {
      const message = this.messageQueue.shift();
      this.send(message);
    }
  }

  /**
   * 安排重连
   */
  scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('[WebSocket] 重连次数已达上限，停止重连');
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectInterval * Math.pow(2, this.reconnectAttempts - 1);

    console.log(`[WebSocket] ${delay}ms 后进行第 ${this.reconnectAttempts} 次重连`);

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, delay);
  }

  /**
   * 清除定时器
   */
  clearTimers() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  /**
   * 显示通知
   */
  showNotification(notification) {
    // 这里可以集成具体的通知组件
    console.log('[Notification]', notification);

    // 如果浏览器支持原生通知
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification(notification.title, {
        body: notification.message,
        icon: '/favicon.ico',
      });
    }

    // 触发通知事件
    this.emit('notification_show', notification);
  }

  /**
   * 添加事件监听器
   */
  on(event, callback) {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, []);
    }
    this.eventListeners.get(event).push(callback);
  }

  /**
   * 移除事件监听器
   */
  off(event, callback) {
    if (this.eventListeners.has(event)) {
      const listeners = this.eventListeners.get(event);
      const index = listeners.indexOf(callback);
      if (index !== -1) {
        listeners.splice(index, 1);
      }
    }
  }

  /**
   * 触发事件
   */
  emit(event, data) {
    if (this.eventListeners.has(event)) {
      this.eventListeners.get(event).forEach(callback => {
        try {
          callback(data);
        } catch (error) {
          console.error(`[WebSocket] 事件处理器错误 (${event}):`, error);
        }
      });
    }
  }

  /**
   * 获取连接状态
   */
  getStatus() {
    return {
      isConnected: this.isConnected.value,
      isConnecting: this.isConnecting.value,
      lastError: this.lastError.value,
      reconnectAttempts: this.reconnectAttempts,
      stats: { ...this.stats },
    };
  }

  /**
   * 请求通知权限
   */
  async requestNotificationPermission() {
    if ('Notification' in window) {
      const permission = await Notification.requestPermission();
      return permission === 'granted';
    }
    return false;
  }
}

// 创建单例实例
const wsManager = new WebSocketManager();

// 组合式函数
export function useWebSocket() {
  return {
    isConnected: wsManager.isConnected,
    isConnecting: wsManager.isConnecting,
    lastError: wsManager.lastError,
    stats: wsManager.stats,

    connect: () => wsManager.connect(),
    disconnect: () => wsManager.disconnect(),
    send: message => wsManager.send(message),
    on: (event, callback) => wsManager.on(event, callback),
    off: (event, callback) => wsManager.off(event, callback),
    getStatus: () => wsManager.getStatus(),
    requestNotificationPermission: () => wsManager.requestNotificationPermission(),
  };
}

export default wsManager;
