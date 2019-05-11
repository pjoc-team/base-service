package service

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/model"
	"github.com/pjoc-team/base-service/pkg/util"
	"github.com/pjoc-team/etcd-config/config"
	"github.com/pjoc-team/etcd-config/config/etcd"
	"strings"
)

type GatewayConfig struct {
	// 集群配置
	PayConfig *model.PayConfig
	// 通知配置
	NoticeConfig *model.NoticeConfig
	// AppId和费率配置
	AppIdAndChannelConfigMap *model.AppIdAndChannelConfigMap
	// AppId和商户配置
	AppIdAndMerchantMap *model.MerchantConfigMap
	// 服务和对应的部署服务名映射
	ServiceMap *model.ServiceConfigMap
	// Channel和对应host配置
	ChannelServiceMap *model.ChannelServiceConfigMap
}

func InitGatewayConfig(etcdPeers string, dirRoot string) *GatewayConfig {
	payConfig := &model.PayConfig{}
	defaultClusterId := "01"
	payConfig.ClusterId = &defaultClusterId
	payConfig.Concurrency = 1000000

	noticeConfig := &model.NoticeConfig{}
	noticeConfig.NoticeIntervalSecond = 1
	noticeConfig.NoticeDelaySecondExpressions = []int{30, 30, 120, 240, 480, 1200, 3600, 7200, 43200, 86400, 172800}

	appIdRateMap := &model.AppIdAndChannelConfigMap{}
	merchantConfigMap := &model.MerchantConfigMap{}
	serviceMap := &model.ServiceConfigMap{}
	channelServiceMap := &model.ChannelServiceConfigMap{}

	InitConfig(etcdPeers, dirRoot, payConfig)
	InitConfig(etcdPeers, dirRoot, appIdRateMap)
	InitConfig(etcdPeers, dirRoot, merchantConfigMap)
	InitConfig(etcdPeers, dirRoot, serviceMap)
	InitConfig(etcdPeers, dirRoot, channelServiceMap)
	InitConfig(etcdPeers, dirRoot, noticeConfig)

	gatewayConfig := &GatewayConfig{}
	gatewayConfig.PayConfig = payConfig
	gatewayConfig.AppIdAndChannelConfigMap = appIdRateMap
	gatewayConfig.AppIdAndMerchantMap = merchantConfigMap
	gatewayConfig.ServiceMap = serviceMap
	gatewayConfig.ChannelServiceMap = channelServiceMap
	gatewayConfig.NoticeConfig = noticeConfig
	return gatewayConfig
}

func BuildPath(dirRoot string, c model.Config) string {
	path := ""
	if c.ConfigType() == model.MAP {
		path = util.AssembleDir(dirRoot, c.DirName())
		path = util.AssembleDir(path, "/")
	} else {
		path = strings.TrimSuffix(util.AssembleDir(dirRoot, c.DirName()), "/")
	}
	return path
}

func InitConfig(etcdPeers string, dirRoot string, c model.Config) {
	path := ""
	if c.ConfigType() == model.MAP {
		path = util.AssembleDir(dirRoot, c.DirName())
		path = util.AssembleDir(path, "/")
	} else {
		path = strings.TrimSuffix(util.AssembleDir(dirRoot, c.DirName()), "/")
	}
	InitConfigByEtcd(etcdPeers, c, path)
}

func InitConfigByEtcd(etcdPeers string, target interface{}, dir string) {
	//etcd://127.0.0.1:2379,127.0.0.1:12379,127.0.0.1:22379/pub/pjoc/pay/config/alipay/
	url := fmt.Sprintf("etcd://%s%s", etcdPeers, dir)
	i := config.New(map[string]interface{}{})
	i.AddBackend(etcd.SCHEMA, &etcd.EtcdBackend{})
	if e := i.Init(config.URL(url), config.WithDefault(target)); e != nil {
		panic(e)
	}
}
