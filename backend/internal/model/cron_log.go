package model

import (
	"time"
)

// Cron Log 表
type CronLog struct {
	ID        uint       `gorm:"primaryKey"`
	Method    string     `gorm:"size:32;" binding:"required,max=32"` // 任务名称
	Param     string     `gorm:"size:128;" binding:"max=128"`        // 任务参数
	ErrMsg    string     `gorm:"size:512;" binding:"max=512"`        // 错误信息
	StartTime *time.Time // 开始时间
	EndTime   *time.Time // 结束时间
	ExecTime  uint       // 执行时间(秒)
	Status    string     `gorm:"size:16;default:Normal"` // 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
