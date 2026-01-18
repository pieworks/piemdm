package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"piemdm/internal/constants"
	"piemdm/internal/handler"
	"piemdm/internal/model"
	mock_service "piemdm/test/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTableFieldHandler_CreateTableField(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success_AllFieldTypes", func(t *testing.T) {
		presets := constants.GetAllFieldTypePresets()
		for fieldType := range presets {
			// Capture range variable
			ft := fieldType
			t.Run(ft, func(t *testing.T) {
				// Prepare request
				reqBody := map[string]interface{}{
					"code":       "test_field_" + ft,
					"table_code": "test_table",
					"name":       "Test Field " + ft,
					"field_type": ft,
					"sort":       10,
				}
				body, _ := json.Marshal(reqBody)
				req, _ := http.NewRequest("POST", "/api/v1/admin/table_fields", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				// Expected model
				expectedField := &model.TableField{
					Code:      "test_field_" + ft,
					TableCode: "test_table",
					Name:      "Test Field " + ft,
					FieldType: ft,
					Sort:      10,
					Status:    "Normal", // Default value
				}

				// Mock expectations for this specific call
				// NOTE: In a loop with one controller, we need valid ordering or separate controllers if running parallel.
				// Here we run sequentially, so we expect one call per iteration.
				mockTableFieldService.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx *gin.Context, field *model.TableField) error {
					assert.Equal(t, expectedField.Code, field.Code, "Code mismatch for type %s", ft)
					assert.Equal(t, expectedField.TableCode, field.TableCode, "TableCode mismatch for type %s", ft)
					assert.Equal(t, expectedField.Name, field.Name, "Name mismatch for type %s", ft)
					assert.Equal(t, expectedField.FieldType, field.FieldType, "FieldType mismatch for type %s", ft)
					assert.Equal(t, expectedField.Status, field.Status, "Status mismatch for type %s", ft)

					// Simulate successful creation
					field.ID = 1
					field.CreatedAt = &time.Time{}
					field.UpdatedAt = &time.Time{}
					return nil
				})

				// Execute
				h.Create(c)

				// Verify
				assert.Equal(t, http.StatusOK, w.Code, "HTTP status mismatch for type %s", ft)

				var data model.TableField
				err := json.Unmarshal(w.Body.Bytes(), &data)
				assert.NoError(t, err, "Unmarshal error for type %s", ft)
				assert.Equal(t, expectedField.Code, data.Code, "Response Code mismatch for type %s", ft)
				assert.Equal(t, uint(1), data.ID, "Response ID mismatch for type %s", ft)
			})
		}
	})

	t.Run("BindingError_MissingRequired", func(t *testing.T) {
		// Prepare request with missing required fields
		reqBody := map[string]interface{}{
			"code": "test_field",
			// Missing table_code, name
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/admin/table_fields", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		h.Create(c)

		// Verify
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		// Prepare request
		reqBody := map[string]interface{}{
			"code":       "test_field",
			"table_code": "test_table",
			"name":       "Test Field",
			"field_type": "text",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/admin/table_fields", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Mock expectations
		mockTableFieldService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("service error"))

		// Execute
		h.Create(c)

		// Verify
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTableFieldHandler_ListTableFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table_fields?table_code=test_table", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockFields := []*model.TableField{{Code: "f1"}, {Code: "f2"}}
		var total int64 = 2

		mockTableFieldService.EXPECT().List(1, 15, gomock.Any(), gomock.Any()).DoAndReturn(
			func(page, pageSize int, outTotal *int64, where map[string]interface{}) ([]*model.TableField, error) {
				*outTotal = total
				assert.Equal(t, "test_table", where["table_code"])
				return mockFields, nil
			})

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("BindingError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table_fields", nil) // Missing table_code
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.List(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTableFieldHandler_GetTableField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table_fields/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockField := &model.TableField{Code: "test"}
		mockTableFieldService.EXPECT().Get(uint(1)).Return(mockField, nil)

		h.Get(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table_fields/1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockTableFieldService.EXPECT().Get(uint(1)).Return(nil, errors.New("not found"))

		h.Get(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTableFieldHandler_UpdateTableField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]any{
			"id":         1,
			"code":       "updated_code",
			"table_code": "table1",
			"name":       "Updated Name",
			"type":       "string",
			"length":     100,
			"status":     "Normal",
			"field_type": "text",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/api/v1/admin/table_fields", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockTableFieldService.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx *gin.Context, field *model.TableField) error {
			assert.Equal(t, uint(1), field.ID)
			assert.Equal(t, "updated_code", field.Code)
			return nil
		})

		h.Update(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTableFieldHandler_BatchUpdateTableFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]any{
			"ids":    []uint{1, 2},
			"status": "Frozen",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/admin/table_fields/batch-update", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockTableFieldService.EXPECT().BatchUpdate(gomock.Any(), []uint{1, 2}, gomock.Any()).DoAndReturn(
			func(ctx *gin.Context, ids []uint, field *model.TableField) error {
				assert.Equal(t, "Frozen", field.Status)
				return nil
			})

		h.BatchUpdate(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTableFieldHandler_BatchDeleteTableFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]any{
			"ids": []uint{1, 2},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/admin/table_fields/batch-delete", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockTableFieldService.EXPECT().BatchDelete(gomock.Any(), []uint{1, 2}).Return(nil)

		h.BatchDelete(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTableFieldHandler_PublicTableField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]any{
			"table_code": "test_table",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/admin/table_fields/public", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockTableFieldService.EXPECT().Public(gomock.Any(), "test_table").Return(nil)

		h.Public(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		// Check that it returns the expected data structure if possible, though exact content depends on implementation
	})
}

func TestTableFieldHandler_GetTableFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table/fields?table_code=test_table", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockFields := []*model.FieldMetadata{{Code: "f1"}}
		mockTableFieldService.EXPECT().GetTableFields("test_table").Return(mockFields, nil)

		h.GetTableFields(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MissingCode", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table/fields", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.GetTableFields(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTableFieldHandler_GetTableOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/table/test_table/options", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "table_code", Value: "test_table"}}

		mockOptions := []map[string]interface{}{{"label": "A", "value": "1"}}
		mockTableFieldService.EXPECT().GetTableOptions("test_table", gomock.Any()).Return(mockOptions, nil)

		h.GetTableOptions(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTableFieldHandler_GetFieldTypePresets(t *testing.T) {
	// This tests a method that uses config.* which is static.
	// We mainly ensure it doesn't panic and returns 200.

	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/field-type-presets", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.GetFieldTypePresets(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTableFieldHandler_GetFieldTypeGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableFieldService := mock_service.NewMockTableFieldService(ctrl)
	h := handler.NewTableFieldHandler(createUserTestHandler(), mockTableFieldService)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/field-type-groups", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		h.GetFieldTypeGroups(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
