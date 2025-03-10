package hservice

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/idontknowtoobrother/hconsul/util"
)

type Service struct {
	ServiceRegistration
	checkId      string
	consulClient *api.Client
	stopCh       chan os.Signal
}

type ServiceRegistration struct {
	Kind           api.ServiceKind
	ID             string
	Service        string
	Port           int
	ProxyEnvoyPort int
	Meta           map[string]string
	TTL            time.Duration
}

func NewService(
	consulAddr string,
	registration ServiceRegistration,
) (*Service, error) {
	client, err := api.NewClient(&api.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	checkId := util.NewCheckID(registration.ID)

	if registration.Kind == "" {
		registration.Kind = api.ServiceKindTypical
	}

	svc := &Service{
		ServiceRegistration: registration,
		checkId:             checkId,
		consulClient:        client,
		stopCh:              make(chan os.Signal, 1),
	}
	signal.Notify(svc.stopCh, os.Interrupt, syscall.SIGTERM)
	return svc, nil
}

func (s *Service) Register() error {
	ttl := s.TTL
	timeout := ttl * 2
	deregisterAfter := timeout * 2

	registration := &api.AgentServiceRegistration{
		Kind: s.Kind,
		ID:   s.ID,
		Name: s.Service,
		Port: s.Port,
		Meta: map[string]string{
			"id": s.ID,
		},
		Tags: []string{
			s.ID,
		},
		Connect: &api.AgentServiceConnect{
			Native: true,
		},
		Check: &api.AgentServiceCheck{
			Name:                           "Agent Healty Kub 😃",
			TTL:                            ttl.String(),
			Timeout:                        timeout.String(),
			Notes:                          fmt.Sprintf("TTL: %s, Timeout: %s, Will be deregistered after: %s", ttl, timeout, deregisterAfter),
			Status:                         api.HealthPassing,
			DeregisterCriticalServiceAfter: deregisterAfter.String(),
		},
	}

	if s.Meta != nil {
		for k, v := range s.Meta {
			if k == "id" {
				continue
			}
			registration.Meta[k] = v
		}
	}

	if err := s.consulClient.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	s.startHeartbeat()
	return nil
}

func (s *Service) Deregister() error {
	return s.consulClient.Agent().ServiceDeregister(s.checkId)
}

func (s *Service) startHeartbeat() {
	interval := s.TTL / 2
	ticker := time.NewTicker(interval)

	if err := s.consulClient.Agent().UpdateTTL(s.checkId, "healthy", api.HealthPassing); err != nil {
		fmt.Printf("failed to update initial TTL: %v\n", err)
	}

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := s.consulClient.Agent().UpdateTTL(s.checkId, "passing", api.HealthPassing); err != nil {
					fmt.Printf("Failed to update TTL: %v\n", err)
				}
			case <-s.stopCh:
				fmt.Println("stopped heartbeat")
				if err := s.Deregister(); err != nil {
					fmt.Println(err)
				}
				return
			}
		}
	}()
}
