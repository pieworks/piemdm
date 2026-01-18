package model

type DelUserRepo struct {
	ID          uint   `gorm:"primarykey"`
	Username    string `gorm:"size:128"`
	DisplayName string `gorm:"size:128"`
	// Status   int8   `gorm:"size:2;default:0;omitempty" `
}

// func (*DelUserRepo) TableName() string {
// 	return "users"
// }
