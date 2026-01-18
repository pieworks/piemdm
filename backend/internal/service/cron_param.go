package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"
)

type CronParamService interface {
	// Base CRUD
	Get(id uint) (*model.CronParam, error)
	Find(sel string, where map[string]any) ([]*model.CronParam, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.CronParam, error)
	Create(cronParam *model.CronParam) error
	Update(cronParam *model.CronParam) error
	Delete(id uint) (*model.CronParam, error)

	// Batch operations
	BatchUpdate(ids []uint, cronParam *model.CronParam) error
	BatchDelete(ids []uint) error
}

type cronParamService struct {
	*Service
	cronParamRepository repository.CronParamRepository
}

func NewCronParamService(service *Service, cronParamRepository repository.CronParamRepository) CronParamService {
	return &cronParamService{
		Service:             service,
		cronParamRepository: cronParamRepository,
	}
}

func (s *cronParamService) Get(id uint) (*model.CronParam, error) {
	return s.cronParamRepository.FindOne(id)
}

func (s *cronParamService) Find(sel string, where map[string]any) ([]*model.CronParam, error) {
	return s.cronParamRepository.Find(sel, where)
}

func (s *cronParamService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.CronParam, error) {
	return s.cronParamRepository.FindPage(page, pageSize, total, where)
}

func (s *cronParamService) Create(cronParam *model.CronParam) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// cronParam.Code = code

	return s.cronParamRepository.Create(cronParam)
}

func (s *cronParamService) Update(cronParam *model.CronParam) error {
	return s.cronParamRepository.Update(cronParam)
}

func (s *cronParamService) BatchUpdate(ids []uint, cronParam *model.CronParam) error {
	return s.cronParamRepository.BatchUpdate(ids, cronParam)
}

func (s *cronParamService) Delete(id uint) (*model.CronParam, error) {
	return s.cronParamRepository.Delete(id)
}

func (s *cronParamService) BatchDelete(ids []uint) error {
	return s.cronParamRepository.BatchDelete(ids)
}
