package repository_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupApprovalRepository(t *testing.T) (repository.ApprovalRepository, sqlmock.Sqlmock) {
	// 设置Gin为测试模式，避免调试信息输出
	gin.SetMode(gin.TestMode)

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

	rdb, _ := redismock.NewClientMock()
	repo := repository.NewRepository(db, rdb, nil)
	base := repository.NewBaseRepository(repo)
	approvalRepo := repository.NewApprovalRepository(repo, base)

	return approvalRepo, mock
}

func TestApprovalRepository_Create(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	now := time.Now()
	approval := &model.Approval{
		ID:              1,
		Code:            "APPROVAL_001",
		Title:           "测试审批申请",
		ApprovalDefCode: "TEST_DEF_001",
		EntityCode:      "test_entity",
		SerialNumber:    "AP20240101000001",
		CurrentTaskID:   "start",
		CurrentTaskName: "开始节点",
		FormData:        `{"field1":"value1"}`,
		FormSchema:      `{"schema":{}}`,
		Priority:        0,
		Urgency:         "Normal",
		Status:          model.ApprovalStatusPending,
		CreatedBy:       "admin",
		UpdatedBy:       "admin",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `approvals`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := approvalRepo.Create(c, approval)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_Update(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	now := time.Now()
	approval := &model.Approval{
		ID:              1,
		Code:            "APPROVAL_001",
		Title:           "测试审批申请-更新",
		ApprovalDefCode: "TEST_DEF_001",
		EntityCode:      "test_entity",
		SerialNumber:    "20240120001",
		CurrentTaskID:   "approval1",
		CurrentTaskName: "审批节点1",
		FormData:        `{"field1":"value1_updated"}`,
		FormSchema:      `{"schema":{}}`,
		Status:          model.ApprovalStatusPending,
		Priority:        1,
		Urgency:         "High",
		CreatedBy:       "user001",
		UpdatedBy:       "user001",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `approvals`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := approvalRepo.Update(c, approval)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_FirstById(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	id := uint(1)

	rows := sqlmock.NewRows([]string{
		"id", "code", "title", "approval_def_code", "entity_code", "entity_id",
		"serial_number", "current_task_id", "current_task_name", "form_data",
		"form_schema", "status", "priority",
		"urgency", "created_by", "updated_by",
		"created_by", "updated_by", "created_at", "updated_at",
	}).AddRow(
		1, "APPROVAL_001", "测试审批申请", "TEST_DEF_001", "test_entity", "entity_001",
		"20240120001", "start", "开始节点", `{"field1":"value1"}`,
		`{"schema":{}}`, "Pending", 0,
		"Normal", "user001", "user001",
		"user001", "user001", time.Now(), time.Now(),
	)
	mock.ExpectQuery("SELECT \\* FROM `approvals`").WillReturnRows(rows)

	approval, err := approvalRepo.FindOne(id)
	assert.NoError(t, err)
	assert.NotNil(t, approval)
	assert.Equal(t, "APPROVAL_001", approval.Code)
	assert.Equal(t, "测试审批申请", approval.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_FirstByCode(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	code := "APPROVAL_001"

	rows := sqlmock.NewRows([]string{
		"id", "code", "title", "approval_def_code", "entity_code",
		"serial_number", "current_task_id", "current_task_name", "form_data",
		"form_schema", "priority", "urgency", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		1, "APPROVAL_001", "测试审批申请", "TEST_DEF_001", "test_entity",
		"20240120001", "start", "开始节点", `{"field1":"value1"}`,
		`{"schema":{}}`, 0, "Normal", "", "Pending",
		"user001", "user001", time.Now(), time.Now(), nil,
	)
	mock.ExpectQuery("SELECT \\* FROM `approvals` WHERE \\(code = \\? AND deleted_at IS NULL\\) AND `approvals`.`deleted_at` IS NULL ORDER BY `approvals`.`id` LIMIT \\?").
		WithArgs(code, 1).
		WillReturnRows(rows)

	approval, err := approvalRepo.FirstByCode(code)
	assert.NoError(t, err)
	assert.NotNil(t, approval)
	assert.Equal(t, "APPROVAL_001", approval.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_FindByApplicantID(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	applicantID := "user001"

	rows := sqlmock.NewRows([]string{
		"id", "code", "title", "approval_def_code", "entity_code",
		"serial_number", "current_task_id", "current_task_name", "form_data",
		"form_schema", "priority", "urgency", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		1, "APPROVAL_001", "测试审批申请1", "TEST_DEF_001", "test_entity",
		"20240120001", "start", "开始节点", `{"field1":"value1"}`,
		`{"schema":{}}`, 0, "Normal", "", "Pending",
		"user001", "user001", time.Now(), time.Now(), nil,
	).AddRow(
		2, "APPROVAL_002", "测试审批申请2", "TEST_DEF_001", "test_entity",
		"20240120002", "approval1", "审批节点1", `{"field1":"value2"}`,
		`{"schema":{}}`, 0, "Normal", "", "Approved",
		"user001", "user001", time.Now(), time.Now(), nil,
	)
	mock.ExpectQuery("SELECT \\* FROM `approvals` WHERE \\(created_by = \\? AND deleted_at IS NULL\\) AND `approvals`.`deleted_at` IS NULL ORDER BY created_at DESC").WillReturnRows(rows)

	approvals, err := approvalRepo.FindByApplicantID(applicantID)
	assert.NoError(t, err)
	assert.NotNil(t, approvals)
	assert.Len(t, approvals, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_FindByStatus(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	status := model.ApprovalStatusPending

	rows := sqlmock.NewRows([]string{
		"id", "code", "title", "approval_def_code", "entity_code",
		"serial_number", "current_task_id", "current_task_name", "form_data",
		"form_schema", "priority", "urgency", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		1, "APPROVAL_001", "测试审批申请", "TEST_DEF_001", "test_entity",
		"20240120001", "start", "开始节点", `{"field1":"value1"}`,
		`{"schema":{}}`, 0, "Normal", "", "Pending",
		"user001", "user001", time.Now(), time.Now(), nil,
	)
	mock.ExpectQuery("SELECT \\* FROM `approvals` WHERE \\(status = \\? AND deleted_at IS NULL\\) AND `approvals`.`deleted_at` IS NULL").WillReturnRows(rows)

	approvals, err := approvalRepo.FindByStatus(status)
	assert.NoError(t, err)
	assert.NotNil(t, approvals)
	assert.Len(t, approvals, 1)
	assert.Equal(t, model.ApprovalStatusPending, approvals[0].Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_UpdateStatus(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	id := uint(1)
	status := model.ApprovalStatusApproved

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `approvals` SET `status`=\\?,`updated_at`=\\? WHERE id = \\? AND `approvals`.`deleted_at` IS NULL").
		WithArgs(status, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := approvalRepo.UpdateStatus(c, id, status)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_UpdateCurrentNode(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	id := uint(1)
	nodeID := "approval1"
	nodeName := "审批节点1"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `approvals` SET `current_task_id`=\\?,`current_task_name`=\\?,`updated_at`=\\? WHERE id = \\? AND `approvals`.`deleted_at` IS NULL").
		WithArgs(nodeID, nodeName, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := approvalRepo.UpdateCurrentNode(c, id, nodeID, nodeName)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_CountByStatus(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	status := model.ApprovalStatusPending

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `approvals` WHERE \\(status = \\? AND deleted_at IS NULL\\) AND `approvals`.`deleted_at` IS NULL").WillReturnRows(rows)

	count, err := approvalRepo.CountByStatus(status)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApprovalRepository_Delete(t *testing.T) {
	approvalRepo, mock := setupApprovalRepository(t)

	id := uint(1)

	// 先查询记录是否存在
	rows := sqlmock.NewRows([]string{
		"id", "code", "title", "approval_def_code", "entity_code",
		"serial_number", "current_task_id", "current_task_name", "form_data",
		"form_schema", "priority", "urgency", "description", "status",
		"created_by", "updated_by", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		1, "APPROVAL_001", "测试审批申请", "TEST_DEF_001", "test_entity",
		"20240120001", "start", "开始节点", `{"field1":"value1"}`,
		`{"schema":{}}`, 0, "Normal", "", "Pending",
		"user001", "user001", time.Now(), time.Now(), nil,
	)
	mock.ExpectQuery("SELECT \\* FROM `approvals` WHERE id = \\? AND `approvals`.`deleted_at` IS NULL ORDER BY `approvals`.`id` LIMIT \\?").
		WithArgs(id, 1).
		WillReturnRows(rows)

	// 模拟软删除操作 - GORM软删除会设置deleted_at字段，同时BeforeDelete钩子会更新状态为Deleted
	mock.ExpectBegin()
	// 首先执行BeforeDelete钩子中的更新操作
	mock.ExpectExec("UPDATE `approvals` SET `status`=\\?,`updated_at`=\\? WHERE id = \\? AND `approvals`.`deleted_at` IS NULL AND `id` = \\?").
		WithArgs("Deleted", sqlmock.AnyArg(), id, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// 然后执行GORM的软删除操作，设置deleted_at字段
	mock.ExpectExec("UPDATE `approvals` SET `deleted_at`=\\? WHERE `approvals`.`id` = \\? AND `approvals`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	deletedApproval, err := approvalRepo.Delete(c, id)
	assert.NoError(t, err)
	assert.NotNil(t, deletedApproval)
	assert.Equal(t, uint(1), deletedApproval.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
