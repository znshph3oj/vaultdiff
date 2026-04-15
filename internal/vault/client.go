package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods for secret versioning.
type Client struct {
	api  *vaultapi.Client
	Mount string
}

// NewClient creates a new Vault client using the provided address and token.
// If addr or token are empty, the SDK falls back to VAULT_ADDR / VAULT_TOKEN env vars.
func NewClient(addr, token, mount string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	if addr != "" {
		cfg.Address = addr
	}

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	if token != "" {
		api.SetToken(token)
	}

	if mount == "" {
		mount = "secret"
	}

	return &Client{api: api, Mount: mount}, nil
}

// SecretVersion holds the data of a specific KV v2 secret version.
type SecretVersion struct {
	Version  int
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

// GetSecretVersion retrieves a specific version of a KV v2 secret.
// Pass version=0 to get the latest version.
func (c *Client) GetSecretVersion(ctx context.Context, path string, version int) (*SecretVersion, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	logical := c.api.Logical()
	secret, err := logical.ReadWithDataWithContext(ctx,
		fmt.Sprintf("%s/data/%s", c.Mount, path),
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q version %d: %w", path, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %q version %d not found", path, version)
	}

	data, _ := secret.Data["data"].(map[string]interface{})
	meta, _ := secret.Data["metadata"].(map[string]interface{})

	v := 0
	if vf, ok := meta["version"].(float64); ok {
		v = int(vf)
	}

	return &SecretVersion{
		Version:  v,
		Data:     data,
		Metadata: meta,
	}, nil
}

// ListSecretVersions returns the metadata for all versions of a KV v2 secret,
// including version numbers, creation times, and deletion status.
func (c *Client) ListSecretVersions(ctx context.Context, path string) (map[string]interface{}, error) {
	logical := c.api.Logical()
	secret, err := logical.ReadWithContext(ctx,
		fmt.Sprintf("%s/metadata/%s", c.Mount, path),
	)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for secret %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for secret %q", path)
	}

	versions, _ := secret.Data["versions"].(map[string]interface{})
	return versions, nil
}
