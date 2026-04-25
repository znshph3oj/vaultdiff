package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SAMLRoleInfo holds configuration for a SAML auth role.
type SAMLRoleInfo struct {
	Name            string   `json:"name"`
	BoundAttributes map[string]string `json:"bound_attributes"`
	BoundSubjects   []string `json:"bound_subjects"`
	TokenPolicies   []string `json:"token_policies"`
	TokenTTL        int      `json:"token_ttl"`
	TokenMaxTTL     int      `json:"token_max_ttl"`
}

// GetSAMLRoleInfo retrieves SAML role configuration from Vault.
func GetSAMLRoleInfo(client *Client, mount, role string) (*SAMLRoleInfo, error) {
	path := fmt.Sprintf("/v1/auth/%s/role/%s", mount, role)
	resp, err := client.RawClient().Get(client.Address() + path)
	if err != nil {
		return nil, fmt.Errorf("saml role request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("saml role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for saml role %q", resp.StatusCode, role)
	}

	var payload struct {
		Data SAMLRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode saml role response: %w", err)
	}
	payload.Data.Name = role
	return &payload.Data, nil
}
