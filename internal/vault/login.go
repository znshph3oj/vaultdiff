package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LoginInfo holds the result of a Vault login operation.
type LoginInfo struct {
	ClientToken   string
	Accessor      string
	Policies      []string
	LeaseDuration time.Duration
	Renewable     bool
	Meta          map[string]string
}

// Login authenticates against Vault using the given method and credentials.
// method is e.g. "userpass", "approle", "token".
func (c *Client) Login(method string, payload map[string]interface{}) (*LoginInfo, error) {
	path := fmt.Sprintf("/v1/auth/%s/login", method)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("login: marshal payload: %w", err)
	}

	resp, err := c.rawPost(path, body)
	if err != nil {
		return nil, fmt.Errorf("login: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("login: auth method %q not found", method)
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("login: permission denied")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken   string            `json:"client_token"`
			Accessor      string            `json:"accessor"`
			Policies      []string          `json:"policies"`
			LeaseDuration int               `json:"lease_duration"`
			Renewable     bool              `json:"renewable"`
			Meta          map[string]string `json:"metadata"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("login: decode response: %w", err)
	}

	a := result.Auth
	return &LoginInfo{
		ClientToken:   a.ClientToken,
		Accessor:      a.Accessor,
		Policies:      a.Policies,
		LeaseDuration: time.Duration(a.LeaseDuration) * time.Second,
		Renewable:     a.Renewable,
		Meta:          a.Meta,
	}, nil
}
