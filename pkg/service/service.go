package service

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pjoc-team/base-service/pkg/discovery"
	"github.com/pjoc-team/base-service/pkg/logger"
	_ "github.com/pjoc-team/etcd-config/config/etcd"
	"github.com/pjoc-team/etcd-config/etcdv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Service struct {
	DiscoveryService      discovery.DiscoveryService
	ListenAddr            string
	ConfigURI             string
	TlsEnable             bool
	LogLevel              string
	LogFormat             string
	CaCert                string
	TlsCert               string
	TlsKey                string
	ServiceName           string
	RegisterServiceToEtcd bool
	EtcdPeers             string
	ServiceDir            string
	logLevel              log.Level
}

type RegisterGrpcFunc = func(server *grpc.Server)

func WithConfigDir(path string) string {
	return filepath.Join(os.Getenv("HOME"), ".cert", path)
}

func (service *Service) parseLog() {
	loglevel, err := log.ParseLevel(service.LogLevel)
	if err != nil {
		log.Fatalf("log level error %s", err)
	}
	logger.Log.SetLevel(service.LogLevel)
	log.SetLevel(loglevel)
	service.logLevel = loglevel
	log.StandardLogger().SetNoLock()
	if service.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func InitService(
	listenAddr string,
	configURI string,
	tlsEnable bool,
	logLevel string,
	logFormat string,
	caCert string,
	tlsCert string,
	tlsKey string,
	serviceName string,
	registerServiceToEtcd bool,
	etcdPeers string,
	serviceDir string,
) (service *Service) {

	service = &Service{}

	service.ListenAddr = listenAddr
	service.ConfigURI = configURI
	service.TlsEnable = tlsEnable
	service.LogLevel = logLevel
	service.LogFormat = logFormat
	service.CaCert = caCert
	service.TlsCert = tlsCert
	service.TlsKey = tlsKey
	service.ServiceName = serviceName
	service.RegisterServiceToEtcd = registerServiceToEtcd
	service.EtcdPeers = etcdPeers
	service.ServiceDir = serviceDir

	if etcdPeers != "" {
		peers := strings.Split(etcdPeers, ",")
		etcdConfig := &etcdv3.EtcdConfig{}
		etcdConfig.Endpoints = peers
		etcdConfig.TimeoutSeconds = 6

		discoveryService := discovery.InitEtcdDiscoveryService(etcdConfig, serviceDir)
		if registerServiceToEtcd && serviceName != "" {
			endpoint := GetEndpoint(listenAddr)
			logger.Log.Infof("Register serviceName: %s endpoint: %v to DiscoveryService", serviceName, endpoint)
			discoveryService.RegisterService(serviceName, endpoint)
		}
		service.DiscoveryService = discoveryService
	}

	return service
}

func (service *Service) TlsGrpcOptions() []grpc.ServerOption {
	var grpcOpts []grpc.ServerOption
	if service.TlsEnable {
		cert, err := tls.LoadX509KeyPair(service.TlsCert, service.TlsKey)
		if err != nil {
			log.Fatal(err)
		}
		rawCaCert, err := ioutil.ReadFile(service.CaCert)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(rawCaCert)
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		})
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}
	return grpcOpts
}

func (service *Service) TlsDialGrpcOptions() []grpc.DialOption {
	var grpcOpts []grpc.DialOption
	if service.TlsEnable {
		creds, _ := credentials.NewClientTLSFromFile(service.CaCert, "")
		transportCredentials := grpc.WithTransportCredentials(creds)
		grpcOpts = append(grpcOpts, transportCredentials)
	} else {
		insecure := grpc.WithInsecure()
		grpcOpts = append(grpcOpts, insecure)
	}
	return grpcOpts
}

func (service *Service) StartGrpc(register RegisterGrpcFunc) {
	service.parseLog()

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}
	logrusEntry := log.NewEntry(log.StandardLogger())
	var grpcOpts []grpc.ServerOption
	if service.logLevel == log.DebugLevel {
		grpcOpts = []grpc.ServerOption{
			grpc_middleware.WithUnaryServerChain(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_logrus.StreamServerInterceptor(logrusEntry, opts...),
			),
		}
	}
	grpc_logrus.ReplaceGrpcLogger(log.NewEntry(log.StandardLogger()))
	options := service.TlsGrpcOptions()
	if options != nil && len(options) > 0 {
		for _, o := range options {
			grpcOpts = append(grpcOpts, o)
		}
	}

	gs := grpc.NewServer(grpcOpts...)
	// pb.RegisterPayChannelServer
	register(gs)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.helloservice", 1)
	healthpb.RegisterHealthServer(gs, healthServer)
	ln, err := net.Listen("tcp", service.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("running grpc on ", service.ListenAddr)
	if err := gs.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
