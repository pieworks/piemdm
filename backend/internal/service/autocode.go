package service

import (
	"fmt"
	"piemdm/internal/model"
	"strings"
	"time"
)

type AutocodeService interface {
	// 业务方法
	GenerateCode(tableCode, fieldCode string, patterns []model.SequencePattern, entityMap map[string]any) (string, error)
}

type autocodeService struct {
	*Service
	globalIdService GlobalIdService
}

func NewAutocodeService(service *Service, globalIdService GlobalIdService) AutocodeService {
	return &autocodeService{
		Service:         service,
		globalIdService: globalIdService,
	}
}

// GenerateCode 根据配置的模式生成自动编码
func (s *autocodeService) GenerateCode(tableCode, fieldCode string, patterns []model.SequencePattern, entityMap map[string]any) (string, error) {
	var result string
	var cycle string
	var start int = 1 // 默认起始值

	// 遍历模式生成编码
	for _, pattern := range patterns {
		switch pattern.Type {
		case "string":
			// 固定字符串
			if value, ok := pattern.Options["value"].(string); ok {
				result += value
			}
		case "date":
			// 日期格式化
			if format, ok := pattern.Options["format"].(string); ok {
				result += s.formatDate(time.Now(), format)
			}
		case "field":
			// 引用数据字段
			if fieldCodeStr, ok := pattern.Options["fieldCode"].(string); ok {
				// 从 entityMap 中获取字段值
				if fieldValue, exists := entityMap[fieldCodeStr]; exists && fieldValue != nil {
					// 转换为字符串
					var strValue string
					switch v := fieldValue.(type) {
					case string:
						strValue = v
					case int, int8, int16, int32, int64:
						strValue = fmt.Sprintf("%d", v)
					case uint, uint8, uint16, uint32, uint64:
						strValue = fmt.Sprintf("%d", v)
					case float64:
						// 转为整数字符串(去掉小数部分)
						strValue = fmt.Sprintf("%.0f", v)
					default:
						return "", fmt.Errorf("field %s has unsupported type for autocode: %T", fieldCodeStr, v)
					}
					// 转大写
					result += strings.ToUpper(strValue)
				} else {
					return "", fmt.Errorf("field %s is required for autocode but not provided or is nil", fieldCodeStr)
				}
			}
		case "integer":
			// 序列号
			if cycleStr, ok := pattern.Options["cycle"].(string); ok {
				cycle = cycleStr
			} else {
				cycle = "none"
			}

			// 获取起始值
			if startVal, ok := pattern.Options["start"].(float64); ok {
				start = int(startVal)
			}

			cycleValue := s.getCycleValue(cycle)
			identifier := s.buildIdentifier(tableCode, fieldCode, cycleValue)

			// 确保 global_id 记录存在
			if err := s.ensureGlobalIdExists(identifier, start, tableCode, fieldCode); err != nil {
				return "", fmt.Errorf("failed to ensure global_id exists: %v", err)
			}

			// 获取下一个序列号
			nextID := s.globalIdService.GetNewID(identifier)

			// 格式化序列号
			digits := 5 // 默认5位
			if digitsVal, ok := pattern.Options["digits"].(float64); ok {
				digits = int(digitsVal)
			}
			result += fmt.Sprintf("%0*d", digits, nextID)
		}
	}

	return result, nil
}

// getCycleValue 根据周期类型获取周期值
func (s *autocodeService) getCycleValue(cycle string) string {
	now := time.Now()
	switch cycle {
	case "daily":
		return now.Format("2006-01-02")
	case "weekly":
		year, week := now.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "monthly":
		return now.Format("2006-01")
	case "yearly":
		return now.Format("2006")
	default: // "none"
		return ""
	}
}

// formatDate 将日期格式化为指定格式
func (s *autocodeService) formatDate(t time.Time, format string) string {
	// 将 YYYYMMDD, YYMM 等格式转换为 Go 的时间格式
	goFormat := format
	goFormat = strings.ReplaceAll(goFormat, "YYYY", "2006")
	goFormat = strings.ReplaceAll(goFormat, "YY", "06")
	goFormat = strings.ReplaceAll(goFormat, "MM", "01")
	goFormat = strings.ReplaceAll(goFormat, "DD", "02")
	goFormat = strings.ReplaceAll(goFormat, "HH", "15")
	goFormat = strings.ReplaceAll(goFormat, "mm", "04")
	goFormat = strings.ReplaceAll(goFormat, "ss", "05")
	return t.Format(goFormat)
}

// buildIdentifier 构建 global_id 的 identifier
func (s *autocodeService) buildIdentifier(tableCode, fieldCode, cycleValue string) string {
	if cycleValue == "" {
		return fmt.Sprintf("%s:%s:", tableCode, fieldCode)
	}
	return fmt.Sprintf("%s:%s:%s", tableCode, fieldCode, cycleValue)
}

// ensureGlobalIdExists 确保 global_id 记录存在
func (s *autocodeService) ensureGlobalIdExists(identifier string, start int, tableCode, fieldCode string) error {
	// 尝试获取现有记录
	where := map[string]any{"identifier": identifier}
	globalIds, err := s.globalIdService.List(1, 1, new(int64), where)
	if err != nil {
		return err
	}

	// 如果已存在则返回
	if len(globalIds) > 0 {
		return nil
	}

	// 创建新记录
	globalId := &model.GlobalId{
		Identifier:  identifier,
		LastID:      uint(start - 1), // LastID 是上一个使用的ID,所以要减1
		Step:        1,
		Description: fmt.Sprintf("自动编码: %s.%s", tableCode, fieldCode),
	}

	return s.globalIdService.Create(globalId)
}
