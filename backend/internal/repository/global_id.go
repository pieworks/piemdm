package repository

import (
	"piemdm/internal/model"
)

type GlobalIdRepository interface {
	// 基础查询
	FindOne(id uint) (*model.GlobalId, error)
	FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.GlobalId, error)

	// Base CRUD
	Create(globalId *model.GlobalId) error
	Update(globalId *model.GlobalId) error
	Delete(id uint) (*model.GlobalId, error)

	// Batch operations
	BatchUpdate(ids []uint, globalId *model.GlobalId) error
	BatchDelete(ids []uint) error

	// 业务方法
	GetNewID(identifier string) uint
}
type globalIdRepository struct {
	*Repository
	source Base
}

func NewGlobalIdRepository(repository *Repository, source Base) GlobalIdRepository {
	return &globalIdRepository{
		Repository: repository,
		source:     source,
	}
}

func (r *globalIdRepository) FindOne(id uint) (*model.GlobalId, error) {
	var globalId model.GlobalId
	if err := r.source.FirstById(&globalId, id); err != nil {
		return nil, err
	}
	return &globalId, nil
}

func (r *globalIdRepository) FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.GlobalId, error) {
	var globalIds []*model.GlobalId
	var globalId model.GlobalId

	// if err := r.db.Where("deleted_at is null").Offset((page - 1) * pageSize).Limit(pageSize).Find(&globalIds).Error; err != nil {
	// 	return nil, err
	// }

	// if err := r.db.Model(globalId).Where("deleted_at is null").Count(total).Error; err != nil {
	// 	r.logger.Error("apptoval repository count err", "err", err)
	// 	return nil, err
	// }
	preloads := []string{}
	err := r.source.FindPage(globalId, &globalIds, page, pageSize, total, where, preloads, "ID desc")
	if err != nil {
		r.logger.Error("获取模型信息失败", "err", err)
	}
	return globalIds, nil
}

func (r *globalIdRepository) Create(globalId *model.GlobalId) error {
	if err := r.source.Create(globalId); err != nil {
		return err
	}
	return nil
}

func (r *globalIdRepository) Update(globalId *model.GlobalId) error {
	if err := r.source.Updates(&globalId, globalId); err != nil {
		return err
	}
	return nil
}

func (r *globalIdRepository) BatchUpdate(ids []uint, globalId *model.GlobalId) error {
	if err := r.db.Model(&globalId).Where("id in ?", ids).Updates(globalId).Error; err != nil {
		return err
	}
	return nil
}

func (r *globalIdRepository) Delete(id uint) (*model.GlobalId, error) {
	var globalId model.GlobalId
	if err := r.db.Where("id = ?", id).First(&globalId).Error; err != nil {
		return nil, err
	}
	return &globalId, nil
}

func (r *globalIdRepository) BatchDelete(ids []uint) error {
	var globalIds []model.GlobalId
	// 使用下面的调用方式，才能把参数传到 Hooks 里面。
	// //单条删除
	// db.Where("id =1").Find(&globalId).Delete(&globalId)
	// //多条删除
	// db.Where("id in ?", ids).Find(&globalIds).Delete(&globalIds)

	// if err := r.db.Delete(&globalId, ids).Error; err != nil {
	if err := r.db.Where("id in ?", ids).Find(&globalIds).Delete(&globalIds).Error; err != nil {
		return err
	}
	return nil
}

// GetNewID 获取数据库实现的递增唯一ID
func (r *globalIdRepository) GetNewID(identifier string) uint {
	var newID uint
	r.db.Exec("UPDATE global_ids SET last_id=LAST_INSERT_ID(last_id+step) WHERE identifier=?", identifier)
	r.db.Raw("SELECT LAST_INSERT_ID()").Scan(&newID)
	return newID
}
