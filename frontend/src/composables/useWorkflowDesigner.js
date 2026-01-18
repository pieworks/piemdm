import { createApprovalDef, findApprovalDef, updateApprovalDef } from "@/api/approval_def";
import { ref } from "vue";

// A map to get human-readable names for node types
const NODE_TYPE_NAMES = {
  APPROVAL: "审批节点",
  CC: "抄送节点",
  CONDITION: "条件分支",
  CONDITION_BRANCH: "条件",
};

/**
 * Creates a new node object.
 * @param {string} type - The type of the node (e.g., 'APPROVAL').
 * @param {object} parentNode - The node that will be the parent of the new node.
 * @returns {object} A new node object.
 */
function createNode(type, parentNode) {
  const parentId = parentNode.id === "start" ? null : parentNode.id;
  const baseNode = {
    id: `node_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`,
    name: NODE_TYPE_NAMES[type] || "新节点",
    nodeType: type,
    parentId: parentId,
    childNode: null,
  };

  if (type === "CONDITION") {
    baseNode.childNodes = [];
  } else if (type === "APPROVAL") {
    baseNode.approverConfig = {}; // Placeholder for approver settings
  } else if (type === "CC") {
    baseNode.ccConfig = {}; // Placeholder for CC settings
  } else if (type === "CONDITION_BRANCH") {
    baseNode.conditionConfig = {}; // Placeholder for condition settings
  }

  return baseNode;
}

/**
 * A Vue 3 composable for managing the workflow designer logic.
 * @param {string} definitionId - The ID of the approval definition to load.
 */
export function useWorkflowDesigner(definitionId) {
  // --- STATE ---
  const approvalDef = ref({ id: definitionId, name: "", description: "", processConfig: null });
  const nodeTree = ref(null);
  const flatNodes = ref(new Map());
  const selectedNode = ref(null);
  const addMenu = ref({ show: false, x: 0, y: 0, parentNode: null });

  // --- PRIVATE HELPERS ---

  /**
   * Converts a flat list of nodes into a tree structure.
   * @param {Array} list - A flat array of node objects.
   * @returns {object} The root node of the tree.
   */
  const buildTreeFromFlatList = (list) => {
    flatNodes.value.clear();
    if (!list || list.length === 0) {
      return { id: "start", name: "发起人", nodeType: "START", childNode: null };
    }

    const nodeMap = new Map();
    list.forEach((node) => {
      const newNode = { ...node, childNode: null, childNodes: [] };
      nodeMap.set(node.id, newNode);
      flatNodes.value.set(node.id, newNode);
    });

    const rootNodes = [];
    list.forEach((node) => {
      if (node.parentId && nodeMap.has(node.parentId)) {
        const parent = nodeMap.get(node.parentId);
        if (parent.nodeType === "CONDITION") {
          parent.childNodes.push(nodeMap.get(node.id));
        } else {
          parent.childNode = nodeMap.get(node.id);
        }
      } else {
        rootNodes.push(nodeMap.get(node.id));
      }
    });

    // TODO: This assumes a single linear root. Needs enhancement for more complex start scenarios.
    return {
      id: "start",
      name: "发起人",
      nodeType: "START",
      childNode: rootNodes[0] || null,
    };
  };

  /**
   * Flattens the tree structure back into a list for saving.
   * @param {object} root - The root node of the tree.
   * @returns {Array} A flat array of node objects.
   */
  const flattenTreeToList = (root) => {
    const list = [];
    const seen = new Set(); // To handle potential cycles, though our structure should prevent them.

    function traverse(node) {
      if (!node || seen.has(node.id) || node.nodeType === "START") {
        return;
      }

      seen.add(node.id);
      const { childNode, childNodes, ...nodeData } = node; // Exclude children from the shallow copy
      list.push(nodeData);

      if (childNode) {
        traverse(childNode);
      }
      if (childNodes && childNodes.length > 0) {
        childNodes.forEach(traverse);
      }
    }

    traverse(root.childNode);
    return list;
  };

  // --- ACTIONS ---

  /**
   * Loads the approval definition from the server.
   */
  const load = async () => {
    if (!definitionId) {
      nodeTree.value = buildTreeFromFlatList([]);
      return;
    }
    try {
      const response = await findApprovalDef({ id: definitionId });
      approvalDef.value = response.data;
      const config = response.data.processConfig
        ? JSON.parse(response.data.processConfig)
        : { nodes: [] };
      const nodes = (config.nodes || []).map((node) => ({
        ...node,
        approverConfig: node.approverConfig ? JSON.parse(node.approverConfig) : {},
        ccConfig: node.ccConfig ? JSON.parse(node.ccConfig) : {},
        conditionConfig: node.conditionConfig ? JSON.parse(node.conditionConfig) : {},
      }));
      nodeTree.value = buildTreeFromFlatList(nodes);
    } catch (error) {
      console.error("Failed to load approval definition:", error);
      alert("加载失败: " + (error.response?.data?.message || error.message));
      nodeTree.value = buildTreeFromFlatList([]);
    }
  };

  /**
   * Saves the current workflow to the server.
   */
  const save = async () => {
    try {
      const nodesToSave = flattenTreeToList(nodeTree.value);
      const processConfig = {
        nodes: nodesToSave.map((n) => ({
          ...n,
          approverConfig: JSON.stringify(n.approverConfig || {}),
          ccConfig: JSON.stringify(n.ccConfig || {}),
          conditionConfig: JSON.stringify(n.conditionConfig || {}),
        })),
      };
      const data = { ...approvalDef.value, processConfig: JSON.stringify(processConfig) };

      if (data.id) {
        await updateApprovalDef(data);
      } else {
        const res = await createApprovalDef(data);
        approvalDef.value.id = res.data.id; // Update id after creation
      }
      alert("保存成功");
    } catch (error) {
      console.error("Failed to save:", error);
      alert("保存失败: " + (error.response?.data?.message || error.message));
    }
  };

  /**
   * Displays the 'Add Node' menu at the specified event location.
   * @param {object} parentNode - The node after which to add the new node.
   * @param {MouseEvent} event - The click event.
   */
  const showAddMenu = (parentNode, event) => {
    addMenu.value = { show: true, x: event.clientX, y: event.clientY, parentNode };
  };

  /**
   * Adds a new node to the tree.
   * @param {string} nodeType - The type of node to add.
   */
  const addNewNode = (nodeType) => {
    const parent = addMenu.value.parentNode;
    if (!parent) return;

    const newNode = createNode(nodeType, parent);
    const oldChild = parent.childNode;

    newNode.childNode = oldChild;
    if (oldChild) oldChild.parentId = newNode.id;

    parent.childNode = newNode;

    if (newNode.nodeType === "CONDITION") {
      const branch1 = createNode("CONDITION_BRANCH", newNode);
      const branch2 = createNode("CONDITION_BRANCH", newNode);
      branch1.name = "分支 1";
      branch2.name = "分支 2";
      branch1.childNode = oldChild; // The rest of the flow goes into the first branch
      if (oldChild) oldChild.parentId = branch1.id;
      newNode.childNode = null; // Condition node itself doesn't have a direct child
      newNode.childNodes.push(branch1, branch2);
    }
  };

  /**
   * Sets the currently selected node.
   * @param {object} nodeToSelect - The node to select.
   */
  const selectNode = (nodeToSelect) => {
    if (selectedNode.value?.id === nodeToSelect.id) {
      selectedNode.value = null;
    } else {
      selectedNode.value = nodeToSelect;
    }
  };

  /**
   * Finds a node's parent in the tree.
   * @param {object} startNode - The node to start searching from.
   * @param {string} childId - The ID of the node whose parent we are looking for.
   * @returns {object|null} The parent node or null if not found.
   */
  const findParent = (startNode, childId) => {
    if (!startNode) return null;

    if (startNode.childNode?.id === childId) return startNode;

    if (startNode.childNodes) {
      for (const branch of startNode.childNodes) {
        if (branch.id === childId) return startNode;
        const found = findParent(branch, childId);
        if (found) return found;
      }
    }

    return findParent(startNode.childNode, childId);
  };

  /**
   * Removes a node from the tree.
   * @param {object} nodeToRemove - The node to be removed.
   */
  const removeNode = (nodeToRemove) => {
    const parent = findParent(nodeTree.value, nodeToRemove.id);
    if (!parent) {
      console.error("Could not find parent for node:", nodeToRemove.id);
      return;
    }

    if (parent.nodeType === "CONDITION") {
      parent.childNodes = parent.childNodes.filter((n) => n.id !== nodeToRemove.id);
    } else {
      parent.childNode = nodeToRemove.childNode;
      if (nodeToRemove.childNode) {
        nodeToRemove.childNode.parentId = parent.id === "start" ? null : parent.id;
      }
    }
  };

  /**
   * Adds a new branch to a condition node.
   * @param {object} conditionNode - The condition node to add a branch to.
   */
  const addBranch = (conditionNode) => {
    if (conditionNode.nodeType !== "CONDITION") return;
    const newBranch = createNode("CONDITION_BRANCH", conditionNode);
    newBranch.name = `分支 ${conditionNode.childNodes.length + 1}`;
    conditionNode.childNodes.push(newBranch);
  };

  const updateNode = (updatedNode) => {
    const parent = findParent(nodeTree.value, updatedNode.id);
    if (!parent) return;

    let nodeToUpdate = null;
    if (parent.nodeType === "CONDITION") {
      nodeToUpdate = parent.childNodes.find((n) => n.id === updatedNode.id);
    } else {
      nodeToUpdate = parent.childNode;
    }

    if (nodeToUpdate) {
      Object.assign(nodeToUpdate, updatedNode);
      // deselect and reselect to force panel refresh if needed
      const currentId = selectedNode.value.id;
      selectedNode.value = null;
      selectedNode.value = findNodeById(nodeTree.value, currentId);
    }
  };

  // --- EXPORTS ---
  return {
    approvalDef,
    nodeTree,
    selectedNode,
    addMenu,
    load,
    save,
    showAddMenu,
    addNewNode,
    selectNode,
    removeNode,
    addBranch,
    updateNode,
  };
}
