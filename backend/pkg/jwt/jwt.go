package jwt

import (
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var logger *slog.Logger

// JWT相关错误常量
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
	ErrEmptyKey     = errors.New("jwt key is empty")
)

type JWT struct {
	key []byte
}

type CustomClaims struct {
	ID       string
	UserName string
	Email    string `json:"email"`
	Admin    string `json:"admin"` // 是否为管理员
	// Roles        []string `json:"roles"`         // 用户角色列表
	// Permissions  []string `json:"permissions"`   // 用户权限列表
	// DepartmentID string   `json:"department_id"` // 部门ID（用于数据权限）
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper) *JWT {
	// return &JWT{key: []byte(conf.GetString("security.jwt.key"))}
	key := conf.GetString("security.jwt.key")
	return &JWT{key: []byte(key)}
}

func (j *JWT) GenToken(userId, userName, email, admin string) (string, error) {
	claims := CustomClaims{
		ID:       userId,
		UserName: userName,
		Email:    email,
		Admin:    admin,
		// Roles:       roles,
		// Permissions: permissions,
		// DepartmentID: departmentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "",
			Subject:   "",
			ID:        "",
			Audience:  []string{},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
// 	re := regexp.MustCompile(`(?i)Bearer `)
// 	tokenString = re.ReplaceAllString(tokenString, "")
// 	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return j.key, nil
// 	})

// 	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// 		return claims, nil
// 	} else {
// 		return nil, err
// 	}
// }

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	re := regexp.MustCompile(`(?i)Bearer `)
	tokenString = re.ReplaceAllString(tokenString, "")

	// 空令牌检查
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	// 检查密钥是否为空
	if j.key == nil || len(j.key) == 0 {
		return nil, ErrEmptyKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		// 根据错误类型返回自定义错误
		if err.Error() == "token is expired" || err.Error() == "token has invalid claims: token is expired" {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if token == nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
