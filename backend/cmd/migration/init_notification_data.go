package main

import (
	"piemdm/internal/model"
	"piemdm/pkg/log"

	"gorm.io/gorm"
)

type InitNotificationData struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewInitNotificationData(db *gorm.DB, logger *log.Logger) *InitNotificationData {
	return &InitNotificationData{
		db:     db,
		logger: logger,
	}
}

func (m *InitNotificationData) Run() error {
	m.logger.Info("Starting notification template data initialization...")

	templates := []model.NotificationTemplate{
		{
			TemplateCode:     "test_flow_template",
			TemplateName:     "Integration Test Full-Link Notification Template",
			TemplateType:     "test",
			NotificationType: "email",
			TitleTemplate:    "Test Notification: {{.title}}",
			ContentTemplate:  "This is a test notification from {{.sender}}, time: {{.time}}",
			Status:           "Normal",
			Description:      "Used for integration test full-link process verification",
			CreatedBy:        "system",
		},
		{
			TemplateCode:     "approval_pending_email",
			TemplateName:     "Approval Pending Email Template",
			TemplateType:     model.TemplateTypeApprovalPending,
			NotificationType: model.NotificationTypeEmail,
			TitleTemplate:    "You have a new approval pending: {{.approval_title}}",
			ContentTemplate:  "Applicant: {{.applicant_name}}, please process in time.",
			Variables:        `{"approval_title": "string", "applicant_name": "string"}`,
			Status:           "Normal",
			CreatedBy:        "system",
		},
	}

	for _, tmpl := range templates {
		var count int64
		m.db.Model(&model.NotificationTemplate{}).Where("template_code = ?", tmpl.TemplateCode).Count(&count)
		if count == 0 {
			if err := m.db.Create(&tmpl).Error; err != nil {
				m.logger.Error("Failed to create notification template: " + tmpl.TemplateCode + ", err: " + err.Error())
				return err
			}
			m.logger.Info("Successfully initialized notification template: " + tmpl.TemplateCode)
		} else {
			m.logger.Info("Notification template already exists, skipping: " + tmpl.TemplateCode)
		}
	}

	return nil
}
