# Open API Guide

PieMDM provides a powerful OpenAPI that allows external systems to access master data through secure authentication. The OpenAPI uses an HMAC-SHA256 signature mechanism based on Canonical Request to ensure the authenticity and integrity of requests.

## 1. Authentication Mechanism

All OpenAPI requests must carry the following authentication parameters in the Header:

| Header Parameter | Description | Example |
| :--- | :--- | :--- |
| `X-App-Id` | App Identifier, obtained in "App Management" | `app_592837482...` |
| `X-Timestamp` | Unix timestamp (seconds) when the request was sent | `1674829374` |
| `X-Nonce` | A random string of at least 16 characters used to prevent replay attacks | `abcdef1234567890` |
| `X-Sign` | Digest value calculated based on the signature algorithm (Hex encoded) | `a1b2c3d4...` |

### 1.1 Signature Algorithm

The signature calculation process is as follows:

1. **Construct Canonical Request String**:
   Format with parts connected by newline `\n`:
   ```text
   HTTPMethod + "\n" +
   CanonicalURI + "\n" +
   CanonicalQueryString + "\n" +
   SHA256(RequestBody) + "\n" +
   X-Timestamp + "\n" +
   X-Nonce
   ```
   - **HTTPMethod**: Uppercase HTTP method, such as `GET` or `POST`.
   - **CanonicalURI**: The full path of the request, must start with `/`, e.g., `/openapi/v1/entities/users`.
   - **CanonicalQueryString**: Sorted query parameter string. Sort by parameter name in ASCII ascending order, format is `key1=value1&key2=value2`.
   - **SHA256(RequestBody)**: SHA256 hash value (lowercase Hex) of the request body. If Body is empty, the hash is: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
   - **X-Timestamp** and **X-Nonce**: Must be strictly consistent with the values in the Header.

2. **Calculate Signature**:
   Use aggregate's `AppSecret` as the Key to perform HMAC-SHA256 calculation on the above canonical string:
   ```text
   Signature = Hex(HMAC-SHA256(AppSecret, CanonicalRequest))
   ```

## 2. Security Policy

In addition to signature verification, OpenAPI also provides multiple security guarantees:

- **IP Whitelist**: Only IP addresses configured in the "App Management" whitelist can successfully call the interface.
- **Permission Audit**: All OpenAPI calls are recorded in detail, including request parameters, response results, call duration, etc.
- **Nonce Anti-Replay**: Each Nonce can only be used once. The system caches recently used Nonces, and repeated submissions will return `401 Unauthorized`.
- **Time Window Verification**: Requests where the deviation between `X-Timestamp` and the server's current time exceeds 5 minutes will be rejected.

## 3. Core Interface

The base path for all interfaces is: `/openapi/v1`

### 3.1 Get Entity List

**GET** `/openapi/v1/entities/{table_code}`

**Query Parameters**:
- `page`: Page number, default 1
- `pageSize`: Number per page, default 15, max 100
- `id`: Optional, exact ID filter
- `status`: Optional, status filter
- `created_at`: Optional, creation time filter

### 3.2 Get Entity Details

**GET** `/openapi/v1/entities/{table_code}/{id}`

## 4. Error Code Reference

| Error Message | Description |
| :--- | :--- |
| `AUTH_FAILED` | `X-App-Id` invalid or missing |
| `SIGNATURE_INVALID` | Signature verification failed |
| `TOKEN_EXPIRED` | Timestamp timeout or Nonce has been used |
| `IP_NOT_ALLOWED` | Client IP is not in the whitelist |
| `PERMISSION_DENIED` | App has no permission to access the specified entity table |

<callout emoji="💡" background-color="light-blue" border-color="blue">
Tip: More interfaces (such as creating and modifying entities) are under development. If you have urgent needs, please contact technical support.
</callout>
