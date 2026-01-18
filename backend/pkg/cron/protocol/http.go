package client

import (
	"encoding/json"
	"io"
	"net/http"

	"piemdm/internal/model"
	"piemdm/internal/service"
)

type Http struct {
	job           *model.Cron
	params        []*model.CronParam
	entityService service.EntityService
}

func NewHttp(job *model.Cron, params []*model.CronParam, entityService service.EntityService) Protocol {
	return &Http{
		job:           job,
		params:        params,
		entityService: entityService,
	}
}

func (s *Http) Run() {
	err := s.Get(1)
	if err != nil {
		logger.Error("Get", "err", err)
	}
}

func (s *Http) Get(page float64) error {
	// 发送请求
	resp, err := http.Get("http://localhost:8088/ping")
	if err != nil {
		logger.Error("Get", "err", err)
	}
	defer resp.Body.Close()

	// 读取返回结果
	body, _ := io.ReadAll(resp.Body)

	// 序列化返回结果
	var data interface{}
	json.Unmarshal(body, &data)

	// fmt.Printf("\ns.job: %#v\n\n", s.job)
	return nil
}
