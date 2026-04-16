package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ReplicationStatus holds the DR and performance replication state.
type ReplicationStatus struct {
	DRMode          string `json:"dr_mode"`
	PerformanceMode string `json:"performance_mode"`
	Primary         bool   `json:"primary"`
	KnownSecondaries []string `json:"known_secondaries"`
}

type replicationResponse struct {
	Data struct {
		DR struct {
			Mode    string   `json:"mode"`
			Primary bool     `json:"primary"`
			KnownSecondaries []string `json:"known_secondaries"`
		} `json:"dr"`
		Performance struct {
			Mode string `json:"mode"`
		} `json:"performance"`
	} `json:"data"`
}

// GetReplicationStatus returns the current replication status from Vault.
func GetReplicationStatus(client *Client) (*ReplicationStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/replication/status", client.Address)
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
		return nil, fmt.Errorf("replication status endpoint not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var r replicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &ReplicationStatus{
		DRMode:           r.Data.DR.Mode,
		PerformanceMode:  r.Data.Performance.Mode,
		Primary:          r.Data.DR.Primary,
		KnownSecondaries: r.Data.DR.KnownSecondaries,
	}, nil
}
