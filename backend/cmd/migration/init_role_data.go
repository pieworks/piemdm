package main

import (
	"context"
	"fmt"
	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitRoleData 初始化角色数据
type InitRoleData struct {
	db             *gorm.DB
	logger         *log.Logger
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	userRepo       repository.UserRepository
}

func NewInitRoleData(db *gorm.DB, logger *log.Logger) *InitRoleData {
	// Initialize repositories
	// Note: RoleRepository needs a Repository and Base. Since we are in main, we might need to verify how these are usually constructed.
	// Looking at migration.go, it uses simple struct initialization if easier, or we can construct properly.
	// For simplicity in migration script, we can use the constructors if available or manual struct init.
	// But `repository.NewRoleRepository` takes `*repository.Repository` and `repository.Base`.

	// Create a base repository wrapper
	baseRepo := repository.NewRepository(db, nil, logger) // No cache for migration
	baseDB := repository.NewBaseRepository(baseRepo)

	return &InitRoleData{
		db:             db,
		logger:         logger,
		roleRepo:       repository.NewRoleRepository(baseRepo, baseDB),
		permissionRepo: repository.NewPermissionRepository(db, logger), // PermissionRepository seems to take db and logger directly
		userRepo:       repository.NewUserRepository(baseRepo, baseDB),
	}
}

func (i *InitRoleData) Run() error {
	ctx := context.Background()
	i.logger.Info("开始初始化角色数据...")

	// 1. 定义角色
	roles := []struct {
		Code        string
		Name        string
		Description string
		IsSuper     bool // 标记是否为超级管理员(拥有所有权限)
		DataScope   string
	}{
		{Code: "superuser", Name: "超级管理员", Description: "拥有系统所有权限", IsSuper: true, DataScope: "All"},
		{Code: "admin", Name: "系统管理员", Description: "拥有大部分系统权限", IsSuper: false, DataScope: "All"},
		{Code: "user", Name: "普通用户", Description: "基础权限", IsSuper: false, DataScope: "Self"},
	}

	for _, r := range roles {
		// 检查角色是否存在
		existingRole, _ := i.roleRepo.Find("*", map[string]any{"code": r.Code})
		var role *model.Role

		if len(existingRole) > 0 {
			i.logger.Info(fmt.Sprintf("角色 %s 已存在", r.Code))
			role = existingRole[0]

			// 更新 DataScope
			if role.DataScope != r.DataScope {
				role.DataScope = r.DataScope
				i.logger.Info(fmt.Sprintf("更新角色 %s 的 DataScope 为 %s", r.Code, r.DataScope))
				// 使用 Repository 的 Update 方法 (假设有) 或直接用 db
				i.db.Save(role)
			}

		} else {
			// 创建角色
			newRole := &model.Role{
				Code:        r.Code,
				Name:        r.Name,
				Description: r.Description,
				Status:      "Normal",
				DataScope:   r.DataScope,
			}
			if err := i.roleRepo.Create(ctx, newRole); err != nil { // use ctx (context.Background)
				i.logger.Error(fmt.Sprintf("创建角色 %s 失败: %v", r.Code, err))
				return err
			}
			i.logger.Info(fmt.Sprintf("创建角色成功: %s", r.Code))
			role = newRole
		}

		// 2. 分配权限
		if err := i.assignPermissions(ctx, role, r.IsSuper); err != nil {
			i.logger.Error(fmt.Sprintf("为角色 %s 分配权限失败: %v", r.Code, err))
			// 继续处理下一个角色
		}
	}

	i.logger.Info("角色数据初始化完成")

	// 3. 为默认 admin 用户分配超级管理员角色
	var adminUser model.User
	if err := i.db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		i.logger.Info("用户 admin 不存在, 开始创建...")
		hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		adminUser = model.User{
			EmployeeID:       "admin",
			Username:         "admin",
			Password:         string(hash),
			FirstName:        "Admin",
			LastName:         "System",
			DisplayName:      "Administrator",
			Email:            "admin@example.com",
			Admin:            "Yes",
			Status:           "Normal",
			SuperiorUsername: "root",
		}
		if err := i.db.Create(&adminUser).Error; err != nil {
			i.logger.Error(fmt.Sprintf("创建 admin 用户失败: %v", err))
			return err
		}
		i.logger.Info("创建 admin 用户成功")
	}

	var totalRole int64
	i.db.Model(&model.UserRole{}).Where("user_id = ?", adminUser.ID).Count(&totalRole)
	if totalRole == 0 {
		var superRole model.Role
		if err := i.db.Where("code = ?", "superuser").First(&superRole).Error; err == nil {
			userRole := model.UserRole{
				UserID: adminUser.ID,
				RoleID: superRole.ID,
			}
			if err := i.db.Create(&userRole).Error; err != nil {
				i.logger.Error(fmt.Sprintf("为 admin 用户分配角色失败: %v", err))
				return err
			}
			i.logger.Info("已为 admin 用户分配 superuser 角色")
		}
	}

	return nil
}

func (i *InitRoleData) assignPermissions(ctx context.Context, role *model.Role, isSuper bool) error {
	// 获取所有权限
	allPermissions, _, err := i.permissionRepo.FindPage(ctx, 1, 1000)
	if err != nil {
		return err
	}

	var permissionIDs []uint

	if isSuper {
		// 超级管理员：拥有所有权限
		for _, p := range allPermissions {
			permissionIDs = append(permissionIDs, p.ID)
		}
	} else if role.Code == "admin" {
		// 管理员：排除部分敏感权限
		for _, p := range allPermissions {
			// 排除权限管理和角色分配权限
			if p.Code == "permission:delete" || p.Code == "role:assign_permission" {
				continue
			}
			permissionIDs = append(permissionIDs, p.ID)
		}
	} else {
		// 普通用户：基础权限
		for _, p := range allPermissions {
			// 只给查看权限
			if p.Action == "list" || p.Action == "read" {
				permissionIDs = append(permissionIDs, p.ID)
			}
		}
	}

	if len(permissionIDs) == 0 {
		return nil
	}

	// 更新角色权限 (先这样简单处理，覆盖更新)
	// RoleRepository added AssignPermissions (Append) and UpdatePermissions (Replace).
	// We should probably use UpdatePermissions to ensure state matches.
	if err := i.roleRepo.UpdatePermissions(ctx, role.ID, permissionIDs); err != nil { // use ctx
		return err
	}

	i.logger.Info(fmt.Sprintf("已为角色 %s 分配 %d 个权限", role.Code, len(permissionIDs)))
	return nil
}
