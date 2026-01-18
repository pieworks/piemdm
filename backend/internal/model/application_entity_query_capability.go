package model

import (
	"time"

	"gorm.io/gorm"
)

// 实体查询能力配置表
type ApplicationEntityQueryCapability struct {
	ID uint `gorm:"primaryKey"`
	// 实体代码
	EntityCode string `gorm:"size:100;not null;uniqueIndex:uk_entity_field" binding:"required,max=100"`
	// 字段代码
	FieldCode string `gorm:"size:100;not null;uniqueIndex:uk_entity_field" binding:"required,max=100"`
	// 允许的操作符，如 ["eq", "in", "range"]
	AllowedOps string `gorm:"type:json" binding:"max=255"`
	// 是否允许模糊搜索
	AllowLike bool `gorm:"default:false"`
	// IN 查询最大数量
	MaxInSize int `gorm:"default:100"`
	// 是否有索引（无索引字段严禁开放 Range/Like）
	IsIndexed bool `gorm:"default:false"`

	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:16;default:Normal"`
	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *ApplicationEntityQueryCapability) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *ApplicationEntityQueryCapability) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *ApplicationEntityQueryCapability) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	return
}
