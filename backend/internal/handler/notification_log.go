package handler

import (
	"net/http"
	"time"

	"piemdm/internal/repository"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

// NotificationLogHandler 通知日志处理器接口
type NotificationLogHandler interface {
	// 基础查询
	List(ctx *gin.Context)
	Get(ctx *gin.Context)
}

// notificationLogHandler 通知日志处理器实现
type notificationLogHandler struct {
	*Handler
	notificationLogService service.NotificationLogService
}

// NewNotificationLogHandler 创建通知日志处理器
func NewNotificationLogHandler(
	handler *Handler,
	notificationLogService service.NotificationLogService,
) NotificationLogHandler {
	return &notificationLogHandler{
		Handler:                handler,
		notificationLogService: notificationLogService,
	}
}

// List 获取通知日志列表
// @Summary 获取通知日志列表
// @Description 分页获取通知日志列表
// @Tags 通知日志管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param approvalId query string false "审批ID"
// @Param taskId query string false "任务ID"
// @Param recipientId query string false "接收人ID"
// @Param notificationType query string false "通知类型"
// @Param status query string false "发送状态"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/notification_logs [get]
func (h *notificationLogHandler) List(ctx *gin.Context) {
	page, pageSize := GetPage(ctx)

	req := &repository.ListNotificationLogRequest{
		Page:             page,
		PageSize:         pageSize,
		ApprovalID:       ctx.Query("approvalId"),
		TaskID:           ctx.Query("taskId"),
		RecipientID:      ctx.Query("recipientId"),
		NotificationType: ctx.Query("notificationType"),
		Status:           ctx.Query("status"),
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

	result, err := h.notificationLogService.List(ctx.Request.Context(), req)
	if err != nil {
		h.logger.Error("获取通知日志列表失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 生成分页链接
	links := resp.GeneratePaginationLinks(ctx.Request, result.Page, result.PageSize, int(result.Total))
	ctx.Header("Link", links.String())

	resp.HandleSuccess(ctx, result.Data)
}

// Get 获取通知日志详情
// @Summary 获取通知日志详情
// @Description 根据ID获取通知日志详情
// @Tags 通知日志管理
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} piemdm_internal_model.NotificationLog
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/notification_logs/{id} [get]
func (h *notificationLogHandler) Get(ctx *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, "日志ID不能为空", nil)
		return
	}

	log, err := h.notificationLogService.Get(ctx.Request.Context(), req.ID)
	if err != nil {
		h.logger.Error("获取通知日志详情失败", "error", err)
		resp.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, log)
}
