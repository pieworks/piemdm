package notification

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// emailProvider 邮件提供者
type emailProvider struct {
	config *EmailConfig
	logger *slog.Logger
}

// NewEmailProvider 创建邮件提供者
func NewEmailProvider(config *EmailConfig, logger *slog.Logger) NotificationProvider {
	return &emailProvider{
		config: config,
		logger: logger,
	}
}

// GetName 获取提供者名称
func (p *emailProvider) GetName() string {
	return "email"
}

// IsEnabled 检查是否启用
func (p *emailProvider) IsEnabled() bool {
	return p.config.Enabled
}

// Send 发送邮件
func (p *emailProvider) Send(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	if !p.IsEnabled() {
		return nil, fmt.Errorf("email provider is disabled")
	}

	return p.sendViaSMTP(ctx, message)
}

// sendViaSMTP 通过SMTP发送邮件
func (p *emailProvider) sendViaSMTP(ctx context.Context, message *NotificationMessage) (*NotificationResult, error) {
	config := p.config

	// 构建邮件内容
	from := config.From
	if config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", config.FromName, config.From)
	}

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(message.To, ", ")
	headers["Subject"] = message.Subject
	headers["MIME-Version"] = "1.0"

	// 设置内容类型
	contentType := "text/plain; charset=UTF-8"
	if message.ContentType == "html" {
		contentType = "text/html; charset=UTF-8"
	}
	headers["Content-Type"] = contentType

	// 构建邮件正文
	var emailBody strings.Builder
	for key, value := range headers {
		emailBody.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	emailBody.WriteString("\r\n")
	emailBody.WriteString(message.Content)

	// SMTP认证 - 根据服务器类型选择合适的认证方式
	var auth smtp.Auth

	// 对于Outlook/Hotmail，使用CRAM-MD5认证
	if strings.Contains(strings.ToLower(config.Host), "outlook") ||
		strings.Contains(strings.ToLower(config.Host), "hotmail") ||
		strings.Contains(strings.ToLower(config.Host), "live") ||
		strings.Contains(strings.ToLower(config.Host), "exchangelabs") {
		// 尝试CRAM-MD5认证
		auth = smtp.CRAMMD5Auth(config.Username, config.Password)
	} else {
		// 其他服务器使用PLAIN认证
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	var err error
	if config.UseTLS {
		err = p.sendMailTLS(addr, auth, config.From, message.To, []byte(emailBody.String()))
	} else {
		err = smtp.SendMail(addr, auth, config.From, message.To, []byte(emailBody.String()))
	}

	// 如果认证失败，尝试其他认证方式
	if err != nil && strings.Contains(err.Error(), "Unrecognized authentication type") {

		// 尝试PLAIN认证作为回退
		fallbackAuth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
		if config.UseTLS {
			err = p.sendMailTLS(addr, fallbackAuth, config.From, message.To, []byte(emailBody.String()))
		} else {
			err = smtp.SendMail(addr, fallbackAuth, config.From, message.To, []byte(emailBody.String()))
		}

		if err != nil {
			p.logger.Error("fallback authentication also failed",
				"error", err,
				"host", config.Host)
		}
	}

	if err != nil {
		p.logger.Error("failed to send email via SMTP",
			"error", err,
			"host", config.Host,
			"port", config.Port)
		return &NotificationResult{
			Success: false,
			Error:   err,
			SentAt:  time.Now(),
		}, err
	}

	messageID := fmt.Sprintf("smtp-%d", time.Now().UnixNano())

	return &NotificationResult{
		Success:   true,
		MessageID: messageID,
		SentAt:    time.Now(),
		Details: map[string]any{
			"provider": "smtp",
			"host":     config.Host,
			"port":     config.Port,
		},
	}, nil
}

// sendMailTLS 使用TLS发送邮件
func (p *emailProvider) sendMailTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 解析主机和端口
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address format: %s", addr)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port: %s", parts[1])
	}

	// 根据端口选择不同的连接方式
	if port == 465 {
		// 465 端口使用 SSL/TLS 直连
		return p.sendMailSSL(addr, auth, from, to, msg)
	} else {
		// 587 端口使用 STARTTLS
		return p.sendMailSTARTTLS(addr, auth, from, to, msg)
	}
}

// sendMailSSL 使用SSL直连发送邮件（适用于465端口）
func (p *emailProvider) sendMailSSL(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	parts := strings.Split(addr, ":")
	host := parts[0]

	// 创建TLS连接配置
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}

	// 对于QQ邮箱，设置更宽松的TLS配置
	if strings.Contains(host, "qq.com") {
		tlsConfig.InsecureSkipVerify = false
		tlsConfig.CipherSuites = []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		}
	}

	// 创建SSL连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		p.logger.Error("failed to create SSL connection",
			"error", err,
			"host", host)
		return fmt.Errorf("SSL connection failed: %w", err)
	}
	defer conn.Close()

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		p.logger.Error("failed to create SMTP client",
			"error", err,
			"host", host)
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Quit()

	// 认证
	if auth != nil {
		p.logger.Debug("attempting SMTP authentication")
		if err = client.Auth(auth); err != nil {
			p.logger.Error("SMTP authentication failed",
				"error", err,
				"host", host)
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
		p.logger.Error("SMTP authentication successful")
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		p.logger.Error("failed to set sender",
			"error", err,
			"from", from)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			p.logger.Error("failed to set recipient",
				"error", err,
				"recipient", recipient)
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		p.logger.Error("failed to get data writer", "error", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write(msg)
	if err != nil {
		p.logger.Error("failed to write message data", "error", err)
		return fmt.Errorf("failed to write message data: %w", err)
	}

	p.logger.Debug("email sent successfully via SSL")
	return nil
}

// sendMailSTARTTLS 使用STARTTLS发送邮件（适用于587端口）
func (p *emailProvider) sendMailSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	parts := strings.Split(addr, ":")
	host := parts[0]

	// 创建普通TCP连接
	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		p.logger.Error("failed to create TCP connection",
			"error", err,
			"host", host)
		return fmt.Errorf("TCP connection failed: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		p.logger.Error("failed to create SMTP client",
			"error", err,
			"host", host)
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Quit()

	// 启动TLS
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		p.logger.Error("failed to start TLS",
			"error", err,
			"host", host)
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			p.logger.Error("SMTP authentication failed",
				"error", err,
				"host", host)
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
		p.logger.Error("SMTP authentication successful")
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		p.logger.Error("failed to set sender",
			"error", err,
			"from", from)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			p.logger.Error("failed to set recipient",
				"error", err,
				"recipient", recipient)
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		p.logger.Error("failed to get data writer", "error", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write(msg)
	if err != nil {
		p.logger.Error("failed to write message data", "error", err)
		return fmt.Errorf("failed to write message data: %w", err)
	}

	p.logger.Debug("email sent successfully via STARTTLS")
	return nil
}

// ValidateConfig 验证配置
func (p *emailProvider) ValidateConfig() error {
	if !p.config.Enabled {
		return nil // 如果未启用，不需要验证
	}

	return p.validateSMTPConfig()
}

// validateSMTPConfig 验证SMTP配置
func (p *emailProvider) validateSMTPConfig() error {
	config := p.config

	if config.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.Port <= 0 {
		return fmt.Errorf("SMTP port must be greater than 0")
	}
	if config.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if config.From == "" {
		return fmt.Errorf("SMTP from address is required")
	}

	return nil
}

// GetSupportedTypes 获取支持的内容类型
func (p *emailProvider) GetSupportedTypes() []string {
	return []string{"text", "html"}
}
