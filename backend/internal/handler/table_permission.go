package handler

import (
	"net/http"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type TablePermissionHandler interface {
	List(c *gin.Context)
	Create(c *gin.Context)
	BatchUpdate(c *gin.Context) // Replaces Freeze/Unfreeze
	BatchDelete(c *gin.Context) // Replaces Delete
}

type tablePermissionHandler struct {
	*Handler
	service service.TablePermissionService
}

func NewTablePermissionHandler(handler *Handler, service service.TablePermissionService) TablePermissionHandler {
	return &tablePermissionHandler{
		Handler: handler,
		service: service,
	}
}

// List 获取角色的表权限列表
// @Summary 获取角色的表权限列表
// @Tags 表权限管理
// @Accept json
// @Produce json
// @Param role_id query int true "角色ID"
// @Success 200 {array} object
// @Router /admin/table_permissions [get]
func (h *tablePermissionHandler) List(c *gin.Context) {
	var params struct {
		RoleID uint `form:"role_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	perms, err := h.service.GetByRoleID(c, params.RoleID)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, perms)
}

// Create 授予表权限
// @Summary 授予表权限
// @Tags 表权限管理
// @Accept json
// @Produce json
// @Param data body object true "权限信息"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_permissions [post]
func (h *tablePermissionHandler) Create(c *gin.Context) {
	var params struct {
		RoleID    uint   `json:"role_id" binding:"required"`
		TableCode string `json:"table_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.service.Grant(c, params.RoleID, params.TableCode); err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, nil)
}

// BatchUpdate 批量更新表权限状态
// @Summary 批量更新表权限状态
// @Tags 表权限管理
// @Accept json
// @Produce json
// @Param data body object true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_permissions/batch [put]
func (h *tablePermissionHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		IDs    []uint `form:"ids" json:"ids"`
		Status string `form:"status" json:"status"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	values := make(map[string]any)
	if req.Status != "" {
		values["status"] = req.Status
	}

	if len(values) == 0 {
		resp.HandleError(c, http.StatusBadRequest, "no values to update", nil)
		return
	}

	if err := h.service.BatchUpdate(c, req.IDs, values); err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// BatchDelete 批量删除表权限
// @Summary 批量删除表权限
// @Tags 表权限管理
// @Accept json
// @Produce json
// @Param data body object true "批量删除请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_permissions/batch [delete]
func (h *tablePermissionHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `form:"ids" json:"ids"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.service.BatchDelete(c, req.IDs); err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}
