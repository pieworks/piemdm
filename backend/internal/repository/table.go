package repository

import (
	"piemdm/internal/model"

	"github.com/gin-gonic/gin"
)

type TableRepository interface {
	First(sel string, tableField model.Table) (*model.Table, error)
	FindOne(id uint) (*model.Table, error)
	Find(sel string, where map[string]any) ([]*model.Table, error)
	FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.Table, error)
	Create(c *gin.Context, table *model.Table) error
	Update(c *gin.Context, table *model.Table) error
	BatchUpdate(c *gin.Context, ids []uint, table *model.Table) error
	Delete(c *gin.Context, id uint) (*model.Table, error)
	BatchDelete(c *gin.Context, ids []uint) error
}
type tableRepository struct {
	*Repository
	base Base
}

func NewTableRepository(repository *Repository, base Base) TableRepository {
	return &tableRepository{
		Repository: repository,
		base:       base,
	}
}

func (r *tableRepository) First(sel string, table model.Table) (*model.Table, error) {
	if sel == "" {
		sel = "*"
	}
	if err := r.base.First(&table, []string{}); err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *tableRepository) FindOne(id uint) (*model.Table, error) {
	var table model.Table
	if err := r.base.FirstById(&table, id); err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *tableRepository) Find(sel string, where map[string]any) ([]*model.Table, error) {
	var tables []*model.Table
	var table model.Table
	if sel == "" {
		sel = "*"
	}

	// err := r.base.FindPage(table, &tables, 1, 20, total, where, []string{}, "ID desc")
	// preloads := []string{}

	err := r.base.Find(table, &tables, sel, where, "sort asc")
	if err != nil {
		r.logger.Error("获取模型信息失败", "err", err)
	}
	return tables, nil
}

func (r *tableRepository) FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.Table, error) {
	var tables []*model.Table
	var table model.Table

	// if err := r.db.Where("deleted_at is null").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tables).Error; err != nil {
	// 	return nil, err
	// }

	// if err := r.db.Model(table).Where("deleted_at is null").Count(total).Error; err != nil {
	// 	r.logger.Error("apptoval repository count err", "err", err)
	// 	return nil, err
	// }
	preloads := []string{}
	err := r.base.FindPage(table, &tables, page, pageSize, total, where, preloads, "ID desc")
	if err != nil {
		r.logger.Error("获取模型信息失败", "err", err)
	}
	return tables, nil
}

func (r *tableRepository) Create(c *gin.Context, table *model.Table) error {
	if err := r.db.WithContext(c).Create(table).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableRepository) Update(c *gin.Context, table *model.Table) error {
	if err := r.db.WithContext(c).Updates(&table).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableRepository) BatchUpdate(c *gin.Context, ids []uint, table *model.Table) error {
	if err := r.db.WithContext(c).Model(&table).Where("id in ?", ids).Updates(table).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableRepository) Delete(c *gin.Context, id uint) (*model.Table, error) {
	var table model.Table
	if err := r.db.WithContext(c).Where("id = ?", id).First(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *tableRepository) BatchDelete(c *gin.Context, ids []uint) error {
	var tables []model.Table
	// 使用下面的调用方式，才能把参数传到 Hooks 里面。
	// //单条删除
	// db.Where("id =1").Find(&table).Delete(&table)
	// //多条删除
	// db.Where("id in ?", ids).Find(&tables).Delete(&tables)

	// if err := r.db.Delete(&table, ids).Error; err != nil {
	if err := r.db.WithContext(c).Where("id in ?", ids).Find(&tables).Delete(&tables).Error; err != nil {
		return err
	}
	return nil
}
