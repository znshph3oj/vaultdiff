package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AzureRoleInfo struct {
	ApplicationObjectID string   `json:"application_object_id"`
	ClientID            string   `json:"client_id"`
	TTL                 string   `json:"ttl"`
	MaxTTL              string   `json:"max_ttl"`
	AzureRoles          []string `json:"azure_roles"`
	AzureGroups         []string `json:"azure_groups"`
}

func GetAzureRoleInfo(client *Client, roleName string) (*AzureRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/azure/roles/%s", client.Address, roleName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("azure role %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data AzureRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
