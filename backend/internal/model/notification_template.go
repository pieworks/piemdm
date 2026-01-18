package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// NotificationTemplate 通知模板模型
type NotificationTemplate struct {
	ID string `gorm:"primarykey;size:64"`
	// 模板编码
	TemplateCode string `gorm:"size:100;uniqueIndex" binding:"required,max=100"`
	// 模板名称
	TemplateName string `gorm:"size:200" binding:"required,max=200"`
	// 模板类型
	TemplateType string `gorm:"size:32" binding:"required,max=32"`
	// 通知类型
	NotificationType string `gorm:"size:16" binding:"required,max=16"`
	// 标题模板
	TitleTemplate string `gorm:"size:500" binding:"required,max=500"`
	// 内容模板
	ContentTemplate string `gorm:"type:text" binding:"required"`
	// 变量定义JSON
	Variables string `gorm:"type:json" json:",omitempty"`
	// 描述
	Description string `gorm:"size:500"`

	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string         `gorm:"size:8;default:Normal"`
	CreatedBy string         `gorm:"size:64" json:",omitempty"`
	UpdatedBy string         `gorm:"size:64" json:",omitempty"`
	CreatedAt *time.Time     `json:",omitempty"`
	UpdatedAt *time.Time     `json:",omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:",omitempty"`
}

// 通知模板类型常量
const (
	TemplateTypeApprovalStart    = "approval_start"    // 审批发起
	TemplateTypeApprovalPending  = "approval_pending"  // 待审批
	TemplateTypeApprovalApproved = "approval_approved" // 审批通过
	TemplateTypeApprovalRejected = "approval_rejected" // 审批拒绝
	TemplateTypeApprovalTimeout  = "approval_timeout"  // 审批超时
	TemplateTypeApprovalCancel   = "approval_cancel"   // 审批取消
	TemplateTypeTaskAssigned     = "task_assigned"     // 任务分配
	TemplateTypeTaskTransferred  = "task_transferred"  // 任务转交
	TemplateTypeTaskReminder     = "task_reminder"     // 任务催办
)

// 通知类型常量
const (
	NotificationTypeEmail    = "email"    // 邮件通知
	NotificationTypeSMS      = "sms"      // 短信通知
	NotificationTypeInternal = "internal" // 站内信
	NotificationTypeWebhook  = "webhook"  // Webhook通知
)

// BeforeCreate 创建前钩子
func (m *NotificationTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	// 设置创建人
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}

	// 验证模板类型
	if !IsValidTemplateType(m.TemplateType) {
		return gorm.ErrInvalidValue
	}

	// 验证通知类型
	if !IsValidNotificationType(m.NotificationType) {
		return gorm.ErrInvalidValue
	}

	return nil
}

// BeforeUpdate 更新前钩子
func (m *NotificationTemplate) BeforeUpdate(tx *gorm.DB) (err error) {
	// 设置更新人
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	// 验证模板类型
	if m.TemplateType != "" && !IsValidTemplateType(m.TemplateType) {
		return gorm.ErrInvalidValue
	}

	// 验证通知类型
	if m.NotificationType != "" && !IsValidNotificationType(m.NotificationType) {
		return gorm.ErrInvalidValue
	}

	return nil
}

// GetVariables 获取变量定义
func (m *NotificationTemplate) GetVariables() map[string]any {
	if m.Variables == "" {
		return make(map[string]any)
	}

	var variables map[string]any
	if err := json.Unmarshal([]byte(m.Variables), &variables); err != nil {
		return make(map[string]any)
	}

	return variables
}

// SetVariables 设置变量定义
func (m *NotificationTemplate) SetVariables(variables map[string]any) error {
	if variables == nil {
		m.Variables = ""
		return nil
	}

	data, err := json.Marshal(variables)
	if err != nil {
		return err
	}

	m.Variables = string(data)
	return nil
}

// IsValidTemplateType 验证模板类型
func IsValidTemplateType(templateType string) bool {
	validTypes := []string{
		TemplateTypeApprovalStart,
		TemplateTypeApprovalPending,
		TemplateTypeApprovalApproved,
		TemplateTypeApprovalRejected,
		TemplateTypeApprovalTimeout,
		TemplateTypeApprovalCancel,
		TemplateTypeTaskAssigned,
		TemplateTypeTaskTransferred,
		TemplateTypeTaskReminder,
	}

	for _, validType := range validTypes {
		if templateType == validType {
			return true
		}
	}
	return false
}

// IsValidNotificationType 验证通知类型
func IsValidNotificationType(notificationType string) bool {
	validTypes := []string{
		NotificationTypeEmail,
		NotificationTypeSMS,
		NotificationTypeInternal,
		NotificationTypeWebhook,
	}

	for _, validType := range validTypes {
		if notificationType == validType {
			return true
		}
	}
	return false
}
