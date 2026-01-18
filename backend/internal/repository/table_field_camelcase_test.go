package repository

import (
	"testing"
)

func TestToCamelCase(t *testing.T) {
	r := &tableFieldRepository{}

	tests := []struct {
		input    string
		expected string
	}{
		// 基本的下划线分隔
		{"user_name", "UserName"},
		{"first_name", "FirstName"},
		{"table_code", "TableCode"},

		// 连字符分隔
		{"fie-1222", "Fie1222"},
		{"test-field", "TestField"},
		{"my-custom-field", "MyCustomField"},

		// 混合分隔符
		{"user_name-v2", "UserNameV2"},
		{"field-test_123", "FieldTest123"},

		// 以数字开头
		{"123field", "Field123field"},
		{"1test", "Field1test"},

		// 特殊字符
		{"field@test", "FieldTest"},
		{"test#field$name", "TestFieldName"},

		// 空字符串
		{"", "UnknownField"},
		{"___", "UnknownField"},

		// 单个单词
		{"name", "Name"},
		{"id", "Id"},

		// 已经是驼峰命名
		{"UserName", "UserName"},
		{"TestField", "TestField"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := r.toCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
