package http

import (
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

var (
	w sync.WaitGroup
)

// StartServer 启动http服务
func StartServer() {
	w.Add(1)
	for k, v := range config.ListenServerList {
		fmt.Printf("Listing addr:%s\n", k)
		go func(sf *config.ServerData) {
			w.Add(1)
			defer w.Done()
			e := AddServer(sf)
			if e != nil {
				fmt.Println(e)
			}
		}(v)
	}
	w.Wait()
}

// AddServer 添加启动http服务
func AddServer(sf *config.ServerData) (e error) {
	c := &ServerContext{
		Data: sf,
		K:    "root",
		h:    HandlerFun(DefaultHandle),
	}
	if sf.SSL {
		e = fasthttp.ListenAndServeTLS(sf.Listen, sf.CertFile, sf.KeyFile, c.RequestHandler)
	} else {
		e = fasthttp.ListenAndServe(sf.Listen, c.RequestHandler)
	}
	return
}
