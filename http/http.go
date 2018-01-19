package http

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/ztgoto/webrouting/http/httphandler"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

var (
	w sync.WaitGroup

	// ListenList http服务监听列表
	ListenList map[string]net.Listener
)

// StartServer 启动服务
func StartServer() {

	initHTTPServer()
	config.DefaultLogger.Info("http server start success!")
	select {
	case <-config.CloseSignal:
		config.DefaultLogger.Info("---close server---")
		CloseServer()
	}
	w.Wait()
	config.DefaultLogger.Info("---all closed---")
}

// CloseServer 关闭服务
func CloseServer() {
	if ListenList != nil && len(ListenList) > 0 {
		for _, v := range ListenList {
			v.Close()
		}
	}
	config.DefaultLogger.Sync()
}

func initHTTPServer() {
	httpConfigList := &config.GlobalConfig.HTTP.Servers

	if len(*httpConfigList) > 0 && ListenList == nil {
		ListenList = make(map[string]net.Listener, len(*httpConfigList))
	}
	for _, v := range *httpConfigList {

		hostMap := toHostMap(&v)
		dispatch := httphandler.NewDefaultDispathc(hostMap)

		listen := v.Listen
		ssl := v.SSL
		var ln net.Listener
		var err error
		if ssl {
			cert := v.Cert
			key := v.Key
			ln, err = createServerTLS(listen, cert, key, dispatch.DoDispatch)

		} else {
			ln, err = createServer(listen, dispatch.DoDispatch)
		}
		if err != nil {
			panic(err)
		}
		ListenList[listen] = ln
		config.DefaultLogger.Info(fmt.Sprintf("http server start [%s]!", listen))
	}
}

func toHostMap(server *config.ServerConfig) map[string][]*config.LocationConfig {
	listen := server.Listen
	hosts := server.Hosts
	if hosts != nil && len(hosts) > 0 {
		hm := make(map[string][]*config.LocationConfig, 32)
		for _, v := range hosts {
			host := strings.TrimSpace(v.Host)
			if len(host) <= 0 {
				continue
			}

			if _, ok := hm[host]; ok {
				config.DefaultLogger.Info(fmt.Sprintf("listen:%s,host:%s,conflicting ignored", listen, host))
				continue
			}
			if v.Locations != nil && len(v.Locations) > 0 {
				lcs := make([]*config.LocationConfig, len(v.Locations))
				for i, loc := range v.Locations {
					lcs[i] = &loc
				}
				hm[host] = lcs
			}
		}
		return hm
	}
	return nil
}

// 创建http服务器
func createServer(addr string, handler fasthttp.RequestHandler) (ln net.Listener, e error) {
	ln, e = net.Listen("tcp4", addr)
	if e != nil {
		return
	}

	go func() {
		w.Add(1)
		e := fasthttp.Serve(ln, handler)
		w.Done()
		config.DefaultLogger.Info(fmt.Sprintf("http server[%s] closed!", addr))
		if e != nil {
			panic(e)
		}
	}()
	config.DefaultLogger.Info(fmt.Sprintf("create Listen [%s]", addr))
	return
}

// 创建https服务器
func createServerTLS(addr, cert, key string, handler fasthttp.RequestHandler) (ln net.Listener, e error) {
	ln, e = net.Listen("tcp4", addr)
	if e != nil {
		return
	}

	go func() {
		w.Add(1)
		e := fasthttp.ServeTLS(ln, cert, key, handler)
		w.Done()
		config.DefaultLogger.Info(fmt.Sprintf("https server[%s] closed!", addr))
		if e != nil {
			panic(e)
		}
	}()
	config.DefaultLogger.Info(fmt.Sprintf("create Listen [%s]\n", addr))
	return
}
