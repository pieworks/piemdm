# User & Role Management

PieMDM provides a comprehensive User and Role Management System (RBAC) for controlling system access and data operation permissions.

## Overview

The system's permission design is based on the **User - Role - Permission** model:
- **User**: The subject of system operations.
- **Role**: A collection of permissions that can be assigned to one or more users.
- **Permission**: Authorization for specific operations, such as "View Menu", "Edit Data", etc.

## User Management

Administrators can manage the full lifecycle of users in the system.

### 1. User List
On the User Management page, you can view all system users and their basic information, including:
- Username
- Email
- Status (Enabled/Disabled)
- Associated Roles
- Creation Time

### 2. Create User
Click the "Create User" button and fill in the necessary information:
- **Username**: Unique identifier used for login.
- **Password**: Initial login password.
- **Email**: Used for notifications and password recovery (optional).
- **Role**: Assign one or more initial roles to the user.

### 3. Edit User
Supports modifying usage basic information and role assignment.
> Note: Modifying roles may require the user to log in again to take effect.

### 4. Disable/Enable User
For accounts of resigned employees or suspended usage, use the "Disable" function to forbid the user from logging in without deleting data.

## Role Management

Through Role Management, you can define different functional permission groups.

### 1. Role List
Displays all roles defined in the system. The system usually includes basic roles like `admin` (Super Administrator) by default.

### 2. Create Role
When defining a new role, you need to specify:
- **Role Name**: E.g., "Data Entry Clerk", "Approval Manager".
- **Role Code**: Unique identifier for backend permission verification (e.g., `data_entry`).
- **Description**: Description of role responsibilities.

### 3. Assign Permissions
In the role details or permission configuration page, you can check specific permission points for the role.
- **Menu Permission**: Controls the left menu items visible after user login.
- **Operation Permission**: Controls buttons within the page (e.g., "Create", "Delete", "Export").

## Best Practices

1. **Principle of Least Privilege**: Upgrade permissions to users only as needed for their work.
2. **Role Reuse**: Try to manage permissions through roles and avoid hardcoding permissions for specific users.
3. **Regular Audit**: Regularly check the assignment of administrator accounts and key roles to ensure system security.
