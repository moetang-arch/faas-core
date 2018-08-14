package faas_service

import "github.com/moetang-arch/faas-api"

type Service struct {
	ProvidedServices0 map[string]interface{}
}

func NewService() *Service {
	service := new(Service)
	service.ProvidedServices0 = faas.GetServiceMap()
	return service
}
