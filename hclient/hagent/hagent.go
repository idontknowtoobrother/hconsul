package hagent

import (
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/idontknowtoobrother/hconsul/hclient"
)

const (
	AllAvaliableId ServiceID = ""
)

type HAgent interface {
	DiscoveServiceAddr(service string, id string) (addr string, err error)
	NewRequest(service Service, id ServiceID) (*Request, error)
}

type Service = string
type ServiceID = string

type hAgent struct {
	client *api.Client
	token  string
}

func NewHAgent(consulAddr, token string) (HAgent, error) {
	client, err := hclient.NewAgentClient(consulAddr)
	if err != nil {
		return nil, err
	}

	hAgent := &hAgent{
		client: client,
		token:  token,
	}

	return hAgent, nil
}

func buidAddr(service *api.CatalogService) string {
	return fmt.Sprintf("http://%s:%d", service.Address, service.ServicePort)
}

func (d *hAgent) DiscoveServiceAddr(service string, id string) (addr string, err error) {
	if service != "" {
		if id != "" {
			return d.discoveServiceWithId(service, id)
		}
		return d.discoveService(service)
	}
	return "", errors.New("not found any service")
}

func (d *hAgent) discoveServiceWithId(service string, id string) (addr string, err error) {
	filter := &Filter{
		Service: service,
		ID:      id,
	}
	filter.Build()

	services, _, err := d.client.Catalog().Service(service, "", &api.QueryOptions{
		Filter: filter.String(),
	})
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("not found any service %s", service)
	}
	return buidAddr(services[0]), nil
}

func (d *hAgent) discoveService(service string) (addr string, err error) {
	services, _, err := d.client.Catalog().Service(service, "", &api.QueryOptions{
		Token: d.token,
	})

	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("not found any service %s", service)
	}

	return buidAddr(services[0]), nil
}
