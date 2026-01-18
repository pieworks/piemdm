package handler

import (
	"net/http"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type CronHandler interface {
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

type cronHandler struct {
	*Handler
	cronService service.CronService
}

func NewCronHandler(handler *Handler, cronService service.CronService) CronHandler {
	return &cronHandler{
		Handler:     handler,
		cronService: cronService,
	}
}

// List 获取定时任务列表
// @Summary 获取定时任务列表
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param code query string false "任务编码"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {array} model.Cron
// @Router /admin/crons [get]
func (h *cronHandler) List(c *gin.Context) {
	var req struct {
		Page      int    `form:"page,default=1"`
		PageSize  int    `form:"pageSize,default=15"`
		Code      string `form:"code"`
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

	if req.Code != "" {
		where["code"] = req.Code
	}
	if req.StartDate != "" {
		where["created_at >="] = req.StartDate
	}
	if req.EndDate != "" {
		where["created_at <="] = req.EndDate
	}

	crons, err := h.cronService.List(page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, crons)
}

// Get 获取定时任务详情
// @Summary 获取定时任务详情
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} model.Cron
// @Router /admin/crons/{id} [get]
func (h *cronHandler) Get(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cron, err := h.cronService.Get(params.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, cron)
}

// Create 创建定时任务
// @Summary 创建定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param data body object true "创建任务请求"
// @Success 200 {object} model.Cron
// @Router /admin/crons [post]
func (h *cronHandler) Create(c *gin.Context) {
	var req struct {
		Code        string `binding:"required,max=8"` // Cron编码
		Expression  string `binding:"max=32"`         // Cron表达式
		Name        string `binding:"max=128"`        // 任务名称
		EntityCode  string `binding:"max=64"`         // 实体编码
		System      string `binding:"max=64"`         // 系统名称
		Url         string `binding:"max=255"`
		Protocol    string `binding:"max=32" ` // 协议类型，Http, Rest, GraphQL, GRPC, Soap, Jwt
		Method      string `binding:"max=16" `
		AppId       string `binding:"max=64" `
		AppKey      string `binding:"max=64" `
		SignType    string `binding:"max=64" `
		Description string `binding:"max=255"`
		Status      string `binding:"oneof=Normal Frozen Deleted"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cron := model.Cron{
		Code:        req.Code,
		Expression:  req.Expression,
		Name:        req.Name,
		EntityCode:  req.EntityCode,
		System:      req.System,
		Url:         req.Url,
		Protocol:    req.Protocol,
		Method:      req.Method,
		AppId:       req.AppId,
		AppKey:      req.AppKey,
		SignType:    req.SignType,
		Description: req.Description,
		Status:      req.Status,
	}

	err := h.cronService.Create(c, &cron)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, cron)
}

// Update 更新定时任务
// @Summary 更新定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param data body object true "更新任务请求"
// @Success 200 {object} model.Cron
// @Router /admin/crons/{id} [put]
func (h *cronHandler) Update(c *gin.Context) {
	var req struct {
		ID          uint   `binding:"required"`
		Code        string `binding:"required,max=8"` // Cron编码
		Expression  string `binding:"max=32"`         // Cron表达式
		Name        string `binding:"max=128"`        // 任务名称
		EntityCode  string `binding:"max=64"`         // 实体编码
		System      string `binding:"max=64"`         // 系统名称
		Url         string `binding:"max=255"`
		Protocol    string `binding:"max=32" ` // 协议类型，Http, Rest, GraphQL, GRPC, Soap, Jwt
		Method      string `binding:"max=16" `
		AppId       string `binding:"max=64" `
		AppKey      string `binding:"max=64" `
		SignType    string `binding:"max=64" `
		Description string `binding:"max=255"`
		Status      string `binding:"oneof=Normal Frozen Deleted"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cron := model.Cron{
		ID:          req.ID,
		Code:        req.Code,
		Expression:  req.Expression,
		Name:        req.Name,
		EntityCode:  req.EntityCode,
		System:      req.System,
		Url:         req.Url,
		Protocol:    req.Protocol,
		Method:      req.Method,
		AppId:       req.AppId,
		AppKey:      req.AppKey,
		SignType:    req.SignType,
		Description: req.Description,
		Status:      req.Status,
	}
	err := h.cronService.Update(c, &cron)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, cron)
}

// Delete 删除定时任务
// @Summary 删除定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/crons/{id} [delete]
func (h *cronHandler) Delete(c *gin.Context) {
	var req struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.cronService.Delete(c, req.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

func (h *cronHandler) BatchCreate(c *gin.Context) {
	id := c.Param("id")
	resp.HandleSuccess(c, gin.H{"message": "Batch Update user", "id": id})
}

// BatchUpdate 批量更新定时任务状态
// @Summary 批量更新定时任务状态
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param data body object true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/crons/batch [put]
func (h *cronHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		Ids    []uint `form:"ids" binding:"required"`
		Status string `form:"status" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var cron model.Cron
	cron.Status = req.Status

	err := h.cronService.BatchUpdate(c, req.Ids, &cron)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// BatchDelete 批量删除定时任务
// @Summary 批量删除定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param data body object true "批量删除请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/crons/batch [delete]
func (h *cronHandler) BatchDelete(c *gin.Context) {
	var req struct {
		Ids []uint `form:"ids" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	err := h.cronService.BatchDelete(c, req.Ids)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}
