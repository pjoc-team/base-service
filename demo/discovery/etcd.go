package main

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/discovery"
	"github.com/pjoc-team/etcd-config/etcdv3"
	"time"
)

func main() {
	etcdConfig := &etcdv3.EtcdConfig{}
	etcdConfig.Endpoints = []string{"127.0.0.1:2379"}
	etcdConfig.TimeoutSeconds = 6
	service := discovery.InitEtcdDiscoveryService(etcdConfig, "/pub/pjoc/pay/services")
	if e := service.RegisterService("testService", discovery.ServiceEndpoint{"127.0.0.1", "8888"}); e != nil {
		fmt.Println("Register with error: ", e.Error())
	}
	time.Sleep(2 * time.Second)
	testService := service.GetService("testService")
	fmt.Println("Found service: ", testService)
	time.Sleep(100 * time.Second)
}
