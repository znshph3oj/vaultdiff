package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// KVEngineInfo holds metadata about a KV engine version.
type KVEngineInfo struct {
	Path    string
	Version int
	Options map[string]string
}

// GetKVEngineInfo returns KV engine info (v1 or v2) for the given mount path.
func GetKVEngineInfo(client *Client, mountPath string) (*KVEngineInfo, error) {
	url := fmt.Sprintf("%s/v1/sys/mounts/%s/tune", client.Address, mountPath)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("mount not found: %s", mountPath)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Options map[string]string `json:"options"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	version := 1
	if v, ok := body.Options["version"]; ok && v == "2" {
		version = 2
	}

	return &KVEngineInfo{
		Path:    mountPath,
		Version: version,
		Options: body.Options,
	}, nil
}
