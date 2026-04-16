package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuthInfo holds information about the current authentication method.
type AuthInfo struct {
	Accessor   string
	DisplayName string
	Policies   []string
	TTL        time.Duration
	Renewable  bool
	Meta       map[string]string
}

// GetAuthInfo returns information about the token's auth context.
func (c *Client) GetAuthInfo() (*AuthInfo, error) {
	path := fmt.Sprintf("%s/v1/auth/token/lookup-self", c.Address)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: check token validity")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		Data struct {
			Accessor    string            `json:"accessor"`
			DisplayName string            `json:"display_name"`
			Policies    []string          `json:"policies"`
			TTL         int               `json:"ttl"`
			Renewable   bool              `json:"renewable"`
			Meta        map[string]string `json:"meta"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &AuthInfo{
		Accessor:    payload.Data.Accessor,
		DisplayName: payload.Data.DisplayName,
		Policies:    payload.Data.Policies,
		TTL:         time.Duration(payload.Data.TTL) * time.Second,
		Renewable:   payload.Data.Renewable,
		Meta:        payload.Data.Meta,
	}, nil
}
