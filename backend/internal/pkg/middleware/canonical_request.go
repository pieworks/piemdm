package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"
	"piemdm/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CanonicalRequestMiddleware Canonical Request 签名验证中间件
func CanonicalRequestMiddleware(
	authService service.OpenApiAuthService,
	logger *log.Logger,
	conf *viper.Viper,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		appId := c.GetHeader("X-App-Id")
		timestamp := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		sign := c.GetHeader("X-Sign")

		if appId == "" || timestamp == "" || nonce == "" || sign == "" {
			logger.Warn("Missing required headers")
			resp.HandleError(c, http.StatusUnauthorized, "AUTH_FAILED: missing required headers (X-App-Id, X-Timestamp, X-Nonce, X-Sign)", nil)
			c.Abort()
			return
		}

		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		app, err := authService.VerifySignature(c, appId, timestamp, nonce, sign, body)
		if err != nil {
			logger.Warn("Signature verification failed", "error", err, "app_id", appId)
			resp.HandleError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		clientIP := c.ClientIP()
		if err := authService.VerifyIPWhitelist(app, clientIP); err != nil {
			logger.Warn("IP whitelist check failed", "error", err, "app_id", appId, "client_ip", clientIP)
			resp.HandleError(c, http.StatusForbidden, err.Error(), nil)
			c.Abort()
			return
		}

		nonceTTL := getNonceTTL(conf)
		if err := authService.CheckAndRecordNonce(nonce, nonceTTL); err != nil {
			logger.Warn("Nonce check failed", "error", err, "nonce", nonce)
			resp.HandleError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		c.Set("application", app)
		logger.Debug("Canonical request verified", "app_id", appId)
		c.Next()
	}
}

// getNonceTTL 从配置读取 Nonce TTL,默认 10 分钟
func getNonceTTL(conf *viper.Viper) time.Duration {
	ttlMinutes := conf.GetInt("openapi.nonce_ttl_minutes")
	if ttlMinutes == 0 {
		ttlMinutes = 10
	}
	return time.Duration(ttlMinutes) * time.Minute
}

// OpenApiEntityPermissionMiddleware Entity 权限校验中间件
func OpenApiEntityPermissionMiddleware(
	authService service.OpenApiAuthService,
	logger *log.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		app, exists := c.Get("application")
		if !exists {
			logger.Error("Application not found in context")
			resp.HandleError(c, http.StatusUnauthorized, "AUTH_FAILED: application not found", nil)
			c.Abort()
			return
		}

		application, ok := app.(*model.Application)
		if !ok {
			logger.Error("Invalid application type in context")
			resp.HandleError(c, http.StatusInternalServerError, "SYSTEM_INTERNAL_ERROR", nil)
			c.Abort()
			return
		}

		tableCode := c.Param("table")
		if tableCode == "" {
			resp.HandleError(c, http.StatusBadRequest, "PARAM_REQUIRED_MISSING: table is required", nil)
			c.Abort()
			return
		}

		if err := authService.VerifyEntityAccess(application.AppId, tableCode); err != nil {
			logger.Warn("Entity access denied", "app_id", application.AppId, "table", tableCode, "error", err)
			resp.HandleError(c, http.StatusForbidden, err.Error(), nil)
			c.Abort()
			return
		}

		logger.Debug("Entity access granted", "app_id", application.AppId, "table", tableCode)
		c.Next()
	}
}
