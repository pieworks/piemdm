package repository_test

import (
	"net/http/httptest"
	"strconv"
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

func setupRepository(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
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
	userRepo := repository.NewUserRepository(repo, base)

	return userRepo, mock
}

func TestUserRepository_Create(t *testing.T) {
	userRepo, mock := setupRepository(t)

	now := time.Now()
	user := &model.User{
		ID:               1,
		EmployeeID:       "EMP001",
		Username:         "test",
		Password:         "password",
		FirstName:        "First",
		LastName:         "Last",
		DisplayName:      "Test",
		Email:            "test@example.com",
		Phone:            "12345678",
		Language:         "zh-cn",
		Sex:              "Male",
		Avatar:           "",
		Description:      "",
		Admin:            "No",
		SuperiorUsername: "admin",
		LastLogin:        &now,
		Status:           "Normal",
		CreatedBy:        "",
		UpdatedBy:        "",
		UpdatedAt:        &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := userRepo.Create(c, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	userRepo, mock := setupRepository(t)

	now := time.Now()
	user := &model.User{
		ID:               1,
		EmployeeID:       "EMP001",
		Username:         "test",
		Password:         "password",
		FirstName:        "First",
		LastName:         "Last",
		DisplayName:      "Test",
		Email:            "test@example.com",
		Phone:            "12345678",
		Language:         "zh-cn",
		Sex:              "Male",
		Avatar:           "",
		Description:      "",
		Admin:            "No",
		SuperiorUsername: "admin",
		LastLogin:        &now,
		Status:           "Normal",
		CreatedBy:        "",
		UpdatedBy:        "",
		UpdatedAt:        &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := userRepo.Update(c, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetById(t *testing.T) {
	userRepo, mock := setupRepository(t)

	userId := uint(123)

	rows := sqlmock.NewRows([]string{
		"id", "employee_id", "username", "password", "first_name", "last_name",
		"display_name", "email", "phone", "language", "sex",
		"avatar", "description", "admin", "superior_username", "last_login", "status", "created_by",
		"updated_by", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		123, "EMP001", "test", "password", "First", "Last",
		"Test", "test@example.com", "12345678", "zh-cn", "Male",
		"", "", "No", "admin", time.Now(), "Normal", "",
		"", time.Now(), time.Now(), nil,
	)
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE user_id = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(strconv.Itoa(int(userId)), 1).
		WillReturnRows(rows)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	user, err := userRepo.FindByUserID(c, strconv.Itoa(int(userId)))
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(123), user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByUsername(t *testing.T) {
	userRepo, mock := setupRepository(t)

	username := "test"

	rows := sqlmock.NewRows([]string{
		"id", "employee_id", "username", "password", "first_name", "last_name",
		"display_name", "email", "phone", "language", "sex",
		"avatar", "description", "admin", "superior_username", "last_login", "status", "created_by",
		"updated_by", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		1, "EMP001", "test", "password", "First", "Last",
		"Test", "test@example.com", "12345678", "zh-cn", "Male",
		"", "", "No", "admin", time.Now(), "Normal", "",
		"", time.Now(), time.Now(), nil,
	)
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE username = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(username, 1).
		WillReturnRows(rows)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	user, err := userRepo.FindByUsername(c, username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test", user.Username)

	assert.NoError(t, mock.ExpectationsWereMet())
}
