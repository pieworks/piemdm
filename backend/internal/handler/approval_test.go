package handler_test

import (
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

func createApprovalTestHandler() *handler.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	testLogger := &log.Logger{
		Logger: logger,
	}
	return handler.NewHandler(testLogger)
}

func TestApprovalHandler_ListApprovals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalService := mock_service.NewMockApprovalService(ctrl)

	expectedList := []*model.Approval{
		{
			ID:              1,
			Code:            "APPROVAL_001",
			Title:           "测试审批申请",
			ApprovalDefCode: "TEST_DEF_001",
			Status:          model.ApprovalStatusPending,
		},
	}

	mockApprovalService.EXPECT().FindPendingByAssignee(gomock.Any(), 1, 10, gomock.Any(), "").Return(expectedList, nil)

	approvalHandler := handler.NewApprovalHandler(createApprovalTestHandler(), mockApprovalService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/approvals", approvalHandler.List)

	req, _ := http.NewRequest("GET", "/approvals?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestApprovalHandler_GetApproval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalService := mock_service.NewMockApprovalService(ctrl)

	expectedApproval := &model.Approval{
		ID:              1,
		Code:            "APPROVAL_001",
		Title:           "测试审批申请",
		ApprovalDefCode: "TEST_DEF_001",
		Status:          model.ApprovalStatusPending,
	}

	mockApprovalService.EXPECT().Get(uint(1)).Return(expectedApproval, nil)

	approvalHandler := handler.NewApprovalHandler(createApprovalTestHandler(), mockApprovalService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/approvals/:id", approvalHandler.Get)

	req, _ := http.NewRequest("GET", "/approvals/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// TestApprovalHandler_StartApproval - 方法未实现，跳过测试
// func TestApprovalHandler_StartApproval(t *testing.T) {
//	 // StartApproval 方法在接口中被注释掉，暂未实现
// }

// TestApprovalHandler_SubmitApproval - 方法未实现，跳过测试
// func TestApprovalHandler_SubmitApproval(t *testing.T) {
//	 // SubmitApproval 方法在接口中被注释掉，暂未实现
// }

// TestApprovalHandler_CancelApproval - 方法未实现，跳过测试
// func TestApprovalHandler_CancelApproval(t *testing.T) {
//	 // CancelApproval 方法在接口中被注释掉，暂未实现
// }

// TestApprovalHandler_GetApprovalByCode - 方法未实现，跳过测试
// func TestApprovalHandler_GetApprovalByCode(t *testing.T) {
//	 // GetApprovalByCode 方法在接口中被注释掉，暂未实现
// }

// TestApprovalHandler_GetApprovalsByApplicant - 方法未实现，跳过测试
// func TestApprovalHandler_GetApprovalsByApplicant(t *testing.T) {
//	 // GetApprovalsByApplicant 方法在接口中被注释掉，暂未实现
// }

// TestApprovalHandler_CreateApproval - 方法已从接口移除，跳过测试
// func TestApprovalHandler_CreateApproval(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockApprovalService := mock_service.NewMockApprovalService(ctrl)
//
// 	approval := &model.Approval{
// 		Code:            "APPROVAL_001",
// 		Title:           "测试审批申请",
// 		ApprovalDefCode: "TEST_DEF_001",
// 		Status:          model.ApprovalStatusPending,
// 	}
//
// 	// CreateApproval 已废弃，不需要 Mock 期望
//
// 	approvalHandler := handler.NewApprovalHandler(createApprovalTestHandler(), mockApprovalService)
//
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	router.POST("/approvals", approvalHandler.Create)
//
// 	body, _ := json.Marshal(approval)
// 	req, _ := http.NewRequest("POST", "/approvals", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)
//
// 	// CreateApproval 已废弃，期望返回 400 错误
// 	assert.Equal(t, http.StatusBadRequest, resp.Code)
// 	assert.Contains(t, resp.Body.String(), "请使用StartApproval接口启动审批流程")
// }

// TestApprovalHandler_UpdateApproval - 方法已从接口移除，跳过测试
// func TestApprovalHandler_UpdateApproval(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockApprovalService := mock_service.NewMockApprovalService(ctrl)
//
// 	approval := &model.Approval{
// 		ID:              1,
// 		Code:            "APPROVAL_001",
// 		Title:           "更新审批申请",
// 		ApprovalDefCode: "TEST_DEF_001",
// 		Status:          model.ApprovalStatusPending,
// 	}
//
// 	// UpdateApproval 不支持直接更新，不需要 Mock 期望
//
// 	approvalHandler := handler.NewApprovalHandler(createApprovalTestHandler(), mockApprovalService)
//
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	router.PUT("/approvals/:id", approvalHandler.Update)
//
// 	body, _ := json.Marshal(approval)
// 	req, _ := http.NewRequest("PUT", "/approvals/1", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)
//
// 	// UpdateApproval 不支持直接更新，期望返回 400 错误
// 	assert.Equal(t, http.StatusBadRequest, resp.Code)
// 	assert.Contains(t, resp.Body.String(), "审批实例不支持直接更新")
// }

func TestApprovalHandler_DeleteApproval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalService := mock_service.NewMockApprovalService(ctrl)

	deletedApproval := &model.Approval{
		ID:              1,
		Code:            "APPROVAL_001",
		Title:           "删除的审批申请",
		ApprovalDefCode: "TEST_DEF_001",
		Status:          model.ApprovalStatusPending,
	}

	mockApprovalService.EXPECT().Delete(gomock.Any(), uint(1)).Return(deletedApproval, nil)

	approvalHandler := handler.NewApprovalHandler(createApprovalTestHandler(), mockApprovalService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/approvals/:id", approvalHandler.Delete)

	req, _ := http.NewRequest("DELETE", "/approvals/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
