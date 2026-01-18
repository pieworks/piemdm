package repository

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"piemdm/internal/model"

	"github.com/gin-gonic/gin"
)

type TableFieldRepository interface {
	First(sel string, tableField model.TableField) (*model.TableField, error)
	FindOne(id uint) (*model.TableField, error)
	Find(sel string, where map[string]any) ([]*model.TableField, error)
	FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.TableField, error)
	Create(c *gin.Context, tableField *model.TableField) error
	Update(c *gin.Context, tableField *model.TableField) error
	BatchUpdate(c *gin.Context, ids []uint, tableField *model.TableField) error
	Delete(c *gin.Context, id uint) (*model.TableField, error)
	BatchDelete(c *gin.Context, ids []uint) error
	Public(tableName string, entity any) error
	BuildEntity(tableCode string) map[string]any
	GetTableOptions(tableCode string, filter map[string]any) ([]map[string]any, error)
}
type tableFieldRepository struct {
	*Repository
	source Base
}

func NewTableFieldRepository(repository *Repository, source Base) TableFieldRepository {
	return &tableFieldRepository{
		Repository: repository,
		source:     source,
	}
}

func (r *tableFieldRepository) First(sel string, tableField model.TableField) (*model.TableField, error) {
	if sel == "" {
		sel = "*"
	}
	if err := r.db.First(&tableField, []string{}).Error; err != nil {
		return nil, err
	}
	return &tableField, nil
}

func (r *tableFieldRepository) FindOne(id uint) (*model.TableField, error) {
	var tableField model.TableField
	if err := r.db.First(&tableField, id).Error; err != nil {
		return nil, err
	}
	return &tableField, nil
}

func (r *tableFieldRepository) Find(sel string, where map[string]any) ([]*model.TableField, error) {
	var tableFields []*model.TableField
	var tableField model.TableField
	if sel == "" {
		sel = "*"
	}

	// err := r.source.FindPage(tableField, &tableFields, 1, 20, total, where, []string{}, "ID desc")
	// preloads := []string{}

	err := r.source.Find(tableField, &tableFields, sel, where, "sort asc")
	if err != nil {
		r.logger.Error("获取模型信息失败", "err", err)
	}
	return tableFields, nil
}

func (r *tableFieldRepository) FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.TableField, error) {
	var tableFields []*model.TableField
	var tableField model.TableField

	// if err := r.db.Where("deleted_at is null").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tableFields).Error; err != nil {
	// 	return nil, err
	// }

	// if err := r.db.Model(tableField).Where("deleted_at is null").Count(total).Error; err != nil {
	// 	r.logger.Error("apptoval repository count err", "err", err)
	// 	return nil, err
	// }
	preloads := []string{}
	err := r.source.FindPage(tableField, &tableFields, page, pageSize, total, where, preloads, "ID desc")
	if err != nil {
		r.logger.Error("获取模型信息失败", "err", err)
	}
	return tableFields, nil
}

func (r *tableFieldRepository) Create(c *gin.Context, tableField *model.TableField) error {
	if err := r.db.WithContext(c).Create(tableField).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableFieldRepository) Update(c *gin.Context, tableField *model.TableField) error {
	// 忽略 Code, Type, TableCode, FieldType 字段，防止修改
	if err := r.db.WithContext(c).Model(tableField).Omit("Code", "Type", "TableCode", "FieldType").Updates(tableField).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableFieldRepository) BatchUpdate(c *gin.Context, ids []uint, tableField *model.TableField) error {
	// 忽略 Code, Type, TableCode, FieldType 字段
	if err := r.db.WithContext(c).Model(&tableField).Where("id in ?", ids).Omit("Code", "Type", "TableCode", "FieldType").Updates(tableField).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableFieldRepository) Delete(c *gin.Context, id uint) (*model.TableField, error) {
	var tableField model.TableField
	if err := r.db.WithContext(c).Where("id = ?", id).First(&tableField).Error; err != nil {
		return nil, err
	}
	return &tableField, nil
}

func (r *tableFieldRepository) BatchDelete(c *gin.Context, ids []uint) error {
	var tableFields []model.TableField
	// 使用下面的调用方式，才能把参数传到 Hooks 里面。
	// //单条删除
	// db.Where("id =1").Find(&tableField).Delete(&tableField)
	// //多条删除
	// db.Where("id in ?", ids).Find(&tableFields).Delete(&tableFields)

	// if err := r.db.Delete(&tableField, ids).Error; err != nil {
	if err := r.db.WithContext(c).Where("id in ?", ids).Find(&tableFields).Delete(&tableFields).Error; err != nil {
		return err
	}
	return nil
}

func (r *tableFieldRepository) Public(tableName string, entity any) error {
	// 如果entity是model.EntityLog,使用GORM的AutoMigrate
	if _, ok := entity.(model.EntityLog); ok {
		if err := r.db.Table(tableName).AutoMigrate(entity); err != nil {
			return err
		}
		return nil
	}

	// 对于map类型,构建临时结构体并使用AutoMigrate
	if _, ok := entity.(map[string]any); !ok {
		return fmt.Errorf("unsupported entity type: %T", entity)
	}

	// 构建用于AutoMigrate的临时结构体
	migrationStruct, err := r.buildMigrationStruct(tableName)
	if err != nil {
		return fmt.Errorf("failed to build migration struct: %v", err)
	}

	// 使用GORM AutoMigrate同步表结构
	if err = r.db.Table(tableName).AutoMigrate(migrationStruct); err != nil {
		return fmt.Errorf("failed to migrate table: %v", err)
	}

	// GORM AutoMigrate 不会修改已存在列的 NOT NULL 约束
	// 需要手动移除业务字段的 NOT NULL 约束
	if err := r.removeNotNullConstraints(tableName); err != nil {
		r.logger.Warn("移除 NOT NULL 约束时出现警告", "table", tableName, "error", err)
		// 不中断流程,继续执行
	}

	r.logger.Info("表结构同步完成", "table", tableName)
	return nil
}

// buildMigrationStruct 构建用于AutoMigrate的临时结构体
func (r *tableFieldRepository) buildMigrationStruct(tableName string) (any, error) {
	// 从tableName提取tableCode
	tableCode := strings.TrimPrefix(tableName, "t_")
	tableCode = strings.TrimSuffix(tableCode, "_draft")
	tableCode = strings.TrimSuffix(tableCode, "_log")

	// 查询table_field定义
	where := make(map[string]any)
	where["status"] = "Normal"
	where["table_code"] = tableCode

	tableFields, err := r.Find("code,name,type,length,required,is_index,is_unique,index_name,index_priority,options", where)
	if err != nil {
		return nil, fmt.Errorf("failed to get table fields: %v", err)
	}

	// 使用reflect动态构建结构体字段
	var structFields []reflect.StructField

	// 添加ID主键
	structFields = append(structFields, reflect.StructField{
		Name: "ID",
		Type: reflect.TypeOf(uint(0)),
		Tag:  reflect.StructTag(`gorm:"primaryKey;column:id"`),
	})

	// 添加业务字段
	for _, field := range tableFields {
		fieldName := r.toCamelCase(field.Code)
		var fieldType reflect.Type
		var gormTag string

		// MySQL数据类型 ---------------------------------------------------------
		// 类型	大小	范围（有符号）	范围（无符号）	用途
		// TINYINT	1 字节	(-128,127)	(0,255)	小整数值
		// SMALLINT	2 字节	(-32768,32767)	(0,65535)	大整数值
		// MEDIUMINT	3 字节	(-8388608,8388607)	(0,16777215)	大整数值
		// INT或INTEGER	4 字节	(-2147483648,2147483647)	(0,4294967295)	大整数值
		// BIGINT	8 字节	(+-9223372036854775807)	(0,18446744073709551615)	极大整数值
		// FLOAT	4 字节	-(+-38位),(0,+(+-38位))	单精度浮点数值
		// DOUBLE	8 字节	-(+-308位),(0,+(+-308位))	双精度浮点数值
		// DECIMAL	对DECIMAL(M,D) ,如果M>D,为M+2否则为D+2	小数值
		// Go数据类型 ---------------------------------------------------------
		// 类型	符号	描述
		// uint8	无符号	8位整型 (0 到 255)
		// uint16	无符号	6位整型 (0 到 65535)
		// uint32	无符号	32位整型 (0 到 4294967295)
		// uint64	无符号	64位整型 (0 到 18446744073709551615)
		// int8	有符号	8位整型 (-128 到 127)
		// int16	有符号	16位整型 (-32768 到 32767)
		// int32	有符号	32位整型 (-2147483648 到 2147483647)
		// int64	有符号	64位整型 (-9223372036854775808 到 9223372036854775807)
		// uint	无符号	特殊 按操作系统位数，是 uint32 或 uint64，不能将 uint 赋值给 uint32 或 uint64
		// int	有符号	特殊 按操作系统位数，是 int32 或 int64，但不能将其赋值给 int32 或 int64 变量
		// uintptr	无符号	无符号整型，用于存放一个指针
		// 根据类型确定Go类型
		switch field.Type {
		case "Date":
			// 日期类型,只存储日期部分
			fieldType = reflect.TypeOf(time.Time{})
			gormTag = fmt.Sprintf(`gorm:"column:%s;type:date;comment:%s"`, field.Code, field.Name)
		case "DateTime":
			// 日期时间类型,存储日期+时间
			fieldType = reflect.TypeOf(time.Time{})
			gormTag = fmt.Sprintf(`gorm:"column:%s;type:datetime(3);comment:%s"`, field.Code, field.Name)
		case "Number":
			// 检查是否为 decimal 类型（通过 options.validation 中的 precision 和 scale 判断）
			if field.Options != nil && field.Options.Validation != nil &&
				field.Options.Validation.Precision != nil && field.Options.Validation.Scale != nil {
				// decimal 类型使用 float64
				fieldType = reflect.TypeOf(float64(0))
				precision := *field.Options.Validation.Precision
				scale := *field.Options.Validation.Scale
				gormTag = fmt.Sprintf(`gorm:"column:%s;type:decimal(%d,%d);default:0;comment:%s"`, field.Code, precision, scale, field.Name)
			} else {
				// integer 类型使用 int
				fieldType = reflect.TypeOf(0)
				gormTag = fmt.Sprintf(`gorm:"column:%s;default:0;comment:%s"`, field.Code, field.Name)
			}
		default: // "Text"
			fieldType = reflect.TypeOf("")
			if field.Length > 0 {
				gormTag = fmt.Sprintf(`gorm:"column:%s;size:%d;comment:%s"`, field.Code, field.Length, field.Name)
			} else {
				gormTag = fmt.Sprintf(`gorm:"column:%s;size:255;comment:%s"`, field.Code, field.Name)
			}
		}

		// 生成索引
		// UniqueScope字段 废弃，全部使用表级索引
		// 索引：`gorm:"index"`
		// 索引+名称：`gorm:"index:idx_name"`
		// 唯一索引：`gorm:"uniqueIndex"` 等效 `gorm:"index:,unique"`
		// 唯一索引+名称：`gorm:"uniqueIndex:idx_name"`
		// 字段优先级：`gorm:"index:idx_member,priority:2"` 默认为10，
		// 多索引：`gorm:"index:idx_id;index:idx_oid,unique"`
		// priority 在创建的时候有效，但是修改的时候，数据库并没有修改。
		// data中可以使用唯一索引
		// data_draft中不能使用唯一索引，为审批流程可能会有很多
		// 添加索引标签
		if field.IsIndex == "Yes" && !strings.HasSuffix(tableName, "_draft") {
			if field.IndexName != "" {
				gormTag += fmt.Sprintf(`;index:%s`, field.IndexName)
			}
		}
		if field.IsUnique == "Yes" && !strings.HasSuffix(tableName, "_draft") {
			if field.IndexName != "" {
				gormTag += fmt.Sprintf(`;uniqueIndex:%s`, field.IndexName)
			}
		}

		structFields = append(structFields, reflect.StructField{
			Name: fieldName,
			Type: fieldType,
			Tag:  reflect.StructTag(gormTag),
		})
	}

	// 添加系统字段
	systemFields := []reflect.StructField{
		// ==============================
		// operationType
		// SAP-ECC中习惯使用的标识方法。
		// 新建: Create
		// 修改: Update
		// 冻结: Freeze
		// 解冻: Unfreeze
		// 锁定: Lock
		// 解锁: Unlock
		// 删除: Delete
		// 扩展: Extension
		// 作废: Void
		// 撤销: Cancel
		// 终止: Terminate
		// 批量创建: BatchCreate
		// 批量修改: BatchUpdate
		// 批量冻结: BatchFreeze
		// 批量解冻: BatchUnfreeze
		// 批量锁定: BatchLock
		// 批量解锁: BatchUnlock
		// 批量删除: BatchDeletion
		// 批量扩展: BatchExtension
		// ==============================
		// actionType
		// I插入，U更新，B冻结，C解冻，D删除

		// C  Create Apply = "I";//新增
		//    Extend Apply_Extend = "Ei";//扩展申请
		// U  Update Modify = "U";//修改
		// B  Freeze = "B";//冻结 (注意：历史代码中可能使用 "F" 表示 Freeze)
		// C  Unfreeze = "C";//解冻 (注意：历史代码中可能使用 "UF" 或 "UL" 表示 Unfreeze)
		// D  Delete Delete = "D";//删除
		// U  Lock = "U";//锁定 (注意：历史代码中可能使用 "L" 表示 Lock)
		// U  Unlock = "U";//解锁 (注意：历史代码中可能使用 "UL" 表示 Unlock)
		//    Cancel Terminate = "T";//终止(作废)
		//    Batch_Request Apply_Batch = "Bi";//批量申请
		//    Batch_Updateapply_Update = "Bu";//批量修改
		//    Px_Apply = "Pi";//平行账新增
		{Name: "Operation", Type: reflect.TypeOf(""), Tag: `gorm:"column:operation;size:16"`},
		// 2020-11-04  应SAP要求，会计科目二次分发时actType为I的都应改为U
		// I插入，U更新，B冻结，C解冻，D删除
		{Name: "Action", Type: reflect.TypeOf(""), Tag: `gorm:"column:action;size:4"`},
		// ===分发状态===（数据状态，审批状态，发布状态 三个状态之一。）
		// TODO 分发状态可以放到分发任务里面。
		// 0分发失败,1成功,2未发送
		// 是否改为0未分发，1分发成功，2分发失败(多个分发任务，有失败的任务，分发失败)。
		{Name: "SendStatus", Type: reflect.TypeOf(0), Tag: `gorm:"column:send_status;default:0"`},
	}

	// 如果是draft表,添加draft相关字段
	if strings.HasSuffix(tableName, "_draft") {
		systemFields = append(systemFields,
			reflect.StructField{Name: "EntityId", Type: reflect.TypeOf(uint(0)), Tag: `gorm:"column:entity_id;default:0"`},
			reflect.StructField{Name: "ApprovalCode", Type: reflect.TypeOf(""), Tag: `gorm:"column:approval_code;size:64;index"`},
			// ===审批状态=== （数据状态，审批状态，发布状态 三个状态之一。）
			// Drafted，Pending，Published 在Draft表。是过程状态，根据流程阶段不同形成
			// 一个entity_id只能有一个Pending的，如果有Pending中的不能提交。Drafted状态的提交也要验证这个逻辑
			reflect.StructField{Name: "DraftStatus", Type: reflect.TypeOf(""), Tag: `gorm:"column:draft_status;size:16"`},
		)
	}

	structFields = append(structFields, systemFields...)

	// 添加GORM标准字段
	structFields = append(structFields,
		// -- 数据公共属性
		// Published，Froze，Deleted，Locked在entity表。都是最终状态，根据操作不同形成
		// ===数据状态===（数据状态，审批状态，发布状态 三个状态之一。）
		// 0草稿，1提交，2已发布，3已冻结，6审核拒绝
		// 0 Drafted 保存为草稿，未提交审批流程，审批拒绝可以修改后重新提交。
		// 1 Pending 提交工作流，审批中
		// 2 Published 已发布，审批完成
		// 3 Froze 已冻结，审批完成
		// Deleted 已删除，审批完成
		// --Locked--，暂时不考虑Lock状态
		// 6 --Rejected-- 审核拒绝，审批完成，应该回到草稿状态
		reflect.StructField{Name: "Status", Type: reflect.TypeOf(""), Tag: `gorm:"column:status;size:16"`},
		reflect.StructField{Name: "CreatedBy", Type: reflect.TypeOf(""), Tag: `gorm:"column:created_by;size:64"`},
		reflect.StructField{Name: "UpdatedBy", Type: reflect.TypeOf(""), Tag: `gorm:"column:updated_by;size:64"`},
		reflect.StructField{Name: "CreatedAt", Type: reflect.TypeOf(time.Time{}), Tag: `gorm:"column:created_at"`},
		reflect.StructField{Name: "UpdatedAt", Type: reflect.TypeOf(time.Time{}), Tag: `gorm:"column:updated_at"`},
		reflect.StructField{Name: "DeletedAt", Type: reflect.TypeOf((*time.Time)(nil)), Tag: `gorm:"column:deleted_at;index"`},
	)

	// 创建结构体类型
	structType := reflect.StructOf(structFields)

	// 创建结构体实例
	instance := reflect.New(structType).Interface()

	return instance, nil
}

// removeNotNullConstraints 移除业务字段的 NOT NULL 约束
// GORM AutoMigrate 不会修改已存在列的约束,需要手动执行 ALTER TABLE
func (r *tableFieldRepository) removeNotNullConstraints(tableName string) error {
	// 从tableName提取tableCode
	tableCode := strings.TrimPrefix(tableName, "t_")
	tableCode = strings.TrimSuffix(tableCode, "_draft")
	tableCode = strings.TrimSuffix(tableCode, "_log")

	// 查询table_field定义
	where := make(map[string]any)
	where["status"] = "Normal"
	where["table_code"] = tableCode

	tableFields, err := r.Find("code,name,type,length,options", where)
	if err != nil {
		return fmt.Errorf("failed to get table fields: %v", err)
	}

	// 检查表是否存在
	var tableExists bool
	err = r.db.Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&tableExists).Error
	if err != nil || !tableExists {
		// 表不存在,无需处理
		return nil
	}

	// 为每个业务字段移除 NOT NULL 约束
	for _, field := range tableFields {
		var columnType string

		// 构建正确的列类型定义
		switch field.Type {
		case "Date":
			columnType = "date"
		case "DateTime":
			columnType = "datetime(3)"
		case "Number":
			// 检查是否为 decimal 类型
			if field.Options != nil && field.Options.Validation != nil &&
				field.Options.Validation.Precision != nil && field.Options.Validation.Scale != nil {
				precision := *field.Options.Validation.Precision
				scale := *field.Options.Validation.Scale
				columnType = fmt.Sprintf("decimal(%d,%d) DEFAULT 0", precision, scale)
			} else {
				columnType = "bigint DEFAULT 0"
			}
		case "Text":
			if field.Length > 0 {
				columnType = fmt.Sprintf("varchar(%d)", field.Length)
			} else {
				columnType = "varchar(255)"
			}
		default:
			columnType = "varchar(255)"
		}

		// 添加 COMMENT
		if field.Name != "" {
			columnType += fmt.Sprintf(" COMMENT '%s'", field.Name)
		}

		// 执行 ALTER TABLE MODIFY COLUMN 移除 NOT NULL
		sql := fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s", tableName, field.Code, columnType)
		if err := r.db.Exec(sql).Error; err != nil {
			r.logger.Warn("移除字段 NOT NULL 约束失败", "table", tableName, "field", field.Code, "error", err)
			// 继续处理其他字段
		}
	}

	return nil
}

// toCamelCase 将字段代码转换为合法的 Go 驼峰命名标识符
// 处理下划线(_)、连字符(-)等分隔符,移除非法字符
// 例如: "fie-1222" -> "Fie1222", "user_name" -> "UserName"
func (r *tableFieldRepository) toCamelCase(s string) string {
	// 替换所有非字母数字字符为下划线,统一处理
	// 例如: "fie-1222" -> "fie_1222"
	var normalized strings.Builder
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			normalized.WriteRune(ch)
		} else {
			normalized.WriteRune('_')
		}
	}

	// 按下划线分割并转换为驼峰命名
	parts := strings.Split(normalized.String(), "_")
	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			// 首字母大写
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}

	fieldName := result.String()

	// 确保字段名以字母开头(Go 标识符要求)
	if len(fieldName) > 0 && fieldName[0] >= '0' && fieldName[0] <= '9' {
		fieldName = "Field" + fieldName
	}

	// 如果结果为空,返回默认值
	if fieldName == "" {
		fieldName = "UnknownField"
	}

	return fieldName
}

func (r *tableFieldRepository) BuildEntity(tableCode string) map[string]any {
	entity := make(map[string]any)

	where := make(map[string]any)
	where["status"] = "Normal"
	where["table_code"] = tableCode

	// 动态生成 struct, 如果获取 draft，需要得到 table_code 获取字段信息
	if strings.HasSuffix(tableCode, "_draft") {
		where["table_code"] = tableCode[0 : len(tableCode)-6]
	}
	tableName := fmt.Sprintf("t_%s", tableCode)

	sel := "code,type,length,required,is_index,is_unique,index_name,index_priority,options"
	tableFields, err := r.Find(sel, where)
	if err != nil {
		r.logger.Error("获取表字段失败", "error", err)
	}

	// 添加ID字段
	entity["id"] = uint(0)

	// 添加业务字段
	for _, tableField := range tableFields {
		fieldCode := tableField.Code

		// MySQL数据类型 ---------------------------------------------------------
		// 类型	大小	范围（有符号）	范围（无符号）	用途
		// TINYINT	1 字节	(-128,127)	(0,255)	小整数值
		// SMALLINT	2 字节	(-32768,32767)	(0,65535)	大整数值
		// MEDIUMINT	3 字节	(-8388608,8388607)	(0,16777215)	大整数值
		// INT或INTEGER	4 字节	(-2147483648,2147483647)	(0,4294967295)	大整数值
		// BIGINT	8 字节	(+-9223372036854775807)	(0,18446744073709551615)	极大整数值
		// FLOAT	4 字节	-(+-38位),(0,+(+-38位))	单精度浮点数值
		// DOUBLE	8 字节	-(+-308位),(0,+(+-308位))	双精度浮点数值
		// DECIMAL	对DECIMAL(M,D) ,如果M>D,为M+2否则为D+2	小数值
		// Go数据类型 ---------------------------------------------------------
		// 类型	符号	描述
		// uint8	无符号	8位整型 (0 到 255)
		// uint16	无符号	6位整型 (0 到 65535)
		// uint32	无符号	32位整型 (0 到 4294967295)
		// uint64	无符号	64位整型 (0 到 18446744073709551615)
		// int8	有符号	8位整型 (-128 到 127)
		// int16	有符号	16位整型 (-32768 到 32767)
		// int32	有符号	32位整型 (-2147483648 到 2147483647)
		// int64	有符号	64位整型 (-9223372036854775808 到 9223372036854775807)
		// uint	无符号	特殊 按操作系统位数，是 uint32 或 uint64，不能将 uint 赋值给 uint32 或 uint64
		// int	有符号	特殊 按操作系统位数，是 int32 或 int64，但不能将其赋值给 int32 或 int64 变量
		// uintptr	无符号	无符号整型，用于存放一个指针

		// 根据类型设置默认值
		switch tableField.Type {
		case "Date", "DateTime":
			entity[fieldCode] = time.Time{}
		case "Number":
			// 检查是否为 decimal 类型
			if tableField.Options != nil && tableField.Options.Validation != nil &&
				tableField.Options.Validation.Precision != nil {
				entity[fieldCode] = float64(0)
			} else {
				entity[fieldCode] = 0
			}
		default: // "Text"
			entity[fieldCode] = ""
		}
	}

	// 添加业务系统字段
	// ==============================
	// operationType
	// SAP-ECC中习惯使用的标识方法。
	// 新建: Create
	// 修改: Update
	// 冻结: Freeze
	// 解冻: Unfreeze
	// 锁定: Lock
	// 解锁: Unlock
	// 删除: Delete
	// 扩展: Extension
	// 作废: Void
	// 撤销: Cancel
	// 终止: Terminate
	// 批量创建: BatchCreate
	// 批量修改: BatchUpdate
	// 批量冻结: BatchFreeze
	// 批量解冻: BatchUnfreeze
	// 批量锁定: BatchLock
	// 批量解锁: BatchUnlock
	// 批量删除: BatchDeletion
	// 批量扩展: BatchExtension
	// ==============================
	// actionType
	// I插入，U更新，B冻结，C解冻，D删除

	// C  Create Apply = "I";//新增
	//    Extend Apply_Extend = "Ei";//扩展申请
	// U  Update Modify = "U";//修改
	// B  Freeze = "B";//冻结 (注意：历史代码中可能使用 "F" 表示 Freeze)
	// C  Unfreeze = "C";//解冻 (注意：历史代码中可能使用 "UF" 或 "UL" 表示 Unfreeze)
	// D  Delete Delete = "D";//删除
	// U  Lock = "U";//锁定 (注意：历史代码中可能使用 "L" 表示 Lock)
	// U  Unlock = "U";//解锁 (注意：历史代码中可能使用 "UL" 表示 Unlock)
	//    Cancel Terminate = "T";//终止(作废)
	//    Batch_Request Apply_Batch = "Bi";//批量申请
	//    Batch_Updateapply_Update = "Bu";//批量修改
	//    Px_Apply = "Pi";//平行账新增
	entity["operation"] = ""
	// 2020-11-04  应SAP要求，会计科目二次分发时actType为I的都应改为U
	// I插入，U更新，B冻结，C解冻，D删除
	entity["action"] = ""
	// ===分发状态===（数据状态，审批状态，发布状态 三个状态之一。）
	// TODO 分发状态可以放到分发任务里面。
	// 0分发失败,1成功,2未发送
	// 是否改为0未分发，1分发成功，2分发失败(多个分发任务，有失败的任务，分发失败)。
	entity["send_status"] = 0

	// 添加流程字段(仅draft表)
	if strings.HasSuffix(tableName, "_draft") {
		entity["entity_id"] = uint(0)
		entity["approval_code"] = ""
		// ===审批状态=== （数据状态，审批状态，发布状态 三个状态之一。）
		// Drafted，Pending，Published 在Draft表。是过程状态，根据流程阶段不同形成
		// 一个entity_id只能有一个Pending的，如果有Pending中的不能提交。Drafted状态的提交也要验证这个逻辑
		entity["draft_status"] = ""
	}

	// 添加GORM标准字段
	// -- 数据公共属性
	// Published，Froze，Deleted，Locked在entity表。都是最终状态，根据操作不同形成
	// ===数据状态===（数据状态，审批状态，发布状态 三个状态之一。）
	// 0草稿，1提交，2已发布，3已冻结，6审核拒绝
	// 0 Drafted 保存为草稿，未提交审批流程，审批拒绝可以修改后重新提交。
	// 1 Pending 提交工作流，审批中
	// 2 Published 已发布，审批完成
	// 3 Froze 已冻结，审批完成
	// Deleted 已删除，审批完成
	// --Locked--，暂时不考虑Lock状态
	// 6 --Rejected-- 审核拒绝，审批完成，应该回到草稿状态

	// 对于主表数据来说，status字段是数据状态，
	// 状态：Normal 正常 Frozen 已冻结 Deleted 已删除
	entity["status"] = ""
	entity["created_by"] = ""
	entity["updated_by"] = ""
	entity["created_at"] = time.Now()
	entity["updated_at"] = time.Now()
	entity["deleted_at"] = nil

	return entity
}

// GetTableOptions 获取表的选项列表（用于关联字段下拉）
// filter: 可选的过滤条件,例如 {"dictionary_class": "DIC0045"}
func (r *tableFieldRepository) GetTableOptions(tableCode string, filter map[string]any) ([]map[string]any, error) {
	tableName := fmt.Sprintf("t_%s", tableCode)

	var options []map[string]any
	query := r.db.Table(tableName).
		Select("code, name").
		Where("status = ?", "Normal")

	// 应用过滤条件
	for field, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	err := query.Order("sort ASC, id ASC").Scan(&options).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get table options: %v", err)
	}

	return options, nil
}
