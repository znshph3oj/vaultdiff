package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SSHRoleInfo struct {
	KeyType        string   `json:"key_type"`
	DefaultUser    string   `json:"default_user"`
	AllowedUsers   string   `json:"allowed_users"`
	TTL            string   `json:"ttl"`
	MaxTTL         string   `json:"max_ttl"`
	AllowedDomains string   `json:"allowed_domains"`
	CIDRList       string   `json:"cidr_list"`
	AllowedExtensions string `json:"allowed_extensions"`
	DefaultExtensions map[string]string `json:"default_extensions"`
	Policies       []string `json:"policies"`
}

func GetSSHRoleInfo(client *Client, roleName string) (*SSHRoleInfo, error) {
	url := fmt.Sprintf("%s/v1/ssh/roles/%s", client.Address, roleName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("ssh role %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data SSHRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
