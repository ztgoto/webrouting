package config

import (
	"io/ioutil"
	"runtime"

	"gopkg.in/yaml.v2"
)

// ApplicationConfig 应用配置
type ApplicationConfig struct {
	Processes int
}

// UpstreamConfig 后端服务配置
type UpstreamConfig struct {
	ID      string
	Balance string
	Timeout int64
	Servers []string
}

// LocationConfig 路由配置
type LocationConfig struct {
	Pattern  string
	Upstream string
	Root     string
	Index    string
	Request  map[string]string
	Response map[string]string
}

// HostMappingConfig host路由配置
type HostMappingConfig struct {
	Host      string
	Locations []LocationConfig
}

// ServerConfig HTTP服务配置
type ServerConfig struct {
	Listen string
	SSL    bool
	Cert   string
	Key    string
	Hosts  []HostMappingConfig
}

// HTTPConfig 全局Http配置
type HTTPConfig struct {
	Servers []ServerConfig
}

// Config 全局配置对象
type Config struct {
	Application ApplicationConfig
	Upstreams   []UpstreamConfig
	HTTP        HTTPConfig
}

const (
	// DefaultConfPath 默认配置文件路径
	DefaultConfPath = "./conf/webrouting.conf"
)

var (
	// GlobalConfig 系统配置
	GlobalConfig = &Config{}

	// ConfPath 配置文件路径
	ConfPath = DefaultConfPath
)

func init() {
	GlobalConfig.Application.Processes = runtime.NumCPU()
}

// ParseConfig 解析Application配置
func ParseConfig(in []byte) (app *Config, err error) {
	app = &Config{}
	err = yaml.Unmarshal(in, app)
	return
}

// LoadConfigFile 读取配置文件
func LoadConfigFile() error {
	content, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, GlobalConfig)
	if err != nil {
		return err
	}
	return nil
}
