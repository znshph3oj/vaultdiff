package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LeaseInfo holds metadata about a Vault secret lease.
type LeaseInfo struct {
	LeaseID       string        `json:"lease_id"`
	Renewable     bool          `json:"renewable"`
	LeaseDuration time.Duration `json:"-"`
	RawDuration   int           `json:"lease_duration"`
	ExpireTime    time.Time     `json:"expire_time,omitempty"`
}

// leaseResponse mirrors the Vault API response for lease lookup.
type leaseResponse struct {
	Data struct {
		ID            string `json:"id"`
		Renewable     bool   `json:"renewable"`
		LeaseDuration int    `json:"ttl"`
		ExpireTime    string `json:"expire_time"`
	} `json:"data"`
}

// GetLeaseInfo retrieves lease information for a given lease ID.
func (c *Client) GetLeaseInfo(leaseID string) (*LeaseInfo, error) {
	path := "/v1/sys/leases/lookup"
	body := map[string]string{"lease_id": leaseID}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal lease request: %w", err)
	}

	req, err := c.newRequest(http.MethodPut, path, data)
	if err != nil {
		return nil, fmt.Errorf("build lease request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lease lookup request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("lease not found: %s", leaseID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for lease lookup", resp.StatusCode)
	}

	var lr leaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return nil, fmt.Errorf("decode lease response: %w", err)
	}

	info := &LeaseInfo{
		LeaseID:       lr.Data.ID,
		Renewable:     lr.Data.Renewable,
		RawDuration:   lr.Data.LeaseDuration,
		LeaseDuration: time.Duration(lr.Data.LeaseDuration) * time.Second,
	}

	if lr.Data.ExpireTime != "" {
		t, err := time.Parse(time.RFC3339, lr.Data.ExpireTime)
		if err == nil {
			info.ExpireTime = t
		}
	}

	return info, nil
}
