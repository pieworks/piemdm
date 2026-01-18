package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint   `gorm:"primarykey"`
	Code        string `gorm:"size:64;unique;not null" form:"code" binding:"required"` // 表名,英文
	Name        string `gorm:"size:128;not null" form:"name" binding:"required"`       // 名称，可以中文
	Description string `gorm:"size:255" form:"desc"`                                   // 描述
	DataScope   string `gorm:"size:32;default:Self" form:"data_scope"`                 // 数据范围: All, Subordinate, Self
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:8;default:Normal"`
	CreatedBy string `gorm:"size:64" json:",omitempty" ` // 创建人
	UpdatedBy string `gorm:"size:64" json:",omitempty" ` // 更新人
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Permissions []*Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (m *Role) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

// GORM 钩子：自动设置操作人
func (m *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
