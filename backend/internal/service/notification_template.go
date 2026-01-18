package service

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/helper/sid"
	"piemdm/pkg/log"
)

// NotificationTemplateService 通知模板服务接口
type NotificationTemplateService interface {
	// Base CRUD
	Create(ctx context.Context, req *CreateNotificationTemplateRequest) (*model.NotificationTemplate, error)
	Update(ctx context.Context, req *UpdateNotificationTemplateRequest) (*model.NotificationTemplate, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*model.NotificationTemplate, error)
	List(ctx context.Context, req *repository.ListNotificationTemplateRequest) (*PageResult[*model.NotificationTemplate], error)

	// 业务方法
	GetByCode(ctx context.Context, templateCode string) (*model.NotificationTemplate, error)
	GetByTypeAndNotification(ctx context.Context, templateType, notificationType string) (*model.NotificationTemplate, error)
	GetActiveTemplates(ctx context.Context, templateType string) ([]*model.NotificationTemplate, error)

	// 模板处理
	RenderTemplate(ctx context.Context, templateID string, variables map[string]any) (*RenderedTemplate, error)
	RenderTemplateByCode(ctx context.Context, templateCode string, variables map[string]any) (*RenderedTemplate, error)
	ValidateTemplate(ctx context.Context, req *ValidateTemplateRequest) error
	PreviewTemplate(ctx context.Context, req *PreviewTemplateRequest) (*RenderedTemplate, error)
	BatchUpdateStatus(ctx context.Context, ids []string, status string) error
}

// CreateNotificationTemplateRequest 创建通知模板请求
type CreateNotificationTemplateRequest struct {
	TemplateCode     string         `json:"TemplateCode" binding:"required,max=100"`
	TemplateName     string         `json:"TemplateName" binding:"required,max=200"`
	TemplateType     string         `json:"TemplateType" binding:"required,max=32"`
	NotificationType string         `json:"NotificationType" binding:"required,max=16"`
	TitleTemplate    string         `json:"TitleTemplate" binding:"required,max=500"`
	ContentTemplate  string         `json:"ContentTemplate" binding:"required"`
	Variables        map[string]any `json:"Variables"`
	Description      string         `json:"Description"`
	Status           string         `json:"Status"`
	CreatedBy        string         `json:"CreatedBy"`
}

// UpdateNotificationTemplateRequest 更新通知模板请求
type UpdateNotificationTemplateRequest struct {
	ID              string         `json:"ID" binding:"required"`
	TemplateName    string         `json:"TemplateName" binding:"required,max=200"`
	TitleTemplate   string         `json:"TitleTemplate" binding:"required,max=500"`
	ContentTemplate string         `json:"ContentTemplate" binding:"required"`
	Variables       map[string]any `json:"Variables"`
	Description     string         `json:"Description"`
	Status          string         `json:"Status"`
	UpdatedBy       string         `json:"UpdatedBy"`
}

// ValidateTemplateRequest 验证模板请求
type ValidateTemplateRequest struct {
	TitleTemplate   string         `json:"TitleTemplate" binding:"required"`
	ContentTemplate string         `json:"ContentTemplate" binding:"required"`
	Variables       map[string]any `json:"Variables"`
}

// PreviewTemplateRequest 预览模板请求
type PreviewTemplateRequest struct {
	TitleTemplate   string         `json:"TitleTemplate" binding:"required"`
	ContentTemplate string         `json:"ContentTemplate" binding:"required"`
	Variables       map[string]any `json:"Variables"`
}

// RenderedTemplate 渲染后的模板
type RenderedTemplate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// notificationTemplateService 通知模板服务实现
type notificationTemplateService struct {
	templateRepo repository.NotificationTemplateRepository
	logger       *log.Logger
}

// NewNotificationTemplateService 创建通知模板服务
func NewNotificationTemplateService(
	templateRepo repository.NotificationTemplateRepository,
	logger *log.Logger,
) NotificationTemplateService {
	return &notificationTemplateService{
		templateRepo: templateRepo,
		logger:       logger,
	}
}

// Create 创建通知模板
func (s *notificationTemplateService) Create(ctx context.Context, req *CreateNotificationTemplateRequest) (*model.NotificationTemplate, error) {
	// 验证模板类型和通知类型
	if !model.IsValidTemplateType(req.TemplateType) {
		return nil, fmt.Errorf("无效的模板类型: %s", req.TemplateType)
	}

	if !model.IsValidNotificationType(req.NotificationType) {
		return nil, fmt.Errorf("无效的通知类型: %s", req.NotificationType)
	}

	// 验证模板语法
	if err := s.validateTemplateSyntax(req.TitleTemplate, req.ContentTemplate); err != nil {
		return nil, fmt.Errorf("模板语法错误: %v", err)
	}

	// 创建模板对象
	sidGen := sid.NewSid()
	templateID, err := sidGen.GenString()
	if err != nil {
		return nil, fmt.Errorf("生成模板ID失败: %v", err)
	}

	template := &model.NotificationTemplate{
		ID:               templateID,
		TemplateCode:     req.TemplateCode,
		TemplateName:     req.TemplateName,
		TemplateType:     req.TemplateType,
		NotificationType: req.NotificationType,
		TitleTemplate:    req.TitleTemplate,
		ContentTemplate:  req.ContentTemplate,
		Description:      req.Description,
		Status:           req.Status,
		CreatedBy:        req.CreatedBy,
	}

	if template.Status == "" {
		template.Status = "Normal"
	}

	// 设置变量
	if err := template.SetVariables(req.Variables); err != nil {
		return nil, fmt.Errorf("设置变量失败: %v", err)
	}

	// 保存到数据库
	if err := s.templateRepo.Create(ctx, template); err != nil {
		if s.logger != nil {
			s.logger.Error("创建通知模板失败", "error", err)
		}
		return nil, fmt.Errorf("创建通知模板失败: %v", err)
	}

	return template, nil
}

// BatchUpdateStatus 批量更新通知模板状态
func (s *notificationTemplateService) BatchUpdateStatus(ctx context.Context, ids []string, status string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids 不能为空")
	}
	if status == "" {
		return fmt.Errorf("状态不能为空")
	}
	return s.templateRepo.UpdateStatus(ctx, ids, status)
}

// Update 更新通知模板
func (s *notificationTemplateService) Update(ctx context.Context, req *UpdateNotificationTemplateRequest) (*model.NotificationTemplate, error) {
	// 获取现有模板
	template, err := s.templateRepo.FindOne(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("获取通知模板失败: %v", err)
	}

	// 验证模板语法
	if err := s.validateTemplateSyntax(req.TitleTemplate, req.ContentTemplate); err != nil {
		return nil, fmt.Errorf("模板语法错误: %v", err)
	}

	// 更新字段
	template.TemplateName = req.TemplateName
	template.TitleTemplate = req.TitleTemplate
	template.ContentTemplate = req.ContentTemplate
	template.Description = req.Description
	if req.Status != "" {
		template.Status = req.Status
	}
	template.UpdatedBy = req.UpdatedBy

	// 设置变量
	if err := template.SetVariables(req.Variables); err != nil {
		return nil, fmt.Errorf("设置变量失败: %v", err)
	}

	// 保存到数据库
	if err := s.templateRepo.Update(ctx, template); err != nil {
		if s.logger != nil {
			s.logger.Error("更新通知模板失败", "error", err)
		}
		return nil, fmt.Errorf("更新通知模板失败: %v", err)
	}

	return template, nil
}

// Delete 删除通知模板
func (s *notificationTemplateService) Delete(ctx context.Context, id string) error {
	if err := s.templateRepo.Delete(ctx, id); err != nil {
		if s.logger != nil {
			s.logger.Error("删除通知模板失败", "error", err)
		}
		return fmt.Errorf("删除通知模板失败: %v", err)
	}

	if s.logger != nil {
		s.logger.Error("删除通知模板成功", "template_id", id)
	}

	return nil
}

// Get 获取通知模板
func (s *notificationTemplateService) Get(ctx context.Context, id string) (*model.NotificationTemplate, error) {
	return s.templateRepo.FindOne(ctx, id)
}

// GetByCode 根据编码获取通知模板
func (s *notificationTemplateService) GetByCode(ctx context.Context, templateCode string) (*model.NotificationTemplate, error) {
	return s.templateRepo.FirstByCode(ctx, templateCode)
}

// GetByTypeAndNotification 根据模板类型和通知类型获取模板
func (s *notificationTemplateService) GetByTypeAndNotification(ctx context.Context, templateType, notificationType string) (*model.NotificationTemplate, error) {
	return s.templateRepo.GetByTypeAndNotification(ctx, templateType, notificationType)
}

// List 获取通知模板列表
func (s *notificationTemplateService) List(ctx context.Context, req *repository.ListNotificationTemplateRequest) (*PageResult[*model.NotificationTemplate], error) {
	templates, err := s.templateRepo.FindPage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取通知模板列表失败: %v", err)
	}

	total, err := s.templateRepo.Count(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("统计通知模板数量失败: %v", err)
	}

	return &PageResult[*model.NotificationTemplate]{
		Data:     templates,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetActiveTemplates 获取激活的模板列表
func (s *notificationTemplateService) GetActiveTemplates(ctx context.Context, templateType string) ([]*model.NotificationTemplate, error) {
	return s.templateRepo.GetActiveTemplates(ctx, templateType)
}

// RenderTemplate 渲染模板
func (s *notificationTemplateService) RenderTemplate(ctx context.Context, templateID string, variables map[string]any) (*RenderedTemplate, error) {
	// 获取模板
	tmpl, err := s.templateRepo.FindOne(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取通知模板失败: %v", err)
	}

	if tmpl.Status != "Normal" {
		return nil, fmt.Errorf("模板未激活")
	}

	return s.renderTemplate(tmpl.TitleTemplate, tmpl.ContentTemplate, variables)
}

// RenderTemplateByCode 根据编码渲染模板
func (s *notificationTemplateService) RenderTemplateByCode(ctx context.Context, templateCode string, variables map[string]any) (*RenderedTemplate, error) {
	// 获取模板
	tmpl, err := s.templateRepo.FirstByCode(ctx, templateCode)
	if err != nil {
		return nil, fmt.Errorf("获取通知模板失败: %v", err)
	}

	return s.renderTemplate(tmpl.TitleTemplate, tmpl.ContentTemplate, variables)
}

// ValidateTemplate 验证模板
func (s *notificationTemplateService) ValidateTemplate(ctx context.Context, req *ValidateTemplateRequest) error {
	return s.validateTemplateSyntax(req.TitleTemplate, req.ContentTemplate)
}

// PreviewTemplate 预览模板
func (s *notificationTemplateService) PreviewTemplate(ctx context.Context, req *PreviewTemplateRequest) (*RenderedTemplate, error) {
	// 验证模板语法
	if err := s.validateTemplateSyntax(req.TitleTemplate, req.ContentTemplate); err != nil {
		return nil, fmt.Errorf("模板语法错误: %v", err)
	}

	return s.renderTemplate(req.TitleTemplate, req.ContentTemplate, req.Variables)
}

// validateTemplateSyntax 验证模板语法
func (s *notificationTemplateService) validateTemplateSyntax(titleTemplate, contentTemplate string) error {
	// 验证标题模板
	if _, err := template.New("title").Parse(titleTemplate); err != nil {
		return fmt.Errorf("标题模板语法错误: %v", err)
	}

	// 验证内容模板
	if _, err := template.New("content").Parse(contentTemplate); err != nil {
		return fmt.Errorf("内容模板语法错误: %v", err)
	}

	return nil
}

// renderTemplate 渲染模板
func (s *notificationTemplateService) renderTemplate(titleTemplate, contentTemplate string, variables map[string]any) (*RenderedTemplate, error) {
	// 如果变量为空，使用空map
	if variables == nil {
		variables = make(map[string]any)
	}

	// 渲染标题
	titleTmpl, err := template.New("title").Parse(titleTemplate)
	if err != nil {
		return nil, fmt.Errorf("解析标题模板失败: %v", err)
	}

	var titleBuf bytes.Buffer
	if err := titleTmpl.Execute(&titleBuf, variables); err != nil {
		return nil, fmt.Errorf("渲染标题模板失败: %v", err)
	}

	// 渲染内容
	contentTmpl, err := template.New("content").Parse(contentTemplate)
	if err != nil {
		return nil, fmt.Errorf("解析内容模板失败: %v", err)
	}

	var contentBuf bytes.Buffer
	if err := contentTmpl.Execute(&contentBuf, variables); err != nil {
		return nil, fmt.Errorf("渲染内容模板失败: %v", err)
	}

	return &RenderedTemplate{
		Title:   strings.TrimSpace(titleBuf.String()),
		Content: strings.TrimSpace(contentBuf.String()),
	}, nil
}

// extractVariables 从模板中提取变量
func (s *notificationTemplateService) extractVariables(templateContent string) []string {
	// 使用正则表达式提取 {{.变量名}} 格式的变量
	re := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)
	matches := re.FindAllStringSubmatch(templateContent, -1)

	var variables []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			variable := match[1]
			if !seen[variable] {
				variables = append(variables, variable)
				seen[variable] = true
			}
		}
	}

	return variables
}
