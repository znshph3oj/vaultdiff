package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// RabbitMQChange represents a single field change between two RabbitMQ role configs.
type RabbitMQChange struct {
	Field string
	From  string
	To    string
}

// CompareRabbitMQRoles returns the list of changes between two RabbitMQ role configs.
func CompareRabbitMQRoles(a, b *vault.RabbitMQRoleInfo) []RabbitMQChange {
	if a == nil || b == nil {
		return nil
	}
	var changes []RabbitMQChange

	if a.Vhost != b.Vhost {
		changes = append(changes, RabbitMQChange{Field: "vhost", From: a.Vhost, To: b.Vhost})
	}
	if a.Tags != b.Tags {
		changes = append(changes, RabbitMQChange{Field: "tags", From: a.Tags, To: b.Tags})
	}

	return changes
}

// PrintRabbitMQDiff writes a diff of two RabbitMQ roles to stdout.
func PrintRabbitMQDiff(a, b *vault.RabbitMQRoleInfo) {
	FprintRabbitMQDiff(os.Stdout, a, b)
}

// FprintRabbitMQDiff writes a diff of two RabbitMQ roles to the given writer.
func FprintRabbitMQDiff(w io.Writer, a, b *vault.RabbitMQRoleInfo) {
	changes := CompareRabbitMQRoles(a, b)
	if len(changes) == 0 {
		fmt.Fprintln(w, "rabbitmq: no changes detected")
		return
	}
	fmt.Fprintln(w, "rabbitmq role diff:")
	for _, c := range changes {
		fmt.Fprintf(w, "  %s: %q -> %q\n", c.Field, c.From, c.To)
	}
}
