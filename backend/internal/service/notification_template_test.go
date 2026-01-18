package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"piemdm/internal/model"
	"piemdm/internal/service"
	mock_repository "piemdm/test/mocks/repository"
	"piemdm/test/testutil"
)

func TestNotificationTemplateService_CreateTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateRepo := mock_repository.NewMockNotificationTemplateRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	templateService := service.NewNotificationTemplateService(
		mockTemplateRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	req := &service.CreateNotificationTemplateRequest{
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Variables: map[string]any{
			"approval_title": "审批标题",
			"applicant_name": "申请人姓名",
		},
		Description: "审批待办邮件通知模板",
		Status:      "Normal",
		CreatedBy:   "admin",
	}

	// 设置 Mock 期望
	mockTemplateRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, template *model.NotificationTemplate) error {
			// 验证模板字段
			assert.Equal(t, req.TemplateCode, template.TemplateCode)
			assert.Equal(t, req.TemplateName, template.TemplateName)
			assert.Equal(t, req.TemplateType, template.TemplateType)
			assert.Equal(t, req.NotificationType, template.NotificationType)
			assert.Equal(t, req.TitleTemplate, template.TitleTemplate)
			assert.Equal(t, req.ContentTemplate, template.ContentTemplate)
			assert.Equal(t, req.Description, template.Description)
			assert.Equal(t, req.Status, template.Status)
			assert.Equal(t, req.CreatedBy, template.CreatedBy)
			assert.NotEmpty(t, template.ID)

			return nil
		})

	// 执行测试
	result, err := templateService.Create(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.TemplateCode, result.TemplateCode)
	assert.Equal(t, req.TemplateName, result.TemplateName)
}

func TestNotificationTemplateService_UpdateTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateRepo := mock_repository.NewMockNotificationTemplateRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	templateService := service.NewNotificationTemplateService(
		mockTemplateRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	templateID := "template_123"
	req := &service.UpdateNotificationTemplateRequest{
		ID:              templateID,
		TemplateName:    "更新后的审批待办邮件模板",
		TitleTemplate:   "您有新的审批待办：{{.approval_title}}（已更新）",
		ContentTemplate: "申请人：{{.applicant_name}}，请及时处理。（已更新）",
		Variables: map[string]any{
			"approval_title": "审批标题",
			"applicant_name": "申请人姓名",
		},
		Description: "更新后的审批待办邮件通知模板",
		Status:      "Normal",
		UpdatedBy:   "admin",
	}

	existingTemplate := &model.NotificationTemplate{
		ID:               templateID,
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Status:           "Normal",
	}

	// 设置 Mock 期望
	mockTemplateRepo.EXPECT().
		FindOne(ctx, templateID).
		Return(existingTemplate, nil)

	mockTemplateRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, template *model.NotificationTemplate) error {
			// 验证更新的字段
			assert.Equal(t, req.TemplateName, template.TemplateName)
			assert.Equal(t, req.TitleTemplate, template.TitleTemplate)
			assert.Equal(t, req.ContentTemplate, template.ContentTemplate)
			assert.Equal(t, req.Description, template.Description)
			assert.Equal(t, req.Status, template.Status)
			assert.Equal(t, req.UpdatedBy, template.UpdatedBy)

			return nil
		})

	// 执行测试
	result, err := templateService.Update(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.TemplateName, result.TemplateName)
	assert.Equal(t, req.TitleTemplate, result.TitleTemplate)
}

func TestNotificationTemplateService_GetTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateRepo := mock_repository.NewMockNotificationTemplateRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	templateService := service.NewNotificationTemplateService(
		mockTemplateRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	templateID := "template_123"
	expectedTemplate := &model.NotificationTemplate{
		ID:               templateID,
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Status:           "Normal",
	}

	// 设置 Mock 期望
	mockTemplateRepo.EXPECT().
		FindOne(ctx, templateID).
		Return(expectedTemplate, nil)

	// 执行测试
	result, err := templateService.Get(ctx, templateID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTemplate.ID, result.ID)
	assert.Equal(t, expectedTemplate.TemplateCode, result.TemplateCode)
	assert.Equal(t, expectedTemplate.TemplateName, result.TemplateName)
}

func TestNotificationTemplateService_RenderTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateRepo := mock_repository.NewMockNotificationTemplateRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	templateService := service.NewNotificationTemplateService(
		mockTemplateRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	templateID := "template_123"
	template := &model.NotificationTemplate{
		ID:               templateID,
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Status:           "Normal",
	}

	variables := map[string]any{
		"approval_title": "请假申请",
		"applicant_name": "张三",
	}

	// 设置 Mock 期望
	mockTemplateRepo.EXPECT().
		FindOne(ctx, templateID).
		Return(template, nil)

	// 执行测试
	result, err := templateService.RenderTemplate(ctx, templateID, variables)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "您有新的审批待办：请假申请", result.Title)
	assert.Equal(t, "申请人：张三，请及时处理。", result.Content)
}

func TestNotificationTemplateService_ValidateTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateRepo := mock_repository.NewMockNotificationTemplateRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	templateService := service.NewNotificationTemplateService(
		mockTemplateRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据 - 有效的模板
	validReq := &service.ValidateTemplateRequest{
		TitleTemplate:   "您有新的审批待办：{{.approval_title}}",
		ContentTemplate: "申请人：{{.applicant_name}}，请及时处理。",
		Variables: map[string]any{
			"approval_title": "请假申请",
			"applicant_name": "张三",
		},
	}

	// 执行测试 - 有效模板
	err := templateService.ValidateTemplate(ctx, validReq)
	assert.NoError(t, err)

	// 测试数据 - 无效的模板（语法错误）
	invalidReq := &service.ValidateTemplateRequest{
		TitleTemplate:   "您有新的审批待办：{{.approval_title", // 缺少闭合括号
		ContentTemplate: "申请人：{{.applicant_name}}，请及时处理。",
		Variables: map[string]any{
			"approval_title": "请假申请",
			"applicant_name": "张三",
		},
	}

	// 执行测试 - 无效模板
	err = templateService.ValidateTemplate(ctx, invalidReq)
	assert.Error(t, err)
}

func TestNotificationTemplateService_DeleteTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockTemplateRepo := mock_repository.NewMockNotificationTemplateRepository(ctrl)
	logger := testutil.CreateTestLogger()

	// 创建服务实例
	templateService := service.NewNotificationTemplateService(
		mockTemplateRepo,
		logger,
	)

	ctx := context.Background()

	// 测试数据
	templateID := "template_123"

	// 设置 Mock 期望
	mockTemplateRepo.EXPECT().
		Delete(ctx, templateID).
		Return(nil)

	// 执行测试
	err := templateService.Delete(ctx, templateID)

	// 验证结果
	assert.NoError(t, err)
}
