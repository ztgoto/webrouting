package http

import (
	"log"
	"net"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

var (
	w sync.WaitGroup

	// ListenList http服务监听列表
	ListenList map[string]net.Listener
)

func initHTTPServer() {
	httpConfigList := &config.GlobalConfig.HTTP.Server
	if len(*httpConfigList) > 0 && ListenList == nil {
		ListenList = make(map[string]net.Listener, len(*httpConfigList))
	}
	for _, v := range *httpConfigList {
		listen := v.Listen
		ssl := v.SSL
		var ln net.Listener
		var err error
		if ssl {
			cert := v.Cert
			key := v.Key
			ln, err = createServerTLS(listen, cert, key, func(ctx *fasthttp.RequestCtx) {

			})

		} else {
			ln, err = createServer(listen, func(ctx *fasthttp.RequestCtx) {

			})
		}
		if err != nil {
			panic(err)
		}
		ListenList[listen] = ln
	}
}

// 创建http服务器
func createServer(addr string, handler fasthttp.RequestHandler) (ln net.Listener, e error) {
	ln, e = net.Listen("tcp4", addr)
	if e != nil {
		return
	}

	go func() {
		e := fasthttp.Serve(ln, handler)
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
		e := fasthttp.ServeTLS(ln, cert, key, handler)
		if e != nil {
			panic(e)
		}
	}()
	log.Printf("create Listen [%s]\n", addr)
	return
}
