import { reactive, render, h } from "vue";
import ToastComponent from "@/components/Toast.vue";

const container = document.createElement('div');
document.body.appendChild(container);

const toast = {
  show(options) {
    const props = {
      show: true,
      message: options.message || "",
      header: options.header || false,
      headerTitle: options.headerTitle || "",
      color: options.color || "primary",
    };

    const vnode = h(ToastComponent, props);
    render(vnode, container);

    let duration = options.duration || 2000;
    setTimeout(() => {
      // Create a vnode with show: false to trigger transition/hide if needed,
      // but simplistic approach is just to unmount or re-render
      // However, to follow Vue pattern with the component we wrote:
      // We should probably just let it unmount or hide.
      // Current Toast.vue uses v-if or class binding for show.

      // Let's just unmount after duration.
      render(null, container);
    }, duration);
  },

  hide() {
    render(null, container);
  },
};
export const AppToast = toast;
