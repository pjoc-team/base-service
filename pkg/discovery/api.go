package discovery

import "sync"

type Service struct {
	ServiceName string            `json:"service_name"`
	Endpoints   []ServiceEndpoint `json:"endpoints"`
	sync.Mutex
}

type ServiceEndpoint struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type DiscoveryService interface {
	RegisterService(serviceName string, endpoint ServiceEndpoint) error

	GetService(serviceName string) *Service
}
