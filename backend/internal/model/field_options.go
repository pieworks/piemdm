package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FieldOptions JSON 结构
type FieldOptions struct {
	UI         *FieldUI         `json:"ui,omitempty"`         // 前端 UI 配置
	Validation *FieldValidation `json:"validation,omitempty"` // 验证规则
	DataSource *FieldDataSource `json:"datasource,omitempty"` // 数据源 (含关联) - 保留兼容
	Behavior   *FieldBehavior   `json:"behavior,omitempty"`   // 行为配置

	// 新增配置
	Relation   *FieldRelation    `json:"relation,omitempty"`   // 关联配置
	DateTime   *FieldDateTime    `json:"datetime,omitempty"`   // 日期时间行为
	Sort       *FieldSort        `json:"sort,omitempty"`       // 排序字段
	RichText   *FieldRichText    `json:"richtext,omitempty"`   // 富文本
	Attachment *FieldAttachment  `json:"attachment,omitempty"` // 附件
	Patterns   []SequencePattern `json:"patterns,omitempty"`   // 自动编码
	Trim       bool              `json:"trim,omitempty"`       // 字符串去空格
}

// FieldUI UI 配置
type FieldUI struct {
	Widget      string         `json:"widget"`                // 组件类型 (必填)
	WidgetProps map[string]any `json:"widgetProps,omitempty"` // 组件参数 (可选)
	Title       string         `json:"title,omitempty"`       // 标题 (可选)
	Placeholder string         `json:"placeholder,omitempty"` // 占位符 (可选)
}

// FieldValidation 验证规则
type FieldValidation struct {
	Required  bool   `json:"required,omitempty"`
	Format    string `json:"format,omitempty"`    // 格式: phone, email, url
	Pattern   string `json:"pattern,omitempty"`   // 自定义正则
	Message   string `json:"message,omitempty"`   // 自定义错误消息
	Validator string `json:"validator,omitempty"` // 预定义验证器: integer, email, phone, url
	Min       *int   `json:"min,omitempty"`       // 最小值
	Max       *int   `json:"max,omitempty"`       // 最大值
	Precision *int   `json:"precision,omitempty"` // 总位数 (仅 decimal 类型)
	Scale     *int   `json:"scale,omitempty"`     // 小数位数 (仅 decimal 类型)
}

// FieldDataSource 数据源定义
type FieldDataSource struct {
	Type string `json:"type"` // static, api, relation
	// Static
	Options []OptionItem `json:"options,omitempty"`
	// API
	API string `json:"api,omitempty"`
	// Relation
	TargetTable string            `json:"targetTable,omitempty"`
	ForeignKey  string            `json:"foreignKey,omitempty"`
	ValueField  string            `json:"valueField,omitempty"`
	LabelField  string            `json:"labelField,omitempty"`
	Filter      map[string]any    `json:"filter,omitempty"`
	FillBack    map[string]string `json:"fillBack,omitempty"`
}

// OptionItem 选项项
type OptionItem struct {
	Label string `json:"label"`
	Value any    `json:"value"`
	Color string `json:"color,omitempty"`
}

// FieldBehavior 行为配置
type FieldBehavior struct {
	DefaultValue any    `json:"defaultValue,omitempty"`
	ReadOnly     bool   `json:"readOnly,omitempty"`
	Show         bool   `json:"show,omitempty"`      // 是否显示 (对应 is_show)
	VisibleOn    string `json:"visibleOn,omitempty"` // 表达式
}

// Value 实现 driver.Valuer 接口
func (o FieldOptions) Value() (driver.Value, error) {
	return json.Marshal(o)
}

// Scan 实现 sql.Scanner 接口
func (o *FieldOptions) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, o)
}

// FieldRelation 关联配置
type FieldRelation struct {
	Target     string         `json:"target"`               // 关联表
	ForeignKey string         `json:"foreignKey,omitempty"` // 外键（可自动生成）
	LabelField string         `json:"labelField,omitempty"` // 显示字段（默认 name）
	ValueField string         `json:"valueField,omitempty"` // 存储字段 = 关联键（默认 code）
	Multiple   bool           `json:"multiple,omitempty"`   // 是否多选
	Filter     map[string]any `json:"filter,omitempty"`     // 过滤条件
}

// FieldDateTime 日期时间行为配置
type FieldDateTime struct {
	Timezone              bool   `json:"timezone,omitempty"`              // 是否支持时区
	DefaultToCurrentTime  bool   `json:"defaultToCurrentTime,omitempty"`  // 创建时默认当前时间
	OnUpdateToCurrentTime bool   `json:"onUpdateToCurrentTime,omitempty"` // 更新时自动更新为当前时间
	ShowTime              bool   `json:"showTime,omitempty"`              // 是否显示时间选择器
	DateFormat            string `json:"dateFormat,omitempty"`            // 日期格式
	TimeFormat            string `json:"timeFormat,omitempty"`            // 时间格式
}

// FieldSort 排序字段配置
type FieldSort struct {
	ScopeKey string `json:"scopeKey,omitempty"` // 排序作用域字段（如按 status 分组排序）
}

// FieldRichText 富文本配置
type FieldRichText struct {
	Length         string   `json:"length,omitempty"`         // short, long
	Toolbar        []string `json:"toolbar,omitempty"`        // 工具栏配置
	FileCollection string   `json:"fileCollection,omitempty"` // 附件存储集合名称
}

// FieldAttachment 附件配置
type FieldAttachment struct {
	Multiple bool     `json:"multiple,omitempty"` // 是否支持多文件上传
	Accept   []string `json:"accept,omitempty"`   // 允许的文件类型
	MaxSize  int      `json:"maxSize,omitempty"`  // 最大文件大小（字节）
}

// SequencePattern 自动编码模式
type SequencePattern struct {
	Type    string         `json:"type"`    // string, integer, date
	Options map[string]any `json:"options"` // 配置选项
}
