package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AuditDevice represents a Vault audit device.
type AuditDevice struct {
	Path        string            `json:"path"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Options     map[string]string `json:"options"`
	Local       bool              `json:"local"`
}

// ListAuditDevices returns all enabled audit devices from Vault.
func (c *Client) ListAuditDevices() (map[string]*AuditDevice, error) {
	req, err := http.NewRequest(http.MethodGet, c.address+"/v1/sys/audit", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("audit devices endpoint not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result map[string]*AuditDevice
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
