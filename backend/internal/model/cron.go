package model

import (
	"time"

	"gorm.io/gorm"
)

// Cron 表
type Cron struct {
	ID         uint   `gorm:"primaryKey"`
	Code       string `gorm:"size:8;unique;" binding:"required,max=8"` // Cron编码
	Expression string `gorm:"size:32;" binding:"max=32"`               // Cron表达式
	Name       string `gorm:"size:128;" binding:"max=128"`             // 任务名称
	EntityCode string `gorm:"size:64" binding:"max=64"`                // 实体编码
	System     string `gorm:"size:64" binding:"max=64" `               // 系统名称
	Url        string `gorm:"size:255" binding:"max=255"`
	// 协议类型，Http, Rest, GraphQL, GRPC, Soap, Jwt
	Protocol    string `gorm:"size:32" binding:"max=32" `
	Method      string `gorm:"size:16" binding:"max=16" `
	AppId       string `gorm:"size:64" binding:"max=64" `
	AppKey      string `gorm:"size:64" binding:"max=64" `
	SignType    string `gorm:"size:64" binding:"max=64" `
	Description string `gorm:"size:255" binding:"max=255"`
	// 状态:Normal 正常, Frozen 已冻结, Deleted 已删除
	Status    string `gorm:"size:8;default:Normal"`
	CreatedBy string `gorm:"size:32"`
	UpdatedBy string `gorm:"size:32"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Cron) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *Cron) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *Cron) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
