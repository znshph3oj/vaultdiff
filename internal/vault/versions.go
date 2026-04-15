package vault

import (
	"fmt"
	"sort"
	"strconv"
)

// VersionMeta holds metadata about a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns metadata for all versions of a KV v2 secret at the given path.
func (c *Client) ListVersions(path string) ([]VersionMeta, error) {
	secretPath := fmt.Sprintf("%s/metadata/%s", c.mount, path)

	secret, err := c.logical.Read(secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("metadata response missing 'versions' key for path %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for versions data")
	}

	var metas []VersionMeta
	for k, v := range versionsMap {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		meta := VersionMeta{}
		// Use the map key as the version number directly when available.
		if vnum, err := strconv.Atoi(k); err == nil {
			meta.Version = vnum
		}
		if ct, ok := entry["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := entry["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if d, ok := entry["destroyed"].(bool); ok {
			meta.Destroyed = d
		}
		metas = append(metas, meta)
	}

	// Sort by version number ascending.
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	return metas, nil
}
