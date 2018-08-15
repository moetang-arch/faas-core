package main

import "github.com/moetang-arch/faas-api"

type Service struct {
	providedServices0 map[string]interface{}
}

func NewService() *Service {
	service := new(Service)
	service.providedServices0 = faas.GetServiceMap()
	return service
}

func (this *Service) Serve() error {
	//TODO
	return nil
}
