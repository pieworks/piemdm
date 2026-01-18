package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"piemdm/internal/service"
	"piemdm/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
)

type EntityHandler interface {
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

	// 历史与日志
	ListEntityHistories(c *gin.Context) // get entity from *draft table
	ListEntityLogs(c *gin.Context)      // get entity from *log table

	// 辅助功能
	Import(ctx *gin.Context)
	Export(c *gin.Context)
	Template(c *gin.Context)
	// 其他
	GetStatistics(c *gin.Context)
}

type entityHandler struct {
	*Handler
	entityService service.EntityService

	tableFieldService              service.TableFieldService
	tableApprovalDefinitionService service.TableApprovalDefinitionService
	tablePermissionService         service.TablePermissionService
	conf                           *viper.Viper
}

func NewEntityHandler(handler *Handler, entityService service.EntityService, tableFieldService service.TableFieldService, tableApprovalDefinitionService service.TableApprovalDefinitionService, tablePermissionService service.TablePermissionService, conf *viper.Viper) EntityHandler {
	return &entityHandler{
		Handler:       handler,
		entityService: entityService,

		tableFieldService:              tableFieldService,
		tableApprovalDefinitionService: tableApprovalDefinitionService,
		tablePermissionService:         tablePermissionService,
		conf:                           conf,
	}
}

func (h *entityHandler) List(c *gin.Context) {
	// page=1&pageSize=15&...
	// get from  table field from api /table_field/find
	page, pageSize := GetPage(c)
	where := map[string]any{}
	var total int64

	entityID := c.Query("entity_id")
	tableCode := c.Query("table_code")
	dictionaryClass := c.Query("dict_code")
	isDraft := c.Query("is_draft")
	approvalCode := c.Query("approval_code")

	if dictionaryClass != "" {
		where["dict_code"] = dictionaryClass
	}
	if entityID != "" {
		where["entity_id"] = entityID
	}
	if approvalCode != "" {
		where["approval_code"] = approvalCode
	}
	if isDraft != "" {
		tableCode = fmt.Sprintf("%s_draft", tableCode)
	}

	for k, v := range c.Request.URL.Query() {
		// v[0] not enpty then use condition
		if v[0] != "" {
			// replace fullwidth character
			str := strings.ReplaceAll(v[0], "，", ",")
			where[k] = strings.TrimSpace(str)
		}
	}

	// delete control parameter
	delete(where, "page")
	delete(where, "pageSize")
	delete(where, "table_code")
	delete(where, "is_draft")

	entities, err := h.entityService.List(c, tableCode, page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, entities)
}

func (h *entityHandler) Get(c *gin.Context) {
	var params struct {
		ID        uint   `uri:"id" binding:"required"`
		TableCode string `uri:"table_code" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// tableCode := c.Query("table_code")
	tableCode := params.TableCode
	fieldWhere := map[string]any{}
	fieldWhere["table_code"] = tableCode
	fieldWhere["status"] = "Normal"
	sel := "table_code,code,name,field_type,type,length,is_unique,required,is_show,options,group_name"
	tableFields, err2 := h.tableFieldService.Find(sel, fieldWhere)
	if err2 != nil {
		resp.HandleError(c, http.StatusInternalServerError, err2.Error(), nil)
		return
	}

	entity, err := h.entityService.Get(c, tableCode, params.ID)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, gin.H{
		"info":        entity,
		"tableFields": tableFields,
	})
}

func (h *entityHandler) Create(c *gin.Context) {
	// 因为是动态模型，所以不能使用固定的request结构接收传入参数
	tableCode := c.Param("table_code")
	req := map[string]any{}
	binding.EnableDecoderUseNumber = true
	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if tableCode == "" || req["table_code"] == nil || req["table_code"] == "" {
		resp.HandleError(c, http.StatusBadRequest, "tableCode is required", nil)
		return
	}
	if req["reason"] == nil || req["reason"] == "" {
		resp.HandleError(c, http.StatusBadRequest, "reason is required", nil)
		return
	}

	reason := req["reason"].(string)
	delete(req, "reason")

	h.logger.Debug("handler-entity-CreateEntity", "req", req)

	// 判断是否关联审批流程，如果有审批流程则提交审批, 没有审批流程, 则直接保存数据
	operation := "Create"
	tableApprovalDefs, err := h.tableApprovalDefinitionService.List(tableCode, operation)
	if err != nil {
		h.logger.Error("err", "err", err)
	}

	// 如果有审批流程则提交审批, 没有审批流程, 则直接保存数据
	if len(tableApprovalDefs) > 0 {
		if err := h.entityService.CreateDraft(c, tableCode, reason, req); err != nil {
			resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
	} else {
		if err := h.entityService.Create(c, tableCode, req); err != nil {
			resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}

	resp.HandleSuccess(c, nil)
}

func (h *entityHandler) Update(c *gin.Context) {
	// 因为是动态模型，所以不能使用固定的request结构接收传入参数
	tableCode := c.Param("table_code")
	// id := c.Param("id")
	formMap := map[string]any{}
	binding.EnableDecoderUseNumber = true
	err := c.ShouldBind(&formMap)
	if err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	h.logger.Debug("UpdateEntity handler", "formMap", formMap)

	// 获取FormData中的值
	// tableCode := formMap["table_code"].(string)

	// 验证reason字段
	reasonVal, ok := formMap["reason"]
	if !ok || reasonVal == nil {
		resp.HandleError(c, http.StatusBadRequest, "reason is required", nil)
		return
	}
	reason, ok := reasonVal.(string)
	if !ok || reason == "" {
		resp.HandleError(c, http.StatusBadRequest, "reason must be a non-empty string", nil)
		return
	}

	idInt64, err := formMap["id"].(json.Number).Int64()
	if err != nil {
		panic(err)
	}

	formMap["id"] = uint(idInt64)
	formMap["table_code"] = tableCode
	// delete(formMap, "table_code")
	delete(formMap, "reason")

	// 序列化数组类型的字段为 JSON 字符串
	for key, value := range formMap {
		if arr, ok := value.([]interface{}); ok {
			jsonBytes, err := json.Marshal(arr)
			if err != nil {
				h.logger.Error("Failed to serialize array field", "key", key, "error", err)
				continue
			}
			formMap[key] = string(jsonBytes)
		}
	}

	if tableCode == "" {
		resp.HandleError(c, http.StatusBadRequest, "Query tableCode is empty。", nil)
		return
	}
	if formMap["id"] == "" {
		resp.HandleError(c, http.StatusBadRequest, "Query id is empty。", nil)
		return
	}
	h.logger.Debug("UpdateDraft", "tableCode", tableCode, "reason", reason, "formMap", formMap)

	if err := h.entityService.UpdateDraft(c, tableCode, reason, formMap); err != nil {
		h.logger.Debug("error", "error", err.Error())
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, nil)
}

func (h *entityHandler) Delete(c *gin.Context) {
}

func (h *entityHandler) BatchCreate(c *gin.Context) {
}

// TODO 冻结数据
// 根据上传的操作获取状态。
func (h *entityHandler) BatchUpdate(c *gin.Context) {
	var params struct {
		IDs []uint `form:"ids" json:"ids"`
		// 未上传字段
		Status    string `form:"status" json:"status"`
		TableCode string `form:"table_code" json:"table_code"`
		// 未上传字段
		Reason    string `form:"reason" json:"reason"`
		Operation string `form:"operation" json:"operation"`
	}
	if err := c.ShouldBind(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	entityMap := map[string]any{}
	entityMap["status"] = params.Status
	entityMap["operation"] = params.Operation
	h.logger.Info("service BatchUpdateEntities", "entityMap", entityMap)

	// params.Reason = "change status..."
	h.logger.Info("service BatchUpdateEntities", "params", params)

	err := h.entityService.BatchUpdate(c, params.TableCode, params.Reason, params.IDs, entityMap)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
	// resp.HandleSuccess(c, nil)
}

func (h *entityHandler) BatchDelete(c *gin.Context) {
	var params struct {
		TableCode string `form:"table_code" json:"table_code"`
		IDs       []uint `form:"ids" json:"ids"`
		Reason    string `form:"reason" json:"reason"`
	}

	if err := c.ShouldBind(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	tableCode := params.TableCode
	reason := params.Reason
	err := h.entityService.BatchDelete(c, tableCode, reason, params.IDs)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp.HandleSuccess(c, gin.H{
		"message": "success",
	})
}

func (h *entityHandler) ListEntityLogs(c *gin.Context) {
	// page=1&pageSize=15&...
	// get from  table field from api /table_field/find
	page, pageSize := GetPage(c)
	where := map[string]any{}
	var total int64

	entityID := c.Query("entity_id")
	tableCode := c.Query("table_code")
	if entityID != "" {
		where["entity_id"] = entityID
	}
	logTableCode := fmt.Sprintf("%s_log", tableCode)
	h.logger.Debug("handler-entity-ListEntities", "logTableCode", logTableCode)

	for k, v := range c.Request.URL.Query() {
		if v[0] != "" {
			str := strings.Replace(v[0], "，", ",", -1)
			where[k] = strings.TrimSpace(str)
		}
	}

	// delete control parameter
	delete(where, "page")
	delete(where, "pageSize")
	delete(where, "table_code")

	// where["table_code"] = tableCode
	entities, err := h.entityService.List(c, logTableCode, page, pageSize, &total, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, entities)
}

func (h *entityHandler) ListEntityHistories(c *gin.Context) {
	// page=1&pageSize=15&...
	// get from  table field from api /table_field/find
	page, pageSize := GetPage(c)
	where := map[string]any{}
	var total int64

	entityID := c.Query("entity_id")
	code := c.Query("code")
	tableCode := c.Query("table_code")
	if entityID != "" {
		where["entity_id"] = entityID
	}
	if code != "" {
		where["code"] = code
	}
	tableCode = fmt.Sprintf("%s_draft", tableCode)

	// Permission check for table_code
	if tableCode != "" {
		userIDStr := c.GetString("user_id")
		if userIDStr != "" {
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				h.logger.Error("Failed to parse user_id", "user_id", userIDStr, "error", err)
				resp.HandleError(c, http.StatusInternalServerError, "Invalid user ID", nil)
				return
			}

			hasPermission, err := h.tablePermissionService.CheckTablePermission(c, uint(userID), tableCode)
			if err != nil {
				h.logger.Error("Failed to check table permission", "error", err)
				resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
				return
			}

			if !hasPermission {
				h.logger.Warn("User has no permission to access table", "user_id", userID, "table_code", tableCode)
				resp.HandleError(c, http.StatusForbidden, "No permission to access this table", nil)
				return
			}
		}
	}

	// delete control parameter
	delete(where, "page")
	delete(where, "pageSize")
	delete(where, "table_code")

	selectString := "id,code,name,entity_id,updated_at"
	entities, err := h.entityService.Find(c, tableCode, selectString, where)
	if err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	links := resp.GeneratePaginationLinks(c.Request, page, pageSize, int(total))
	c.Header("Link", links.String())

	resp.HandleSuccess(c, entities)
}

func (h *entityHandler) Import(c *gin.Context) {
	file, err := c.FormFile("file") // 获取上传的文件
	if err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	reader, err := file.Open() // 打开文件，获取 io.Reader

	// file, err := c.Request.MultipartForm.File["file"][0].Open()
	tableCode := c.PostForm("table_code")
	operation := c.PostForm("operation")
	reason := c.PostForm("reason")
	if err != nil || tableCode == "" || operation == "" {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.Debug("h Import: ",
		"tableCode", tableCode,
		"reason", reason,
		"operation", operation)
	if err = h.entityService.Import(c, tableCode, reason, operation, reader); err != nil {
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, nil)
}

func (h *entityHandler) Export(c *gin.Context) {
	// 可以使用 c.QueryMap("ids") 接收query array
	// var params struct {
	// 	TableCode  string                 `form:"table_code" json:"table_code"`
	// 	Filter     string                 `form:"filter" json:"filter"`
	// 	Ids        string                `form:"ids" json:"ids"`
	// 	SearchData map[string]any `form:"searchData" json:"searchData"`
	// }
	// if err := c.ShouldBindQuery(&params); err != nil {
	// 	resp.HandleError(c, http.StatusBadRequest,  err.Error(), nil)
	// 	return
	// }

	// tableCode := params.TableCode
	// filter := params.Filter
	tableCode := c.Query("table_code")
	filter := c.Query("filter")
	ids := c.Query("ids")
	where := make(map[string]any)
	switch filter {
	case "selected":
		where["id in"] = ids
	case "filtered":
		// 根据查询条件下载的时候需要把参数收集起来
		for k, v := range c.Request.URL.Query() {
			if v[0] != "" {
				where[k] = v[0]
			}
		}
		delete(where, "ids")
	default:
		// all
	}

	delete(where, "table_code")
	delete(where, "filter")

	filename, err := h.entityService.Export(c, tableCode, filter, where)
	if err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	path := h.conf.GetString("app.export-save-path")
	fullURL := h.conf.GetString("app.prefix-url") + "/" + path + filename
	resp.HandleSuccess(c, gin.H{
		"name":            filename,
		"export_url":      fullURL,
		"export_save_url": path,
	})
}

func (h *entityHandler) Template(c *gin.Context) {
	// 可以使用 c.QueryMap("ids") 接收query array
	var params struct {
		TableCode string `form:"table_code" json:"table_code"`
		Operation string `form:"operation" json:"operation"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	filename, err := h.entityService.Template(c, params.TableCode, params.Operation)
	if err != nil {
		resp.HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	path := h.conf.GetString("app.export-save-path")
	fullURL := h.conf.GetString("app.prefix-url") + "/" + path + filename
	resp.HandleSuccess(c, gin.H{
		"name":            filename,
		"export_url":      fullURL,
		"export_save_url": path,
	})

	// resp.HandleSuccess(c, nil)
}
func (h *entityHandler) GetStatistics(c *gin.Context) {
	statistics, err := h.entityService.GetEntitiesStatistics(c)
	if err != nil {
		h.logger.Error("获取实体统计失败", "error", err)
		resp.HandleError(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.HandleSuccess(c, statistics)
}
