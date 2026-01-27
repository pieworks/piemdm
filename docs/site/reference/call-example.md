# OpenAPI Call Examples

This document provides detailed examples of how to call PieMDM OpenAPI, covering the official Go SDK, manual signature calculation, and common command-line tool usage.

## 1. Using the Official Go SDK

The official Go SDK is located in the `packages/go/openapi` directory of the project source code. It encapsulates signature calculation, request retries, and error handling logic.

### 1.1 Installation & Dependencies
If you are developing in the PieMDM Monorepo environment, you can reference it directly via `go.work`. If referencing across projects, please ensure the project can access the SDK package.

```bash
go get github.com/pieteams/piemdm/packages/go/openapi
```

### 1.2 SDK Usage Example

```go
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/pieteams/piemdm/packages/go/openapi/client"
)

func main() {
	// 1. Initialize Client
	cli := client.NewClient(client.Config{
		BaseURL:   "http://localhost:8787", // API Base URL
		AppID:     "your_app_id",           // App ID
		AppSecret: "your_app_secret",       // App Secret
		Timeout:   10 * time.Second,        // Request Timeout
	})

	// 2. Call List Query Interface (Using Product Table as Example)
	// Supports pagination parameters page, pageSize
	resp, err := cli.Get("/openapi/v1/entities/product?page=1&pageSize=10")
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("List Response: %s\n", string(body))

	// 3. Call Detail Query Interface
	respDetail, err := cli.Get("/openapi/v1/entities/product/1")
	if err != nil {
		fmt.Printf("Query Detail Failed: %v\n", err)
		return
	}
	defer respDetail.Body.Close()

	bodyDetail, _ := io.ReadAll(respDetail.Body)
	fmt.Printf("Detail Response: %s\n", string(bodyDetail))
}
```

## 2. Manual Signature Calculation

If you use other programming languages (such as Java, Python, PHP), you need to manually implement the signature logic.

### 2.1 Core Signature Logic

The signature is based on **Canonical Request**. The concatenation rule is as follows:

```text
CanonicalRequest = 
    HTTPMethod + "\n" +
    Path + "\n" +
    SortedQueryString + "\n" +
    RequestBodyHash + "\n" +
    Timestamp + "\n" +
    Nonce
```

### 2.2 Python Example

```python
import hashlib
import hmac
import time
import uuid
import requests

def compute_signature(secret, method, path, query_params, body, timestamp, nonce):
    # 1. Sort and concatenate query parameters
    sorted_query = "&".join(f"{k}={v}" for k, v in sorted(query_params.items()))
    
    # 2. SHA256 hash the Body (empty string if empty)
    body_hash = hashlib.sha256(body.encode('utf-8')).hexdigest() if body else ""
    
    # 3. Construct Canonical Request text
    canonical_request = "\n".join([
        method.upper(),
        path,
        sorted_query,
        body_hash,
        str(timestamp),
        nonce
    ])
    
    # 4. Calculate HMAC-SHA256 signature
    sign = hmac.new(secret.encode('utf-8'), canonical_request.encode('utf-8'), hashlib.sha256).hexdigest()
    return sign

# Usage Example
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

## 3. Using cURL

Although it is difficult to generate dynamic signatures directly with cURL, during debugging, you can first generate a valid signature via the SDK and then use cURL for replay (note that it must be within the `X-Timestamp` window and `X-Nonce` must not be reused).

```bash
curl -X GET "http://localhost:8787/openapi/v1/entities/product?page=1&pageSize=10" \
     -H "X-App-Id: test_app_001" \
     -H "X-Timestamp: 1705478400" \
     -H "X-Nonce: abc123xyz789" \
     -H "X-Sign: a1b2c3d4e5f6g7h8..."
```

## 4. Common Troubleshooting

- **Signature Mismatch**: Please check if the concatenation order of `CanonicalRequest` is accurate, especially the newline character `\n`.
- **Timestamp Error**: Ensure strict synchronization between local time and server time, with a maximum allowable error of 5 minutes.
- **Nonce Used**: Nonces cannot be repeated within 10 minutes. Please ensure a fresh UUID is generated for each request.
- **IP Forbidden**: Check in the management backend if the AppID has IP restrictions configured and if your exit IP is in the whitelist.
