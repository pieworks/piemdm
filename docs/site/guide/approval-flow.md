# Approval Flow Management

PieMDM has a powerful built-in approval engine that allows enterprises to flexibly configure approval processes based on business needs. Whether it is the creation of material master data, the modification of supplier profiles, or the freezing of financial data, approval flows can ensure data accuracy and compliance.

We also support integration with external approval systems:
- Feishu (Supported)
- DingTalk
- WeChat Work
- Others

## Core Concepts

The system's approval architecture consists of the following parts:
- **Approval Definition**: The template of the process, defining the approval name, associated forms, and node sequence.
- **Approval Node**: Specific steps in the process, defining who approves and how.
- **Approval Instance**: The specific execution record of a process, corresponding to a specific approval document.
- **Approval Task**: To-do items assigned to specific executors.

## Approval Node Types

When configuring approval flows, you can combine multiple types of nodes:

| Node Type | Description |
| :--- | :--- |
| **Start Node (START)** | The entry point of the flow, usually initiated by the applicant. |
| **Approval Node (APPROVAL)** | Core node, requiring the approver to pass or reject. |
| **CC Node (CC)** | Notifies relevant personnel, no action required, for viewing only. |
| **Condition Node (CONDITION)** | Logic branch, automatically selecting a path based on form data (e.g., Amount > 10,000 goes to Manager Approval). |
| **Parallel/Merge Node** | Supports simultaneous approval of multiple steps, continuing after convergence. |
| **End Node (END)** | The endpoint of the flow; data typically becomes effective after passing. |

## Approval Modes

For approval nodes, you can configure different decision modes:
- **OR**: Multiple approvers in the node, only one needs to pass to proceed to the next step.
- **AND**: All approvers in the node must pass for the flow to continue.
- **SEQUENTIAL**: Approvals are performed one by one in a predefined order.

## Approver Settings

Approvers can be configured based on multiple dimensions:
- **Specific User**: Select specific system users.
- **Specific Role**: Assign to all users with a specific role (e.g., "Finance Manager").
- **Supervisor**: Dynamically retrieve the direct or indirect supervisor of the applicant.
- **Self-Selection**: The approver of the previous step temporarily specifies the executor of the next step.

## Data Flow Mechanism

### 1. Draft Data Isolation
When a table is associated with an approval flow, newly added or modified data is first stored in "Draft Status". At this time, the data is invisible in the business system and only circulates in the approval module.

### 2. Activation Strategy
- **Approval Passed**: After the flow reaches the "End Node" and passes, the system automatically synchronizes the draft data to the formal business table, and the master data officially takes effect.
- **Approval Rejected**: The flow terminates. The applicant can modify and resubmit based on approval comments, or withdraw the application.

## Best Practices

1. **Simplify Paths**: For routine, low-risk operations, try to shorten the approval path to improve efficiency.
2. **Clear Node Naming**: Use business-meaningful node names (e.g., "First Review: Dept Supervisor") to help participants understand progress.
3. **Backup Plans**: Considering personnel turnover or leave, it is recommended to assign approval tasks by "Role" rather than "Individual".
