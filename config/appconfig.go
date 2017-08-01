package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/ztgoto/webrouting/utils"
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
	Servers   []StreamConfig `json:"servers"`
}

// StreamConfig 路由列表配置
type StreamConfig struct {
	Addr   string `json:"addr"`
	Status int    `json:"status"`
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

// ServerData 用于存储解析后的数据信息
type ServerData struct {
	Listen   string
	SSL      bool
	CertFile string
	KeyFile  string
	Rules    map[string][]LocationConfig
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

	// ListenServerList 解析处理后的HTTP服务配置
	ListenServerList map[string]*ServerData
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

// PrepareSetting 初始设置
func PrepareSetting() error {
	err := LoadConf()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", AppConf)
	runtime.GOMAXPROCS(AppConf.MaxProcs)

	err = resolverHTTPServer()
	if err != nil {
		return err
	}
	return nil
}

// resolverHTTPServer 处理解析HTTP服务信息
func resolverHTTPServer() error {
	httpServers := AppConf.HTTP.Servers
	if ListenServerList == nil {
		httpServers := AppConf.HTTP.Servers
		ListenServerList = make(map[string]*ServerData, len(httpServers))
	}
	for _, server := range httpServers {
		fmt.Println(server)
		listen := utils.SpaceRegexp.ReplaceAllString(server.Listen, "")
		if len(listen) == 0 {
			return errors.New("http.servers.listen is empty")
		}

		sname := utils.SpaceRegexp.ReplaceAllString(server.ServerName, "")
		names := strings.Split(sname, ",")

		if v, ok := ListenServerList[listen]; ok {
			v.SSL = server.SSL
			v.CertFile = server.CertFile
			v.KeyFile = server.KeyFile

		} else {
			s := &ServerData{
				Listen:   listen,
				SSL:      server.SSL,
				CertFile: server.CertFile,
				KeyFile:  server.KeyFile,
			}

			var rules map[string][]LocationConfig
			for _, name := range names {
				if len(name) > 0 {
					if rules == nil {
						rules = make(map[string][]LocationConfig, len(names))
					}
					rules[name] = server.Locations
				}
			}
			s.Rules = rules
			ListenServerList[listen] = s
		}

	}

	fmt.Printf("%+v\n", ListenServerList)
	return nil
}

// mergeLocations 合并Location 信息
// 将s2合并到s1
func mergeLocations(s1, s2 []LocationConfig) []LocationConfig {
	cap := len(s1) + len(s2)
	if cap == 0 {
		return nil
	}
	result := make([]LocationConfig, 0, cap)
	ut := make(map[string]int, cap)

	for _, s := range s1 {
		pattern := strings.TrimSpace(s.Pattern)
		if len(pattern) == 0 {
			panic("pattern is empty")
		}
		if index, ok := ut[pattern]; ok {
			lc := result[index]
			mapCopy(lc.RequestHeaders, s.RequestHeaders)
			mapCopy(lc.ResponseHeaders, s.ResponseHeaders)
		} else {
			result = append(result, s)
			ut[pattern] = len(result) - 1
		}
	}

	for _, s := range s2 {
		pattern := strings.TrimSpace(s.Pattern)
		if len(pattern) == 0 {
			panic("pattern is empty")
		}
		if index, ok := ut[pattern]; ok {
			lc := result[index]
			mapCopy(lc.RequestHeaders, s.RequestHeaders)
			mapCopy(lc.ResponseHeaders, s.ResponseHeaders)
		} else {
			result = append(result, s)
			ut[pattern] = len(result) - 1
		}
	}

	return result
}

// mapCopy map拷贝
func mapCopy(target, source map[string]string) {
	if target != nil && source != nil {
		for k, v := range source {
			target[k] = v
		}
	}
}
