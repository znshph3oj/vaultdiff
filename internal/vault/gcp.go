package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GCPRoleInfo holds information about a GCP secrets engine role.
type GCPRoleInfo struct {
	Name        string   `json:"name"`
	RoleType    string   `json:"role_type"`
	Project     string   `json:"project"`
	Bindings    string   `json:"bindings"`
	TokenTTL    int      `json:"token_ttl"`
	SecretType  string   `json:"secret_type"`
	Scopes      []string `json:"token_scopes"`
}

// GetGCPRoleInfo fetches information about a GCP role from Vault.
func GetGCPRoleInfo(client *Client, roleName string) (*GCPRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/gcp/roleset/%s", client.Address, roleName)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("gcp role %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data GCPRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &wrapper.Data, nil
}
