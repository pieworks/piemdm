package model

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;not null"`
	PermissionID uint `gorm:"primaryKey;not null"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
