package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RabbitMQRoleInfo holds configuration for a RabbitMQ secret engine role.
type RabbitMQRoleInfo struct {
	Name  string            `json:"name"`
	Vhost string            `json:"vhost"`
	Tags  string            `json:"tags"`
	VhostTopics map[string]map[string]string `json:"vhost_topics"`
}

// GetRabbitMQRoleInfo fetches the RabbitMQ role configuration from Vault.
func GetRabbitMQRoleInfo(client *Client, mount, role string) (*RabbitMQRoleInfo, error) {
	path := fmt.Sprintf("/v1/%s/roles/%s", mount, role)
	req, err := http.NewRequest(http.MethodGet, client.Address+path, nil)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("rabbitmq: role %q not found", role)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rabbitmq: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Data RabbitMQRoleInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("rabbitmq: decode response: %w", err)
	}
	payload.Data.Name = role
	return &payload.Data, nil
}
