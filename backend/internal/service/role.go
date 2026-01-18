package service

import (
	"context"
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/casbin/casbin/v2"
)

type RoleService interface {
	Get(id uint) (*model.Role, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.Role, error)
	Find(sel string, where map[string]any) ([]*model.Role, error)
	Create(c context.Context, role *model.Role) error
	Update(c context.Context, role *model.Role) error
	BatchUpdate(c context.Context, ids []uint, role *model.Role) error
	Delete(c context.Context, id uint) (*model.Role, error)
	BatchDelete(c context.Context, ids []uint) error

	// Permission management
	GetRolePermissions(c context.Context, roleID uint) ([]*model.Permission, error)
	AssignPermissions(c context.Context, roleID uint, permissionIDs []uint) error
	RemovePermissions(c context.Context, roleID uint, permissionIDs []uint) error
	UpdatePermissions(c context.Context, roleID uint, permissionIDs []uint) error

	// User management
	GetRoleUsers(roleID uint) ([]model.User, error)
	UpdateRoleUsers(roleID uint, userIDs []uint) error
}

type roleService struct {
	*Service
	roleRepository repository.RoleRepository
	userRoleRepo   repository.UserRoleRepository
	enforcer       *casbin.Enforcer
}

func NewRoleService(service *Service, roleRepository repository.RoleRepository, userRoleRepo repository.UserRoleRepository, enforcer *casbin.Enforcer) RoleService {
	return &roleService{
		Service:        service,
		roleRepository: roleRepository,
		userRoleRepo:   userRoleRepo,
		enforcer:       enforcer,
	}
}

func (s *roleService) Get(id uint) (*model.Role, error) {
	return s.roleRepository.FindOne(id)
}

func (s *roleService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.Role, error) {
	return s.roleRepository.FindPage(page, pageSize, total, where)
}

func (s *roleService) Find(sel string, where map[string]any) ([]*model.Role, error) {
	return s.roleRepository.Find(sel, where)
}

func (s *roleService) Create(c context.Context, role *model.Role) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// cron.Code = code

	return s.roleRepository.Create(c, role)
}

func (s *roleService) Update(c context.Context, role *model.Role) error {
	return s.roleRepository.Update(c, role)
}

func (s *roleService) BatchUpdate(c context.Context, ids []uint, role *model.Role) error {
	return s.roleRepository.BatchUpdate(c, ids, role)
}

func (s *roleService) Delete(c context.Context, id uint) (*model.Role, error) {
	return s.roleRepository.Delete(c, id)
}

func (s *roleService) BatchDelete(c context.Context, ids []uint) error {
	return s.roleRepository.BatchDelete(c, ids)
}

func (s *roleService) GetRolePermissions(c context.Context, roleID uint) ([]*model.Permission, error) {
	return s.roleRepository.GetPermissions(c, roleID)
}

func (s *roleService) AssignPermissions(c context.Context, roleID uint, permissionIDs []uint) error {
	if err := s.roleRepository.AssignPermissions(c, roleID, permissionIDs); err != nil {
		return err
	}
	// Reload policy
	return s.enforcer.LoadPolicy()
}

func (s *roleService) RemovePermissions(c context.Context, roleID uint, permissionIDs []uint) error {
	if err := s.roleRepository.RemovePermissions(c, roleID, permissionIDs); err != nil {
		return err
	}
	return s.enforcer.LoadPolicy()
}

func (s *roleService) UpdatePermissions(c context.Context, roleID uint, permissionIDs []uint) error {
	if err := s.roleRepository.UpdatePermissions(c, roleID, permissionIDs); err != nil {
		return err
	}
	return s.enforcer.LoadPolicy()
}

// GetRoleUsers 获取角色的用户
func (s *roleService) GetRoleUsers(roleID uint) ([]model.User, error) {
	return s.userRoleRepo.FindUsersByRoleID(roleID)
}

// UpdateRoleUsers 更新角色的用户
func (s *roleService) UpdateRoleUsers(roleID uint, userIDs []uint) error {
	return s.userRoleRepo.UpdateRoleUsers(roleID, userIDs)
}
