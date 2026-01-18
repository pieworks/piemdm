package notification

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// smsProvider 短信提供者（空实现）
type smsProvider struct {
	config *SMSConfig
	logger *slog.Logger
}

// NewSMSProvider 创建短信提供者
func NewSMSProvider(config *SMSConfig, logger *slog.Logger) NotificationProvider {
	return &smsProvider{
		config: config,
		logger: logger,
	}
}

func (p *smsProvider) GetName() string {
	return "sms"
}

func (p *smsProvider) IsEnabled() bool {
	return p.config.Enabled
}

func (p *smsProvider) Send(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	p.logger.Info("SMS provider not implemented yet",
		"to", message.To,
		"content", message.Content)

	// TODO: 实现短信发送逻辑
	// 根据配置的provider选择对应的短信服务商
	// - Aliyun SMS
	// - Tencent SMS
	// - Custom API

	return &NotificationResult{
		Success: false,
		Error:   fmt.Errorf("SMS provider not implemented"),
		SentAt:  time.Now(),
		Details: map[string]any{
			"provider": "sms",
			"status":   "not_implemented",
		},
	}, fmt.Errorf("SMS provider not implemented")
}

func (p *smsProvider) ValidateConfig() error {
	if !p.config.Enabled {
		return nil
	}

	// TODO: 根据不同的provider验证配置
	p.logger.Info("SMS config validation not implemented")
	return nil
}

func (p *smsProvider) GetSupportedTypes() []string {
	return []string{"text"}
}

// feishuProvider 飞书提供者（空实现）
type feishuProvider struct {
	config *FeishuConfig
	logger *slog.Logger
}

// NewFeishuProvider 创建飞书提供者
func NewFeishuProvider(config *FeishuConfig, logger *slog.Logger) NotificationProvider {
	return &feishuProvider{
		config: config,
		logger: logger,
	}
}

func (p *feishuProvider) GetName() string {
	return "feishu"
}

func (p *feishuProvider) IsEnabled() bool {
	return p.config.Enabled
}

func (p *feishuProvider) Send(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	p.logger.Info("Feishu provider not implemented yet",
		"to", message.To,
		"content", message.Content)

	// TODO: 实现飞书消息发送逻辑
	// 使用飞书开放平台API
	// https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN

	return &NotificationResult{
		Success: false,
		Error:   fmt.Errorf("Feishu provider not implemented"),
		SentAt:  time.Now(),
		Details: map[string]any{
			"provider": "feishu",
			"status":   "not_implemented",
		},
	}, fmt.Errorf("Feishu provider not implemented")
}

func (p *feishuProvider) ValidateConfig() error {
	if !p.config.Enabled {
		return nil
	}

	if p.config.AppID == "" {
		return fmt.Errorf("Feishu app ID is required")
	}
	if p.config.AppSecret == "" {
		return fmt.Errorf("Feishu app secret is required")
	}

	return nil
}

func (p *feishuProvider) GetSupportedTypes() []string {
	return []string{"text", "markdown"}
}

// dingtalkProvider 钉钉提供者（空实现）
type dingtalkProvider struct {
	config *DingtalkConfig
	logger *slog.Logger
}

// NewDingtalkProvider 创建钉钉提供者
func NewDingtalkProvider(config *DingtalkConfig, logger *slog.Logger) NotificationProvider {
	return &dingtalkProvider{
		config: config,
		logger: logger,
	}
}

func (p *dingtalkProvider) GetName() string {
	return "dingtalk"
}

func (p *dingtalkProvider) IsEnabled() bool {
	return p.config.Enabled
}

func (p *dingtalkProvider) Send(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	p.logger.Info("Dingtalk provider not implemented yet",
		"to", message.To,
		"content", message.Content)

	// TODO: 实现钉钉消息发送逻辑
	// 使用钉钉开放平台API
	// https://developers.dingtalk.com/document/app/custom-robot-access

	return &NotificationResult{
		Success: false,
		Error:   fmt.Errorf("Dingtalk provider not implemented"),
		SentAt:  time.Now(),
		Details: map[string]any{
			"provider": "dingtalk",
			"status":   "not_implemented",
		},
	}, fmt.Errorf("Dingtalk provider not implemented")
}

func (p *dingtalkProvider) ValidateConfig() error {
	if !p.config.Enabled {
		return nil
	}

	if p.config.AppKey == "" {
		return fmt.Errorf("Dingtalk app key is required")
	}
	if p.config.AppSecret == "" {
		return fmt.Errorf("Dingtalk app secret is required")
	}

	return nil
}

func (p *dingtalkProvider) GetSupportedTypes() []string {
	return []string{"text", "markdown"}
}

// customProvider 自定义提供者（空实现）
type customProvider struct {
	config *CustomConfig
	logger *slog.Logger
}

// NewCustomProvider 创建自定义提供者
func NewCustomProvider(config *CustomConfig, logger *slog.Logger) NotificationProvider {
	return &customProvider{
		config: config,
		logger: logger,
	}
}

func (p *customProvider) GetName() string {
	return "custom"
}

func (p *customProvider) IsEnabled() bool {
	return p.config.Enabled
}

func (p *customProvider) Send(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	p.logger.Info("Custom provider not implemented yet",
		"to", message.To,
		"content", message.Content)

	// TODO: 实现自定义HTTP接口调用逻辑
	// 支持多个自定义端点
	// 支持不同的HTTP方法和认证方式

	return &NotificationResult{
		Success: false,
		Error:   fmt.Errorf("Custom provider not implemented"),
		SentAt:  time.Now(),
		Details: map[string]any{
			"provider": "custom",
			"status":   "not_implemented",
		},
	}, fmt.Errorf("Custom provider not implemented")
}

func (p *customProvider) ValidateConfig() error {
	if !p.config.Enabled {
		return nil
	}

	if len(p.config.Endpoints) == 0 {
		return fmt.Errorf("at least one custom endpoint is required")
	}

	for i, endpoint := range p.config.Endpoints {
		if endpoint.Name == "" {
			return fmt.Errorf("endpoint %d: name is required", i)
		}
		if endpoint.URL == "" {
			return fmt.Errorf("endpoint %d: URL is required", i)
		}
		if endpoint.Method == "" {
			endpoint.Method = "POST" // 默认POST方法
		}
	}

	return nil
}

func (p *customProvider) GetSupportedTypes() []string {
	return []string{"text", "html", "json"}
}
