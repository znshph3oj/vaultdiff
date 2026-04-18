package audit

import (
	"fmt"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// RecordKVEngineChange logs a KV engine diff event to the audit logger.
func RecordKVEngineChange(logger *Logger, path string, changes []diff.KVEngineChange) error {
	if len(changes) == 0 {
		return nil
	}
	data := map[string]string{}
	for _, c := range changes {
		data[c.Field] = fmt.Sprintf("%s -> %s", c.OldValue, c.NewValue)
	}
	return logger.Record(Entry{
		Operation: "kv-engine-diff",
		Path:      path,
		Data:      data,
	})
}
