package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"piemdm/internal/model"
)

func TestApprovalNode_BeforeCreate(t *testing.T) {
	tests := []struct {
		name    string
		node    *model.ApprovalNode
		wantErr bool
		errMsg  string
	}{
		{
			name: "审批节点-有效的审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: model.ApproverTypeUsers,
			},
			wantErr: false,
		},
		{
			name: "审批节点-AUTO_REJECT类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: model.ApproverTypeAutoReject,
			},
			wantErr: false,
		},
		{
			name: "审批节点-AUTO_APPROVE类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: model.ApproverTypeAutoApprove,
			},
			wantErr: false,
		},
		{
			name: "审批节点-无效的审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: "INVALID_TYPE",
			},
			wantErr: true,
			errMsg:  "应该返回验证错误",
		},
		{
			name: "开始节点-有审批人类型但不验证",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeStart,
				ApproverType: "INVALID_TYPE", // 开始节点不需要审批人，所以不验证
			},
			wantErr: false,
		},
		{
			name: "结束节点-有审批人类型但不验证",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeEnd,
				ApproverType: "INVALID_TYPE", // 结束节点不需要审批人，所以不验证
			},
			wantErr: false,
		},
		{
			name: "条件节点-有审批人类型但不验证",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeCondition,
				ApproverType: "INVALID_TYPE", // 条件节点不需要审批人，所以不验证
			},
			wantErr: false,
		},
		{
			name: "审批节点-空的审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: "", // 空值不验证
			},
			wantErr: false,
		},
		{
			name: "无效的节点类型",
			node: &model.ApprovalNode{
				NodeType:     "INVALID_NODE_TYPE",
				ApproverType: model.ApproverTypeUsers,
			},
			wantErr: true,
			errMsg:  "应该返回节点类型验证错误",
		},
		{
			name: "审批节点-所有有效的审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: model.ApproverTypeRoles,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.BeforeCreate(&gorm.DB{})

			if tt.wantErr {
				assert.Error(t, err, tt.errMsg)
				assert.Equal(t, gorm.ErrInvalidValue, err)
			} else {
				assert.NoError(t, err)
				// 验证默认状态设置
				if tt.node.Status == "" {
					assert.Equal(t, model.ApprovalDefStatusNormal, tt.node.Status)
				}
			}
		})
	}
}

func TestApprovalNode_BeforeUpdate(t *testing.T) {
	tests := []struct {
		name    string
		node    *model.ApprovalNode
		wantErr bool
		errMsg  string
	}{
		{
			name: "审批节点-有效的审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: model.ApproverTypeDepartments,
			},
			wantErr: false,
		},
		{
			name: "审批节点-无效的审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeApproval,
				ApproverType: "INVALID_TYPE",
			},
			wantErr: true,
			errMsg:  "应该返回验证错误",
		},
		{
			name: "非审批节点-不验证审批人类型",
			node: &model.ApprovalNode{
				NodeType:     model.NodeTypeCC,
				ApproverType: "INVALID_TYPE", // 抄送节点不需要审批人，所以不验证
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.BeforeUpdate(&gorm.DB{})

			if tt.wantErr {
				assert.Error(t, err, tt.errMsg)
				assert.Equal(t, gorm.ErrInvalidValue, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidApproverType(t *testing.T) {
	tests := []struct {
		name         string
		approverType string
		expected     bool
	}{
		{"USERS", model.ApproverTypeUsers, true},
		{"ROLES", model.ApproverTypeRoles, true},
		{"DEPARTMENTS", model.ApproverTypeDepartments, true},
		{"POSITIONS", model.ApproverTypePositions, true},
		{"SELF_SELECT", model.ApproverTypeSelfSelect, true},
		{"EXPRESSION", model.ApproverTypeExpression, true},
		{"SUPERIOR", model.ApproverTypeSuperior, true},
		{"DEPT_MANAGER", model.ApproverTypeDeptManager, true},
		{"AUTO_REJECT", model.ApproverTypeAutoReject, true},
		{"AUTO_APPROVE", model.ApproverTypeAutoApprove, true},
		{"INVALID", "INVALID", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.IsValidApproverType(tt.approverType)
			assert.Equal(t, tt.expected, result,
				"IsValidApproverType(%s) = %v, expected %v",
				tt.approverType, result, tt.expected)
		})
	}
}

func TestApprovalNode_NodeTypeCheckers(t *testing.T) {
	tests := []struct {
		name     string
		nodeType string
		checker  func(*model.ApprovalNode) bool
		expected bool
	}{
		{"IsStartNode-START", model.NodeTypeStart, (*model.ApprovalNode).IsStartNode, true},
		{"IsStartNode-APPROVAL", model.NodeTypeApproval, (*model.ApprovalNode).IsStartNode, false},
		{"IsApprovalNode-APPROVAL", model.NodeTypeApproval, (*model.ApprovalNode).IsApprovalNode, true},
		{"IsApprovalNode-START", model.NodeTypeStart, (*model.ApprovalNode).IsApprovalNode, false},
		{"IsEndNode-END", model.NodeTypeEnd, (*model.ApprovalNode).IsEndNode, true},
		{"IsConditionNode-CONDITION", model.NodeTypeCondition, (*model.ApprovalNode).IsConditionNode, true},
		{"IsCCNode-CC", model.NodeTypeCC, (*model.ApprovalNode).IsCCNode, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &model.ApprovalNode{NodeType: tt.nodeType}
			result := tt.checker(node)
			assert.Equal(t, tt.expected, result)
		})
	}
}
