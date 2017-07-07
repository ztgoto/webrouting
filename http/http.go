package http

import (
	"bufio"
	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
	"log"
	"net"
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

	log.Printf("request path:%q\n", ctx.Path())
	log.Printf("Headers:\n%s\n", ctx.Request.Header.Header())
	log.Printf("Body:%v\n\n\n", ctx.Request.Body())

	serverConn, err := net.Dial("tcp", "192.168.4.132:8026")

	if err != nil {
		log.Printf("end server exception%v\n", err)
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBodyString("server error")
		return
	}
	defer serverConn.Close()

	ctx.Request.WriteTo(serverConn)

	ctx.Response.Read(bufio.NewReader(serverConn))

}
