package model

import (
	"time"

	"gorm.io/gorm"
)

// Cron Map 表
type CronAction struct {
	ID          uint   `gorm:"primaryKey"`
	CronCode    string `gorm:"size:8;" binding:"required,max=8"` // 任务编码
	Root        string `gorm:"size:64;" binding:"max=64"`        // 数据根部
	InField     string `gorm:"size:64;" binding:"max=64"`        // 源系统字段名
	OutField    string `gorm:"size:64;" binding:"max=64"`        // 本系统字段名
	Action      string `gorm:"size:16;" binding:"max=16"`        // 动作
	ActionParam string `gorm:"size:16;" binding:"max=16"`        // 动作参数
	ActionSort  uint   `gorm:"size:4;" binding:"max=999"`        // 动作排序
	Status      string `gorm:"size:16;default:Normal"`           // 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	UpdatedAt   time.Time
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
