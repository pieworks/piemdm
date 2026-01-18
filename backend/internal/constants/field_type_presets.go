package constants

import (
	"encoding/json"
	"piemdm/internal/model"
)

// FieldPreset 字段类型预设配置
type FieldPreset struct {
	// 显示名称
	Label string `json:"label"`
	// 分组
	Group string `json:"group"`
	// 数据类型: Text, Number, Date, Boolean
	DataType string `json:"dataType"`
	// 字段长度（仅 string 类型）
	Length int `json:"length,omitempty"`
	// 精度（仅 decimal 类型）
	Precision int `json:"precision,omitempty"`
	// 小数位数（仅 decimal 类型）
	Scale int `json:"scale,omitempty"`
	// UI 配置
	UI *model.FieldUI `json:"ui"`
	// 验证规则
	Validation *model.FieldValidation `json:"validation,omitempty"`
	// 是否需要数据源配置
	RequireDatasource bool `json:"requireDatasource,omitempty"`
	// 是否需要关联配置
	RequireRelation bool `json:"requireRelation,omitempty"`
}

// FieldTypeGroup 字段类型分组
type FieldTypeGroup struct {
	Name  string   `json:"name"`  // 分组名称
	Label string   `json:"label"` // 显示名称
	Types []string `json:"types"` // 包含的字段类型
}

// 辅助函数：创建 int 指针
func intPtr(i int) *int {
	return &i
}

// fieldTypePresets 字段类型预设配置映射
var fieldTypePresets = map[string]*FieldPreset{
	// ===== 基本类型 (Basic) =====
	"text": {
		Label:    "单行文本",
		Group:    "basic",
		DataType: "Text",
		Length:   255,
		UI: &model.FieldUI{
			Widget: "Input",
		},
		Validation: &model.FieldValidation{
			Max: intPtr(255),
		},
	},
	"textarea": {
		Label:    "多行文本",
		Group:    "basic",
		DataType: "Text",
		Length:   2000,
		UI: &model.FieldUI{
			Widget: "Textarea",
		},
		Validation: &model.FieldValidation{
			Max: intPtr(2000),
		},
	},
	"phone": {
		Label:    "手机号码",
		Group:    "basic",
		DataType: "Text",
		Length:   20,
		UI: &model.FieldUI{
			Widget: "Input",
			WidgetProps: map[string]any{
				"type": "tel",
			},
		},
		Validation: &model.FieldValidation{
			Format:  "phone",
			Pattern: "^1[3-9]\\d{9}$",
			Max:     intPtr(20),
		},
	},
	"email": {
		Label:    "电子邮箱",
		Group:    "basic",
		DataType: "Text",
		Length:   100,
		UI: &model.FieldUI{
			Widget: "Input",
			WidgetProps: map[string]any{
				"type": "email",
			},
		},
		Validation: &model.FieldValidation{
			Format: "email",
			Max:    intPtr(100),
		},
	},
	"url": {
		Label:    "URL",
		Group:    "basic",
		DataType: "Text",
		Length:   255,
		UI: &model.FieldUI{
			Widget: "Input",
			WidgetProps: map[string]any{
				"type": "url",
			},
		},
		Validation: &model.FieldValidation{
			Format: "url",
			Max:    intPtr(255),
		},
	},
	"integer": {
		Label:    "整数",
		Group:    "basic",
		DataType: "Number",
		UI: &model.FieldUI{
			Widget: "InputNumber",
			WidgetProps: map[string]any{
				"step": 1,
			},
		},
		Validation: &model.FieldValidation{
			Validator: "integer",
		},
	},
	"decimal": {
		Label:     "小数",
		Group:     "basic",
		DataType:  "Number",
		Precision: 10,
		Scale:     2,
		UI: &model.FieldUI{
			Widget: "InputNumber",
			WidgetProps: map[string]any{
				"step": 0.01,
			},
		},
	},
	"percent": {
		Label:     "百分比",
		Group:     "basic",
		DataType:  "Number",
		Precision: 5,
		Scale:     2,
		UI: &model.FieldUI{
			Widget: "InputNumber",
			WidgetProps: map[string]any{
				"step":   0.01,
				"min":    0,
				"max":    100,
				"suffix": "%",
			},
		},
		Validation: &model.FieldValidation{
			Min: intPtr(0),
			Max: intPtr(100),
		},
	},
	"password": {
		Label:    "密码",
		Group:    "basic",
		DataType: "Text",
		Length:   128,
		UI: &model.FieldUI{
			Widget: "Input",
			WidgetProps: map[string]any{
				"type": "password",
			},
		},
		Validation: &model.FieldValidation{
			Min: intPtr(6),
			Max: intPtr(128),
		},
	},

	// ===== 选择类型 (Choices) =====
	"checkbox": {
		Label:    "勾选",
		Group:    "choices",
		DataType: "Boolean",
		UI: &model.FieldUI{
			Widget: "Checkbox",
		},
	},
	"select": {
		Label:    "下拉单选",
		Group:    "choices",
		DataType: "Text",
		Length:   64,
		UI: &model.FieldUI{
			Widget: "Select",
		},
		RequireDatasource: true,
		RequireRelation:   true,
	},
	"multiselect": {
		Label:    "下拉多选",
		Group:    "choices",
		DataType: "Text",
		Length:   500,
		UI: &model.FieldUI{
			Widget: "MultiSelect",
		},
		RequireDatasource: true,
		RequireRelation:   true,
	},
	"radio": {
		Label:    "单选框",
		Group:    "choices",
		DataType: "Text",
		Length:   64,
		UI: &model.FieldUI{
			Widget: "RadioGroup",
		},
		RequireDatasource: true,
		RequireRelation:   true,
	},
	"checkboxgroup": {
		Label:    "复选框组",
		Group:    "choices",
		DataType: "Text",
		Length:   500,
		UI: &model.FieldUI{
			Widget: "CheckboxGroup",
		},
		RequireDatasource: true,
		RequireRelation:   true,
	},

	// ===== 日期时间 (DateTime) =====
	"date": {
		Label:    "日期",
		Group:    "datetime",
		DataType: "Date",
		UI: &model.FieldUI{
			Widget: "DatePicker",
			WidgetProps: map[string]any{
				"format": "YYYY-MM-DD",
			},
		},
	},
	"time": {
		Label:    "时间",
		Group:    "datetime",
		DataType: "Text",
		Length:   8,
		UI: &model.FieldUI{
			Widget: "TimePicker",
			WidgetProps: map[string]any{
				"format": "HH:mm:ss",
			},
		},
	},
	"datetime": {
		Label:    "日期时间",
		Group:    "datetime",
		DataType: "DateTime", // 使用 DateTime 区分 DATE 和 DATETIME
		UI: &model.FieldUI{
			Widget: "DateTimePicker",
			WidgetProps: map[string]any{
				"format":   "YYYY-MM-DD HH:mm:ss",
				"showTime": true,
			},
		},
	},

	// ===== 关系类型 (Relation) =====
	"belongsto": {
		Label:    "多对一",
		Group:    "relation",
		DataType: "Text",
		Length:   64,
		UI: &model.FieldUI{
			Widget: "Select",
		},
		RequireRelation: true,
	},
	"hasmany": {
		Label:    "一对多",
		Group:    "relation",
		DataType: "Text",
		Length:   500,
		UI: &model.FieldUI{
			Widget: "SubTable",
		},
		RequireRelation: true,
	},
	"manytomany": {
		Label:    "多对多",
		Group:    "relation",
		DataType: "Text",
		Length:   500,
		UI: &model.FieldUI{
			Widget: "MultiSelect",
		},
		RequireRelation: true,
	},

	// ===== 高级类型 (Advanced) =====
	"autocode": {
		Label:    "自动编码",
		Group:    "advanced",
		DataType: "Text",
		Length:   64,
		UI: &model.FieldUI{
			Widget: "Input",
			WidgetProps: map[string]any{
				"disabled": true,
			},
		},
	},
	// "formula": {
	// 	Label:    "公式",
	// 	Group:    "advanced",
	// 	DataType: "Text",
	// 	Length:   255,
	// 	UI: &model.FieldUI{
	// 		Widget: "Input",
	// 		WidgetProps: map[string]any{
	// 			"disabled": true,
	// 		},
	// 	},
	// },
	// "json": {  // 暂时注释掉,以后实现
	// 	Label:    "JSON",
	// 	Group:    "advanced",
	// 	DataType: "Text",
	// 	Length:   2000,
	// 	UI: &model.FieldUI{
	// 		Widget: "Textarea",
	// 		WidgetProps: map[string]any{
	// 			"rows": 10,
	// 		},
	// 	},
	// },
	"attachment": {
		Label:    "附件",
		Group:    "advanced",
		DataType: "Text",
		Length:   500,
		UI: &model.FieldUI{
			Widget: "Upload",
		},
	},
	// "richtext": {
	// 	Label:    "富文本",
	// 	Group:    "advanced",
	// 	DataType: "Text",
	// 	Length:   10000,
	// 	UI: &model.FieldUI{
	// 		Widget: "RichTextEditor",
	// 	},
	// },
	// "sort": {
	// 	Label:    "排序",
	// 	Group:    "advanced",
	// 	DataType: "Number",
	// 	UI: &model.FieldUI{
	// 		Widget: "InputNumber",
	// 	},
	// },
}

// fieldTypeGroups 字段类型分组
var fieldTypeGroups = []FieldTypeGroup{
	{
		Name:  "basic",
		Label: "基本类型",
		Types: []string{"text", "textarea", "phone", "email", "url", "integer", "decimal", "percent", "password"},
	},
	{
		Name:  "choices",
		Label: "选择类型",
		Types: []string{"checkbox", "select", "multiselect", "radio", "checkboxgroup"},
	},
	{
		Name:  "datetime",
		Label: "日期时间",
		Types: []string{"date", "time", "datetime"},
	},
	{
		Name:  "relation",
		Label: "关系类型",
		Types: []string{"belongsto", "hasmany", "manytomany"},
	},
	{
		Name:  "advanced",
		Label: "高级类型",
		Types: []string{"autocode", "formula", "attachment", "richtext", "sort"}, // json 暂时注释掉
	},
}

// GetFieldPreset 获取字段类型预设配置
func GetFieldPreset(fieldType string) (*FieldPreset, bool) {
	preset, ok := fieldTypePresets[fieldType]
	if !ok {
		return nil, false
	}
	// 深拷贝以避免修改原始配置
	return preset.Clone(), true
}

// GetAllFieldTypePresets 获取所有字段类型预设
func GetAllFieldTypePresets() map[string]*FieldPreset {
	result := make(map[string]*FieldPreset)
	for k, v := range fieldTypePresets {
		result[k] = v.Clone()
	}
	return result
}

// GetFieldTypeGroups 获取字段类型分组
func GetFieldTypeGroups() []FieldTypeGroup {
	// 返回副本
	groups := make([]FieldTypeGroup, len(fieldTypeGroups))
	copy(groups, fieldTypeGroups)
	return groups
}

// Clone 深拷贝字段预设配置
func (p *FieldPreset) Clone() *FieldPreset {
	if p == nil {
		return nil
	}

	clone := &FieldPreset{
		Label:             p.Label,
		Group:             p.Group,
		DataType:          p.DataType,
		Length:            p.Length,
		Precision:         p.Precision,
		Scale:             p.Scale,
		RequireDatasource: p.RequireDatasource,
		RequireRelation:   p.RequireRelation,
	}

	// 深拷贝 UI 配置
	if p.UI != nil {
		clone.UI = &model.FieldUI{
			Widget:      p.UI.Widget,
			Title:       p.UI.Title,
			Placeholder: p.UI.Placeholder,
		}
		if p.UI.WidgetProps != nil {
			clone.UI.WidgetProps = make(map[string]any)
			// 深拷贝 map
			data, _ := json.Marshal(p.UI.WidgetProps)
			json.Unmarshal(data, &clone.UI.WidgetProps)
		}
	}

	// 深拷贝 Validation 配置
	if p.Validation != nil {
		clone.Validation = &model.FieldValidation{
			Format:    p.Validation.Format,
			Pattern:   p.Validation.Pattern,
			Message:   p.Validation.Message,
			Validator: p.Validation.Validator,
		}
		if p.Validation.Min != nil {
			min := *p.Validation.Min
			clone.Validation.Min = &min
		}
		if p.Validation.Max != nil {
			max := *p.Validation.Max
			clone.Validation.Max = &max
		}
	}

	return clone
}
