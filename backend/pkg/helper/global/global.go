package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	// 系统配置
	// Conf Configuration
	// 日志
	// Log    *log.Logger
	// Logger *log.Logger
	// mysql实例
	// Mysql *gorm.DB
	// Casbin实例
	// CasbinACLEnforcer *casbin.SyncedEnforcer
	// validation.v10校验器
	Validate *validator.Validate
	// validation.v10相关翻译器
	Translator ut.Translator
)
