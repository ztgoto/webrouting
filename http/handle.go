package http

import (
	"bufio"
	"log"
	"net"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

// RootHandler 基础http请求处理器
func RootHandler(ctx *fasthttp.RequestCtx) {

	log.Printf("request path:%q\n", ctx.Path())
	log.Printf("Request Headers:\n%s\n", ctx.Request.Header.Header())
	log.Printf("Request Body:%v\n\n\n", ctx.Request.Body())

	serverConn, err := net.Dial("tcp", "127.0.0.1:80")

	if err != nil {
		log.Printf("end server exception%v\n", err)
		ctx.Response.SetStatusCode(config.HTTPStatusBadGateway)
		ctx.Response.SetBodyString("Bad Gateway")
		return
	}
	defer serverConn.Close()

	ctx.Request.WriteTo(serverConn)

	ctx.Response.Read(bufio.NewReader(serverConn))
	log.Printf("Response Headers:\n%s\n", ctx.Response.Header.Header())
	ctx.Response.Header.Set("Server", "fasthttp")
}
