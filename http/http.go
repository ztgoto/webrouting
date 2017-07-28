package http

import (
	"github.com/ztgoto/webrouting/config"
)

// StartServer 启动http服务
func StartServer(cf *config.AppConfig) {

	// Start HTTP server.
	// if len((*cf).Addr) > 0 {

	// 	log.Printf("Starting HTTP server on %q", (*cf).Addr)
	// 	if err := fasthttp.ListenAndServe((*cf).Addr, RootHandler); err != nil {
	// 		log.Fatalf("error in ListenAndServe: %s", err)
	// 	}

	// }

}
