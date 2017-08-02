package http

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

// StartServer 启动http服务
func StartServer() {

	// Start HTTP server.
	// if len((*cf).Addr) > 0 {

	// 	log.Printf("Starting HTTP server on %q", (*cf).Addr)
	// 	if err := fasthttp.ListenAndServe((*cf).Addr, RootHandler); err != nil {
	// 		log.Fatalf("error in ListenAndServe: %s", err)
	// 	}

	// }

}

// AddServer 添加启动http服务
func AddServer(sf *config.ServerData) {
	if sf.SSL {
		if err := fasthttp.ListenAndServeTLS(sf.Listen, sf.CertFile, sf.KeyFile, RootHandler); err != nil {
			fmt.Println(err)
		}
	} else {

		if err := fasthttp.ListenAndServe(sf.Listen, RootHandler); err != nil {
			fmt.Println(err)
		}
	}
}
