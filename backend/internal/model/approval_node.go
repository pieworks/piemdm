package model

import (
	"time"

	"gorm.io/gorm"
)

// ApprovalNode 审批节点模型
type ApprovalNode struct {
	ID              uint   `gorm:"primarykey"`
	ApprovalDefCode string `gorm:"size:128;index" binding:"required,max=128"` // 审批定义编码
	NodeCode        string `gorm:"size:128;index" binding:"max=128"`          // 节点编码
	NodeName        string `gorm:"size:128" binding:"max=128"`                // 节点名称
	// 节点类型：START/APPROVAL/CONDITION/CC/END/PARALLEL/MERGE
	NodeType    string `gorm:"size:16" binding:"max=16"`
	Description string `gorm:"size:500"`  // 节点描述，审批要点
	SortOrder   int    `gorm:"default:0"` // 排序

	// 审批配置
	ApproverType    string `gorm:"size:32" binding:"max=32"` // 审批类型
	ApproverConfig  string `gorm:"type:text"`                // 审批人配置JSON
	ConditionConfig string `gorm:"type:text"`                // 条件配置JSON

	// 状态:Normal 正常, Frozen 已冻结, Deleted 已删除
	Status    string         `gorm:"size:8;default:Normal"`
	CreatedBy string         `gorm:"size:64" json:",omitempty"` // 创建人
	UpdatedBy string         `gorm:"size:64" json:",omitempty"` // 更新人
	CreatedAt time.Time      `json:",omitempty"`                // 创建时间
	UpdatedAt time.Time      `json:",omitempty"`                // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:",omitempty"`   // 删除时间
}

// TableName 指定表名
// func (ApprovalNode) TableName() string {
// 	return "approval_node"
// }

// BeforeCreate 创建前钩子
func (m *ApprovalNode) BeforeCreate(tx *gorm.DB) (err error) {
	// 设置默认状态
	if m.Status == "" {
		m.Status = ApprovalDefStatusNormal
	}

	// 验证节点类型
	if m.NodeType != "" && !IsValidNodeType(m.NodeType) {
		return gorm.ErrInvalidValue
	}

	// // 验证审批模式
	// if m.ApprovalMode != "" && !IsValidApprovalMode(m.ApprovalMode) {
	// 	return gorm.ErrInvalidValue
	// }

	// 验证审批人类型 - 只对审批节点进行验证
	if m.NodeType == NodeTypeApproval && m.ApproverType != "" && !IsValidApproverType(m.ApproverType) {
		return gorm.ErrInvalidValue
	}

	return nil
}

// BeforeUpdate 更新前钩子
func (m *ApprovalNode) BeforeUpdate(tx *gorm.DB) (err error) {
	// 验证节点类型
	if m.NodeType != "" && !IsValidNodeType(m.NodeType) {
		return gorm.ErrInvalidValue
	}

	// // 验证审批模式
	// if m.ApprovalMode != "" && !IsValidApprovalMode(m.ApprovalMode) {
	// 	return gorm.ErrInvalidValue
	// }

	// 验证审批人类型 - 只对审批节点进行验证
	if m.NodeType == NodeTypeApproval && m.ApproverType != "" && !IsValidApproverType(m.ApproverType) {
		return gorm.ErrInvalidValue
	}

	return nil
}

// IsActive 检查是否为激活状态
func (m *ApprovalNode) IsActive() bool {
	return m.Status == ApprovalDefStatusNormal
}

// IsStartNode 检查是否为开始节点
func (m *ApprovalNode) IsStartNode() bool {
	return m.NodeType == NodeTypeStart
}

// IsEndNode 检查是否为结束节点
func (m *ApprovalNode) IsEndNode() bool {
	return m.NodeType == NodeTypeEnd
}

// IsApprovalNode 检查是否为审批节点
func (m *ApprovalNode) IsApprovalNode() bool {
	return m.NodeType == NodeTypeApproval
}

// IsConditionNode 检查是否为条件节点
func (m *ApprovalNode) IsConditionNode() bool {
	return m.NodeType == NodeTypeCondition
}

// IsCCNode 检查是否为抄送节点
func (m *ApprovalNode) IsCCNode() bool {
	return m.NodeType == NodeTypeCC
}

// // IsParallelNode 检查是否为并行节点
// func (m *ApprovalNode) IsParallelNode() bool {
// 	return m.NodeType == NodeTypeParallel
// }

// // IsMergeNode 检查是否为合并节点
// func (m *ApprovalNode) IsMergeNode() bool {
// 	return m.NodeType == NodeTypeMerge
// }
