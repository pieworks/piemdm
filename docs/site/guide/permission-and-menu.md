# Permission & Menu Management

PieMDM's menu system is tightly integrated with the permission system, supporting dynamic menu rendering. This means that after logging in, users can only see the menu items they have permission for.

## Menu Management

The menu management function allows administrators to configure the system's navigation structure.

### 1. Menu List
Displays all menus and buttons of the system in a tree structure.
- **Directory**: Top-level navigation, usually does not directly correspond to a page.
- **Menu**: Specific page entry.
- **Button**: Operations within the page and their corresponding resource identifiers.

### 2. Add/Edit Menu
When configuring a menu, the main fields include:

- **Parent Menu**: Select parent node.
- **Menu Type**: Directory / Menu / Button.
- **Menu Name**: Title displayed in the sidebar.
- **Route Path**: Frontend route path (e.g., `/system/user`).
- **Component Path**: File path of the Vue component (e.g., `@/views/system/role/index.vue`).
- **Permission Key**: Key identifier for backend authentication (e.g., `system:user:list`).
- **Visible**: Controls whether the menu is displayed in the sidebar (some functional pages may need to be hidden).
- **Sort**: Controls the display order of the menu.

## Permission Control Mechanism

### Menu Visibility
When the frontend application initializes, it calls the API to get the current user's permission list. The system filters the accessible routing table based on the user's permission keys and dynamically mounts them. If a user does not have permission for a menu, that menu will not be rendered in the sidebar.

### Button-Level Permission
Inside the page, directives or functions control the display of buttons.
For example, only users with `system:user:add` permission can see the "Create User" button.

```html
<!-- Example: Vue directive controlling button permission -->
<el-button v-hasPerm="['system:user:add']">Create User</el-button>
```

### API Interface Authentication
Frontend permission control is only for improving user experience; core security depends on backend interface authentication.
Backend middleware intercepts requests and verifies whether the current user's role contains the **Permission Key** required for the request. If permission is insufficient, `403 Forbidden` is returned.

## Configuration Suggestions

1. **Clear Hierarchy**: It is recommended to organize permission keys according to the "Module -> Function -> Operation" hierarchy (e.g., `system` -> `user` -> `add`).
2. **Unique Key**: Ensure that the permission key for each menu and button is unique in the system.
3. **Synchronous Maintenance**: When adding backend interfaces, synchronize the addition of corresponding button/resource permissions in menu management and assign them to appropriate roles.
