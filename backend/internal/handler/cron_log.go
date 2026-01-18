package handler

import (
	"net/http"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type CronLogHandler interface {
	// Base CRUD
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)

	// Batch operations
	BatchCreate(c *gin.Context)
	BatchUpdate(c *gin.Context)
	BatchDelete(c *gin.Context)
}

type cronLogHandler struct {
	*Handler
	cronLogService service.CronLogService
}

func NewCronLogHandler(handler *Handler, cronLogService service.CronLogService) CronLogHandler {
	return &cronLogHandler{
		Handler:        handler,
		cronLogService: cronLogService,
	}
}

// List 获取定时任务日志列表
// @Summary 获取定时任务日志列表
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param id query string false "日志ID"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {array} model.CronLog
// @Router /admin/cron_logs [get]
func (h *cronLogHandler) List(c *gin.Context) {
	var req struct {
		Page      int    `form:"page,default=1"`
		PageSize  int    `form:"pageSize,default=15"`
		ID        string `form:"id"`
		StartDate string `form:"startDate"`
		EndDate   string `form:"endDate"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	page := req.Page
	pageSize := req.PageSize
	where := make(map[string]any)
	var total int64

	if req.ID != "" {
		where["id"] = req.ID
	}
	if req.StartDate != "" {
		where["created_at >="] = req.StartDate
	}
	if req.EndDate != "" {
		where["created_at <="] = req.EndDate
	}

	cronLogs, err := h.cronLogService.List(page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, cronLogs)
}

// Get 获取定时任务日志详情
// @Summary 获取定时任务日志详情
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} model.CronLog
// @Router /admin/cron_logs/{id} [get]
func (h *cronLogHandler) Get(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cronLog, err := h.cronLogService.Get(params.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, cronLog)
}

// Create 创建定时任务日志
// @Summary 创建定时任务日志
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param data body object true "创建日志请求"
// @Success 200 {object} model.CronLog
// @Router /admin/cron_logs [post]
func (h *cronLogHandler) Create(c *gin.Context) {
	var req struct {
		Method string `binding:"required,max=32"` // 任务名称
		Param  string `binding:"max=128"`         // 任务参数
		ErrMsg string `binding:"max=512"`         // 错误信息
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cronLog := model.CronLog{
		Method: req.Method,
		Param:  req.Param,
		ErrMsg: req.ErrMsg,
	}

	err := h.cronLogService.Create(&cronLog)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, cronLog)
}

// Update 更新定时任务日志
// @Summary 更新定时任务日志
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Param data body object true "更新日志请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/cron_logs/{id} [put]
func (h *cronLogHandler) Update(c *gin.Context) {
	var req struct {
		ID     uint   `binding:"required"`
		Method string `binding:"required,max=32"` // 任务名称
		Param  string `binding:"max=128"`         // 任务参数
		ErrMsg string `binding:"max=512"`         // 错误信息
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	cronLog := model.CronLog{
		ID:     req.ID,
		Method: req.Method,
		Param:  req.Param,
		ErrMsg: req.ErrMsg,
	}

	err := h.cronLogService.Update(&cronLog)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

// Delete 删除定时任务日志
// @Summary 删除定时任务日志
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/cron_logs/{id} [delete]
func (h *cronLogHandler) Delete(c *gin.Context) {
	resp.HandleSuccess(c, nil)
}

func (h *cronLogHandler) BatchCreate(c *gin.Context) {
	resp.HandleSuccess(c, nil)
}

// BatchUpdate 批量更新定时任务日志状态
// @Summary 批量更新定时任务日志状态
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param data body object true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/cron_logs/batch [put]
func (h *cronLogHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		Ids    []uint `form:"ids" binding:"required"`
		Status string `form:"status" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var cronLog model.CronLog
	cronLog.Status = req.Status

	err := h.cronLogService.BatchUpdate(req.Ids, &cronLog)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// BatchDelete 批量删除定时任务日志
// @Summary 批量删除定时任务日志
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param data body object true "批量删除请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/cron_logs/batch [delete]
func (h *cronLogHandler) BatchDelete(c *gin.Context) {
	var req struct {
		Ids []uint `form:"ids" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	err := h.cronLogService.BatchDelete(req.Ids)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}
