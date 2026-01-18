package model

import (
	"time"

	"gorm.io/gorm"
)

// 应用可以访问的实体的属性
type ApplicationEntityField struct {
	ID uint `gorm:"primaryKey"`
	// appID
	AppID string `gorm:"size:64;index:idx_app;" binding:"required,max=64"`
	// entityCode
	EntityCode string `gorm:"size:64;index:idx_entity_code;" binding:"required,max=64"`
	// fieldCode
	FieldCode string `gorm:"size:64" binding:"required,max=64"`

	Status    string `gorm:"size:8;default:Normal"` // 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
