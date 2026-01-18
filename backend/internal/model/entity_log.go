package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityLog struct {
	ID           uint   `gorm:"primaryKey"`
	EntityID     uint   `gorm:"not null;index" binding:"required"`          // 关联的动态表记录ID
	FieldCode    string `gorm:"size:63;not null" binding:"required,max=64"` // 修改的字段
	FieldName    string `gorm:"size:127;" binding:"max=128"`                // 修改的字段
	BeforeUpdate string `gorm:"size:255"`                                   // 修改前
	AfterUpdate  string `gorm:"size:255"`                                   // 修改后
	Reason       string `gorm:"size:255"`                                   // 原因
	UpdateBy     string `gorm:"size:64"`                                    // 修改人
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
