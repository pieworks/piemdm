import { render, h } from "vue";

import myLoad from "@/components/Loading.vue";

const load = {
  loadingCount: 0,

  show(title = "加载中...") {
    if (this.loadingCount === 0) {
      let conf = {
        show: true,
        title: title,
      };
      const vnode = h(myLoad, conf);
      render(vnode, document.body);
    }
    this.loadingCount++;
  },

  hide() {
    if (this.loadingCount <= 0) return;

    this.loadingCount--;
    if (this.loadingCount === 0) {
      render(null, document.body);
    }
  },
};
export { load };
