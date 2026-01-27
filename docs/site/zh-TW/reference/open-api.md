# 開放接口 (OpenAPI) 指南

PieMDM 提供了強大的 OpenAPI，允許外部系統通過安全認證的方式訪問主數據。OpenAPI 採用了基於規範請求（Canonical Request）的 HMAC-SHA256 簽名機制，確保請求的真實性和完整性。

## 1. 認證機制

所有 OpenAPI 請求必須在 Header 中攜帶以下認證參數：

| Header 參數 | 說明 | 示例 |
| :--- | :--- | :--- |
| `X-App-Id` | 應用標識，在“應用管理”中獲取 | `app_592837482...` |
| `X-Timestamp` | 請求發送時的 Unix 時間戳（秒） | `1674829374` |
| `X-Nonce` | 長度至少 16 位的隨機字符串，用於防重放攻擊 | `abcdef1234567890` |
| `X-Sign` | 根據簽名算法計算出的摘要值（Hex 編碼） | `a1b2c3d4...` |

### 1.1 簽名算法

簽名的計算過程如下：

1. **構造規範請求字符串 (Canonical Request)**：
   格式為各部分由換行符 `\n` 連接：
   ```text
   HTTPMethod + "\n" +
   CanonicalURI + "\n" +
   CanonicalQueryString + "\n" +
   SHA256(RequestBody) + "\n" +
   X-Timestamp + "\n" +
   X-Nonce
   ```
   - **HTTPMethod**: 大寫的 HTTP 方法，如 `GET` 或 `POST`。
   - **CanonicalURI**: 請求的完整路徑，必須以 `/` 開頭，如 `/openapi/v1/entities/users`。
   - **CanonicalQueryString**: 排序後的查詢參數字符串。按參數名 ASCII 升序排列，格式為 `key1=value1&key2=value2`。
   - **SHA256(RequestBody)**: 請求體的 SHA256 哈希值（小寫 Hex）。若 Body 為空，則哈希值為：`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`。
   - **X-Timestamp** 和 **X-Nonce**: 與 Header 中的值保持嚴格一致。

2. **計算簽名**：
   使用應用的 `AppSecret` 作為 Key，對上述規範字符串進行 HMAC-SHA256 計算：
   ```text
   Signature = Hex(HMAC-SHA256(AppSecret, CanonicalRequest))
   ```

## 2. 安全策略

除了簽名驗證，OpenAPI 還提供了多重安全保障：

- **IP 白名單**：只有在“應用管理”中配置的白名單 IP 才能成功調用接口。
- **權限審計**：所有 OpenAPI 的調用都會被詳細記錄，包括請求參數、響應結果、調用耗時等。
- **Nonce 防重放**：每個 Nonce 僅能使用一次。系統會緩存最近使用的 Nonce，重複提交將返回 `401 Unauthorized`。
- **時間窗口校驗**：`X-Timestamp` 與服務器當前時間的偏差超過 5 分鐘將被拒絕。

## 3. 核心接口

所有接口的基礎路徑為：`/openapi/v1`

### 3.1 獲取實體列表

**GET** `/openapi/v1/entities/{table_code}`

**查詢參數**：
- `page`: 頁碼，默認 1
- `pageSize`: 每頁數量，默認 15，最大 100
- `id`: 可選，精確過濾 ID
- `status`: 可選，狀態過濾
- `created_at`: 可選，創建時間過濾

### 3.2 獲取實體詳情

**GET** `/openapi/v1/entities/{table_code}/{id}`

## 4. 錯誤碼參考

| 錯誤消息 | 說明 |
| :--- | :--- |
| `AUTH_FAILED` | `X-App-Id` 無效或缺失 |
| `SIGNATURE_INVALID` | 簽名驗證失敗 |
| `TOKEN_EXPIRED` | 時間戳超時或 Nonce 已被使用 |
| `IP_NOT_ALLOWED` | 客戶端 IP 不在白名單中 |
| `PERMISSION_DENIED` | 應用無權訪問指定的實體表 |

<callout emoji="💡" background-color="light-blue" border-color="blue">
提示：更多接口（如創建、修改實體）正在開發中。如有緊急需求，請聯繫技術支持。
</callout>
