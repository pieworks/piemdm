package handler

import (
	"net/http"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type RoleHandler interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	BatchCreate(c *gin.Context)
	BatchUpdate(c *gin.Context)
	BatchDelete(c *gin.Context)

	// Permission management
	GetRolePermissions(c *gin.Context)
	AssignPermissions(c *gin.Context)
	RemovePermissions(c *gin.Context)
	UpdatePermissions(c *gin.Context)

	// User management
	GetRoleUsers(c *gin.Context)
	UpdateRoleUsers(c *gin.Context)
}

type roleHandler struct {
	*Handler
	roleService service.RoleService
}

func NewRoleHandler(handler *Handler, roleService service.RoleService) RoleHandler {
	return &roleHandler{
		Handler:     handler,
		roleService: roleService,
	}
}

// List 获取角色列表
// @Summary 获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param code query string false "角色编码"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {array} model.Role
// @Router /admin/roles [get]
func (h *roleHandler) List(c *gin.Context) {
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

	roles, err := h.roleService.List(page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, roles)
}

// Get 获取角色详情
// @Summary 获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} model.Role
// @Router /admin/roles/{id} [get]
func (h *roleHandler) Get(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	role, err := h.roleService.Get(params.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, role)
}

// Create 创建角色
// @Summary 创建角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body object true "创建角色请求"
// @Success 200 {object} model.Role
// @Router /admin/roles [post]
func (h *roleHandler) Create(c *gin.Context) {
	var req struct {
		Code        string `binding:"required"`
		Name        string `binding:"required"`
		Description string `binding:"max=255"`
		Status      string `binding:"oneof=Normal Frozen Deleted"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	role := model.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	err := h.roleService.Create(c, &role)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, role)
}

// Update 更新角色
// @Summary 更新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body object true "更新角色请求"
// @Success 200 {object} model.Role
// @Router /admin/roles/{id} [put]
func (h *roleHandler) Update(c *gin.Context) {
	var req struct {
		ID          uint   `binding:"required"`
		Code        string `binding:"required"`
		Name        string `binding:"required"`
		Description string `binding:"max=255"`
		Status      string `binding:"oneof=Normal Frozen Deleted"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	role := model.Role{
		ID:          req.ID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	err := h.roleService.Update(c, &role)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, role)
}

// Delete 删除角色
// @Summary 删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/{id} [delete]
func (h *roleHandler) Delete(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.roleService.Delete(c, params.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

func (h *roleHandler) BatchCreate(c *gin.Context) {
}

// BatchUpdate 批量更新角色状态
// @Summary 批量更新角色状态
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body object true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/batch [put]
func (h *roleHandler) BatchUpdate(c *gin.Context) {
	var params struct {
		IDs    []uint `form:"ids"`
		Status string `form:"status"`
	}
	if err := c.ShouldBind(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var role model.Role
	role.Status = params.Status

	err := h.roleService.BatchUpdate(c, params.IDs, &role)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
	// resp.HandleSuccess(c, nil)
}

// BatchDelete 批量删除角色
// @Summary 批量删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body object true "批量删除请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/batch [delete]
func (h *roleHandler) BatchDelete(c *gin.Context) {
	var params struct {
		IDs []uint `form:"ids"`
	}
	if err := c.ShouldBind(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	err := h.roleService.BatchDelete(c, params.IDs)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// GetRolePermissions 获取角色的权限列表
// @Summary 获取角色的权限列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {array} model.Permission
// @Router /admin/roles/{id}/permissions [get]
func (h *roleHandler) GetRolePermissions(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	permissions, err := h.roleService.GetRolePermissions(c, params.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, permissions)
}

// AssignPermissions 为角色分配权限
// @Summary 为角色分配权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body object true "权限ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/{id}/permissions [post]
func (h *roleHandler) AssignPermissions(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.roleService.AssignPermissions(c, params.Id, req.PermissionIDs)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

// RemovePermissions 移除角色的权限
// @Summary 移除角色的权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body object true "权限ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/{id}/permissions [delete]
func (h *roleHandler) RemovePermissions(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.roleService.RemovePermissions(c, params.Id, req.PermissionIDs)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

// UpdatePermissions 更新角色的权限
// @Summary 更新角色的权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body object true "权限ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/{id}/permissions [put]
func (h *roleHandler) UpdatePermissions(c *gin.Context) {
	var params struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.roleService.UpdatePermissions(c, params.Id, req.PermissionIDs)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

// GetRoleUsers 获取角色的用户
// GetRoleUsers 获取角色的用户列表
// @Summary 获取角色的用户列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {array} model.User
// @Router /admin/roles/{id}/users [get]
func (h *roleHandler) GetRoleUsers(c *gin.Context) {
	var params struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	users, err := h.roleService.GetRoleUsers(params.ID)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, users)
}

// UpdateRoleUsers 更新角色的用户
// UpdateRoleUsers 更新角色的用户列表
// @Summary 更新角色的用户列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body object true "用户ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /admin/roles/{id}/users [put]
func (h *roleHandler) UpdateRoleUsers(c *gin.Context) {
	var params struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.roleService.UpdateRoleUsers(params.ID, req.UserIDs); err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, gin.H{"message": "success"})
}
