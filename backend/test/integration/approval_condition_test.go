package integration

import (
	"testing"

	"piemdm/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApprovalConditionEvaluation tests the logic for evaluating conditional branches in an approval workflow.
func TestApprovalConditionEvaluation(t *testing.T) {
	// The conditional logic configuration.
	conditionConfigJSON := `{
		"branches": [
			{
				"name": "条件分支 1",
				"condition": {
					"fieldName": "price",
					"operator": "gte",
					"fieldValue": "3000"
				},
				"nodes": [
					{"nodeCode": "finance-approval"}
				]
			},
			{
				"name": "其他情况",
				"condition": {
					"fieldName": "",
					"operator": "eq",
					"fieldValue": ""
				},
				"nodes": [
					{"nodeCode": "auto-reject"}
				]
			}
		]
	}`

	// The node that contains the conditional logic.
	conditionNode := &model.ApprovalNode{
		NodeCode:        "condition-node-1",
		NodeName:        "价格条件节点",
		NodeType:        model.NodeTypeCondition,
		ConditionConfig: conditionConfigJSON,
	}

	// The list of all possible nodes in the workflow.
	approvalNodes := []*model.ApprovalNode{
		conditionNode,
		{NodeCode: "finance-approval", NodeName: "财务审批", NodeType: model.NodeTypeApproval},
		{NodeCode: "auto-reject", NodeName: "自动驳回", NodeType: model.NodeTypeApproval},
	}

	tests := []struct {
		name             string
		formData         map[string]any
		expectedNodeCode string
		description      string
	}{
		{
			name:             "价格大于等于3000-应该走财务审批",
			formData:         map[string]any{"price": 3500},
			expectedNodeCode: "finance-approval",
			description:      "价格3500大于等于3000，应该走财务审批节点",
		},
		{
			name:             "价格等于3000-应该走财务审批",
			formData:         map[string]any{"price": 3000},
			expectedNodeCode: "finance-approval",
			description:      "价格3000等于3000，应该走财务审批节点",
		},
		{
			name:             "价格小于3000-应该走自动驳回",
			formData:         map[string]any{"price": 2500},
			expectedNodeCode: "auto-reject",
			description:      "价格2500小于3000，应该走其他情况（自动驳回）节点",
		},
		{
			name:             "没有价格字段-应该走自动驳回",
			formData:         map[string]any{"name": "测试商品"},
			expectedNodeCode: "auto-reject",
			description:      "没有价格字段，应该走其他情况（自动驳回）节点",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("测试场景: %s", tt.description)

			// 调用条件评估方法
			nextNode, err := evaluateConditionNode(approvalNodes, conditionNode, tt.formData)

			// 验证结果
			require.NoError(t, err, "条件评估不应该返回错误")
			require.NotNil(t, nextNode, "应该返回下一个节点")
			assert.Equal(t, tt.expectedNodeCode, nextNode.NodeCode, "节点代码应该匹配预期: %s", tt.description)

			t.Logf("✅ 测试通过: 选择了节点 %s (%s)", nextNode.NodeCode, nextNode.NodeName)
		})
	}
}
