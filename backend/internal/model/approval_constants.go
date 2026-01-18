package model

// 审批定义状态常量
const (
	ApprovalDefStatusNormal  = "Normal"  // 启用
	ApprovalDefStatusFrozen  = "Frozen"  // 停用
	ApprovalDefStatusDeleted = "Deleted" // 删除
)

// 审批实例状态常量
const (
	ApprovalStatusPending  = "Pending"  // 审批中
	ApprovalStatusApproved = "Approved" // 已通过
	ApprovalStatusRejected = "Rejected" // 已拒绝
	ApprovalStatusCanceled = "Canceled" // 已撤回
	ApprovalStatusDeleted  = "Deleted"  // 已删除
	ApprovalStatusExpired  = "Expired"  // 已过期
)

// 审批任务状态常量
const (
	TaskStatusPending     = "Pending"     // 待审批
	TaskStatusApproved    = "Approved"    // 已同意
	TaskStatusRejected    = "Rejected"    // 已拒绝
	TaskStatusTransferred = "Transferred" // 已转交
	TaskStatusDone        = "Done"        // 已完成
	TaskStatusCanceled    = "Canceled"    // 已取消
	TaskStatusExpired     = "Expired"     // 已过期
	TaskStatusTimeout     = "Timeout"     // 超时
)

// 操作类型常量
const (
	OperationSubmit   = "Submit"   // 提交申请
	OperationApprove  = "Approve"  // 审批同意
	OperationReject   = "Reject"   // 审批拒绝
	OperationTransfer = "Transfer" // 转交任务
	OperationCancel   = "Cancel"   // 撤回申请
	OperationRemind   = "Remind"   // 催办提醒
	OperationComment  = "Comment"  // 添加评论
	OperationView     = "View"     // 查看详情
)

// 业务操作类型常量 (SAP-ECC习惯)
const (
	OperationCreate        = "Create"        // 新建
	OperationUpdate        = "Update"        // 修改
	OperationFreeze        = "Freeze"        // 冻结
	OperationUnfreeze      = "Unfreeze"      // 解冻
	OperationLock          = "Lock"          // 锁定
	OperationUnlock        = "Unlock"        // 解锁
	OperationDelete        = "Delete"        // 删除
	OperationExtend        = "Extend"        // 扩展
	OperationVoid          = "Void"          // 作废
	OperationCancelBiz     = "Cancel"        // 撤销 (业务撤销，注意与上面的 OperationCancel 区分)
	OperationTerminate     = "Terminate"     // 终止
	OperationBatchCreate   = "BatchCreate"   // 批量创建
	OperationBatchUpdate   = "BatchUpdate"   // 批量更新
	OperationBatchFreeze   = "BatchFreeze"   // 批量冻结
	OperationBatchUnfreeze = "BatchUnfreeze" // 批量解冻
	OperationBatchLock     = "BatchLock"     // 批量锁定
	OperationBatchUnlock   = "BatchUnlock"   // 批量解锁
	OperationBatchDelete   = "BatchDelete"   // 批量删除
	OperationBatchExtend   = "BatchExtend"   // 批量扩展
)

// Action 类型常量 (针对已有系统的操作类型)
const (
	ActionInsert   = "I" // 插入
	ActionUpdate   = "U" // 更新
	ActionFreeze   = "B" // 冻结
	ActionUnfreeze = "C" // 解冻
	ActionDelete   = "D" // 删除
)

// 状态常量
const (
	StatusNormal     = "Normal"     // 正常
	StatusFrozen     = "Frozen"     // 已冻结
	StatusLocked     = "Locked"     // 已锁定
	StatusDeleted    = "Deleted"    // 已删除
	StatusVoided     = "Voided"     // 已作废
	StatusTerminated = "Terminated" // 已终止
	StatusExtended   = "Extended"   // 已扩展
)

// 历史遗留的 operation 代码常量 (为了向后兼容)
const (
	// 注意：这些是历史遗留代码，新代码应该使用上面的 OperationXXX 常量
	LegacyOperationCodeCreate    = "C"  // 历史遗留的 Create 代码
	LegacyOperationCodeUpdate    = "U"  // 历史遗留的 Update 代码
	LegacyOperationCodeFreeze    = "F"  // 历史遗留的 Freeze 代码
	LegacyOperationCodeUnfreeze  = "UF" // 历史遗留的 Unfreeze 代码
	LegacyOperationCodeLock      = "L"  // 历史遗留的 Lock 代码
	LegacyOperationCodeUnlock    = "UL" // 历史遗留的 Unlock 代码
	LegacyOperationCodeDelete    = "D"  // 历史遗留的 Delete 代码
	LegacyOperationCodeTerminate = "T"  // 历史遗留的 Terminate 代码
	// 带 "M" 前缀的表示Batch operations
	LegacyOperationCodeMCreate    = "MC"  // 历史遗留的批量 Create 代码
	LegacyOperationCodeMUpdate    = "MU"  // 历史遗留的批量 Update 代码
	LegacyOperationCodeMFreeze    = "MF"  // 历史遗留的批量 Freeze 代码
	LegacyOperationCodeMUnfreeze  = "MUF" // 历史遗留的批量 Unfreeze 代码
	LegacyOperationCodeMLock      = "ML"  // 历史遗留的批量 Lock 代码
	LegacyOperationCodeMUnlock    = "MUL" // 历史遗留的批量 Unlock 代码
	LegacyOperationCodeMDelete    = "MD"  // 历史遗留的批量 Delete 代码
	LegacyOperationCodeMTerminate = "MT"  // 历史遗留的批量 Terminate 代码
)

// 节点类型常量
const (
	NodeTypeStart       = "START"        // 开始节点
	NodeTypeApproval    = "APPROVAL"     // 审批节点
	NodeTypeCondition   = "CONDITION"    // 条件节点
	NodeTypeCC          = "CC"           // 抄送节点
	NodeTypeEnd         = "END"          // 结束节点
	NodeTypeAutoReject  = "AUTO_REJECT"  // 自动拒绝
	NodeTypeAutoApprove = "AUTO_APPROVE" // 自动通过
	// NodeTypeParallel  = "PARALLEL"  // 并行节点
	// NodeTypeMerge     = "MERGE"     // 合并节点
)

// 审批人配置类型常量
const (
	ApproverTypeUsers       = "USERS"        // 指定用户
	ApproverTypeRoles       = "ROLES"        // 指定角色
	ApproverTypeDepartments = "DEPARTMENTS"  // 指定部门
	ApproverTypePositions   = "POSITIONS"    // 指定岗位
	ApproverTypeSelfSelect  = "SELF_SELECT"  // 发起人自选
	ApproverTypeExpression  = "EXPRESSION"   // 表达式配置
	ApproverTypeSuperior    = "SUPERIOR"     // 直属上级
	ApproverTypeDeptManager = "DEPT_MANAGER" // 部门负责人
	ApproverTypeAutoReject  = "AUTO_REJECT"  // 自动拒绝
	ApproverTypeAutoApprove = "AUTO_APPROVE" // 自动通过
	ApproverTypeSystem      = "SYSTEM"       // 系统
)

// 审批模式常量
const (
	ApprovalModeOR         = "OR"         // 或签（任意一人同意即可）
	ApprovalModeAND        = "AND"        // 会签（所有人都需要同意）
	ApprovalModeSequential = "SEQUENTIAL" // 依次审批（按顺序逐个审批）
)

// 状态转换验证函数
func IsValidApprovalDefStatus(status string) bool {
	switch status {
	case ApprovalDefStatusNormal, ApprovalDefStatusFrozen, ApprovalDefStatusDeleted:
		return true
	default:
		return false
	}
}

func IsValidApprovalStatus(status string) bool {
	switch status {
	case ApprovalStatusPending, ApprovalStatusApproved, ApprovalStatusRejected,
		ApprovalStatusCanceled, ApprovalStatusDeleted, ApprovalStatusExpired:
		return true
	default:
		return false
	}
}

func IsValidTaskStatus(status string) bool {
	switch status {
	case TaskStatusPending, TaskStatusApproved, TaskStatusRejected,
		TaskStatusTransferred, TaskStatusDone, TaskStatusCanceled, TaskStatusExpired:
		return true
	default:
		return false
	}
}

func IsValidNodeType(nodeType string) bool {
	switch nodeType {
	case NodeTypeStart, NodeTypeApproval, NodeTypeCondition,
		// NodeTypeCC, NodeTypeEnd, NodeTypeParallel, NodeTypeMerge:
		NodeTypeCC, NodeTypeEnd, NodeTypeAutoReject, NodeTypeAutoApprove:
		return true
	default:
		return false
	}
}

func IsValidApproverType(approverType string) bool {
	switch approverType {
	case ApproverTypeUsers, ApproverTypeRoles, ApproverTypeDepartments,
		ApproverTypePositions, ApproverTypeSelfSelect, ApproverTypeExpression,
		ApproverTypeSuperior, ApproverTypeDeptManager, ApproverTypeAutoReject,
		ApproverTypeAutoApprove:
		return true
	default:
		return false
	}
}

func IsValidApprovalMode(mode string) bool {
	switch mode {
	case ApprovalModeOR, ApprovalModeAND, ApprovalModeSequential:
		return true
	default:
		return false
	}
}

// 验证函数
func IsValidOperation(operation string) bool {
	switch operation {
	case OperationCreate, OperationUpdate, OperationFreeze, OperationUnfreeze,
		OperationLock, OperationUnlock, OperationDelete, OperationExtend,
		OperationVoid, OperationCancelBiz, OperationTerminate,
		OperationBatchCreate, OperationBatchUpdate, OperationBatchFreeze,
		OperationBatchUnfreeze, OperationBatchLock, OperationBatchUnlock,
		OperationBatchDelete, OperationBatchExtend:
		return true
	default:
		return false
	}
}

func IsValidAction(action string) bool {
	switch action {
	case ActionInsert, ActionUpdate, ActionFreeze, ActionUnfreeze, ActionDelete:
		return true
	default:
		return false
	}
}

func IsValidStatus(status string) bool {
	switch status {
	case StatusNormal, StatusFrozen, StatusLocked, StatusDeleted,
		StatusVoided, StatusTerminated, StatusExtended:
		return true
	default:
		return false
	}
}

// 映射函数
// GetActionByOperation 根据 operation 获取对应的 action
func GetActionByOperation(operation string) string {
	switch operation {
	case OperationCreate, OperationExtend, OperationBatchCreate, OperationBatchExtend:
		return ActionInsert
	case OperationUpdate, OperationLock, OperationUnlock, OperationVoid,
		OperationCancelBiz, OperationTerminate, OperationBatchUpdate,
		OperationBatchLock, OperationBatchUnlock:
		return ActionUpdate
	case OperationFreeze, OperationBatchFreeze:
		return ActionFreeze
	case OperationUnfreeze, OperationBatchUnfreeze:
		return ActionUnfreeze
	case OperationDelete, OperationBatchDelete:
		return ActionDelete
	default:
		return ActionUpdate // 默认返回 Update
	}
}

// GetStatusByOperation 根据 operation 获取对应的状态
func GetStatusByOperation(operation string) string {
	switch operation {
	case OperationCreate, OperationUpdate, OperationUnfreeze, OperationUnlock,
		OperationCancelBiz, OperationBatchCreate, OperationBatchUpdate,
		OperationBatchUnfreeze, OperationBatchUnlock, OperationBatchExtend:
		return StatusNormal
	case OperationFreeze, OperationBatchFreeze:
		return StatusFrozen
	case OperationLock, OperationBatchLock:
		return StatusLocked
	case OperationDelete, OperationBatchDelete:
		return StatusDeleted
	case OperationVoid:
		return StatusVoided
	case OperationTerminate:
		return StatusTerminated
	case OperationExtend:
		return StatusExtended
	default:
		return StatusNormal // 默认返回 Normal
	}
}

// GetOperationNameByOperation 根据 operation 获取操作名称（中文）
func GetOperationNameByOperation(operation string) string {
	switch operation {
	case OperationCreate:
		return "新建"
	case OperationUpdate:
		return "修改"
	case OperationFreeze:
		return "冻结"
	case OperationUnfreeze:
		return "解冻"
	case OperationLock:
		return "锁定"
	case OperationUnlock:
		return "解锁"
	case OperationDelete:
		return "删除"
	case OperationExtend:
		return "扩展"
	case OperationVoid:
		return "作废"
	case OperationCancelBiz:
		return "撤销"
	case OperationTerminate:
		return "终止"
	case OperationBatchCreate:
		return "批量创建"
	case OperationBatchUpdate:
		return "批量更新"
	case OperationBatchFreeze:
		return "批量冻结"
	case OperationBatchUnfreeze:
		return "批量解冻"
	case OperationBatchLock:
		return "批量锁定"
	case OperationBatchUnlock:
		return "批量解锁"
	case OperationBatchDelete:
		return "批量删除"
	case OperationBatchExtend:
		return "批量扩展"
	default:
		return "未知操作"
	}
}

// IsLegacyOperationCode 判断是否是历史遗留的 operation 代码
func IsLegacyOperationCode(code string) bool {
	switch code {
	case LegacyOperationCodeCreate, LegacyOperationCodeUpdate,
		LegacyOperationCodeFreeze, LegacyOperationCodeUnfreeze,
		LegacyOperationCodeLock, LegacyOperationCodeUnlock,
		LegacyOperationCodeDelete, LegacyOperationCodeTerminate,
		LegacyOperationCodeMCreate, LegacyOperationCodeMUpdate,
		LegacyOperationCodeMFreeze, LegacyOperationCodeMUnfreeze,
		LegacyOperationCodeMLock, LegacyOperationCodeMUnlock,
		LegacyOperationCodeMDelete, LegacyOperationCodeMTerminate:
		return true
	default:
		return false
	}
}

// ConvertLegacyCodeToOperation 将历史遗留的代码转换为标准的 operation
func ConvertLegacyCodeToOperation(code string) string {
	switch code {
	case LegacyOperationCodeCreate, LegacyOperationCodeMCreate:
		return OperationCreate
	case LegacyOperationCodeUpdate, LegacyOperationCodeMUpdate:
		return OperationUpdate
	case LegacyOperationCodeFreeze, LegacyOperationCodeMFreeze:
		return OperationFreeze
	case LegacyOperationCodeUnfreeze:
		return OperationUnfreeze
	case LegacyOperationCodeMUnfreeze:
		return OperationBatchUnfreeze
	case LegacyOperationCodeLock, LegacyOperationCodeMLock:
		return OperationLock
	case LegacyOperationCodeUnlock:
		return OperationUnlock
	case LegacyOperationCodeMUnlock:
		return OperationBatchUnlock
	case LegacyOperationCodeDelete, LegacyOperationCodeMDelete:
		return OperationDelete
	case LegacyOperationCodeTerminate, LegacyOperationCodeMTerminate:
		return OperationTerminate
	default:
		return code // 如果不是历史遗留代码，原样返回
	}
}

// 状态转换规则验证
func CanTransitionApprovalStatus(from, to string) bool {
	switch from {
	case ApprovalStatusPending:
		return to == ApprovalStatusApproved || to == ApprovalStatusRejected || to == ApprovalStatusCanceled
	case ApprovalStatusApproved, ApprovalStatusRejected, ApprovalStatusCanceled:
		return to == ApprovalStatusDeleted
	default:
		return false
	}
}

func CanTransitionTaskStatus(from, to string) bool {
	switch from {
	case TaskStatusPending:
		return to == TaskStatusApproved || to == TaskStatusRejected ||
			to == TaskStatusTransferred || to == TaskStatusCanceled
	case TaskStatusApproved, TaskStatusRejected:
		return to == TaskStatusDone
	case TaskStatusTransferred:
		return to == TaskStatusPending
	default:
		return false
	}
}
