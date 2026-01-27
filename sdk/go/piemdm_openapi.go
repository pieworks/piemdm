package openapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EMPTY_PAYLOAD_HASH 空 Body 的 SHA256 哈希值(固定值)
// SHA256("") = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
const EMPTY_PAYLOAD_HASH = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Client OpenAPI 客户端
type Client struct {
	AppID      string
	AppSecret  string
	BaseURL    string
	HTTPClient *http.Client
}

// Response API 响应结构
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Total   int64           `json:"total,omitempty"`
}

// NewClient 创建新的 OpenAPI 客户端
func NewClient(appID, appSecret, baseURL string) *Client {
	return &Client{
		AppID:     appID,
		AppSecret: appSecret,
		BaseURL:   strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// buildCanonicalRequest 构建规范请求字符串
func (c *Client) buildCanonicalRequest(method, path string, query url.Values, body []byte, timestamp, nonce string) string {
	// 1. HTTP 方法
	canonicalRequest := method + "\n"

	// 2. 路径
	canonicalRequest += path + "\n"

	// 3. 查询参数(按字典序排序)
	if len(query) > 0 {
		keys := make([]string, 0, len(query))
		for k := range query {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		params := make([]string, 0, len(keys))
		for _, k := range keys {
			params = append(params, fmt.Sprintf("%s=%s", k, query.Get(k)))
		}
		canonicalRequest += strings.Join(params, "&")
	}
	canonicalRequest += "\n"

	// 4. 请求体 SHA256 哈希
	var bodyHash string
	if len(body) == 0 {
		// 空 Body 使用固定哈希值
		bodyHash = EMPTY_PAYLOAD_HASH
	} else {
		// 计算 Body 的 SHA256
		h := sha256.New()
		h.Write(body)
		bodyHash = hex.EncodeToString(h.Sum(nil))
	}
	canonicalRequest += bodyHash + "\n"

	// 5. 时间戳
	canonicalRequest += timestamp + "\n"

	// 6. Nonce
	canonicalRequest += nonce

	return canonicalRequest
}

// computeSignature 计算 HMAC-SHA256 签名
func (c *Client) computeSignature(canonicalRequest string) string {
	h := hmac.New(sha256.New, []byte(c.AppSecret))
	h.Write([]byte(canonicalRequest))
	return hex.EncodeToString(h.Sum(nil))
}

// request 发送 HTTP 请求
func (c *Client) request(method, path string, query url.Values, body []byte) (*Response, error) {
	// 生成时间戳和 nonce
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := uuid.New().String()

	// 构建规范请求
	canonicalRequest := c.buildCanonicalRequest(method, path, query, body, timestamp, nonce)

	// 计算签名
	signature := c.computeSignature(canonicalRequest)

	// 构建完整 URL
	fullURL := c.BaseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, fullURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Id", c.AppID)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Sign", signature)

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查错误
	if apiResp.Code != 200 {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	return &apiResp, nil
}

// List 查询实体列表
func (c *Client) List(table string, params map[string]string) (*Response, error) {
	path := fmt.Sprintf("/openapi/v1/entities/%s", table)

	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	return c.request("GET", path, query, nil)
}

// Get 查询实体详情
func (c *Client) Get(table string, id int) (*Response, error) {
	path := fmt.Sprintf("/openapi/v1/entities/%s/%d", table, id)
	return c.request("GET", path, nil, nil)
}

// Create 创建实体 (Phase 2)
func (c *Client) Create(table string, data map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/openapi/v1/entities/%s", table)

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	return c.request("POST", path, nil, body)
}

// Update 更新实体 (Phase 2)
func (c *Client) Update(table string, id int, data map[string]interface{}) (*Response, error) {
	path := fmt.Sprintf("/openapi/v1/entities/%s/%d", table, id)

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	return c.request("PUT", path, nil, body)
}

// Delete 删除实体 (Phase 2)
func (c *Client) Delete(table string, id int) (*Response, error) {
	path := fmt.Sprintf("/openapi/v1/entities/%s/%d", table, id)
	return c.request("DELETE", path, nil, nil)
}
