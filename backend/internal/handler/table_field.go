package handler

import (
	"encoding/json"
	"net/http"

	"piemdm/internal/constants"
	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type TableFieldHandler interface {
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

	// 特殊操作
	Public(c *gin.Context)
	GetTableFields(c *gin.Context)
	GetFieldTypePresets(c *gin.Context)
	GetFieldTypeGroups(c *gin.Context)
	GetTableOptions(c *gin.Context)
}

type tableFieldHandler struct {
	*Handler
	tableFieldService service.TableFieldService
}

func NewTableFieldHandler(handler *Handler, tableFieldService service.TableFieldService) TableFieldHandler {
	return &tableFieldHandler{
		Handler:           handler,
		tableFieldService: tableFieldService,
	}
}

// List 获取表字段列表
// @Summary 获取表字段列表
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(15)
// @Param table_code query string false "表编码"
// @Success 200 {array} model.TableField
// @Router /admin/table_fields [get]
func (h *tableFieldHandler) List(c *gin.Context) {
	var req struct {
		TableCode string `form:"table_code" binding:"required"`
		Page      int    `form:"page,default=1"`
		PageSize  int    `form:"pageSize,default=15"`
		StartDate string `form:"startDate"`
		EndDate   string `form:"endDate"`
		Code      string `form:"code"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	page := req.Page
	pageSize := req.PageSize
	where := make(map[string]any)
	var total int64

	where["table_code"] = req.TableCode
	if req.Code != "" {
		where["code"] = req.Code
	}
	if req.StartDate != "" {
		where["created_at >="] = req.StartDate
	}
	if req.EndDate != "" {
		where["created_at <="] = req.EndDate
	}

	tableFields, err := h.tableFieldService.List(page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, tableFields)
}

// Get 获取表字段详情
// @Summary 获取表字段详情
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param id path int true "字段ID"
// @Success 200 {object} model.TableField
// @Router /admin/table_fields/{id} [get]
func (h *tableFieldHandler) Get(c *gin.Context) {
	var req struct {
		Id uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	approval, err := h.tableFieldService.Get(req.Id)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, approval)
}

// GetTableFields 获取表的所有字段（包括系统字段）
// GetTableFields 获取表的所有字段(包括系统字段)
// @Summary 获取表的所有字段
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param table_code query string true "表编码"
// @Success 200 {array} model.TableField
// @Router /admin/table_fields/fields [get]
func (h *tableFieldHandler) GetTableFields(c *gin.Context) {
	tableCode := c.Query("table_code")
	if tableCode == "" {
		resp.HandleError(c, http.StatusBadRequest, "table_code参数不能为空", nil)
		return
	}

	fields, err := h.tableFieldService.GetTableFields(tableCode)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, fields)
}

// Create 创建表字段
// @Summary 创建表字段
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param data body object true "创建请求"
// @Success 200 {object} model.TableField
// @Router /admin/table_fields [post]
func (h *tableFieldHandler) Create(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required,max=64"`
		TableCode   string `json:"table_code" binding:"required,max=64"`
		Name        string `json:"name" binding:"required,max=128"`
		FieldType   string `json:"field_type" binding:"required,max=32"` // 必填,用于推断 Type/Length
		IsUnique    string `json:"is_unique" binding:"max=8"`
		Description string `json:"description" binding:"max=256"`
		Required    string `json:"required" binding:"max=8"`
		IndexName   string `json:"index_name" binding:"max=128"`
		IsFilter    string `json:"is_filter" binding:"max=8"`
		IsShow      string `json:"is_show" binding:"max=8"`
		Sort        uint   `json:"sort"`
		GroupName   string `json:"group_name" binding:"max=64"`
		Options     string `json:"options" binding:"max=2000"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 从 FieldType 推断 Type 和 Length
	var fieldType string
	var length int
	if req.FieldType != "" {
		preset, ok := constants.GetFieldPreset(req.FieldType)
		if ok && preset != nil {
			fieldType = preset.DataType
			length = preset.Length
		} else {
			// 如果没有预设,使用默认值
			fieldType = "Text"
			length = 255
		}
	} else {
		fieldType = "Text"
		length = 255
	}

	tableField := model.TableField{
		Code:        req.Code,
		TableCode:   req.TableCode,
		Name:        req.Name,
		FieldType:   req.FieldType,
		Type:        fieldType,
		Length:      length,
		IsUnique:    req.IsUnique,
		Description: req.Description,
		Required:    req.Required,
		IndexName:   req.IndexName,
		IsFilter:    req.IsFilter,
		IsShow:      req.IsShow,
		Sort:        req.Sort,
		GroupName:   req.GroupName,
		Status:      "Normal", // 默认状态
	}

	// 解析 Options JSON
	if req.Options != "" {
		var options model.FieldOptions
		if err := json.Unmarshal([]byte(req.Options), &options); err != nil {
			resp.HandleError(c, http.StatusBadRequest, "Invalid options JSON: "+err.Error(), nil)
			return
		}
		tableField.Options = &options
	}

	err := h.tableFieldService.Create(c, &tableField)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, tableField)
}

// Update 更新表字段
// @Summary 更新表字段
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param id path int true "字段ID"
// @Param data body object true "更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_fields/{id} [put]
func (h *tableFieldHandler) Update(c *gin.Context) {
	var req struct {
		Id          uint   `binding:"required"`
		Code        string `binding:"required,max=64"`
		TableCode   string `json:"table_code" binding:"required,max=64"` // 前端发送 table_code
		Name        string `binding:"required,max=128"`
		FieldType   string `json:"field_type" binding:"required,max=32"` // 必填,用于推断 Type/Length
		IsUnique    string `binding:"max=8"`
		Description string `binding:"max=256"`
		Required    string `binding:"max=8"`
		IndexName   string `binding:"max=128"`
		IsFilter    string `binding:"max=8"`
		IsShow      string `binding:"max=8"`
		Sort        uint   ``
		GroupName   string `json:"group_name" binding:"max=64"`
		Options     string `json:"options" binding:"max=2000"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 从 FieldType 推断 Type 和 Length
	var fieldType string
	var length int
	if req.FieldType != "" {
		preset, ok := constants.GetFieldPreset(req.FieldType)
		if ok && preset != nil {
			fieldType = preset.DataType
			length = preset.Length
		} else {
			// 如果没有预设,使用默认值
			fieldType = "Text"
			length = 255
		}
	} else {
		fieldType = "Text"
		length = 255
	}

	tableField := model.TableField{
		ID:          req.Id,
		Code:        req.Code,
		TableCode:   req.TableCode,
		Name:        req.Name,
		Type:        fieldType,
		Length:      length,
		IsUnique:    req.IsUnique,
		Description: req.Description,
		Required:    req.Required,
		IndexName:   req.IndexName,
		IsFilter:    req.IsFilter,
		IsShow:      req.IsShow,
		Sort:        req.Sort,
		Status:      "Normal", // 默认状态
		GroupName:   req.GroupName,
		FieldType:   req.FieldType,
	}

	// 解析 Options JSON
	if req.Options != "" {
		var options model.FieldOptions
		if err := json.Unmarshal([]byte(req.Options), &options); err != nil {
			resp.HandleError(c, http.StatusBadRequest, "Invalid options JSON: "+err.Error(), nil)
			return
		}
		tableField.Options = &options
	}

	err := h.tableFieldService.Update(c, &tableField)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, nil)
}

func (h *tableFieldHandler) Delete(c *gin.Context) {
}

func (h *tableFieldHandler) BatchCreate(c *gin.Context) {
}

// BatchUpdate 批量更新表字段状态
// @Summary 批量更新表字段状态
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param data body object true "批量更新请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_fields/batch [put]
func (h *tableFieldHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		Ids    []uint `form:"ids" binding:"required"`
		Status string `form:"status" binding:"required,oneof=Normal Frozen"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var tableField model.TableField
	tableField.Status = req.Status

	err := h.tableFieldService.BatchUpdate(c, req.Ids, &tableField)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// BatchDelete 批量删除表字段
// @Summary 批量删除表字段
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param data body object true "批量删除请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_fields/batch_delete [post]
func (h *tableFieldHandler) BatchDelete(c *gin.Context) {
	var req struct {
		Ids []uint `form:"ids" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	err := h.tableFieldService.BatchDelete(c, req.Ids)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

// Public 发布表结构
// @Summary 发布表结构
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param data body object true "发布请求"
// @Success 200 {object} map[string]interface{}
// @Router /admin/table_fields/public [post]
func (h *tableFieldHandler) Public(c *gin.Context) {
	var req struct {
		TableCode string `form:"table_code" binding:"required,max=64" json:"table_code"`
	}

	if err := c.ShouldBind(&req); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.tableFieldService.Public(c, req.TableCode)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 返回成功信息,说明三个表都已更新
	resp.HandleSuccess(c, gin.H{
		"message": "表结构发布成功",
		"tables": []string{
			"t_" + req.TableCode,
			"t_" + req.TableCode + "_draft",
			"t_" + req.TableCode + "_log",
		},
	})
}

// GetFieldTypePresets 获取所有字段类型预设
// GET /api/admin/field-type-presets
// GetFieldTypePresets 获取所有字段类型预设
// @Summary 获取字段类型预设列表
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Success 200 {array} object
// @Router /admin/field-type-presets [get]
func (h *tableFieldHandler) GetFieldTypePresets(c *gin.Context) {
	presets := constants.GetAllFieldTypePresets()
	resp.HandleSuccess(c, presets)
}

// GetFieldTypeGroups 获取字段类型分组
// GET /api/admin/field-type-groups
// GetFieldTypeGroups 获取字段类型分组
// @Summary 获取字段类型分组
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Success 200 {array} object
// @Router /admin/field-type-groups [get]
func (h *tableFieldHandler) GetFieldTypeGroups(c *gin.Context) {
	groups := constants.GetFieldTypeGroups()
	resp.HandleSuccess(c, groups)
}

// GetTableOptions 获取表的选项列表（用于关联字段下拉）
// GET /api/admin/table/:table_code/options?filter={"field":"value"}
// GetTableOptions 获取表的选项列表（用于关联字段下拉）
// @Summary 获取表的选项列表
// @Tags 表字段管理
// @Accept json
// @Produce json
// @Param table_code path string true "表编码"
// @Param filter query string false "过滤条件(JSON格式)"
// @Success 200 {array} object
// @Router /admin/table/{table_code}/options [get]
func (h *tableFieldHandler) GetTableOptions(c *gin.Context) {
	tableCode := c.Param("table_code")
	if tableCode == "" {
		resp.HandleError(c, http.StatusBadRequest, "table_code参数不能为空", nil)
		return
	}

	// 解析 filter 参数 (JSON 格式)
	var filter map[string]any
	filterStr := c.Query("filter")
	if filterStr != "" {
		if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
			resp.HandleError(c, http.StatusBadRequest, "Invalid filter JSON: "+err.Error(), nil)
			return
		}
	}

	options, err := h.tableFieldService.GetTableOptions(tableCode, filter)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, options)
}
