package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GCPKMSKeyInfo holds metadata about a GCP KMS key managed by Vault.
type GCPKMSKeyInfo struct {
	Name          string `json:"name"`
	KeyRing       string `json:"key_ring"`
	CryptoKey     string `json:"crypto_key"`
	Algorithm     string `json:"algorithm"`
	ProtectionLevel string `json:"protection_level"`
	RotationPeriod string `json:"rotation_period"`
	MinVersion    int    `json:"min_version"`
}

// GetGCPKMSKeyInfo retrieves GCP KMS key configuration from Vault.
func GetGCPKMSKeyInfo(client *Client, mount, keyName string) (*GCPKMSKeyInfo, error) {
	path := fmt.Sprintf("/v1/%s/keys/%s", mount, keyName)
	resp, err := client.RawClient().RawRequest(client.RawClient().NewRequest("GET", path))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("GCP KMS key %q not found", keyName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var wrapper struct {
		Data GCPKMSKeyInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return &wrapper.Data, nil
}
