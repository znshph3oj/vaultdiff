package vault

import (
	"context"
	"fmt"
	"strings"
)

// PolicyAccess represents the access level for a given path.
type PolicyAccess struct {
	Path         string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

// GetPolicyForPath returns the capabilities a token has for a given secret path.
func (c *Client) GetPolicyForPath(ctx context.Context, path string) ([]PolicyAccess, error) {
	token := c.vault.Token()
	if token == "urn nil, fmt.Errorf("no vault token configured")
	}

	// Use sys/capabilities-self to check what the current token can do.
	body := map[string]interface{}{
		"paths": []string{path},
	}

	secret, err := c.vault.Logical().WriteWithContext(ctx, "sys/capabilities-self", body)
	if err != nil {
		return nil, fmt.Errorf("capabilities lookup failed for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("empty response from capabilities endpoint")
	}

	var accesses []PolicyAccess
	for key, val := range secret.Data {
		if key == "capabilities" {
			continue
		}
		caps, ok := val.([]interface{})
		if !ok {
			continue
		}
		strCaps := make([]string, 0, len(caps))
		for _, c := range caps {
			if s, ok := c.(string); ok {
				strCaps = append(strCaps, s)
			}
		}
		accesses = append(accesses, PolicyAccess{
			Path:         key,
			Capabilities: strCaps,
		})
	}
	return accesses, nil
}

// CanRead returns true if the given capabilities slice includes "read" or "sudo".
func CanRead(caps []string) bool {
	return hasCapability(caps, "read") || hasCapability(caps, "sudo")
}

// CanWrite returns true if the given capabilities slice includes "create", "update", or "sudo".
func CanWrite(caps []string) bool {
	return hasCapability(caps, "create") || hasCapability(caps, "update") || hasCapability(caps, "sudo")
}

func hasCapability(caps []string, target string) bool {
	for _, c := range caps {
		if strings.EqualFold(c, target) {
			return true
		}
	}
	return false
}
