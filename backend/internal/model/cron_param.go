package model

import (
	"time"

	"gorm.io/gorm"
)

// Cron Map 表
type CronParam struct {
	ID        uint   `gorm:"primaryKey"`
	CronCode  string `gorm:"size:8;" binding:"required,max=8"` // 任务编码
	Name      string `gorm:"size:64;" binding:"max=64"`        // 参数名称
	Type      string `gorm:"size:8;" binding:"max=8"`          // 参数类型 Text，Number，Boolean
	Value     string `gorm:"size:64;" binding:"max=64"`        // 参数值，text，now(), Hours(), Minute(), Day(), Month()
	Sort      uint   `gorm:"size:4;" binding:"max=999"`        // 动作排序
	Status    string `gorm:"size:16;default:Normal"`           // 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
