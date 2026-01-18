import { useI18n } from 'vue-i18n';

// 对Date的扩展，将 Date 转化为指定格式的String
// 月(M)、日(d)、小时(h)、分(m)、秒(s)、季度(q) 可以用 1-2 个占位符，
// 年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字)
// (new Date()).Format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423
// (new Date()).Format("yyyy-M-d h:m:s.S")      ==> 2006-7-2 8:9:4.18
// eslint-disable-next-line no-extend-native
Date.prototype.Format = function (fmt) {
  var o = {
    'M+': this.getMonth() + 1, // 月份
    'd+': this.getDate(), // 日
    'h+': this.getHours(), // 小时
    'm+': this.getMinutes(), // 分
    's+': this.getSeconds(), // 秒
    'q+': Math.floor((this.getMonth() + 3) / 3), // 季度
    'S': this.getMilliseconds(), // 毫秒
  };
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (this.getFullYear() + '').substr(4 - RegExp.$1.length));
  }
  for (var k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length === 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length)
      );
    }
  }
  return fmt;
};

export function formatTimeToStr(times, pattern) {
  var d = new Date(times).Format('yyyy-MM-dd hh:mm:ss');
  if (pattern) {
    d = new Date(times).Format(pattern);
  }
  return d.toLocaleString();
}

/**
 * 检查是否是Go的零值时间
 * @param {string|Date} dateValue
 * @returns {boolean}
 */
export function isGoZeroTime(dateValue) {
  if (!dateValue) return true;

  const zeroTimePatterns = [
    '0001-01-01T00:00:00Z',
    '0001-01-01 00:00:00',
    '0001-01-01T00:00:00.000Z',
  ];

  if (typeof dateValue === 'string') {
    return zeroTimePatterns.some(pattern => dateValue.startsWith(pattern.substring(0, 10)));
  }

  if (dateValue instanceof Date) {
    return dateValue.getFullYear() === 1;
  }

  return false;
}

/**
 * 安全的日期格式化函数
 * @param {string|Date} dateValue
 * @param {string} format
 * @returns {string}
 */
export function formatDate(dateValue, format = 'long') {
  if (!dateValue || isGoZeroTime(dateValue)) {
    return '';
  }

  try {
    const { d } = useI18n();
    const date = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;

    if (isNaN(date.getTime())) {
      return '';
    }

    return d(date, format);
  } catch (error) {
    console.warn('日期格式化失败:', error);
    return '';
  }
}

/**
 * 检查日期是否有效（非空且非零值）
 * @param {string|Date} dateValue
 * @returns {boolean}
 */
export function isValidDate(dateValue) {
  if (!dateValue || isGoZeroTime(dateValue)) {
    return false;
  }

  const date = typeof dateValue === 'string' ? new Date(dateValue) : dateValue;
  return !isNaN(date.getTime());
}
