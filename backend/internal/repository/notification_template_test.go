package repository_test

import (
	"context"
	"testing"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupNotificationTemplateRepository(t *testing.T) (repository.NotificationTemplateRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	templateRepo := repository.NewNotificationTemplateRepository(db)

	return templateRepo, mock
}

func TestNotificationTemplateRepository_Create(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	now := time.Now()
	template := &model.NotificationTemplate{
		ID:               "template_001",
		TemplateCode:     "approval_pending_email",
		TemplateName:     "审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时处理。",
		Variables:        `{"approval_title":"审批标题","applicant_name":"申请人姓名"}`,
		Description:      "审批待办邮件通知模板",
		Status:           "Normal",
		CreatedBy:        "admin",
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `notification_templates`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := templateRepo.Create(context.Background(), template)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_GetByID(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	templateID := "template_001"

	rows := sqlmock.NewRows([]string{
		"id", "template_code", "template_name", "template_type", "notification_type",
		"title_template", "content_template", "variables", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at",
	}).AddRow(
		"template_001", "approval_pending_email", "审批待办邮件模板", "approval_pending", "email",
		"您有新的审批待办：{{.approval_title}}", "申请人：{{.applicant_name}}，请及时处理。",
		`{"approval_title":"审批标题","applicant_name":"申请人姓名"}`, "审批待办邮件通知模板", "Normal",
		"admin", "admin", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_templates` WHERE id = \\? AND `notification_templates`.`deleted_at` IS NULL ORDER BY `notification_templates`.`id` LIMIT \\?").
		WithArgs(templateID, 1).
		WillReturnRows(rows)

	template, err := templateRepo.FindOne(context.Background(), templateID)
	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, "template_001", template.ID)
	assert.Equal(t, "approval_pending_email", template.TemplateCode)
	assert.Equal(t, "审批待办邮件模板", template.TemplateName)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_GetByCode(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	templateCode := "approval_pending_email"

	rows := sqlmock.NewRows([]string{
		"id", "template_code", "template_name", "template_type", "notification_type",
		"title_template", "content_template", "variables", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at",
	}).AddRow(
		"template_001", "approval_pending_email", "审批待办邮件模板", "approval_pending", "email",
		"您有新的审批待办：{{.approval_title}}", "申请人：{{.applicant_name}}，请及时处理。",
		`{"approval_title":"审批标题","applicant_name":"申请人姓名"}`, "审批待办邮件通知模板", "Normal",
		"admin", "admin", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_templates` WHERE \\(template_code = \\? AND status = \\?\\) AND `notification_templates`.`deleted_at` IS NULL ORDER BY `notification_templates`.`id` LIMIT \\?").
		WithArgs(templateCode, "Normal", 1).
		WillReturnRows(rows)

	template, err := templateRepo.FirstByCode(context.Background(), templateCode)
	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, "approval_pending_email", template.TemplateCode)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_GetByTypeAndNotification(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	templateType := model.TemplateTypeApprovalPending
	notificationType := model.NotificationTypeEmail

	rows := sqlmock.NewRows([]string{
		"id", "template_code", "template_name", "template_type", "notification_type",
		"title_template", "content_template", "variables", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at",
	}).AddRow(
		"template_001", "approval_pending_email", "审批待办邮件模板", "approval_pending", "email",
		"您有新的审批待办：{{.approval_title}}", "申请人：{{.applicant_name}}，请及时处理。",
		`{"approval_title":"审批标题","applicant_name":"申请人姓名"}`, "审批待办邮件通知模板", "Normal",
		"admin", "admin", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_templates` WHERE \\(template_type = \\? AND notification_type = \\? AND status = \\?\\) AND `notification_templates`.`deleted_at` IS NULL ORDER BY `notification_templates`.`id` LIMIT \\?").
		WithArgs(templateType, notificationType, "Normal", 1).
		WillReturnRows(rows)

	template, err := templateRepo.GetByTypeAndNotification(context.Background(), templateType, notificationType)
	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, templateType, template.TemplateType)
	assert.Equal(t, notificationType, template.NotificationType)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_Update(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	now := time.Now()
	template := &model.NotificationTemplate{
		ID:               "template_001",
		TemplateCode:     "approval_pending_email",
		TemplateName:     "更新后的审批待办邮件模板",
		TemplateType:     model.TemplateTypeApprovalPending,
		NotificationType: model.NotificationTypeEmail,
		TitleTemplate:    "您有新的审批待办：{{.approval_title}}",
		ContentTemplate:  "申请人：{{.applicant_name}}，请及时查看详情。",
		Variables:        `{"approval_title":"审批标题","applicant_name":"申请人姓名"}`,
		Description:      "更新后的审批待办邮件通知模板",
		Status:           "Normal",
		UpdatedBy:        "admin",
		UpdatedAt:        &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `notification_templates`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := templateRepo.Update(context.Background(), template)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_Delete(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	templateID := "template_001"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `notification_templates` SET `deleted_at`=\\?,`status`=\\?,`updated_at`=\\? WHERE id = \\? AND `notification_templates`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), templateID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := templateRepo.Delete(context.Background(), templateID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_List(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	req := &repository.ListNotificationTemplateRequest{
		Page:     1,
		PageSize: 10,
	}

	rows := sqlmock.NewRows([]string{
		"id", "template_code", "template_name", "template_type", "notification_type",
		"title_template", "content_template", "variables", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at",
	}).AddRow(
		"template_001", "approval_pending_email", "审批待办邮件模板", "approval_pending", "email",
		"您有新的审批待办：{{.approval_title}}", "申请人：{{.applicant_name}}，请及时处理。",
		`{"approval_title":"审批标题","applicant_name":"申请人姓名"}`, "审批待办邮件通知模板", "Normal",
		"admin", "admin", time.Now(), time.Now(),
	).AddRow(
		"template_002", "approval_approved_email", "审批通过邮件模板", "approval_approved", "email",
		"您的审批已通过：{{.approval_title}}", "申请人：{{.applicant_name}}，审批已通过。",
		`{"approval_title":"审批标题","applicant_name":"申请人姓名"}`, "审批通过邮件通知模板", "Normal",
		"admin", "admin", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_templates` WHERE `notification_templates`.`deleted_at` IS NULL ORDER BY created_at DESC LIMIT \\?").
		WithArgs(10).
		WillReturnRows(rows)

	templates, err := templateRepo.FindPage(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, templates)
	assert.Equal(t, 2, len(templates))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_Count(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	req := &repository.ListNotificationTemplateRequest{
		TemplateType: model.TemplateTypeApprovalPending,
	}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `notification_templates` WHERE template_type = \\? AND `notification_templates`.`deleted_at` IS NULL").
		WithArgs("approval_pending").
		WillReturnRows(rows)

	count, err := templateRepo.Count(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationTemplateRepository_GetActiveTemplates(t *testing.T) {
	templateRepo, mock := setupNotificationTemplateRepository(t)

	templateType := model.TemplateTypeApprovalPending

	rows := sqlmock.NewRows([]string{
		"id", "template_code", "template_name", "template_type", "notification_type",
		"title_template", "content_template", "variables", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at",
	}).AddRow(
		"template_001", "approval_pending_email", "审批待办邮件模板", "approval_pending", "email",
		"您有新的审批待办：{{.approval_title}}", "申请人：{{.applicant_name}}，请及时处理。",
		`{"approval_title":"审批标题","applicant_name":"申请人姓名"}`, "审批待办邮件通知模板", "Normal",
		"admin", "admin", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_templates` WHERE status = \\? AND template_type = \\? AND `notification_templates`.`deleted_at` IS NULL ORDER BY notification_type, created_at DESC").
		WithArgs("Normal", templateType).
		WillReturnRows(rows)

	templates, err := templateRepo.GetActiveTemplates(context.Background(), templateType)
	assert.NoError(t, err)
	assert.NotNil(t, templates)
	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "Normal", templates[0].Status)
	assert.Equal(t, templateType, templates[0].TemplateType)

	assert.NoError(t, mock.ExpectationsWereMet())
}
