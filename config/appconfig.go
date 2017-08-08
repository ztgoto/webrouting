package config

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
)

// AppConfig 系统配置
type AppConfig struct {
	MaxProcs  int                       `json:"max_procs"`
	UpStreams map[string]UpStreamConfig `json:"upstreams"`
	HTTP      HTTPConfig                `json:"http"`
}

// UpStreamConfig 后端服务配置
type UpStreamConfig struct {
	Algorithm string         `json:"algorithm"`
	Timeout   int64          `json:"timeout"`
	Retries   int            `json:"retries"`
	Servers   []StreamConfig `json:"servers"`
}

// StreamConfig 路由列表配置
type StreamConfig struct {
	Addr   string `json:"addr"`
	Status int    `json:"status"`
	Weight int    `json:"weight"`
}

// HTTPConfig  http model config
type HTTPConfig struct {
	Servers []HTTPServerConfig `json:"servers"`
}

// HTTPServerConfig http服务配置
type HTTPServerConfig struct {
	Listen     string           `json:"listen"`
	ServerName string           `json:"server_name"`
	SSL        bool             `json:"ssl"`
	CertFile   string           `json:"cert_file"`
	KeyFile    string           `json:"key_file"`
	Locations  []LocationConfig `json:"locations"`
}

// LocationConfig 路径映射
type LocationConfig struct {
	Pattern         string            `json:"pattern"`
	ProxyPass       string            `json:"proxy_pass"`
	Root            string            `json:"root"`
	Index           string            `json:"index"`
	RequestHeaders  map[string]string `json:"request_headers"`
	ResponseHeaders map[string]string `json:"response_headers"`
}

const (
	// DefaultConfPath 默认配置文件路径
	DefaultConfPath = "./conf/webrouting.conf"
)

var (
	// AppConf 系统配置
	AppConf = &AppConfig{}

	// ConfPath 配置文件路径
	ConfPath = DefaultConfPath
)

func init() {
	AppConf.MaxProcs = runtime.NumCPU()
}

// LoadConf 从配置文件加载配置
func LoadConf() error {

	content, err := ioutil.ReadFile(ConfPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, AppConf)
	if err != nil {
		return err
	}
	return nil
}
