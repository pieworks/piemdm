# Entity Management (Data Maintenance)

PieMDM provides flexible data maintenance capabilities, supporting automatic generation of dynamic forms based on metadata configuration. Users can maintain business entity data through page operations, batch imports, and other methods.

## 1. Core Process

In PieMDM, data maintenance follows this core process:

1. **Permission Verification**: The system first checks if the current user has maintenance permission for the table (configured via "Table Permission Management").
2. **Form Rendering**: The frontend automatically renders corresponding input controls based on the "Field Management" information defined for the table.
3. **Metadata Driven**: Controls like dropdowns, file uploads, and date pickers are all driven by the `options` configuration of the fields.
4. **Change Reason**: Whether creating or modifying, the system **mandatorily requires** filling in a "Change Reason", which is crucial for auditing and approval.
5. **Approval Flow Judgment**:
   - If the operation for the table (such as `Create` or `Update`) is configured with an [Approval Flow](./approval-flow.md), the data will first be saved to the `_draft` table and automatically trigger the approval process.
   - If no approval flow is configured, the data will be updated directly to the official business table.

## 2. Operation Guide

### 2.1 Create Entity

1. Enter the corresponding business menu.
2. Click the **"Create"** button at the top of the page.
3. **Fill in Change Reason**: At the top of the form, you must input the reason for this operation.
4. **Fill in Form Data**:
   - **Basic Fields**: Directly input text, numbers, booleans, etc.
   - **Auto Code**: If a field is defined as "Auto Code", no manual intervention is needed; the system automatically generates it upon saving (displayed as "Auto Generate").
   - **Relation Select**: The system automatically loads options based on the configured related table. Search filtering is supported if there are many options.
   - **Date/Time**: Use the built-in date/time picker, and the system usually automatically initializes the current time as the default value.
   - **File Upload**: Supports file drag-and-drop or click-to-upload, and configuration for multi-selection and file type restrictions.
5. **Submit**: Click the "Submit" button, and the system will proceed with subsequent processing based on the approval configuration.

### 2.2 Edit and Update

1. Find the target data in the list page and click **"Edit"**.
2. The system will extract the latest current data to populate the form.
3. After modification, you must fill in a **new change reason**.
4. After submission, if there is an approval flow, since the data is in "Pending Approval" status, the list page usually indicates that the data is in process, prohibiting further editing until the process ends.

## 3. Advanced Maintenance Functions

### 3.1 Batch Import

PieMDM supports quick import of large amounts of data via Excel:

1. **Download Template**: Click "Download Template", and the system will generate an Excel template containing all non-system fields.
2. **Fill Data**: Fill in data in Excel according to the agreed format (Note: Relation fields usually require filling in Code or ID).
3. **Upload Import**: Select the file and fill in "Import Reason", then click confirm.
4. The system will validate data legitimacy row by row and trigger corresponding approval or direct save logic.

### 3.2 Batch Update and Delete

- **Batch Update Status**: Check multiple records in the list page to batch modify their "Status" field (such as Freeze, Inactive, etc.).
- **Batch Delete**: After checking data, click "Batch Delete", which also requires filling in a deletion reason. The system supports a "Soft Delete" mechanism, where data records are retained in the database but marked as deleted.

## 4. FAQ

<callout emoji="💡" background-color="light-blue">
Why doesn't data appear in the list page immediately after I click submit?
</callout>

This is usually because the table is configured with an **Approval Flow**. Data currently exists in the "Approval Draft". The system will only synchronize data from the draft table to the official business table after the approval process is finally approved. You can check the approval progress in "My Applications".

<callout emoji="⚠️" background-color="light-yellow">
Why do some dropdowns fail to load content?
</callout>

Please contact the system administrator to check the "Data Source/Relation Config" in "Field Management". Ensure the related table is configured correctly and the current user has permission to access data in that related table.
