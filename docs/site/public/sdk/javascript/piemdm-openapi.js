/**
 * PieMDM OpenAPI JavaScript SDK
 *
 * @version 1.0.0
 * @author PieMDM Team
 */

const crypto = require('crypto');

/**
 * 空 Body 的 SHA256 哈希值(固定值)
 * SHA256("") = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
 */
const EMPTY_PAYLOAD_HASH = 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855';

class PiemdmOpenApi {
  /**
   * 构造函数
   *
   * @param {string} appId - 应用 ID
   * @param {string} appSecret - 应用密钥
   * @param {string} baseUrl - API 基础 URL
   * @param {number} timeout - 超时时间(毫秒),默认 30000
   */
  constructor(appId, appSecret, baseUrl, timeout = 30000) {
    this.appId = appId;
    this.appSecret = appSecret;
    this.baseUrl = baseUrl.replace(/\/$/, '');
    this.timeout = timeout;
  }

  /**
   * 构建规范请求字符串
   *
   * @private
   */
  buildCanonicalRequest(method, path, query, body, timestamp, nonce) {
    let canonical = `${method}\n`;
    canonical += `${path}\n`;

    // 查询参数按字典序排序
    if (query && Object.keys(query).length > 0) {
      const sortedKeys = Object.keys(query).sort();
      const params = sortedKeys.map(key => `${key}=${query[key]}`);
      canonical += params.join('&');
    }
    canonical += '\n';

    // 请求体 SHA256 哈希
    let bodyHash;
    if (!body || body === '') {
      bodyHash = EMPTY_PAYLOAD_HASH;
    } else {
      bodyHash = crypto.createHash('sha256').update(body).digest('hex');
    }
    canonical += `${bodyHash}\n`;
    canonical += `${timestamp}\n`;
    canonical += nonce;

    return canonical;
  }

  /**
   * 计算 HMAC-SHA256 签名
   *
   * @private
   */
  computeSignature(canonicalRequest) {
    const hmac = crypto.createHmac('sha256', this.appSecret);
    hmac.update(canonicalRequest);
    return hmac.digest('hex');
  }

  /**
   * 生成随机 nonce
   *
   * @private
   */
  generateNonce() {
    return crypto.randomBytes(16).toString('hex');
  }

  /**
   * 发送 HTTP 请求
   *
   * @private
   */
  async request(method, path, query = {}, data = null) {
    const timestamp = Math.floor(Date.now() / 1000).toString();
    const nonce = this.generateNonce();
    const body = data ? JSON.stringify(data) : '';

    // 构建规范请求
    const canonicalRequest = this.buildCanonicalRequest(
      method,
      path,
      query,
      body,
      timestamp,
      nonce
    );

    // 计算签名
    const signature = this.computeSignature(canonicalRequest);

    // 构建完整 URL
    let url = this.baseUrl + path;
    if (Object.keys(query).length > 0) {
      const queryString = Object.entries(query)
        .map(([key, value]) => `${key}=${encodeURIComponent(value)}`)
        .join('&');
      url += '?' + queryString;
    }

    // 设置请求选项
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json',
        'X-App-Id': this.appId,
        'X-Timestamp': timestamp,
        'X-Nonce': nonce,
        'X-Sign': signature,
      },
    };

    if (body) {
      options.body = body;
    }

    // 发送请求
    try {
      const response = await fetch(url, options);
      const result = await response.json();

      if (result.code !== 200) {
        throw new Error(`API error: ${result.message}`);
      }

      return result;
    } catch (error) {
      throw new Error(`Request failed: ${error.message}`);
    }
  }

  /**
   * 查询实体列表
   *
   * @param {string} table - 表名
   * @param {Object} params - 查询参数
   * @returns {Promise<Object>}
   */
  async list(table, params = {}) {
    const path = `/openapi/v1/entities/${table}`;
    return this.request('GET', path, params);
  }

  /**
   * 查询实体详情
   *
   * @param {string} table - 表名
   * @param {number} id - 记录 ID
   * @returns {Promise<Object>}
   */
  async get(table, id) {
    const path = `/openapi/v1/entities/${table}/${id}`;
    return this.request('GET', path);
  }

  /**
   * 创建实体 (Phase 2)
   *
   * @param {string} table - 表名
   * @param {Object} data - 数据
   * @returns {Promise<Object>}
   */
  async create(table, data) {
    const path = `/openapi/v1/entities/${table}`;
    return this.request('POST', path, {}, data);
  }

  /**
   * 更新实体 (Phase 2)
   *
   * @param {string} table - 表名
   * @param {number} id - 记录 ID
   * @param {Object} data - 数据
   * @returns {Promise<Object>}
   */
  async update(table, id, data) {
    const path = `/openapi/v1/entities/${table}/${id}`;
    return this.request('PUT', path, {}, data);
  }

  /**
   * 删除实体 (Phase 2)
   *
   * @param {string} table - 表名
   * @param {number} id - 记录 ID
   * @returns {Promise<Object>}
   */
  async delete(table, id) {
    const path = `/openapi/v1/entities/${table}/${id}`;
    return this.request('DELETE', path);
  }
}

// Node.js 环境导出
if (typeof module !== 'undefined' && module.exports) {
  module.exports = PiemdmOpenApi;
}

// 浏览器环境导出
if (typeof window !== 'undefined') {
  window.PiemdmOpenApi = PiemdmOpenApi;
}

// 使用示例
/*
// Node.js 环境
const PiemdmOpenApi = require('./piemdm-openapi');

const client = new PiemdmOpenApi(
  'test_app_001',
  'test_secret_123456',
  'https://your-domain.com'
);

// 列表查询
client.list('product', { page: 1, pageSize: 10 })
  .then(result => console.log(result))
  .catch(error => console.error(error));

// 详情查询
client.get('product', 1)
  .then(result => console.log(result))
  .catch(error => console.error(error));

// 使用 async/await
async function example() {
  try {
    const result = await client.list('product', { page: 1, pageSize: 10 });
    console.log(result);
  } catch (error) {
    console.error(error);
  }
}
*/
