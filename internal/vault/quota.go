package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// QuotaInfo holds rate limit quota details for a Vault path.
type QuotaInfo struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Type        string  `json:"type"`
	Rate        float64 `json:"rate"`
	Interval    float64 `json:"interval"`
	BlockInterval float64 `json:"block_interval"`
}

// GetQuota retrieves the rate limit quota for the given quota name.
func (c *Client) GetQuota(name string) (*QuotaInfo, error) {
	path := fmt.Sprintf("/v1/sys/quotas/rate-limit/%s", name)
	req, err := http.NewRequest(http.MethodGet, c.address+path, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("quota %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data QuotaInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &wrapper.Data, nil
}
