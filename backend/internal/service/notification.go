package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"
)

//go:generate mockgen -source=notification.go -destination=../../test/mocks/service/notification_service.go

// NotificationService 通知服务接口
type NotificationService interface {
	// 发送通知
	SendApproval(ctx context.Context, req *ApprovalNotificationRequest) error
	SendBatch(ctx context.Context, reqs []*ApprovalNotificationRequest) error
	SendByTemplate(ctx context.Context, req *TemplateNotificationRequest) error
	Test(ctx context.Context, req *TestNotificationRequest) error

	// 统计
	GetStatistics(ctx context.Context, req *repository.NotificationStatisticsRequest) (*repository.NotificationStatistics, error)

	// 任务处理
	ProcessPending(ctx context.Context, limit int) error
	ProcessRetry(ctx context.Context, limit int) error

	// 渠道管理
	RegisterChannel(channel NotificationChannel) error
	GetEnabledChannels() []string
}

// ApprovalNotificationRequest 审批通知请求
type ApprovalNotificationRequest struct {
	ApprovalID       string         `json:"approval_id"`       // 审批实例ID
	TaskID           string         `json:"task_id"`           // 任务ID（可选）
	TemplateType     string         `json:"template_type"`     // 模板类型
	NotificationType string         `json:"notification_type"` // 通知类型
	RecipientID      string         `json:"recipient_id"`      // 接收人ID
	RecipientType    string         `json:"recipient_type"`    // 接收人类型
	Variables        map[string]any `json:"variables"`         // 模板变量
	Priority         int            `json:"priority"`          // 优先级
	ExpireTime       *time.Time     `json:"expire_time"`       // 过期时间
}

// TemplateNotificationRequest 模板通知请求
type TemplateNotificationRequest struct {
	TemplateCode     string         `json:"template_code"`     // 模板编码
	NotificationType string         `json:"notification_type"` // 通知类型
	RecipientID      string         `json:"recipient_id"`      // 接收人ID
	RecipientType    string         `json:"recipient_type"`    // 接收人类型
	Variables        map[string]any `json:"variables"`         // 模板变量
	ApprovalID       string         `json:"approval_id"`       // 审批实例ID（可选）
	TaskID           string         `json:"task_id"`           // 任务ID（可选）
	Priority         int            `json:"priority"`          // 优先级
	ExpireTime       *time.Time     `json:"expire_time"`       // 过期时间
}

// TestNotificationRequest 测试通知请求
type TestNotificationRequest struct {
	NotificationType string         `json:"notification_type"` // 通知类型
	RecipientID      string         `json:"recipient_id"`      // 接收人ID
	Title            string         `json:"title"`             // 标题
	Content          string         `json:"content"`           // 内容
	Variables        map[string]any `json:"variables"`         // 变量
}

// NotificationLogListResult 通知日志列表结果 (非泛型,用于mockgen兼容)
type NotificationLogListResult struct {
	Data     []*model.NotificationLog `json:"data"`      // 数据列表
	Total    int64                    `json:"total"`     // 总数
	Page     int                      `json:"page"`      // 当前页码
	PageSize int                      `json:"page_size"` // 页面大小
}

// notificationService 通知服务实现
type notificationService struct {
	templateService NotificationTemplateService
	logRepo         repository.NotificationLogRepository
	channelManager  *NotificationChannelManager
	logger          *log.Logger
	mu              sync.RWMutex
}

// NewNotificationService 创建通知服务
func NewNotificationService(
	templateService NotificationTemplateService,
	logRepo repository.NotificationLogRepository,
	logger *log.Logger,
) NotificationService {
	return &notificationService{
		templateService: templateService,
		logRepo:         logRepo,
		// 创建通知渠道管理器
		channelManager: NewNotificationChannelManager(logger),
		logger:         logger,
	}
}

// SendApproval 发送审批通知
func (s *notificationService) SendApproval(ctx context.Context, req *ApprovalNotificationRequest) error {

	// 1. 获取通知模板
	template, err := s.templateService.GetByTypeAndNotification(ctx, req.TemplateType, req.NotificationType)
	if err != nil {
		return fmt.Errorf("获取通知模板失败: %v", err)
	}

	// 2. 渲染模板
	rendered, err := s.templateService.RenderTemplate(ctx, template.ID, req.Variables)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}

	// 3. 创建通知日志
	notificationLog := &model.NotificationLog{
		ApprovalID:       req.ApprovalID,
		TaskID:           req.TaskID,
		RecipientID:      req.RecipientID,
		RecipientType:    req.RecipientType,
		NotificationType: req.NotificationType,
		TemplateID:       template.ID,
		TemplateCode:     template.TemplateCode,
		Title:            rendered.Title,
		Content:          rendered.Content,
		Status:           model.NotificationStatusPending,
		MaxRetryCount:    3,
	}

	// 设置额外数据
	if req.Variables != nil {
		extraData, _ := json.Marshal(req.Variables)
		notificationLog.ExtraData = string(extraData)
	}

	// 保存通知日志
	if err := s.logRepo.Create(ctx, notificationLog); err != nil {
		return fmt.Errorf("创建通知日志失败: %v", err)
	}

	// 4. 异步发送通知
	go s.sendNotificationAsync(ctx, notificationLog, rendered)

	return nil
}

// SendBatch 批量发送通知
func (s *notificationService) SendBatch(ctx context.Context, reqs []*ApprovalNotificationRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(reqs))

	for _, req := range reqs {
		wg.Add(1)
		go func(r *ApprovalNotificationRequest) {
			defer wg.Done()
			if err := s.SendApproval(ctx, r); err != nil {
				errChan <- fmt.Errorf("发送通知失败[%s]: %v", r.RecipientID, err)
			}
		}(req)
	}

	wg.Wait()
	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {

		return fmt.Errorf("批量发送通知失败: %d个错误", len(errors))
	}

	return nil
}

// SendByTemplate 根据模板发送通知
func (s *notificationService) SendByTemplate(ctx context.Context, req *TemplateNotificationRequest) error {

	// 1. 获取通知模板
	template, err := s.templateService.GetByCode(ctx, req.TemplateCode)
	if err != nil {
		return fmt.Errorf("获取通知模板失败: %v", err)
	}

	// 2. 验证通知类型匹配
	if template.NotificationType != req.NotificationType {
		return fmt.Errorf("模板通知类型不匹配: 期望%s, 实际%s", req.NotificationType, template.NotificationType)
	}

	// 3. 渲染模板
	rendered, err := s.templateService.RenderTemplate(ctx, template.ID, req.Variables)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}

	// 4. 创建通知日志
	notificationLog := &model.NotificationLog{
		ApprovalID:       req.ApprovalID,
		TaskID:           req.TaskID,
		RecipientID:      req.RecipientID,
		RecipientType:    req.RecipientType,
		NotificationType: req.NotificationType,
		TemplateID:       template.ID,
		TemplateCode:     template.TemplateCode,
		Title:            rendered.Title,
		Content:          rendered.Content,
		Status:           model.NotificationStatusPending,
		MaxRetryCount:    3,
	}

	// 设置额外数据
	if req.Variables != nil {
		extraData, _ := json.Marshal(req.Variables)
		notificationLog.ExtraData = string(extraData)
	}

	// 保存通知日志
	if err := s.logRepo.Create(ctx, notificationLog); err != nil {
		return fmt.Errorf("创建通知日志失败: %v", err)
	}

	// 5. 异步发送通知
	go s.sendNotificationAsync(ctx, notificationLog, rendered)

	return nil
}

// ProcessPending 处理待发送通知
func (s *notificationService) ProcessPending(ctx context.Context, limit int) error {
	// 获取待发送通知
	logs, err := s.logRepo.GetPendingLogs(ctx, limit)
	if err != nil {
		return fmt.Errorf("获取待发送通知失败: %v", err)
	}

	if len(logs) == 0 {
		return nil
	}

	// 并发处理
	var wg sync.WaitGroup
	for _, log := range logs {
		wg.Add(1)
		go func(l *model.NotificationLog) {
			defer wg.Done()
			s.processNotificationLog(ctx, l)
		}(log)
	}

	wg.Wait()

	return nil
}

// ProcessRetry 处理重试通知
func (s *notificationService) ProcessRetry(ctx context.Context, limit int) error {
	// 获取需要重试的通知
	logs, err := s.logRepo.GetRetryLogs(ctx, limit)
	if err != nil {
		return fmt.Errorf("获取重试通知失败: %v", err)
	}

	if len(logs) == 0 {
		return nil
	}

	// 并发处理
	var wg sync.WaitGroup
	for _, log := range logs {
		wg.Add(1)
		go func(l *model.NotificationLog) {
			defer wg.Done()
			s.processNotificationLog(ctx, l)
		}(log)
	}

	wg.Wait()

	return nil
}

// GetStatistics 获取通知统计
func (s *notificationService) GetStatistics(ctx context.Context, req *repository.NotificationStatisticsRequest) (*repository.NotificationStatistics, error) {
	return s.logRepo.GetStatistics(ctx, req)
}

// RegisterChannel 注册通知渠道
func (s *notificationService) RegisterChannel(channel NotificationChannel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.channelManager.RegisterChannel(channel)
}

// GetEnabledChannels 获取启用的通知渠道
func (s *notificationService) GetEnabledChannels() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channels := s.channelManager.GetEnabledChannels()
	var types []string
	for _, channel := range channels {
		types = append(types, channel.GetType())
	}
	return types
}

// Test 测试通知发送
func (s *notificationService) Test(ctx context.Context, req *TestNotificationRequest) error {

	// 构建通知请求
	notificationReq := &NotificationRequest{
		RecipientID:   req.RecipientID,
		RecipientType: model.RecipientTypeUser,
		Title:         req.Title,
		Content:       req.Content,
		Variables:     req.Variables,
		Priority:      1,
	}

	// 发送通知
	err := s.channelManager.SendNotification(ctx, req.NotificationType, notificationReq)
	if err != nil {
		return fmt.Errorf("测试通知发送失败: %v", err)
	}

	return nil
}

// sendNotificationAsync 异步发送通知
func (s *notificationService) sendNotificationAsync(ctx context.Context, log *model.NotificationLog, rendered *RenderedTemplate) {
	s.processNotificationLog(ctx, log)
}

// processNotificationLog 处理通知日志
func (s *notificationService) processNotificationLog(ctx context.Context, log *model.NotificationLog) {
	// 构建通知请求
	notificationReq := &NotificationRequest{
		RecipientID:   log.RecipientID,
		RecipientType: log.RecipientType,
		Title:         log.Title,
		Content:       log.Content,
		Priority:      1,
	}

	// 解析额外数据
	if log.ExtraData != "" {
		var extraData map[string]any
		if err := json.Unmarshal([]byte(log.ExtraData), &extraData); err == nil {
			notificationReq.Variables = extraData
		}
	}

	// 发送通知
	err := s.channelManager.SendNotification(ctx, log.NotificationType, notificationReq)
	if err != nil {
		// 标记为失败
		log.MarkAsFailed(err.Error())

	} else {
		// 标记为成功
		log.MarkAsSent()

	}

	// 更新日志状态
	if updateErr := s.logRepo.Update(ctx, log); updateErr != nil {
		s.logger.Error("更新通知日志失败",
			"log_id", log.ID,
			"error", updateErr)
	}
}
