package httphandler

import (
	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

// "bufio"

// "net/http" type Transport

// Handler 处理器
type Handler interface {
	Handle(*fasthttp.RequestCtx)
}

// HandlerMapping 处理器映射
type HandlerMapping interface {
	GetHandler(*fasthttp.RequestCtx) *HandlerExecutionChain
}

// HandlerInterceptor 拦截器
type HandlerInterceptor interface {
	PreHandle(*fasthttp.RequestCtx) bool
	PostHandle(*fasthttp.RequestCtx)
	AfterCompletion(*fasthttp.RequestCtx)
}

// HandlerExecutionChain 执行链
type HandlerExecutionChain struct {
	interceptorIndex int
	handler          Handler
	interceptors     []HandlerInterceptor
}

func (hec *HandlerExecutionChain) applyPreHandle(ctx *fasthttp.RequestCtx) bool {
	if hec.interceptors != nil && len(hec.interceptors) > 0 {
		for i, v := range hec.interceptors {
			if !v.PreHandle(ctx) {
				hec.triggerAfterCompletion(ctx)
				return false
			}
			hec.interceptorIndex = i
		}
	}
	return true
}

func (hec *HandlerExecutionChain) applyPostHandle(ctx *fasthttp.RequestCtx) {
	if hec.interceptors != nil && len(hec.interceptors) > 0 {
		for i := len(hec.interceptors); i >= 0; i-- {
			hi := hec.interceptors[i]
			hi.PostHandle(ctx)
		}
	}
}

func (hec *HandlerExecutionChain) triggerAfterCompletion(ctx *fasthttp.RequestCtx) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	if hec.interceptors != nil && len(hec.interceptors) > 0 {
		for i := hec.interceptorIndex; i >= 0; i-- {
			hi := hec.interceptors[i]
			hi.AfterCompletion(ctx)
		}
	}
}

// Dispatch 路由分发
type Dispatch struct {
	handlerMappings []HandlerMapping
}

// DoDispatch 处理器
func (rd *Dispatch) DoDispatch(ctx *fasthttp.RequestCtx) {
	// httphost := string(ctx.Request.Host())
	// path := string(ctx.Path())
	// log.Printf("httphost:%s,uri:%s\n", httphost, path)

	hec := rd.getHandler(ctx)

	if hec == nil {
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	handler := hec.handler

	if !hec.applyPreHandle(ctx) {
		return
	}
	handler.Handle(ctx)
	hec.applyPostHandle(ctx)
}

func (rd *Dispatch) getHandler(ctx *fasthttp.RequestCtx) *HandlerExecutionChain {
	if rd.handlerMappings != nil && len(rd.handlerMappings) > 0 {
		for _, v := range rd.handlerMappings {
			handler := v.GetHandler(ctx)
			if handler != nil {
				return handler
			}
		}
	}
	return nil
}

// NewDefaultDispathc 创建dispatch
func NewDefaultDispathc(lc map[string][]*config.LocationConfig) *Dispatch {

	return &Dispatch{
		handlerMappings: []HandlerMapping{
			&RoutingHandlerMapping{
				LocConfig: lc,
			},
			&StaticFileHandlerMapping{
				LocConfig: lc,
			},
		},
	}
}
