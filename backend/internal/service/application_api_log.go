package service

import (
	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/gin-gonic/gin"
)

type ApplicationApiLogService interface {
	// 创建审计日志
	Create(c *gin.Context, log *model.ApplicationApiLog) error

	// 分页查询审计日志
	List(page, pageSize int, where map[string]any) ([]*model.ApplicationApiLog, int64, error)

	// 根据 RequestID 查询
	GetByRequestID(requestID string) (*model.ApplicationApiLog, error)
}

type applicationApiLogService struct {
	*Service
	repo repository.ApplicationApiLogRepository
}

func NewApplicationApiLogService(
	service *Service,
	repo repository.ApplicationApiLogRepository,
) ApplicationApiLogService {
	return &applicationApiLogService{
		Service: service,
		repo:    repo,
	}
}

func (s *applicationApiLogService) Create(c *gin.Context, log *model.ApplicationApiLog) error {
	return s.repo.Create(c, log)
}

func (s *applicationApiLogService) List(page, pageSize int, where map[string]any) ([]*model.ApplicationApiLog, int64, error) {
	var total int64
	logs, err := s.repo.FindPage(page, pageSize, &total, where)
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (s *applicationApiLogService) GetByRequestID(requestID string) (*model.ApplicationApiLog, error) {
	return s.repo.FindByRequestID(requestID)
}
