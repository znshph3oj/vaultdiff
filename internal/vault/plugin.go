package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PluginInfo holds metadata about a registered Vault plugin.
type PluginInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Builtin bool   `json:"builtin"`
}

// ListPlugins returns all registered plugins of the given type ("auth", "secret", "database").
// Pass an empty string to list all types.
func (c *Client) ListPlugins(pluginType string) ([]PluginInfo, error) {
	path := "/v1/sys/plugins/catalog"
	if pluginType != "" {
		path = fmt.Sprintf("/v1/sys/plugins/catalog/%s", pluginType)
	}

	resp, err := c.http.Get(c.addr + path)
	if err != nil {
		return nil, fmt.Errorf("list plugins: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin catalog not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list plugins: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Detailed []PluginInfo `json:"detailed"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode plugins: %w", err)
	}
	return body.Data.Detailed, nil
}
