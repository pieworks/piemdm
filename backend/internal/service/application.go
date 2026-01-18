package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type ApplicationService interface {
	// Base CRUD
	Get(id uint) (*model.Application, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.Application, error)
	Create(c *gin.Context, application *model.Application) error
	Update(c *gin.Context, application *model.Application) error
	Delete(c *gin.Context, id uint) error

	// Batch operations
	BatchUpdate(c *gin.Context, ids []uint, application *model.Application) error
	BatchDelete(c *gin.Context, ids []uint) error
}

type applicationService struct {
	*Service
	applicationRepository repository.ApplicationRepository
}

func NewApplicationService(service *Service, applicationRepository repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		Service:               service,
		applicationRepository: applicationRepository,
	}
}

func (s *applicationService) Get(id uint) (*model.Application, error) {
	return s.applicationRepository.FindOne(id)
}

func (s *applicationService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.Application, error) {
	return s.applicationRepository.FindPage(page, pageSize, total, where)
}

func (s *applicationService) Create(c *gin.Context, application *model.Application) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// application.Code = code

	return s.applicationRepository.Create(c, application)
}

func (s *applicationService) Update(c *gin.Context, application *model.Application) error {
	return s.applicationRepository.Update(c, application)
}

func (s *applicationService) BatchUpdate(c *gin.Context, ids []uint, application *model.Application) error {
	return s.applicationRepository.BatchUpdate(c, ids, application)
}

func (s *applicationService) Delete(c *gin.Context, id uint) error {
	return s.applicationRepository.Delete(c, id)
}

func (s *applicationService) BatchDelete(c *gin.Context, ids []uint) error {
	return s.applicationRepository.BatchDelete(c, ids)
}
