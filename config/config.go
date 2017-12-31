package config

import (
	"gopkg.in/yaml.v2"
)

// ApplicationConfig 应用配置
type ApplicationConfig struct {
	Processes int
}

type UpstreamConfig struct {
	Id string
	Balance string
	Server []string
}

type LocationConfig struct {
	Pattern string
	Proxy string
	Root string
	Index string
	Request map[string]string
	Response map[string]string
}

type ServerConfig struct {
	Listen int
	Host string
	SSL        bool
	Cert   string
	Key    string
	Location []LocationConfig
}

type HttpConfig struct {
	Server []ServerConfig
}

type Config struct {
	Application ApplicationConfig
	Upstream []UpstreamConfig
	Http  HttpConfig
}



// ParseConfig 解析Application配置
func ParseConfig(in []byte) (app *Config, err error) {
	app=&Config{}
	err = yaml.Unmarshal(in,app);
	return
}
