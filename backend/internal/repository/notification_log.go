package repository

import (
	"context"
	"time"

	"piemdm/internal/model"

	"gorm.io/gorm"
)

//go:generate mockgen -source=notification_log_repository.go -destination=../../test/mocks/repository/notification_log.go

// NotificationLogRepository 通知日志仓储接口
type NotificationLogRepository interface {
	// 基础查询
	FindOne(ctx context.Context, id uint) (*model.NotificationLog, error)
	FindPage(ctx context.Context, req *ListNotificationLogRequest) ([]*model.NotificationLog, error)
	Count(ctx context.Context, req *ListNotificationLogRequest) (int64, error)
	GetStatistics(ctx context.Context, req *NotificationStatisticsRequest) (*NotificationStatistics, error)

	// Base CRUD
	Create(ctx context.Context, log *model.NotificationLog) error
	Update(ctx context.Context, log *model.NotificationLog) error

	// Batch operations
	BatchUpdateStatus(ctx context.Context, ids []uint, status string) error

	// 业务查询方法
	GetPendingLogs(ctx context.Context, limit int) ([]*model.NotificationLog, error)
	GetRetryLogs(ctx context.Context, limit int) ([]*model.NotificationLog, error)
	GetByApprovalID(ctx context.Context, approvalID string) ([]*model.NotificationLog, error)
	GetByTaskID(ctx context.Context, taskID string) ([]*model.NotificationLog, error)
	GetByRecipient(ctx context.Context, recipientID, recipientType string, limit int) ([]*model.NotificationLog, error)

	// 辅助方法
	DeleteExpiredLogs(ctx context.Context, expiredBefore time.Time) error
}

// ListNotificationLogRequest 通知日志列表请求
type ListNotificationLogRequest struct {
	ApprovalID       string     `json:"approval_id"`       // 审批ID
	TaskID           string     `json:"task_id"`           // 任务ID
	RecipientID      string     `json:"recipient_id"`      // 接收人ID
	RecipientType    string     `json:"recipient_type"`    // 接收人类型
	NotificationType string     `json:"notification_type"` // 通知类型
	Status           string     `json:"status"`            // 状态
	StartTime        *time.Time `json:"start_time"`        // 开始时间
	EndTime          *time.Time `json:"end_time"`          // 结束时间
	Page             int        `json:"page"`              // 页码
	PageSize         int        `json:"page_size"`         // 页大小
}

// NotificationStatisticsRequest 通知统计请求
type NotificationStatisticsRequest struct {
	StartTime        *time.Time `json:"start_time"`        // 开始时间
	EndTime          *time.Time `json:"end_time"`          // 结束时间
	NotificationType string     `json:"notification_type"` // 通知类型
	RecipientType    string     `json:"recipient_type"`    // 接收人类型
}

// NotificationStatistics 通知统计信息
type NotificationStatistics struct {
	TotalCount   int64            `json:"total_count"`   // 总数
	SentCount    int64            `json:"sent_count"`    // 已发送数
	FailedCount  int64            `json:"failed_count"`  // 失败数
	PendingCount int64            `json:"pending_count"` // 待发送数
	RetryCount   int64            `json:"retry_count"`   // 重试数
	ByType       map[string]int64 `json:"by_type"`       // 按类型统计
	ByStatus     map[string]int64 `json:"by_status"`     // 按状态统计
	ByDate       map[string]int64 `json:"by_date"`       // 按日期统计
}

// notificationLogRepository 通知日志仓储实现
type notificationLogRepository struct {
	db *gorm.DB
}

// NewNotificationLogRepository 创建通知日志仓储
func NewNotificationLogRepository(db *gorm.DB) NotificationLogRepository {
	return &notificationLogRepository{
		db: db,
	}
}

// Create 创建通知日志
func (r *notificationLogRepository) Create(ctx context.Context, log *model.NotificationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// Update 更新通知日志
func (r *notificationLogRepository) Update(ctx context.Context, log *model.NotificationLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

// FindOne 根据ID获取通知日志
func (r *notificationLogRepository) FindOne(ctx context.Context, id uint) (*model.NotificationLog, error) {
	var log model.NotificationLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// FindPage 获取通知日志列表
func (r *notificationLogRepository) FindPage(ctx context.Context, req *ListNotificationLogRequest) ([]*model.NotificationLog, error) {
	var logs []*model.NotificationLog

	query := r.db.WithContext(ctx).Model(&model.NotificationLog{})

	// 添加过滤条件
	query = r.buildListQuery(query, req)

	// 分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// 排序
	query = query.Order("created_at DESC")

	err := query.Find(&logs).Error
	return logs, err
}

// Count 统计通知日志数量
func (r *notificationLogRepository) Count(ctx context.Context, req *ListNotificationLogRequest) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.NotificationLog{})

	// 添加过滤条件
	query = r.buildListQuery(query, req)

	err := query.Count(&count).Error
	return count, err
}

// GetPendingLogs 获取待发送的通知日志
func (r *notificationLogRepository) GetPendingLogs(ctx context.Context, limit int) ([]*model.NotificationLog, error) {
	var logs []*model.NotificationLog

	query := r.db.WithContext(ctx).Where("status = ?", model.NotificationStatusPending).
		Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// GetRetryLogs 获取需要重试的通知日志
func (r *notificationLogRepository) GetRetryLogs(ctx context.Context, limit int) ([]*model.NotificationLog, error) {
	var logs []*model.NotificationLog

	now := time.Now()
	query := r.db.WithContext(ctx).Where("status = ? AND next_retry_time <= ? AND retry_count < max_retry_count",
		model.NotificationStatusRetry, now).
		Order("next_retry_time ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// GetByApprovalID 根据审批ID获取通知日志
func (r *notificationLogRepository) GetByApprovalID(ctx context.Context, approvalID string) ([]*model.NotificationLog, error) {
	var logs []*model.NotificationLog
	err := r.db.WithContext(ctx).Where("approval_id = ?", approvalID).
		Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetByTaskID 根据任务ID获取通知日志
func (r *notificationLogRepository) GetByTaskID(ctx context.Context, taskID string) ([]*model.NotificationLog, error) {
	var logs []*model.NotificationLog
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).
		Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetByRecipient 根据接收人获取通知日志
func (r *notificationLogRepository) GetByRecipient(ctx context.Context, recipientID, recipientType string, limit int) ([]*model.NotificationLog, error) {
	var logs []*model.NotificationLog

	query := r.db.WithContext(ctx).Where("recipient_id = ? AND recipient_type = ?", recipientID, recipientType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// GetStatistics 获取通知统计信息
func (r *notificationLogRepository) GetStatistics(ctx context.Context, req *NotificationStatisticsRequest) (*NotificationStatistics, error) {
	stats := &NotificationStatistics{
		ByType:   make(map[string]int64),
		ByStatus: make(map[string]int64),
		ByDate:   make(map[string]int64),
	}

	query := r.db.WithContext(ctx).Model(&model.NotificationLog{})

	// 添加时间过滤
	if req.StartTime != nil {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("created_at <= ?", req.EndTime)
	}
	if req.NotificationType != "" {
		query = query.Where("notification_type = ?", req.NotificationType)
	}
	if req.RecipientType != "" {
		query = query.Where("recipient_type = ?", req.RecipientType)
	}

	// 总数统计
	query.Count(&stats.TotalCount)

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	r.db.WithContext(ctx).Model(&model.NotificationLog{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	for _, stat := range statusStats {
		stats.ByStatus[stat.Status] = stat.Count
		switch stat.Status {
		case model.NotificationStatusSent:
			stats.SentCount = stat.Count
		case model.NotificationStatusFailed:
			stats.FailedCount = stat.Count
		case model.NotificationStatusPending:
			stats.PendingCount = stat.Count
		case model.NotificationStatusRetry:
			stats.RetryCount = stat.Count
		}
	}

	// 按类型统计
	var typeStats []struct {
		NotificationType string `json:"notification_type"`
		Count            int64  `json:"count"`
	}
	r.db.WithContext(ctx).Model(&model.NotificationLog{}).
		Select("notification_type, COUNT(*) as count").
		Group("notification_type").
		Scan(&typeStats)

	for _, stat := range typeStats {
		stats.ByType[stat.NotificationType] = stat.Count
	}

	return stats, nil
}

// BatchUpdateStatus 批量更新状态
func (r *notificationLogRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status string) error {
	return r.db.WithContext(ctx).Model(&model.NotificationLog{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// DeleteExpiredLogs 删除过期日志
func (r *notificationLogRepository) DeleteExpiredLogs(ctx context.Context, expiredBefore time.Time) error {
	return r.db.WithContext(ctx).Where("created_at < ? AND status IN ?",
		expiredBefore, []string{model.NotificationStatusSent, model.NotificationStatusExpired}).
		Delete(&model.NotificationLog{}).Error
}

// buildListQuery 构建列表查询条件
func (r *notificationLogRepository) buildListQuery(query *gorm.DB, req *ListNotificationLogRequest) *gorm.DB {
	if req.ApprovalID != "" {
		query = query.Where("approval_id = ?", req.ApprovalID)
	}

	if req.TaskID != "" {
		query = query.Where("task_id = ?", req.TaskID)
	}

	if req.RecipientID != "" {
		query = query.Where("recipient_id = ?", req.RecipientID)
	}

	if req.RecipientType != "" {
		query = query.Where("recipient_type = ?", req.RecipientType)
	}

	if req.NotificationType != "" {
		query = query.Where("notification_type = ?", req.NotificationType)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.StartTime != nil {
		query = query.Where("created_at >= ?", req.StartTime)
	}

	if req.EndTime != nil {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	return query
}
