package config

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/config/yaml"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pjoc-team/base-service/pkg/logger"
)

var Log = logger.Log

type InitHandle func(*Config) error

type Config struct {
	Key     string
	Mu      *sync.Mutex
	Init    InitHandle
	EtcdKey string
	Target  interface{}
}

type ConfigManage struct {
	Source map[string]interface{}
	Target []*Config
}

func NewConfigManage() *ConfigManage {
	return &ConfigManage{Source: make(map[string]interface{}), Target: make([]*Config, 0)}
}
func (c *ConfigManage) RegisterAfter(key string, mu *sync.Mutex, target interface{}, handle InitHandle) {
	if mu == nil {
		mu = &sync.Mutex{}
	}
	c.Target = append(c.Target, &Config{Target: target, Mu: mu, Key: key, Init: handle})
}

func (c *ConfigManage) RegisterBefore(key string, mu *sync.Mutex, target interface{}, handle InitHandle) {
	if mu == nil {
		mu = &sync.Mutex{}
	}
	c.Target = append([]*Config{&Config{Target: target, Mu: mu, Key: key, Init: handle}}, c.Target...)
}
func (c *ConfigManage) Register(key string, mu *sync.Mutex, target interface{}, handle InitHandle) {
	c.RegisterAfter(key, mu, target, handle)
}
func (c *ConfigManage) GetEtcdKey(config *Config) {
	key := config.Key
	v, ok := c.Source[key]
	if !ok {
		return
	}
	value, ok := v.(map[interface{}]interface{})
	if ok {
		etcd := value["EtcdKey"]
		if etcd != nil {
			config.EtcdKey, ok = etcd.(string)
		}
	}
	etcdKey, ok := c.Source[config.Key+"EtcdKey"]
	if ok {
		config.EtcdKey, _ = etcdKey.(string)
	}

}

func (c *ConfigManage) UpdateValue(config *Config) {
	if config.Mu != nil {
		config.Mu.Lock()
		defer config.Mu.Unlock()
	}
	key := config.Key
	v, ok := c.Source[key]
	if !ok {
		return
	}
	c.GetEtcdKey(config)
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	valueRef := reflect.ValueOf(v)
	targetType := reflect.TypeOf(config.Target)
	for targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	//只适配Slice 数组 和 map 的情况
	if valueRef.Kind() == reflect.Slice && targetType.Kind() == reflect.Map {
		newMap := map[interface{}]interface{}{}
		arrayIter := v.([]interface{})
		for i, v := range arrayIter {
			newMap[fmt.Sprintf("%d", i)] = v
		}
		bytes, err = yaml.Marshal(newMap)
		v = newMap
	}

	target := config.Target
	err = yaml.Unmarshal(bytes, target)
	if err != nil {
		panic(fmt.Sprintf("%s %s %s %s", key, err, valueRef.Kind(), targetType.Kind()))
	}
	Log.Infof("config updatevalue from file %s %#v", config.Key, *config)
	return
}

//InitTarget  初始化各个目标
func (c *ConfigManage) InitTarget(target *Config) {
	v := target
	if v.Init != nil {
		v.Mu.Lock()
		v.Init(v)
		v.Mu.Unlock()
	}
	Log.Debugf("config end etcd %s etcd key [%s] %#v", v.Key, v.EtcdKey, v.Target)
}

func (c *ConfigManage) Init(file string) {
	if err := yaml.UnmarshalFromFile(file, c.Source); err != nil {
		panic(err)
	}
	for _, config := range c.Target {
		c.UpdateValue(config)
	}
	for _, v := range c.Target {
		c.InitTarget(v)
	}
	//update etcd
}

var StdConfig = NewConfigManage()

//Register 注册配置文件事件
//更新配置以及调用handle 会加锁
func Register(key string, mu *sync.Mutex, target interface{}, handle InitHandle) {
	StdConfig.Register(key, mu, target, handle)
}
func RegisterBefore(key string, mu *sync.Mutex, target interface{}, handle InitHandle) {
	StdConfig.RegisterBefore(key, mu, target, handle)
}

func ParseDuration(t string, def time.Duration) (time.Duration, error) {
	if t == "" {
		return def, nil
	}
	unit := strings.TrimLeft(t, "1234567890")
	itime, err := strconv.Atoi(strings.TrimRight(t, "ms"))
	if err != nil {
		return def, fmt.Errorf("error format :%s %s", t, err)
	}
	switch unit {
	case "ms":
		return time.Duration(itime) * time.Millisecond, nil
	case "s":
		return time.Duration(itime) * time.Second, nil
	default:
		return def, fmt.Errorf("time unknow unit:%s", unit)
	}

}

func GetDuration(t string, def time.Duration) time.Duration {
	o, e := ParseDuration(t, def)
	if e != nil {
		Log.Errorf("parse Duration error %s %s", t, e)
		return def
	}
	return o
}
