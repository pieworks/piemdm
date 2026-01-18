package integration

import (
	"encoding/json"
	"fmt"
	"strconv"

	"piemdm/internal/model"
)

// evaluateConditionNode evaluates the conditions in a CONDITION node and returns the next node.
func evaluateConditionNode(approvalNodes []*model.ApprovalNode, conditionNode *model.ApprovalNode, formData map[string]any) (*model.ApprovalNode, error) {
	var conditionConfig struct {
		Branches []struct {
			Condition struct {
				FieldName  string `json:"fieldName"`
				Operator   string `json:"operator"`
				FieldValue string `json:"fieldValue"`
			} `json:"condition"`
			Nodes []struct {
				NodeCode string `json:"nodeCode"`
			} `json:"nodes"`
		} `json:"branches"`
	}

	if err := json.Unmarshal([]byte(conditionNode.ConditionConfig), &conditionConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal condition config: %w", err)
	}

	for _, branch := range conditionConfig.Branches {
		if evaluateCondition(branch.Condition, formData) {
			if len(branch.Nodes) > 0 {
				nextNodeCode := branch.Nodes[0].NodeCode
				for _, node := range approvalNodes {
					if node.NodeCode == nextNodeCode {
						return node, nil
					}
				}
				return nil, fmt.Errorf("next node with code '%s' not found", nextNodeCode)
			}
		}
	}

	return nil, nil // No condition met
}

// evaluateCondition evaluates a single condition against the form data.
func evaluateCondition(condition struct {
	FieldName  string `json:"fieldName"`
	Operator   string `json:"operator"`
	FieldValue string `json:"fieldValue"`
}, formData map[string]any,
) bool {
	if condition.FieldName == "" {
		return condition.Operator == "eq" && condition.FieldValue == ""
	}

	fieldValue, exists := formData[condition.FieldName]
	if !exists {
		return false
	}

	fieldValueStr := convertToString(fieldValue)
	expectedValue := condition.FieldValue

	switch condition.Operator {
	case "gte", ">=":
		fieldNum, err1 := strconv.ParseFloat(fieldValueStr, 64)
		expectedNum, err2 := strconv.ParseFloat(expectedValue, 64)
		if err1 == nil && err2 == nil {
			return fieldNum >= expectedNum
		}
		return fieldValueStr >= expectedValue
	default:
		return false
	}
}

// convertToString converts an interface{} to a string.
func convertToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// parseNumber parses a string into a float64.
func parseNumber(s string) (float64, error) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return float64(i), nil
	}
	return strconv.ParseFloat(s, 64)
}
