package model

import (
	"time"

	"gorm.io/gorm"
)

// Webhook Request Structure
type WebhookReq struct {
	HookID       uint
	ApprovalCode string
	TableCode    string
	EntityID     uint
}

// 系统表
type Webhook struct {
	ID        uint   `gorm:"primaryKey"`
	Url       string `gorm:"size:256;" binding:"required,max=256"` // 调用Url
	TableCode string `gorm:"size:64;" binding:"max=64"`            // 主数据Code
	Username  string `gorm:"size:64;" binding:"max=64"`            // 用户名
	// 请求方式：POST，GET
	// 数据类型：plain，application/json，application/xxx-form-urlencoded，application/xml
	// Data：string
	// Headers；[]Key，Value
	// Response示例
	// Basic Auth 用户名
	// Basic Auth 密码
	ContentType string `gorm:"size:128" binding:"max=128"`           // 系统名称
	Secret      string `gorm:"size:128" binding:"required,max=128" ` // AppSecret
	Events      string `gorm:"size:256" binding:"required,max=256" ` // 监控什么事件
	Description string `gorm:"size:255" binding:"max=255"`
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:8;default:Normal"`
	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Webhook) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *Webhook) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *Webhook) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
