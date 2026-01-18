package model

import (
	"time"

	"gorm.io/gorm"
)

// func (TableField) TableName() string {
// 	return "table_fields"
// }

type TableField struct {
	// ===== field attr
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"size:64;uniqueIndex:idx_table_field_table_code_code;" binding:"required,max=64"` // 编码
	TableCode string `gorm:"size:64;uniqueIndex:idx_table_field_table_code_code;" binding:"required,max=64"` // 所属模型

	Name string `gorm:"size:128" binding:"required,max=128"` // 属性名称
	// FieldType 目前会保存字段格式，在后端进行字段默认配置，字段校验等
	// 后端的 FieldPreset 预设系统设计得非常完善，每个 FieldType 已经包含：
	//   DataType (数据库类型)
	//   Length (字段长度)
	//   UI.Widget (UI 组件)
	//   Validation (验证规则)
	// 结论：不需要移除 FieldType，它是一个很好的设计！
	FieldType     string `gorm:"size:32" binding:"max=32"`          // 业务字段类型(UI层): text, phone, email, select, autocode 等
	Type          string `gorm:"size:16" binding:"required,max=16"` // 数据库字段类型(DB层): Text, Integer, Date, Decimal 等
	Length        int    `binding:"required"`                       // 属性长度
	Required      string `gorm:"size:8" binding:"max=8"`            // 是否必填 Yes：必填 No:不必填
	IsIndex       string `gorm:"size:8" binding:"max=8"`            // 是否创建索引，Yes,No
	IsUnique      string `gorm:"size:8" binding:"max=8"`            // 唯一标识，Yes,No
	IndexName     string `gorm:"size:128" binding:"max=128"`        // 索引名称，扩展视图使用，不对外显示
	IndexPriority int    `binding:"max=3"`                          // 索引名称，扩展视图使用，不对外显示
	Description   string `gorm:"size:256" binding:"max=256"`        // 属性描述

	// ===== show style
	IsFilter  string `gorm:"size:8" binding:"max=8"`              // 是否作为过滤条件  Yes：是  No-否
	IsShow    string `gorm:"size:8" binding:"max=8"`              // 是否展示字段  Yes：是  No-否
	GroupName string `gorm:"size:64;comment:组名" binding:"max=64"` // 组名
	Sort      uint   `gorm:"size:10;default:0" form:"sort"`       // 显示顺序

	// ========== Options JSON 字段
	Options *FieldOptions `gorm:"type:json;comment:字段配置(JSON格式)"` // 新增: JSON 配置

	//  ========== 状态
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:8;default:Normal"`
	CreatedBy string `gorm:"size:64"` // 创建人
	UpdatedBy string `gorm:"size:64"` // 更新人
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *TableField) BeforeDelete(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}
	tx.Model(m).Where("id = ?", m.ID).Updates(map[string]any{
		"status":     "Deleted",
		"updated_by": m.UpdatedBy,
	})
	return
}

func (m *TableField) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *TableField) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}
