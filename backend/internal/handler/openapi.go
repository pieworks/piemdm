package handler

import (
	"net/http"
	"strconv"

	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"
	"piemdm/pkg/log"

	"github.com/gin-gonic/gin"
)

type OpenApiHandler interface {
	// 简单列表(仅支持基础筛选)
	List(c *gin.Context)

	// 详情
	Get(c *gin.Context)

	// 创建 (Phase 1 暂不实现)
	// Create(c *gin.Context)

	// 全量更新 (Phase 1 暂不实现)
	// Update(c *gin.Context)

	// 局部更新 (Phase 1 暂不实现)
	// Patch(c *gin.Context)

	// 删除 (Phase 1 暂不实现)
	// Delete(c *gin.Context)
}

type openApiHandler struct {
	logger        *log.Logger
	entityService service.EntityService
	entityRepo    repository.EntityRepository
}

func NewOpenApiHandler(
	logger *log.Logger,
	entityService service.EntityService,
	entityRepo repository.EntityRepository,
) OpenApiHandler {
	return &openApiHandler{
		logger:        logger,
		entityService: entityService,
		entityRepo:    entityRepo,
	}
}

// List 查询实体列表
// @Summary 查询实体列表 (OpenAPI)
// @Description 通过 OpenAPI 查询指定实体的列表数据
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Param table path string true "实体表名"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /openapi/v1/entities/{table} [get]
// @Security ApiKeyAuth
func (h *openApiHandler) List(c *gin.Context) {
	// 从 URL 参数获取表代码
	tableCode := c.Param("table")
	if tableCode == "" {
		resp.HandleError(c, http.StatusBadRequest, "PARAM_REQUIRED_MISSING: table is required", nil)
		return
	}

	// 从 Context 获取已认证的 Application
	app, exists := c.Get("application")
	if !exists {
		resp.HandleError(c, http.StatusUnauthorized, "AUTH_FAILED: application not found in context", nil)
		return
	}

	application, ok := app.(*model.Application)
	if !ok {
		h.logger.Error("Invalid application type in context")
		resp.HandleError(c, http.StatusInternalServerError, "SYSTEM_INTERNAL_ERROR", nil)
		return
	}

	h.logger.Debug("OpenAPI List request", "app_id", application.AppId, "table", tableCode)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "15"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}

	// 构建查询条件
	where := make(map[string]any)

	// 支持基础筛选参数 (id, status, created_at 等)
	if id := c.Query("id"); id != "" {
		where["id"] = id
	}
	if status := c.Query("status"); status != "" {
		where["status"] = status
	}
	if createdAt := c.Query("created_at"); createdAt != "" {
		where["created_at"] = createdAt
	}

	var total int64

	// 调用 EntityRepository 直接获取数据 (OpenAPI 不需要用户权限检查)
	entities, err := h.entityRepo.FindPage(tableCode, page, pageSize, &total, where)
	if err != nil {
		h.logger.Error("Failed to list entities", "error", err, "table", tableCode)
		resp.HandleError(c, http.StatusInternalServerError, "SYSTEM_INTERNAL_ERROR: "+err.Error(), nil)
		return
	}

	// 生成分页链接
	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	// 返回成功响应
	resp.HandleSuccess(c, gin.H{
		"data":     entities,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// Get 查询实体详情
// @Summary 查询实体详情 (OpenAPI)
// @Description 通过 OpenAPI 查询指定实体的详情数据
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Param table path string true "实体表名"
// @Param id path int true "记录 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /openapi/v1/entities/{table}/{id} [get]
// @Security ApiKeyAuth
func (h *openApiHandler) Get(c *gin.Context) {
	// 从 URL 参数获取表代码和 ID
	tableCode := c.Param("table")
	idStr := c.Param("id")

	if tableCode == "" {
		resp.HandleError(c, http.StatusBadRequest, "PARAM_REQUIRED_MISSING: table is required", nil)
		return
	}

	if idStr == "" {
		resp.HandleError(c, http.StatusBadRequest, "PARAM_REQUIRED_MISSING: id is required", nil)
		return
	}

	// 解析 ID
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.HandleError(c, http.StatusBadRequest, "PARAM_TYPE_INVALID: id must be a valid integer", nil)
		return
	}

	// 从 Context 获取已认证的 Application
	app, exists := c.Get("application")
	if !exists {
		resp.HandleError(c, http.StatusUnauthorized, "AUTH_FAILED: application not found in context", nil)
		return
	}

	application, ok := app.(*model.Application)
	if !ok {
		h.logger.Error("Invalid application type in context")
		resp.HandleError(c, http.StatusInternalServerError, "SYSTEM_INTERNAL_ERROR", nil)
		return
	}

	h.logger.Debug("OpenAPI Get request", "app_id", application.AppId, "table", tableCode, "id", id)

	// 调用 EntityRepository 直接获取数据 (OpenAPI 不需要用户权限检查)
	entity, err := h.entityRepo.FindOne(tableCode, uint(id))
	if err != nil {
		h.logger.Error("Failed to get entity", "error", err, "table", tableCode, "id", id)

		// 根据错误类型返回不同的错误码
		if err.Error() == "record not found" {
			resp.HandleError(c, http.StatusNotFound, "DATA_NOT_FOUND: entity not found", nil)
			return
		}

		resp.HandleError(c, http.StatusInternalServerError, "SYSTEM_INTERNAL_ERROR: "+err.Error(), nil)
		return
	}

	// 返回成功响应
	resp.HandleSuccess(c, entity)
}
