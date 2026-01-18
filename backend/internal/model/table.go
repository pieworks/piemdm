package model

import (
	"time"

	"gorm.io/gorm"
)

type Table struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"size:64;unique;not null" binding:"required,max=64,containsany=abcdefghijklmnopqrstuvwxyz0123456789_,lowercase" label:"表名"` // 表名,英文
	// 名称，可以中文
	Name string `gorm:"size:128;not null" binding:"required,max=128"`
	// 展示模式：List 列表, Tree 树形
	DisplayMode string `gorm:"size:16;default:List;" json:"DisplayMode" validate:"oneof=List Tree"`
	// 表类型：Entity 主表实体, Item 行项目
	TableType string `gorm:"size:8;default:Entity" json:"TableType" validate:"oneof=Entity Item"`
	// 父表Code（仅当TableType=Item时有效，不带t_前缀）
	ParentTable string `gorm:"size:64;" json:"ParentTable" binding:"max=64"`
	// 父表关联字段
	ParentField string `gorm:"size:64;" json:"ParentField" binding:"max=64"`
	// 本表关联字段
	SelfField   string `gorm:"size:64;" json:"SelfField" binding:"max=64"`
	Sort        uint   `gorm:"size:10;default:0" binding:"max=9999"`
	Description string `gorm:"size:255" binding:"max=255"`
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	Status string `gorm:"size:8;default:Normal"`
	// 这些都是服务端参数，不用校验。
	CreatedBy string `gorm:"size:64"` // 创建人
	UpdatedBy string `gorm:"size:64"` // 更新人
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Table) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Model(m).Where("id = ?", m.ID).Update("status", "Deleted")
	return
}

// GORM 钩子：自动设置操作人
func (m *Table) BeforeCreate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.CreatedBy = user
	}
	return
}

func (m *Table) BeforeUpdate(tx *gorm.DB) (err error) {
	if user, ok := tx.Statement.Context.Value("user_name").(string); ok {
		m.UpdatedBy = user
	}

	return
}

// //BeforeCreate 在创建Article之前，先把创建时间赋值
// func (st *Table) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("CreatedOn", time.Now().Unix())
// 	return nil
// }

// //BeforeUpdate 在更新Article之前，先把更新时间赋值
// func (st *Table) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	return nil
// }

// func (st *Table) BeforeCreate(tx *gorm.DB) (err error) {

// 	if st.Type == 3 {
// 		err = errors.New("can't save invalid data")
// 	}
// 	return
// }

// func (st *Table) AfterCreate(tx *gorm.DB) (err error) {
// 	if st.ID == 1 {
// 		// tx.Model(st).Update("role", "admin")
// 	}
// 	return
// }

// func (m *Table) BeforeDelete(tx *gorm.DB) (err error) {
// 	// m.Status = "Deleted"
// 	// tx.Statement.SetColumn("status", "Deleted")
// 	tx.Model(m).Where("id = ?", m.ID).Update("status", "Deleted")
// 	// mm := tx.Statement.Dest.(map[string]any)

// 	return
// }

// `id` int(11) NOT NULL AUTO_INCREMENT,
// `code` varchar(32) DEFAULT NULL COMMENT '模型英文名称code',
// `name` varchar(32) DEFAULT NULL COMMENT '模型名称',
// `type` varchar(2) DEFAULT NULL COMMENT '模型类型。0,列表;1,父子树;2,分类树 3.附属模型 4.关系模型',
// `parent_id` varchar(32) DEFAULT NULL COMMENT '从表所属id',
// `view` varchar(2) DEFAULT NULL COMMENT '是否多视图 0-是 1-否',
// `view_name` varchar(200) DEFAULT NULL COMMENT '视图名称',
// `group` varchar(2) DEFAULT NULL COMMENT '是否分组。0,是;1,否',
// `group_name` varchar(200) DEFAULT NULL COMMENT '分组名称',
// `status` varchar(2) DEFAULT '1' COMMENT '模型状态。0,活动;1,冻结 2-逻辑删除',
// `created_by` varchar(10) DEFAULT NULL COMMENT '创建人',
// `field` varchar(100) DEFAULT NULL COMMENT '模型存到mongo中的conllection 名称',
// `order` int(10) DEFAULT NULL COMMENT '模型的添加顺序标识',
// `desc` varchar(255) DEFAULT NULL COMMENT '模型描述',
// `category` varchar(2) DEFAULT '0' COMMENT '区分分类（0，模型；1，字典表；2；内置页）',
// `page` varchar(1000) DEFAULT NULL,
// `max_level` varchar(32) DEFAULT NULL COMMENT '最大层级（父子树）',
// `page_whether` varchar(2) DEFAULT '0' COMMENT '配置是否需要附属模型0：纯附属模型1：主数据有附属模型；3：历史模型',
// `model_level` varchar(2) DEFAULT '1' COMMENT '附属模型的级别：1.一级 2. 二级',
// `fsmodel_id` varchar(32) DEFAULT NULL COMMENT '二级附属模型：所关联的一级附属模型的modelId',
// `fsmodel_dttrId` varchar(1000) DEFAULT NULL COMMENT '二级附属模型关联一级附属模型的字段',
