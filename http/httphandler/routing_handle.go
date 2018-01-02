package httphandler

import (
	"log"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
)

// RoutingHandlerMapping 请求映射
type RoutingHandlerMapping struct {
	RoutingConfig map[string][]config.LocationConfig
}

// GetHandler 根据路径规则获取对应的处理器
func (rhm *RoutingHandlerMapping) GetHandler(ctx *fasthttp.RequestCtx) *HandlerExecutionChain {
	host := string(ctx.Request.Host())
	// host, _, e := net.SplitHostPort(string(ctx.Request.Host()))

	// if e != nil {
	// 	return nil
	// }
	log.Printf("host:%s\n", host)
	return nil
}
