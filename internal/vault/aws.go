package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AWSRoleInfo holds configuration for a Vault AWS secret engine role.
type AWSRoleInfo struct {
	Name           string   `json:"name"`
	CredentialType string   `json:"credential_type"`
	PolicyARNs     []string `json:"policy_arns"`
	RoleARNs       []string `json:"role_arns"`
	DefaultTTL     int      `json:"default_ttl"`
	MaxTTL         int      `json:"max_ttl"`
}

// GetAWSRoleInfo retrieves AWS secret engine role configuration from Vault.
func GetAWSRoleInfo(client *Client, role string) (*AWSRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/aws/roles/%s", client.Address, role)

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
		return nil, fmt.Errorf("aws role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data AWSRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	result.Data.Name = role
	return &result.Data, nil
}
