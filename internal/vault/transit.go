package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TransitKeyInfo holds metadata about a Vault transit key.
type TransitKeyInfo struct {
	Name            string
	Type            string
	DeletionAllowed bool
	Exportable      bool
	LatestVersion   int
	MinDecryptVersion int
}

// GetTransitKeyInfo fetches metadata for a named transit key.
func GetTransitKeyInfo(client *Client, keyName string) (*TransitKeyInfo, error) {
	url := fmt.Sprintf("%s/v1/transit/keys/%s", client.Address, keyName)

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
		return nil, fmt.Errorf("transit key %q not found", keyName)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Type              string `json:"type"`
			DeletionAllowed   bool   `json:"deletion_allowed"`
			Exportable        bool   `json:"exportable"`
			LatestVersion     int    `json:"latest_version"`
			MinDecryptVersion int    `json:"min_decryption_version"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &TransitKeyInfo{
		Name:              keyName,
		Type:              body.Data.Type,
		DeletionAllowed:   body.Data.DeletionAllowed,
		Exportable:        body.Data.Exportable,
		LatestVersion:     body.Data.LatestVersion,
		MinDecryptVersion: body.Data.MinDecryptVersion,
	}, nil
}
