package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"piemdm/internal/model"
)

// TestBatchCreateApprovalNodes 测试批量创建审批节点
func TestBatchCreateApprovalNodes(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 准备测试数据
	testNodes := []model.ApprovalNode{
		{
			ApprovalDefCode: "test_approval_001",
			NodeCode:        "start_node",
			NodeName:        "开始节点",
			NodeType:        model.NodeTypeStart,
			SortOrder:       1,
		},
		{
			ApprovalDefCode: "test_approval_001",
			NodeCode:        "approval_node_001",
			NodeName:        "部门经理审批",
			NodeType:        model.NodeTypeApproval,
			ApproverType:    model.ApproverTypeUsers,
			ApproverConfig:  `{"users":["user1","user2"]}`,
			SortOrder:       2,
		},
		{
			ApprovalDefCode: "test_approval_001",
			NodeCode:        "approval_node_002",
			NodeName:        "总经理审批",
			NodeType:        model.NodeTypeApproval,
			ApproverType:    model.ApproverTypeRoles,
			ApproverConfig:  `{"roles":["general_manager"]}`,
			SortOrder:       3,
		},
		{
			ApprovalDefCode: "test_approval_001",
			NodeCode:        "cc_node",
			NodeName:        "抄送节点",
			NodeType:        model.NodeTypeCC,
			ApproverConfig:  `{"users":["hr_manager"]}`,
			SortOrder:       4,
		},
		{
			ApprovalDefCode: "test_approval_001",
			NodeCode:        "end_node",
			NodeName:        "结束节点",
			NodeType:        model.NodeTypeEnd,
			SortOrder:       5,
		},
	}

	// 转换为JSON
	jsonData, err := json.Marshal(testNodes)
	require.NoError(t, err)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "/api/v1/approval-nodes/batch", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 这里应该使用实际的路由器，但为了测试目的，我们模拟一个简单的处理器
	// 在实际项目中，你需要设置完整的路由和依赖注入
	router := gin.New()
	router.POST("/api/v1/approval-nodes/batch", func(c *gin.Context) {
		var nodes []model.ApprovalNode
		if err := c.ShouldBindJSON(&nodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证每个节点
		for i, node := range nodes {
			// 模拟BeforeCreate钩子验证
			if err := node.BeforeCreate(nil); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      "validation failed",
					"node_index": i,
					"node_code":  node.NodeCode,
				})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "批量创建成功",
			"count":   len(nodes),
		})
	})

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "批量创建成功", response["message"])
	assert.Equal(t, float64(5), response["count"]) // JSON数字默认为float64
}

// TestBatchCreateApprovalNodes_ValidationError 测试批量创建时的验证错误
func TestBatchCreateApprovalNodes_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 准备包含无效数据的测试数据
	testNodes := []model.ApprovalNode{
		{
			ApprovalDefCode: "test_approval_002",
			NodeCode:        "start_node",
			NodeName:        "开始节点",
			NodeType:        model.NodeTypeStart,
			SortOrder:       1,
		},
		{
			ApprovalDefCode: "test_approval_002",
			NodeCode:        "invalid_approval_node",
			NodeName:        "无效审批节点",
			NodeType:        model.NodeTypeApproval,
			ApproverType:    "INVALID_APPROVER_TYPE", // 无效的审批人类型
			SortOrder:       2,
		},
	}

	jsonData, err := json.Marshal(testNodes)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/approval-nodes/batch", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/approval-nodes/batch", func(c *gin.Context) {
		var nodes []model.ApprovalNode
		if err := c.ShouldBindJSON(&nodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证每个节点
		for i, node := range nodes {
			if err := node.BeforeCreate(nil); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      "validation failed",
					"node_index": i,
					"node_code":  node.NodeCode,
				})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "批量创建成功",
			"count":   len(nodes),
		})
	})

	router.ServeHTTP(w, req)

	// 验证应该返回验证错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "validation failed", response["error"])
	assert.Equal(t, float64(1), response["node_index"]) // 第二个节点（索引1）验证失败
	assert.Equal(t, "invalid_approval_node", response["node_code"])
}

// TestBatchCreateApprovalNodes_MixedNodeTypes 测试混合节点类型的批量创建
func TestBatchCreateApprovalNodes_MixedNodeTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试各种节点类型的组合
	testNodes := []model.ApprovalNode{
		{
			ApprovalDefCode: "test_approval_003",
			NodeCode:        "start",
			NodeName:        "开始",
			NodeType:        model.NodeTypeStart,
			ApproverType:    "SOME_INVALID_TYPE", // 开始节点不验证审批人类型
		},
		{
			ApprovalDefCode: "test_approval_003",
			NodeCode:        "condition",
			NodeName:        "条件判断",
			NodeType:        model.NodeTypeCondition,
			ApproverType:    "ANOTHER_INVALID_TYPE", // 条件节点不验证审批人类型
		},
		{
			ApprovalDefCode: "test_approval_003",
			NodeCode:        "parallel",
			NodeName:        "并行审批",
			// NodeType:        model.NodeTypeParallel,
			ApproverType: "YET_ANOTHER_INVALID_TYPE", // 并行节点不验证审批人类型
		},
		{
			ApprovalDefCode: "test_approval_003",
			NodeCode:        "approval",
			NodeName:        "实际审批",
			NodeType:        model.NodeTypeApproval,
			ApproverType:    "DEPARTMENTS", // 只有审批节点验证审批人类型
			ApproverConfig:  `{"departments":["dept1","dept2"]}`,
		},
		{
			ApprovalDefCode: "test_approval_003",
			NodeCode:        "merge",
			NodeName:        "合并节点",
			// NodeType:        model.NodeTypeMerge,
			ApproverType: "FINAL_INVALID_TYPE", // 合并节点不验证审批人类型
		},
		{
			ApprovalDefCode: "test_approval_003",
			NodeCode:        "end",
			NodeName:        "结束",
			NodeType:        model.NodeTypeEnd,
			ApproverType:    "END_INVALID_TYPE", // 结束节点不验证审批人类型
		},
	}

	jsonData, err := json.Marshal(testNodes)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/approval-nodes/batch", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/approval-nodes/batch", func(c *gin.Context) {
		var nodes []model.ApprovalNode
		if err := c.ShouldBindJSON(&nodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证每个节点
		for i, node := range nodes {
			if err := node.BeforeCreate(nil); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      "validation failed",
					"node_index": i,
					"node_code":  node.NodeCode,
				})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "批量创建成功",
			"count":   len(nodes),
		})
	})

	router.ServeHTTP(w, req)

	// 验证应该成功，因为只有审批节点会验证审批人类型
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "批量创建成功", response["message"])
	assert.Equal(t, float64(6), response["count"])
}
