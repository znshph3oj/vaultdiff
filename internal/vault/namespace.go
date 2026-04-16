package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NamespaceInfo holds metadata about a Vault namespace.
type NamespaceInfo struct {
	Path        string            `json:"path"`
	ID          string            `json:"id"`
	CustomMeta  map[string]string `json:"custom_metadata"`
}

// ListNamespaces returns the child namespaces under the given prefix.
func (c *Client) ListNamespaces(prefix string) ([]NamespaceInfo, error) {
	path := fmt.Sprintf("/v1/sys/namespaces/%s", prefix)
	req, err := http.NewRequest(http.MethodGet, c.address+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("namespace not found: %s", prefix)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	namespaces := make([]NamespaceInfo, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		namespaces = append(namespaces, NamespaceInfo{Path: prefix + k})
	}
	return namespaces, nil
}
