package service

import (
	"errors"
	"strings"

	"piemdm/internal/model"
	"piemdm/internal/repository"

	"github.com/google/uuid"
)

type ApprovalTaskService interface {
	// Base CRUD
	Get(id uint) (*model.ApprovalTask, error)
	List(page, pageSize int, total *int64, where map[string]any, order string) ([]*model.ApprovalTask, error)
	Create(approvalTask *model.ApprovalTask) error
	Update(approvalTask *model.ApprovalTask) error
	Delete(id uint) (*model.ApprovalTask, error)

	// Batch operations
	BatchUpdate(ids []uint, approvalTask *model.ApprovalTask) error
	BatchDelete(ids []uint) error

	// 业务查询方法
	GetByCode(code string) (*model.ApprovalTask, error)
	First(where map[string]any) (*model.ApprovalTask, error)
	GetByApprovalCode(approvalCode string) ([]*model.ApprovalTask, error)
	GetByAssignee(assigneeID string) ([]*model.ApprovalTask, error)
	GetPendingByAssignee(assigneeID string) ([]*model.ApprovalTask, error)
	GetByNodeCode(nodeCode string) ([]*model.ApprovalTask, error)
	GetByStatus(status string) ([]*model.ApprovalTask, error)
	GetOverdueTasks() ([]*model.ApprovalTask, error)
	GetExpiredTasks() ([]*model.ApprovalTask, error)

	// 任务处理
	ApproveTask(taskID uint, assigneeID, comment, reason string) error
	RejectTask(taskID uint, assigneeID, comment, reason string) error
	TransferTask(taskID uint, fromUserID, toUserID, fromUserName, toUserName, reason string) error
	CompleteTask(taskID uint, result, comment, reason string) error

	// 任务状态管理
	StartTask(taskID uint) error
	CancelTask(taskID uint, reason string) error
	UpdateTaskStatus(taskID uint, status string) error

	// 催办管理
	RemindTask(taskID uint) error
	BatchRemindTasks(assigneeID string) error
	GetTasksNeedRemind() ([]*model.ApprovalTask, error)

	// 验证方法
	ValidateTask(task *model.ApprovalTask) error
	CanProcess(taskID uint, assigneeID string) (bool, error)
}

type approvalTaskService struct {
	*Service
	approvalTaskRepository repository.ApprovalTaskRepository
	approvalRepository     repository.ApprovalRepository
}

func NewApprovalTaskService(
	service *Service,
	approvalTaskRepository repository.ApprovalTaskRepository,
	approvalRepository repository.ApprovalRepository,
) ApprovalTaskService {
	return &approvalTaskService{
		Service:                service,
		approvalTaskRepository: approvalTaskRepository,
		approvalRepository:     approvalRepository,
	}
}

// 基础CRUD操作
func (s *approvalTaskService) Get(id uint) (*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindOne(id)
}

func (s *approvalTaskService) GetByCode(code string) (*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FirstByCode(code)
}

func (s *approvalTaskService) First(where map[string]any) (*model.ApprovalTask, error) {
	return s.approvalTaskRepository.First(where)
}

func (s *approvalTaskService) List(page, pageSize int, total *int64, where map[string]any, order string) ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindPage(page, pageSize, total, where, order)
}

func (s *approvalTaskService) Create(approvalTask *model.ApprovalTask) error {
	// 验证任务
	if err := s.ValidateTask(approvalTask); err != nil {
		s.logger.Error("任务验证失败", "error", err)
		return err
	}

	// 生成任务编码
	if approvalTask.TaskCode == "" {
		uuid := uuid.New()
		approvalTask.TaskCode = strings.ToUpper(uuid.String())
	}

	return s.approvalTaskRepository.Create(approvalTask)
}

func (s *approvalTaskService) Update(approvalTask *model.ApprovalTask) error {
	// 验证任务
	if err := s.ValidateTask(approvalTask); err != nil {
		s.logger.Error("任务验证失败", "error", err)
		return err
	}

	return s.approvalTaskRepository.Update(approvalTask)
}

func (s *approvalTaskService) BatchUpdate(ids []uint, approvalTask *model.ApprovalTask) error {
	return s.approvalTaskRepository.BatchUpdate(ids, approvalTask)
}

func (s *approvalTaskService) Delete(id uint) (*model.ApprovalTask, error) {
	// 先获取记录
	task, err := s.approvalTaskRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	// 执行删除
	err = s.approvalTaskRepository.BatchDelete([]uint{id})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *approvalTaskService) BatchDelete(ids []uint) error {
	return s.approvalTaskRepository.BatchDelete(ids)
}

// 业务查询方法
func (s *approvalTaskService) GetByApprovalCode(approvalCode string) ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindByApprovalCode(approvalCode)
}

func (s *approvalTaskService) GetByAssignee(assigneeID string) ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindByAssigneeID(assigneeID)
}

func (s *approvalTaskService) GetPendingByAssignee(assigneeID string) ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindPendingByAssignee(assigneeID)
}

func (s *approvalTaskService) GetByNodeCode(nodeCode string) ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindByNodeCode(nodeCode)
}

func (s *approvalTaskService) GetByStatus(status string) ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindByStatus(status)
}

func (s *approvalTaskService) GetOverdueTasks() ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindOverdueTasks()
}

func (s *approvalTaskService) GetExpiredTasks() ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindExpiredTasks()
}

// 任务处理
func (s *approvalTaskService) ApproveTask(taskID uint, assigneeID, comment, reason string) error {
	// 验证权限
	canProcess, err := s.CanProcess(taskID, assigneeID)
	if err != nil {
		return err
	}
	if !canProcess {
		return errors.New("无权限处理此任务")
	}

	// 获取任务
	task, err := s.approvalTaskRepository.FindOne(uint(taskID))
	if err != nil {
		return err
	}

	if !task.IsPending() {
		return errors.New("任务已处理")
	}

	// 标记任务开始处理
	// task.MarkAsStarted()

	// 完成任务
	if err := s.approvalTaskRepository.CompleteTask(taskID, model.TaskStatusApproved, comment, reason); err != nil {
		return err
	}

	return nil
}

func (s *approvalTaskService) RejectTask(taskID uint, assigneeID, comment, reason string) error {
	// 验证权限
	canProcess, err := s.CanProcess(taskID, assigneeID)
	if err != nil {
		return err
	}
	if !canProcess {
		return errors.New("无权限处理此任务")
	}

	// 获取任务
	task, err := s.approvalTaskRepository.FindOne(uint(taskID))
	if err != nil {
		return err
	}

	if !task.IsPending() {
		return errors.New("任务已处理")
	}

	// if !task.CanReject() {
	// 	return errors.New("此任务不允许拒绝")
	// }

	// 标记任务开始处理
	// task.MarkAsStarted()

	// 完成任务
	if err := s.approvalTaskRepository.CompleteTask(taskID, model.TaskStatusRejected, comment, reason); err != nil {
		return err
	}

	return nil
}

func (s *approvalTaskService) TransferTask(taskID uint, fromUserID, toUserID, fromUserName, toUserName, reason string) error {
	// 验证权限
	// canTransfer, err := s.CanTransfer(taskID)
	// if err != nil {
	// 	return err
	// }
	// if !canTransfer {
	// 	return errors.New("此任务不允许转交")
	// }

	// 获取任务
	task, err := s.approvalTaskRepository.FindOne(uint(taskID))
	if err != nil {
		return err
	}

	if !task.IsPending() {
		return errors.New("只有待处理的任务才能转交")
	}

	if task.AssigneeID != fromUserID {
		return errors.New("只能转交自己的任务")
	}

	// 执行转交
	if err := s.approvalTaskRepository.TransferTask(taskID, fromUserID, toUserID, fromUserName, toUserName, reason); err != nil {
		return err
	}

	return nil
}

func (s *approvalTaskService) CompleteTask(taskID uint, result, comment, reason string) error {
	task, err := s.approvalTaskRepository.FindOne(taskID)
	if err != nil {
		return err
	}

	if task.IsCompleted() {
		return errors.New("任务已完成")
	}

	// 标记任务完成
	// task.MarkAsCompleted()

	return s.approvalTaskRepository.CompleteTask(taskID, result, comment, reason)
}

// 任务状态管理
func (s *approvalTaskService) StartTask(taskID uint) error {
	task, err := s.approvalTaskRepository.FindOne(taskID)
	if err != nil {
		return err
	}

	if !task.IsPending() {
		return errors.New("只有待处理的任务才能开始")
	}

	// 标记任务开始
	// task.MarkAsStarted()
	return s.approvalTaskRepository.Update(task)
}

func (s *approvalTaskService) CancelTask(taskID uint, reason string) error {
	task, err := s.approvalTaskRepository.FindOne(taskID)
	if err != nil {
		return err
	}

	if task.IsCompleted() {
		return errors.New("已完成的任务不能取消")
	}

	return s.approvalTaskRepository.UpdateStatus(taskID, model.TaskStatusCanceled)
}

func (s *approvalTaskService) UpdateTaskStatus(taskID uint, status string) error {
	if !model.IsValidTaskStatus(status) {
		return errors.New("无效的任务状态")
	}
	return s.approvalTaskRepository.UpdateStatus(taskID, status)
}

// 催办管理
func (s *approvalTaskService) RemindTask(taskID uint) error {
	task, err := s.approvalTaskRepository.FindOne(taskID)
	if err != nil {
		return err
	}

	if !task.IsPending() {
		return errors.New("只有待处理的任务才能催办")
	}

	// 更新催办次数
	newCount := task.RemindCount + 1
	if err := s.approvalTaskRepository.UpdateRemindCount(taskID, newCount); err != nil {
		return err
	}

	return nil
}

func (s *approvalTaskService) BatchRemindTasks(assigneeID string) error {
	// 获取用户的待处理任务
	tasks, err := s.approvalTaskRepository.FindPendingByAssignee(assigneeID)
	if err != nil {
		return err
	}

	remindCount := 0
	for _, task := range tasks {
		// 检查是否需要催办（距离上次催办超过24小时）
		// if task.LastRemindAt == nil || time.Since(*task.LastRemindAt) > 24*time.Hour {
		if err := s.RemindTask(task.ID); err != nil {
			s.logger.Error("催办任务失败", "error", err, "taskID", task.ID)
		} else {
			remindCount++
		}
		// }
	}

	return nil
}

func (s *approvalTaskService) GetTasksNeedRemind() ([]*model.ApprovalTask, error) {
	return s.approvalTaskRepository.FindTasksNeedRemind()
}

// 验证方法
func (s *approvalTaskService) ValidateTask(task *model.ApprovalTask) error {
	if task.ApprovalCode == "" {
		return errors.New("审批实例编码不能为空")
	}
	if task.NodeCode == "" {
		return errors.New("节点编码不能为空")
	}
	if task.AssigneeID == "" {
		return errors.New("审批人不能为空")
	}
	if task.Status != "" && !model.IsValidTaskStatus(task.Status) {
		return errors.New("无效的任务状态")
	}
	return nil
}

func (s *approvalTaskService) CanProcess(taskID uint, assigneeID string) (bool, error) {
	task, err := s.approvalTaskRepository.FindOne(taskID)
	if err != nil {
		return false, err
	}

	// 检查是否是任务的审批人
	if task.AssigneeID != assigneeID {
		return false, nil
	}

	// 检查任务状态
	if !task.IsPending() {
		return false, nil
	}

	// 检查审批实例状态
	approval, err := s.approvalRepository.FirstByCode(task.ApprovalCode)
	if err != nil {
		return false, err
	}

	if !approval.IsPending() {
		return false, nil
	}

	return true, nil
}

// func (s *approvalTaskService) CanTransfer(taskID uint) (bool, error) {
// 	task, err := s.approvalTaskRepository.FindOne(taskID)
// 	if err != nil {
// 		return false, err
// 	}
// 	// return task.CanTransfer(), nil
// 	return false, nil
// }

// 统计方法
// func (s *approvalTaskService) GetTaskStatistics(assigneeID string) (map[string]int64, error) {
// 	stats := make(map[string]int64)

// 	// 统计各状态的任务数量
// 	statuses := []string{
// 		model.TaskStatusPending,
// 		model.TaskStatusApproved,
// 		model.TaskStatusRejected,
// 		model.TaskStatusTransferred,
// 		model.TaskStatusDone,
// 		model.TaskStatusCanceled,
// 	}

// 	for _, status := range statuses {
// 		count, err := s.approvalTaskRepository.CountByStatus(status)
// 		if err != nil {
// 			return nil, err
// 		}
// 		stats[status] = count
// 	}

// 	// 统计用户相关的任务
// 	userTaskCount, err := s.approvalTaskRepository.CountByAssignee(assigneeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	stats["user_total"] = userTaskCount

// 	userPendingCount, err := s.approvalTaskRepository.CountPendingByAssignee(assigneeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	stats["user_pending"] = userPendingCount

// 	return stats, nil
// }

// func (s *approvalTaskService) GetProcessDuration(taskID uint) (int, error) {
// 	task, err := s.approvalTaskRepository.FindOne(taskID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	// return task.GetProcessDuration(), nil
// 	return 0, nil
// }
