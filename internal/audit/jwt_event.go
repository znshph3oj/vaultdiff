package audit

import (
	"fmt"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// RecordJWTRoleChange writes an audit log entry for JWT role changes.
func RecordJWTRoleChange(logger *Logger, user, mount, role string, changes []diff.JWTRoleChange) error {
	if len(changes) == 0 {
		return nil
	}
	return logger.Record(Entry{
		User:      user,
		Operation: "jwt_role_diff",
		Path:      fmt.Sprintf("%s/role/%s", mount, role),
		Data:      formatJWTChanges(changes),
	})
}

func formatJWTChanges(changes []diff.JWTRoleChange) map[string]string {
	result := make(map[string]string, len(changes))
	for _, c := range changes {
		result[c.Field] = fmt.Sprintf("%s -> %s",
			truncate(c.From, 64),
			truncate(c.To, 64))
	}
	return result
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
