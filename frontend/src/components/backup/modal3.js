import { render, h } from "vue";

import myModal from "@/components/Modal.vue";

const modal = {
  show(title = "加载中...") {
    let conf = {
      show: true,
      display: "block",
    };
    const vnode = h(myModal, conf);
    render(vnode, document.body);
  },

  hide() {
    render(null, document.body);
  },
};
export { modal };
