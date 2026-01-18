package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID          uint   `gorm:"primarykey"`
	Code        string `gorm:"size:64;unique;not null" json:"code" binding:"required"` // 权限标识,如 user:list
	Name        string `gorm:"size:64;not null" json:"name" binding:"required"`        // 权限名称
	Resource    string `gorm:"size:64" json:"resource"`                                // 资源类型,如 user, role
	Action      string `gorm:"size:64" json:"action"`                                  // 操作类型,如 read, write
	ParentID    uint   `gorm:"default:0" json:"parent_id"`                             // 父权限ID,用于构建树状结构
	Description string `gorm:"size:255" json:"description"`                            // 描述

	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:8;default:Normal"`
	CreatedBy string `gorm:"size:64" json:",omitempty"` // 创建人
	UpdatedBy string `gorm:"size:64" json:",omitempty"` // 更新人
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联关系
	Children []*Permission `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子权限
}

func (m *Permission) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *Permission) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	return
}
