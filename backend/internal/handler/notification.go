package handler

import (
	"net/http"
	"strconv"
	"time"

	"piemdm/internal/repository"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

// NotificationHandler 通知处理器接口
type NotificationHandler interface {
	// 通知发送
	SendApproval(ctx *gin.Context)
	SendBatch(ctx *gin.Context)
	SendByTemplate(ctx *gin.Context)
	Test(ctx *gin.Context)

	// 统计
	GetStatistics(ctx *gin.Context)

	// 处理
	ProcessPending(ctx *gin.Context)
	ProcessRetry(ctx *gin.Context)
}

// notificationHandler 通知处理器实现
type notificationHandler struct {
	*Handler
	notificationService service.NotificationService
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(
	handler *Handler,
	notificationService service.NotificationService,
) NotificationHandler {
	return &notificationHandler{
		Handler:             handler,
		notificationService: notificationService,
	}
}

// 通知发送相关API

// SendApproval 发送审批通知
// @Summary 发送审批通知
// @Description 发送审批相关的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param notification body service.ApprovalNotificationRequest true "通知信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/approval [post]
func (h *notificationHandler) SendApproval(ctx *gin.Context) {
	var req service.ApprovalNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("参数绑定失败", "error", err)
		resp.HandleError(ctx, http.StatusBadRequest, "参数格式错误", nil)
		return
	}

	err := h.notificationService.SendApproval(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("发送审批通知失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, map[string]string{"message": "通知发送成功"})
}

// SendBatch 批量发送通知
// @Summary 批量发送通知
// @Description 批量发送多个通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param notifications body []service.ApprovalNotificationRequest true "通知列表"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/batch [post]
func (h *notificationHandler) SendBatch(ctx *gin.Context) {
	var reqs []*service.ApprovalNotificationRequest
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		h.logger.Error("参数绑定失败", "error", err)
		resp.HandleError(ctx, http.StatusBadRequest, "参数格式错误", nil)
		return
	}

	err := h.notificationService.SendBatch(ctx.Request.Context(), reqs)
	if err != nil {
		h.logger.Error("批量发送通知失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, map[string]string{"message": "批量通知发送成功"})
}

// SendByTemplate 使用模板发送通知
// @Summary 使用模板发送通知
// @Description 使用指定模板发送通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param notification body service.TemplateNotificationRequest true "模板通知信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/template [post]
func (h *notificationHandler) SendByTemplate(ctx *gin.Context) {
	var req service.TemplateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("参数绑定失败", "error", err)
		resp.HandleError(ctx, http.StatusBadRequest, "参数格式错误", nil)
		return
	}

	err := h.notificationService.SendByTemplate(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("使用模板发送通知失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, map[string]string{"message": "模板通知发送成功"})
}

// Test 测试通知
// @Summary 测试通知
// @Description 测试通知渠道是否正常工作
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param test body service.TestNotificationRequest true "测试信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/test [post]
func (h *notificationHandler) Test(ctx *gin.Context) {
	var req service.TestNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("参数绑定失败", "error", err)
		resp.HandleError(ctx, http.StatusBadRequest, "参数格式错误", nil)
		return
	}

	err := h.notificationService.Test(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("测试通知失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, map[string]string{"message": "测试通知发送成功"})
}

// GetStatistics 获取通知统计
// @Summary 获取通知统计
// @Description 获取通知发送统计信息
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Param groupBy query string false "分组方式" Enums(day,week,month)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/statistics [get]
func (h *notificationHandler) GetStatistics(ctx *gin.Context) {
	req := &repository.NotificationStatisticsRequest{
		NotificationType: ctx.Query("notificationType"),
		RecipientType:    ctx.Query("recipientType"),
	}

	// 解析时间参数
	if startTimeStr := ctx.Query("startTime"); startTimeStr != "" {
		if startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
			req.StartTime = &startTime
		}
	}
	if endTimeStr := ctx.Query("endTime"); endTimeStr != "" {
		if endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
			req.EndTime = &endTime
		}
	}

	result, err := h.notificationService.GetStatistics(ctx.Request.Context(), req)
	if err != nil {
		h.logger.Error("获取通知统计失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, result)
}

// ProcessPending 处理待发送通知
// @Summary 处理待发送通知
// @Description 手动触发处理待发送的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param limit query int false "处理数量限制" default(100)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/process-pending [post]
func (h *notificationHandler) ProcessPending(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	limit := 100 // 默认值
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	err := h.notificationService.ProcessPending(ctx.Request.Context(), limit)
	if err != nil {
		h.logger.Error("处理待发送通知失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, map[string]string{"message": "待发送通知处理完成"})
}

// ProcessRetry 处理重试通知
// @Summary 处理重试通知
// @Description 处理重试发送失败的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param limit query int false "重试数量限制" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/notifications/retry-failed [post]
func (h *notificationHandler) ProcessRetry(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	limit := 50 // 默认值
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	err := h.notificationService.ProcessRetry(ctx.Request.Context(), limit)
	if err != nil {
		h.logger.Error("处理重试通知失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, map[string]string{"message": "重试通知处理完成"})
}
