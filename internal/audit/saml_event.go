package audit

import (
	"fmt"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// RecordSAMLRoleChange writes an audit log entry for a SAML role diff.
// It skips writing if there are no changes.
func RecordSAMLRoleChange(logger *Logger, user, mount, role string, changes []diff.SAMLRoleChange) error {
	if len(changes) == 0 {
		return nil
	}

	return logger.Record(Entry{
		User:      user,
		Operation: "saml_role_diff",
		Path:      fmt.Sprintf("auth/%s/role/%s", mount, role),
		Data:      formatSAMLChanges(changes),
	})
}

func formatSAMLChanges(changes []diff.SAMLRoleChange) map[string]string {
	out := make(map[string]string, len(changes))
	for _, c := range changes {
		out[c.Field] = fmt.Sprintf("%s -> %s", truncateSAML(c.From, 64), truncateSAML(c.To, 64))
	}
	return out
}

func truncateSAML(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
