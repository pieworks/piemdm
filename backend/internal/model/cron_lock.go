package model

import (
	"time"

	"gorm.io/gorm"
)

// Cron Lock 表
type CronLock struct {
	ID         uint   `gorm:"primaryKey"`
	LockMethod string `gorm:"size:32;" binding:"required,max=32"` // 任务名称
	ExpireTime uint   // 过期时间
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
