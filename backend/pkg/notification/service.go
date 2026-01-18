package notification

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// notificationService 通知服务实现
type notificationService struct {
	providers map[string]NotificationProvider
	config    *NotificationConfig
	logger    *slog.Logger
	queue     NotificationQueue
	mu        sync.RWMutex
}

// NewNotificationService 创建通知服务
func NewNotificationService(config *NotificationConfig, logger *slog.Logger) NotificationService {
	return &notificationService{
		providers: make(map[string]NotificationProvider),
		config:    config,
		logger:    logger,
	}
}

// RegisterProvider 注册通知提供者
func (s *notificationService) RegisterProvider(provider NotificationProvider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	name := provider.GetName()
	if name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	// 验证配置
	if err := provider.ValidateConfig(); err != nil {
		s.logger.Warn("provider config validation failed",
			"provider", name,
			"error", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.providers[name] = provider

	return nil
}

// GetProvider 获取指定的通知提供者
func (s *notificationService) GetProvider(name string) (NotificationProvider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, exists := s.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}

	return provider, nil
}

// GetEnabledProviders 获取所有启用的通知提供者
func (s *notificationService) GetEnabledProviders() []NotificationProvider {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var enabled []NotificationProvider
	for _, provider := range s.providers {
		if provider.IsEnabled() {
			enabled = append(enabled, provider)
		}
	}

	return enabled
}

// Send 发送通知（使用默认提供者）
func (s *notificationService) Send(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	// 获取第一个启用的提供者作为默认提供者
	enabledProviders := s.GetEnabledProviders()
	if len(enabledProviders) == 0 {
		return nil, fmt.Errorf("no enabled notification providers found")
	}

	return s.SendWithProvider(ctx, enabledProviders[0].GetName(), message)
}

// SendWithProvider 使用指定提供者发送通知
func (s *notificationService) SendWithProvider(ctx context.Context, providerName string, message *NotificationMessage) (*NotificationResult, error) {
	provider, err := s.GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	if !provider.IsEnabled() {
		return nil, fmt.Errorf("provider '%s' is disabled", providerName)
	}

	result, err := provider.Send(ctx, message)
	if err != nil {
		s.logger.Error("failed to send notification",
			"provider", providerName,
			"error", err)
		return result, err
	}

	return result, nil
}

// SendToMultipleProviders 使用多个提供者发送通知
func (s *notificationService) SendToMultipleProviders(ctx context.Context, providerNames []string, message *NotificationMessage) ([]*NotificationResult, error) {
	var results []*NotificationResult
	var errors []error

	for _, providerName := range providerNames {
		result, err := s.SendWithProvider(ctx, providerName, message)
		if err != nil {
			errors = append(errors, fmt.Errorf("provider '%s': %w", providerName, err))
			// 创建失败结果
			result = &NotificationResult{
				Success: false,
				Error:   err,
			}
		}
		results = append(results, result)
	}

	// 如果所有提供者都失败，返回错误
	if len(errors) == len(providerNames) {
		return results, fmt.Errorf("all providers failed: %v", errors)
	}

	return results, nil
}

// SendAsync 异步发送通知
func (s *notificationService) SendAsync(ctx context.Context, message *NotificationMessage) error {
	if s.queue == nil {
		// 如果没有队列，直接同步发送
		_, err := s.Send(ctx, message)
		return err
	}

	// 使用第一个启用的提供者
	enabledProviders := s.GetEnabledProviders()
	if len(enabledProviders) == 0 {
		return fmt.Errorf("no enabled notification providers found")
	}

	return s.queue.Enqueue(ctx, message, enabledProviders[0].GetName())
}

// SetQueue 设置通知队列
func (s *notificationService) SetQueue(queue NotificationQueue) {
	s.queue = queue
}
