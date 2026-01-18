package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type CronService interface {
	// Base CRUD
	Get(id uint) (*model.Cron, error)
	Find(sel string, where map[string]any) ([]*model.Cron, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.Cron, error)
	Create(c *gin.Context, cron *model.Cron) error
	Update(c *gin.Context, cron *model.Cron) error
	Delete(c *gin.Context, id uint) error

	// Batch operations
	BatchUpdate(c *gin.Context, ids []uint, cron *model.Cron) error
	BatchDelete(c *gin.Context, ids []uint) error
}

type cronService struct {
	*Service
	cronRepository repository.CronRepository
}

func NewCronService(service *Service, cronRepository repository.CronRepository) CronService {
	return &cronService{
		Service:        service,
		cronRepository: cronRepository,
	}
}

func (s *cronService) Get(id uint) (*model.Cron, error) {
	return s.cronRepository.FindOne(id)
}

func (s *cronService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.Cron, error) {
	return s.cronRepository.FindPage(page, pageSize, total, where)
}

func (s *cronService) Find(sel string, where map[string]any) ([]*model.Cron, error) {
	return s.cronRepository.Find(sel, where)
}

func (s *cronService) Create(c *gin.Context, cron *model.Cron) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// cron.Code = code

	return s.cronRepository.Create(c, cron)
}

func (s *cronService) Update(c *gin.Context, cron *model.Cron) error {
	return s.cronRepository.Update(c, cron)
}

func (s *cronService) BatchUpdate(c *gin.Context, ids []uint, cron *model.Cron) error {
	return s.cronRepository.BatchUpdate(c, ids, cron)
}

func (s *cronService) Delete(c *gin.Context, id uint) error {
	return s.cronRepository.Delete(c, id)
}

func (s *cronService) BatchDelete(c *gin.Context, ids []uint) error {
	return s.cronRepository.BatchDelete(c, ids)
}
