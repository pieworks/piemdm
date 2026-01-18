package model

// 权限规则
type DelRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"` // 策略类型
	V0    string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"` // 用户 sub
	V1    string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"` // 资源 obj
	V2    string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"` // 动作 act
	V3    string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"` // 效果 eft
	V4    string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"`
	V5    string `gorm:"size:100;uniqueIndex:idx_rule" binding:"max=100"`
	//  PType  string `json:"p_type" gorm:"column:p_type" description:"策略类型"`
	//  RoleId string `json:"role_id" gorm:"column:v0" description:"角色ID"`
	//  Path   string `json:"path" gorm:"column:v1" description:"api路径"`
	//  Method string `json:"method" gorm:"column:v2" description:"访问方法"`
}
