package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DatabaseRoleInfo holds information about a Vault database role.
type DatabaseRoleInfo struct {
	Name               string
	DBName             string
	CreationStatements []string
	DefaultTTL         string
	MaxTTL             string
}

// GetDatabaseRoleInfo fetches database role info from Vault.
func GetDatabaseRoleInfo(client *Client, role string) (*DatabaseRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/database/roles/%s", client.Address, role)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("database role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			DBName             string   `json:"db_name"`
			CreationStatements []string `json:"creation_statements"`
			DefaultTTL         string   `json:"default_ttl"`
			MaxTTL             string   `json:"max_ttl"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &DatabaseRoleInfo{
		Name:               role,
		DBName:             body.Data.DBName,
		CreationStatements: body.Data.CreationStatements,
		DefaultTTL:         body.Data.DefaultTTL,
		MaxTTL:             body.Data.MaxTTL,
	}, nil
}
