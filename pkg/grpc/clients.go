package grpc

import (
	"errors"
	"fmt"
	"github.com/pjoc-team/base-service/pkg/logger"
	"github.com/pjoc-team/base-service/pkg/model"
	"github.com/pjoc-team/base-service/pkg/service"
	pb "github.com/pjoc-team/pay-proto/go"
	"strconv"
	"sync"
)

type GrpcClientFactory struct {
	payChannelClientMap map[string]pb.PayChannelClient
	dbServiceClient     pb.PayDatabaseServiceClient
	settlementClient    pb.SettlementGatewayClient

	channelServiceConfigMap *model.ChannelServiceConfigMap
	serviceConfigMap        *model.ServiceConfigMap
	sync.Mutex
	service.Service
}

func InitGrpFactory(svc service.Service, gatewayConfig *service.GatewayConfig) *GrpcClientFactory {
	factory := &GrpcClientFactory{Service: svc}
	factory.serviceConfigMap = gatewayConfig.ServiceMap
	factory.channelServiceConfigMap = gatewayConfig.ChannelServiceMap
	factory.payChannelClientMap = make(map[string]pb.PayChannelClient)
	return factory
}

func (svc *GrpcClientFactory) GetChannelClient(channelId string) (client pb.PayChannelClient, e error) {
	client = svc.payChannelClientMap[channelId]
	if client == nil {
		svc.Lock()
		if client, e = svc.InitChannelServiceClient(channelId); e != nil {
			svc.payChannelClientMap[channelId] = client
			return nil, e
		}
		defer svc.Unlock()
	}
	return client, nil
}

func (svc *GrpcClientFactory) InitChannelServiceClient(channelId string) (pb.PayChannelClient, error) {
	if svc.channelServiceConfigMap == nil {
		e := errors.New("failed to found channel config")
		logger.Log.Errorf(e.Error())
		return nil, e
	}
	configs := *svc.channelServiceConfigMap

	if channelConfig, exists := configs[channelId]; !exists {
		e := fmt.Errorf("could'nt found such channel: %s config", channelId)
		logger.Log.Errorf(e.Error())
		return nil, e
	} else {
		options := svc.TlsDialGrpcOptions()
		if conn, e := svc.GetGrpcConnectionByDiscoveryService(channelConfig.Host, strconv.Itoa(channelConfig.Port), options...); e != nil {
			logger.Log.Errorf("Failed to connect channel service client! error: %s", e.Error())
			return nil, e
		} else {
			client := pb.NewPayChannelClient(conn)
			return client, nil
		}
	}

}

func (svc *GrpcClientFactory) GetDatabaseClient() (pb.PayDatabaseServiceClient, error) {
	if svc.dbServiceClient == nil {
		return svc.InitDatabaseServiceClient()
	} else {
		return svc.dbServiceClient, nil
	}
}

func (svc *GrpcClientFactory) InitDatabaseServiceClient() (pb.PayDatabaseServiceClient, error) {
	if svc.serviceConfigMap == nil {
		e := errors.New("failed to found db config")
		logger.Log.Errorf(e.Error())
		return nil, e
	}
	options := svc.TlsDialGrpcOptions()
	if cfg, exists := (*svc.serviceConfigMap)["pay-database-service"]; !exists {
		e := errors.New("system error!")
		logger.Log.Errorf("could'nt found db service config!")
		return nil, e
	} else if conn, e := svc.GetGrpcConnectionByDiscoveryService(cfg.Host, strconv.Itoa(cfg.Port), options...); e != nil {
		logger.Log.Errorf("Failed to connect database service client! error: %s", e.Error())
		return nil, e
	} else {
		logger.Log.Infof("Found endpoint: %v and get conn: %v", cfg, conn)
		client := pb.NewPayDatabaseServiceClient(conn)
		svc.dbServiceClient = client
		return svc.dbServiceClient, nil
	}
}

func (svc *GrpcClientFactory) GetSettlementClient() (pb.SettlementGatewayClient, error) {
	if svc.settlementClient == nil {
		return svc.InitSettlementClient()
	} else {
		return svc.settlementClient, nil
	}
}

func (svc *GrpcClientFactory) InitSettlementClient() (pb.SettlementGatewayClient, error) {
	if svc.serviceConfigMap == nil {
		e := errors.New("failed to found db config")
		logger.Log.Errorf(e.Error())
		return nil, e
	}
	options := svc.TlsDialGrpcOptions()
	if cfg, exists := (*svc.serviceConfigMap)["settlement-gateway"]; !exists {
		e := errors.New("system error")
		logger.Log.Errorf("could'nt found settlement service config! exists services: %v", svc.serviceConfigMap)
		return nil, e
	} else if conn, e := svc.GetGrpcConnectionByDiscoveryService(cfg.Host, strconv.Itoa(cfg.Port), options...); e != nil {
		logger.Log.Errorf("Failed to connect settlement service client! error: %s", e.Error())
		return nil, e
	} else {
		logger.Log.Infof("Found endpoint: %v and get conn: %v", cfg, conn)
		client := pb.NewSettlementGatewayClient(conn)
		svc.settlementClient = client
		return svc.settlementClient, nil
	}
}
