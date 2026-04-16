package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SealStatus represents the seal state of a Vault instance.
type SealStatus struct {
	Sealed      bool   `json:"sealed"`
	Initialized bool   `json:"initialized"`
	T           int    `json:"t"`
	N           int    `json:"n"`
	Progress    int    `json:"progress"`
	Version     string `json:"version"`
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_id"`
}

// GetSealStatus fetches the current seal status from Vault.
func (c *Client) GetSealStatus() (*SealStatus, error) {
	req, err := http.NewRequest(http.MethodGet, c.address+"/v1/sys/seal-status", nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting seal status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("seal status endpoint not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var status SealStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &status, nil
}
