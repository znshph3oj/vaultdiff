package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TokenInfo holds metadata about a Vault token.
type TokenInfo struct {
	Accessor   string    `json:"accessor"`
	Policies   []string  `json:"policies"`
	TTL        int       `json:"ttl"`
	ExpireTime time.Time `json:"expire_time"`
	Renewable  bool      `json:"renewable"`
	Meta       map[string]string `json:"meta"`
}

// LookupToken retrieves information about the token currently used by the client.
func (c *Client) LookupToken() (*TokenInfo, error) {
	url := fmt.Sprintf("%s/v1/auth/token/lookup-self", c.Address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("token lookup unauthorized (status %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			Accessor   string            `json:"accessor"`
			Policies   []string          `json:"policies"`
			TTL        int               `json:"ttl"`
			ExpireTime string            `json:"expire_time"`
			Renewable  bool              `json:"renewable"`
			Meta       map[string]string `json:"meta"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	info := &TokenInfo{
		Accessor:  envelope.Data.Accessor,
		Policies:  envelope.Data.Policies,
		TTL:       envelope.Data.TTL,
		Renewable: envelope.Data.Renewable,
		Meta:      envelope.Data.Meta,
	}
	if envelope.Data.ExpireTime != "" {
		t, err := time.Parse(time.RFC3339, envelope.Data.ExpireTime)
		if err == nil {
			info.ExpireTime = t
		}
	}
	return info, nil
}
