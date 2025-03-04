package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/idontknowtoobrother/hconsul/service"
)

func main() {
	mailService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindTypical,
			Service: "mail-service",
			ID:      "mail-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10000,
			ProxyEnvoyPort: 8181,
			TTL:            1 * time.Minute,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = mailService.Register(); err != nil {
		log.Fatal(err)
	}

	factoryService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindTypical,
			Service: "factory-service",
			ID:      "factory-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10001,
			ProxyEnvoyPort: 8181,
			TTL:            1 * time.Minute,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = factoryService.Register(); err != nil {
		log.Fatal(err)
	}

	paymentService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindTypical,
			Service: "payment-service",
			ID:      "payment-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10002,
			ProxyEnvoyPort: 8181,
			TTL:            1 * time.Minute,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	if err = paymentService.Register(); err != nil {
		log.Fatal(err)
	}

	receiptionService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindTypical,
			Service: "receiption-service",
			ID:      "receiption-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10003,
			ProxyEnvoyPort: 8181,
			TTL:            1 * time.Minute,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	if err = receiptionService.Register(); err != nil {
		log.Fatal(err)
	}

	hrmService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindTypical,
			Service: "hrm-service",
			ID:      "hrm-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10005,
			ProxyEnvoyPort: 8182,
			TTL:            1 * time.Minute,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = hrmService.Register(); err != nil {
		log.Fatal(err)
	}

	privateGatewayService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindAPIGateway,
			Service: "private-gateway",
			ID:      "gateway-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10004,
			ProxyEnvoyPort: 8181,
			TTL:            1 * time.Minute,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = privateGatewayService.Register(); err != nil {
		log.Fatal(err)
	}

	saasGatewayService, err := service.NewService(
		"localhost:8500",
		service.ServiceRegistration{
			Kind:    api.ServiceKindAPIGateway,
			Service: "saas-gateway",
			ID:      "saas-1",
			Meta: map[string]string{
				"api_version": "v1",
			},
			Port:           10006,
			ProxyEnvoyPort: 8182,
			TTL:            1 * time.Minute,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = saasGatewayService.Register(); err != nil {
		log.Fatal(err)
	}

	go func() {
		// health check
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		http.ListenAndServe(":8080", nil)
	}()

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")

	if err = mailService.Deregister(); err != nil {
		log.Fatal(err)
	}

	if err = factoryService.Deregister(); err != nil {
		log.Fatal(err)
	}

	if err = paymentService.Deregister(); err != nil {
		log.Fatal(err)
	}

	if err = receiptionService.Deregister(); err != nil {
		log.Fatal(err)
	}

	if err = privateGatewayService.Deregister(); err != nil {
		log.Fatal(err)
	}

	if err = hrmService.Deregister(); err != nil {
		log.Fatal(err)
	}

	if err = saasGatewayService.Deregister(); err != nil {
		log.Fatal(err)
	}

	log.Println("Service stopped")
}
