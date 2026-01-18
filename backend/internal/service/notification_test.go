package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"piemdm/internal/model"
	"piemdm/internal/service"
	mock_repository "piemdm/test/mocks/repository"
	mock_service "piemdm/test/mocks/service"
	"piemdm/test/testutil"
)

func TestNotificationService_SendApprovalNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateService := mock_service.NewMockNotificationTemplateService(ctrl)
	mockLogRepo := mock_repository.NewMockNotificationLogRepository(ctrl)

	logger := testutil.CreateTestLogger()

	// 创建服务实例
	notificationService := service.NewNotificationService(
		mockTemplateService,
		mockLogRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	req := &service.ApprovalNotificationRequest{
		ApprovalID:       "approval_123",
		TaskID:           "task_456",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		RecipientID:      "user_789",
		RecipientType:    "user",
		Variables: map[string]any{
			"applicant_name": "张三",
			"approval_title": "请假申请",
		},
		Priority: 1,
	}

	template := &model.NotificationTemplate{
		ID:               "template_123",
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Status:           "Normal",
	}

	rendered := &service.RenderedTemplate{
		Title:   "您有新的审批待办：请假申请",
		Content: "申请人：张三，请及时处理。",
	}

	// 设置 Mock 期望
	mockTemplateService.EXPECT().
		GetByTypeAndNotification(ctx, req.TemplateType, req.NotificationType).
		Return(template, nil)

	mockTemplateService.EXPECT().
		RenderTemplate(ctx, template.ID, req.Variables).
		Return(rendered, nil)

	mockLogRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, log *model.NotificationLog) error {
			// 验证通知日志的字段
			assert.Equal(t, req.ApprovalID, log.ApprovalID)
			assert.Equal(t, req.TaskID, log.TaskID)
			assert.Equal(t, req.RecipientID, log.RecipientID)
			assert.Equal(t, req.RecipientType, log.RecipientType)
			assert.Equal(t, req.NotificationType, log.NotificationType)
			assert.Equal(t, template.ID, log.TemplateID)
			assert.Equal(t, template.TemplateCode, log.TemplateCode)
			assert.Equal(t, rendered.Title, log.Title)
			assert.Equal(t, rendered.Content, log.Content)
			assert.Equal(t, model.NotificationStatusPending, log.Status)
			assert.Equal(t, 3, log.MaxRetryCount)

			// 设置ID模拟数据库自增
			log.ID = 1
			return nil
		})

	// 为异步处理添加Update期望（可能会被调用）
	mockLogRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		AnyTimes()

	// 执行测试
	err := notificationService.SendApproval(ctx, req)

	// 验证结果
	assert.NoError(t, err)
}

func TestNotificationService_SendBatchNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateService := mock_service.NewMockNotificationTemplateService(ctrl)
	mockLogRepo := mock_repository.NewMockNotificationLogRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	notificationService := service.NewNotificationService(
		mockTemplateService,
		mockLogRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	reqs := []*service.ApprovalNotificationRequest{
		{
			ApprovalID:       "approval_123",
			TemplateType:     model.TemplateTypeApprovalPending,
			NotificationType: model.NotificationTypeEmail,
			RecipientID:      "user_789",
			RecipientType:    "user",
			Variables: map[string]any{
				"applicant_name": "张三",
				"approval_title": "请假申请",
			},
		},
		{
			ApprovalID:       "approval_456",
			TemplateType:     model.TemplateTypeApprovalPending,
			NotificationType: model.NotificationTypeEmail,
			RecipientID:      "user_101",
			RecipientType:    "user",
			Variables: map[string]any{
				"applicant_name": "李四",
				"approval_title": "报销申请",
			},
		},
	}

	template := &model.NotificationTemplate{
		ID:               "template_123",
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Status:           "Normal",
	}

	rendered1 := &service.RenderedTemplate{
		Title:   "您有新的审批待办：请假申请",
		Content: "申请人：张三，请及时处理。",
	}

	rendered2 := &service.RenderedTemplate{
		Title:   "您有新的审批待办：报销申请",
		Content: "申请人：李四，请及时处理。",
	}

	// 设置 Mock 期望 - 每个请求都会调用
	mockTemplateService.EXPECT().
		GetByTypeAndNotification(ctx, model.TemplateTypeApprovalPending, model.NotificationTypeEmail).
		Return(template, nil).
		Times(2)

	mockTemplateService.EXPECT().
		RenderTemplate(ctx, template.ID, reqs[0].Variables).
		Return(rendered1, nil)

	mockTemplateService.EXPECT().
		RenderTemplate(ctx, template.ID, reqs[1].Variables).
		Return(rendered2, nil)

	mockLogRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil).
		Times(2)

	// 为异步处理添加Update期望（可能会被调用）
	mockLogRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		AnyTimes()

	// 执行测试
	err := notificationService.SendBatch(ctx, reqs)

	// 验证结果
	assert.NoError(t, err)
}

func TestNotificationService_SendNotificationByTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateService := mock_service.NewMockNotificationTemplateService(ctrl)
	mockLogRepo := mock_repository.NewMockNotificationLogRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	notificationService := service.NewNotificationService(
		mockTemplateService,
		mockLogRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	req := &service.TemplateNotificationRequest{
		TemplateCode:     "approval_pending_email",
		NotificationType: model.NotificationTypeEmail,
		RecipientID:      "user_789",
		RecipientType:    "user",
		Variables: map[string]any{
			"applicant_name": "张三",
			"approval_title": "请假申请",
		},
		ApprovalID: "approval_123",
		TaskID:     "task_456",
		Priority:   1,
	}

	template := &model.NotificationTemplate{
		ID:               "template_123",
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Status:           "Normal",
	}

	rendered := &service.RenderedTemplate{
		Title:   "您有新的审批待办：请假申请",
		Content: "申请人：张三，请及时处理。",
	}

	// 设置 Mock 期望
	mockTemplateService.EXPECT().
		GetByCode(ctx, req.TemplateCode).
		Return(template, nil)

	mockTemplateService.EXPECT().
		RenderTemplate(ctx, template.ID, req.Variables).
		Return(rendered, nil)

	mockLogRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, log *model.NotificationLog) error {
			// 验证通知日志的字段
			assert.Equal(t, req.ApprovalID, log.ApprovalID)
			assert.Equal(t, req.TaskID, log.TaskID)
			assert.Equal(t, req.RecipientID, log.RecipientID)
			assert.Equal(t, req.RecipientType, log.RecipientType)
			assert.Equal(t, req.NotificationType, log.NotificationType)
			assert.Equal(t, template.ID, log.TemplateID)
			assert.Equal(t, template.TemplateCode, log.TemplateCode)
			assert.Equal(t, rendered.Title, log.Title)
			assert.Equal(t, rendered.Content, log.Content)
			assert.Equal(t, model.NotificationStatusPending, log.Status)

			// 设置ID模拟数据库自增
			log.ID = 1
			return nil
		})

	// 为异步处理添加Update期望（可能会被调用）
	mockLogRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		AnyTimes()

	// 执行测试
	err := notificationService.SendByTemplate(ctx, req)

	// 验证结果
	assert.NoError(t, err)
}

func TestNotificationService_ProcessPendingNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateService := mock_service.NewMockNotificationTemplateService(ctrl)
	mockLogRepo := mock_repository.NewMockNotificationLogRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	notificationService := service.NewNotificationService(
		mockTemplateService,
		mockLogRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	pendingLogs := []*model.NotificationLog{
		{
			ID:               1,
			ApprovalID:       "approval_123",
			RecipientID:      "user_789",
			RecipientType:    "user",
			NotificationType: model.NotificationTypeEmail,
			TemplateID:       "template_123",
			TemplateCode:     "approval_pending_email",
			Title:            "您有新的审批待办：请假申请",
			Content:          "申请人：张三，请及时处理。",
			Status:           model.NotificationStatusPending,
			RetryCount:       0,
			MaxRetryCount:    3,
		},
	}

	// 设置 Mock 期望 - 使用正确的方法名
	mockLogRepo.EXPECT().
		GetPendingLogs(ctx, 10).
		Return(pendingLogs, nil)

	// 为异步处理添加Update期望（可能会被调用）
	mockLogRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		AnyTimes()

	// 由于异步处理，我们只验证获取待发送通知的调用
	// 实际的发送逻辑在异步协程中执行

	// 执行测试
	err := notificationService.ProcessPending(ctx, 10)

	// 验证结果
	assert.NoError(t, err)
}

func TestNotificationService_TestNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateService := mock_service.NewMockNotificationTemplateService(ctrl)
	mockLogRepo := mock_repository.NewMockNotificationLogRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	notificationService := service.NewNotificationService(
		mockTemplateService,
		mockLogRepo,
		logger,
	)

	// 注册邮件通知渠道
	emailConfig := service.EmailConfig{
		Enabled:   true,
		SMTPHost:  "smtp.test.com",
		SMTPPort:  587,
		Username:  "test@test.com",
		Password:  "password",
		FromEmail: "test@test.com",
		FromName:  "Test System",
	}
	emailChannel := service.NewEmailChannel(emailConfig, logger)
	notificationService.RegisterChannel(emailChannel)

	ctx := context.Background()

	// 测试数据
	req := &service.TestNotificationRequest{
		NotificationType: model.NotificationTypeEmail,
		RecipientID:      "user_789",
		Title:            "测试通知",
		Content:          "这是一条测试通知",
		Variables: map[string]any{
			"test_var": "test_value",
		},
	}

	// 执行测试
	err := notificationService.Test(ctx, req)

	// 验证结果 - 测试通知应该成功（即使没有实际发送）
	assert.NoError(t, err)
}

func TestNotificationService_RegisterChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateService := mock_service.NewMockNotificationTemplateService(ctrl)
	mockLogRepo := mock_repository.NewMockNotificationLogRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	notificationService := service.NewNotificationService(
		mockTemplateService,
		mockLogRepo,
		logger,
	)

	// 创建一个测试通知渠道 - 使用正确的构造函数
	emailConfig := service.EmailConfig{
		Enabled:   true,
		SMTPHost:  "smtp.test.com",
		SMTPPort:  587,
		Username:  "test@test.com",
		Password:  "password",
		FromEmail: "test@test.com",
		FromName:  "Test System",
	}
	testChannel := service.NewEmailChannel(emailConfig, testutil.CreateTestLogger())

	// 执行测试
	err := notificationService.RegisterChannel(testChannel)

	// 验证结果
	assert.NoError(t, err)

	// 验证渠道已注册
	channels := notificationService.GetEnabledChannels()
	assert.Contains(t, channels, model.NotificationTypeEmail)
}
