# Entity Approval

Approval management is a core link in PieMDM data quality assurance. Through configurable approval flows, the system ensures that changes to critical data are reviewed by relevant responsible persons.

## 1. Approval Center Overview

Users can enter the Approval Center through "Approval Management" in the navigation bar. This integrates all approval tasks related to the current user, mainly divided into the following three dimensions:

- **Pending My Approval**: Approval tasks that currently need to be processed by you.
- **I Have Processed**: Historical approval records that you have finished processing.
- **My Applications**: All data change applications initiated by you and their current progress.

## 2. Processing Approval Tasks

When there is a new application that needs your approval, you can follow these steps:

1. **View List**: In the "Pending My Approval" tab, the system displays information such as application title, initiator, application time, and current node.
2. **View Details**: Click "Details" to enter the task processing page.
3. **Check Form**: The page will display detailed data of this change. If it is a modify operation, the system usually highlights the changed fields.
4. **History Trace**: At the bottom of the page, you can view the processing trajectory and opinions prior to this application.
5. **Execute Operation**:
   - **Approve**: Data flows to the next node or takes final effect.
   - **Reject**: The application is returned to the initiator, who needs to modify and resubmit.
   - **Fill Opinion**: Regardless of the operation performed, it is recommended to fill in specific processing opinions.

## 3. Initiator's Follow-up Operations

As an applicant, you can track progress in real-time in "My Applications":

- **Track Status**: View which approval node the application is currently staying at, and who is the current approver.
- **Withdraw Application**: Before the process ends, you can withdraw a submitted application (depending on permission configuration).
- **Process End**:
  - **Approved**: Data will be automatically migrated from the `_draft` table to the official table, and the version record will be updated.
  - **Rejected**: You can click "Resubmit" to modify based on the original application data and re-initiate approval.

## 4. Approval and Data Effect Logic

<callout emoji="⚙️" background-color="pale-gray">
The system backend automatically associates business tables with approval definitions. When you click submit in "Entity Management", the backend detects if the table is configured with mandatory approval. If so, data will be inserted into the <code>_draft</code> suffix table of that table. The system performs the <code>MOVE TARGET</code> operation to overwrite or merge the draft into official data only after the process is marked as <code>Approved</code>.
</callout>

## 5. FAQ

<callout emoji="❓" background-color="light-purple">
If I find an error in the data during approval, can I modify it directly?
</callout>

Approvers usually only have view and decision rights and cannot directly modify other people's application data. If you think the data is incorrect, please "Reject" it with modification suggestions attached, and let the initiator correct it themselves.

<callout emoji="❓" background-color="light-purple">
How to view historical approval records?
</callout>

You can check "I Have Processed" or view its "Change Log (LOG)" on specific entity detail pages. The log records detailed data versions after each approval.
