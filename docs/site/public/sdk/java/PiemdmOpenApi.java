package com.piemdm.openapi;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.*;

import com.google.gson.Gson;
import com.google.gson.JsonObject;

/**
 * PieMDM OpenAPI Java SDK
 *
 * @version 1.0.0
 */
public class PiemdmOpenApi {

    /**
     * 空 Body 的 SHA256 哈希值(固定值)
     * SHA256("") = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
     */
    private static final String EMPTY_PAYLOAD_HASH = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";

    private final String appId;
    private final String appSecret;
    private final String baseUrl;
    private final int timeout;
    private final Gson gson;

    /**
     * 构造函数
     *
     * @param appId 应用 ID
     * @param appSecret 应用密钥
     * @param baseUrl API 基础 URL
     */
    public PiemdmOpenApi(String appId, String appSecret, String baseUrl) {
        this(appId, appSecret, baseUrl, 30000);
    }

    /**
     * 构造函数
     *
     * @param appId 应用 ID
     * @param appSecret 应用密钥
     * @param baseUrl API 基础 URL
     * @param timeout 超时时间(毫秒)
     */
    public PiemdmOpenApi(String appId, String appSecret, String baseUrl, int timeout) {
        this.appId = appId;
        this.appSecret = appSecret;
        this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1) : baseUrl;
        this.timeout = timeout;
        this.gson = new Gson();
    }

    /**
     * 构建规范请求字符串
     */
    private String buildCanonicalRequest(String method, String path, Map<String, String> query,
                                        String body, String timestamp, String nonce) {
        StringBuilder canonical = new StringBuilder();
        canonical.append(method).append("\n");
        canonical.append(path).append("\n");

        // 查询参数按字典序排序
        if (query != null && !query.isEmpty()) {
            List<String> keys = new ArrayList<>(query.keySet());
            Collections.sort(keys);

            List<String> params = new ArrayList<>();
            for (String key : keys) {
                params.add(key + "=" + query.get(key));
            }
            canonical.append(String.join("&", params));
        }
        canonical.append("\n");

        // 请求体 SHA256 哈希
        String bodyHash;
        if (body == null || body.isEmpty()) {
            bodyHash = EMPTY_PAYLOAD_HASH;
        } else {
            try {
                MessageDigest digest = MessageDigest.getInstance("SHA-256");
                byte[] hash = digest.digest(body.getBytes(StandardCharsets.UTF_8));
                StringBuilder hexString = new StringBuilder();
                for (byte b : hash) {
                    String hex = Integer.toHexString(0xff & b);
                    if (hex.length() == 1) hexString.append('0');
                    hexString.append(hex);
                }
                bodyHash = hexString.toString();
            } catch (Exception e) {
                throw new RuntimeException("Failed to compute SHA-256", e);
            }
        }
        canonical.append(bodyHash).append("\n");
        canonical.append(timestamp).append("\n");
        canonical.append(nonce);

        return canonical.toString();
    }

    /**
     * 计算 HMAC-SHA256 签名
     */
    private String computeSignature(String canonicalRequest) throws Exception {
        Mac mac = Mac.getInstance("HmacSHA256");
        SecretKeySpec secretKeySpec = new SecretKeySpec(appSecret.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
        mac.init(secretKeySpec);
        byte[] hash = mac.doFinal(canonicalRequest.getBytes(StandardCharsets.UTF_8));

        StringBuilder hexString = new StringBuilder();
        for (byte b : hash) {
            String hex = Integer.toHexString(0xff & b);
            if (hex.length() == 1) hexString.append('0');
            hexString.append(hex);
        }
        return hexString.toString();
    }

    /**
     * 发送 HTTP 请求
     */
    private JsonObject request(String method, String path, Map<String, String> query, Map<String, Object> data)
            throws Exception {
        String timestamp = String.valueOf(System.currentTimeMillis() / 1000);
        String nonce = UUID.randomUUID().toString().replace("-", "");
        String body = data != null ? gson.toJson(data) : "";

        // 构建规范请求
        String canonicalRequest = buildCanonicalRequest(method, path, query, body, timestamp, nonce);

        // 计算签名
        String signature = computeSignature(canonicalRequest);

        // 构建完整 URL
        StringBuilder urlBuilder = new StringBuilder(baseUrl).append(path);
        if (query != null && !query.isEmpty()) {
            urlBuilder.append("?");
            List<String> params = new ArrayList<>();
            for (Map.Entry<String, String> entry : query.entrySet()) {
                params.add(entry.getKey() + "=" + entry.getValue());
            }
            urlBuilder.append(String.join("&", params));
        }

        URL url = new URL(urlBuilder.toString());
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod(method);
        conn.setConnectTimeout(timeout);
        conn.setReadTimeout(timeout);
        conn.setRequestProperty("Content-Type", "application/json");
        conn.setRequestProperty("X-App-Id", appId);
        conn.setRequestProperty("X-Timestamp", timestamp);
        conn.setRequestProperty("X-Nonce", nonce);
        conn.setRequestProperty("X-Sign", signature);

        if (body != null && !body.isEmpty()) {
            conn.setDoOutput(true);
            try (OutputStream os = conn.getOutputStream()) {
                os.write(body.getBytes(StandardCharsets.UTF_8));
            }
        }

        // 读取响应
        int responseCode = conn.getResponseCode();
        BufferedReader in = new BufferedReader(
            new InputStreamReader(
                responseCode >= 400 ? conn.getErrorStream() : conn.getInputStream(),
                StandardCharsets.UTF_8
            )
        );

        StringBuilder response = new StringBuilder();
        String line;
        while ((line = in.readLine()) != null) {
            response.append(line);
        }
        in.close();

        JsonObject result = gson.fromJson(response.toString(), JsonObject.class);

        if (result.get("code").getAsInt() != 200) {
            throw new Exception("API error: " + result.get("message").getAsString());
        }

        return result;
    }

    /**
     * 查询实体列表
     */
    public JsonObject list(String table, Map<String, String> params) throws Exception {
        String path = "/openapi/v1/entities/" + table;
        return request("GET", path, params, null);
    }

    /**
     * 查询实体详情
     */
    public JsonObject get(String table, int id) throws Exception {
        String path = "/openapi/v1/entities/" + table + "/" + id;
        return request("GET", path, null, null);
    }

    /**
     * 创建实体 (Phase 2)
     */
    public JsonObject create(String table, Map<String, Object> data) throws Exception {
        String path = "/openapi/v1/entities/" + table;
        return request("POST", path, null, data);
    }

    /**
     * 更新实体 (Phase 2)
     */
    public JsonObject update(String table, int id, Map<String, Object> data) throws Exception {
        String path = "/openapi/v1/entities/" + table + "/" + id;
        return request("PUT", path, null, data);
    }

    /**
     * 删除实体 (Phase 2)
     */
    public JsonObject delete(String table, int id) throws Exception {
        String path = "/openapi/v1/entities/" + table + "/" + id;
        return request("DELETE", path, null, null);
    }

    /**
     * 使用示例
     */
    public static void main(String[] args) {
        try {
            PiemdmOpenApi client = new PiemdmOpenApi(
                "test_app_001",
                "test_secret_123456",
                "https://your-domain.com"
            );

            // 列表查询
            Map<String, String> params = new HashMap<>();
            params.put("page", "1");
            params.put("pageSize", "10");
            JsonObject result = client.list("product", params);
            System.out.println(result);

            // 详情查询
            JsonObject detail = client.get("product", 1);
            System.out.println(detail);

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
