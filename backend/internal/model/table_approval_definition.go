package model

import (
	"time"

	"gorm.io/gorm"
)

type TableApprovalDefinition struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	EntityCode      string `gorm:"uniqueIndex:idx_entity_code_operation;size:64;not null" json:"entity_code"`
	Operation       string `gorm:"uniqueIndex:idx_entity_code_operation;size:64;not null" json:"operation"`
	ApprovalDefCode string `gorm:"size:64;not null" json:"approval_def_code"`
	Description     string `gorm:"size:255" json:"description"`
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string         `gorm:"size:8;default:Normal" json:"status"`
	CreatedBy string         `gorm:"size:32" json:"created_by"`
	UpdatedBy string         `gorm:"size:32" json:"updated_by"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *TableApprovalDefinition) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *TableApprovalDefinition) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *TableApprovalDefinition) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
