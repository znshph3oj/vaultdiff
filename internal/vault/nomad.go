package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NomadRoleInfo holds configuration for a Vault Nomad secrets engine role.
type NomadRoleInfo struct {
	Name     string   `json:"name"`
	Policies []string `json:"policies"`
	Global   bool     `json:"global"`
	Type     string   `json:"type"`
	Lease    string   `json:"lease"`
}

// GetNomadRoleInfo retrieves the Nomad role configuration from Vault.
func GetNomadRoleInfo(client *Client, role string) (*NomadRoleInfo, error) {
	path := fmt.Sprintf("/v1/nomad/role/%s", role)
	req, err := http.NewRequest(http.MethodGet, client.Address+path, nil)
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
		return nil, fmt.Errorf("nomad role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for nomad role %q", resp.StatusCode, role)
	}

	var envelope struct {
		Data NomadRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	envelope.Data.Name = role
	return &envelope.Data, nil
}
