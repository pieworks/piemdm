package model

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey;not null"`
	RoleID uint `gorm:"primaryKey;not null"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
