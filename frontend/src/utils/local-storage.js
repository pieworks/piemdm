export const localStorage = {
  get(key) {
    return window.localStorage.getItem(key);
  },
  set(key, value) {
    return window.localStorage.setItem(key, value);
  },
  remove(key) {
    return window.localStorage.removeItem(key);
  },
  clear() {
    return window.localStorage.clear();
  },
};
