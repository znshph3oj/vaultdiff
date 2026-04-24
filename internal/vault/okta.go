package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OktaRoleInfo holds configuration for an Okta auth method role.
type OktaRoleInfo struct {
	Policies        []string `json:"policies"`
	TTL             string   `json:"ttl"`
	MaxTTL          string   `json:"max_ttl"`
	BoundGroups     []string `json:"bound_groups"`
	BoundUsers      []string `json:"bound_users"`
}

// GetOktaRoleInfo retrieves Okta role configuration from Vault.
func GetOktaRoleInfo(client *Client, mount, roleName string) (*OktaRoleInfo, error) {
	path := fmt.Sprintf("%s/v1/auth/%s/groups/%s", client.Address, mount, roleName)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("okta role %q not found at mount %q", roleName, mount)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for okta role %q", resp.StatusCode, roleName)
	}

	var result struct {
		Data OktaRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result.Data, nil
}
