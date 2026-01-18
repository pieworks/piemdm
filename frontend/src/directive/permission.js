// web/src/directive/permission.js
import { usePermission } from "@/services/permission";
export default {
  mounted(el, binding) {
    const { value } = binding;
    const permission = usePermission();
    if (!permission.canPerformAction(value)) {
      el.parentNode?.removeChild(el);
    }
  },
};
