/**
 * 字段值格式化工具
 * 提供统一的字段值格式化功能
 */

import dictionaryService from '@/services/dictionaryService';

/**
 * 格式化字段值(统一入口)
 * @param {any} value - 字段值
 * @param {Object} field - 字段配置
 * @param {Object} options - 格式化选项
 * @returns {Promise<string|null>} 格式化后的值
 */
export async function formatFieldValue(value, field, options = {}) {
  // 处理 null/undefined
  if (value === null || value === undefined) {
    return null;
  }

  // 处理空字符串
  if (value === '') {
    return null;
  }

  // 处理空数组
  if (Array.isArray(value) && value.length === 0) {
    return null;
  }



  // 关联字段(select, radio, multiselect, checkboxgroup)
  if (field.relation || field.field_type === 'select' || field.field_type === 'radio' ||
      field.field_type === 'multiselect' || field.field_type === 'checkboxgroup') {
    return formatRelationValue(value, field);
  }

  // 日期字段
  if (field.field_type === 'date' || field.type === 'Date') {
    return formatDateValue(value);
  }

  // 日期时间字段
  if (field.field_type === 'datetime' || field.type === 'DateTime') {
    return formatDateTimeValue(value);
  }

  // 时间字段
  if (field.field_type === 'time') {
    return value;
  }

  // 其他字段:直接返回原始值
  return value;
}

/**
 * 格式化关联字段值
 * @param {string|Array} value - code 或 code 数组
 * @param {Object} field - 字段配置
 * @returns {Promise<string|null>} 格式化后的值
 */
export async function formatRelationValue(value, field) {
  if (!value) return null;



  // 处理后端 JOIN 返回的对象格式 (兼容旧的后端实现)
  if (typeof value === 'object' && !Array.isArray(value)) {
    // 如果是 {code: null, name: null},返回 null
    if (value.code === null && value.name === null) {
      return null;
    }
    // 返回 "code name" 格式
    if (value.code && value.name) {
      return `${value.code} ${value.name}`;
    }
    return value.name || value.code || null;
  }

  // 处理 JSON 字符串格式的对象
  if (typeof value === 'string' && value.startsWith('{')) {
    try {
      const obj = JSON.parse(value);
      if (obj.code === null && obj.name === null) {
        return null;
      }
      // 返回 "code name" 格式
      if (obj.code && obj.name) {
        return `${obj.code} ${obj.name}`;
      }
      return obj.name || obj.code || null;
    } catch (e) {
      // 解析失败,继续按 code 处理
    }
  }

  // 处理 JSON 字符串格式的数组 (多选字段)
  if (typeof value === 'string' && value.startsWith('[')) {
    try {
      const arr = JSON.parse(value);
      if (Array.isArray(arr)) {
        value = arr;  // 转换为数组继续处理
      }
    } catch (e) {
      // 解析失败,继续按原值处理
    }
  }

  // 获取字典类别
  // 获取字典编码
  const dictCode = field.relation?.filter?.dict_code ||
                          field.Options?.relation?.filter?.dict_code;

  // 获取关联表配置
  const target = field.relation?.target || field.Options?.relation?.target;
  const valueField = field.relation?.valueField || field.Options?.relation?.valueField || 'code';
  const labelField = field.relation?.labelField || field.Options?.relation?.labelField || 'name';
  const filter = field.relation?.filter || field.Options?.relation?.filter || {};

  let dict = [];

  if (dictCode) {
    // 字典表关联:使用 dict 过滤
    dict = await dictionaryService.getDictionary(dictCode);
  } else if (target && target !== 'dict_item') {
    // 非字典表关联:使用通用方法查询
    dict = await dictionaryService.getTableData(target, filter, valueField, labelField);
  } else {
    // 没有配置关联信息,直接返回原值
    return Array.isArray(value) ? value.join(', ') : value;
  }

  // 统一处理:单选和多选
  const codes = Array.isArray(value) ? value : [value];
  const formatted = codes
    .filter(code => code)  // 过滤空值
    .map(code => {
      const item = dict.find(d => d.code === code || d.code === String(code));
      if (item) {
        // 返回 "code name" 格式
        return `${item.code} ${item.name}`;
      }
      return code;  // 找不到时只显示 code
    });

  // 返回格式化结果
  return Array.isArray(value) ? formatted.join(', ') : formatted[0] || null;
}

/**
 * 格式化日期值
 * @param {string} value - 日期字符串
 * @returns {string|null} 格式化后的日期
 */
export function formatDateValue(value) {
  if (!value) return null;

  // 只取日期部分,去除时间和时区
  if (typeof value === 'string') {
    // 移除时区信息和时间部分
    const dateOnly = value.split('T')[0].split(' ')[0];
    return dateOnly;
  }

  return value;
}

/**
 * 格式化日期时间值
 * @param {string} value - 日期时间字符串
 * @returns {string|null} 格式化后的日期时间
 */
export function formatDateTimeValue(value) {
  if (!value) return null;

  if (typeof value === 'string') {
    // 移除时区信息,保留日期和时间
    return value.replace(/\+\d{2}:\d{2}$/, '').replace('T', ' ');
  }

  return value;
}

/**
 * 批量预加载字段所需的字典数据
 * @param {Array} fields - 字段配置数组
 */
export async function preloadFieldDictionaries(fields) {
  const dictCodes = fields
    .filter(field => field.relation || field.field_type === 'select' ||
                     field.field_type === 'radio' || field.field_type === 'multiselect' ||
                     field.field_type === 'checkboxgroup')
    .map(field => field.relation?.filter?.dict_code ||
                  field.Options?.relation?.filter?.dict_code)
    .filter(dc => dc);

  // 去重
  const uniqueCodes = [...new Set(dictCodes)];

  // 批量预加载
  await dictionaryService.preloadDictionaries(uniqueCodes);
}
