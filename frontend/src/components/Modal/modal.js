import Modal from "@/components/Modal/Modal.vue";
import { h, render } from "vue";

function mountModal(props, slots = {}) {
  const container = document.createElement("div");
  document.body.appendChild(container);

  return new Promise((resolve) => {
    let autoCloseTimer = null;

    const cleanup = () => {
      if (autoCloseTimer) clearTimeout(autoCloseTimer);

      // 清理 Vue 组件
      render(null, container);
      if (document.body.contains(container)) {
        document.body.removeChild(container);
      }

      // 作为安全措施,清理可能残留的 Bootstrap Modal 副作用
      // 正常情况下 Modal.dispose() 应该已经处理了这些
      const backdrops = document.querySelectorAll('.modal-backdrop');
      if (backdrops.length > 0) {
        backdrops.forEach(backdrop => backdrop.remove());
        document.body.classList.remove('modal-open');
        document.body.style.paddingRight = '';
      }
    };

    const vnode = h(
      Modal,
      {
        ...props,
        show: true,
        "onUpdate:show": (val) => {
          if (!val) {
            resolve(false); // Treat closing via backdrop/X as cancel
            cleanup();
          }
        },
        onConfirm: () => {
          resolve(true);
          // cleanup(); // Cleanup happens in onUpdate:show (hidden.bs.modal)
        },
        onCancel: () => {
          resolve(false);
          // cleanup(); // Cleanup happens in onUpdate:show (hidden.bs.modal)
        },
      },
      slots
    );
    render(vnode, container);

    if (props.autoClose && typeof props.autoClose === 'number') {
      autoCloseTimer = setTimeout(() => {
        resolve(false); // Auto-close acts as dismissal
        cleanup();
      }, props.autoClose);
    }
  });
}

const AppModal = {
  alert(options = {}) {
    return mountModal(
      {
        title: options.title || "提示",
        okTitle: options.okTitle || "确定",
        // Alert: No cancel button, auto-close by default
        showCancelButton: false,
        autoClose: options.autoClose !== undefined ? options.autoClose : 5000,
        centered: options.centered !== false,
        size: options.size || "md",
        scrollable: options.scrollable !== false,
        animation: options.animation !== false,
        ...options,
      },
      {
        default:
          typeof options.content === "function"
            ? options.content
            : () => options.content || options.bodyContent,
      }
    );
  },
  confirm(options = {}) {
    return mountModal(
      {
        title: options.title || "确认",
        okTitle: options.okTitle || "确定",
        cancelTitle: options.cancelTitle || "取消",
        // Confirm: Show cancel button, no auto-close
        showCancelButton: true,
        centered: options.centered !== false,
        size: options.size || "md",
        scrollable: options.scrollable !== false,
        animation: options.animation !== false,
        ...options,
      },
      {
        default:
          typeof options.content === "function"
            ? options.content
            : () => options.content || options.bodyContent,
      }
    );
  },
  info(options = {}) {
    return mountModal(
      {
        title: options.title || "信息",
        okTitle: options.okTitle || "确定",
        // Info: No cancel button, manual close only (no auto-close)
        showCancelButton: false,
        centered: options.centered !== false,
        size: options.size || "md",
        scrollable: options.scrollable !== false,
        animation: options.animation !== false,
        ...options,
      },
      {
        default:
          typeof options.content === "function"
            ? options.content
            : () => options.content || options.bodyContent,
      }
    );
  },
};

export { AppModal };
