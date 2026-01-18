package model

import (
	"time"

	"gorm.io/gorm"
)

type TableRelation struct {
	ID            uint   `gorm:"primaryKey"`
	TableCode     string `gorm:"size:64" binding:"max=64"`   // 模型
	RelationTable string `gorm:"size:64" binding:"max=64"`   // 关联模型
	RelationName  string `gorm:"size:128" binding:"max=128"` // 视图名称
	RelationCode  string `gorm:"size:64" binding:"max=64"`   // 关联编码
	Sort          uint   `gorm:"size:10;default:0"`          // 显示顺序
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:8;default:Normal"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *TableRelation) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Model(m).Where("id = ?", m.ID).Update("status", "Deleted")
	return
}
