# OpenAPI 調用示例

本文檔提供了如何調用 PieMDM OpenAPI 的詳細示例，涵蓋了官方 Go SDK、手動簽名計算以及常用的命令行工具調用方式。

## 1. 使用官方 Go SDK

官方 Go SDK 位於項目源代碼的 `packages/go/openapi` 目錄下。它已封裝了簽名計算、請求重試以及錯誤處理邏輯。

### 1.1 安裝與依賴
如果您在 PieMDM 的 Monorepo 環境中開發，可以通過 `go.work` 直接引用。如果是跨項目引用，請確保項目可以訪問該 SDK 包。

```bash
go get github.com/pieworks/piemdm/packages/go/openapi
```

### 1.2 SDK 調用示例

```go
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/pieworks/piemdm/packages/go/openapi/client"
)

func main() {
	// 1. 初始化客戶端
	cli := client.NewClient(client.Config{
		BaseURL:   "http://localhost:8787", // API 基礎地址
		AppID:     "your_app_id",           // 應用 ID
		AppSecret: "your_app_secret",       // 應用密鑰
		Timeout:   10 * time.Second,        // 請求超時時間
	})

	// 2. 調用列表查詢接口 (以產品表為例)
	// 支持分頁參數 page, pageSize
	resp, err := cli.Get("/openapi/v1/entities/product?page=1&pageSize=10")
	if err != nil {
		fmt.Printf("請求失敗: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("列表響應: %s\n", string(body))

	// 3. 調用詳情查詢接口
	respDetail, err := cli.Get("/openapi/v1/entities/product/1")
	if err != nil {
		fmt.Printf("查詢詳情失敗: %v\n", err)
		return
	}
	defer respDetail.Body.Close()

	bodyDetail, _ := io.ReadAll(respDetail.Body)
	fmt.Printf("詳情響應: %s\n", string(bodyDetail))
}
```

## 2. 手動實現簽名計算

如果您使用其他編程語言（如 Java, Python, PHP），需要手動實現簽名邏輯。

### 2.1 簽名核心邏輯

簽名基於 **Canonical Request**。拼接規則如下：

```text
CanonicalRequest = 
    HTTPMethod + "\n" +
    Path + "\n" +
    SortedQueryString + "\n" +
    RequestBodyHash + "\n" +
    Timestamp + "\n" +
    Nonce
```

### 2.2 Python 示例

```python
import hashlib
import hmac
import time
import uuid
import requests

def compute_signature(secret, method, path, query_params, body, timestamp, nonce):
    # 1. 對查詢參數進行排序並拼接
    sorted_query = "&".join(f"{k}={v}" for k, v in sorted(query_params.items()))
    
    # 2. 對 Body 進行 SHA256 哈希 (空則為 "")
    body_hash = hashlib.sha256(body.encode('utf-8')).hexdigest() if body else ""
    
    # 3. 構建規範請求文本
    canonical_request = "\n".join([
        method.upper(),
        path,
        sorted_query,
        body_hash,
        str(timestamp),
        nonce
    ])
    
    # 4. HMAC-SHA256 計算簽名
    sign = hmac.new(secret.encode('utf-8'), canonical_request.encode('utf-8'), hashlib.sha256).hexdigest()
    return sign

# 使用示例
app_id = "test_app"
app_secret = "test_secret"
timestamp = int(time.time())
nonce = str(uuid.uuid4())
method = "GET"
path = "/openapi/v1/entities/product"
query = {"page": "1", "pageSize": "10"}

signature = compute_signature(app_secret, method, path, query, "", timestamp, nonce)

headers = {
    "X-App-Id": app_id,
    "X-Timestamp": str(timestamp),
    "X-Nonce": nonce,
    "X-Sign": signature
}

url = f"http://localhost:8787{path}?page=1&pageSize=10"
response = requests.get(url, headers=headers)
print(response.json())
```

## 3. 使用 cURL 調用

雖然 cURL 難以直接生成動態簽名，但在聯調時，您可以先通過 SDK 生成一個有效的簽名，然後使用 cURL 進行重放（注意要在 `X-Timestamp` 窗口期內，且 `X-Nonce` 未被再次使用）。

```bash
curl -X GET "http://localhost:8787/openapi/v1/entities/product?page=1&pageSize=10" \
     -H "X-App-Id: test_app_001" \
     -H "X-Timestamp: 1705478400" \
     -H "X-Nonce: abc123xyz789" \
     -H "X-Sign: a1b2c3d4e5f6g7h8..."
```

## 4. 常見問題排查

- **簽名不匹配**：請檢查 `CanonicalRequest` 的拼接順序是否準確，尤其是換行符 `\n`。
- **時間戳報錯**：確保本地時間與服務器時間同步，允許的最大誤差為 5 分鐘。
- **Nonce 被佔用**：Nonce 在 10 分鐘內不可重複，請確保每次請求生成全新的 UUID。
- **IP 禁止訪問**：請檢查管理後台中該 AppID 是否配置了限制 IP，且您的出口 IP 在白名單內。
