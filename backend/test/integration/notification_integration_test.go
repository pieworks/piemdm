package integration

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"piemdm/pkg/log"
	"piemdm/pkg/notification"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationServiceIntegration(t *testing.T) {
	// 创建测试配置
	config := viper.New()
	config.Set("notification.email.enabled", true)
	config.Set("notification.email.host", "smtp.example.com")
	config.Set("notification.email.port", 587)
	config.Set("notification.email.username", "test@example.com")
	config.Set("notification.email.password", "password")
	config.Set("notification.email.from", "test@example.com")
	config.Set("notification.email.from_name", "Test System")
	config.Set("notification.email.use_tls", true)

	// 创建logger
	config.Set("log.log_level", "debug")
	config.Set("log.encoding", "console")
	config.Set("log.log_file_name", "./test.log")
	config.Set("log.max_size", 1024)
	config.Set("log.max_backups", 30)
	config.Set("log.max_age", 7)
	config.Set("log.compress", true)
	config.Set("env", "test")

	// 使用 discard logger 来抑制预期的错误日志
	// 因为我们使用的是假 SMTP 配置，连接错误是预期的，不需要打印到控制台
	discardHandler := slog.NewTextHandler(io.Discard, nil)
	logger := &log.Logger{Logger: slog.New(discardHandler)}

	// 测试通知服务创建
	service := notification.NewNotificationServiceProvider(config, logger)
	require.NotNil(t, service, "通知服务应该成功创建")

	// 测试获取启用的提供者
	providers := service.GetEnabledProviders()
	assert.Len(t, providers, 1, "应该有一个启用的提供者")
	assert.Equal(t, "email", providers[0].GetName(), "应该是邮件提供者")

	// 测试创建通知消息
	message := notification.CreateCCNotificationMessage(
		[]string{"test@example.com"},
		"测试审批",
		"张三",
		"已完成",
	)

	assert.NotNil(t, message, "通知消息应该成功创建")
	assert.Equal(t, []string{"test@example.com"}, message.To)
	assert.Contains(t, message.Subject, "测试审批")
	assert.Contains(t, message.Content, "张三")
	assert.Equal(t, "html", message.ContentType)

	// 测试发送通知（这里会失败，因为SMTP配置是假的，但不应该panic）
	ctx := context.Background()
	result, err := service.Send(ctx, message)

	// 预期会失败，因为SMTP配置是测试用的
	assert.Error(t, err, "使用假SMTP配置应该失败")
	assert.NotNil(t, result, "即使失败也应该返回结果")
	assert.False(t, result.Success, "结果应该标记为失败")

	t.Log("通知服务集成测试完成")
}

func TestNotificationMessageCreation(t *testing.T) {
	// 测试审批通知消息创建
	approvalMessage := notification.CreateApprovalNotificationMessage(
		[]string{"approver@example.com"},
		"用户权限申请",
		"李四",
		"部门经理审批",
	)

	assert.NotNil(t, approvalMessage)
	assert.Equal(t, []string{"approver@example.com"}, approvalMessage.To)
	assert.Contains(t, approvalMessage.Subject, "审批通知")
	assert.Contains(t, approvalMessage.Content, "用户权限申请")
	assert.Contains(t, approvalMessage.Content, "李四")
	assert.Contains(t, approvalMessage.Content, "部门经理审批")
	assert.Equal(t, "html", approvalMessage.ContentType)
	assert.Equal(t, 3, approvalMessage.Priority) // 高优先级

	// 测试抄送通知消息创建
	ccMessage := notification.CreateCCNotificationMessage(
		[]string{"cc1@example.com", "cc2@example.com"},
		"设备采购申请",
		"王五",
		"已通过",
	)

	assert.NotNil(t, ccMessage)
	assert.Equal(t, []string{"cc1@example.com", "cc2@example.com"}, ccMessage.To)
	assert.Contains(t, ccMessage.Subject, "抄送通知")
	assert.Contains(t, ccMessage.Content, "设备采购申请")
	assert.Contains(t, ccMessage.Content, "王五")
	assert.Contains(t, ccMessage.Content, "已通过")
	assert.Equal(t, "html", ccMessage.ContentType)
	assert.Equal(t, 2, ccMessage.Priority) // 中等优先级

	// 检查元数据
	assert.Equal(t, "approval", approvalMessage.Metadata["type"])
	assert.Equal(t, "cc", ccMessage.Metadata["type"])
}
