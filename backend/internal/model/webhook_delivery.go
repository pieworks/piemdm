package model

import (
	"time"
)

type WebhookDelivery struct {
	ID           uint   `gorm:"primaryKey"`
	HookID       uint   `gorm:"not null" binding:"required"`                      // WebHook编号
	DeliveryCode string `gorm:"size:64;not null;index" binding:"required,max=64"` // 调用编号
	// Event name: Release 发布(所有变更都会发布)
	Event           string     `gorm:"size:32;" binding:"max=32"`  // 事件名称
	EntityID        uint       `gorm:"size:64" binding:"required"` // 实体编码
	RequestHeaders  string     `gorm:"type:text"`                  // 请求头
	RequestPayload  string     `gorm:"type:mediumtext"`            // 请求参数
	ResponseStatus  int        `gorm:"size:10" binding:"max=999"`  // 回应状态
	ResponseMessage string     `gorm:"type:mediumtext"`            // 回应消息
	ResponseHeaders string     `gorm:"type:text"`                  // 回应头
	ResponseBody    string     `gorm:"type:mediumtext"`            // 回应信息
	DeliveredAt     *time.Time // 投递时间
	CompletedAt     *time.Time // 完成时间
}
