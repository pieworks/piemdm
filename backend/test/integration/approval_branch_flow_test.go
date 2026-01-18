package integration

import (
	"testing"

	"piemdm/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovalBranchFlowControl(t *testing.T) {
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

	approvalNodes := []*model.ApprovalNode{
		{NodeCode: "start", NodeName: "开始", NodeType: "START", SortOrder: 1},
		{NodeCode: "condition-price", NodeName: "价格条件判断", NodeType: "CONDITION", SortOrder: 2, ConditionConfig: conditionConfigJSON},
		{NodeCode: "finance-approval", NodeName: "财务审批", NodeType: "APPROVAL", SortOrder: 100},
		{NodeCode: "auto-reject", NodeName: "自动驳回", NodeType: "APPROVAL", SortOrder: 101},
		{NodeCode: "final-review", NodeName: "最终审核", NodeType: "APPROVAL", SortOrder: 3},
		{NodeCode: "end", NodeName: "结束", NodeType: "END", SortOrder: 4},
	}

	tests := []struct {
		name             string
		currentNodeCode  string
		formData         map[string]any
		expectedNextNode string
	}{
		{
			name:             "条件节点-价格>=3000-应该跳转到财务审批",
			currentNodeCode:  "condition-price",
			formData:         map[string]any{"price": 3500},
			expectedNextNode: "finance-approval",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var currentNode *model.ApprovalNode
			for _, node := range approvalNodes {
				if node.NodeCode == tt.currentNodeCode {
					currentNode = node
					break
				}
			}
			require.NotNil(t, currentNode, "should find current node")

			nextNode, err := evaluateConditionNode(approvalNodes, currentNode, tt.formData)

			require.NoError(t, err)
			require.NotNil(t, nextNode)
			assert.Equal(t, tt.expectedNextNode, nextNode.NodeCode)
		})
	}
}
