# Approval Process Definition & Management

PieMDM provides a highly configurable approval engine. You can define personalized approval processes for different entities (tables) and operations (such as Create, Update, Delete) based on business needs.

## 1. Approval Definition

"Approval Definition" is the template of the process, describing how a complete approval flow runs.

### 1.1 Basic Configuration
- **Definition Name & Code**: Each approval flow has a unique Code and an easy-to-understand Name.
- **Status Management**:
    - **Draft**: Under editing, cannot be referenced.
    - **Active**: Officially effective, can be bound to entity operations.
    - **Inactive**: Temporarily offline; initiated flows will continue to execute, but new operations will no longer trigger this flow.
- **Publish & Version Control**: The system supports multi-version management. When you modify a published approval flow and choose "Republish", the system generates a new version number to ensure traceability of process changes.

### 1.2 Approval Platform
PieMDM supports multiple approval delivery methods:
- **Internal**: Process directly on the PieMDM web interface.
- **Feishu**: The system enables synchronization of Feishu approval forms and nodes, allowing users to complete approvals directly in the Feishu dialog box.

## 2. Approval Node

An approval flow is composed of multiple nodes in order.

### 2.1 Node Types
- **Start Node (Start)**: The starting point of the flow, usually triggered by the user submitting data.
- **Approval Node (Approval)**: Core manual review node.
- **CC Node (CC)**: Notifies relevant personnel, strictly for information purposes, does not affect flow progression.
- **Condition Branch (Condition)**: Automatically routes to different approval branches based on data content (e.g., Amount, Department).
- **End Node (End)**: The endpoint of the flow, marking that data officially takes effect.

### 2.2 Approver Configuration
You can flexibly define approvers for each node:
- **Specific Personnel**: Select specific people from the user list.
- **Specific Role**: E.g., "Finance Director"; all users belonging to this role can approve.
- **Dynamic Relationship**: E.g., "Direct Supervisor", "Department Head".

## 3. Table Binding

After defining the approval flow, it needs to be associated with specific entity table operations to function.

You can configure this in **"Model Management > Approval Settings"** or **"Process Management > Table Binding"**:
- **Entity Table**: Select the target entity (e.g., "Supplier Information").
- **Operation Type**:
    - **Create**: Does new data entry require approval?
    - **Update**: Does editing existing data require approval?
    - **Delete**: Does voiding data require approval?
- **Binding Flow**: Select an approval definition in "Active" status.

## 4. Process Publishing Mechanism

<callout emoji="🚀" background-color="light-green">
When you complete the node configuration, you must click **"Verify and Publish"**. The system checks if the flow is closed (i.e., has an end point), has isolated nodes, and if the approver configuration is valid.
</callout>

## 5. FAQ

<callout emoji="❓" background-color="light-purple">
Will running tasks change after modifying the approval definition?
</callout>

No. Approval instances that have already started will strictly run according to the process definition version **at the time of initiation**. Only newly initiated applications will use the latest published version.

<callout emoji="❓" background-color="light-purple">
Can a table be bound to multiple approval flows?
</callout>

For the same operation (e.g., "Update"), a table can only be bound to one **Active** flow. However, you can bind different flows for "Create" and "Update" respectively.
