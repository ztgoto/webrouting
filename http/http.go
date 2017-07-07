package http

import (
	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
	"log"
)

func StartServer(cf *config.HttpConfig) {

	// Start HTTP server.
	if len((*cf).Addr) > 0 {

		log.Printf("Starting HTTP server on %q", (*cf).Addr)
		if err := fasthttp.ListenAndServe((*cf).Addr, request); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}

	}

}

func request(ctx *fasthttp.RequestCtx) {
	log.Println(ctx.Path())
}
