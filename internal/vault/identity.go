package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IdentityEntity struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Policies []string          `json:"policies"`
	Metadata map[string]string `json:"metadata"`
	Disabled bool              `json:"disabled"`
}

func (c *Client) GetIdentityEntity(entityID string) (*IdentityEntity, error) {
	url := fmt.Sprintf("%s/v1/identity/entity/id/%s", c.Address, entityID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("entity %q not found", entityID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data IdentityEntity `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
