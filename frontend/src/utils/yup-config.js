import { i18n } from '@/lang/index';
import * as yup from 'yup';

// 获取当前的 t 函数
const getCurrentTranslation = () => {
  return i18n.global.t;
};

// 设置 Yup 的全局错误信息
export const setupYupLocale = () => {
  yup.setLocale({
    mixed: {
      required: ({ path }) => {
        const t = getCurrentTranslation();
        // 根据字段名获取对应的翻译，如果没有则使用字段名本身
        const fieldName = t(path) || path;
        return fieldName + t('is required');
      },
      notType: ({ path, type }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return `${fieldName} must be a ${type}`;
      },
    },
    string: {
      min: ({ path, min }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return fieldName + t('must be at least {count} characters', { count: min });
      },
      max: ({ path, max }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return fieldName + t('cannot exceed {count} characters', { count: max });
      },
      email: ({ path }) => {
        const t = getCurrentTranslation();
        return t('Please enter a valid') + t('email');
      },
      url: ({ path }) => {
        const t = getCurrentTranslation();
        return t('Please enter a valid') + t('URL');
      },
    },
    number: {
      min: ({ path, min }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return fieldName + t('must be at least') + ' ' + min;
      },
      max: ({ path, max }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return fieldName + t('cannot exceed') + ' ' + max;
      },
      positive: ({ path }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return fieldName + t('must be positive');
      },
      integer: ({ path }) => {
        const t = getCurrentTranslation();
        const fieldName = t(path) || path;
        return fieldName + t('must be integer');
      },
    },
  });
};

// 创建一个工厂函数来创建带有自定义错误信息的 yup schema
export const createYupSchema = () => {
  return yup;
};

// 创建一个辅助函数来创建带有默认错误信息的字段
export const createField = (type = 'string') => {
  switch (type) {
    case 'string':
      return yup.string();
    case 'number':
      return yup.number();
    case 'boolean':
      return yup.boolean();
    case 'date':
      return yup.date();
    case 'array':
      return yup.array();
    case 'object':
      return yup.object();
    default:
      return yup.mixed();
  }
};

// 导出配置好的 yup
export { yup };
