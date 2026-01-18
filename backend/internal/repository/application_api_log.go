package repository

import (
	"piemdm/internal/model"

	"github.com/gin-gonic/gin"
)

type ApplicationApiLogRepository interface {
	// 创建审计日志
	Create(c *gin.Context, log *model.ApplicationApiLog) error

	// 分页查询审计日志
	FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.ApplicationApiLog, error)

	// 根据 RequestID 查询
	FindByRequestID(requestID string) (*model.ApplicationApiLog, error)
}

type applicationApiLogRepository struct {
	*Repository
	source Base
}

func NewApplicationApiLogRepository(repository *Repository, source Base) ApplicationApiLogRepository {
	return &applicationApiLogRepository{
		Repository: repository,
		source:     source,
	}
}

func (r *applicationApiLogRepository) Create(c *gin.Context, log *model.ApplicationApiLog) error {
	if err := r.db.WithContext(c).Create(log).Error; err != nil {
		return err
	}
	return nil
}

func (r *applicationApiLogRepository) FindPage(page, pageSize int, total *int64, where map[string]any) ([]*model.ApplicationApiLog, error) {
	var logs []*model.ApplicationApiLog
	var log model.ApplicationApiLog

	preloads := []string{}
	err := r.source.FindPage(log, &logs, page, pageSize, total, where, preloads, "created_at desc")
	if err != nil {
		r.logger.Error("获取审计日志失败", "err", err)
		return nil, err
	}

	return logs, nil
}

func (r *applicationApiLogRepository) FindByRequestID(requestID string) (*model.ApplicationApiLog, error) {
	var log model.ApplicationApiLog
	if err := r.db.Where("request_id = ?", requestID).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
