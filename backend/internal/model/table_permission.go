package model

import (
	"time"

	"gorm.io/gorm"
)

type TablePermission struct {
	ID        uint   `gorm:"primaryKey"`
	RoleID    uint   `gorm:"index;not null" binding:"required"`         // 关联 Role 表
	TableCode string `gorm:"size:64;index;not null" binding:"required"` // 关联 Table Code
	Status    string `gorm:"size:8;default:Normal"`                     // Normal 正常, Frozen 已冻结
	Operation string `gorm:"size:32;default:All"`                       // 权限范围: All, Create, Update, Delete etc. 目前主要用All

	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate 钩子
func (m *TablePermission) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

// BeforeUpdate 钩子
func (m *TablePermission) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	return
}
