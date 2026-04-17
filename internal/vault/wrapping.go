package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WrappingInfo holds metadata about a wrapped token response.
type WrappingInfo struct {
	Token          string `json:"token"`
	Accessor       string `json:"accessor"`
	TTL            int    `json:"ttl"`
	CreationTime   string `json:"creation_time"`
	CreationPath   string `json:"creation_path"`
	WrappedAccessor string `json:"wrapped_accessor"`
}

// LookupWrappingToken queries Vault for metadata about a wrapping token.
func (c *Client) LookupWrappingToken(token string) (*WrappingInfo, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/wrapping/lookup", nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("X-Vault-Wrap-TTL-Token", token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("wrapping token not found or expired")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data WrappingInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result.Data, nil
}
