package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type OpenApiAuthService interface {
	// 验证 Canonical Request 签名
	VerifySignature(c *gin.Context, appId, timestamp, nonce, sign string, body []byte) (*model.Application, error)

	// 验证 IP 白名单
	VerifyIPWhitelist(app *model.Application, clientIP string) error

	// 验证 Entity 访问权限
	VerifyEntityAccess(appId, entityCode string) error

	// 检查并记录 Nonce(防重放)
	CheckAndRecordNonce(nonce string, ttl time.Duration) error
}

type openApiAuthService struct {
	logger                *log.Logger
	applicationRepo       repository.ApplicationRepository
	applicationEntityRepo repository.ApplicationEntityRepository
	redisClient           *redis.Client
	conf                  *viper.Viper
}

func NewOpenApiAuthService(
	logger *log.Logger,
	applicationRepo repository.ApplicationRepository,
	applicationEntityRepo repository.ApplicationEntityRepository,
	redisClient *redis.Client,
	conf *viper.Viper,
) OpenApiAuthService {
	return &openApiAuthService{
		logger:                logger,
		applicationRepo:       applicationRepo,
		applicationEntityRepo: applicationEntityRepo,
		redisClient:           redisClient,
		conf:                  conf,
	}
}

// VerifySignature 验证 Canonical Request 签名
func (s *openApiAuthService) VerifySignature(
	c *gin.Context,
	appId, timestamp, nonce, sign string,
	body []byte,
) (*model.Application, error) {
	app, err := s.applicationRepo.FindByAppId(appId)
	if err != nil {
		s.logger.Warn("Application not found", "app_id", appId, "error", err)
		return nil, fmt.Errorf("AUTH_SIGNATURE_INVALID: application not found")
	}

	ts, err := parseTimestamp(timestamp)
	if err != nil {
		return nil, fmt.Errorf("AUTH_TOKEN_EXPIRED: invalid timestamp format")
	}

	if !isTimestampValid(ts, getTimestampWindow(s.conf)) {
		return nil, fmt.Errorf("AUTH_TOKEN_EXPIRED: timestamp expired")
	}

	canonicalRequest := buildCanonicalRequest(
		c.Request.Method,
		c.Request.URL.Path,
		c.Request.URL.Query(),
		body,
		timestamp,
		nonce,
	)

	s.logger.Debug("Canonical Request", "request", canonicalRequest)

	expectedSign := computeSignature(canonicalRequest, app.AppSecret)

	if !hmac.Equal([]byte(sign), []byte(expectedSign)) {
		s.logger.Warn("Signature mismatch",
			"app_id", appId,
			"expected", expectedSign,
			"got", sign)
		return nil, fmt.Errorf("AUTH_SIGNATURE_INVALID: signature mismatch")
	}

	return app, nil
}

// getTimestampWindow 从配置读取时间窗口,默认 5 分钟
func getTimestampWindow(conf *viper.Viper) time.Duration {
	windowMinutes := conf.GetInt("openapi.timestamp_window_minutes")
	if windowMinutes == 0 {
		windowMinutes = 5
	}
	return time.Duration(windowMinutes) * time.Minute
}

// isTimestampValid 验证时间戳是否在允许的窗口内
func isTimestampValid(ts time.Time, window time.Duration) bool {
	now := time.Now()
	return now.Sub(ts) <= window && ts.Sub(now) <= window
}

// VerifyIPWhitelist 验证 IP 白名单
func (s *openApiAuthService) VerifyIPWhitelist(app *model.Application, clientIP string) error {
	if app.IP == "" {
		// 如果没有配置 IP 白名单,则允许所有 IP
		return nil
	}

	// 支持逗号分隔的多个 IP
	allowedIPs := strings.Split(app.IP, ",")
	for _, ip := range allowedIPs {
		if strings.TrimSpace(ip) == clientIP {
			return nil
		}
	}

	s.logger.Warn("IP not in whitelist", "app_id", app.AppId, "client_ip", clientIP, "allowed_ips", app.IP)
	return fmt.Errorf("AUTH_IP_NOT_ALLOWED: IP not in whitelist")
}

// VerifyEntityAccess 验证 Entity 访问权限
func (s *openApiAuthService) VerifyEntityAccess(appId, entityCode string) error {
	// 查询 ApplicationEntity
	entity, err := s.applicationEntityRepo.FindByAppIdAndEntityCode(appId, entityCode)
	if err != nil {
		s.logger.Warn("Entity access denied", "app_id", appId, "entity_code", entityCode, "error", err)
		return fmt.Errorf("PERMISSION_DENIED: no access to this entity")
	}

	// 检查状态
	if entity.Status != "Normal" {
		return fmt.Errorf("PERMISSION_DENIED: entity access is not active")
	}

	// 检查过期时间
	if !entity.ExpiredAt.IsZero() && time.Now().After(entity.ExpiredAt) {
		return fmt.Errorf("PERMISSION_DENIED: entity access has expired")
	}

	return nil
}

// CheckAndRecordNonce 检查并记录 Nonce(防重放)
func (s *openApiAuthService) CheckAndRecordNonce(nonce string, ttl time.Duration) error {
	key := fmt.Sprintf("openapi:nonce:%s", nonce)

	// 尝试设置 Nonce,如果已存在则返回错误
	ctx := context.Background()
	success, err := s.redisClient.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		s.logger.Error("Failed to check nonce", "error", err)
		return fmt.Errorf("SYSTEM_INTERNAL_ERROR: failed to check nonce")
	}

	if !success {
		s.logger.Warn("Nonce already used", "nonce", nonce)
		return fmt.Errorf("AUTH_SIGNATURE_INVALID: nonce already used (replay attack)")
	}

	return nil
}

// buildCanonicalRequest 构建规范请求字符串
func buildCanonicalRequest(
	method, path string,
	queryParams url.Values,
	body []byte,
	timestamp, nonce string,
) string {
	// 1. HTTP 方法
	canonicalMethod := strings.ToUpper(method)

	// 2. URI 路径
	canonicalPath := path

	// 3. 排序后的查询参数
	canonicalQuery := sortQueryString(queryParams)

	// 4. 请求体哈希
	bodyHash := hashRequestBody(body)

	// 5. 组合 Canonical Request
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		canonicalMethod,
		canonicalPath,
		canonicalQuery,
		bodyHash,
		timestamp,
		nonce,
	)
}

// sortQueryString 对查询参数排序
func sortQueryString(queryParams url.Values) string {
	if len(queryParams) == 0 {
		return ""
	}

	// 获取所有键并排序
	keys := make([]string, 0, len(queryParams))
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建排序后的查询字符串
	var parts []string
	for _, k := range keys {
		for _, v := range queryParams[k] {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
		}
	}

	return strings.Join(parts, "&")
}

// hashRequestBody 计算请求体的 SHA256 哈希
func hashRequestBody(body []byte) string {
	if len(body) == 0 {
		body = []byte("")
	}

	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}

// computeSignature 计算 HMAC-SHA256 签名
func computeSignature(canonicalRequest, appSecret string) string {
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(canonicalRequest))
	return hex.EncodeToString(h.Sum(nil))
}

// parseTimestamp 解析时间戳
func parseTimestamp(timestamp string) (time.Time, error) {
	// 支持 Unix 时间戳(秒)
	ts, err := time.Parse("1136239445", timestamp)
	if err == nil {
		return ts, nil
	}

	// 尝试解析为整数
	var unixTime int64
	_, err = fmt.Sscanf(timestamp, "%d", &unixTime)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unixTime, 0), nil
}
