package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ConsulRoleInfo holds configuration for a Vault Consul secrets engine role.
type ConsulRoleInfo struct {
	Name      string   `json:"name"`
	Policies  []string `json:"policies"`
	TokenType string   `json:"token_type"`
	TTL       string   `json:"ttl"`
	MaxTTL    string   `json:"max_ttl"`
	Local     bool     `json:"local"`
}

// GetConsulRoleInfo retrieves Consul role configuration from Vault.
func GetConsulRoleInfo(client *Client, mount, role string) (*ConsulRoleInfo, error) {
	path := fmt.Sprintf("/v1/%s/roles/%s", mount, role)
	resp, err := client.RawClient().NewRequest("GET", path)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	httpResp, err := client.RawClient().RawRequestWithContext(client.Context(), resp)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("consul role %q not found at mount %q", role, mount)
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for consul role %q", httpResp.StatusCode, role)
	}

	var wrapper struct {
		Data ConsulRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	wrapper.Data.Name = role
	return &wrapper.Data, nil
}
