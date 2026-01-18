import { format, formatDistance, isValid } from "date-fns";
import { enUS, zhCN, zhTW } from "date-fns/locale";

// 语言映射配置
export const LANGUAGE_MAPS = {
  // vue-i18n 语言代码到 date-fns 语言包的映射
  "en-US": enUS,
  "zh-CN": zhCN,
  "zh-TW": zhTW,
};

// 获取 date-fns 对应的语言包
export const getDateFnsLocale = (i18nLocale) => {
  return LANGUAGE_MAPS[i18nLocale] || zhCN;
};

// 语言选项
export const LANGUAGE_OPTIONS = [
  { code: "en-US", label: "English" },
  { code: "zh-CN", label: "简体中文" },
  { code: "zh-TW", label: "繁體中文" },
];

// 获取语言标签
export const getLanguageLabel = (code) => {
  const option = LANGUAGE_OPTIONS.find((opt) => opt.code === code);
  return option ? option.label : code;
};

// 日期格式化函数
export const formatDate = (dateString) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return isValid(date) ? format(date, "yyyy-MM-dd") : "";
};

// 相对时间格式化函数
export const formatDateDistance = (dateString, locale = zhCN) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return isValid(date)
    ? formatDistance(date, new Date(), { addSuffix: true, locale })
    : "";
};
