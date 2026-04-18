package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PKICertInfo holds metadata about a PKI certificate.
type PKICertInfo struct {
	SerialNumber string
	CommonName   string
	Issuer       string
	NotBefore    string
	NotAfter     string
	Revoked      bool
	Mount        string
}

// GetPKICertInfo fetches certificate metadata from a PKI secrets engine.
func GetPKICertInfo(client *Client, mount, serial string) (*PKICertInfo, error) {
	url := fmt.Sprintf("%s/v1/%s/cert/%s", client.Address, mount, serial)
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
		return nil, fmt.Errorf("certificate %q not found at mount %q", serial, mount)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Certificate  string `json:"certificate"`
			SerialNumber string `json:"serial_number"`
			RevocationTime int64 `json:"revocation_time"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	info := &PKICertInfo{
		SerialNumber: body.Data.SerialNumber,
		Revoked:      body.Data.RevocationTime > 0,
		Mount:        mount,
	}
	return info, nil
}
