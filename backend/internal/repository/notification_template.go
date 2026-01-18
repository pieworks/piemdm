package repository

import (
	"context"
	"time"

	"piemdm/internal/model"

	"gorm.io/gorm"
)

//go:generate mockgen -source=notification_template_repository.go -destination=../../test/mocks/repository/notification_template.go

// NotificationTemplateRepository 通知模板仓储接口
type NotificationTemplateRepository interface {
	// 基础查询
	FindOne(ctx context.Context, id string) (*model.NotificationTemplate, error)
	FirstByCode(ctx context.Context, templateCode string) (*model.NotificationTemplate, error)
	GetByTypeAndNotification(ctx context.Context, templateType, notificationType string) (*model.NotificationTemplate, error)
	FindPage(ctx context.Context, req *ListNotificationTemplateRequest) ([]*model.NotificationTemplate, error)
	Count(ctx context.Context, req *ListNotificationTemplateRequest) (int64, error)
	GetActiveTemplates(ctx context.Context, templateType string) ([]*model.NotificationTemplate, error)

	// Base CRUD
	Create(ctx context.Context, template *model.NotificationTemplate) error
	Update(ctx context.Context, template *model.NotificationTemplate) error
	UpdateStatus(ctx context.Context, ids []string, status string) error
	Delete(ctx context.Context, id string) error
}

// ListNotificationTemplateRequest 通知模板列表请求
type ListNotificationTemplateRequest struct {
	TemplateType     string `json:"template_type"`     // 模板类型
	NotificationType string `json:"notification_type"` // 通知类型
	TemplateCode     string `json:"template_code"`     // 模板编码
	TemplateName     string `json:"template_name"`     // 模板名称
	Status           string `json:"status"`            // 状态
	Keyword          string `json:"keyword"`           // 关键词搜索
	Page             int    `json:"page"`              // 页码
	PageSize         int    `json:"page_size"`         // 页大小
}

// notificationTemplateRepository 通知模板仓储实现
type notificationTemplateRepository struct {
	db *gorm.DB
}

// NewNotificationTemplateRepository 创建通知模板仓储
func NewNotificationTemplateRepository(db *gorm.DB) NotificationTemplateRepository {
	return &notificationTemplateRepository{
		db: db,
	}
}

// Create 创建通知模板
func (r *notificationTemplateRepository) Create(ctx context.Context, template *model.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// Update 更新通知模板
func (r *notificationTemplateRepository) Update(ctx context.Context, template *model.NotificationTemplate) error {
	return r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).Where("id = ?", template.ID).Updates(template).Error
}

// UpdateStatus 批量更新通知模板状态
func (r *notificationTemplateRepository) UpdateStatus(ctx context.Context, ids []string, status string) error {
	return r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).Where("id IN ?", ids).Update("status", status).Error
}

// Delete 删除通知模板
func (r *notificationTemplateRepository) Delete(ctx context.Context, id string) error {
	// 使用 Updates 手动进行软删除，以便同时更新 Status 字段
	// 注意：这也将触发 BeforeUpdate 钩子，而不是 BeforeDelete
	return r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.StatusDeleted,
			"deleted_at": time.Now(),
		}).Error
}

// FindOne 根据ID获取通知模板
func (r *notificationTemplateRepository) FindOne(ctx context.Context, id string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FirstByCode 根据模板编码获取通知模板
func (r *notificationTemplateRepository) FirstByCode(ctx context.Context, templateCode string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	err := r.db.WithContext(ctx).Where("template_code = ? AND status = ?", templateCode, "Normal").First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByTypeAndNotification 根据模板类型和通知类型获取模板
func (r *notificationTemplateRepository) GetByTypeAndNotification(ctx context.Context, templateType, notificationType string) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	err := r.db.WithContext(ctx).Where("template_type = ? AND notification_type = ? AND status = ?",
		templateType, notificationType, "Normal").First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindPage 获取通知模板列表
func (r *notificationTemplateRepository) FindPage(ctx context.Context, req *ListNotificationTemplateRequest) ([]*model.NotificationTemplate, error) {
	var templates []*model.NotificationTemplate

	query := r.db.WithContext(ctx).Model(&model.NotificationTemplate{})

	// 添加过滤条件
	query = r.buildListQuery(query, req)

	// 分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// 排序
	query = query.Order("created_at DESC")

	err := query.Find(&templates).Error
	return templates, err
}

// Count 统计通知模板数量
func (r *notificationTemplateRepository) Count(ctx context.Context, req *ListNotificationTemplateRequest) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.NotificationTemplate{})

	// 添加过滤条件
	query = r.buildListQuery(query, req)

	err := query.Count(&count).Error
	return count, err
}

// GetActiveTemplates 获取激活的模板列表
func (r *notificationTemplateRepository) GetActiveTemplates(ctx context.Context, templateType string) ([]*model.NotificationTemplate, error) {
	var templates []*model.NotificationTemplate

	query := r.db.WithContext(ctx).Where("status = ?", "Normal")

	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
	}

	err := query.Order("notification_type, created_at DESC").Find(&templates).Error
	return templates, err
}

// buildListQuery 构建列表查询条件
func (r *notificationTemplateRepository) buildListQuery(query *gorm.DB, req *ListNotificationTemplateRequest) *gorm.DB {
	if req.TemplateType != "" {
		query = query.Where("template_type = ?", req.TemplateType)
	}

	if req.NotificationType != "" {
		query = query.Where("notification_type = ?", req.NotificationType)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.TemplateCode != "" {
		query = query.Where("template_code LIKE ?", "%"+req.TemplateCode+"%")
	}

	if req.TemplateName != "" {
		query = query.Where("template_name LIKE ?", "%"+req.TemplateName+"%")
	}

	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("template_name LIKE ? OR template_code LIKE ? OR description LIKE ?",
			keyword, keyword, keyword)
	}

	return query
}
