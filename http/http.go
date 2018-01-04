package http

import (
	"log"
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
	log.Println("http server start success!")
	select {
	case <-config.CloseSignal:
		log.Println("---close server---")
		CloseServer()
	}
	w.Wait()
	log.Println("---all closed---")
}

// CloseServer 关闭服务
func CloseServer() {
	if ListenList != nil && len(ListenList) > 0 {
		for _, v := range ListenList {
			v.Close()
		}
	}
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
		log.Printf("http server start [%s]!\n", listen)
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
				log.Printf("listen:%s,host:%s,conflicting ignored", listen, host)
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
		log.Printf("http server[%s] closed!", addr)
		if e != nil {
			panic(e)
		}
	}()
	log.Printf("create Listen [%s]\n", addr)
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
		log.Printf("http server[%s] closed!", addr)
		if e != nil {
			panic(e)
		}
	}()
	log.Printf("create Listen [%s]\n", addr)
	return
}
