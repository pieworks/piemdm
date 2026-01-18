import { localStorage } from "@/utils/local-storage";
import { createI18n } from "vue-i18n";
import enUS from "./en-US.js";
import zhCN from "./zh-CN.js";
import zhTW from "./zh-TW.js";

// 定义日期时间格式
const datetimeFormats = {
  "zh-CN": {
    short: {
      year: "numeric",
      month: "short",
      day: "numeric",
    },
    medium: {
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
      hour12: false,
    },
    long: {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
      hour12: false,
    },
    full: {
      year: "numeric",
      month: "long",
      day: "numeric",
      weekday: "long",
      hour: "numeric",
      minute: "numeric",
      second: "numeric",
      hour12: false,
      timeZoneName: "short",
    },
  },
  "en-US": {
    short: {
      year: "numeric",
      month: "short",
      day: "numeric",
    },
    medium: {
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
      hour12: false,
    },
    long: {
      year: "numeric",
      month: "short",
      day: "numeric",
      weekday: "short",
      hour: "numeric",
      minute: "numeric",
      hour12: true,
    },
    full: {
      year: "numeric",
      month: "long",
      day: "numeric",
      weekday: "long",
      hour: "numeric",
      minute: "numeric",
      second: "numeric",
      hour12: true,
      timeZoneName: "short",
    },
  },
  "zh-TW": {
    short: {
      year: "numeric",
      month: "short",
      day: "numeric",
    },
    medium: {
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
      hour12: false,
    },
    long: {
      year: "numeric",
      month: "short",
      day: "numeric",
      weekday: "short",
      hour: "numeric",
      minute: "numeric",
      hour12: false,
    },
    full: {
      year: "numeric",
      month: "long",
      day: "numeric",
      weekday: "long",
      hour: "numeric",
      minute: "numeric",
      second: "numeric",
      hour12: false,
      timeZoneName: "short",
    },
  },
};

const i18n = createI18n({
  locale: localStorage.get("lang") || "zh-CN",
  legacy: false,
  globalInjection: true,
  messages: {
    "zh-CN": {
      ...zhCN,
    },
    "en-US": {
      ...enUS,
    },
    "zh-TW": {
      ...zhTW,
    },
  },
  datetimeFormats, // 添加日期时间格式配置
});

export { i18n };
