package model

import (
	"time"

	"gorm.io/gorm"
)

// 应用可以访问的实体
type ApplicationEntity struct {
	ID uint `gorm:"primaryKey"`
	// AppId
	AppId string `gorm:"size:64;unique;" binding:"required,max=64"`
	// EntityId
	EntityCode string `gorm:"size:64" binding:"required,max=64"`
	// Expire time
	ExpiredAt time.Time

	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status    string `gorm:"size:16;default:Normal"`
	CreatedBy string `gorm:"size:64"`
	UpdatedBy string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
