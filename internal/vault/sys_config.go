package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SysConfig holds Vault system configuration settings.
type SysConfig struct {
	DefaultLeaseTTL string `json:"default_lease_ttl"`
	MaxLeaseTTL     string `json:"max_lease_ttl"`
	ForceNoCache    bool   `json:"force_no_cache"`
}

// GetSysConfig fetches the current Vault system configuration.
func (c *Client) GetSysConfig() (*SysConfig, error) {
	url := fmt.Sprintf("%s/v1/sys/config/state/sanitized", c.Address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("sys config not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data SysConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result.Data, nil
}
