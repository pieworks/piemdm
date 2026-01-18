package repository

import (
	"piemdm/internal/model"

	"github.com/gin-gonic/gin"
)

type ApplicationEntityRepository interface {
	// 根据 AppId 和 EntityCode 查询
	FindByAppIdAndEntityCode(appId, entityCode string) (*model.ApplicationEntity, error)

	// 基础查询
	FindOne(id uint) (*model.ApplicationEntity, error)
	FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.ApplicationEntity, error)

	// Base CRUD
	Create(c *gin.Context, entity *model.ApplicationEntity) error
	Update(c *gin.Context, entity *model.ApplicationEntity) error
	Delete(c *gin.Context, id uint) error

	// Batch operations
	BatchUpdate(c *gin.Context, ids []uint, entity *model.ApplicationEntity) error
	BatchDelete(c *gin.Context, ids []uint) error
}

type applicationEntityRepository struct {
	*Repository
	source Base
}

func NewApplicationEntityRepository(repository *Repository, source Base) ApplicationEntityRepository {
	return &applicationEntityRepository{
		Repository: repository,
		source:     source,
	}
}

func (r *applicationEntityRepository) FindByAppIdAndEntityCode(appId, entityCode string) (*model.ApplicationEntity, error) {
	var entity model.ApplicationEntity
	if err := r.db.Where("app_id = ? AND entity_code = ? AND status = ?", appId, entityCode, "Normal").First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *applicationEntityRepository) FindOne(id uint) (*model.ApplicationEntity, error) {
	var entity model.ApplicationEntity
	if err := r.source.FirstById(&entity, id); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *applicationEntityRepository) FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.ApplicationEntity, error) {
	var entities []*model.ApplicationEntity
	var entity model.ApplicationEntity

	preloads := []string{}
	err := r.source.FindPage(entity, &entities, page, pageSize, total, where, preloads, "ID desc")
	if err != nil {
		r.logger.Error("获取 ApplicationEntity 失败", "err", err)
		return nil, err
	}

	return entities, nil
}

func (r *applicationEntityRepository) Create(c *gin.Context, entity *model.ApplicationEntity) error {
	if err := r.db.WithContext(c).Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *applicationEntityRepository) Update(c *gin.Context, entity *model.ApplicationEntity) error {
	if err := r.db.WithContext(c).Updates(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *applicationEntityRepository) BatchUpdate(c *gin.Context, ids []uint, entity *model.ApplicationEntity) error {
	if err := r.db.WithContext(c).Model(&entity).Where("id in ?", ids).Updates(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *applicationEntityRepository) Delete(c *gin.Context, id uint) error {
	if err := r.db.WithContext(c).Where("id = ?", id).Delete(&model.ApplicationEntity{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *applicationEntityRepository) BatchDelete(c *gin.Context, ids []uint) error {
	var entities []model.ApplicationEntity
	if err := r.db.WithContext(c).Where("id in ?", ids).Find(&entities).Delete(&entities).Error; err != nil {
		return err
	}
	return nil
}
