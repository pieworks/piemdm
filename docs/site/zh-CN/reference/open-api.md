# 开放接口 (OpenAPI) 指南

PieMDM 提供了强大的 OpenAPI，允许外部系统通过安全认证的方式访问主数据。OpenAPI 采用了基于规范请求（Canonical Request）的 HMAC-SHA256 签名机制，确保请求的真实性和完整性。

## 1. 认证机制

所有 OpenAPI 请求必须在 Header 中携带以下认证参数：

| Header 参数 | 说明 | 示例 |
| :--- | :--- | :--- |
| `X-App-Id` | 应用标识，在“应用管理”中获取 | `app_592837482...` |
| `X-Timestamp` | 请求发送时的 Unix 时间戳（秒） | `1674829374` |
| `X-Nonce` | 长度至少 16 位的随机字符串，用于防重放攻击 | `abcdef1234567890` |
| `X-Sign` | 根据签名算法计算出的摘要值（Hex 编码） | `a1b2c3d4...` |

### 1.1 签名算法

签名的计算过程如下：

1. **构造规范请求字符串 (Canonical Request)**：
   格式为各部分由换行符 `\n` 连接：
   ```text
   HTTPMethod + "\n" +
   CanonicalURI + "\n" +
   CanonicalQueryString + "\n" +
   SHA256(RequestBody) + "\n" +
   X-Timestamp + "\n" +
   X-Nonce
   ```
   - **HTTPMethod**: 大写的 HTTP 方法，如 `GET` 或 `POST`。
   - **CanonicalURI**: 请求的完整路径，必须以 `/` 开头，如 `/openapi/v1/entities/users`。
   - **CanonicalQueryString**: 排序后的查询参数字符串。按参数名 ASCII 升序排列，格式为 `key1=value1&key2=value2`。
   - **SHA256(RequestBody)**: 请求体的 SHA256 哈希值（小写 Hex）。若 Body 为空，则哈希值为：`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`。
   - **X-Timestamp** 和 **X-Nonce**: 与 Header 中的值保持严格一致。

2. **计算签名**：
   使用应用的 `AppSecret` 作为 Key，对上述规范字符串进行 HMAC-SHA256 计算：
   ```text
   Signature = Hex(HMAC-SHA256(AppSecret, CanonicalRequest))
   ```

## 2. 安全策略

除了签名验证，OpenAPI 还提供了多重安全保障：

- **IP 白名单**：只有在“应用管理”中配置的白名单 IP 才能成功调用接口。
- **权限审计**：所有 OpenAPI 的调用都会被详细记录，包括请求参数、响应结果、调用耗时等。
- **Nonce 防重放**：每个 Nonce 仅能使用一次。系统会缓存最近使用的 Nonce，重复提交将返回 `401 Unauthorized`。
- **时间窗口校验**：`X-Timestamp` 与服务器当前时间的偏差超过 5 分钟将被拒绝。

## 3. 核心接口

所有接口的基础路径为：`/openapi/v1`

### 3.1 获取实体列表

**GET** `/openapi/v1/entities/{table_code}`

**查询参数**：
- `page`: 页码，默认 1
- `pageSize`: 每页数量，默认 15，最大 100
- `id`: 可选，精确过滤 ID
- `status`: 可选，状态过滤
- `created_at`: 可选，创建时间过滤

### 3.2 获取实体详情

**GET** `/openapi/v1/entities/{table_code}/{id}`

## 4. 错误码参考

| 错误消息 | 说明 |
| :--- | :--- |
| `AUTH_FAILED` | `X-App-Id` 无效或缺失 |
| `SIGNATURE_INVALID` | 签名验证失败 |
| `TOKEN_EXPIRED` | 时间戳超时或 Nonce 已被使用 |
| `IP_NOT_ALLOWED` | 客户端 IP 不在白名单中 |
| `PERMISSION_DENIED` | 应用无权访问指定的实体表 |

<callout emoji="💡" background-color="light-blue" border-color="blue">
提示：更多接口（如创建、修改实体）正在开发中。如有紧急需求，请联系技术支持。
</callout>
