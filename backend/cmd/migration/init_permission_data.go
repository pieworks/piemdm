package main

import (
	"context"
	"fmt"
	"piemdm/internal/model"
	"piemdm/internal/repository"
	"piemdm/pkg/log"

	"gorm.io/gorm"
)

// InitPermissionData 初始化权限数据
type InitPermissionData struct {
	db     *gorm.DB
	logger *log.Logger
	repo   repository.PermissionRepository
}

func NewInitPermissionData(db *gorm.DB, logger *log.Logger) *InitPermissionData {
	return &InitPermissionData{
		db:     db,
		logger: logger,
		repo:   repository.NewPermissionRepository(db, logger),
	}
}

func (i *InitPermissionData) Run() error {
	ctx := context.Background()

	i.logger.Info("开始初始化权限数据...")

	// 定义基础权限数据
	permissions := []model.Permission{
		// 用户管理权限
		{Code: "user", Name: "用户管理", Resource: "user", Action: "", ParentID: 0, Description: "用户管理模块"},
		{Code: "user:list", Name: "查看用户", Resource: "user", Action: "list", ParentID: 0, Description: "查看用户列表"},
		{Code: "user:create", Name: "创建用户", Resource: "user", Action: "create", ParentID: 0, Description: "创建新用户"},
		{Code: "user:update", Name: "更新用户", Resource: "user", Action: "update", ParentID: 0, Description: "更新用户信息"},
		{Code: "user:delete", Name: "删除用户", Resource: "user", Action: "delete", ParentID: 0, Description: "删除用户"},

		// 角色管理权限
		{Code: "role", Name: "角色管理", Resource: "role", Action: "", ParentID: 0, Description: "角色管理模块"},
		{Code: "role:list", Name: "查看角色", Resource: "role", Action: "list", ParentID: 0, Description: "查看角色列表"},
		{Code: "role:create", Name: "创建角色", Resource: "role", Action: "create", ParentID: 0, Description: "创建新角色"},
		{Code: "role:update", Name: "更新角色", Resource: "role", Action: "update", ParentID: 0, Description: "更新角色信息"},
		{Code: "role:delete", Name: "删除角色", Resource: "role", Action: "delete", ParentID: 0, Description: "删除角色"},
		{Code: "role:assign_permission", Name: "分配权限", Resource: "role", Action: "assign_permission", ParentID: 0, Description: "为角色分配权限"},

		// 权限管理权限
		{Code: "permission", Name: "权限管理", Resource: "permission", Action: "", ParentID: 0, Description: "权限管理模块"},
		{Code: "permission:list", Name: "查看权限", Resource: "permission", Action: "list", ParentID: 0, Description: "查看权限列表"},
		{Code: "permission:create", Name: "创建权限", Resource: "permission", Action: "create", ParentID: 0, Description: "创建新权限"},
		{Code: "permission:update", Name: "更新权限", Resource: "permission", Action: "update", ParentID: 0, Description: "更新权限信息"},
		{Code: "permission:delete", Name: "删除权限", Resource: "permission", Action: "delete", ParentID: 0, Description: "删除权限"},

		// 审批管理权限
		{Code: "approval", Name: "审批管理", Resource: "approval", Action: "", ParentID: 0, Description: "审批管理模块"},
		{Code: "approval:list", Name: "查看审批", Resource: "approval", Action: "list", ParentID: 0, Description: "查看审批列表"},
		{Code: "approval:create", Name: "创建审批", Resource: "approval", Action: "create", ParentID: 0, Description: "创建审批申请"},
		{Code: "approval:approve", Name: "审批通过", Resource: "approval", Action: "approve", ParentID: 0, Description: "审批通过"},
		{Code: "approval:reject", Name: "审批拒绝", Resource: "approval", Action: "reject", ParentID: 0, Description: "审批拒绝"},

		// 表管理权限
		{Code: "table", Name: "表管理", Resource: "table", Action: "", ParentID: 0, Description: "表管理模块"},
		{Code: "table:list", Name: "查看表", Resource: "table", Action: "list", ParentID: 0, Description: "查看表列表"},
		{Code: "table:create", Name: "创建表", Resource: "table", Action: "create", ParentID: 0, Description: "创建新表"},
		{Code: "table:update", Name: "更新表", Resource: "table", Action: "update", ParentID: 0, Description: "更新表结构"},
		{Code: "table:delete", Name: "删除表", Resource: "table", Action: "delete", ParentID: 0, Description: "删除表"},

		// 应用管理权限
		{Code: "application", Name: "应用管理", Resource: "application", Action: "", ParentID: 0, Description: "应用管理模块"},
		{Code: "application:list", Name: "查看应用", Resource: "application", Action: "list", ParentID: 0, Description: "查看应用列表"},
		{Code: "application:create", Name: "创建应用", Resource: "application", Action: "create", ParentID: 0, Description: "创建新应用"},
		{Code: "application:update", Name: "更新应用", Resource: "application", Action: "update", ParentID: 0, Description: "更新应用信息"},
		{Code: "application:delete", Name: "删除应用", Resource: "application", Action: "delete", ParentID: 0, Description: "删除应用"},

		// 定时任务权限
		{Code: "cron", Name: "定时任务", Resource: "cron", Action: "", ParentID: 0, Description: "定时任务模块"},
		{Code: "cron:list", Name: "查看任务", Resource: "cron", Action: "list", ParentID: 0, Description: "查看定时任务列表"},
		{Code: "cron:create", Name: "创建任务", Resource: "cron", Action: "create", ParentID: 0, Description: "创建定时任务"},
		{Code: "cron:update", Name: "更新任务", Resource: "cron", Action: "update", ParentID: 0, Description: "更新定时任务"},
		{Code: "cron:delete", Name: "删除任务", Resource: "cron", Action: "delete", ParentID: 0, Description: "删除定时任务"},

		// 定时任务日志权限
		{Code: "cron_log", Name: "定时任务日志", Resource: "cron_log", Action: "", ParentID: 0, Description: "定时任务日志模块"},
		{Code: "cron_log:list", Name: "查看任务日志", Resource: "cron_log", Action: "list", ParentID: 0, Description: "查看定时任务日志列表"},
		{Code: "cron_log:delete", Name: "删除任务日志", Resource: "cron_log", Action: "delete", ParentID: 0, Description: "删除定时任务日志"},

		// Webhook 权限
		{Code: "webhook", Name: "Webhook", Resource: "webhook", Action: "", ParentID: 0, Description: "Webhook 模块"},
		{Code: "webhook:list", Name: "查看 Webhook", Resource: "webhook", Action: "list", ParentID: 0, Description: "查看 Webhook 列表"},
		{Code: "webhook:create", Name: "创建 Webhook", Resource: "webhook", Action: "create", ParentID: 0, Description: "创建 Webhook"},
		{Code: "webhook:update", Name: "更新 Webhook", Resource: "webhook", Action: "update", ParentID: 0, Description: "更新 Webhook"},
		{Code: "webhook:delete", Name: "删除 Webhook", Resource: "webhook", Action: "delete", ParentID: 0, Description: "删除 Webhook"},

		// Webhook 投递记录权限
		{Code: "webhook_delivery", Name: "Webhook投递", Resource: "webhook_delivery", Action: "", ParentID: 0, Description: "Webhook投递记录模块"},
		{Code: "webhook_delivery:list", Name: "查看投递记录", Resource: "webhook_delivery", Action: "list", ParentID: 0, Description: "查看投递记录列表"},

		// 审批定义权限
		{Code: "approval_def", Name: "审批定义", Resource: "approval_def", Action: "", ParentID: 0, Description: "审批定义模块"},
		{Code: "approval_def:list", Name: "查看审批定义", Resource: "approval_def", Action: "list", ParentID: 0, Description: "查看审批定义列表"},
		{Code: "approval_def:create", Name: "创建审批定义", Resource: "approval_def", Action: "create", ParentID: 0, Description: "创建审批定义"},
		{Code: "approval_def:update", Name: "更新审批定义", Resource: "approval_def", Action: "update", ParentID: 0, Description: "更新审批定义"},
		{Code: "approval_def:delete", Name: "删除审批定义", Resource: "approval_def", Action: "delete", ParentID: 0, Description: "删除审批定义"},

		// 审批节点权限
		{Code: "approval_node", Name: "审批节点", Resource: "approval_node", Action: "", ParentID: 0, Description: "审批节点模块"},
		{Code: "approval_node:list", Name: "查看审批节点", Resource: "approval_node", Action: "list", ParentID: 0, Description: "查看审批节点列表"},
		{Code: "approval_node:create", Name: "创建审批节点", Resource: "approval_node", Action: "create", ParentID: 0, Description: "创建审批节点"},
		{Code: "approval_node:update", Name: "更新审批节点", Resource: "approval_node", Action: "update", ParentID: 0, Description: "更新审批节点"},
		{Code: "approval_node:delete", Name: "删除审批节点", Resource: "approval_node", Action: "delete", ParentID: 0, Description: "删除审批节点"},

		// 审批任务权限
		{Code: "approval_task", Name: "审批任务", Resource: "approval_task", Action: "", ParentID: 0, Description: "审批任务模块"},
		{Code: "approval_task:list", Name: "查看审批任务", Resource: "approval_task", Action: "list", ParentID: 0, Description: "查看审批任务列表"},
		{Code: "approval_task:create", Name: "创建审批任务", Resource: "approval_task", Action: "create", ParentID: 0, Description: "创建审批任务"},
		{Code: "approval_task:update", Name: "更新审批任务", Resource: "approval_task", Action: "update", ParentID: 0, Description: "更新审批任务"},
		{Code: "approval_task:delete", Name: "删除审批任务", Resource: "approval_task", Action: "delete", ParentID: 0, Description: "删除审批任务"},

		// 通知权限
		{Code: "notification", Name: "通知管理", Resource: "notification", Action: "", ParentID: 0, Description: "通知管理模块"},
		{Code: "notification:list", Name: "查看通知", Resource: "notification", Action: "list", ParentID: 0, Description: "查看通知列表"},
		{Code: "notification:create", Name: "发送通知", Resource: "notification", Action: "create", ParentID: 0, Description: "发送通知"},
		{Code: "notification:delete", Name: "删除通知", Resource: "notification", Action: "delete", ParentID: 0, Description: "删除通知"},

		// 通知模板权限
		{Code: "notification_template", Name: "通知模板", Resource: "notification_template", Action: "", ParentID: 0, Description: "通知模板模块"},
		{Code: "notification_template:list", Name: "查看模板", Resource: "notification_template", Action: "list", ParentID: 0, Description: "查看通知模板列表"},
		{Code: "notification_template:create", Name: "创建模板", Resource: "notification_template", Action: "create", ParentID: 0, Description: "创建通知模板"},
		{Code: "notification_template:update", Name: "更新模板", Resource: "notification_template", Action: "update", ParentID: 0, Description: "更新通知模板"},
		{Code: "notification_template:delete", Name: "删除模板", Resource: "notification_template", Action: "delete", ParentID: 0, Description: "删除通知模板"},

		// 通知日志权限
		{Code: "notification_log", Name: "通知日志", Resource: "notification_log", Action: "", ParentID: 0, Description: "通知日志模块"},
		{Code: "notification_log:list", Name: "查看日志", Resource: "notification_log", Action: "list", ParentID: 0, Description: "查看通知日志列表"},
		{Code: "notification_log:delete", Name: "删除日志", Resource: "notification_log", Action: "delete", ParentID: 0, Description: "删除通知日志"},

		// 规则权限
		{Code: "rule", Name: "规则管理", Resource: "rule", Action: "", ParentID: 0, Description: "规则管理模块"},
		{Code: "rule:list", Name: "查看规则", Resource: "rule", Action: "list", ParentID: 0, Description: "查看规则列表"},
		{Code: "rule:create", Name: "创建规则", Resource: "rule", Action: "create", ParentID: 0, Description: "创建规则"},
		{Code: "rule:update", Name: "更新规则", Resource: "rule", Action: "update", ParentID: 0, Description: "更新规则"},
		{Code: "rule:delete", Name: "删除规则", Resource: "rule", Action: "delete", ParentID: 0, Description: "删除规则"},

		// 表字段权限
		{Code: "table_field", Name: "表字段", Resource: "table_field", Action: "", ParentID: 0, Description: "表字段模块"},
		{Code: "table_field:list", Name: "查看字段", Resource: "table_field", Action: "list", ParentID: 0, Description: "查看表字段列表"},
		{Code: "table_field:create", Name: "创建字段", Resource: "table_field", Action: "create", ParentID: 0, Description: "创建表字段"},
		{Code: "table_field:update", Name: "更新字段", Resource: "table_field", Action: "update", ParentID: 0, Description: "更新表字段"},
		{Code: "table_field:delete", Name: "删除字段", Resource: "table_field", Action: "delete", ParentID: 0, Description: "删除表字段"},

		// 表扩展权限
		{Code: "table_ext", Name: "表扩展", Resource: "table_ext", Action: "", ParentID: 0, Description: "表扩展模块"},
		{Code: "table_ext:list", Name: "查看扩展", Resource: "table_ext", Action: "list", ParentID: 0, Description: "查看表扩展列表"},
		{Code: "table_ext:create", Name: "创建扩展", Resource: "table_ext", Action: "create", ParentID: 0, Description: "创建表扩展"},
		{Code: "table_ext:update", Name: "更新扩展", Resource: "table_ext", Action: "update", ParentID: 0, Description: "更新表扩展"},
		{Code: "table_ext:delete", Name: "删除扩展", Resource: "table_ext", Action: "delete", ParentID: 0, Description: "删除表扩展"},

		// 表审批定义权限
		{Code: "table_approval_def", Name: "表审批关联", Resource: "table_approval_def", Action: "", ParentID: 0, Description: "表审批关联模块"},
		{Code: "table_approval_def:list", Name: "查看关联", Resource: "table_approval_def", Action: "list", ParentID: 0, Description: "查看表审批关联"},
		{Code: "table_approval_def:create", Name: "创建关联", Resource: "table_approval_def", Action: "create", ParentID: 0, Description: "创建表审批关联"},
		{Code: "table_approval_def:update", Name: "更新关联", Resource: "table_approval_def", Action: "update", ParentID: 0, Description: "更新表审批关联"},
		{Code: "table_approval_def:delete", Name: "删除关联", Resource: "table_approval_def", Action: "delete", ParentID: 0, Description: "删除表审批关联"},

		// 数据权限
		{Code: "table_permission", Name: "数据权限", Resource: "table_permission", Action: "", ParentID: 0, Description: "数据权限模块"},
		{Code: "table_permission:list", Name: "查看数据权限", Resource: "table_permission", Action: "list", ParentID: 0, Description: "查看数据权限列表"},
		{Code: "table_permission:create", Name: "分配数据权限", Resource: "table_permission", Action: "create", ParentID: 0, Description: "分配数据权限"},
		{Code: "table_permission:update", Name: "更新数据权限", Resource: "table_permission", Action: "update", ParentID: 0, Description: "更新数据权限"},
		{Code: "table_permission:delete", Name: "撤销数据权限", Resource: "table_permission", Action: "delete", ParentID: 0, Description: "撤销数据权限"},
	}
	// 创建权限
	createdCount := 0
	skippedCount := 0

	for _, perm := range permissions {
		// 检查是否已存在
		existing, err := i.repo.FindByCode(ctx, perm.Code)
		if err == nil && existing != nil {
			i.logger.Info(fmt.Sprintf("权限 %s 已存在,跳过", perm.Code))
			skippedCount++
			continue
		}

		// 创建权限
		if err := i.repo.Create(ctx, &perm); err != nil {
			i.logger.Error(fmt.Sprintf("创建权限 %s 失败: %s", perm.Code, err.Error()))
			return err
		}

		i.logger.Info(fmt.Sprintf("创建权限: %s - %s", perm.Code, perm.Name))
		createdCount++
	}

	i.logger.Info(fmt.Sprintf("权限数据初始化完成! 创建: %d, 跳过: %d", createdCount, skippedCount))

	// 构建树状结构 (将子权限关联到父权限)
	if err := i.buildPermissionTree(ctx); err != nil {
		i.logger.Error("构建权限树失败: " + err.Error())
		return err
	}

	return nil
}

func (i *InitPermissionData) buildPermissionTree(ctx context.Context) error {
	i.logger.Info("开始构建权限树...")

	// 定义父子关系
	parentChildMap := map[string][]string{
		"user":                  {"user:list", "user:create", "user:update", "user:delete"},
		"role":                  {"role:list", "role:create", "role:update", "role:delete", "role:assign_permission"},
		"permission":            {"permission:list", "permission:create", "permission:update", "permission:delete"},
		"approval":              {"approval:list", "approval:create", "approval:approve", "approval:reject"},
		"table":                 {"table:list", "table:create", "table:update", "table:delete"},
		"application":           {"application:list", "application:create", "application:update", "application:delete"},
		"cron":                  {"cron:list", "cron:create", "cron:update", "cron:delete"},
		"cron_log":              {"cron_log:list", "cron_log:delete"},
		"webhook":               {"webhook:list", "webhook:create", "webhook:update", "webhook:delete"},
		"webhook_delivery":      {"webhook_delivery:list"},
		"approval_def":          {"approval_def:list", "approval_def:create", "approval_def:update", "approval_def:delete"},
		"approval_node":         {"approval_node:list", "approval_node:create", "approval_node:update", "approval_node:delete"},
		"approval_task":         {"approval_task:list", "approval_task:create", "approval_task:update", "approval_task:delete"},
		"notification":          {"notification:list", "notification:create", "notification:delete"},
		"notification_template": {"notification_template:list", "notification_template:create", "notification_template:update", "notification_template:delete"},
		"notification_log":      {"notification_log:list", "notification_log:delete"},
		"rule":                  {"rule:list", "rule:create", "rule:update", "rule:delete"},
		"table_field":           {"table_field:list", "table_field:create", "table_field:update", "table_field:delete"},
		"table_ext":             {"table_ext:list", "table_ext:create", "table_ext:update", "table_ext:delete"},
		"table_approval_def":    {"table_approval_def:list", "table_approval_def:create", "table_approval_def:update", "table_approval_def:delete"},
		"table_permission":      {"table_permission:list", "table_permission:create", "table_permission:update", "table_permission:delete"},
	}

	for parentCode, childCodes := range parentChildMap {
		// 获取父权限
		parent, err := i.repo.FindByCode(ctx, parentCode)
		if err != nil {
			i.logger.Warn(fmt.Sprintf("找不到父权限: %s", parentCode))
			continue
		}

		// 更新子权限的 parent_id
		for _, childCode := range childCodes {
			child, err := i.repo.FindByCode(ctx, childCode)
			if err != nil {
				i.logger.Warn(fmt.Sprintf("找不到子权限: %s", childCode))
				continue
			}

			child.ParentID = parent.ID
			if err := i.repo.Update(ctx, child); err != nil {
				i.logger.Error(fmt.Sprintf("更新权限 %s 的父级失败: %s", childCode, err.Error()))
				return err
			}
		}

		i.logger.Info(fmt.Sprintf("已设置 %s 的 %d 个子权限", parentCode, len(childCodes)))
	}

	i.logger.Info("权限树构建完成!")
	return nil
}
