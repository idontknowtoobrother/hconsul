package client

import (
	"github.com/hashicorp/consul/api"
)

func NewAgentClient(consulAddr string) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
