package repository

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCondition_MultiValue(t *testing.T) {
	tests := []struct {
		name          string
		where         map[string]any
		wantWhereSql  string
		wantValuesLen int
		wantErr       bool
	}{
		{
			name: "single value",
			where: map[string]any{
				"code": "TEST001",
			},
			wantWhereSql:  "code = ?",
			wantValuesLen: 1,
			wantErr:       false,
		},
		{
			name: "multi value with newline",
			where: map[string]any{
				"code": "TEST001\nTEST002\nTEST003",
			},
			wantWhereSql:  "code in ?",
			wantValuesLen: 1,
			wantErr:       false,
		},
		{
			name: "explicit equal with newline",
			where: map[string]any{
				"code =": "TEST001\nTEST002",
			},
			wantWhereSql:  "code in ?",
			wantValuesLen: 1,
			wantErr:       false,
		},
		{
			name: "multi value with newline and empty lines",
			where: map[string]any{
				"code": "TEST001\n\nTEST002\n  \nTEST003",
			},
			wantWhereSql:  "code in ?",
			wantValuesLen: 1,
			wantErr:       false,
		},
		{
			name: "mixed conditions",
			where: map[string]any{
				"name like": "%Test%",
				"status":    "active\npending",
			},
			// Map iteration order is random, so we check contains
			wantWhereSql:  "",
			wantValuesLen: 2,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWhereSql, gotValues, err := BuildCondition(tt.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "mixed conditions" {
				assert.Contains(t, gotWhereSql, "name like ?")
				assert.Contains(t, gotWhereSql, "status in ?")
				assert.Len(t, gotValues, tt.wantValuesLen)
			} else {
				assert.Equal(t, tt.wantWhereSql, gotWhereSql)
				assert.Len(t, gotValues, tt.wantValuesLen)

				if strings.Contains(gotWhereSql, "in ?") {
					// Verify that the value is a slice
					assert.IsType(t, []string{}, gotValues[0])
					sliceVal := gotValues[0].([]string)
					assert.True(t, len(sliceVal) > 1)
				}
			}
		})
	}
}
