package service

import (
	"context"
	"fmt"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"
)

//go:generate mockgen -source=notification_log.go -destination=../../test/mocks/service/notification_log_service.go

// NotificationLogService 通知日志服务接口
type NotificationLogService interface {
	// 基础查询
	List(ctx context.Context, req *repository.ListNotificationLogRequest) (*NotificationLogListResult, error)
	Get(ctx context.Context, id uint) (*model.NotificationLog, error)
}

// notificationLogService 通知日志服务实现
type notificationLogService struct {
	logRepo repository.NotificationLogRepository
	logger  *log.Logger
}

// NewNotificationLogService 创建通知日志服务
func NewNotificationLogService(
	logRepo repository.NotificationLogRepository,
	logger *log.Logger,
) NotificationLogService {
	return &notificationLogService{
		logRepo: logRepo,
		logger:  logger,
	}
}

// List 获取通知日志列表
func (s *notificationLogService) List(ctx context.Context, req *repository.ListNotificationLogRequest) (*NotificationLogListResult, error) {
	logs, err := s.logRepo.FindPage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取通知日志列表失败: %v", err)
	}

	total, err := s.logRepo.Count(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("统计通知日志数量失败: %v", err)
	}

	return &NotificationLogListResult{
		Data:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// Get 获取通知日志详情
func (s *notificationLogService) Get(ctx context.Context, id uint) (*model.NotificationLog, error) {
	log, err := s.logRepo.FindOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取通知日志详情失败: %v", err)
	}

	if log == nil {
		return nil, fmt.Errorf("通知日志不存在: id=%d", id)
	}

	return log, nil
}
