package service

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/discovery"
	"github.com/pjoc-team/base-service/pkg/util"
	"google.golang.org/grpc"
)

func GetEndpoint(addr string) discovery.ServiceEndpoint {
	ip := util.GetHostIP()
	port := util.GetPortByListenAddr(addr)
	endpoint := discovery.ServiceEndpoint{ip, port}
	return endpoint
}

func (svc *Service) GetGrpcConnectionByDiscoveryService(serviceName string, port string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	service := svc.DiscoveryService.GetService(serviceName)
	if service == nil {
		return svc.Dial(serviceName, port, options...)
	}
	endpoints := service.Endpoints
	if endpoints == nil {
		return svc.Dial(serviceName, port)
	}
	for _, endpoint := range endpoints {
		if conn, e := svc.Dial(endpoint.Ip, endpoint.Port, options...); e != nil {
			continue
		} else {
			return conn, e
		}
	}
	return svc.Dial(serviceName, port, options...)
}

func (svc *Service) Dial(host string, port string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	return grpc.Dial(address, options...)
}
