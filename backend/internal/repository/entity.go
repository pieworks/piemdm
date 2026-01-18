package repository

import (
	"fmt"
	"strings"

	"piemdm/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// Repository/数据库
// FindOne(article *models.Article) error
//    First， ordered by primary key
//    Last，ordered by primary key
//    Take， no specified order
// FindPage(page int, pageSize int, total *int64, where interface{}) ([]*models.Article, error)
//    FindInBatches，分批查询
// Find(where interface{}) ([]*models.Article,  error)
//    Find，Retrieving all objects
//    Find in，Find(&users, []int{1,2,3})
//    Pluck，查询单个列，并将结果扫描到切片，多列，您应该使用 Select 和 Scan
// Create(article *models.Article) error
//    FirstOrCreate，如果没有找到记录，可以使用包含更多的属性的结构体创建记录，
//    批量创建[]User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
//    Create From Map
// Update(article *models.Article) error
//    Save，保存所有的字段，即使字段是零值
//    Update， 更新单列
//    Updates，更新多列
//    批量更新要使用Model更新struct，使用table更新map
//    检查字段是否有变更，使用BeforeUpdate中的changed
// Delete(article *models.Article) error
// UpdateByIds(ids []int64) error
// DeleteByIds(ids []int64) error

// Controller/handler
// index
// view
// create
// update

// GetTables(PageNum, PageSize int, where interface{}) []*models.Article
// GetArticle(where interface{}) *models.Article
// AddArticle(article *models.Article) bool
// GetArticles(PageNum int, PageSize int, total *uint64, where interface{}) []*models.Article
// EditArticle(article models.Article) bool
// DeleteArticle(id int) bool
// ExistArticleByID(id int) bool
// GetArticleTotal(maps map[string]any) (count int)

// CheckUser(where interface{}) bool
// GetUserAvatar(sel *string, where interface{}) *string
// GetUserID(sel *string, where interface{}) int
// GetUsers(PageNum int, PageSize int, total *uint64, where interface{}) []*models.User
// AddUser(user *models.User) bool
// ExistUserByName(where interface{}) bool
// UpdateUser(user *models.User, role *models.Role) bool
// DeleteUser(id int) bool
// GetUserByID(id int) *models.User

type EntityRepository interface {
	// 基础查询
	FindOne(tableCode string, id uint) (map[string]any, error)
	FindPage(tableCode string, page, pageSize int, total *int64, where map[string]any) ([]map[string]any, error)
	FindLogPage(tableCode string, page, pageSize int, total *int64, where map[string]any) ([]map[string]any, error)
	Find(tableCode string, selectString string, where map[string]any) ([]map[string]any, error)

	// Base CRUD
	Create(c *gin.Context, tableCode string, entityMap any) error
	Update(c *gin.Context, tableCode string, entity any, where map[string]any) error
	Delete(c *gin.Context, tableCode string, id uint) error

	// Batch operations
	BatchUpdate(c *gin.Context, tableCode string, ids []uint, entityMap map[string]any) error
	BatchDelete(c *gin.Context, tableCode string, ids []uint) error
	// 统计操作
	GetStatisticsByStatus(tableCode string) (map[string]int64, error)
}

type entityRepository struct {
	*Repository
	source Base
	stfr   TableFieldRepository
}

func NewEntityRepository(repository *Repository, source Base, stfr TableFieldRepository) EntityRepository {
	return &entityRepository{
		Repository: repository,
		source:     source,
		stfr:       stfr,
	}
}

// getTableName 构建数据表名称
func (r *entityRepository) getTableName(tableCode string) string {
	return fmt.Sprintf("t_%s", tableCode)
}

// func (r *entityRepository) FindById(id uint) (*model.Entity, error) {
// 	var entity model.Entity
// 	entity.ID = uint(id)
// 	if err := r.source.First(&entity, []string{}, "*"); err != nil {
// 		return nil, err
// 	}
// 	//
// 	// if err := r.db.Where("id = ?", id).First(&entity).Error; err != nil {
// 	// 	return nil, err
// 	// }
// 	return &entity, nil
// }

func (r *entityRepository) FindOne(tableCode string, id uint) (map[string]any, error) {
	entity := map[string]any{}

	// 添加表前缀
	table := r.getTableName(tableCode)

	if err := r.db.Table(table).Where("id = ?", id).Take(&entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *entityRepository) FindPage(tableCode string, page, pageSize int, total *int64, where map[string]any) ([]map[string]any, error) {
	tableName := tableCode
	table := r.getTableName(tableName)

	// 前端统一查询方案: 移除 JOIN 逻辑
	// - 后端只返回原始 code 值(单选)或 code 数组(多选)
	// - 前端负责查询字典并格式化显示为 "code name"
	// - 优点: 简化后端,提高性能,统一处理方式

	// 构建条件
	conditionString, values, _ := BuildCondition(where)

	// 执行查询 - 直接查询主表,不做 JOIN
	var entities []map[string]any
	query := r.db.Table(table + " t").
		Select("t.*").
		Where("t.deleted_at is null")

	// 添加条件
	if conditionString != "" {
		query = query.Where(conditionString, values...)
	}

	// 分页查询
	if err := query.Order("t.id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&entities).Error; err != nil {
		// Check for MySQL Error 1146: Table doesn't exist
		if strings.Contains(err.Error(), "Error 1146") {
			r.logger.Warn("Table does not exist, suppressing error", "table", table)
			return []map[string]any{}, nil
		}
		r.logger.Error("查询失败", "err", err)
		return nil, err
	}

	// 查询总数
	countQuery := r.db.Table(table + " t").
		Where("t.deleted_at is null")

	if conditionString != "" {
		countQuery = countQuery.Where(conditionString, values...)
	}

	if err := countQuery.Count(total).Error; err != nil {
		r.logger.Error("查询总数出错", "err", err)
	}

	// 前端统一查询方案: 不再需要处理 JSON_OBJECT
	// 后端直接返回原始 code 值,由前端查询字典并格式化

	return entities, nil
}

func (r *entityRepository) FindLogPage(tableCode string, page, pageSize int, total *int64, where map[string]any) ([]map[string]any, error) {
	tableName := tableCode
	table := r.getTableName(tableName)

	var entities []map[string]any

	if err := r.db.Table(table).Where("deleted_at is null").Where(where).Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		panic(err)
	}

	err := r.db.Table(table).Where("deleted_at is null").Where(where).Count(total).Error
	if err != nil {
		r.logger.Error("查询总数出错", "err", err)
	}

	return entities, nil
}

func (r *entityRepository) Find(tableCode, selectString string, where map[string]any) ([]map[string]any, error) {
	table := r.getTableName(tableCode)
	var entities []map[string]any

	// 构建条件
	conditionString, values, _ := BuildCondition(where)

	if err := r.db.Table(table).Select(selectString).Where("deleted_at is null").Where(conditionString, values...).Find(&entities).Error; err != nil {
		r.logger.Error("Can't find Entity", "err", err)
	}

	return entities, nil
}

// TODO 可以传入map，统一转为 struct
func (r *entityRepository) Create(c *gin.Context, tableCode string, entity any) error {
	table := r.getTableName(tableCode)

	switch entityMap := entity.(type) {
	case map[string]any:
		// 生成entity map
		entityNew := r.stfr.BuildEntity(tableCode)

		// 过滤掉非数据库字段
		excludeFields := map[string]bool{
			"table_code": true,
			"reason":     true,
		}

		// 将输入的 entityMap 合并到 entityNew,跳过非数据库字段
		for k, v := range entityMap {
			if !excludeFields[k] {
				entityNew[k] = v
			}
		}

		// 查询配置了unique index的字段
		uniqueFields, err := r.stfr.Find("code,is_unique,index_name", map[string]any{
			"table_code": tableCode,
			"is_unique":  "Yes",
		})
		if err != nil {
			r.logger.Error("查询唯一索引字段失败", "error", err)
		}

		// 构建OnConflict的Columns
		var conflictColumns []clause.Column
		if len(uniqueFields) > 0 {
			// 使用unique fields作为冲突检测列
			for _, field := range uniqueFields {
				conflictColumns = append(conflictColumns, clause.Column{Name: field.Code})
			}
		} else {
			// 如果没有配置unique index,使用id作为默认冲突检测列
			conflictColumns = []clause.Column{{Name: "id"}}
		}

		// 构建要更新的列列表(排除冲突检测列)
		updateColumns := make([]string, 0)
		conflictColumnNames := make(map[string]bool)
		for _, col := range conflictColumns {
			conflictColumnNames[col.Name] = true
		}

		// 添加所有非冲突列到更新列表
		for k := range entityNew {
			if !conflictColumnNames[k] && k != "id" && k != "created_at" && k != "deleted_at" {
				updateColumns = append(updateColumns, k)
			}
		}

		// 配置OnConflict,指定冲突列和要更新的列
		if err := r.db.WithContext(c).Table(table).Clauses(clause.OnConflict{
			Columns:   conflictColumns,
			DoUpdates: clause.AssignmentColumns(updateColumns),
		}).Create(entityNew).Error; err != nil {
			return err
		}
	default:
		if err := r.db.WithContext(c).Table(table).Create(entity).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *entityRepository) Update(c *gin.Context, tableCode string, entity any, where map[string]any) error {
	table := r.getTableName(tableCode)
	err := r.db.WithContext(c).Table(table).Where(where).Updates(entity).Error
	if err != nil {
		r.logger.Error("更新失败", "err", err)
		return err
	}

	return nil
}

func (r *entityRepository) BatchUpdate(c *gin.Context, tableCode string, ids []uint, entityMap map[string]any) error {
	table := r.getTableName(tableCode)
	if err := r.db.WithContext(c).Table(table).Where("id in ?", ids).Updates(entityMap).Error; err != nil {
		return err
	}
	return nil
}

func (r *entityRepository) Delete(c *gin.Context, tableCode string, id uint) error {
	var entity model.Entity
	if err := r.db.WithContext(c).Where("id = ?", id).First(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *entityRepository) BatchDelete(c *gin.Context, tableCode string, ids []uint) error {
	var entities []model.Entity
	if err := r.db.WithContext(c).Where("id in ?", ids).Find(&entities).Delete(&entities).Error; err != nil {
		return err
	}
	return nil
}
func (r *entityRepository) GetStatisticsByStatus(tableCode string) (map[string]int64, error) {
	if tableCode == "" {
		return nil, fmt.Errorf("tableCode 不能为空")
	}
	table := r.getTableName(tableCode)
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Table(table).
		Select("status, count(*) as count").
		Where("deleted_at is null").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, res := range results {
		if res.Status == "" {
			stats["Normal"] += res.Count // 默认为 Normal
		} else {
			stats[res.Status] = res.Count
		}
	}
	return stats, nil
}
