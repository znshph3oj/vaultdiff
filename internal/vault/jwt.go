package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JWTRoleInfo holds configuration for a JWT/OIDC auth role.
type JWTRoleInfo struct {
	Name            string   `json:"name"`
	RoleType        string   `json:"role_type"`
	BoundAudiences  []string `json:"bound_audiences"`
	BoundSubject    string   `json:"bound_subject"`
	UserClaim       string   `json:"user_claim"`
	GroupsClaim     string   `json:"groups_claim"`
	TTL             string   `json:"ttl"`
	MaxTTL          string   `json:"max_ttl"`
	TokenPolicies   []string `json:"token_policies"`
}

// GetJWTRoleInfo retrieves JWT/OIDC role configuration from Vault.
func (c *Client) GetJWTRoleInfo(mount, role string) (*JWTRoleInfo, error) {
	path := fmt.Sprintf("/v1/%s/role/%s", mount, role)
	resp, err := c.http.Get(c.addr + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("jwt role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var envelope struct {
		Data JWTRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	envelope.Data.Name = role
	return &envelope.Data, nil
}
