package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// KubernetesRoleInfo holds configuration for a Vault Kubernetes auth role.
type KubernetesRoleInfo struct {
	Name                          string   `json:"name"`
	BoundServiceAccountNames      []string `json:"bound_service_account_names"`
	BoundServiceAccountNamespaces []string `json:"bound_service_account_namespaces"`
	TTL                           string   `json:"ttl"`
	MaxTTL                        string   `json:"max_ttl"`
	Policies                      []string `json:"policies"`
}

// GetKubernetesRoleInfo retrieves a Kubernetes auth role from Vault.
func GetKubernetesRoleInfo(client *Client, roleName string) (*KubernetesRoleInfo, error) {
	path := fmt.Sprintf("/v1/auth/kubernetes/role/%s", roleName)
	req, err := http.NewRequest(http.MethodGet, client.Address+path, nil)
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
		return nil, fmt.Errorf("kubernetes role %q not found", roleName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data KubernetesRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	result.Data.Name = roleName
	return &result.Data, nil
}
