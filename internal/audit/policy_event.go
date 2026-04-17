package audit

import (
	"time"

	"github.com/user/vaultdiff/internal/diff"
)

func RecordPolicyDiff(logger *Logger, user string, changes []diff.PolicyChange) error {
	for _, c := range changes {
		if c.Status == "unchanged" {
			continue
		}
		entry := Entry{
			Timestamp: time.Now(),
			User:      user,
			Operation: "policy-" + c.Status,
			Path:      "sys/policies/acl/" + c.Name,
			Data: map[string]string{
				"policy": c.Name,
				"status": c.Status,
			},
		}
		if err := logger.Record(entry); err != nil {
			return err
		}
	}
	return nil
}
