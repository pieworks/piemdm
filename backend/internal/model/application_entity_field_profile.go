package model

import (
	"time"

	"gorm.io/gorm"
)

// 字段权限模板表
type ApplicationEntityFieldProfile struct {
	ID uint `gorm:"primaryKey"`
	// 模板代码，如 READ_BASIC, READ_FULL
	Code string `gorm:"size:32;unique;not null" binding:"required,max=32"`
	// 模板名称
	Name string `gorm:"size:64;not null" binding:"required,max=64"`
	// 模板描述
	Description string `gorm:"size:255" binding:"max=255"`
	// 适用的实体代码
	EntityCode string `gorm:"size:64;not null;index" binding:"required,max=64"`
	// 包含的字段代码列表 JSON 数组
	FieldCodes string `gorm:"type:json;not null"`

	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:16;default:Normal"`
	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *ApplicationEntityFieldProfile) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *ApplicationEntityFieldProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *ApplicationEntityFieldProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	return
}
