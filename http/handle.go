package http

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

// Handler 处理器
type Handler interface {
	Handle(c context.Context, ctx *fasthttp.RequestCtx) error
}

// HandlerFun 处理函数
type HandlerFun func(c context.Context, ctx *fasthttp.RequestCtx) error

// Handle 处理器
func (f HandlerFun) Handle(c context.Context, ctx *fasthttp.RequestCtx) error {
	e := f(c, ctx)
	return e
}

// ServerContext http服务上下文数据
type ServerContext struct {
	Data *config.ServerData
	K    interface{}
	h    Handler
}

// RequestHandler 数据处理
func (s *ServerContext) RequestHandler(ctx *fasthttp.RequestCtx) {
	c := context.WithValue(context.Background(), s.K, s.Data)
	e := s.h.Handle(c, ctx)
	if e != nil {
		fmt.Println(e)
	}
}

// DefaultHandle 基础http请求处理器
func DefaultHandle(c context.Context, ctx *fasthttp.RequestCtx) error {

	// log.Printf("request path:%q\n", ctx.Path())
	// log.Printf("Request Headers:\n%s\n", ctx.Request.Header.Header())
	// log.Printf("Request Body:%v\n\n\n", ctx.Request.Body())
	// log.Printf("Request Host:-%s-\n\n\n", ctx.Request.Host())

	// serverConn, err := net.Dial("tcp", "127.0.0.1:80")

	// if err != nil {
	// 	log.Printf("end server exception%v\n", err)
	// 	ctx.Response.SetStatusCode(config.HTTPStatusBadGateway)
	// 	ctx.Response.SetBodyString("Bad Gateway")
	// 	return err
	// }
	// defer serverConn.Close()

	// ctx.Request.WriteTo(serverConn)

	// ctx.Response.Read(bufio.NewReader(serverConn))
	// log.Printf("Response Headers:\n%s\n", ctx.Response.Header.Header())
	// ctx.Response.Header.Set("Server", "fasthttp")
	return nil
}
