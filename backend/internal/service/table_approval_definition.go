package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type TableApprovalDefinitionService interface {
	List(entityCode, operation string) ([]model.TableApprovalDefinition, error)
	Get(id uint) (*model.TableApprovalDefinition, error)
	Create(c *gin.Context, item *model.TableApprovalDefinition) error
	Update(c *gin.Context, item *model.TableApprovalDefinition) error
	Delete(c *gin.Context, id uint) error
	BatchCreate(c *gin.Context, list []model.TableApprovalDefinition) error
	BatchDelete(c *gin.Context, ids []uint) error
}

type tableApprovalDefinitionService struct {
	*Service
	tableApprovalDefinitionRepository repository.TableApprovalDefinitionRepository
}

func NewTableApprovalDefinitionService(service *Service, tableApprovalDefinitionRepository repository.TableApprovalDefinitionRepository) TableApprovalDefinitionService {
	return &tableApprovalDefinitionService{
		Service:                           service,
		tableApprovalDefinitionRepository: tableApprovalDefinitionRepository,
	}
}

func (s *tableApprovalDefinitionService) List(entityCode, operation string) ([]model.TableApprovalDefinition, error) {
	return s.tableApprovalDefinitionRepository.List(entityCode, operation)
}

func (s *tableApprovalDefinitionService) Get(id uint) (*model.TableApprovalDefinition, error) {
	return s.tableApprovalDefinitionRepository.Get(id)
}

func (s *tableApprovalDefinitionService) Create(c *gin.Context, item *model.TableApprovalDefinition) error {
	return s.tableApprovalDefinitionRepository.Create(c, item)
}

func (s *tableApprovalDefinitionService) Update(c *gin.Context, item *model.TableApprovalDefinition) error {
	return s.tableApprovalDefinitionRepository.Update(c, item)
}

func (s *tableApprovalDefinitionService) Delete(c *gin.Context, id uint) error {
	return s.tableApprovalDefinitionRepository.Delete(c, id)
}

func (s *tableApprovalDefinitionService) BatchCreate(c *gin.Context, list []model.TableApprovalDefinition) error {
	return s.tableApprovalDefinitionRepository.BatchCreate(c, list)
}

func (s *tableApprovalDefinitionService) BatchDelete(c *gin.Context, ids []uint) error {
	return s.tableApprovalDefinitionRepository.BatchDelete(c, ids)
}
