package audit

import (
	"github.com/yourusername/vaultdiff/internal/diff"
)

// RecordUserpassRoleChange logs a userpass role diff to the audit logger.
// It skips writing if there are no changes.
func RecordUserpassRoleChange(logger *Logger, username string, changes []diff.UserpassChange) error {
	if len(changes) == 0 {
		return nil
	}

	data := map[string]interface{}{
		"username": username,
		"changes":  formatUserpassChanges(changes),
	}

	return logger.Record(Entry{
		Operation: "userpass_role_diff",
		Path:      "auth/userpass/users/" + username,
		Data:      data,
	})
}

func formatUserpassChanges(changes []diff.UserpassChange) []map[string]string {
	out := make([]map[string]string, 0, len(changes))
	for _, c := range changes {
		out = append(out, map[string]string{
			"field":     c.Field,
			"old_value": c.OldValue,
			"new_value": c.NewValue,
		})
	}
	return out
}
