package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MountInfo holds configuration details for a single secrets engine mount.
type MountInfo struct {
	Path        string            `json:"path"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Options     map[string]string `json:"options"`
	Local       bool              `json:"local"`
	SealWrap    bool              `json:"seal_wrap"`
}

// ListMounts returns all secret engine mounts from Vault.
func (c *Client) ListMounts() (map[string]*MountInfo, error) {
	req, err := http.NewRequest(http.MethodGet, c.address+"/v1/sys/mounts", nil)
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
		return nil, fmt.Errorf("mounts endpoint not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	mounts := make(map[string]*MountInfo)
	for key, val := range raw {
		var info MountInfo
		if err := json.Unmarshal(val, &info); err != nil {
			continue
		}
		info.Path = key
		mounts[key] = &info
	}
	return mounts, nil
}
