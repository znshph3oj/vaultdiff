package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// IdentityGroup represents a Vault identity group.
type IdentityGroup struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Policies []string          `json:"policies"`
	Metadata map[string]string `json:"metadata"`
	MemberEntityIDs []string   `json:"member_entity_ids"`
}

// GetIdentityGroup fetches a Vault identity group by ID.
func GetIdentityGroup(client *Client, groupID string) (*IdentityGroup, error) {
	url := fmt.Sprintf("%s/v1/identity/group/id/%s", client.Address, groupID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("group %q not found", groupID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		Data IdentityGroup `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &payload.Data, nil
}
