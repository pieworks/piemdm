package model

import (
	"time"

	"gorm.io/gorm"
)

// 审批定义
type TableFieldGroup struct {
	ID        uint   `gorm:"primarykey"`
	Code      string `gorm:"size:64;unique;not null" binding:"required,max=64"` // 字段组code
	Name      string `gorm:"size:128;not null" binding:"required,max=128"`      // 字段组名称
	TableCode string `gorm:"size:64" binding:"required,max=64"`                 // 实体编码
	View      string `gorm:"size:32" binding:"max=32"`                          // 视图编码
	Sort      uint   `gorm:"size:10;default:0"`                                 // 显示顺序
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status      string `gorm:"size:8;default:Normal"`
	Description string `gorm:"size:255" binding:"max=255"` // 备注
	CreatedBy   string `gorm:"size:64" json:",omitempty"`
	UpdatedBy   string `gorm:"size:64" json:",omitempty"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *TableFieldGroup) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *TableFieldGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *TableFieldGroup) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
