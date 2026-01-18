package service_test

import (
	"errors"
	"testing"

	"piemdm/internal/model"
	"piemdm/internal/service"
	mock_repository "piemdm/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupApplicationServiceTest(t *testing.T) (*gomock.Controller, *mock_repository.MockApplicationRepository, service.ApplicationService) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockApplicationRepository(ctrl)
	svc := service.NewApplicationService(&service.Service{}, mockRepo)
	return ctrl, mockRepo, svc
}

func TestApplicationService_Get(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	expectedApp := &model.Application{ID: 1, Name: "Test App"}

	mockRepo.EXPECT().FindOne(uint(1)).Return(expectedApp, nil)

	app, err := svc.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedApp, app)

	// Test error scenario
	mockRepo.EXPECT().FindOne(uint(2)).Return(nil, errors.New("not found"))
	app, err = svc.Get(2)
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestApplicationService_List(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	page := 1
	pageSize := 10
	total := int64(2)
	where := map[string]any{"status": "Normal"}
	expectedApps := []*model.Application{
		{ID: 1, Name: "App1"},
		{ID: 2, Name: "App2"},
	}

	mockRepo.EXPECT().FindPage(page, pageSize, &total, where).Return(expectedApps, nil)

	apps, err := svc.List(page, pageSize, &total, where)
	assert.NoError(t, err)
	assert.Equal(t, expectedApps, apps)

	// Test error scenario
	mockRepo.EXPECT().FindPage(page, pageSize, &total, where).Return(nil, errors.New("db error"))
	apps, err = svc.List(page, pageSize, &total, where)
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func TestApplicationService_Create(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	app := &model.Application{Name: "New App"}

	mockRepo.EXPECT().Create(gomock.Any(), app).Return(nil)

	err := svc.Create(nil, app)
	assert.NoError(t, err)

	// Test error scenario
	mockRepo.EXPECT().Create(gomock.Any(), app).Return(errors.New("create error"))
	err = svc.Create(nil, app)
	assert.Error(t, err)
}

func TestApplicationService_Update(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	app := &model.Application{ID: 1, Name: "Updated App"}

	mockRepo.EXPECT().Update(gomock.Any(), app).Return(nil)

	err := svc.Update(nil, app)
	assert.NoError(t, err)

	// Test error scenario
	mockRepo.EXPECT().Update(gomock.Any(), app).Return(errors.New("update error"))
	err = svc.Update(nil, app)
	assert.Error(t, err)
}

func TestApplicationService_BatchUpdate(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	ids := []uint{1, 2}
	app := &model.Application{Status: "Frozen"}

	mockRepo.EXPECT().BatchUpdate(gomock.Any(), ids, app).Return(nil)

	err := svc.BatchUpdate(nil, ids, app)
	assert.NoError(t, err)

	// Test error scenario
	mockRepo.EXPECT().BatchUpdate(gomock.Any(), ids, app).Return(errors.New("batch update error"))
	err = svc.BatchUpdate(nil, ids, app)
	assert.Error(t, err)
}

func TestApplicationService_Delete(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	id := uint(1)

	mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	err := svc.Delete(nil, id)
	assert.NoError(t, err)

	// Test error scenario
	mockRepo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("delete error"))
	err = svc.Delete(nil, id)
	assert.Error(t, err)
}

func TestApplicationService_BatchDelete(t *testing.T) {
	ctrl, mockRepo, svc := setupApplicationServiceTest(t)
	defer ctrl.Finish()

	ids := []uint{1, 2}

	mockRepo.EXPECT().BatchDelete(gomock.Any(), ids).Return(nil)

	err := svc.BatchDelete(nil, ids)
	assert.NoError(t, err)

	// Test error scenario
	mockRepo.EXPECT().BatchDelete(gomock.Any(), ids).Return(errors.New("batch delete error"))
	err = svc.BatchDelete(nil, ids)
	assert.Error(t, err)
}
