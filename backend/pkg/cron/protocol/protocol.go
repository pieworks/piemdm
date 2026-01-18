package client

import (
	"log/slog"
	"piemdm/internal/model"
	"piemdm/internal/service"
)

var logger = slog.Default()

// API is interface
type Protocol interface {
	Get(page float64) error
	Run()
}

// NewAPI return Api instance by type
func NewProtocol(protocolType string, job *model.Cron, param []*model.CronParam, entityService service.EntityService) Protocol {
	switch protocolType {
	case "Http":
		return NewHttp(job, param, entityService)
	}
	return nil
}
