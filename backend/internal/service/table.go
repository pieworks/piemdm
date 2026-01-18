package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type TableService interface {
	Get(id uint) (*model.Table, error)
	Find(sel string, where map[string]any) ([]*model.Table, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.Table, error)
	Create(c *gin.Context, table *model.Table) error
	Update(c *gin.Context, table *model.Table) error
	BatchUpdate(c *gin.Context, ids []uint, table *model.Table) error
	Delete(c *gin.Context, id uint) (*model.Table, error)
	BatchDelete(c *gin.Context, ids []uint) error
}

type tableService struct {
	*Service
	tableRepository      repository.TableRepository
	tableFieldRepository repository.TableFieldRepository
}

func NewTableService(service *Service, tableRepository repository.TableRepository, tableFieldRepository repository.TableFieldRepository) TableService {
	return &tableService{
		Service:              service,
		tableRepository:      tableRepository,
		tableFieldRepository: tableFieldRepository,
	}
}

func (s *tableService) Get(id uint) (*model.Table, error) {
	return s.tableRepository.FindOne(id)
}

func (s *tableService) Find(sel string, where map[string]any) ([]*model.Table, error) {
	return s.tableRepository.Find(sel, where)
}

func (s *tableService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.Table, error) {
	return s.tableRepository.FindPage(page, pageSize, total, where)
}

func (s *tableService) Create(c *gin.Context, table *model.Table) error {
	return s.tableRepository.Create(c, table)
}

func (s *tableService) Update(c *gin.Context, table *model.Table) error {
	// if table.Type == "Extension" {
	// 	// 1. 检查table_field表中是否存在 LocalField，LocalPrimary。改成直接查询两个字段。
	// 	// fieldsWhere := map[string]any{}
	// 	// fieldsWhere["table_code"] = table.Code
	// 	// fieldsWhere["code"] = []string{table.LocalField, table.LocalPrimary}
	// 	// fieldsWhere["status"] = "Normal"
	// 	// fields, err := s.tableFieldRepository.Find("", fieldsWhere)
	// 	// if err != nil {
	// 	// 	return nil
	// 	// }
	// 	// locals := map[string]int{}
	// 	// for _, field := range fields {
	// 	// 	if field.Code == table.LocalField {
	// 	// 		locals["LocalField"] = 1
	// 	// 	}
	// 	// 	if field.Code == table.LocalPrimary {
	// 	// 		locals["LocalPrimary"] = 1
	// 	// 	}
	// 	// }
	// 	// 2. LocalField 使用 BasicField 信息，如果存在了不增加
	// 	localField := model.TableField{
	// 		TableCode: table.Code,
	// 		Code:      table.LocalField,
	// 	}
	// 	lfield, err := s.tableFieldRepository.First("", localField)
	// 	localFieldIndexName := strings.ToLower("idx_" + table.Code + "_" + table.LocalField + "_" + table.LocalPrimary)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if lfield != nil {
	// 		lfield.IndexName = localFieldIndexName
	// 		err = s.tableFieldRepository.Update(lfield)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		localField.IndexName = localFieldIndexName
	// 		err = s.tableFieldRepository.Create(&localField)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	// 3. LocalPrimary 如果不存在增加一个64位长度的默认文本字段，如果存在了不增加。
	// 	localPrimary := model.TableField{
	// 		TableCode: table.Code,
	// 		Code:      table.LocalPrimary,
	// 	}
	// 	pfield, err := s.tableFieldRepository.First("", localPrimary)
	// 	localPrimaryIndexName := strings.ToLower("idx_" + table.Code + "_" + table.LocalField + "_" + table.LocalPrimary)
	// 	if pfield != nil {
	// 		pfield.IndexName = localPrimaryIndexName
	// 		err = s.tableFieldRepository.Update(pfield)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		localPrimary.IndexName = localPrimaryIndexName
	// 		err = s.tableFieldRepository.Create(&localPrimary)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	// 4. 创建 LocalField，LocalPrimary 联合索引
	// 	// 只创建相同的索引名称，需要发布模型才会更新数据库。
	// 	// TODO可能模型并没有发布，最好通过模型发布创建联合索引
	// 	// tableField 表增加索引名称字段，不用对外展示，通过索引名称字段创建联合索引。
	// 	// 这样每次发布的时候就可以跟着发布模型，同时创建索引。
	// }
	return s.tableRepository.Update(c, table)
	// return nil
}

func (s *tableService) BatchUpdate(c *gin.Context, ids []uint, table *model.Table) error {
	return s.tableRepository.BatchUpdate(c, ids, table)
}

func (s *tableService) Delete(c *gin.Context, id uint) (*model.Table, error) {
	return s.tableRepository.Delete(c, id)
}

func (s *tableService) BatchDelete(c *gin.Context, ids []uint) error {
	return s.tableRepository.BatchDelete(c, ids)
}
