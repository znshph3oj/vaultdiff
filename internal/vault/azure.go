package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AzureRoleInfo holds the configuration details for an Azure secrets engine role.
type AzureRoleInfo struct {
	ApplicationObjectID string   `json:"application_object_id"`
	ClientID            string   `json:"client_id"`
	TTL                 string   `json:"ttl"`
	MaxTTL              string   `json:"max_ttl"`
	AzureRoles          []string `json:"azure_roles"`
	AzureGroups         []string `json:"azure_groups"`
}

// GetAzureRoleInfo retrieves the configuration for a named Azure secrets engine
// role from the given Vault client. Returns an error if the role does not exist
// or the request fails.
func GetAzureRoleInfo(client *Client, roleName string) (*AzureRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/azure/roles/%s", client.Address, roleName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for azure role %q: %w", roleName, err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting azure role %q: %w", roleName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("azure role %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for azure role %q", resp.StatusCode, roleName)
	}

	var wrapper struct {
		Data AzureRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding azure role %q response: %w", roleName, err)
	}
	return &wrapper.Data, nil
}
