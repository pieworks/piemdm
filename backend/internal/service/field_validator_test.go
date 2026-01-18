package service

import (
	"piemdm/internal/model"
	"testing"
)

func TestValidateFieldValue_Required(t *testing.T) {
	field := &model.TableField{
		Code:     "test_field",
		Required: "Yes",
	}

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"empty string should fail", "", true},
		{"nil should fail", nil, true},
		{"valid string should pass", "test", false},
		{"zero number should pass", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldValue(field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_MinMax(t *testing.T) {
	minVal := 5
	maxVal := 10
	field := &model.TableField{
		Code: "test_field",
		Options: &model.FieldOptions{
			Validation: &model.FieldValidation{
				Min: &minVal,
				Max: &maxVal,
			},
		},
	}

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"string too short", "abc", true},
		{"string too long", "12345678901", true},
		{"string valid length", "12345", false},
		{"number too small", 3, true},
		{"number too large", 15, true},
		{"number valid", 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldValue(field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_Pattern(t *testing.T) {
	field := &model.TableField{
		Code: "phone",
		Options: &model.FieldOptions{
			Validation: &model.FieldValidation{
				Pattern: `^1[3-9]\d{9}$`,
			},
		},
	}

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"valid phone", "13812345678", false},
		{"invalid phone - too short", "138123456", true},
		{"invalid phone - wrong prefix", "12812345678", true},
		{"invalid phone - not number", "1381234567a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldValue(field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_Format(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		value   string
		wantErr bool
	}{
		{"valid email", "email", "test@example.com", false},
		{"invalid email", "email", "invalid-email", true},
		{"valid phone", "phone", "13812345678", false},
		{"invalid phone", "phone", "12345", true},
		{"valid url", "url", "https://example.com", false},
		{"invalid url", "url", "not-a-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &model.TableField{
				Code: "test_field",
				Options: &model.FieldOptions{
					Validation: &model.FieldValidation{
						Format: tt.format,
					},
				},
			}
			err := ValidateFieldValue(field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_Validator(t *testing.T) {
	field := &model.TableField{
		Code: "test_field",
		Options: &model.FieldOptions{
			Validation: &model.FieldValidation{
				Validator: "integer",
			},
		},
	}

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"integer should pass", 123, false},
		{"float should fail", 123.45, true},
		{"string number should fail", "123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldValue(field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
