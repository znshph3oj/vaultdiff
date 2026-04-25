package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AliasInfo represents a Vault identity alias.
type AliasInfo struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	MountAccessor string            `json:"mount_accessor"`
	MountType     string            `json:"mount_type"`
	CanonicalID   string            `json:"canonical_id"`
	Metadata      map[string]string `json:"metadata"`
}

// GetAliasInfo retrieves an identity alias by ID from Vault.
func GetAliasInfo(client *Client, aliasID string) (*AliasInfo, error) {
	path := fmt.Sprintf("/v1/identity/entity-alias/id/%s", aliasID)
	resp, err := client.RawClient().Get(client.Address() + path)
	if err != nil {
		return nil, fmt.Errorf("alias request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("alias %q not found", aliasID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for alias %q", resp.StatusCode, aliasID)
	}

	var payload struct {
		Data AliasInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode alias response: %w", err)
	}
	return &payload.Data, nil
}
