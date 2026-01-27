# OpenAPI 调用示例

本文档提供了如何调用 PieMDM OpenAPI 的详细示例，涵盖了官方 Go SDK、手动签名计算以及常用的命令行工具调用方式。

## 1. 使用官方 Go SDK

官方 Go SDK 位于项目源代码的 `packages/go/openapi` 目录下。它已封装了签名计算、请求重试以及错误处理逻辑。

### 1.1 安装与依赖
如果您在 PieMDM 的 Monorepo 环境中开发，可以通过 `go.work` 直接引用。如果是跨项目引用，请确保项目可以访问该 SDK 包。

```bash
go get github.com/pieteams/piemdm/packages/go/openapi
```

### 1.2 SDK 调用示例

```go
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/pieteams/piemdm/packages/go/openapi/client"
)

func main() {
	// 1. 初始化客户端
	cli := client.NewClient(client.Config{
		BaseURL:   "http://localhost:8787", // API 基础地址
		AppID:     "your_app_id",           // 应用 ID
		AppSecret: "your_app_secret",       // 应用密钥
		Timeout:   10 * time.Second,        // 请求超时时间
	})

	// 2. 调用列表查询接口 (以产品表为例)
	// 支持分页参数 page, pageSize
	resp, err := cli.Get("/openapi/v1/entities/product?page=1&pageSize=10")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("列表响应: %s\n", string(body))

	// 3. 调用详情查询接口
	respDetail, err := cli.Get("/openapi/v1/entities/product/1")
	if err != nil {
		fmt.Printf("查询详情失败: %v\n", err)
		return
	}
	defer respDetail.Body.Close()

	bodyDetail, _ := io.ReadAll(respDetail.Body)
	fmt.Printf("详情响应: %s\n", string(bodyDetail))
}
```

## 2. 手动实现签名计算

如果您使用其他编程语言（如 Java, Python, PHP），需要手动实现签名逻辑。

### 2.1 签名核心逻辑

签名基于 **Canonical Request**。拼接规则如下：

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
    # 1. 对查询参数进行排序并拼接
    sorted_query = "&".join(f"{k}={v}" for k, v in sorted(query_params.items()))
    
    # 2. 对 Body 进行 SHA256 哈希 (空则为 "")
    body_hash = hashlib.sha256(body.encode('utf-8')).hexdigest() if body else ""
    
    # 3. 构建规范请求文本
    canonical_request = "\n".join([
        method.upper(),
        path,
        sorted_query,
        body_hash,
        str(timestamp),
        nonce
    ])
    
    # 4. HMAC-SHA256 计算签名
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

## 3. 使用 cURL 调用

虽然 cURL 难以直接生成动态签名，但在联调时，您可以先通过 SDK 生成一个有效的签名，然后使用 cURL 进行重放（注意要在 `X-Timestamp` 窗口期内，且 `X-Nonce` 未被再次使用）。

```bash
curl -X GET "http://localhost:8787/openapi/v1/entities/product?page=1&pageSize=10" \
     -H "X-App-Id: test_app_001" \
     -H "X-Timestamp: 1705478400" \
     -H "X-Nonce: abc123xyz789" \
     -H "X-Sign: a1b2c3d4e5f6g7h8..."
```

## 4. 常见问题排查

- **签名不匹配**：请检查 `CanonicalRequest` 的拼接顺序是否准确，尤其是换行符 `\n`。
- **时间戳报错**：确保本地时间与服务器时间同步，允许的最大误差为 5 分钟。
- **Nonce 被占用**：Nonce 在 10 分钟内不可重复，请确保每次请求生成全新的 UUID。
- **IP 禁止访问**：请检查管理后台中该 AppID 是否配置了限制 IP，且您的出口 IP 在白名单内。
