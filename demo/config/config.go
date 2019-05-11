package config

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/service"
)

func main() {
	etcdPeers := "127.0.0.1:2379"
	config := service.InitGatewayConfig(etcdPeers, "/pub/pjoc/pay")
	fmt.Println(config.AppIdAndMerchantMap)
	fmt.Println(config.ServiceMap)
}
