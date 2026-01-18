package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type EntityLogService interface {
	// Base CRUD
	Get(tableCode string, id uint) (*model.EntityLog, error)
	List(tableCode string, page, pageSize int, total *int64, where map[string]any) ([]*model.EntityLog, error)
	Create(c *gin.Context, tableCode string, entityLog *model.EntityLog) error
	Update(c *gin.Context, tableCode string, entityLog *model.EntityLog) error
	Delete(c *gin.Context, tableCode string, id uint) (*model.EntityLog, error)

	// Batch operations
	BatchUpdate(c *gin.Context, tableCode string, ids []uint, entityLog *model.EntityLog) error
	BatchDelete(c *gin.Context, tableCode string, ids []uint) error
}

type entityLogService struct {
	*Service
	entityLogRepository repository.EntityLogRepository
}

func NewEntityLogService(service *Service, entityLogRepository repository.EntityLogRepository) EntityLogService {
	return &entityLogService{
		Service:             service,
		entityLogRepository: entityLogRepository,
	}
}

func (s *entityLogService) Get(tableCode string, id uint) (*model.EntityLog, error) {
	return s.entityLogRepository.FindOne(tableCode, id)
}

func (s *entityLogService) List(tableCode string, page, pageSize int, total *int64, where map[string]any) ([]*model.EntityLog, error) {
	return s.entityLogRepository.FindPage(tableCode, page, pageSize, total, where)
}

func (s *entityLogService) Create(c *gin.Context, tableCode string, entityLog *model.EntityLog) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// entityLog.Code = code

	return s.entityLogRepository.Create(c, tableCode, entityLog)
}

func (s *entityLogService) Update(c *gin.Context, tableCode string, entityLog *model.EntityLog) error {
	return s.entityLogRepository.Update(c, tableCode, entityLog)
}

func (s *entityLogService) BatchUpdate(c *gin.Context, tableCode string, ids []uint, entityLog *model.EntityLog) error {
	return s.entityLogRepository.BatchUpdate(c, tableCode, ids, entityLog)
}

func (s *entityLogService) Delete(c *gin.Context, tableCode string, id uint) (*model.EntityLog, error) {
	return s.entityLogRepository.FindOne(tableCode, id)
}

func (s *entityLogService) BatchDelete(c *gin.Context, tableCode string, ids []uint) error {
	return s.entityLogRepository.BatchDelete(c, tableCode, ids)
}
