<?php

namespace PieMDM\OpenAPI;

/**
 * PieMDM OpenAPI PHP SDK
 *
 * @package PieMDM\OpenAPI
 * @version 1.0.0
 */
class Client
{
    /**
     * 空 Body 的 SHA256 哈希值(固定值)
     * SHA256("") = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
     */
    const EMPTY_PAYLOAD_HASH = 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855';

    private string $appId;
    private string $appSecret;
    private string $baseUrl;
    private int $timeout;

    /**
     * 构造函数
     *
     * @param string $appId 应用 ID
     * @param string $appSecret 应用密钥
     * @param string $baseUrl API 基础 URL
     * @param int $timeout 超时时间(秒)
     */
    public function __construct(string $appId, string $appSecret, string $baseUrl, int $timeout = 30)
    {
        $this->appId = $appId;
        $this->appSecret = $appSecret;
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->timeout = $timeout;
    }

    /**
     * 构建规范请求字符串
     *
     * @param string $method HTTP 方法
     * @param string $path 请求路径
     * @param array $query 查询参数
     * @param string $body 请求体
     * @param string $timestamp 时间戳
     * @param string $nonce Nonce
     * @return string
     */
    private function buildCanonicalRequest(
        string $method,
        string $path,
        array $query,
        string $body,
        string $timestamp,
        string $nonce
    ): string {
        $canonicalRequest = $method . "\n";
        $canonicalRequest .= $path . "\n";

        // 查询参数按字典序排序
        if (!empty($query)) {
            ksort($query);
            $params = [];
            foreach ($query as $key => $value) {
                $params[] = $key . '=' . $value;
            }
            $canonicalRequest .= implode('&', $params);
        }
        $canonicalRequest .= "\n";

        // 请求体 SHA256 哈希
        if (empty($body)) {
            $bodyHash = self::EMPTY_PAYLOAD_HASH;
        } else {
            $bodyHash = hash('sha256', $body);
        }
        $canonicalRequest .= $bodyHash . "\n";
        $canonicalRequest .= $timestamp . "\n";
        $canonicalRequest .= $nonce;

        return $canonicalRequest;
    }

    /**
     * 计算 HMAC-SHA256 签名
     *
     * @param string $canonicalRequest 规范请求字符串
     * @return string
     */
    private function computeSignature(string $canonicalRequest): string
    {
        return hash_hmac('sha256', $canonicalRequest, $this->appSecret);
    }

    /**
     * 发送 HTTP 请求
     *
     * @param string $method HTTP 方法
     * @param string $path 请求路径
     * @param array $query 查询参数
     * @param array|null $data 请求数据
     * @return array
     * @throws \Exception
     */
    private function request(string $method, string $path, array $query = [], ?array $data = null): array
    {
        $timestamp = (string)time();
        $nonce = bin2hex(random_bytes(16));
        $body = $data ? json_encode($data) : '';

        // 构建规范请求
        $canonicalRequest = $this->buildCanonicalRequest($method, $path, $query, $body, $timestamp, $nonce);

        // 计算签名
        $signature = $this->computeSignature($canonicalRequest);

        // 构建完整 URL
        $url = $this->baseUrl . $path;
        if (!empty($query)) {
            $url .= '?' . http_build_query($query);
        }

        // 设置请求头
        $headers = [
            'Content-Type: application/json',
            'X-App-Id: ' . $this->appId,
            'X-Timestamp: ' . $timestamp,
            'X-Nonce: ' . $nonce,
            'X-Sign: ' . $signature,
        ];

        // 发送请求
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, $this->timeout);

        if ($body) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
        }

        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $error = curl_error($ch);
        curl_close($ch);

        if ($error) {
            throw new \Exception("cURL error: $error");
        }

        $result = json_decode($response, true);
        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new \Exception("Failed to parse JSON response: " . json_last_error_msg());
        }

        if ($result['code'] !== 200) {
            throw new \Exception("API error: " . $result['message']);
        }

        return $result;
    }

    /**
     * 查询实体列表
     *
     * @param string $table 表名
     * @param array $params 查询参数
     * @return array
     * @throws \Exception
     */
    public function list(string $table, array $params = []): array
    {
        $path = "/openapi/v1/entities/$table";
        return $this->request('GET', $path, $params);
    }

    /**
     * 查询实体详情
     *
     * @param string $table 表名
     * @param int $id 记录 ID
     * @return array
     * @throws \Exception
     */
    public function get(string $table, int $id): array
    {
        $path = "/openapi/v1/entities/$table/$id";
        return $this->request('GET', $path);
    }

    /**
     * 创建实体 (Phase 2)
     *
     * @param string $table 表名
     * @param array $data 数据
     * @return array
     * @throws \Exception
     */
    public function create(string $table, array $data): array
    {
        $path = "/openapi/v1/entities/$table";
        return $this->request('POST', $path, [], $data);
    }

    /**
     * 更新实体 (Phase 2)
     *
     * @param string $table 表名
     * @param int $id 记录 ID
     * @param array $data 数据
     * @return array
     * @throws \Exception
     */
    public function update(string $table, int $id, array $data): array
    {
        $path = "/openapi/v1/entities/$table/$id";
        return $this->request('PUT', $path, [], $data);
    }

    /**
     * 删除实体 (Phase 2)
     *
     * @param string $table 表名
     * @param int $id 记录 ID
     * @return array
     * @throws \Exception
     */
    public function delete(string $table, int $id): array
    {
        $path = "/openapi/v1/entities/$table/$id";
        return $this->request('DELETE', $path);
    }
}

// 使用示例
/*
$client = new \PieMDM\OpenAPI\Client(
    'test_app_001',
    'test_secret_123456',
    'https://your-domain.com'
);

try {
    // 列表查询
    $result = $client->list('product', [
        'page' => 1,
        'pageSize' => 10
    ]);
    print_r($result);

    // 详情查询
    $detail = $client->get('product', 1);
    print_r($detail);

} catch (\Exception $e) {
    echo "Error: " . $e->getMessage();
}
*/
