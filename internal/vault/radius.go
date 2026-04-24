package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RADIUSRoleInfo holds configuration for a RADIUS auth role.
type RADIUSRoleInfo struct {
	Policies []string `json:"policies"`
	TTL      string   `json:"ttl"`
	MaxTTL   string   `json:"max_ttl"`
}

// GetRADIUSRoleInfo fetches RADIUS role info from Vault for the given role name.
func GetRADIUSRoleInfo(client *Client, role string) (*RADIUSRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/auth/radius/users/%s", client.Address, role)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("radius: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("radius: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("radius: role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("radius: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data RADIUSRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("radius: decode response: %w", err)
	}
	return &result.Data, nil
}
