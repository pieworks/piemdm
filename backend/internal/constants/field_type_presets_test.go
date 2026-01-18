package constants

import (
	"testing"
)

func TestGetFieldPreset(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		wantNil   bool
	}{
		{"text preset exists", "text", false},
		{"phone preset exists", "phone", false},
		{"select preset exists", "select", false},
		{"belongsto preset exists", "belongsto", false},
		{"invalid preset", "invalid_type", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset, ok := GetFieldPreset(tt.fieldType)
			if tt.wantNil && ok {
				t.Errorf("GetFieldPreset(%s) should return false, got true", tt.fieldType)
			}
			if !tt.wantNil && !ok {
				t.Errorf("GetFieldPreset(%s) should return true, got false", tt.fieldType)
			}
			if tt.wantNil && preset != nil {
				t.Errorf("GetFieldPreset(%s) should return nil preset", tt.fieldType)
			}
			if !tt.wantNil && preset == nil {
				t.Errorf("GetFieldPreset(%s) should not return nil preset", tt.fieldType)
			}
		})
	}
}

func TestFieldPresetProperties(t *testing.T) {
	// Test text preset
	textPreset, ok := GetFieldPreset("text")
	if !ok || textPreset == nil {
		t.Fatal("text preset should exist")
	}
	if textPreset.DataType != "Text" {
		t.Errorf("text preset DataType = %s, want Text", textPreset.DataType)
	}
	if textPreset.Length != 255 {
		t.Errorf("text preset Length = %d, want 255", textPreset.Length)
	}
	if textPreset.UI.Widget != "Input" {
		t.Errorf("text preset UI.Widget = %s, want Input", textPreset.UI.Widget)
	}

	// Test phone preset
	phonePreset, ok := GetFieldPreset("phone")
	if !ok || phonePreset == nil {
		t.Fatal("phone preset should exist")
	}
	if phonePreset.Validation.Format != "phone" {
		t.Errorf("phone preset Validation.Format = %s, want phone", phonePreset.Validation.Format)
	}
	if phonePreset.Validation.Pattern == "" {
		t.Error("phone preset should have a pattern")
	}

	// Test integer preset
	integerPreset, ok := GetFieldPreset("integer")
	if !ok || integerPreset == nil {
		t.Fatal("integer preset should exist")
	}
	if integerPreset.DataType != "Number" {
		t.Errorf("integer preset DataType = %s, want Number", integerPreset.DataType)
	}
	if integerPreset.Validation.Validator != "integer" {
		t.Errorf("integer preset Validation.Validator = %s, want integer", integerPreset.Validation.Validator)
	}
}

func TestGetFieldTypeGroups(t *testing.T) {
	groups := GetFieldTypeGroups()
	if len(groups) == 0 {
		t.Fatal("GetFieldTypeGroups() should return at least one group")
	}

	expectedGroups := map[string]bool{
		"basic":    false,
		"choices":  false,
		"datetime": false,
		"relation": false,
		"advanced": false,
	}

	for _, group := range groups {
		if _, exists := expectedGroups[group.Name]; exists {
			expectedGroups[group.Name] = true
		}
		if len(group.Types) == 0 {
			t.Errorf("Group %s should have at least one type", group.Name)
		}
	}

	for groupName, found := range expectedGroups {
		if !found {
			t.Errorf("Expected group %s not found", groupName)
		}
	}
}

func TestGetAllFieldTypePresets(t *testing.T) {
	presets := GetAllFieldTypePresets()
	if len(presets) == 0 {
		t.Fatal("GetAllFieldTypePresets() should return at least one preset")
	}

	// Check that all presets have required fields
	for fieldType, preset := range presets {
		if preset.DataType == "" {
			t.Errorf("Preset %s should have DataType", fieldType)
		}
		if preset.UI.Widget == "" {
			t.Errorf("Preset %s should have UI.Widget", fieldType)
		}
	}
}
