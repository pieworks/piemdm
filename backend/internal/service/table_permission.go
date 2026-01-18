package service

import (
	"errors"
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TablePermissionService interface {
	Grant(ctx *gin.Context, roleId uint, tableCode string) error

	// Valid Replacement:
	BatchUpdate(ctx *gin.Context, ids []uint, values map[string]any) error
	BatchDelete(ctx *gin.Context, ids []uint) error

	GetByRoleID(ctx *gin.Context, roleId uint) ([]*model.TablePermission, error)
	HasPermission(ctx *gin.Context, userId uint, tableCode string) (bool, error)
	GetAllowedTableCodes(ctx *gin.Context, userId uint) ([]string, error)
	CheckTablePermission(ctx *gin.Context, userId uint, tableCode string) (bool, error)
}

type tablePermissionService struct {
	repo         repository.TablePermissionRepository
	tableRepo    repository.TableRepository
	userRoleRepo repository.UserRoleRepository
}

func NewTablePermissionService(
	repo repository.TablePermissionRepository,
	tableRepo repository.TableRepository,
	userRoleRepo repository.UserRoleRepository,
) TablePermissionService {
	return &tablePermissionService{
		repo:         repo,
		tableRepo:    tableRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (s *tablePermissionService) Grant(ctx *gin.Context, roleId uint, tableCode string) error {
	// Check if exists
	existing, err := s.repo.FindOne(map[string]any{
		"role_id":    roleId,
		"table_code": tableCode,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		if existing.Status != "Normal" {
			// Reuse UpdateByIds for single update
			return s.repo.UpdateByIds([]uint{existing.ID}, map[string]any{"status": "Normal"})
		}
		return nil
	}

	perm := &model.TablePermission{
		RoleID:    roleId,
		TableCode: tableCode,
		Status:    "Normal",
	}
	// Context user handling
	if user, exists := ctx.Get("user_name"); exists {
		perm.CreatedBy = user.(string)
		perm.UpdatedBy = user.(string)
	}

	return s.repo.Create(perm)
}

func (s *tablePermissionService) BatchUpdate(ctx *gin.Context, ids []uint, values map[string]any) error {
	// Should filter allowed fields to update? e.g. status
	// For now, allow all map values passed from handler (handler should filter)
	return s.repo.UpdateByIds(ids, values)
}

func (s *tablePermissionService) BatchDelete(ctx *gin.Context, ids []uint) error {
	return s.repo.DeleteByIds(ids)
}

func (s *tablePermissionService) GetByRoleID(ctx *gin.Context, roleId uint) ([]*model.TablePermission, error) {
	return s.repo.FindByRoleID(roleId)
}

func (s *tablePermissionService) HasPermission(ctx *gin.Context, userId uint, tableCode string) (bool, error) {
	// 1. Get user roles
	roles, err := s.userRoleRepo.FindRolesByUserID(userId)
	if err != nil {
		return false, err
	}
	if len(roles) == 0 {
		return false, nil
	}

	roleIds := make([]uint, len(roles))
	for i, r := range roles {
		roleIds[i] = r.ID
	}

	// 2. Check if any role has permission
	// GORM Where with slice performs IN query
	perms, err := s.repo.Find(map[string]any{
		"role_id":    roleIds,
		"table_code": tableCode,
	})

	if err != nil {
		return false, err
	}

	// Check if any active permission exists
	for _, p := range perms {
		if p.Status == "Normal" {
			return true, nil
		}
	}

	return false, nil
}

// GetAllowedTableCodes returns a list of table codes that the user has permission to access.
// Returns nil if the user is a superuser or admin (meaning all tables are allowed).
func (s *tablePermissionService) GetAllowedTableCodes(ctx *gin.Context, userId uint) ([]string, error) {
	// 1. Get user roles
	roles, err := s.userRoleRepo.FindRolesByUserID(userId)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return []string{}, nil // No roles = no permissions
	}

	// 2. Check if user is superuser or admin with DataScope="All"
	for _, role := range roles {
		if role.Code == "superuser" || (role.Code == "admin" && role.DataScope == "All") {
			return nil, nil // nil means all tables are allowed
		}
	}

	// 3. Get all table permissions for user's roles
	roleIds := make([]uint, len(roles))
	for i, r := range roles {
		roleIds[i] = r.ID
	}

	perms, err := s.repo.Find(map[string]any{
		"role_id": roleIds,
		"status":  "Normal",
	})
	if err != nil {
		return nil, err
	}

	// 4. Extract unique table codes
	tableCodeMap := make(map[string]bool)
	for _, p := range perms {
		tableCodeMap[p.TableCode] = true
	}

	tableCodes := make([]string, 0, len(tableCodeMap))
	for code := range tableCodeMap {
		tableCodes = append(tableCodes, code)
	}

	return tableCodes, nil
}

// CheckTablePermission checks if a user has permission to access a table.
// For _draft and _log tables, it checks the base table permission.
// Returns true if user has permission, false otherwise.
func (s *tablePermissionService) CheckTablePermission(ctx *gin.Context, userId uint, tableCode string) (bool, error) {
	allowedTableCodes, err := s.GetAllowedTableCodes(ctx, userId)
	if err != nil {
		return false, err
	}

	// If nil, user is superuser/admin (all tables allowed)
	if allowedTableCodes == nil {
		return true, nil
	}

	// Check direct permission
	for _, code := range allowedTableCodes {
		if code == tableCode {
			return true, nil
		}
	}

	// Check if this is a _draft or _log table
	// If so, check permission for the base table
	baseTableCode := tableCode
	if len(tableCode) > 6 && tableCode[len(tableCode)-6:] == "_draft" {
		baseTableCode = tableCode[:len(tableCode)-6]
	} else if len(tableCode) > 4 && tableCode[len(tableCode)-4:] == "_log" {
		baseTableCode = tableCode[:len(tableCode)-4]
	}

	// If we stripped a suffix, check base table permission
	if baseTableCode != tableCode {
		for _, code := range allowedTableCodes {
			if code == baseTableCode {
				return true, nil
			}
		}
	}

	return false, nil
}
