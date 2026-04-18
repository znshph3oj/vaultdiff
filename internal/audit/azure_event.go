package audit

import (
	"github.com/yourusername/vaultdiff/internal/diff"
)

func RecordAzureRoleChange(logger *Logger, path string, changes []diff.AzureRoleChange) error {
	if len(changes) == 0 {
		return nil
	}
	fields := make(map[string]interface{}, len(changes))
	for _, c := range changes {
		fields[c.Field] = map[string]string{"from": c.From, "to": c.To}
	}
	return logger.Record(Entry{
		Operation: "azure_role_diff",
		Path:      path,
		Data:      fields,
	})
}
