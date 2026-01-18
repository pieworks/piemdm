package service_test

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"piemdm/internal/model"
	"piemdm/internal/service"
	"piemdm/pkg/helper/sid"
	"piemdm/pkg/jwt"
	"piemdm/pkg/log"
	mock_repository "piemdm/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupApprovalNodeService(t *testing.T) (service.ApprovalNodeService, *mock_repository.MockApprovalNodeRepository, *mock_repository.MockApprovalDefinitionRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	logger := &log.Logger{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	mockNodeRepo := mock_repository.NewMockApprovalNodeRepository(ctrl)
	mockDefRepo := mock_repository.NewMockApprovalDefinitionRepository(ctrl)

	baseService := service.NewService(logger, &sid.Sid{}, &jwt.JWT{})

	nodeService := service.NewApprovalNodeService(baseService, mockNodeRepo, mockDefRepo)
	return nodeService, mockNodeRepo, mockDefRepo, ctrl
}

func TestApprovalNodeService_Create(t *testing.T) {
	nodeService, mockNodeRepo, mockDefRepo, ctrl := setupApprovalNodeService(t)
	defer ctrl.Finish()

	now := time.Now()
	node := &model.ApprovalNode{
		ID:              1,
		ApprovalDefCode: "TEST_DEF_001",
		NodeCode:        "node_001",
		NodeType:        "APPROVAL",
		NodeName:        "审批节点1",
		ApproverType:    "USERS",
		ApproverConfig:  `{"users":["user001"]}`,
		SortOrder:       1,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Mock 审批定义存在性检查
	approvalDef := &model.ApprovalDefinition{
		ID:   1,
		Code: "TEST_DEF_001",
		Name: "测试审批定义",
	}
	mockDefRepo.EXPECT().FirstByCode("TEST_DEF_001").Return(approvalDef, nil)

	// Mock 创建节点
	mockNodeRepo.EXPECT().Create(node).Return(nil)

	err := nodeService.Create(node)
	assert.NoError(t, err)
}

func TestApprovalNodeService_Get(t *testing.T) {
	nodeService, mockNodeRepo, _, ctrl := setupApprovalNodeService(t)
	defer ctrl.Finish()

	id := uint(1)
	expectedNode := &model.ApprovalNode{
		ID:              1,
		ApprovalDefCode: "TEST_DEF_001",
		NodeCode:        "node_001",
		NodeType:        "APPROVAL",
		NodeName:        "审批节点1",
		ApproverType:    "USERS",
		ApproverConfig:  `{"users":["user001"]}`,
		SortOrder:       1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockNodeRepo.EXPECT().FindOne(id).Return(expectedNode, nil)

	node, err := nodeService.Get(id)
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "node_001", node.NodeCode)
	assert.Equal(t, "审批节点1", node.NodeName)
}

func TestApprovalNodeService_GetByDefCode(t *testing.T) {
	nodeService, mockNodeRepo, _, ctrl := setupApprovalNodeService(t)
	defer ctrl.Finish()

	defCode := "TEST_DEF_001"
	expectedNodes := []*model.ApprovalNode{
		{
			ID:              1,
			ApprovalDefCode: "TEST_DEF_001",
			NodeCode:        "start",
			NodeType:        "START",
			NodeName:        "开始节点",
			SortOrder:       0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              2,
			ApprovalDefCode: "TEST_DEF_001",
			NodeCode:        "node_001",
			NodeType:        "APPROVAL",
			NodeName:        "审批节点1",
			ApproverType:    "USERS",
			ApproverConfig:  `{"users":["user001"]}`,
			SortOrder:       1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockNodeRepo.EXPECT().FindByApprovalDefCode(defCode).Return(expectedNodes, nil)

	nodes, err := nodeService.GetByApprovalDefCode(defCode)
	assert.NoError(t, err)
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 2)
	assert.Equal(t, "start", nodes[0].NodeCode)
	assert.Equal(t, "node_001", nodes[1].NodeCode)
}

func TestApprovalNodeService_Delete(t *testing.T) {
	nodeService, mockNodeRepo, mockDefRepo, ctrl := setupApprovalNodeService(t)
	defer ctrl.Finish()

	id := uint(1)
	node := &model.ApprovalNode{
		ID:              id,
		ApprovalDefCode: "TEST_DEF_001",
		NodeCode:        "node_001",
		NodeType:        "APPROVAL",
	}

	approvalDef := &model.ApprovalDefinition{
		ID:     1,
		Code:   "TEST_DEF_001",
		Status: "Draft", // 草稿状态允许删除
	}

	// CanDeleteNode 检查：第一次调用 FindOne
	mockNodeRepo.EXPECT().FindOne(id).Return(node, nil)
	// CanDeleteNode 检查：调用 FirstByCode 检查审批定义状态
	mockDefRepo.EXPECT().FirstByCode("TEST_DEF_001").Return(approvalDef, nil)
	// Delete 方法：第二次调用 FindOne 获取删除前的记录
	mockNodeRepo.EXPECT().FindOne(id).Return(node, nil)
	// Delete 方法：调用 BatchDelete 执行实际删除
	mockNodeRepo.EXPECT().BatchDelete([]uint{id}).Return(nil)

	deletedNode, err := nodeService.Delete(id)
	assert.NoError(t, err)
	assert.NotNil(t, deletedNode)
}
