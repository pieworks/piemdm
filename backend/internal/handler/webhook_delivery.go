package handler

import (
	"net/http"
	"strconv"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type WebhookDeliveryHandler interface {
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

type webhookDeliveryHandler struct {
	*Handler
	webhookDeliveryService service.WebhookDeliveryService
	tablePermissionService service.TablePermissionService
}

func NewWebhookDeliveryHandler(handler *Handler, webhookDeliveryService service.WebhookDeliveryService, tablePermissionService service.TablePermissionService) WebhookDeliveryHandler {
	return &webhookDeliveryHandler{
		Handler:                handler,
		webhookDeliveryService: webhookDeliveryService,
		tablePermissionService: tablePermissionService,
	}
}

// List 获取Webhook投递记录列表
// @Summary 获取Webhook投递记录列表
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param deliveryCode query string false "投递代码"
// @Param table_code query string false "表代码"
// @Param entity_id query string false "实体ID"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {array} model.WebhookDelivery
// @Router /admin/webhook_deliveries [get]
func (h *webhookDeliveryHandler) List(c *gin.Context) {
	var req struct {
		DeliveryCode string `form:"deliveryCode"`
		TableCode    string `form:"table_code"`
		StartDate    string `form:"startDate"`
		EndDate      string `form:"endDate"`
		EntityId     string `form:"entity_id"`
		Page         int    `form:"page,default=1"`
		PageSize     int    `form:"pageSize,default=15"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Permission check for table_code
	if req.TableCode != "" {
		userIDStr := c.GetString("user_id")
		if userIDStr != "" {
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				h.logger.Error("Failed to parse user_id", "user_id", userIDStr, "error", err)
				resp.HandleError(c, http.StatusInternalServerError, "Invalid user ID", nil)
				return
			}

			allowedTableCodes, err := h.tablePermissionService.GetAllowedTableCodes(c, uint(userID))
			if err != nil {
				h.logger.Error("Failed to get allowed table codes", "error", err)
				resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
				return
			}

			// If not superuser/admin, check table permission
			if allowedTableCodes != nil {
				hasPermission := false
				for _, code := range allowedTableCodes {
					if code == req.TableCode {
						hasPermission = true
						break
					}
				}
				if !hasPermission {
					h.logger.Warn("User has no permission to access table", "user_id", userID, "table_code", req.TableCode)
					resp.HandleError(c, http.StatusForbidden, "No permission to access this table", nil)
					return
				}
			}
		}
	}

	page := req.Page
	pageSize := req.PageSize
	where := make(map[string]any)
	var total int64

	if req.DeliveryCode != "" {
		where["delivery_code"] = req.DeliveryCode
	}
	if req.StartDate != "" {
		where["created_at >="] = req.StartDate
	}
	if req.EndDate != "" {
		where["created_at <="] = req.EndDate
	}
	if req.EntityId != "" {
		where["entity_id"] = req.EntityId
	}
	webhooks, err := h.webhookDeliveryService.List(page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, webhooks)
}

// Get 获取Webhook投递记录详情
// @Summary 获取Webhook投递记录详情
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param id path int true "投递记录ID"
// @Success 200 {object} model.WebhookDelivery
// @Router /admin/webhook_deliveries/{id} [get]
func (h *webhookDeliveryHandler) Get(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	webhook, err := h.webhookDeliveryService.Get(params.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, webhook)
}

// Create 创建Webhook投递记录
// @Summary 创建Webhook投递记录
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param data body object true "投递记录信息"
// @Success 200 {object} model.WebhookDelivery
// @Router /admin/webhook_deliveries [post]
func (h *webhookDeliveryHandler) Create(c *gin.Context) {
	var req struct {
		HookID          uint   `binding:"required"`
		DeliveryCode    string `binding:"required,max=64"`
		Event           string `binding:"max=32"`
		EntityID        uint   `binding:"required"`
		RequestHeaders  string `binding:"max=256"`
		RequestPayload  string `binding:"max=256"`
		ResponseStatus  int    `binding:"max=999"`
		ResponseMessage string `binding:"max=256"`
		ResponseHeaders string `binding:"max=256"`
		ResponseBody    string `binding:"max=256"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	now := time.Now()
	webhookDelivery := model.WebhookDelivery{
		HookID:          req.HookID,
		DeliveryCode:    req.DeliveryCode,
		Event:           req.Event,
		EntityID:        req.EntityID,
		RequestHeaders:  req.RequestHeaders,
		RequestPayload:  req.RequestPayload,
		ResponseStatus:  req.ResponseStatus,
		ResponseMessage: req.ResponseMessage,
		ResponseHeaders: req.ResponseHeaders,
		ResponseBody:    req.ResponseBody,
		DeliveredAt:     &now,
		CompletedAt:     &now,
	}

	err := h.webhookDeliveryService.Create(&webhookDelivery)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, webhookDelivery)
}

// Update 更新Webhook投递记录
// @Summary 更新Webhook投递记录
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param id path int true "投递记录ID"
// @Param data body object true "投递记录信息"
// @Success 200 {object} map[string]interface{}
// @Router /admin/webhook_deliveries/{id} [put]
func (h *webhookDeliveryHandler) Update(c *gin.Context) {
	var req struct {
		ID              uint   `binding:"required"`
		HookID          uint   `binding:"required"`
		DeliveryCode    string `binding:"required,max=64"`
		Event           string `binding:"max=32"`
		EntityID        uint   `binding:"required"`
		RequestHeaders  string `binding:"max=256"`
		RequestPayload  string `binding:"max=256"`
		ResponseStatus  int    `binding:"max=999"`
		ResponseMessage string `binding:"max=256"`
		ResponseHeaders string `binding:"max=256"`
		ResponseBody    string `binding:"max=256"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	webhookDelivery := model.WebhookDelivery{
		ID:              req.ID,
		HookID:          req.HookID,
		DeliveryCode:    req.DeliveryCode,
		Event:           req.Event,
		EntityID:        req.EntityID,
		RequestHeaders:  req.RequestHeaders,
		RequestPayload:  req.RequestPayload,
		ResponseStatus:  req.ResponseStatus,
		ResponseMessage: req.ResponseMessage,
		ResponseHeaders: req.ResponseHeaders,
		ResponseBody:    req.ResponseBody,
	}

	err := h.webhookDeliveryService.Update(&webhookDelivery)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

// Delete 删除Webhook投递记录
// @Summary 删除Webhook投递记录
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param id path int true "投递记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/webhook_deliveries/{id} [delete]
func (h *webhookDeliveryHandler) Delete(c *gin.Context) {
	var req struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.webhookDeliveryService.Delete(req.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

func (h *webhookDeliveryHandler) BatchCreate(c *gin.Context) {
	resp.HandleSuccess(c, nil)
}

// BatchUpdate 批量更新Webhook投递记录
// @Summary 批量更新Webhook投递记录
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param data body object true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/webhook_deliveries/batch [put]
func (h *webhookDeliveryHandler) BatchUpdate(c *gin.Context) {
	var params struct {
		Ids    []uint `form:"ids" binding:"required"`
		Action string `form:"action" binding:"required"`
	}
	if err := c.ShouldBind(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var webhookDelivery model.WebhookDelivery
	err := h.webhookDeliveryService.BatchUpdate(params.Ids, &webhookDelivery)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// BatchDelete 批量删除Webhook投递记录
// @Summary 批量删除Webhook投递记录
// @Tags Webhook投递记录
// @Accept json
// @Produce json
// @Param data body object true "批量删除请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/webhook_deliveries/batch [delete]
func (h *webhookDeliveryHandler) BatchDelete(c *gin.Context) {
	var req struct {
		Ids []uint `form:"ids" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	err := h.webhookDeliveryService.BatchDelete(req.Ids)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}
