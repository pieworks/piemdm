package model

import (
	"time"

	"gorm.io/gorm"
)

// 系统应用表
type Application struct {
	ID          uint   `gorm:"primaryKey"`
	AppId       string `gorm:"size:64;unique;" binding:"max=64"` // AppId
	AppSecret   string `gorm:"size:128" binding:"max=128" `      // AppSecret
	Name        string `gorm:"size:128" binding:"max=128"`       // 系统名称
	IP          string ``                                        // 系统来源IP
	Description string `gorm:"size:255" binding:"max=255"`       // 简介

	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:16;default:Normal"`
	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Application) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *Application) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *Application) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
