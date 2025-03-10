package hagent

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/idontknowtoobrother/hconsul/hclient"
)

const (
	AllAvaliableId ServiceID = ""
)

type HAgent interface {
	GetAddr() string
	GetToken() string
	DiscoveryService(service string) (addr string, err error)
	DiscoveryServiceId(service string, id string) (addr string, err error)
	NewRequest(service Service, id ServiceID) (*Request, error)
}

type Service = string
type ServiceID = string

type hAgent struct {
	client *api.Client
	addr   string
	token  string
}

func NewHAgent(consulAddr, token string) (HAgent, error) {
	client, err := hclient.NewAgentClient(consulAddr, token)
	if err != nil {
		return nil, err
	}

	hAgent := &hAgent{
		client: client,
		addr:   consulAddr,
		token:  token,
	}

	return hAgent, nil
}

func (d *hAgent) GetAddr() string {
	return d.addr
}

func (d *hAgent) GetToken() string {
	return d.token
}

func buidAddr(service *api.CatalogService) string {
	if service.ServicePort == 0 {
		return service.Address
	}
	return fmt.Sprintf("%s:%d", service.Address, service.ServicePort)
}

func (d *hAgent) DiscoveryServiceId(service string, id string) (addr string, err error) {
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

func (d *hAgent) DiscoveryService(service string) (addr string, err error) {
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
