package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// KubernetesRoleBinding holds the details of a Kubernetes auth role.
type KubernetesRoleBinding struct {
	Name                          string   `json:"name"`
	BoundServiceAccountNames      []string `json:"bound_service_account_names"`
	BoundServiceAccountNamespaces []string `json:"bound_service_account_namespaces"`
	TTL                           string   `json:"ttl"`
	MaxTTL                        string   `json:"max_ttl"`
	Policies                      []string `json:"token_policies"`
}

// GetKubernetesRoleBinding fetches a Kubernetes auth role binding from Vault.
func (c *Client) GetKubernetesRoleBinding(mount, role string) (*KubernetesRoleBinding, error) {
	path := fmt.Sprintf("/v1/auth/%s/role/%s", mount, role)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("kubernetes role %q not found on mount %q", role, mount)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for kubernetes role %q", resp.StatusCode, role)
	}

	var envelope struct {
		Data KubernetesRoleBinding `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	envelope.Data.Name = role
	return &envelope.Data, nil
}
