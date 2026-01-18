package service

import (
	"context"
	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"

	"github.com/pkg/errors"
)

// PermissionService 权限服务接口
type PermissionService interface {
	// Base CRUD
	Create(ctx context.Context, permission *model.Permission) error
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
	GetAll(ctx context.Context) ([]*model.Permission, error)

	// 树状结构
	GetTree(ctx context.Context) ([]*model.Permission, error)
	GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error)

	// 分页查询
	List(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error)

	// 权限验证
	HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error)
}

type permissionService struct {
	repo   repository.PermissionRepository
	logger *log.Logger
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(repo repository.PermissionRepository, logger *log.Logger) PermissionService {
	return &permissionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *permissionService) Create(ctx context.Context, permission *model.Permission) error {
	// 检查 code 是否已存在
	existing, err := s.repo.FindByCode(ctx, permission.Code)
	if err == nil && existing != nil {
		return errors.New("权限代码已存在")
	}

	return s.repo.Create(ctx, permission)
}

func (s *permissionService) Update(ctx context.Context, permission *model.Permission) error {
	// 检查权限是否存在
	_, err := s.repo.FindByID(ctx, permission.ID)
	if err != nil {
		return errors.Wrap(err, "权限不存在")
	}

	// 检查 code 是否与其他权限冲突
	existing, err := s.repo.FindByCode(ctx, permission.Code)
	if err == nil && existing != nil && existing.ID != permission.ID {
		return errors.New("权限代码已被其他权限使用")
	}

	return s.repo.Update(ctx, permission)
}

func (s *permissionService) Delete(ctx context.Context, id uint) error {
	// 检查是否有子权限
	children, err := s.repo.FindChildren(ctx, id)
	if err != nil {
		return errors.Wrap(err, "查询子权限失败")
	}

	if len(children) > 0 {
		return errors.New("该权限下还有子权限,无法删除")
	}

	return s.repo.Delete(ctx, id)
}

func (s *permissionService) Get(ctx context.Context, id uint) (*model.Permission, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *permissionService) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	return s.repo.FindByCode(ctx, code)
}

func (s *permissionService) GetAll(ctx context.Context) ([]*model.Permission, error) {
	return s.repo.FindAll(ctx)
}

func (s *permissionService) GetTree(ctx context.Context) ([]*model.Permission, error) {
	return s.repo.FindTree(ctx)
}

func (s *permissionService) GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error) {
	return s.repo.FindChildren(ctx, parentID)
}

func (s *permissionService) List(ctx context.Context, page, pageSize int) ([]*model.Permission, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.FindPage(ctx, page, pageSize)
}

func (s *permissionService) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	// TODO: 实现用户权限检查逻辑
	// 1. 获取用户的所有角色
	// 2. 获取这些角色的所有权限
	// 3. 检查是否包含指定的权限代码

	// 这部分需要在实现 Role Service 后完成
	s.logger.Warn("HasPermission 方法尚未完全实现")
	return false, errors.New("功能尚未实现")
}
