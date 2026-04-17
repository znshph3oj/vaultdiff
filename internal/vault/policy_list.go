package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PolicyInfo struct {
	Name  string
	Rules string
}

func (c *Client) ListPolicies() ([]string, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/v1/sys/policies/acl?list=true", c.addr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var out struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Data.Keys, nil
}

func (c *Client) GetPolicy(name string) (*PolicyInfo, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/v1/sys/policies/acl/%s", c.addr, name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var out struct {
		Data struct {
			Name  string `json:"name"`
			Rules string `json:"rules"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &PolicyInfo{Name: out.Data.Name, Rules: out.Data.Rules}, nil
}
