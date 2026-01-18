package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"piemdm/internal/handler"
	"piemdm/internal/model"
	"piemdm/pkg/log"
	mock_service "piemdm/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupApplicationHandlerTest(t *testing.T) (*gin.Engine, *mock_service.MockApplicationService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockService := mock_service.NewMockApplicationService(ctrl)

	logger := &log.Logger{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	baseHandler := handler.NewHandler(logger)
	sysAppHandler := handler.NewApplicationHandler(baseHandler, mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/applications", sysAppHandler.List)
	router.GET("/api/v1/applications/:id", sysAppHandler.Get)
	router.POST("/api/v1/applications", sysAppHandler.Create)
	router.PUT("/api/v1/applications/:id", sysAppHandler.Update)
	router.DELETE("/api/v1/applications/:id", sysAppHandler.Delete)
	router.POST("/api/v1/applications/batch", sysAppHandler.BatchCreate)
	router.PUT("/api/v1/applications/batch", sysAppHandler.BatchUpdate)
	router.DELETE("/api/v1/admin/applications/batch", sysAppHandler.BatchDelete)

	return router, mockService, ctrl
}

func TestApplicationHandler_List(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful listing
	expectedApps := []*model.Application{
		{ID: 1, Name: "App1"},
		{ID: 2, Name: "App2"},
	}
	mockService.EXPECT().List(1, 10, gomock.Any(), gomock.Any()).Return(expectedApps, nil)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/applications?page=1&pageSize=10", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var apps []*model.Application
	_ = json.Unmarshal(rec.Body.Bytes(), &apps)
	assert.Len(t, apps, 2)
	assert.Equal(t, "App1", apps[0].Name)
	assert.Equal(t, "App2", apps[1].Name)

	// Test service error
	mockService.EXPECT().List(1, 10, gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/applications?page=1&pageSize=10", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Equal(t, "service error", errorRes["message"])

	// Test invalid request (e.g., invalid page size)
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/applications?page=1&pageSize=abc", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Contains(t, errorRes["message"].(string), "invalid")
}

func TestApplicationHandler_Get(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful retrieval
	expectedApp := &model.Application{ID: 1, Name: "Test App"}
	mockService.EXPECT().Get(uint(1)).Return(expectedApp, nil)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/applications/1", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var app model.Application
	_ = json.Unmarshal(rec.Body.Bytes(), &app)
	assert.Equal(t, uint(1), app.ID)
	assert.Equal(t, "Test App", app.Name)

	// Test service error (not found)
	mockService.EXPECT().Get(uint(2)).Return(nil, errors.New("not found"))
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/applications/2", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Equal(t, "not found", errorRes["message"])

	// Test invalid ID in URI
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/applications/abc", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Contains(t, errorRes["message"].(string), "invalid")
}

func TestApplicationHandler_Create(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful creation
	newApp := map[string]string{
		"name":        "New App",
		"status":      "Normal",
		"description": "A new application",
	}

	mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	jsonBody, _ := json.Marshal(newApp)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/applications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Test invalid request (missing required fields)
	invalidApp := map[string]string{"description": "test"}
	jsonBody, _ = json.Marshal(invalidApp)
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/applications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Contains(t, errorRes["message"].(string), "required")
}

func TestApplicationHandler_Update(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful update
	updatedApp := map[string]any{
		"id":          1,
		"name":        "Updated App Name",
		"status":      "Frozen",
		"description": "Updated description",
	}

	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	jsonBody, _ := json.Marshal(updatedApp)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/applications/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Test invalid request (missing required fields)
	invalidApp := map[string]any{"id": 1}
	jsonBody, _ = json.Marshal(invalidApp)
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/applications/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Contains(t, errorRes["message"].(string), "required")
}

func TestApplicationHandler_Delete(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful deletion
	mockService.EXPECT().Delete(gomock.Any(), uint(1)).Return(nil)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/applications/1", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Empty(t, result)

	// Test service error
	mockService.EXPECT().Delete(gomock.Any(), uint(2)).Return(errors.New("service delete error"))
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/applications/2", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Equal(t, "service delete error", errorRes["message"])

	// Test invalid ID in URI
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/applications/abc", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Contains(t, errorRes["message"].(string), "invalid")
}

func TestApplicationHandler_BatchCreate(t *testing.T) {
	router, _, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/applications/batch", nil)
	router.ServeHTTP(rec, req)

	// BatchCreate is not implemented, returns empty response
	assert.Equal(t, http.StatusOK, rec.Code)
	var result map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Empty(t, result)
}

func TestApplicationHandler_BatchUpdate(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful batch update
	ids := []uint{1, 2}
	updateReq := struct {
		Ids    []uint `form:"ids"`
		Status string `form:"status"`
	}{
		Ids:    ids,
		Status: "Frozen",
	}
	mockService.EXPECT().BatchUpdate(gomock.Any(), ids, gomock.Any()).DoAndReturn(func(ctx interface{}, ids []uint, app *model.Application) error {
		assert.Equal(t, "Frozen", app.Status)
		return nil
	}).Times(1)

	formData := strings.NewReader(fmt.Sprintf("ids=%d&ids=%d&status=%s", ids[0], ids[1], updateReq.Status))
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/applications/batch", formData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var res map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, "success", res["message"])

	// Test service error
	mockService.EXPECT().BatchUpdate(gomock.Any(), ids, gomock.Any()).Return(errors.New("service batch update error")).Times(1)
	formData2 := strings.NewReader(fmt.Sprintf("ids=%d&ids=%d&status=%s", ids[0], ids[1], updateReq.Status))
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/applications/batch", formData2)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Equal(t, "service batch update error", errorRes["message"])

	// Test invalid request (missing status) - status 是必填的,应该返回 400
	invalidFormData := strings.NewReader(fmt.Sprintf("ids=%d&ids=%d", ids[0], ids[1]))
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/applications/batch", invalidFormData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code) // status 是必填的
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Contains(t, errorRes["message"].(string), "required")
}

func TestApplicationHandler_BatchDelete(t *testing.T) {
	router, mockService, ctrl := setupApplicationHandlerTest(t)
	defer ctrl.Finish()

	// Test successful batch deletion
	ids := []uint{1, 2}
	mockService.EXPECT().BatchDelete(gomock.Any(), ids).Return(nil).Times(1)

	reqBody := map[string]any{
		"ids": ids,
	}
	jsonBody, _ := json.Marshal(reqBody)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/applications/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var res map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, "success", res["message"])

	// Test service error
	mockService.EXPECT().BatchDelete(gomock.Any(), ids).Return(errors.New("service batch delete error")).Times(1)
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/admin/applications/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var errorRes map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &errorRes)
	assert.Equal(t, "service batch delete error", errorRes["message"])

	// Test invalid request (empty ids array)
	emptyIds := []uint{}
	mockService.EXPECT().BatchDelete(gomock.Any(), emptyIds).Return(nil).Times(1)

	emptyReqBody := map[string]any{
		"ids": emptyIds,
	}
	emptyJsonBody, _ := json.Marshal(emptyReqBody)
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/admin/applications/batch", bytes.NewBuffer(emptyJsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // 会成功，即使 ids 为空
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, "success", res["message"])
}
