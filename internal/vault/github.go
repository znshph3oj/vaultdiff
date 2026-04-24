package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GitHubRoleInfo holds configuration for a Vault GitHub auth role (team mapping).
type GitHubRoleInfo struct {
	TeamName string   `json:"team_name"`
	Policies []string `json:"policies"`
	TTL      string   `json:"ttl"`
	MaxTTL   string   `json:"max_ttl"`
}

// GetGitHubRoleInfo retrieves the GitHub auth role configuration for a given team.
func GetGitHubRoleInfo(client *Client, mount, team string) (*GitHubRoleInfo, error) {
	path := fmt.Sprintf("/v1/auth/%s/map/teams/%s", mount, team)
	req, err := http.NewRequest(http.MethodGet, client.Address+path, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("github role %q not found at mount %q", team, mount)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var wrapper struct {
		Data GitHubRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &wrapper.Data, nil
}
