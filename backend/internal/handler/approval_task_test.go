package handler_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"piemdm/internal/handler"
	"piemdm/internal/model"
	"piemdm/pkg/log"
	mock_service "piemdm/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createTaskTestHandler() *handler.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	testLogger := &log.Logger{
		Logger: logger,
	}
	return handler.NewHandler(testLogger)
}

func TestApprovalTaskHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalTaskService := mock_service.NewMockApprovalTaskService(ctrl)

	expectedList := []*model.ApprovalTask{
		{
			ID:           1,
			ApprovalCode: "APPROVAL_001",
			NodeCode:     "APPROVAL_001",
			Status:       model.TaskStatusPending,
		},
	}

	mockApprovalTaskService.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedList, nil)

	// 传递 nil 作为 ApprovalService,因为这个测试不需要它
	approvalTaskHandler := handler.NewApprovalTaskHandler(createTaskTestHandler(), mockApprovalTaskService, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/approval-tasks", approvalTaskHandler.List)

	req, _ := http.NewRequest("GET", "/approval-tasks?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestApprovalTaskHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalTaskService := mock_service.NewMockApprovalTaskService(ctrl)

	expectedTask := &model.ApprovalTask{
		ID:           1,
		ApprovalCode: "APPROVAL_001",
		NodeCode:     "APPROVAL_001",
		Status:       model.TaskStatusPending,
	}

	mockApprovalTaskService.EXPECT().Get(uint(1)).Return(expectedTask, nil)

	approvalTaskHandler := handler.NewApprovalTaskHandler(createTaskTestHandler(), mockApprovalTaskService, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/approval-tasks/:id", approvalTaskHandler.Get)

	req, _ := http.NewRequest("GET", "/approval-tasks/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// TestApprovalTaskHandler_RemindTask - 方法未实现，跳过测试
// func TestApprovalTaskHandler_RemindTask(t *testing.T) {
//	 // RemindTask 方法在接口中被注释掉，暂未实现
// }

// TestApprovalTaskHandler_BatchRemindTasks - 方法未实现，跳过测试
// func TestApprovalTaskHandler_BatchRemindTasks(t *testing.T) {
//	 // BatchRemindTasks 方法在接口中被注释掉，暂未实现
// }

// TestApprovalTaskHandler_GetTasksByAssignee - 方法未实现，跳过测试
// func TestApprovalTaskHandler_GetTasksByAssignee(t *testing.T) {
//	 // GetTasksByAssignee 方法在接口中被注释掉，暂未实现
// }

// TestApprovalTaskHandler_GetPendingTasksByAssignee - 方法未实现，跳过测试
// func TestApprovalTaskHandler_GetPendingTasksByAssignee(t *testing.T) {
//	 // GetPendingTasksByAssignee 方法在接口中被注释掉，暂未实现
// }

// TestApprovalTaskHandler_GetTasksByApproval - 方法未实现，跳过测试
// func TestApprovalTaskHandler_GetTasksByApproval(t *testing.T) {
//	 // GetTasksByApproval 方法在接口中被注释掉，暂未实现
// }

// TestApprovalTaskHandler_ApproveTaskWithoutTaskIdInBody - 方法未实现，跳过测试
// func TestApprovalTaskHandler_ApproveTaskWithoutTaskIdInBody(t *testing.T) {
//	 // ApproveTask 方法在接口中被注释掉，暂未实现
// }

// TestApprovalTaskHandler_ProcessTaskRequiresTaskId - 方法未实现，跳过测试
// func TestApprovalTaskHandler_ProcessTaskRequiresTaskId(t *testing.T) {
//	 // ProcessTask 方法在接口中被注释掉，暂未实现
// }

func TestApprovalTaskHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalTaskService := mock_service.NewMockApprovalTaskService(ctrl)

	task := &model.ApprovalTask{
		ApprovalCode: "APPROVAL_001",
		NodeCode:     "NODE_001",
		Status:       model.TaskStatusPending,
	}

	mockApprovalTaskService.EXPECT().Create(gomock.Any()).Return(nil)

	approvalTaskHandler := handler.NewApprovalTaskHandler(createTaskTestHandler(), mockApprovalTaskService, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/approval-tasks", approvalTaskHandler.Create)

	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/approval-tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestApprovalTaskHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalTaskService := mock_service.NewMockApprovalTaskService(ctrl)

	task := &model.ApprovalTask{
		ID:           1,
		ApprovalCode: "APPROVAL_001",
		NodeCode:     "NODE_001",
		Status:       model.TaskStatusPending,
	}

	mockApprovalTaskService.EXPECT().Update(gomock.Any()).Return(nil)

	approvalTaskHandler := handler.NewApprovalTaskHandler(createTaskTestHandler(), mockApprovalTaskService, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/approval-tasks/:id", approvalTaskHandler.Update)

	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("PUT", "/approval-tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestApprovalTaskHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalTaskService := mock_service.NewMockApprovalTaskService(ctrl)

	deletedTask := &model.ApprovalTask{
		ID:           1,
		ApprovalCode: "APPROVAL_001",
		NodeCode:     "NODE_001",
		Status:       model.TaskStatusPending,
	}

	mockApprovalTaskService.EXPECT().Delete(uint(1)).Return(deletedTask, nil)

	approvalTaskHandler := handler.NewApprovalTaskHandler(createTaskTestHandler(), mockApprovalTaskService, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/approval-tasks/:id", approvalTaskHandler.Delete)

	req, _ := http.NewRequest("DELETE", "/approval-tasks/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
