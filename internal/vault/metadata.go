package vault

import (
	"context"
	"fmt"
	"time"
)

// VersionMetadata holds metadata for a single secret version.
type VersionMetadata struct {
	Version      int       `json:"version"`
	CreatedTime  time.Time `json:"created_time"`
	DeletionTime time.Time `json:"deletion_time,omitempty"`
	Destroyed    bool      `json:"destroyed"`
	CreatedBy    string    `json:"created_by,omitempty"`
}

// SecretMetadata holds full metadata for a KV v2 secret path.
type SecretMetadata struct {
	Path            string                     `json:"path"`
	CurrentVersion  int                        `json:"current_version"`
	OldestVersion   int                        `json:"oldest_version"`
	CreatedTime     time.Time                  `json:"created_time"`
	UpdatedTime     time.Time                  `json:"updated_time"`
	Versions        map[int]*VersionMetadata   `json:"versions"`
}

// GetSecretMetadata fetches full metadata for a secret path from Vault KV v2.
func (c *Client) GetSecretMetadata(ctx context.Context, path string) (*SecretMetadata, error) {
	apiPath := fmt.Sprintf("%s/metadata/%s", c.mount, path)

	secret, err := c.logical.ReadWithContext(ctx, apiPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", path)
	}

	meta := &SecretMetadata{
		Path:    path,
		Versions: make(map[int]*VersionMetadata),
	}

	if v, ok := secret.Data["current_version"].(float64); ok {
		meta.CurrentVersion = int(v)
	}
	if v, ok := secret.Data["oldest_version"].(float64); ok {
		meta.OldestVersion = int(v)
	}
	if v, ok := secret.Data["created_time"].(string); ok {
		meta.CreatedTime, _ = time.Parse(time.RFC3339Nano, v)
	}
	if v, ok := secret.Data["updated_time"].(string); ok {
		meta.UpdatedTime, _ = time.Parse(time.RFC3339Nano, v)
	}

	if versions, ok := secret.Data["versions"].(map[string]interface{}); ok {
		for key, raw := range versions {
			var vNum int
			fmt.Sscanf(key, "%d", &vNum)
			if vMap, ok := raw.(map[string]interface{}); ok {
				vm := &VersionMetadata{Version: vNum}
				if ct, ok := vMap["created_time"].(string); ok {
					vm.CreatedTime, _ = time.Parse(time.RFC3339Nano, ct)
				}
				if dt, ok := vMap["deletion_time"].(string); ok && dt != "" {
					vm.DeletionTime, _ = time.Parse(time.RFC3339Nano, dt)
				}
				if d, ok := vMap["destroyed"].(bool); ok {
					vm.Destroyed = d
				}
				meta.Versions[vNum] = vm
			}
		}
	}

	return meta, nil
}
