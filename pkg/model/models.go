package model

type ConfigType string

var (
	// 数组的形式，存储为key/model的形式，即：如果在key目录下新增对象，则会解析成map，
	MAP ConfigType = "MAP"
	// 对象的形式，会将设置反序列化为对象
	MODEL ConfigType = "MODEL"
)

type Config interface {
	DirName() string
	ConfigType() ConfigType
}

func (b *BaseConfig) FullPath() string {
	return b.baseDir + b.dirName
}

func (b *BaseConfig) Key() string {
	return b.key
}

func (b *BaseConfig) BaseDir() string {
	return b.baseDir
}

var models = make([]Config, 0)

func RegisterModel(config Config) {
	models = append(models, config)
}

type BaseConfig struct {
	key      string `json:"-"`
	baseDir  string `json:"-"`
	dirName  string `json:"-"`
	fullPath string `json:"-"`
}

type PayConfig struct {
	ClusterId   *string `json:"cluster_id"`
	Concurrency int     `json:"concurrency"`
	// 通知地址的正则，必须包含{gateway_order_id}
	NotifyUrlPattern string `json:"notify_url_pattern"`
	// 跳转地址的正则，必须包含{gateway_order_id}
	ReturnUrlPattern string `json:"return_url_pattern"`
}

func (g *PayConfig) ConfigType() ConfigType {
	return MODEL
}

func (g *PayConfig) DirName() string {
	return "/cluster"
}

// ################## notice ##################
type NoticeConfig struct {
	NoticeIntervalSecond int `json:"notice_interval_second"`
	// 通知间隔
	//
	// 例如: [30, 30, 120, 240, 480, 1200, 3600, 7200, 43200, 86400, 172800]
	// 表示如果通知失败，则会隔 30s, 30s, 2min, 4min, 8min, 20min, 1H, 2H, 12H, 24H, 48H 通知
	NoticeDelaySecondExpressions []int `json:"notice_expressions"`
}

func (g *NoticeConfig) ConfigType() ConfigType {
	return MODEL
}

func (g *NoticeConfig) DirName() string {
	return "/notice"
}

// ################## service ##################
type ServiceConfigMap map[string]ServiceConfig
type ServiceConfig struct {
	ServiceName string `json:"service_name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
}

func (g *ServiceConfigMap) ConfigType() ConfigType {
	return MODEL
}

func (g *ServiceConfigMap) DirName() string {
	return "/service_map/"
}

// ################## channel ##################
type ChannelServiceConfigMap map[string]ChannelServiceConfig

type ChannelServiceConfig struct {
	BaseConfig
	ChannelId   string `json:"channel_id"`
	ServiceName string `json:"service_name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
}

func (g *ChannelServiceConfigMap) ConfigType() ConfigType {
	return MAP
}

func (g *ChannelServiceConfigMap) DirName() string {
	return "/channels/"
}

// ################## merchant ##################

type MerchantConfigMap map[string]MerchantConfig

type MerchantConfig struct {
	AppId                string `json:"app_id"`
	GatewayRSAPublicKey  string `json:"gateway_rsa_public_key"`
	GatewayRSAPrivateKey string `json:"gateway_rsa_private_key"`
	MerchantRSAPublicKey string `json:"merchant_rsa_public_key"`
	Md5Key               string `json:"md5_key"`
}

func (g *MerchantConfigMap) ConfigType() ConfigType {
	return MAP
}

func (g *MerchantConfigMap) DirName() string {
	return "/merchants/"
}

// ################## appid ##################
type AppIdAndChannelConfigMap map[string]AppIdAndChannelConfigs

type AppIdAndChannelConfigs struct {
	AppId          string               `json:"app_id"`
	ChannelConfigs []AppIdChannelConfig `json:"channel_configs"`
}

func (g *AppIdAndChannelConfigMap) ConfigType() ConfigType {
	return MAP
}

func (g *AppIdAndChannelConfigMap) DirName() string {
	return "/rates/"
}

type AppIdChannelConfig struct {
	RatePercent    float32 `json:"rate_percent"`
	Method         string  `json:"method"`
	ChannelAccount string  `json:"channel_account"`
	Available      bool    `json:"available"`
	ChannelId      string  `json:"channel_id"`
}

// ################## personal merchant ##################
type PersonalMerchantConfigMap map[string]PersonalMerchant

type PersonalMerchant struct {
	AppId                string `json:"app_id"`
	GatewayRSAPublicKey  string `json:"gateway_rsa_public_key"`
	GatewayRSAPrivateKey string `json:"gateway_rsa_private_key"`
	MerchantRSAPublicKey string `json:"merchant_rsa_public_key"`
	Md5Key               string `json:"md5_key"`
}

func (g *PersonalMerchantConfigMap) ConfigType() ConfigType {
	return MAP
}

func (g *PersonalMerchantConfigMap) DirName() string {
	return "/personal_merchants/"
}
