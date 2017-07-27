package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
)

// AppConfig 系统配置
type AppConfig struct {
	MaxProcs int    `json:"max_procs"`
	Addr     string `json:"http_addr"`
	Vhost    bool   `json:"vhost"`
}

const (
	// DefaultAddr 默认http监听地址
	DefaultAddr = "0.0.0.0:8080"
	// DefaultVhost 启用虚拟主机
	DefaultVhost = false
	// DefaultConfPath 默认配置文件路径
	DefaultConfPath = "./conf/webrouting.conf"
)

var (
	// AppConf 系统配置
	AppConf = &AppConfig{
		Addr:  DefaultAddr,
		Vhost: DefaultVhost,
	}

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

// PrepareSetting 初始设置
func PrepareSetting() {
	fmt.Printf("%+v\n", AppConf)
	runtime.GOMAXPROCS(AppConf.MaxProcs)
}
