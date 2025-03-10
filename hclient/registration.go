package hclient

import (
	"github.com/hashicorp/consul/api"
)

func NewAgentClient(
	consulAddr string,
	consulToken string,
) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: consulAddr,
		Token:   consulToken,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
