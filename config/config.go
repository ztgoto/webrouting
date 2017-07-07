package config

import (
	"sync"
)

type SysConfig struct {
	MaxProcs int "MaxProcs"
}

type HttpConfig struct {
	Addr  string "Addr"
	Vhost bool   "Vhost"
}

const (
	DefaultAddr  = "0.0.0.0:8080"
	DefaultVhost = false
)

var (
	Wait sync.WaitGroup
)

var ServerConfig HttpConfig = HttpConfig{
	Addr:  DefaultAddr,
	Vhost: DefaultVhost,
}

var SC SysConfig = SysConfig{
	MaxProcs: 1,
}
