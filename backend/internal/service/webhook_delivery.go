package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"
)

type WebhookDeliveryService interface {
	// Base CRUD
	Get(id uint) (*model.WebhookDelivery, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.WebhookDelivery, error)
	Create(webhookDelivery *model.WebhookDelivery) error
	Update(webhookDelivery *model.WebhookDelivery) error
	Delete(id uint) error

	// Batch operations
	BatchUpdate(ids []uint, webhookDelivery *model.WebhookDelivery) error
	BatchDelete(ids []uint) error
}

type webhookDeliveryService struct {
	*Service
	webhookDeliveryRepository repository.WebhookDeliveryRepository
}

func NewWebhookDeliveryService(service *Service, webhookDeliveryRepository repository.WebhookDeliveryRepository) WebhookDeliveryService {
	return &webhookDeliveryService{
		Service:                   service,
		webhookDeliveryRepository: webhookDeliveryRepository,
	}
}

func (s *webhookDeliveryService) Get(id uint) (*model.WebhookDelivery, error) {
	return s.webhookDeliveryRepository.FindOne(id)
}

func (s *webhookDeliveryService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.WebhookDelivery, error) {
	return s.webhookDeliveryRepository.FindPage(page, pageSize, total, where)
}

func (s *webhookDeliveryService) Create(webhookDelivery *model.WebhookDelivery) error {
	// uuid := uuid.New()
	// code := strings.ToUpper(uuid.String())
	// webhookDelivery.Code = code

	return s.webhookDeliveryRepository.Create(webhookDelivery)
}

func (s *webhookDeliveryService) Update(webhookDelivery *model.WebhookDelivery) error {
	return s.webhookDeliveryRepository.Update(webhookDelivery)
}

func (s *webhookDeliveryService) BatchUpdate(ids []uint, webhookDelivery *model.WebhookDelivery) error {
	return s.webhookDeliveryRepository.BatchUpdate(ids, webhookDelivery)
}

func (s *webhookDeliveryService) Delete(id uint) error {
	return s.webhookDeliveryRepository.Delete(id)
}

func (s *webhookDeliveryService) BatchDelete(ids []uint) error {
	return s.webhookDeliveryRepository.BatchDelete(ids)
}
