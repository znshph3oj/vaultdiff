package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TerraformRoleInfo holds configuration for a Vault Terraform Cloud secret engine role.
type TerraformRoleInfo struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	TeamID       string `json:"team_id"`
	TTL          string `json:"ttl"`
	MaxTTL       string `json:"max_ttl"`
	TokenType    string `json:"token_type"`
}

// GetTerraformRoleInfo retrieves a Terraform Cloud secret engine role from Vault.
func GetTerraformRoleInfo(client *Client, mount, role string) (*TerraformRoleInfo, error) {
	if role == "" {
		return nil, fmt.Errorf("role name must not be empty")
	}
	path := fmt.Sprintf("%s/role/%s", mount, role)
	resp, err := client.RawClient().RawRequest(
		client.RawClient().NewRequest("GET", "/v1/"+path),
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// continue
	case http.StatusNotFound:
		return nil, fmt.Errorf("terraform role %q not found", role)
	default:
		return nil, fmt.Errorf("unexpected status %d for terraform role %q", resp.StatusCode, role)
	}

	var wrapper struct {
		Data TerraformRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	wrapper.Data.Name = role
	return &wrapper.Data, nil
}
