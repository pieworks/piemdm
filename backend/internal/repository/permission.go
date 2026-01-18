package repository

import (
	"context"
	"piemdm/internal/model"
	"piemdm/pkg/log"

	"gorm.io/gorm"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	// Base CRUD
	Create(ctx context.Context, permission *model.Permission) error
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Permission, error)
	FindByCode(ctx context.Context, code string) (*model.Permission, error)
	FindAll(ctx context.Context) ([]*model.Permission, error)

	// 树状结构查询
	FindTree(ctx context.Context) ([]*model.Permission, error)
	FindChildren(ctx context.Context, parentID uint) ([]*model.Permission, error)

	// 分页查询
	FindPage(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error)
}

type permissionRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

// NewPermissionRepository 创建权限仓储实例
func NewPermissionRepository(db *gorm.DB, logger *log.Logger) PermissionRepository {
	return &permissionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

func (r *permissionRepository) FindByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Order("parent_id, id").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindTree(ctx context.Context) ([]*model.Permission, error) {
	// 获取所有权限
	var allPermissions []*model.Permission
	err := r.db.WithContext(ctx).Order("parent_id, id").Find(&allPermissions).Error
	if err != nil {
		return nil, err
	}

	// 构建树状结构
	permissionMap := make(map[uint]*model.Permission)
	var rootPermissions []*model.Permission

	// 第一遍遍历,建立 ID 到权限的映射
	for _, p := range allPermissions {
		permissionMap[p.ID] = p
		p.Children = []*model.Permission{} // 初始化子节点切片
	}

	// 第二遍遍历,构建父子关系
	for _, p := range allPermissions {
		if p.ParentID == 0 {
			// 根节点
			rootPermissions = append(rootPermissions, p)
		} else {
			// 子节点,添加到父节点的 Children 中
			if parent, exists := permissionMap[p.ParentID]; exists {
				parent.Children = append(parent.Children, p)
			}
		}
	}

	return rootPermissions, nil
}

func (r *permissionRepository) FindChildren(ctx context.Context, parentID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("id").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindPage(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&model.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("parent_id, id").
		Find(&permissions).Error

	return permissions, total, err
}
