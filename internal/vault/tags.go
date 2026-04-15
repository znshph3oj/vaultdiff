package vault

import (
	"fmt"
	"net/http"
)

// SecretTags represents the custom metadata (tags) attached to a KV v2 secret.
type SecretTags map[string]string

// GetSecretTags retrieves the custom_metadata field from a KV v2 secret's metadata.
func (c *Client) GetSecretTags(path string) (SecretTags, error) {
	secret, err := c.Logical().Read(fmt.Sprintf("%s/metadata/%s", c.mount, path))
	if err != nil {
		return nil, fmt.Errorf("reading tags for %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q not found", path)
	}

	raw, ok := secret.Data["custom_metadata"]
	if !ok || raw == nil {
		return SecretTags{}, nil
	}

	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for custom_metadata: %T", raw)
	}

	tags := make(SecretTags, len(m))
	for k, v := range m {
		tags[k] = fmt.Sprintf("%v", v)
	}
	return tags, nil
}

// SetSecretTags writes custom_metadata tags to a KV v2 secret's metadata.
func (c *Client) SetSecretTags(path string, tags SecretTags) error {
	data := make(map[string]interface{}, len(tags))
	for k, v := range tags {
		data[k] = v
	}

	_, err := c.Logical().JSONMergePatch(
		nil,
		fmt.Sprintf("%s/metadata/%s", c.mount, path),
		map[string]interface{}{"custom_metadata": data},
	)
	if err != nil {
		// fallback: use Write for Vault versions that don't support PATCH
		_, werr := c.Logical().Write(
			fmt.Sprintf("%s/metadata/%s", c.mount, path),
			map[string]interface{}{"custom_metadata": data},
		)
		if werr != nil {
			return fmt.Errorf("setting tags for %q: %w", path, werr)
		}
	}
	return nil
}

// httpStatus is a helper used in tests to inspect mock responses.
func httpStatus(code int) string {
	return http.StatusText(code)
}
