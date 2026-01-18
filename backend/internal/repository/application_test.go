package repository_test

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gorm_mysql "gorm.io/driver/mysql"
	gorm "gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupApplicationRepositoryTest(t *testing.T) (sqlmock.Sqlmock, repository.ApplicationRepository, *gorm.DB, *log.Logger) {
	// 设置Gin为测试模式，避免调试信息输出
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gdb, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // 表名前缀
			SingularTable: false, // 使用单数表名
		},
	})
	assert.NoError(t, err)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Use the actual constructor for Repository
	repoInstance := repository.NewRepository(gdb, nil, &log.Logger{Logger: logger})
	baseRepo := repository.NewBaseRepository(repoInstance)
	repo := repository.NewApplicationRepository(repoInstance, baseRepo)

	return mock, repo, gdb, &log.Logger{Logger: logger}
}

func TestApplicationRepository_Create(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	app := &model.Application{
		AppId:     "new_app",
		Name:      "New Application",
		Status:    "Normal",
		CreatedBy: "test_user",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `applications` (`app_id`,`app_secret`,`name`,`ip`,`description`,`status`,`created_by`,`updated_by`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(app.AppId, app.AppSecret, app.Name, app.IP, app.Description, app.Status, app.CreatedBy, app.UpdatedBy, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := repo.Create(c, app)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplicationRepository_FindOne(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	expectedApp := &model.Application{
		ID:        1,
		AppId:     "test_app_id",
		Name:      "Test App",
		Status:    "Normal",
		CreatedBy: "test_user",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "app_id", "name", "status", "created_by", "created_at"}).
		AddRow(expectedApp.ID, expectedApp.AppId, expectedApp.Name, expectedApp.Status, expectedApp.CreatedBy, expectedApp.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `applications` WHERE `applications`.`id` = ? ORDER BY `applications`.`id` LIMIT ?")).
		WithArgs(1, 1).WillReturnRows(rows)

	app, err := repo.FindOne(1)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, expectedApp.Name, app.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplicationRepository_FindPage(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	page := 1
	pageSize := 10
	var total int64
	where := map[string]any{"status": "Normal"}

	expectedApps := []*model.Application{
		{ID: 1, AppId: "app1", Name: "App One", Status: "Normal"},
		{ID: 2, AppId: "app2", Name: "App Two", Status: "Normal"},
	}

	// Expect COUNT query
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `applications` WHERE status = ?")).
		WithArgs("Normal").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Expect SELECT query
	rows := sqlmock.NewRows([]string{"id", "app_id", "name", "status"}).
		AddRow(expectedApps[0].ID, expectedApps[0].AppId, expectedApps[0].Name, expectedApps[0].Status).
		AddRow(expectedApps[1].ID, expectedApps[1].AppId, expectedApps[1].Name, expectedApps[1].Status)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `applications` WHERE status = ? ORDER BY ID desc LIMIT ?")).
		WithArgs("Normal", pageSize).WillReturnRows(rows)

	apps, err := repo.FindPage(page, pageSize, &total, where)
	assert.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, int64(2), total)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplicationRepository_Update(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	app := &model.Application{
		ID:   1,
		Name: "Updated App Name",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `name`=?,`updated_at`=? WHERE `applications`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs(app.Name, sqlmock.AnyArg(), app.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := repo.Update(c, app)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplicationRepository_BatchUpdate(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	ids := []uint{1, 2, 3}
	app := &model.Application{
		Status: "Frozen",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `status`=?,`updated_at`=? WHERE id in (?,?,?) AND `applications`.`deleted_at` IS NULL")).
		WithArgs(app.Status, sqlmock.AnyArg(), ids[0], ids[1], ids[2]).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := repo.BatchUpdate(c, ids, app)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplicationRepository_Delete(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	id := uint(1)

	mock.ExpectBegin()
	// BeforeDelete hook: 先更新状态为 Deleted 和 updated_by
	// 注意：当使用 Where("id = ?", id).Delete(&model.Application{}) 时，
	// GORM会创建一个新的Application对象，其ID为默认值0，
	// 所以BeforeDelete钩子中的m.ID是0，导致WHERE条件中的id=0
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `status`=?,`updated_by`=?,`updated_at`=? WHERE id = ? AND `applications`.`deleted_at` IS NULL")).
		WithArgs("Deleted", "", sqlmock.AnyArg(), 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// 然后执行软删除
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `deleted_at`=? WHERE id = ? AND `applications`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := repo.Delete(c, id)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplicationRepository_BatchDelete(t *testing.T) {
	mock, repo, gdb, _ := setupApplicationRepositoryTest(t)
	defer func() {
		if sqlDB, err := gdb.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	ids := []uint{1, 2, 3}

	// GORM will first find the records, then update deleted_at
	rows := sqlmock.NewRows([]string{"id", "app_id"}).
		AddRow(1, "app1").
		AddRow(2, "app2").
		AddRow(3, "app3")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `applications` WHERE id in (?,?,?) AND `applications`.`deleted_at` IS NULL")).
		WithArgs(ids[0], ids[1], ids[2]).WillReturnRows(rows)
	mock.ExpectBegin()

	// BeforeDelete hook for each record: 先更新状态为 Deleted 和 updated_by (注意实际的 SQL 有额外的 AND `id` = ? 条件)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `status`=?,`updated_by`=?,`updated_at`=? WHERE id = ? AND `applications`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs("Deleted", "", sqlmock.AnyArg(), ids[0], ids[0]).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `status`=?,`updated_by`=?,`updated_at`=? WHERE id = ? AND `applications`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs("Deleted", "", sqlmock.AnyArg(), ids[1], ids[1]).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `status`=?,`updated_by`=?,`updated_at`=? WHERE id = ? AND `applications`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs("Deleted", "", sqlmock.AnyArg(), ids[2], ids[2]).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 然后批量软删除 (注意实际的 SQL 有额外的 IN 条件)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `applications` SET `deleted_at`=? WHERE id in (?,?,?) AND `applications`.`deleted_at` IS NULL AND `applications`.`id` IN (?,?,?)")).
		WithArgs(sqlmock.AnyArg(), ids[0], ids[1], ids[2], ids[0], ids[1], ids[2]).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := repo.BatchDelete(c, ids)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
