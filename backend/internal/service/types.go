package service

// PageResult 分页结果通用类型
type PageResult[T any] struct {
	Data     []T   `json:"data"`      // 数据列表
	Total    int64 `json:"total"`     // 总数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 页面大小
}
