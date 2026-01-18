import { computed, onMounted, onUnmounted, ref } from "vue";

/**
 * 响应式设计组合式函数
 * 提供设备检测、屏幕尺寸监听等功能
 */
export function useResponsive() {
  // 响应式数据
  const windowWidth = ref(window.innerWidth);
  const windowHeight = ref(window.innerHeight);

  // 断点定义
  const breakpoints = {
    xs: 0,
    sm: 576,
    md: 768,
    lg: 992,
    xl: 1200,
    xxl: 1400,
  };

  // 计算属性：当前屏幕类型
  const screenType = computed(() => {
    const width = windowWidth.value;
    if (width >= breakpoints.xxl) return "xxl";
    if (width >= breakpoints.xl) return "xl";
    if (width >= breakpoints.lg) return "lg";
    if (width >= breakpoints.md) return "md";
    if (width >= breakpoints.sm) return "sm";
    return "xs";
  });

  // 计算属性：设备类型判断
  const isMobile = computed(() => windowWidth.value < breakpoints.md);
  const isTablet = computed(
    () => windowWidth.value >= breakpoints.md && windowWidth.value < breakpoints.lg
  );
  const isDesktop = computed(() => windowWidth.value >= breakpoints.lg);
  const isSmallScreen = computed(() => windowWidth.value < breakpoints.lg);

  // 计算属性：方向判断
  const isLandscape = computed(() => windowWidth.value > windowHeight.value);
  const isPortrait = computed(() => windowWidth.value <= windowHeight.value);

  // 计算属性：触摸设备检测
  const isTouchDevice = computed(() => {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  });

  // 计算属性：高分辨率屏幕检测
  const isHighDPI = computed(() => {
    return window.devicePixelRatio > 1;
  });

  // 计算属性：暗色模式检测
  const isDarkMode = computed(() => {
    return window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
  });

  // 计算属性：减少动画偏好检测
  const prefersReducedMotion = computed(() => {
    return window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  });

  // 窗口大小变化处理函数
  const handleResize = () => {
    windowWidth.value = window.innerWidth;
    windowHeight.value = window.innerHeight;
  };

  // 生命周期管理
  onMounted(() => {
    window.addEventListener("resize", handleResize);
    handleResize();
  });

  onUnmounted(() => {
    window.removeEventListener("resize", handleResize);
  });

  // 工具函数：检查是否匹配指定断点
  const matchBreakpoint = (breakpoint) => {
    return windowWidth.value >= breakpoints[breakpoint];
  };

  // 工具函数：检查是否在指定断点范围内
  const matchBreakpointRange = (min, max) => {
    const width = windowWidth.value;
    const minWidth = breakpoints[min] || 0;
    const maxWidth = breakpoints[max] || Infinity;
    return width >= minWidth && width < maxWidth;
  };

  // 工具函数：获取响应式类名
  const getResponsiveClass = (baseClass, breakpointClasses = {}) => {
    const classes = [baseClass];

    Object.entries(breakpointClasses).forEach(([breakpoint, className]) => {
      if (matchBreakpoint(breakpoint)) {
        classes.push(className);
      }
    });

    return classes.join(" ");
  };

  // 工具函数：获取响应式值
  const getResponsiveValue = (values) => {
    const sortedBreakpoints = Object.keys(breakpoints).sort(
      (a, b) => breakpoints[b] - breakpoints[a]
    );

    for (const breakpoint of sortedBreakpoints) {
      if (values[breakpoint] !== undefined && matchBreakpoint(breakpoint)) {
        return values[breakpoint];
      }
    }

    return values.default || values.xs || null;
  };

  // 工具函数：获取移动端优化的表格配置
  const getTableConfig = () => {
    if (isMobile.value) {
      return {
        size: "sm",
        responsive: true,
        className: "mobile-table",
      };
    }

    return {
      size: "default",
      responsive: true,
      className: "desktop-table",
    };
  };

  // 工具函数：获取移动端优化的分页配置
  const getPaginationConfig = () => {
    if (isMobile.value) {
      return {
        size: "sm",
        showQuickJumper: false,
        showSizeChanger: false,
        showTotal: false,
        simple: true,
      };
    }

    return {
      size: "default",
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: true,
      simple: false,
    };
  };

  // 工具函数：获取移动端优化的模态框配置
  const getModalConfig = () => {
    if (isMobile.value) {
      return {
        size: "fullscreen",
        centered: false,
        scrollable: true,
        className: "mobile-modal",
      };
    }

    return {
      size: "lg",
      centered: true,
      scrollable: true,
      className: "desktop-modal",
    };
  };

  // 工具函数：获取移动端优化的表单配置
  const getFormConfig = () => {
    if (isMobile.value) {
      return {
        layout: "vertical",
        size: "default",
        className: "mobile-form",
      };
    }

    return {
      layout: "horizontal",
      size: "default",
      className: "desktop-form",
    };
  };

  return {
    windowWidth,
    windowHeight,
    screenType,
    isMobile,
    isTablet,
    isDesktop,
    isSmallScreen,
    isLandscape,
    isPortrait,
    isTouchDevice,
    isHighDPI,
    isDarkMode,
    prefersReducedMotion,
    matchBreakpoint,
    matchBreakpointRange,
    getResponsiveClass,
    getResponsiveValue,
    getTableConfig,
    getPaginationConfig,
    getModalConfig,
    getFormConfig,
    breakpoints,
  };
}

/**
 * 移动端优化的表格组合式函数
 */
export function useMobileTable() {
  const { isMobile, isTablet } = useResponsive();

  // 获取移动端优化的列配置
  const getMobileColumns = (columns) => {
    if (!isMobile.value) return columns;

    // 移动端只显示关键列
    return columns
      .filter((col) => col.mobile !== false)
      .map((col) => ({
        ...col,
        width: col.mobileWidth || col.width,
        ellipsis: true,
      }));
  };

  // 获取移动端优化的操作列配置
  const getMobileActions = (actions) => {
    if (!isMobile.value) return actions;

    // 移动端合并操作到下拉菜单
    const primaryActions = actions.filter((action) => action.primary);
    const secondaryActions = actions.filter((action) => !action.primary);

    if (secondaryActions.length > 2) {
      return [
        ...primaryActions.slice(0, 1),
        {
          type: "dropdown",
          text: "更多",
          actions: secondaryActions,
        },
      ];
    }

    return actions;
  };

  return {
    isMobile,
    isTablet,
    getMobileColumns,
    getMobileActions,
  };
}

/**
 * 移动端优化的导航组合式函数
 */
export function useMobileNavigation() {
  const { isMobile } = useResponsive();
  const sidebarVisible = ref(false);

  const toggleSidebar = () => {
    sidebarVisible.value = !sidebarVisible.value;
  };

  const closeSidebar = () => {
    sidebarVisible.value = false;
  };

  const openSidebar = () => {
    sidebarVisible.value = true;
  };

  // 移动端自动关闭侧边栏
  const handleRouteChange = () => {
    if (isMobile.value) {
      closeSidebar();
    }
  };

  return {
    isMobile,
    sidebarVisible,
    toggleSidebar,
    closeSidebar,
    openSidebar,
    handleRouteChange,
  };
}
