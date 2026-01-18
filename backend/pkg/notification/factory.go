package notification

import (
	"fmt"
	"log/slog"
)

// NotificationFactory 通知工厂
type NotificationFactory struct {
	config *NotificationConfig
	logger *slog.Logger
}

// NewNotificationFactory 创建通知工厂
func NewNotificationFactory(config *NotificationConfig, logger *slog.Logger) *NotificationFactory {
	return &NotificationFactory{
		config: config,
		logger: logger,
	}
}

// CreateService 创建通知服务
func (f *NotificationFactory) CreateService() (NotificationService, error) {
	service := NewNotificationService(f.config, f.logger)

	// 注册所有提供者
	if err := f.registerProviders(service); err != nil {
		return nil, fmt.Errorf("failed to register providers: %w", err)
	}

	return service, nil
}

// registerProviders 注册所有提供者
func (f *NotificationFactory) registerProviders(service NotificationService) error {
	// 注册邮件提供者
	if f.config.Email.Enabled {
		emailProvider := NewEmailProvider(&f.config.Email, f.logger)
		if err := service.RegisterProvider(emailProvider); err != nil {
			f.logger.Error("failed to register email provider", "error", err)
			return err
		}
	}

	// 注册短信提供者
	if f.config.SMS.Enabled {
		smsProvider := NewSMSProvider(&f.config.SMS, f.logger)
		if err := service.RegisterProvider(smsProvider); err != nil {
			f.logger.Error("failed to register SMS provider", "error", err)
			return err
		}
	}

	// 注册飞书提供者
	if f.config.Feishu.Enabled {
		feishuProvider := NewFeishuProvider(&f.config.Feishu, f.logger)
		if err := service.RegisterProvider(feishuProvider); err != nil {
			f.logger.Error("failed to register Feishu provider", "error", err)
			return err
		}
	}

	// 注册钉钉提供者
	if f.config.Dingtalk.Enabled {
		dingtalkProvider := NewDingtalkProvider(&f.config.Dingtalk, f.logger)
		if err := service.RegisterProvider(dingtalkProvider); err != nil {
			f.logger.Error("failed to register Dingtalk provider", "error", err)
			return err
		}
	}

	// 注册自定义提供者
	if f.config.Custom.Enabled {
		customProvider := NewCustomProvider(&f.config.Custom, f.logger)
		if err := service.RegisterProvider(customProvider); err != nil {
			f.logger.Error("failed to register custom provider", "error", err)
			return err
		}
	}

	return nil
}

// CreateEmailMessage 创建邮件消息
func CreateEmailMessage(to []string, subject, content string) *NotificationMessage {
	return &NotificationMessage{
		To:          to,
		Subject:     subject,
		Content:     content,
		ContentType: "html",
		Priority:    2, // 中等优先级
	}
}

// CreateTextMessage 创建文本消息
func CreateTextMessage(to []string, content string) *NotificationMessage {
	return &NotificationMessage{
		To:          to,
		Content:     content,
		ContentType: "text",
		Priority:    2, // 中等优先级
	}
}

// CreateApprovalNotificationMessage 创建审批通知消息
func CreateApprovalNotificationMessage(to []string, approvalTitle, applicant, currentNode string) *NotificationMessage {
	subject := fmt.Sprintf("【审批通知】%s", approvalTitle)

	content := fmt.Sprintf(`
<html>
<body>
<h3>审批通知</h3>
<p>您好，</p>
<p>您有一个新的审批任务需要处理：</p>
<ul>
<li><strong>审批标题：</strong>%s</li>
<li><strong>申请人：</strong>%s</li>
<li><strong>当前节点：</strong>%s</li>
<li><strong>通知时间：</strong>%s</li>
</ul>
<p>请及时登录系统处理。</p>
<p>此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
	`, approvalTitle, applicant, currentNode, fmt.Sprintf("%s", "现在"))

	return &NotificationMessage{
		To:          to,
		Subject:     subject,
		Content:     content,
		ContentType: "html",
		Priority:    3, // 高优先级
		Metadata: map[string]string{
			"type":         "approval",
			"approval":     approvalTitle,
			"applicant":    applicant,
			"current_node": currentNode,
		},
	}
}

// CreateCCNotificationMessage 创建抄送通知消息
func CreateCCNotificationMessage(to []string, approvalTitle, applicant, status string) *NotificationMessage {
	subject := fmt.Sprintf("【抄送通知】%s", approvalTitle)

	content := fmt.Sprintf(`
<html>
<body>
<h3>抄送通知</h3>
<p>您好，</p>
<p>以下审批流程已完成，特此抄送：</p>
<ul>
<li><strong>审批标题：</strong>%s</li>
<li><strong>申请人：</strong>%s</li>
<li><strong>审批状态：</strong>%s</li>
<li><strong>通知时间：</strong>%s</li>
</ul>
<p>如需查看详情，请登录系统查看。</p>
<p>此邮件由系统自动发送，请勿回复。</p>
</body>
</html>
	`, approvalTitle, applicant, status, fmt.Sprintf("%s", "现在"))

	return &NotificationMessage{
		To:          to,
		Subject:     subject,
		Content:     content,
		ContentType: "html",
		Priority:    2, // 中等优先级
		Metadata: map[string]string{
			"type":      "cc",
			"approval":  approvalTitle,
			"applicant": applicant,
			"status":    status,
		},
	}
}
