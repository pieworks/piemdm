package model

import (
	"time"
)

type GlobalId struct {
	ID          uint   `gorm:"primaryKey"`
	Identifier  string `gorm:"size:64;unique;not null" binding:"required"` // 业务名称
	LastID      uint   `gorm:"not null;default:1" binding:"required"`
	Step        uint   `gorm:"default:0"`
	Description string `gorm:"size:255"` // 审批实例描述
	CreatedBy   string `gorm:"size:64"`  // 创建人
	CreatedAt   *time.Time
}
