package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"
)

type GlobalIdService interface {
	// Base CRUD
	Get(id uint) (*model.GlobalId, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.GlobalId, error)
	Create(globalId *model.GlobalId) error
	Update(globalId *model.GlobalId) error
	Delete(id uint) (*model.GlobalId, error)

	// Batch operations
	BatchUpdate(ids []uint, globalId *model.GlobalId) error
	BatchDelete(ids []uint) error

	// 业务方法
	GetNewID(identifier string) uint
}

type globalIdService struct {
	*Service
	globalIdRepository repository.GlobalIdRepository
}

func NewGlobalIdService(service *Service, globalIdRepository repository.GlobalIdRepository) GlobalIdService {
	return &globalIdService{
		Service:            service,
		globalIdRepository: globalIdRepository,
	}
}

func (s *globalIdService) Get(id uint) (*model.GlobalId, error) {
	return s.globalIdRepository.FindOne(id)
}

func (s *globalIdService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.GlobalId, error) {
	return s.globalIdRepository.FindPage(page, pageSize, total, where)
}

func (s *globalIdService) Create(globalId *model.GlobalId) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// globalId.Code = code

	return s.globalIdRepository.Create(globalId)
}

func (s *globalIdService) Update(globalId *model.GlobalId) error {
	return s.globalIdRepository.Update(globalId)
}

func (s *globalIdService) BatchUpdate(ids []uint, globalId *model.GlobalId) error {
	return s.globalIdRepository.BatchUpdate(ids, globalId)
}

func (s *globalIdService) Delete(id uint) (*model.GlobalId, error) {
	return s.globalIdRepository.Delete(id)
}

func (s *globalIdService) BatchDelete(ids []uint) error {
	return s.globalIdRepository.BatchDelete(ids)
}

// 获取新id
func (s *globalIdService) GetNewID(identifier string) uint {
	return s.globalIdRepository.GetNewID(identifier)
}
