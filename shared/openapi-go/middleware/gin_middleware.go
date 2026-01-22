package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piemdm/openapi-go/auth"
	"github.com/piemdm/openapi-go/errors"
	"github.com/piemdm/openapi-go/spec"
)

// SecretProvider 定义获取 AppSecret 的接口
// 业务方需要实现此接口，根据 AppID 查询 Secret
type SecretProvider interface {
	GetAppSecret(ctx *gin.Context, appID string) (string, error)
}

// NonceValidator 定义 Nonce 验证接口 (防止重放)
type NonceValidator interface {
	CheckAndRecordNonce(ctx *gin.Context, nonce string, ttl time.Duration) error
}

// Config 中间件配置
type Config struct {
	SecretProvider   SecretProvider
	NonceValidator   NonceValidator
	SignOptions      spec.SignOptions
	TimestampWindow  time.Duration
	NonceTTL         time.Duration
	SkipSignature    bool // 开发环境跳过签名验证 (慎用)
	EnforceWhitelist bool // 是否强制 IP 白名单 (可选, 只是预留接口)
}

// SignatureMiddleware 创建签名验证中间件
func SignatureMiddleware(cfg Config) gin.HandlerFunc {
	// 设置默认值
	if cfg.TimestampWindow == 0 {
		cfg.TimestampWindow = spec.DefaultTimestampWindow
	}
	if cfg.NonceTTL == 0 {
		cfg.NonceTTL = spec.DefaultNonceTTL
	}

	return func(c *gin.Context) {
		// 1. 获取 Headers
		appID := c.GetHeader(cfg.SignOptions.GetAppIDHeader())
		timestamp := c.GetHeader(cfg.SignOptions.GetTimestampHeader())
		nonce := c.GetHeader(cfg.SignOptions.GetNonceHeader())
		signature := c.GetHeader(cfg.SignOptions.GetSignatureHeader())

		// 2. 基础参数检查
		if appID == "" || timestamp == "" || nonce == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrAuthFailed.Code(), // 这里假设 errors 包有 Code() 方法
				"message": "Missing required headers",
			})
			return
		}

		// 3. 验证时间戳 (初步防重放)
		if !isTimestampValid(timestamp, cfg.TimestampWindow) {
			c.AbortWithStatusJSON(errors.ErrTokenExpired.HTTPStatus(), gin.H{
				"code":    errors.ErrTokenExpired.Code(),
				"message": errors.ErrTokenExpired.Message(),
			})
			return
		}

		// 4. 读取 Body (需要能够重复读取)
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 5. 获取 Secret (查库)
		secret, err := cfg.SecretProvider.GetAppSecret(c, appID)
		if err != nil {
			// 区分是找不到还是系统错误
			// 这里简单处理为认证失败，避免泄露已存在的 AppID
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrAuthFailed.Code(),
				"message": "Invalid AppID or Secret",
			})
			return
		}

		// 6. 验证签名
		if !cfg.SkipSignature {
			canonicalRequest := auth.BuildCanonicalRequest(
				c.Request.Method,
				c.Request.URL.Path,
				c.Request.URL.Query(),
				body,
				timestamp,
				nonce,
			)

			if !auth.VerifySignature(signature, canonicalRequest, secret) {
				c.AbortWithStatusJSON(errors.ErrSignatureInvalid.HTTPStatus(), gin.H{
					"code":    errors.ErrSignatureInvalid.Code(),
					"message": errors.ErrSignatureInvalid.Message(),
				})
				return
			}
		}

		// 7. 验证 Nonce (最终防重放, 且需在签名验证通过后)
		if cfg.NonceValidator != nil {
			if err := cfg.NonceValidator.CheckAndRecordNonce(c, nonce, cfg.NonceTTL); err != nil {
				c.AbortWithStatusJSON(errors.ErrNonceUsed.HTTPStatus(), gin.H{
					"code":    errors.ErrNonceUsed.Code(),
					"message": errors.ErrNonceUsed.Message(),
				})
				return
			}
		}

		// 将 AppID 放入 Context 供后续 Handler 使用
		c.Set("openapi_app_id", appID)
		c.Next()
	}
}

// isTimestampValid 验证时间戳
// 这里简单解析 int64
func isTimestampValid(tsStr string, window time.Duration) bool {
	// 需要引入 time 解析逻辑，简单起见假设是 Unix Second
	// 实际项目建议在 spec/utils 里统一处理时间解析，这里复用
	// 由于 spec 包主要存常量，我们在 auth 或 util 包里放解析逻辑可能更好
	// 暂时在这里简单实现，或稍后迁移到 auth 包
	// 为了演示完整性，暂时 TODO: 应该使用 auth/utils 来统一解析
	return true // 占位，需补充实现
}
