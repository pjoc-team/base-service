package discovery

import (
	"context"
	"encoding/json"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/pjoc-team/base-service/pkg/logger"
	"github.com/pjoc-team/base-service/pkg/util"
	"github.com/pjoc-team/etcd-config/etcdv3"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const SERVICE_DIR = "/pub/pjoc/pay/services"

type EtcdDiscovery struct {
	client     etcdv3.EtcdClient
	services   map[string]interface{}
	ServiceDir string
}

func (d *EtcdDiscovery) RegisterService(serviceName string, endpoint ServiceEndpoint) error {
	service := d.GetService(serviceName)
	if service == nil {
		service = &Service{}
		service.ServiceName = serviceName
		service.Endpoints = []ServiceEndpoint{endpoint}
	} else {
		service.Endpoints = append(service.Endpoints, endpoint)
	}
	defer d.AddShutdownHook(serviceName, endpoint)
	return d.Update(service)
}

func (d *EtcdDiscovery) AddShutdownHook(serviceName string, endpoint ServiceEndpoint) {
	go func() {
		logger.Log.Warnf("Listening signals...")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		logger.Log.Warnf("Unregister serviceName: %s endpoint: %s", serviceName, endpoint)
		d.Delete(serviceName, endpoint)
		logger.Log.Warnf("Done unregister serviceName: %v and endpoint: %v", serviceName, endpoint)
		os.Exit(0)
	}()
}

func (d *EtcdDiscovery) GetService(serviceName string) *Service {
	path := SERVICE_DIR + "/" + serviceName
	i := d.services[path]
	if i == nil {
		logger.Log.Warnf("Not find service: %s", serviceName)
		return nil
	} else if service, ok := i.(*Service); !ok {
		logger.Log.Warnf("Could'nt convert interface: %v to service... \n", i)
		return nil
	} else {
		return service
	}
}

func (d *EtcdDiscovery) Update(service *Service) error {
	var jsonData []byte
	var e error
	path := SERVICE_DIR + "/" + service.ServiceName
	if jsonData, e = json.Marshal(service); e != nil {
		return e
	}
	jsonStr := string(jsonData)
	timeout, _ := context.WithTimeout(context.TODO(), time.Duration(d.client.EtcdConfig.TimeoutSeconds)*time.Second)
	if _, e := d.client.Client.Put(timeout, path, jsonStr); e != nil {
		return e
	}
	return nil
}

func (d *EtcdDiscovery) Delete(serviceName string, endpoint ServiceEndpoint) error {
	service := d.GetService(serviceName)
	endpoints := removeEndpoint(service.Endpoints, endpoint)
	service.Endpoints = endpoints

	d.Update(service)
	return nil
}

func removeEndpoint(endpoints []ServiceEndpoint, endpoint ServiceEndpoint) []ServiceEndpoint {
	arri := make([]interface{}, len(endpoints))
	for i, e := range endpoints {
		arri[i] = e
	}
	elements := util.RemoveElements(arri, endpoint)
	result := make([]ServiceEndpoint, len(elements))
	for i, e := range elements {
		result[i] = e.(ServiceEndpoint)
	}
	return result
}

func bindServiceFunc() interface{} {
	return &Service{}
}

func InitEtcdDiscoveryService(etcd *etcdv3.EtcdConfig, serviceDir string) DiscoveryService {
	client := etcdv3.Init(etcd)
	if serviceDir == "" {
		serviceDir = SERVICE_DIR
	}
	serviceDir = strings.TrimSuffix(serviceDir, "/")
	result := client.BindWithMultiResult(bindServiceFunc, serviceDir, []clientv3.OpOption{clientv3.WithPrefix()}, etcdv3.JsonBindHandle)

	discovery := &EtcdDiscovery{}
	discovery.client = client
	discovery.services = result
	return discovery
}
