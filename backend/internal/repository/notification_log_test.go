package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"piemdm/internal/model"
	"piemdm/internal/repository"
)

func setupNotificationLogRepository(t *testing.T) (repository.NotificationLogRepository, sqlmock.Sqlmock) {
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

	logRepo := repository.NewNotificationLogRepository(db)

	return logRepo, mock
}

func TestNotificationLogRepository_Create(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	now := time.Now()
	log := &model.NotificationLog{
		ID:               1,
		ApprovalID:       "approval_001",
		TaskID:           "task_001",
		RecipientID:      "user_001",
		RecipientType:    model.RecipientTypeUser,
		NotificationType: model.NotificationTypeEmail,
		TemplateID:       "template_001",
		Title:            "审批待办通知",
		Content:          "您有新的审批待办，请及时处理。",
		Status:           model.NotificationStatusPending,
		RetryCount:       0,
		MaxRetryCount:    5,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `notification_logs`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := logRepo.Create(context.Background(), log)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_GetByID(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	logID := uint(1)

	rows := sqlmock.NewRows([]string{
		"id", "approval_id", "task_id", "recipient_id", "recipient_type",
		"notification_type", "template_id", "title", "content", "status",
		"sent_at", "error_message", "retry_count", "max_retry_count", "next_retry_at",
		"created_at", "updated_at",
	}).AddRow(
		1, "approval_001", "task_001", "user_001", "user",
		"email", "template_001", "审批待办通知", "您有新的审批待办，请及时处理。", "pending",
		nil, "", 0, 5, nil,
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_logs` WHERE id = \\? ORDER BY `notification_logs`.`id` LIMIT \\?").
		WithArgs(logID, 1).
		WillReturnRows(rows)

	log, err := logRepo.FindOne(context.Background(), logID)
	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, uint(1), log.ID)
	assert.Equal(t, "approval_001", log.ApprovalID)
	assert.Equal(t, "审批待办通知", log.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_Update(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	now := time.Now()
	log := &model.NotificationLog{
		ID:               1,
		ApprovalID:       "approval_001",
		RecipientID:      "user_001",
		RecipientType:    model.RecipientTypeUser,
		NotificationType: model.NotificationTypeEmail,
		Title:            "审批待办通知",
		Content:          "您有新的审批待办，请及时处理。",
		Status:           model.NotificationStatusSent,
		SendTime:         &now,
		UpdatedAt:        &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `notification_logs`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := logRepo.Update(context.Background(), log)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// 删除测试 - Repository接口没有Delete方法，移除此测试

func TestNotificationLogRepository_List(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	req := &repository.ListNotificationLogRequest{
		Page:     1,
		PageSize: 10,
	}

	rows := sqlmock.NewRows([]string{
		"id", "approval_id", "task_id", "recipient_id", "recipient_type",
		"notification_type", "template_id", "title", "content", "status",
		"sent_at", "error_message", "retry_count", "max_retry_count", "next_retry_at",
		"created_at", "updated_at",
	}).AddRow(
		1, "approval_001", "task_001", "user_001", "user",
		"email", "template_001", "审批待办通知1", "您有新的审批待办，请及时处理。", "sent",
		time.Now(), "", 0, 5, nil,
		time.Now(), time.Now(),
	).AddRow(
		2, "approval_002", "task_002", "user_002", "user",
		"sms", "template_002", "审批待办通知2", "您有新的审批待办，请及时处理。", "pending",
		nil, "", 0, 5, nil,
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_logs` ORDER BY created_at DESC LIMIT \\?").
		WithArgs(10).
		WillReturnRows(rows)

	logs, err := logRepo.FindPage(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.Equal(t, 2, len(logs))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_Count(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	req := &repository.ListNotificationLogRequest{
		ApprovalID: "approval_001",
	}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `notification_logs` WHERE approval_id = \\?").
		WithArgs("approval_001").
		WillReturnRows(rows)

	count, err := logRepo.Count(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_GetPendingLogs(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	limit := 10

	rows := sqlmock.NewRows([]string{
		"id", "approval_id", "task_id", "recipient_id", "recipient_type",
		"notification_type", "template_id", "title", "content", "status",
		"sent_at", "error_message", "retry_count", "max_retry_count", "next_retry_at",
		"created_at", "updated_at",
	}).AddRow(
		1, "approval_001", "task_001", "user_001", "user",
		"email", "template_001", "待发送通知", "您有新的审批待办，请及时处理。", "pending",
		nil, "", 0, 5, nil,
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_logs` WHERE status = \\? ORDER BY created_at ASC LIMIT \\?").
		WithArgs("pending", 10).
		WillReturnRows(rows)

	logs, err := logRepo.GetPendingLogs(context.Background(), limit)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.Equal(t, 1, len(logs))
	assert.Equal(t, model.NotificationStatusPending, logs[0].Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_GetRetryLogs(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	limit := 10

	rows := sqlmock.NewRows([]string{
		"id", "approval_id", "task_id", "recipient_id", "recipient_type",
		"notification_type", "template_id", "title", "content", "status",
		"sent_at", "error_message", "retry_count", "max_retry_count", "next_retry_at",
		"created_at", "updated_at",
	}).AddRow(
		1, "approval_001", "task_001", "user_001", "user",
		"email", "template_001", "重试通知", "您有新的审批待办，请及时处理。", "failed",
		nil, "发送失败", 2, 5, time.Now().Add(-1*time.Hour),
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_logs` WHERE status = \\? AND next_retry_time <= \\? AND retry_count < max_retry_count ORDER BY next_retry_time ASC LIMIT \\?").
		WithArgs("retry", sqlmock.AnyArg(), 10).
		WillReturnRows(rows)

	logs, err := logRepo.GetRetryLogs(context.Background(), limit)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.Equal(t, 1, len(logs))
	assert.Equal(t, model.NotificationStatusFailed, logs[0].Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_GetByApprovalID(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	approvalID := "approval_001"

	rows := sqlmock.NewRows([]string{
		"id", "approval_id", "task_id", "recipient_id", "recipient_type",
		"notification_type", "template_id", "title", "content", "status",
		"sent_at", "error_message", "retry_count", "max_retry_count", "next_retry_at",
		"created_at", "updated_at",
	}).AddRow(
		1, "approval_001", "task_001", "user_001", "user",
		"email", "template_001", "审批通知1", "您有新的审批待办，请及时处理。", "sent",
		time.Now(), "", 0, 5, nil,
		time.Now(), time.Now(),
	).AddRow(
		2, "approval_001", "task_002", "user_002", "user",
		"sms", "template_002", "审批通知2", "您有新的审批待办，请及时处理。", "pending",
		nil, "", 0, 5, nil,
		time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `notification_logs` WHERE approval_id = \\? ORDER BY created_at DESC").
		WithArgs(approvalID).
		WillReturnRows(rows)

	logs, err := logRepo.GetByApprovalID(context.Background(), approvalID)
	assert.NoError(t, err)
	assert.NotNil(t, logs)
	assert.Equal(t, 2, len(logs))
	for _, log := range logs {
		assert.Equal(t, approvalID, log.ApprovalID)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_BatchUpdateStatus(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	logIDs := []uint{1, 2, 3}
	status := model.NotificationStatusSent

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `notification_logs`").WillReturnResult(sqlmock.NewResult(3, 3))
	mock.ExpectCommit()

	err := logRepo.BatchUpdateStatus(context.Background(), logIDs, status)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationLogRepository_DeleteExpiredLogs(t *testing.T) {
	logRepo, mock := setupNotificationLogRepository(t)

	expiredBefore := time.Now().Add(-7 * 24 * time.Hour) // 7天前

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `notification_logs` WHERE created_at < \\? AND status IN \\(\\?,\\?\\)").
		WithArgs(expiredBefore, "sent", "expired").
		WillReturnResult(sqlmock.NewResult(5, 5))
	mock.ExpectCommit()

	err := logRepo.DeleteExpiredLogs(context.Background(), expiredBefore)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
