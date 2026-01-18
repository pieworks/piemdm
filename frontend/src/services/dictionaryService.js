/**
 * 字典管理服务
 * 提供字典数据的缓存、批量加载和缓存失效管理
 */

import { getEntityList } from '@/api/entity';

class DictionaryService {
  constructor() {
    // 字典数据缓存: Map<dictCode, Array<{code, name}>>
    this.cache = new Map();

    // 加载中的请求: Map<dictCode, Promise>
    this.loading = new Map();
  }

  /**
   * 获取字典数据(带缓存)
   * @param {string} dictCode - 字典类别,如 'DIC0064'
   * @returns {Promise<Array>} 字典数据数组
   */
  async getDictionary(dictCode) {
    if (!dictCode) {
      console.warn('getDictionary: dictCode is required');
      return [];
    }

    // 1. 检查缓存
    if (this.cache.has(dictCode)) {
      return this.cache.get(dictCode);
    }

    // 2. 检查是否正在加载(避免重复请求)
    if (this.loading.has(dictCode)) {
      return this.loading.get(dictCode);
    }

    // 3. 加载字典数据
    const loadPromise = this.loadDictionary(dictCode);
    this.loading.set(dictCode, loadPromise);

    try {
      const data = await loadPromise;
      this.cache.set(dictCode, data);
      return data;
    } finally {
      this.loading.delete(dictCode);
    }
  }

  /**
   * 从后端加载字典数据
   * @private
   */
  async loadDictionary(dictCode) {
    try {
      const queryParams = {
        table_code: 'dict_item',
        dict_code: dictCode,
        pageSize: 1000, // 字典数据通常不多,一次加载全部
      };

      const res = await getEntityList(queryParams);

      if (res && res.data) {
        return res.data.map(item => ({
          code: item.code,
          name: item.name,
        }));
      }

      return [];
    } catch (error) {
      console.error(`Failed to load dictionary ${dictCode}:`, error);
      return [];
    }
  }

  /**
   * 批量预加载字典数据
   * @param {Array<string>} dictCodes - 字典类别数组
   */
  async preloadDictionaries(dictCodes) {
    const promises = dictCodes
      .filter(dc => dc && !this.cache.has(dc))
      .map(dc => this.getDictionary(dc));

    await Promise.all(promises);
  }

  /**
   * 清除指定字典的缓存
   * @param {string} dictCode - 字典类别
   */
  clearCache(dictCode) {
    if (dictCode) {
      this.cache.delete(dictCode);
    } else {
      // 清除所有缓存
      this.cache.clear();
    }
  }

  /**
   * 获取关联表数据(通用方法,支持非字典表)
   * @param {string} tableCode - 表代码
   * @param {Object} filter - 过滤条件
   * @param {string} valueField - 值字段名 (默认 code)
   * @param {string} labelField - 显示字段名 (默认 name)
   * @returns {Promise<Array>} 数据数组
   */
  async getTableData(tableCode, filter = {}, valueField = 'code', labelField = 'name') {
    if (!tableCode) {
      console.warn('getTableData: tableCode is required');
      return [];
    }

    // 生成缓存 key
    const cacheKey = `table:${tableCode}:${JSON.stringify(filter)}`;

    // 1. 检查缓存
    if (this.cache.has(cacheKey)) {
      return this.cache.get(cacheKey);
    }

    // 2. 检查是否正在加载
    if (this.loading.has(cacheKey)) {
      return this.loading.get(cacheKey);
    }

    // 3. 加载数据
    const loadPromise = this.loadTableData(tableCode, filter, valueField, labelField);
    this.loading.set(cacheKey, loadPromise);

    try {
      const data = await loadPromise;
      this.cache.set(cacheKey, data);
      return data;
    } finally {
      this.loading.delete(cacheKey);
    }
  }

  /**
   * 从后端加载关联表数据
   * @private
   */
  async loadTableData(tableCode, filter, valueField, labelField) {
    try {
      const res = await getEntityList({
        table_code: tableCode,
        ...filter,
        pageSize: 1000,
      });

      if (res && res.data) {
        return res.data.map(item => ({
          code: item[valueField],
          name: item[labelField],
        }));
      }

      return [];
    } catch (error) {
      console.error(`Failed to load table data ${tableCode}:`, error);
      return [];
    }
  }

  /**
   * 获取缓存统计信息
   */
  getCacheStats() {
    return {
      size: this.cache.size,
      keys: Array.from(this.cache.keys()),
    };
  }
}

// 导出单例
export default new DictionaryService();
