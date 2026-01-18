package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationLog 通知日志模型
type NotificationLog struct {
	ID uint `gorm:"primarykey" json:"ID"`
	// 审批实例ID
	ApprovalID string `gorm:"size:64;index" json:"ApprovalID"`
	// 任务ID
	TaskID string `gorm:"size:64;index" json:"TaskID"`
	// 接收人ID
	RecipientID string `gorm:"size:64;index" binding:"required,max=64" json:"RecipientID"`
	// 接收人类型
	RecipientType string `gorm:"size:16" binding:"required,max=16" json:"RecipientType"`
	// 通知类型
	NotificationType string `gorm:"size:16" binding:"required,max=16" json:"NotificationType"`
	// 模板ID
	TemplateID string `gorm:"size:64" json:"TemplateID"`
	// 模板编码
	TemplateCode string `gorm:"size:100" json:"TemplateCode"`
	// 通知标题
	Title string `gorm:"size:200" binding:"required,max=200" json:"Title"`
	// 通知内容
	Content string `gorm:"type:text" binding:"required" json:"Content"`
	// 发送状态
	Status string `gorm:"size:16;default:pending;index" json:"Status"`
	// 发送时间
	SendTime *time.Time `gorm:"index" json:"SendTime"`
	// 错误信息
	ErrorMessage string `gorm:"type:text" json:"ErrorMessage"`
	// 重试次数
	RetryCount int `gorm:"default:0" json:"RetryCount"`
	// 最大重试次数
	MaxRetryCount int `gorm:"default:3" json:"MaxRetryCount"`
	// 下次重试时间
	NextRetryTime *time.Time `json:"NextRetryTime,omitempty"`
	// 额外数据JSON
	ExtraData string     `gorm:"type:json" json:"ExtraData,omitempty"`
	CreatedAt *time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt *time.Time `json:"UpdatedAt,omitempty"`
}

// 接收人类型常量
const (
	RecipientTypeUser       = "user"       // 用户
	RecipientTypeRole       = "role"       // 角色
	RecipientTypeDepartment = "department" // 部门
	RecipientTypeGroup      = "group"      // 用户组
)

// 通知状态常量
const (
	NotificationStatusPending = "pending" // 待发送
	NotificationStatusSent    = "sent"    // 已发送
	NotificationStatusFailed  = "failed"  // 发送失败
	NotificationStatusRetry   = "retry"   // 重试中
	NotificationStatusExpired = "expired" // 已过期
)

// BeforeCreate 创建前钩子
func (m *NotificationLog) BeforeCreate(tx *gorm.DB) (err error) {
	// 验证接收人类型
	if !IsValidRecipientType(m.RecipientType) {
		return gorm.ErrInvalidValue
	}

	// 验证通知类型
	if !IsValidNotificationType(m.NotificationType) {
		return gorm.ErrInvalidValue
	}

	// 验证通知状态
	if m.Status != "" && !IsValidNotificationStatus(m.Status) {
		return gorm.ErrInvalidValue
	}

	// 设置默认状态
	if m.Status == "" {
		m.Status = NotificationStatusPending
	}

	return nil
}

// BeforeUpdate 更新前钩子
func (m *NotificationLog) BeforeUpdate(tx *gorm.DB) (err error) {
	// 验证接收人类型
	if m.RecipientType != "" && !IsValidRecipientType(m.RecipientType) {
		return gorm.ErrInvalidValue
	}

	// 验证通知类型
	if m.NotificationType != "" && !IsValidNotificationType(m.NotificationType) {
		return gorm.ErrInvalidValue
	}

	// 验证通知状态
	if m.Status != "" && !IsValidNotificationStatus(m.Status) {
		return gorm.ErrInvalidValue
	}

	return nil
}

// IsPending 检查是否为待发送状态
func (m *NotificationLog) IsPending() bool {
	return m.Status == NotificationStatusPending
}

// IsSent 检查是否已发送
func (m *NotificationLog) IsSent() bool {
	return m.Status == NotificationStatusSent
}

// IsFailed 检查是否发送失败
func (m *NotificationLog) IsFailed() bool {
	return m.Status == NotificationStatusFailed
}

// IsRetrying 检查是否在重试中
func (m *NotificationLog) IsRetrying() bool {
	return m.Status == NotificationStatusRetry
}

// IsExpired 检查是否已过期
func (m *NotificationLog) IsExpired() bool {
	return m.Status == NotificationStatusExpired
}

// CanRetry 检查是否可以重试
func (m *NotificationLog) CanRetry() bool {
	return (m.IsFailed() || m.IsRetrying()) && m.RetryCount < m.MaxRetryCount
}

// MarkAsSent 标记为已发送
func (m *NotificationLog) MarkAsSent() {
	m.Status = NotificationStatusSent
	now := time.Now()
	m.SendTime = &now
	m.ErrorMessage = ""
}

// MarkAsFailed 标记为发送失败
func (m *NotificationLog) MarkAsFailed(errorMessage string) {
	m.Status = NotificationStatusFailed
	m.ErrorMessage = errorMessage
	m.RetryCount++

	// 设置下次重试时间（指数退避）
	if m.CanRetry() {
		retryDelay := time.Duration(1<<m.RetryCount) * time.Minute // 2^n 分钟
		nextRetry := time.Now().Add(retryDelay)
		m.NextRetryTime = &nextRetry
		m.Status = NotificationStatusRetry
	}
}

// MarkAsExpired 标记为已过期
func (m *NotificationLog) MarkAsExpired() {
	m.Status = NotificationStatusExpired
}

// IsValidRecipientType 验证接收人类型
func IsValidRecipientType(recipientType string) bool {
	validTypes := []string{
		RecipientTypeUser,
		RecipientTypeRole,
		RecipientTypeDepartment,
		RecipientTypeGroup,
	}

	for _, validType := range validTypes {
		if recipientType == validType {
			return true
		}
	}
	return false
}

// IsValidNotificationStatus 验证通知状态
func IsValidNotificationStatus(status string) bool {
	validStatuses := []string{
		NotificationStatusPending,
		NotificationStatusSent,
		NotificationStatusFailed,
		NotificationStatusRetry,
		NotificationStatusExpired,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}
