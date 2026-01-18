package repository

import (
	"piemdm/internal/model"

	"gorm.io/gorm"
)

type UserRoleRepository interface {
	// FindRolesByUserID 获取用户的所有角色
	FindRolesByUserID(userID uint) ([]model.Role, error)
	// FindUsersByRoleID 获取角色的所有用户
	FindUsersByRoleID(roleID uint) ([]model.User, error)
	// UpdateUserRoles 更新用户的角色(覆盖)
	UpdateUserRoles(userID uint, roleIDs []uint) error
	// UpdateRoleUsers 更新角色的用户(覆盖)
	UpdateRoleUsers(roleID uint, userIDs []uint) error
}

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

// FindRolesByUserID 获取用户的所有角色
func (r *userRoleRepository) FindRolesByUserID(userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// FindUsersByRoleID 获取角色的所有用户
func (r *userRoleRepository) FindUsersByRoleID(roleID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users).Error
	return users, err
}

// UpdateUserRoles 更新用户的角色(覆盖)
func (r *userRoleRepository) UpdateUserRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除用户的所有角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 2. 如果roleIDs为空,直接返回
		if len(roleIDs) == 0 {
			return nil
		}

		// 3. 批量插入新的角色关联
		userRoles := make([]model.UserRole, len(roleIDs))
		for i, roleID := range roleIDs {
			userRoles[i] = model.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
		}
		return tx.Create(&userRoles).Error
	})
}

// UpdateRoleUsers 更新角色的用户(覆盖)
func (r *userRoleRepository) UpdateRoleUsers(roleID uint, userIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除角色的所有用户关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 2. 如果userIDs为空,直接返回
		if len(userIDs) == 0 {
			return nil
		}

		// 3. 批量插入新的用户关联
		userRoles := make([]model.UserRole, len(userIDs))
		for i, userID := range userIDs {
			userRoles[i] = model.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
		}
		return tx.Create(&userRoles).Error
	})
}
