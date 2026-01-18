package integration

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/internal/service"
	"piemdm/pkg/config"
	"piemdm/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// OpenAPI 核心逻辑集成测试
// 基于 test_openapi.sh 的测试用例,只测试核心的签名验证和权限逻辑

const (
	testAppID     = "test_app_001"
	testAppSecret = "test_secret_123456"
	testEntity    = "list_tree"
)

// OpenAPITestSuite OpenAPI 核心逻辑测试套件
type OpenAPITestSuite struct {
	suite.Suite
	db                       *gorm.DB
	redisClient              *redis.Client
	conf                     *viper.Viper
	logger                   *log.Logger
	openApiAuthService       service.OpenApiAuthService
	applicationApiLogService service.ApplicationApiLogService
}

// SetupSuite 测试套件初始化
func (s *OpenAPITestSuite) SetupSuite() {
	// 设置配置文件路径
	configPath := "../../config/local.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		s.T().Skip("配置文件不存在，跳过集成测试")
	}
	os.Setenv("APP_CONF", configPath)

	// 从配置文件加载配置
	s.conf = config.NewConfig()

	// 初始化日志
	s.logger = log.NewLog(s.conf)

	// 初始化数据库
	var err error
	s.db, err = gorm.Open(mysql.Open(s.conf.GetString("data.mysql.dsn")), &gorm.Config{})
	require.NoError(s.T(), err, "数据库连接失败")

	// 初始化 Redis
	s.redisClient = redis.NewClient(&redis.Options{
		Addr:     s.conf.GetString("data.redis.addr"),
		Password: s.conf.GetString("data.redis.password"),
		DB:       s.conf.GetInt("data.redis.db"),
	})

	// 初始化 Repository
	repo := repository.NewRepository(s.db, s.redisClient, s.logger)
	baseRepo := repository.NewBaseRepository(repo)
	applicationRepo := repository.NewApplicationRepository(repo, baseRepo)
	applicationEntityRepo := repository.NewApplicationEntityRepository(repo, baseRepo)
	applicationApiLogRepo := repository.NewApplicationApiLogRepository(repo, baseRepo)

	// 初始化 Service
	svc := service.NewService(s.logger, nil, nil)
	s.openApiAuthService = service.NewOpenApiAuthService(s.logger, applicationRepo, applicationEntityRepo, s.redisClient, s.conf)
	s.applicationApiLogService = service.NewApplicationApiLogService(svc, applicationApiLogRepo)

	// 准备测试数据
	s.prepareTestData()
}

// TearDownSuite 测试套件清理
func (s *OpenAPITestSuite) TearDownSuite() {
	// 清理测试数据
	s.cleanupTestData()
}

// prepareTestData 准备测试数据
func (s *OpenAPITestSuite) prepareTestData() {
	// 确保测试 Application 存在
	var app model.Application
	result := s.db.Where("app_id = ?", testAppID).First(&app)
	if result.Error != nil {
		// 创建测试 Application
		app = model.Application{
			AppId:     testAppID,
			AppSecret: testAppSecret,
			Name:      "Test Application",
			IP:        "127.0.0.1,::1",
			Status:    "Normal",
		}
		s.db.Create(&app)
	}

	// 确保 ApplicationEntity 存在
	var appEntity model.ApplicationEntity
	result = s.db.Where("app_id = ? AND entity_code = ?", testAppID, testEntity).First(&appEntity)
	if result.Error != nil {
		appEntity = model.ApplicationEntity{
			AppId:      testAppID,
			EntityCode: testEntity,
			Status:     "Normal",
		}
		s.db.Create(&appEntity)
	}
}

// cleanupTestData 清理测试数据
func (s *OpenAPITestSuite) cleanupTestData() {
	// 清理审计日志
	// s.db.Exec("DELETE FROM application_api_logs WHERE application_id = ?", testAppID)
}

// buildCanonicalRequest 构建 Canonical Request 字符串
func buildCanonicalRequest(method, path, query, bodyHash, timestamp, nonce string) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method,
		path,
		query,
		bodyHash,
		timestamp,
		nonce,
	)
}

// computeSignature 计算 HMAC-SHA256 签名
func computeSignature(canonicalRequest, appSecret string) string {
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(canonicalRequest))
	return hex.EncodeToString(h.Sum(nil))
}

// hashRequestBody 计算请求体的 SHA256 哈希
func hashRequestBody(body []byte) string {
	if len(body) == 0 {
		body = []byte("")
	}
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}

// TestSignatureVerification 测试签名验证
func (s *OpenAPITestSuite) TestSignatureVerification() {
	s.Run("正确的签名应该验证通过", func() {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		nonce := uuid.New().String()
		method := "GET"
		path := "/openapi/v1/entities/" + testEntity
		query := ""
		body := []byte("")

		bodyHash := hashRequestBody(body)
		canonicalRequest := buildCanonicalRequest(method, path, query, bodyHash, timestamp, nonce)
		signature := computeSignature(canonicalRequest, testAppSecret)

		// 创建 gin.Context 并设置 Request
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, nil)
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 验证签名
		app, err := s.openApiAuthService.VerifySignature(c, testAppID, timestamp, nonce, signature, body)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), app)
		assert.Equal(s.T(), testAppID, app.AppId)
	})

	s.Run("错误的签名应该验证失败", func() {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		nonce := uuid.New().String()
		invalidSignature := "invalid_signature_12345"

		// 创建 gin.Context 并设置 Request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/openapi/v1/entities/"+testEntity, nil)
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 验证签名
		app, err := s.openApiAuthService.VerifySignature(c, testAppID, timestamp, nonce, invalidSignature, []byte(""))

		assert.Error(s.T(), err)
		assert.Nil(s.T(), app)
		assert.Contains(s.T(), err.Error(), "signature mismatch")
	})

	s.Run("过期的时间戳应该验证失败", func() {
		// 使用 10 分钟前的时间戳 (超过 5 分钟窗口)
		timestamp := fmt.Sprintf("%d", time.Now().Add(-10*time.Minute).Unix())
		nonce := uuid.New().String()
		method := "GET"
		path := "/openapi/v1/entities/" + testEntity

		bodyHash := hashRequestBody([]byte(""))
		canonicalRequest := buildCanonicalRequest(method, path, "", bodyHash, timestamp, nonce)
		signature := computeSignature(canonicalRequest, testAppSecret)

		// 创建 gin.Context 并设置 Request
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, nil)
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 验证签名
		app, err := s.openApiAuthService.VerifySignature(c, testAppID, timestamp, nonce, signature, []byte(""))

		assert.Error(s.T(), err)
		assert.Nil(s.T(), app)
		assert.Contains(s.T(), err.Error(), "timestamp")
	})
}

// TestIPWhitelist 测试 IP 白名单验证
func (s *OpenAPITestSuite) TestIPWhitelist() {
	// 获取测试 Application
	var app model.Application
	s.db.Where("app_id = ?", testAppID).First(&app)

	s.Run("白名单中的 IP 应该通过", func() {
		err := s.openApiAuthService.VerifyIPWhitelist(&app, "127.0.0.1")
		assert.NoError(s.T(), err)
	})

	s.Run("不在白名单中的 IP 应该被拒绝", func() {
		err := s.openApiAuthService.VerifyIPWhitelist(&app, "192.168.1.100")
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "AUTH_IP_NOT_ALLOWED")
	})
}

// TestEntityAccess 测试 Entity 访问权限
func (s *OpenAPITestSuite) TestEntityAccess() {
	s.Run("有权限的 Entity 应该通过", func() {
		err := s.openApiAuthService.VerifyEntityAccess(testAppID, testEntity)
		assert.NoError(s.T(), err)
	})

	s.Run("无权限的 Entity 应该被拒绝", func() {
		err := s.openApiAuthService.VerifyEntityAccess(testAppID, "non_existent_entity")
		assert.Error(s.T(), err)
	})
}

// TestNonceReplay 测试 Nonce 防重放
func (s *OpenAPITestSuite) TestNonceReplay() {
	s.Run("Nonce 重放应该被拒绝", func() {
		nonce := uuid.New().String()
		ttl := 10 * time.Minute

		// 第一次使用 Nonce
		err := s.openApiAuthService.CheckAndRecordNonce(nonce, ttl)
		assert.NoError(s.T(), err)

		// 第二次使用相同的 Nonce (重放攻击)
		err = s.openApiAuthService.CheckAndRecordNonce(nonce, ttl)
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "already used")
	})
}

// TestAuditLog 测试审计日志
func (s *OpenAPITestSuite) TestAuditLog() {
	s.Run("应该能够创建审计日志", func() {
		auditLog := &model.ApplicationApiLog{
			RequestId:     uuid.New().String(),
			ApplicationId: testAppID,
			LogType:       "ACCESS",
			HttpMethod:    "GET",
			RequestPath:   "/openapi/v1/entities/" + testEntity,
			QueryParams:   "",
			ClientIp:      "127.0.0.1",
			HttpStatus:    200,
			DurationMs:    100,
			Outcome:       "SUCCESS",
		}

		// 创建 gin.Context
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		err := s.applicationApiLogService.Create(c, auditLog)
		assert.NoError(s.T(), err)

		// 验证日志已创建
		var count int64
		s.db.Model(&model.ApplicationApiLog{}).Where("request_id = ?", auditLog.RequestId).Count(&count)
		assert.Equal(s.T(), int64(1), count)
	})
}

// TestCanonicalRequestBuilding 测试 Canonical Request 构建
func (s *OpenAPITestSuite) TestCanonicalRequestBuilding() {
	s.Run("应该正确构建 Canonical Request", func() {
		method := "GET"
		path := "/openapi/v1/entities/product"
		query := "page=1&pageSize=10"
		body := []byte(`{"name":"test"}`)
		timestamp := "1234567890"
		nonce := "test-nonce-123"

		bodyHash := hashRequestBody(body)
		canonicalRequest := buildCanonicalRequest(method, path, query, bodyHash, timestamp, nonce)

		expected := "GET\n/openapi/v1/entities/product\npage=1&pageSize=10\n" + bodyHash + "\n1234567890\ntest-nonce-123"
		assert.Equal(s.T(), expected, canonicalRequest)
	})

	s.Run("空 body 应该正确哈希", func() {
		emptyHash := hashRequestBody([]byte(""))
		expectedHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // SHA256 of empty string
		assert.Equal(s.T(), expectedHash, emptyHash)
	})
}

// TestSignatureComputation 测试签名计算
func (s *OpenAPITestSuite) TestSignatureComputation() {
	s.Run("应该正确计算 HMAC-SHA256 签名", func() {
		canonicalRequest := "GET\n/api/v1/openapi/v1/entities/product\n\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n1234567890\ntest-nonce"
		appSecret := "test_secret"

		signature := computeSignature(canonicalRequest, appSecret)

		// 验证签名是 64 个字符的十六进制字符串 (SHA256 输出)
		assert.Len(s.T(), signature, 64)
		assert.Regexp(s.T(), "^[a-f0-9]{64}$", signature)
	})

	s.Run("相同输入应该产生相同签名", func() {
		canonicalRequest := "test-request"
		appSecret := "test-secret"

		sig1 := computeSignature(canonicalRequest, appSecret)
		sig2 := computeSignature(canonicalRequest, appSecret)

		assert.Equal(s.T(), sig1, sig2)
	})

	s.Run("不同输入应该产生不同签名", func() {
		appSecret := "test-secret"

		sig1 := computeSignature("request1", appSecret)
		sig2 := computeSignature("request2", appSecret)

		assert.NotEqual(s.T(), sig1, sig2)
	})
}

// TestOpenAPIIntegration 运行测试套件
func TestOpenAPIIntegration(t *testing.T) {
	suite.Run(t, new(OpenAPITestSuite))
}
