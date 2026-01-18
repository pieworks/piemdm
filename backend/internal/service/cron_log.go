package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"
)

type CronLogService interface {
	// Base CRUD
	Get(id uint) (*model.CronLog, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.CronLog, error)
	Create(cronLog *model.CronLog) error
	Update(cronLog *model.CronLog) error
	Delete(id uint) (*model.CronLog, error)

	// Batch operations
	BatchUpdate(ids []uint, cronLog *model.CronLog) error
	BatchDelete(ids []uint) error
}

type cronLogService struct {
	*Service
	cronLogRepository repository.CronLogRepository
}

func NewCronLogService(service *Service, cronLogRepository repository.CronLogRepository) CronLogService {
	return &cronLogService{
		Service:           service,
		cronLogRepository: cronLogRepository,
	}
}

func (s *cronLogService) Get(id uint) (*model.CronLog, error) {
	return s.cronLogRepository.FindOne(id)
}

func (s *cronLogService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.CronLog, error) {
	return s.cronLogRepository.FindPage(page, pageSize, total, where)
}

func (s *cronLogService) Create(cronLog *model.CronLog) error {
	return s.cronLogRepository.Create(cronLog)
}

func (s *cronLogService) Update(cronLog *model.CronLog) error {
	return s.cronLogRepository.Update(cronLog)
}

func (s *cronLogService) BatchUpdate(ids []uint, cronLog *model.CronLog) error {
	return s.cronLogRepository.BatchUpdate(ids, cronLog)
}

func (s *cronLogService) Delete(id uint) (*model.CronLog, error) {
	return s.cronLogRepository.Delete(id)
}

func (s *cronLogService) BatchDelete(ids []uint) error {
	return s.cronLogRepository.BatchDelete(ids)
}
