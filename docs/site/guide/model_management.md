# Model Management (Data Modeling)

Model Management is the core function of PieMDM, defining the logical structure of master data. Through a visual interface, you can easily create entity tables and define their attribute fields, supporting highly flexible extensions.

## 1. Concept Definition

In PieMDM, a "Model" typically corresponds to a set of tables in the underlying database. When you define an entity table in the backend, the system automatically maintains the following three types of tables:

- **Main Table (t_{code})**: Stores effective, standard master data.
- **Draft Table (t_{code}_draft)**: Stores temporary data that is under approval or not yet published.
- **Log Table (t_{code}_log)**: Records all change trails of the entity data.

## 2. Entity Table Management

Click **"Basic Settings > Model Management"** to enter the list page:

1. **Create Entity**: Enter the entity name and unique "Entity Code".
2. **Naming Convention**: It is recommended to use UpperCamelCase or snake_case (e.g., `User` or `material_info`). **Note: Once an entity code is created, it cannot be modified.**
3. **Relationships**: You can set the hierarchical relationship of tables, such as setting "Office Address" as a child table of "Organization".

## 3. Field (Attribute) Definition

In the entity detail page, you can manage all attribute fields of the entity.

### 3.1 Core Attributes
- **Field Type**: The system has built-in types:
  - **Basic**: Single-line text, Multi-line text, Number, Boolean.
  - **Options**: Select, Checkbox, Relation (Dropdown).
  - **Media**: Image upload, Attachment upload.
  - **Date**: Date (DatePicker), Time, DateTime.
- **Required**: If checked, this field cannot be empty during entity data maintenance.
- **Unique**: If checked, the system will automatically create unique indexes in the main and draft tables to prevent duplicate data.
- **Filter**: If checked, this field will appear in the search area of the entity data page.

### 3.2 Relation Fields
Relation fields are advanced modeling features that support referencing data from other entity tables:
- **Related Table**: Select the target entity table.
- **Display Field**: Select the field shown to the user in the dropdown list (e.g., select by ID but display "Name").

## 4. Schema Publish

**Key Operation**: After modifying fields (adding, deleting, or adjusting length), the model does not take effect in the database immediately.

Click the **"Publish"** button in the top right corner of the page. The system will compare the current definition with the physical database structure and automatically execute `ALTER TABLE` statements to synchronize the main table, draft table, and log table.

<callout emoji="⚠️" background-color="light-yellow">
Note: Before executing the "Publish" operation, please ensure that all approval tasks for the current entity have been processed to avoid data structure conflicts.
</callout>

## 5. FAQ

<callout emoji="❓" background-color="light-purple">
Will modifying a field code cause data loss?
</callout>

Field codes directly correspond to database column names. Modifying a published field code will cause the system to execute "Delete old column, Create new column" operations, resulting in the **loss of data in the old column**. It is recommended to confirm in the development environment before making changes.

<callout emoji="❓" background-color="light-purple">
What is the effect of field sorting?
</callout>

The sort value (Sort) of a field determines the order in which the field appears in the "Entity Management" form, as well as the column order in the table view.
