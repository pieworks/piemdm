package constants

// SystemFieldCodes System Reserved Fields
// These fields are used internally by the system. Users cannot create fields with the same code as these fields.
var SystemFieldCodes = []string{
	// GORM Standard Fields
	"id",
	"created_at",
	"updated_at",
	"deleted_at",
	"created_by",
	"updated_by",

	// Business System Fields
	"status",
	"operation",
	"action",
	"send_status",

	// Workflow Fields (draft table specific)
	"entity_id",
	"approval_code",
	"draft_status",
}

// IsSystemFieldCode Checks if it is a system field code
func IsSystemFieldCode(code string) bool {
	for _, systemCode := range SystemFieldCodes {
		if code == systemCode {
			return true
		}
	}
	return false
}
