package repository

import (
	"piemdm/internal/model"

	"gorm.io/gorm"
)

type TablePermissionRepository interface {
	// Base CRUD
	Create(permission *model.TablePermission) error
	Delete(id uint) error
	// DeleteBy(userId uint, tableCode string) error // Remove in favor of Find then Delete or generic Delete
	Find(where map[string]any) ([]*model.TablePermission, error)
	FindOne(where map[string]any) (*model.TablePermission, error) // Renamed from First

	// Batch operations
	UpdateByIds(ids []uint, values map[string]any) error // Generic batch update
	DeleteByIds(ids []uint) error

	// 扩展查询
	FindByRoleID(roleId uint) ([]*model.TablePermission, error)
}

type tablePermissionRepository struct {
	db *gorm.DB
}

func NewTablePermissionRepository(db *gorm.DB) TablePermissionRepository {
	return &tablePermissionRepository{db: db}
}

func (r *tablePermissionRepository) Create(permission *model.TablePermission) error {
	return r.db.Create(permission).Error
}

func (r *tablePermissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.TablePermission{}, id).Error
}

// func (r *tablePermissionRepository) DeleteBy(userId uint, tableCode string) error {
// 	return r.db.Where("user_id = ? AND table_code = ?", userId, tableCode).Delete(&model.TablePermission{}).Error
// }

func (r *tablePermissionRepository) Find(where map[string]any) ([]*model.TablePermission, error) {
	var permissions []*model.TablePermission
	err := r.db.Where(where).Find(&permissions).Error
	return permissions, err
}

func (r *tablePermissionRepository) FindByRoleID(roleId uint) ([]*model.TablePermission, error) {
	var permissions []*model.TablePermission
	err := r.db.Where("role_id = ?", roleId).Find(&permissions).Error
	return permissions, err
}

func (r *tablePermissionRepository) FindOne(where map[string]any) (*model.TablePermission, error) {
	var permission model.TablePermission
	err := r.db.Where(where).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *tablePermissionRepository) UpdateByIds(ids []uint, values map[string]any) error {
	return r.db.Model(&model.TablePermission{}).Where("id IN ?", ids).Updates(values).Error
}

func (r *tablePermissionRepository) DeleteByIds(ids []uint) error {
	return r.db.Delete(&model.TablePermission{}, ids).Error
}
