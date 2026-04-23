package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UserpassRoleInfo holds configuration for a userpass auth role.
type UserpassRoleInfo struct {
	Username      string   `json:"username"`
	Policies      []string `json:"token_policies"`
	TTL           string   `json:"token_ttl"`
	MaxTTL        string   `json:"token_max_ttl"`
	BoundCIDRs    []string `json:"token_bound_cidrs"`
}

// GetUserpassRoleInfo fetches userpass user configuration from Vault.
func GetUserpassRoleInfo(client *Client, username string) (*UserpassRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/auth/userpass/users/%s", client.Address, username)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("userpass: create request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userpass: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("userpass: user %q not found", username)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userpass: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Data UserpassRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("userpass: decode response: %w", err)
	}
	payload.Data.Username = username
	return &payload.Data, nil
}
