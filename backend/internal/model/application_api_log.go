package model

import (
	"time"
)

// Application API 请求日志表
type ApplicationApiLog struct {
	ID uint `gorm:"primaryKey"`
	// idx_request_id: 快速查询单个请求
	// 全链路 Trace ID
	RequestId string `gorm:"size:50;index:idx_request_id" binding:"required,max=50"`
	// idx_app_time: 应用 + 时间复合索引，支持按应用查询历史
	// Application ID
	ApplicationId string `gorm:"size:50;index:idx_app_time" binding:"required,max=50"`
	// idx_log_type: 按日志类型过滤
	// 日志类型: ACCESS / SECURITY / ERROR
	LogType string `gorm:"size:20;index:idx_log_type" binding:"required,max=20"`

	// -- 请求信息
	// HTTP 方法
	HttpMethod string `gorm:"size:10" binding:"required,max=10"`
	// 请求路径
	RequestPath string `gorm:"size:255" binding:"required,max=255"`
	// Query 参数（可选）
	QueryParams string `gorm:"type:text"`
	// 请求载荷（已脱敏，可选）
	RequestPayload string `gorm:"type:text"`
	// 客户端 IP
	ClientIp string `gorm:"size:50" binding:"max=50"`
	// User-Agent
	UserAgent string `gorm:"size:255" binding:"max=255"`

	// -- 响应信息
	// idx_http_status: 按 HTTP 状态码分析
	// HTTP 状态码 (协议层结果，如 200/400/401/403/429/500)
	HttpStatus int `gorm:"index:idx_http_status"`
	// 响应载荷（错误时记录）
	ResponsePayload string `gorm:"type:text"`
	// 请求耗时(ms)
	DurationMs int

	// -- 结果与错误
	// idx_outcome: 按结果类型统计
	// SUCCESS / FAILED / RATE_LIMITED / AUTH_FAILED
	Outcome string `gorm:"size:20;not null;index:idx_outcome" binding:"required,max=20"`
	// 错误码 (平台/业务层原因，机器可读，稳定，如 AUTH_SIGNATURE_INVALID)
	ErrorCode string `gorm:"size:50" binding:"max=50"`
	// 错误信息 (人可读，可变，如 "Invalid signature")
	ErrorMessage string `gorm:"size:255" binding:"max=255"`

	// -- 影响的数据
	// 影响的数据 ID 列表
	AffectedResourceIds *string `gorm:"type:json"`

	// 创建时间
	CreatedAt time.Time `gorm:"index:idx_app_time"`
}
