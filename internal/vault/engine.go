package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EngineInfo holds metadata about a mounted secrets engine.
type EngineInfo struct {
	Path        string
	Type        string
	Description string
	Options     map[string]string
	Local       bool
	SealWrap    bool
}

// ListEngines returns all mounted secrets engines from the Vault server.
func (c *Client) ListEngines(ctx context.Context) ([]EngineInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/v1/sys/mounts", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
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
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var engines []EngineInfo
	for path, data := range raw {
		var m struct {
			Type        string            `json:"type"`
			Description string            `json:"description"`
			Options     map[string]string `json:"options"`
			Local       bool              `json:"local"`
			SealWrap    bool              `json:"seal_wrap"`
		}
		if err := json.Unmarshal(data, &m); err != nil {
			continue
		}
		engines = append(engines, EngineInfo{
			Path:        path,
			Type:        m.Type,
			Description: m.Description,
			Options:     m.Options,
			Local:       m.Local,
			SealWrap:    m.SealWrap,
		})
	}
	return engines, nil
}
