package service

import (
	"errors"
	"fmt"

	"piemdm/internal/model"
	"piemdm/internal/repository"
)

type ApprovalNodeService interface {
	// Base CRUD
	Get(id uint) (*model.ApprovalNode, error)
	List(page, pageSize int, total *int64, where map[string]any) ([]*model.ApprovalNode, error)
	Create(approvalNode *model.ApprovalNode) error
	Update(approvalNode *model.ApprovalNode) error
	Delete(id uint) (*model.ApprovalNode, error)

	// Batch operations
	BatchUpdate(ids []uint, approvalNode *model.ApprovalNode) error
	BatchDelete(ids []uint) error

	// 业务方法
	GetByCode(code string) (*model.ApprovalNode, error)
	First(where map[string]any) (*model.ApprovalNode, error)
	GetByApprovalDefCode(approvalDefCode string) ([]*model.ApprovalNode, error)
	GetActiveByApprovalDefCode(approvalDefCode string) ([]*model.ApprovalNode, error)
	GetByNodeType(nodeType string) ([]*model.ApprovalNode, error)
	GetStartNode(approvalDefCode string) (*model.ApprovalNode, error)
	GetEndNodes(approvalDefCode string) ([]*model.ApprovalNode, error)
	GetNextNodes(approvalDefCode, currentNodeCode string) ([]*model.ApprovalNode, error)

	// 流程验证
	ValidateWorkflow(approvalDefCode string) error
	ValidateNode(node *model.ApprovalNode) error
	CanDeleteNode(id uint) (bool, error)

	// 节点配置
	ConfigureApprovers(nodeId uint, approverConfig string) error
	ConfigureConditions(nodeId uint, conditionConfig string) error
	ConfigureTimeouts(nodeId uint, timeoutHours int) error

	// 状态管理
	ActivateNode(id uint) error
	DeactivateNode(id uint) error
	BatchSyncNodes(approvalDefCode string, nodes []*model.ApprovalNode) (string, error)
}

type approvalNodeService struct {
	*Service
	approvalNodeRepository       repository.ApprovalNodeRepository
	approvalDefinitionRepository repository.ApprovalDefinitionRepository
}

func NewApprovalNodeService(
	service *Service,
	approvalNodeRepository repository.ApprovalNodeRepository,
	approvalDefinitionRepository repository.ApprovalDefinitionRepository,
) ApprovalNodeService {
	return &approvalNodeService{
		Service:                      service,
		approvalNodeRepository:       approvalNodeRepository,
		approvalDefinitionRepository: approvalDefinitionRepository,
	}
}

// 基础CRUD操作
func (s *approvalNodeService) Get(id uint) (*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindOne(id)
}

func (s *approvalNodeService) GetByCode(code string) (*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FirstByCode(code)
}

func (s *approvalNodeService) First(where map[string]any) (*model.ApprovalNode, error) {
	return s.approvalNodeRepository.First(where)
}

func (s *approvalNodeService) List(page, pageSize int, total *int64, where map[string]any) ([]*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindPage(page, pageSize, total, where)
}

func (s *approvalNodeService) Create(approvalNode *model.ApprovalNode) error {
	// 验证节点配置
	if err := s.ValidateNode(approvalNode); err != nil {
		s.logger.Error("节点验证失败", "error", err)
		return err
	}

	// 验证审批定义是否存在
	_, err := s.approvalDefinitionRepository.FirstByCode(approvalNode.ApprovalDefCode)
	if err != nil {
		return errors.New("审批定义不存在")
	}

	return s.approvalNodeRepository.Create(approvalNode)
}

func (s *approvalNodeService) Update(approvalNode *model.ApprovalNode) error {
	// 验证节点配置
	if err := s.ValidateNode(approvalNode); err != nil {
		s.logger.Error("节点验证失败", "error", err)
		return err
	}

	return s.approvalNodeRepository.Update(approvalNode)
}

func (s *approvalNodeService) BatchUpdate(ids []uint, approvalNode *model.ApprovalNode) error {
	return s.approvalNodeRepository.BatchUpdate(ids, approvalNode)
}

func (s *approvalNodeService) Delete(id uint) (*model.ApprovalNode, error) {
	// 检查是否可以删除
	canDelete, err := s.CanDeleteNode(id)
	if err != nil {
		return nil, err
	}
	if !canDelete {
		return nil, errors.New("当前节点不允许删除")
	}

	// 先获取记录
	node, err := s.approvalNodeRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	// 执行删除
	err = s.approvalNodeRepository.BatchDelete([]uint{id})
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *approvalNodeService) BatchDelete(ids []uint) error {
	// 检查每个节点是否可以删除
	for _, id := range ids {
		canDelete, err := s.CanDeleteNode(id)
		if err != nil {
			return err
		}
		if !canDelete {
			return errors.New("存在不允许删除的节点")
		}
	}

	return s.approvalNodeRepository.BatchDelete(ids)
}

// 业务方法
func (s *approvalNodeService) GetByApprovalDefCode(approvalDefCode string) ([]*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindByApprovalDefCode(approvalDefCode)
}

func (s *approvalNodeService) GetActiveByApprovalDefCode(approvalDefCode string) ([]*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindActiveByApprovalDefCode(approvalDefCode)
}

func (s *approvalNodeService) GetByNodeType(nodeType string) ([]*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindByNodeType(nodeType)
}

func (s *approvalNodeService) GetStartNode(approvalDefCode string) (*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindStartNode(approvalDefCode)
}

func (s *approvalNodeService) GetEndNodes(approvalDefCode string) ([]*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindEndNodes(approvalDefCode)
}

func (s *approvalNodeService) GetNextNodes(approvalDefCode, currentNodeCode string) ([]*model.ApprovalNode, error) {
	return s.approvalNodeRepository.FindNextNodes(approvalDefCode, currentNodeCode)
}

// 流程验证
func (s *approvalNodeService) ValidateWorkflow(approvalDefCode string) error {
	nodes, err := s.approvalNodeRepository.FindActiveByApprovalDefCode(approvalDefCode)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return errors.New("工作流缺少节点")
	}

	// 检查开始节点
	startNodes := 0
	endNodes := 0
	for _, node := range nodes {
		if node.IsStartNode() {
			startNodes++
		}
		if node.IsEndNode() {
			endNodes++
		}
	}

	if startNodes == 0 {
		return errors.New("工作流缺少开始节点")
	}
	if startNodes > 1 {
		return errors.New("工作流只能有一个开始节点")
	}
	if endNodes == 0 {
		return errors.New("工作流缺少结束节点")
	}

	return nil
}

func (s *approvalNodeService) ValidateNode(node *model.ApprovalNode) error {
	if node.ApprovalDefCode == "" {
		return errors.New("审批定义编码不能为空")
	}
	if node.NodeName == "" {
		return errors.New("节点名称不能为空")
	}
	if node.NodeType == "" {
		return errors.New("节点类型不能为空")
	}
	if !model.IsValidNodeType(node.NodeType) {
		return errors.New("无效的节点类型")
	}

	// 审批节点必须配置审批人
	if node.IsApprovalNode() {
		if node.ApproverType == "" {
			return errors.New("审批节点必须配置审批人类型")
		}
		if !model.IsValidApproverType(node.ApproverType) {
			return errors.New("无效的审批人类型")
		}
		// if node.ApprovalMode != "" && !model.IsValidApprovalMode(node.ApprovalMode) {
		// 	return errors.New("无效的审批模式")
		// }
	}

	return nil
}

func (s *approvalNodeService) CanDeleteNode(id uint) (bool, error) {
	node, err := s.approvalNodeRepository.FindOne(id)
	if err != nil {
		return false, err
	}

	// 检查审批定义状态
	def, err := s.approvalDefinitionRepository.FirstByCode(node.ApprovalDefCode)
	if err != nil {
		return false, err
	}

	// 如果审批定义是激活状态，不允许删除节点
	if def.Status == model.ApprovalDefStatusNormal {
		return false, nil
	}

	return true, nil
}

// 节点配置
func (s *approvalNodeService) ConfigureApprovers(nodeId uint, approverConfig string) error {
	node, err := s.approvalNodeRepository.FindOne(nodeId)
	if err != nil {
		return err
	}

	if !node.IsApprovalNode() {
		return errors.New("只有审批节点才能配置审批人")
	}

	node.ApproverConfig = approverConfig
	return s.approvalNodeRepository.Update(node)
}

func (s *approvalNodeService) ConfigureConditions(nodeId uint, conditionConfig string) error {
	node, err := s.approvalNodeRepository.FindOne(nodeId)
	if err != nil {
		return err
	}

	if !node.IsConditionNode() {
		return errors.New("只有条件节点才能配置条件")
	}

	node.ConditionConfig = conditionConfig
	return s.approvalNodeRepository.Update(node)
}

func (s *approvalNodeService) ConfigureTimeouts(nodeId uint, timeoutHours int) error {
	node, err := s.approvalNodeRepository.FindOne(nodeId)
	if err != nil {
		return err
	}

	// node.TimeoutHours = timeoutHours
	return s.approvalNodeRepository.Update(node)
}

// 状态管理
func (s *approvalNodeService) ActivateNode(id uint) error {
	return s.approvalNodeRepository.UpdateStatus(id, model.ApprovalDefStatusNormal)
}

func (s *approvalNodeService) DeactivateNode(id uint) error {
	return s.approvalNodeRepository.UpdateStatus(id, model.ApprovalDefStatusFrozen)
}

// BatchSyncNodes 批量同步节点（增删改）
func (s *approvalNodeService) BatchSyncNodes(approvalDefCode string, nodes []*model.ApprovalNode) (string, error) {
	// 简化实现：逐个处理节点，不使用事务
	// 在生产环境中，建议使用事务管理器来确保数据一致性

	// 1. 获取现有节点
	existingNodes, err := s.approvalNodeRepository.FindByApprovalDefCode(approvalDefCode)
	if err != nil {
		return "", fmt.Errorf("获取现有节点失败: %w", err)
	}

	// 2. 构建现有节点映射（以NodeCode为key）
	existingMap := make(map[string]*model.ApprovalNode)
	for _, node := range existingNodes {
		existingMap[node.NodeCode] = node
	}

	// 3. 构建新节点映射
	newNodeMap := make(map[string]*model.ApprovalNode)
	for _, node := range nodes {
		newNodeMap[node.NodeCode] = node
	}

	// 4. 结果统计
	type SyncResult struct {
		Created []uint `json:"created"`
		Updated []uint `json:"updated"`
		Deleted []uint `json:"deleted"`
		Total   int    `json:"total"`
	}
	result := &SyncResult{
		Created: make([]uint, 0),
		Updated: make([]uint, 0),
		Deleted: make([]uint, 0),
	}

	// 5. 处理新增和更新
	for _, node := range nodes {
		if existingNode, exists := existingMap[node.NodeCode]; exists {
			// 更新现有节点
			existingNode.NodeName = node.NodeName
			existingNode.NodeType = node.NodeType
			existingNode.ApproverType = node.ApproverType
			existingNode.ApproverConfig = node.ApproverConfig
			existingNode.ConditionConfig = node.ConditionConfig
			existingNode.SortOrder = node.SortOrder

			if err := s.approvalNodeRepository.Update(existingNode); err != nil {
				return "", fmt.Errorf("更新节点失败: %w", err)
			}
			result.Updated = append(result.Updated, existingNode.ID)
		} else {
			// 创建新节点
			newNode := &model.ApprovalNode{
				ApprovalDefCode: approvalDefCode,
				NodeCode:        node.NodeCode,
				NodeName:        node.NodeName,
				NodeType:        node.NodeType,
				ApproverType:    node.ApproverType,
				ApproverConfig:  node.ApproverConfig,
				ConditionConfig: node.ConditionConfig,
				SortOrder:       node.SortOrder,
				Status:          model.ApprovalDefStatusNormal,
			}

			// 验证节点
			if err := s.ValidateNode(newNode); err != nil {
				return "", fmt.Errorf("节点验证失败: %w", err)
			}

			if err := s.approvalNodeRepository.Create(newNode); err != nil {
				return "", fmt.Errorf("创建节点失败: %w", err)
			}
			result.Created = append(result.Created, newNode.ID)
		}
	}

	// 6. 处理删除（现有节点中不在新节点列表中的）
	for nodeCode, existingNode := range existingMap {
		if _, exists := newNodeMap[nodeCode]; !exists {
			if err := s.approvalNodeRepository.BatchDelete([]uint{existingNode.ID}); err != nil {
				return "", fmt.Errorf("删除节点失败: %w", err)
			}
			result.Deleted = append(result.Deleted, existingNode.ID)
		}
	}

	result.Total = len(result.Created) + len(result.Updated) + len(result.Deleted)

	s.logger.Info("BatchSyncNodes",
		"approvalDefCode", approvalDefCode,
		"created", len(result.Created),
		"updated", len(result.Updated),
		"deleted", len(result.Deleted),
	)

	return fmt.Sprintf("同步完成: 创建%d个, 更新%d个, 删除%d个",
		len(result.Created), len(result.Updated), len(result.Deleted)), nil
}
