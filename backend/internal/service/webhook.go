package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type WebhookService interface {
	// Base CRUD
	Get(id uint) (*model.Webhook, error)
	Find(sel string, where map[string]any) ([]*model.Webhook, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.Webhook, error)
	Create(c *gin.Context, webhook *model.Webhook) error
	Update(c *gin.Context, webhook *model.Webhook) error
	Delete(c *gin.Context, id uint) error

	// Batch operations
	BatchUpdate(c *gin.Context, ids []uint, webhook *model.Webhook) error
	BatchDelete(c *gin.Context, ids []uint) error
}

type webhookService struct {
	*Service
	webhookRepository repository.WebhookRepository
}

func NewWebhookService(service *Service, webhookRepository repository.WebhookRepository) WebhookService {
	return &webhookService{
		Service:           service,
		webhookRepository: webhookRepository,
	}
}

func (s *webhookService) Get(id uint) (*model.Webhook, error) {
	return s.webhookRepository.FindOne(id)
}

func (s *webhookService) Find(sel string, where map[string]any) ([]*model.Webhook, error) {
	return s.webhookRepository.Find(sel, where)
}

func (s *webhookService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.Webhook, error) {
	return s.webhookRepository.FindPage(page, pageSize, total, where)
}

func (s *webhookService) Create(c *gin.Context, webhook *model.Webhook) error {
	return s.webhookRepository.Create(c, webhook)
}

func (s *webhookService) Update(c *gin.Context, webhook *model.Webhook) error {
	return s.webhookRepository.Update(c, webhook)
}

func (s *webhookService) BatchUpdate(c *gin.Context, ids []uint, webhook *model.Webhook) error {
	return s.webhookRepository.BatchUpdate(c, ids, webhook)
}

func (s *webhookService) Delete(c *gin.Context, id uint) error {
	return s.webhookRepository.Delete(c, id)
}

func (s *webhookService) BatchDelete(c *gin.Context, ids []uint) error {
	return s.webhookRepository.BatchDelete(c, ids)
}
