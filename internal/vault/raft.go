package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RaftStatus represents the Raft cluster configuration status.
type RaftStatus struct {
	LeaderID        string        `json:"leader_id"`
	AppliedIndex    uint64        `json:"applied_index"`
	CommitIndex     uint64        `json:"commit_index"`
	Servers         []RaftServer  `json:"servers"`
}

// RaftServer represents a single node in the Raft cluster.
type RaftServer struct {
	NodeID   string `json:"node_id"`
	Address  string `json:"address"`
	Leader   bool   `json:"leader"`
	Voter    bool   `json:"voter"`
	Protocol string `json:"protocol_version"`
}

// GetRaftStatus returns the current Raft cluster status from Vault.
func (c *Client) GetRaftStatus() (*RaftStatus, error) {
	req, err := http.NewRequest(http.MethodGet, c.address+"/v1/sys/storage/raft/configuration", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("raft status not found (is raft storage enabled?)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data RaftStatus `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
