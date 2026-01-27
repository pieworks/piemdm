# Table & Field Management

PieMDM allows administrators to customize data models based on business needs. Through "Table & Field Management", you can define master data entities and their attributes.

## Table Management

In PieMDM, a "Table" represents a master data object (e.g., "Supplier", "Material", "Customer").

### 1. Table Attributes
When creating a new table, you need to configure the following key attributes:

- **Table Code (Code)**: Unique identifier in the system, restricted to lowercase letters, numbers, and underscores (e.g., `mdm_material`). The system automatically adds a `t_` prefix to the physical table name.
- **Name**: Displayed business name (e.g., "Material Master").
- **Display Mode**:
  - `List`: List mode, suitable for most flat data.
  - `Tree`: Tree mode, suitable for hierarchical data (e.g., "Organization", "Category Dictionary").
- **Table Type**:
  - `Entity`: Independent master entity table.
  - `Item`: Line item table, must form a relationship with a parent table.

### 2. Relation Configuration (Item Type Only)
If the table type is `Item`, additional configuration is required:
- **Parent Table**: The main table this line item belongs to.
- **Relation Field**: Define the foreign key between the child table and the parent table.

## Field Management

Fields define the specific data items stored in the table. PieMDM provides rich field types and validation rules.

### 1. Field Type Presets
Field types are divided into several categories:

| Group | Types | Description |
| :--- | :--- | :--- |
| **Basic** | Text, Multi-line Text, Number, Integer, Percentage, Password, URL, Email, Phone | Basic data storage |
| **Selection** | Select, Multi-Select, Radio Group, Checkbox Group, Switch | Predefined option data |
| **DateTime** | Date, Time, DateTime | Time dimension data |
| **Relation** | ManyToOne (BelongsTo), OneToMany (HasMany), ManyToMany | Cross-table data association |
| **Advanced** | AutoCode, Attachment, RichText, Formula | Complex business scenario support |

### 2. Field Configuration Items
- **Field Code**: Column name in the database.
- **Display Name**: Label in forms and lists.
- **UI Component**: Defines how the field is presented on the interface (e.g., `Input`, `Select`, `DatePicker`).
- **Validation Rules**:
  - **Required**: Whether it can be empty.
  - **Length Limit**: Max/Min character count.
  - **Numeric Range**: Max/Min value.
  - **Regex**: Custom format validation.

## System Reserved Fields

To support auditing, version management, and approval workflows, PieMDM automatically adds the following reserved fields to each table. Users must not create custom fields with the same codes as these fields.

- **Audit Fields**: `id`, `created_at`, `updated_at`, `deleted_at`, `created_by`, `updated_by`.
- **Business Status**: `status` (Data Status), `draft_status` (Draft Status).
- **Process Tracking**: `approval_code` (Approval No.), `entity_id` (Original Entity ID).

## Operation Process

1. **New Table**: Define basic table information and storage engine configuration.
2. **Define Fields**: Add fields one by one according to business attributes, configuring their types and validation rules.
3. **Generate Physical Table**: After saving the configuration, the system automatically creates or updates the corresponding physical table structure in the database.
4. **Permission Assignment**: Configure access entry for the new table in "Menu Management" and assign view and edit permissions to corresponding roles.
