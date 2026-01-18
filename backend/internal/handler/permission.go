package handler

import (
	"net/http"
	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

// PermissionHandler 权限处理器接口
type PermissionHandler interface {
	List(c *gin.Context)
	GetTree(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// permissionHandler 权限处理器实现
type permissionHandler struct {
	*Handler
	permissionService service.PermissionService
}

// NewPermissionHandler 创建权限处理器实例
func NewPermissionHandler(handler *Handler, permissionService service.PermissionService) PermissionHandler {
	return &permissionHandler{
		Handler:           handler,
		permissionService: permissionService,
	}
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required,max=64"`
	Name        string `json:"name" binding:"required,max=64"`
	Resource    string `json:"resource" binding:"max=64"`
	Action      string `json:"action" binding:"max=64"`
	ParentID    uint   `json:"parent_id"`
	Description string `json:"description" binding:"max=255"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Code        string `json:"code" binding:"required,max=64"`
	Name        string `json:"name" binding:"required,max=64"`
	Resource    string `json:"resource" binding:"max=64"`
	Action      string `json:"action" binding:"max=64"`
	ParentID    uint   `json:"parent_id"`
	Description string `json:"description" binding:"max=255"`
}

// Create 创建权限
// @Summary 创建权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param data body CreatePermissionRequest true "创建权限请求"
// @Success 200 {object} model.Permission
// @Router /admin/permissions [post]
func (h *permissionHandler) Create(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	permission := &model.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		ParentID:    req.ParentID,
		Description: req.Description,
		Status:      "Normal",
	}

	if err := h.permissionService.Create(c.Request.Context(), permission); err != nil {
		h.logger.Error("创建权限失败: " + err.Error())
		resp.HandleError(c, http.StatusInternalServerError, "创建权限失败", err.Error())
		return
	}

	resp.HandleSuccess(c, permission)
}

// Update 更新权限
// @Summary 更新权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param data body UpdatePermissionRequest true "更新权限请求"
// @Success 200 {object} model.Permission
// @Router /admin/permissions/{id} [put]
func (h *permissionHandler) Update(c *gin.Context) {
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	permission := &model.Permission{
		ID:          req.ID,
		Code:        req.Code,
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		ParentID:    req.ParentID,
		Description: req.Description,
	}

	if err := h.permissionService.Update(c.Request.Context(), permission); err != nil {
		h.logger.Error("更新权限失败: " + err.Error())
		resp.HandleError(c, http.StatusInternalServerError, "更新权限失败", err.Error())
		return
	}

	resp.HandleSuccess(c, permission)
}

// Delete 删除权限
// @Summary 删除权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/permissions/{id} [delete]
func (h *permissionHandler) Delete(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.permissionService.Delete(c.Request.Context(), params.Id); err != nil {
		h.logger.Error("删除权限失败: " + err.Error())
		resp.HandleError(c, http.StatusInternalServerError, "删除权限失败", err.Error())
		return
	}

	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// Get 获取单个权限
// @Summary 获取权限详情
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} model.Permission
// @Router /admin/permissions/{id} [get]
func (h *permissionHandler) Get(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	permission, err := h.permissionService.Get(c.Request.Context(), params.Id)
	if err != nil {
		h.logger.Error("获取权限失败: " + err.Error())
		resp.HandleError(c, http.StatusNotFound, "权限不存在", err.Error())
		return
	}

	resp.HandleSuccess(c, permission)
}

// List 获取权限列表
// @Summary 获取权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param code query string false "权限编码"
// @Param resource query string false "资源类型"
// @Success 200 {array} model.Permission
// @Router /admin/permissions [get]
func (h *permissionHandler) List(c *gin.Context) {
	var req struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"pageSize,default=15"`
		Code     string `form:"code"`
		Resource string `form:"resource"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	page := req.Page
	pageSize := req.PageSize

	permissions, total, err := h.permissionService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("获取权限列表失败: " + err.Error())
		resp.HandleError(c, http.StatusInternalServerError, "获取权限列表失败", err.Error())
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, permissions)
}

// GetTree 获取权限树
// @Summary 获取权限树形结构
// @Tags 权限管理
// @Accept json
// @Produce json
// @Success 200 {array} model.Permission
// @Router /admin/permissions/tree [get]
func (h *permissionHandler) GetTree(c *gin.Context) {
	tree, err := h.permissionService.GetTree(c.Request.Context())
	if err != nil {
		h.logger.Error("获取权限树失败: " + err.Error())
		resp.HandleError(c, http.StatusInternalServerError, "获取权限树失败", err.Error())
		return
	}

	resp.HandleSuccess(c, tree)
}
