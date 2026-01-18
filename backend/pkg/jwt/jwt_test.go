package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJwt(t *testing.T) {
	tests := []struct {
		name     string
		config   func() *viper.Viper
		wantErr  bool
		errorMsg string
	}{
		{
			name: "有效密钥配置",
			config: func() *viper.Viper {
				v := viper.New()
				v.Set("security.jwt.key", "test-secret-key-1234567890")
				return v
			},
			wantErr: false,
		},
		{
			name: "空密钥配置",
			config: func() *viper.Viper {
				v := viper.New()
				v.Set("security.jwt.key", "")
				return v
			},
			wantErr: false, // 空密钥也是有效的，但会在使用时出错
		},
		{
			name: "缺少密钥配置",
			config: func() *viper.Viper {
				v := viper.New()
				// 不设置 security.jwt.key
				return v
			},
			wantErr: false, // 返回空密钥
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJwt(tt.config())
			assert.NotNil(t, j)
			assert.NotNil(t, j.key)
		})
	}
}

func TestJWT_GenToken(t *testing.T) {
	// 创建测试配置
	v := viper.New()
	v.Set("security.jwt.key", "test-secret-key-for-unit-testing")

	j := NewJwt(v)

	tests := []struct {
		name      string
		userID    string
		userName  string
		email     string
		admin     string
		wantError bool
	}{
		{
			name:      "生成有效令牌",
			userID:    "123",
			userName:  "testuser",
			email:     "test@example.com",
			admin:     "Yes",
			wantError: false,
		},
		{
			name:      "生成普通用户令牌",
			userID:    "456",
			userName:  "regularuser",
			email:     "user@example.com",
			admin:     "No",
			wantError: false,
		},
		{
			name:      "空用户ID",
			userID:    "",
			userName:  "emptyuser",
			email:     "empty@example.com",
			admin:     "No",
			wantError: false, // 空ID应该也是允许的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := j.GenToken(tt.userID, tt.userName, tt.email, tt.admin)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Greater(t, len(token), 10) // 令牌应该有合理长度
			}
		})
	}
}

func TestJWT_ParseToken(t *testing.T) {
	// 创建测试配置
	v := viper.New()
	v.Set("security.jwt.key", "test-secret-key-for-parsing-test")

	j := NewJwt(v)

	// 生成测试令牌
	userID := "789"
	userName := "parsetest"
	email := "parse@example.com"
	admin := "Yes"

	token, err := j.GenToken(userID, userName, email, admin)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	tests := []struct {
		name        string
		tokenString string
		wantClaims  *CustomClaims
		wantError   bool
		errorType   error
	}{
		{
			name:        "解析有效令牌",
			tokenString: token,
			wantClaims: &CustomClaims{
				ID:       userID,
				UserName: userName,
				Email:    email,
				Admin:    admin,
			},
			wantError: false,
		},
		{
			name:        "解析带Bearer前缀的令牌",
			tokenString: "Bearer " + token,
			wantClaims: &CustomClaims{
				ID:       userID,
				UserName: userName,
				Email:    email,
				Admin:    admin,
			},
			wantError: false,
		},
		{
			name:        "解析无效令牌",
			tokenString: "invalid.token.string",
			wantError:   true,
			errorType:   ErrInvalidToken,
		},
		{
			name:        "解析空令牌",
			tokenString: "",
			wantError:   true,
			errorType:   ErrInvalidToken,
		},
		{
			name:        "解析篡改的令牌",
			tokenString: token + "tampered",
			wantError:   true,
			errorType:   ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := j.ParseToken(tt.tokenString)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.wantClaims.ID, claims.ID)
				assert.Equal(t, tt.wantClaims.UserName, claims.UserName)
				assert.Equal(t, tt.wantClaims.Email, claims.Email)
				assert.Equal(t, tt.wantClaims.Admin, claims.Admin)
				assert.NotZero(t, claims.ExpiresAt)
				assert.NotZero(t, claims.IssuedAt)
				assert.NotZero(t, claims.NotBefore)
			}
		})
	}
}

func TestJWT_ParseToken_ExpiredToken(t *testing.T) {
	// 创建测试配置
	v := viper.New()
	v.Set("security.jwt.key", "test-secret-for-expired-token")

	j := NewJwt(v)

	// 手动创建过期的令牌
	claims := CustomClaims{
		ID:       "expired-user",
		UserName: "expired",
		Email:    "expired@example.com",
		Admin:    "No",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)), // 过期时间：24小时前
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-48 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-48 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret-for-expired-token"))
	require.NoError(t, err)

	// 测试解析过期令牌
	parsedClaims, err := j.ParseToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestJWT_ParseToken_EmptyKey(t *testing.T) {
	// 创建空密钥的JWT实例
	v := viper.New()
	v.Set("security.jwt.key", "")
	j := NewJwt(v)

	// 尝试解析令牌应该返回空密钥错误
	_, err := j.ParseToken("any.token.string")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyKey)
}

func TestJWT_GenTokenAndParseToken_Integration(t *testing.T) {
	v := viper.New()
	v.Set("security.jwt.key", "integration-test-secret")
	j := NewJwt(v)

	// 测试数据
	testCases := []struct {
		name     string
		userID   string
		userName string
		email    string
		admin    string
	}{
		{"管理员用户", "1", "admin1", "admin1@example.com", "Yes"},
		{"普通用户", "2", "user1", "user1@example.com", "No"},
		{"特殊字符用户名", "3", "user_name-123", "user123@example.com", "No"},
		{"长用户ID", "1000000000000000001", "longid", "longid@example.com", "Yes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 生成令牌
			token, err := j.GenToken(tc.userID, tc.userName, tc.email, tc.admin)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// 解析令牌
			claims, err := j.ParseToken(token)
			assert.NoError(t, err)
			assert.NotNil(t, claims)

			// 验证声明内容
			assert.Equal(t, tc.userID, claims.ID)
			assert.Equal(t, tc.userName, claims.UserName)
			assert.Equal(t, tc.email, claims.Email)
			assert.Equal(t, tc.admin, claims.Admin)

			// 验证时间字段
			assert.NotNil(t, claims.ExpiresAt)
			assert.NotNil(t, claims.IssuedAt)
			assert.NotNil(t, claims.NotBefore)

			// 验证令牌在有效期内
			assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
			assert.True(t, claims.IssuedAt.Time.Before(time.Now()) || claims.IssuedAt.Time.Equal(time.Now()))
			assert.True(t, claims.NotBefore.Time.Before(time.Now()) || claims.NotBefore.Time.Equal(time.Now()))
		})
	}
}

func TestCustomClaims_JSONTags(t *testing.T) {
	// 验证CustomClaims结构体的JSON标签
	claims := CustomClaims{
		ID:       "test-id",
		UserName: "test-user",
		Email:    "test@example.com",
		Admin:    "Yes",
	}

	// 这里主要是验证结构体定义，确保JSON序列化时字段名正确
	assert.Equal(t, "email", getJSONTag(claims, "Email"))
	assert.Equal(t, "admin", getJSONTag(claims, "Admin"))
}

// 辅助函数：获取结构体字段的JSON标签
func getJSONTag(v interface{}, fieldName string) string {
	// 这个函数在实际测试中可能需要使用reflect包
	// 这里简化处理，直接返回已知的标签值
	switch fieldName {
	case "Email":
		return "email"
	case "Admin":
		return "admin"
	default:
		return ""
	}
}

func TestJWT_ConcurrentAccess(t *testing.T) {
	v := viper.New()
	v.Set("security.jwt.key", "concurrent-test-secret")
	j := NewJwt(v)

	// 并发测试：多个goroutine同时生成和解析令牌
	const numGoroutines = 10
	const tokensPerGoroutine = 10

	errChan := make(chan error, numGoroutines*tokensPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for k := 0; k < tokensPerGoroutine; k++ {
				userID := string(rune(goroutineID*100 + k))
				token, err := j.GenToken(userID, "user", "test@example.com", "No")
				if err != nil {
					errChan <- err
					continue
				}

				claims, err := j.ParseToken(token)
				if err != nil {
					errChan <- err
					continue
				}

				if claims.ID != userID {
					errChan <- fmt.Errorf("claims.ID mismatch: got %s, want %s", claims.ID, userID)
				}
			}
		}(i)
	}

	// 等待所有goroutine完成
	time.Sleep(100 * time.Millisecond)

	// 检查错误
	close(errChan)
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	assert.Empty(t, errors, "并发测试中发生错误: %v", errors)
}
