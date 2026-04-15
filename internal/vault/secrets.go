package vault

import (
	"context"
	"fmt"
	"strconv"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretVersion holds the data and metadata for a specific secret version.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Metadata map[string]interface{}
}

// GetSecretVersion retrieves a specific version of a KV v2 secret.
// If version is 0, the latest version is returned.
func (c *Client) GetSecretVersion(ctx context.Context, path string, version int) (*SecretVersion, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{strconv.Itoa(version)}
	}

	secret, err := c.logical.ReadWithDataWithContext(ctx,
		fmt.Sprintf("%s/data/%s", c.mount, path),
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q version %d: %w", path, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %q version %d not found", path, version)
	}

	rawData, ok := secret.Data["data"]
	if !ok || rawData == nil {
		return nil, fmt.Errorf("secret %q has no data field", path)
	}

	rawMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("secret %q data field has unexpected type", path)
	}

	data := make(map[string]string, len(rawMap))
	for k, v := range rawMap {
		data[k] = fmt.Sprintf("%v", v)
	}

	meta, _ := secret.Data["metadata"].(map[string]interface{})

	var resolvedVersion int
	if meta != nil {
		if v, ok := meta["version"]; ok {
			switch val := v.(type) {
			case json.Number:
				n, _ := val.Int64()
				resolvedVersion = int(n)
			case float64:
				resolvedVersion = int(val)
			}
		}
	}

	return &SecretVersion{
		Version:  resolvedVersion,
		Data:     data,
		Metadata: meta,
	}, nil
}

// logical is an interface subset of *vaultapi.Logical used for testing.
type logicalClient interface {
	ReadWithDataWithContext(ctx context.Context, path string, data map[string][]string) (*vaultapi.Secret, error)
}
