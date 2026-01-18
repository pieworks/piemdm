package service_test

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/sid"
	"piemdm/pkg/jwt"
	"piemdm/pkg/log"
	mock_repository "piemdm/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupApprovalTaskService(t *testing.T) (service.ApprovalTaskService, *mock_repository.MockApprovalTaskRepository, *mock_repository.MockApprovalRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	logger := &log.Logger{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	mockTaskRepo := mock_repository.NewMockApprovalTaskRepository(ctrl)
	mockApprovalRepo := mock_repository.NewMockApprovalRepository(ctrl)

	baseService := service.NewService(logger, &sid.Sid{}, &jwt.JWT{})

	taskService := service.NewApprovalTaskService(baseService, mockTaskRepo, mockApprovalRepo)
	return taskService, mockTaskRepo, mockApprovalRepo, ctrl
}

func TestApprovalTaskService_Create(t *testing.T) {
	taskService, mockTaskRepo, _, ctrl := setupApprovalTaskService(t)
	defer ctrl.Finish()

	now := time.Now()
	task := &model.ApprovalTask{
		ID:           1,
		TaskCode:     "TASK_001",
		ApprovalCode: "APPROVAL_001",
		NodeCode:     "node_001",
		NodeName:     "审批节点1",
		AssigneeID:   "user001",
		AssigneeName: "张三",
		Status:       model.TaskStatusPending,
		// Priority:     0,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	mockTaskRepo.EXPECT().Create(task).Return(nil)

	err := taskService.Create(task)
	assert.NoError(t, err)
}

func TestApprovalTaskService_Get(t *testing.T) {
	taskService, mockTaskRepo, _, ctrl := setupApprovalTaskService(t)
	defer ctrl.Finish()

	id := uint(1)
	now := time.Now()
	expectedTask := &model.ApprovalTask{
		ID:           id,
		TaskCode:     "TASK_001",
		ApprovalCode: "APPROVAL_001",
		NodeCode:     "node_001",
		NodeName:     "审批节点1",
		AssigneeID:   "user001",
		AssigneeName: "张三",
		Status:       model.TaskStatusPending,
		// Priority:     0,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	mockTaskRepo.EXPECT().FindOne(id).Return(expectedTask, nil)

	task, err := taskService.Get(id)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "TASK_001", task.TaskCode)
	assert.Equal(t, "审批节点1", task.NodeName)
}

func TestApprovalTaskService_GetByAssignee(t *testing.T) {
	taskService, mockTaskRepo, _, ctrl := setupApprovalTaskService(t)
	defer ctrl.Finish()

	assigneeID := "user001"
	now := time.Now()
	expectedTasks := []*model.ApprovalTask{
		{
			ID:           uint(1),
			TaskCode:     "TASK_001",
			ApprovalCode: "APPROVAL_001",
			NodeCode:     "node_001",
			NodeName:     "审批节点1",
			AssigneeID:   "user001",
			AssigneeName: "张三",
			Status:       model.TaskStatusPending,
			// Priority:     0,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			ID:           uint(2),
			TaskCode:     "TASK_002",
			ApprovalCode: "APPROVAL_002",
			NodeCode:     "node_001",
			NodeName:     "审批节点1",
			AssigneeID:   "user001",
			AssigneeName: "张三",
			Status:       model.TaskStatusPending,
			//  Priority:     0,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	mockTaskRepo.EXPECT().FindByAssigneeID(assigneeID).Return(expectedTasks, nil)

	tasks, err := taskService.GetByAssignee(assigneeID)
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "user001", tasks[0].AssigneeID)
}

func TestApprovalTaskService_Delete(t *testing.T) {
	taskService, mockTaskRepo, _, ctrl := setupApprovalTaskService(t)
	defer ctrl.Finish()

	id := uint(1)
	task := &model.ApprovalTask{
		ID:       id,
		TaskCode: "TASK_001",
		Status:   model.TaskStatusPending,
	}

	// 第一次调用 FindOne 获取删除前的记录
	mockTaskRepo.EXPECT().FindOne(id).Return(task, nil)
	// 调用 BatchDelete 执行实际删除
	mockTaskRepo.EXPECT().BatchDelete([]uint{id}).Return(nil)

	deletedTask, err := taskService.Delete(id)
	assert.NoError(t, err)
	assert.NotNil(t, deletedTask)
}
