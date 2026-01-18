import { render, h } from 'vue'
import Message from './Message.vue'

function getInstance(config) {
  const conf = Object.assign({}, config);
  if (conf.duration === undefined) {
    conf.duration = 3000;
  }

  const close = () => {
    render(null, document.body)
  };
  conf.close = close;

  const vnode = h(Message, conf);
  render(vnode, document.body);

  setTimeout(() => {
    render(null, document.body)
  }, conf.duration);

  return {close}

}

const message = (opts) => {
  return getInstance(opts);
};

message.install = (app) => {
  // 这样就可以值选项式api中通过this.$message的形式调用了 返回一个包含close的对象，可以手动关闭，也支持自动关闭
  app.config.globalProperties.$message = message;
};

message.success = (msg) => {
  let conf = {
    type: "success",
    message: msg,
  };
  return getInstance(conf);
};

message.error = (msg) => {
  let conf = {
    type: "error",
    message: msg,
  };
  return getInstance(conf);
};

export const AppMessage = message;
