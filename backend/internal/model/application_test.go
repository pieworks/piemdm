package model_test

import (
	"context"
	"testing"

	"piemdm/internal/model"

	"github.com/stretchr/testify/assert"
	gorm_sqlite "gorm.io/driver/sqlite"
	gorm "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(gorm_sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.Application{})
	assert.NoError(t, err)

	return db
}

func TestApplication_BeforeDelete(t *testing.T) {
	db := setupTestDB(t)

	app := &model.Application{
		AppId:     "test_app_id",
		Name:      "Test App",
		Status:    "Normal",
		CreatedBy: "test_user",
	}

	err := db.Create(app).Error
	assert.NoError(t, err)

	// Simulate BeforeDelete hook by calling Delete
	err = db.Delete(app).Error
	assert.NoError(t, err)

	// Retrieve the application to check its status
	var updatedApp model.Application
	db.Unscoped().First(&updatedApp, app.ID)

	assert.Equal(t, "Deleted", updatedApp.Status)
	assert.NotNil(t, updatedApp.DeletedAt.Valid)
}

func TestApplication_BeforeCreate(t *testing.T) {
	db := setupTestDB(t)

	// Create a context with user_name
	ctx := context.WithValue(context.Background(), "user_name", "test_creator")
	dbWithContext := db.WithContext(ctx)

	app := &model.Application{
		AppId:     "new_app_id",
		Name:      "New App",
		Status:    "Normal",
		AppSecret: "secret",
	}

	err := dbWithContext.Create(app).Error
	assert.NoError(t, err)

	assert.Equal(t, "test_creator", app.CreatedBy)
	assert.NotNil(t, app.CreatedAt)
}

func TestApplication_BeforeUpdate(t *testing.T) {
	db := setupTestDB(t)

	// First, create an application
	app := &model.Application{
		AppId:     "update_app_id",
		Name:      "Original App",
		Status:    "Normal",
		CreatedBy: "initial_user",
		AppSecret: "secret",
	}
	err := db.Create(app).Error
	assert.NoError(t, err)

	// Update the application with a new context user
	ctx := context.WithValue(context.Background(), "user_name", "updater_user")
	dbWithContext := db.WithContext(ctx)

	app.Name = "Updated App Name"
	err = dbWithContext.Save(app).Error
	assert.NoError(t, err)

	assert.Equal(t, "updater_user", app.UpdatedBy)
	assert.NotNil(t, app.UpdatedAt)
	assert.True(t, app.UpdatedAt.After(app.CreatedAt))
}
