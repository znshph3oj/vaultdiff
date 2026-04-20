package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LDAPRoleInfo struct {
	RoleName        string   `json:"role_name"`
	GroupFilter     string   `json:"groupfilter"`
	GroupDN         string   `json:"groupdn"`
	GroupAttr       string   `json:"groupattr"`
	UserDN          string   `json:"userdn"`
	UserAttr        string   `json:"userattr"`
	BindDN          string   `json:"binddn"`
	TTL             string   `json:"ttl"`
	MaxTTL          string   `json:"max_ttl"`
	Policies        []string `json:"policies"`
}

func GetLDAPRoleInfo(client *Client, mount, roleName string) (*LDAPRoleInfo, error) {
	path := fmt.Sprintf("/v1/%s/groups/%s", mount, roleName)
	resp, err := client.RawClient().NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("ldap role %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data LDAPRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	wrapper.Data.RoleName = roleName
	return &wrapper.Data, nil
}
