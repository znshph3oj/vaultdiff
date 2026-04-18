package audit

import (
	"time"

	"github.com/your-org/vaultdiff/internal/diff"
)

// RecordDatabaseRoleChange logs a database role diff event to the audit logger.
func RecordDatabaseRoleChange(logger *Logger, path, user string, changes []diff.DatabaseRoleChange) error {
	if len(changes) == 0 {
		return nil
	}
	data := make(map[string]string, len(changes))
	for _, c := range changes {
		data[c.Field] = c.OldValue + " -> " + c.NewValue
	}
	return logger.Record(Entry{
		Timestamp: time.Now().UTC(),
		Path:      path,
		User:      user,
		Operation: "database_role_diff",
		Data:      data,
	})
}
