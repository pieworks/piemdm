package service

import (
	"fmt"
	"regexp"
	"strconv"

	"piemdm/internal/model"
)

// ValidateFieldValue 验证字段值
func ValidateFieldValue(field *model.TableField, value any) error {
	strValue := fmt.Sprintf("%v", value)
	// 必填验证
	if field.Required == "Yes" && (value == nil || strValue == "") {
		return fmt.Errorf("字段 '%s' 是必填项", field.Name)
	}

	if value == nil || strValue == "" {
		return nil // 非必填且为空，跳过其他验证
	}

	rules := field.Options
	if rules == nil || rules.Validation == nil {
		return nil
	}

	validation := rules.Validation

	// 长度/数值范围验证
	isNumber := false
	var numVal float64

	switch v := value.(type) {
	case int:
		numVal = float64(v)
		isNumber = true
	case int8:
		numVal = float64(v)
		isNumber = true
	case int16:
		numVal = float64(v)
		isNumber = true
	case int32:
		numVal = float64(v)
		isNumber = true
	case int64:
		numVal = float64(v)
		isNumber = true
	case float32:
		numVal = float64(v)
		isNumber = true
	case float64:
		numVal = float64(v)
		isNumber = true
	case uint, uint8, uint16, uint32, uint64:
		str := fmt.Sprintf("%v", v)
		if val, err := strconv.ParseFloat(str, 64); err == nil {
			numVal = val
			isNumber = true
		}
	}

	if isNumber {
		if validation.Max != nil && numVal > float64(*validation.Max) {
			return fmt.Errorf("字段 '%s' 数值不能大于 %d", field.Name, *validation.Max)
		}
		if validation.Min != nil && numVal < float64(*validation.Min) {
			return fmt.Errorf("字段 '%s' 数值不能小于 %d", field.Name, *validation.Min)
		}
	} else {
		if validation.Max != nil && len(strValue) > *validation.Max {
			return fmt.Errorf("字段 '%s' 长度不能超过 %d", field.Name, *validation.Max)
		}
		if validation.Min != nil && len(strValue) < *validation.Min {
			return fmt.Errorf("字段 '%s' 长度不能少于 %d", field.Name, *validation.Min)
		}
	}

	// 正则验证
	if validation.Pattern != "" {
		re, err := regexp.Compile(validation.Pattern)
		if err != nil {
			return fmt.Errorf("正则表达式错误: %v", err)
		}
		if !re.MatchString(strValue) {
			msg := validation.Message
			if msg == "" {
				msg = fmt.Sprintf("字段 '%s' 格式不正确", field.Name)
			}
			return fmt.Errorf("%s", msg)
		}
	}

	// 格式验证
	if validation.Format != "" {
		if err := validateFormat(validation.Format, strValue, field.Name); err != nil {
			return err
		}
	}

	// 预定义验证器
	if validation.Validator != "" {
		if err := validateByValidator(validation.Validator, value, strValue, field.Name); err != nil {
			return err
		}
	}

	return nil
}

// validateFormat 验证预定义格式
func validateFormat(format, value, fieldName string) error {
	patterns := map[string]string{
		"email": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"phone": `^1[3-9]\d{9}$`,
		"url":   `^https?://[^\s]+$`,
	}

	pattern, ok := patterns[format]
	if !ok {
		return nil // 未知格式，跳过验证
	}

	re := regexp.MustCompile(pattern)
	if !re.MatchString(value) {
		return fmt.Errorf("字段 '%s' 格式不正确（期望: %s）", fieldName, format)
	}

	return nil
}

// validateByValidator 使用预定义验证器验证
func validateByValidator(validator string, rawValue any, value, fieldName string) error {
	switch validator {
	case "integer":
		// 验证是否为整数
		// 首先检查原始值的类型
		switch rawValue.(type) {
		case string:
			return fmt.Errorf("字段 '%s' 必须是整数类型", fieldName)
		}

		re := regexp.MustCompile(`^-?\d+$`)
		if !re.MatchString(value) {
			return fmt.Errorf("字段 '%s' 必须是整数", fieldName)
		}
	case "email":
		return validateFormat("email", value, fieldName)
	case "phone":
		return validateFormat("phone", value, fieldName)
	case "url":
		return validateFormat("url", value, fieldName)
	}

	return nil
}
