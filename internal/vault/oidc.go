package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OIDCRoleInfo holds configuration for an OIDC auth role.
type OIDCRoleInfo struct {
	RoleName        string   `json:"role_name"`
	BoundAudiences  []string `json:"bound_audiences"`
	AllowedRedirects []string `json:"allowed_redirect_uris"`
	UserClaim       string   `json:"user_claim"`
	TokenTTL        int      `json:"token_ttl"`
	TokenMaxTTL     int      `json:"token_max_ttl"`
	TokenPolicies   []string `json:"token_policies"`
}

// GetOIDCRoleInfo retrieves OIDC role configuration from Vault.
func GetOIDCRoleInfo(client *Client, mount, role string) (*OIDCRoleInfo, error) {
	if role == "" {
		return nil, fmt.Errorf("role name must not be empty")
	}
	path := fmt.Sprintf("/v1/auth/%s/role/%s", mount, role)
	resp, err := client.RawClient().Get(client.Address() + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return nil, fmt.Errorf("oidc role %q not found", role)
	default:
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data OIDCRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	wrapper.Data.RoleName = role
	return &wrapper.Data, nil
}
