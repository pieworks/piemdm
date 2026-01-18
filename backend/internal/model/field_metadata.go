package model

// FieldMetadata 字段元数据
// 用于返回表的所有字段信息,包括用户定义的业务字段和系统字段
type FieldMetadata struct {
	Code      string `json:"code"`              // 字段代码
	Name      string `json:"name"`              // 字段名称
	FieldType string `json:"field_type"`        // 字段类型: text/integer/autocode 等
	Type      string `json:"type"`              // 数据类型: Text/Number/Date
	Length    int    `json:"length,omitempty"`  // 字段长度(仅Text类型)
	Required  bool   `json:"required"`          // 是否必填
	IsShow    bool   `json:"is_show"`           // 是否显示
	IsSystem  bool   `json:"is_system"`         // 是否系统字段
	IsFilter  bool   `json:"is_filter"`         // 是否可过滤
	Sort      int    `json:"sort"`              // 排序
	Options   any    `json:"options,omitempty"` // 字段配置选项(JSON)
}
