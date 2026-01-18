import { defineStore } from "pinia";

export const useLayoutStore = defineStore("layout", {
  state: () => ({
    activePageName: "",
  }),
  actions: {
    setActivePageName(name) {
      this.activePageName = name;
    },
  },
});
