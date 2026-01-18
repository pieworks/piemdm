# OpenAPI 集成测试说明

## 概述

本目录包含 OpenAPI Phase 1 的集成测试,基于 `test_openapi.sh` 脚本整理而成。

## 测试文件

- `openapi_test.go` - OpenAPI 集成测试主文件

## 测试用例

### 1. List 接口测试
- ✅ 正常的 List 请求
- ✅ 带分页参数的 List 请求

### 2. Get 接口测试
- ✅ 正常的 Get 请求

### 3. 安全测试
- ✅ 错误的签名验证
- ✅ 缺少必需的请求头
- ✅ Nonce 重放攻击防护

## 运行测试

### 使用 Shell 脚本 (推荐用于快速验证)

```bash
cd backend
./test_openapi.sh
```

### 使用 Go 测试 (需要完整的测试环境)

```bash
cd backend
go test -v ./test/integration/openapi_test.go
```

**注意**: Go 集成测试目前需要完整的测试环境初始化,包括:
- 数据库连接
- Redis 连接
- 测试数据准备
- 完整的依赖注入

## 测试环境要求

1. **数据库**
   - MySQL 服务运行中
   - 已创建 `application_api_logs` 表
   - 已创建测试 Application: `test_app_001`
   - 已授权访问实体: `list_tree`

2. **Redis**
   - Redis 服务运行中
   - 用于 Nonce 防重放验证

3. **服务**
   - 后端服务运行在 `localhost:8787`

## 测试数据准备

```sql
-- 创建测试 Application
INSERT INTO applications (app_id, app_secret, name, ip, status, created_at, updated_at)
VALUES ('test_app_001', 'test_secret_123456', 'Test Application', '127.0.0.1,::1', 'Normal', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    app_secret = 'test_secret_123456',
    ip = '127.0.0.1,::1',
    status = 'Normal',
    updated_at = NOW();

-- 授权访问实体
INSERT INTO application_entities (app_id, entity_code, status, created_at, updated_at)
VALUES ('test_app_001', 'list_tree', 'Normal', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    status = 'Normal',
    updated_at = NOW();
```

## 签名算法

### Canonical Request 格式

```
HTTP_METHOD + "\n" +
URI_PATH + "\n" +
SORTED_QUERY_STRING + "\n" +
SHA256(REQUEST_BODY) + "\n" +
TIMESTAMP + "\n" +
NONCE
```

### 签名计算

```
HMAC-SHA256(Canonical_Request, App_Secret)
```

## 测试覆盖率

| 测试用例 | Shell 脚本 | Go 测试 | 状态 |
|---------|-----------|---------|------|
| List 请求 | ✅ | ⏳ | Shell 已实现 |
| 分页 List | ✅ | ⏳ | Shell 已实现 |
| Get 请求 | ✅ | ⏳ | Shell 已实现 |
| 错误签名 | ✅ | ⏳ | Shell 已实现 |
| 缺少请求头 | ✅ | ⏳ | Shell 已实现 |
| Nonce 重放 | ✅ | ⏳ | Shell 已实现 |

## 下一步

1. 完善 Go 集成测试的环境初始化
2. 实现测试数据的自动创建和清理
3. 添加更多边界条件测试
4. 集成到 CI/CD 流程

## 参考

- [External API Integration Design](../../docs/04-features/external-api/external-api-integration.md)
- [Phase 1 Implementation](../../docs/04-features/external-api/phase1-implementation.md)
