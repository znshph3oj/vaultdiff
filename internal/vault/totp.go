package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TOTPKeyInfo holds information about a TOTP key in Vault.
type TOTPKeyInfo struct {
	AccountName string `json:"account_name"`
	Algorithm   string `json:"algorithm"`
	Digits      int    `json:"digits"`
	Issuer      string `json:"issuer"`
	Period      int    `json:"period"`
	QRSize      int    `json:"qr_size"`
}

// GetTOTPKeyInfo retrieves TOTP key configuration from Vault at the given mount and key name.
func GetTOTPKeyInfo(client *Client, mount, keyName string) (*TOTPKeyInfo, error) {
	path := fmt.Sprintf("/v1/%s/keys/%s", mount, keyName)
	resp, err := client.RawClient().NewRequest("GET", path)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	httpResp, err := client.RawClient().RawRequestWithContext(nil, resp)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("totp key %q not found at mount %q", keyName, mount)
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for totp key %q", httpResp.StatusCode, keyName)
	}

	var payload struct {
		Data TOTPKeyInfo `json:"data"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &payload.Data, nil
}
