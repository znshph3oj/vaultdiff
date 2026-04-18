package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppRoleInfo holds metadata about a Vault AppRole.
type AppRoleInfo struct {
	RoleID        string   `json:"role_id"`
	BindSecretID  bool     `json:"bind_secret_id"`
	LocalSecretIDs bool    `json:"local_secret_ids"`
	Policies      []string `json:"token_policies"`
	TTL           int      `json:"token_ttl"`
	MaxTTL        int      `json:"token_max_ttl"`
}

// GetAppRoleInfo fetches metadata for the given AppRole name.
func GetAppRoleInfo(client *Client, roleName string) (*AppRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/auth/approle/role/%s", client.Address, roleName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("approle %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data AppRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
