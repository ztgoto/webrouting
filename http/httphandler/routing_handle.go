package httphandler

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
	"github.com/ztgoto/webrouting/utils"
)

var (
	// RegexpCache 正则匹配缓存对象
	RegexpCache = utils.NewConcurrentMap(32)
	clientsMap  = make(map[string][]*fasthttp.HostClient, 32)
)

// RoutingHandlerMapping 反向代理请求映射
type RoutingHandlerMapping struct {
	LocConfig  map[string][]*config.LocationConfig
	handlerMap map[*config.LocationConfig]*RoutingHandler
}

// GetHandler 根据路径规则获取对应的处理器
func (rhm *RoutingHandlerMapping) GetHandler(ctx *fasthttp.RequestCtx) *HandlerExecutionChain {
	httphost := string(ctx.Request.Host())
	path := string(ctx.Path())
	if rhm.LocConfig == nil || len(rhm.LocConfig) <= 0 {
		return nil
	}

	host := strings.Split(httphost, ":")[0]

	lcs, ok := rhm.LocConfig[host]
	if !ok || lcs == nil || len(lcs) <= 0 {
		return nil
	}

	hitlc := matchLocationConfig(lcs, path)

	if hitlc == nil {
		return nil
	}

	if rhm.handlerMap == nil {
		rhm.handlerMap = make(map[*config.LocationConfig]*RoutingHandler, 32)
	}

	rh, ok := rhm.handlerMap[hitlc]

	if !ok {
		proxy := hitlc.Upstream

		if len(strings.TrimSpace(proxy)) == 0 {
			return nil
		}

		uc := findUpstreamConfig(proxy)
		if uc == nil {
			panic(fmt.Sprintf("Upstream id[%s] not found", proxy))
		}

		rh = NewRoutingHandler(hitlc, uc)
		if rh == nil {
			return nil
		}
		rhm.handlerMap[hitlc] = rh
	}

	return &HandlerExecutionChain{
		interceptorIndex: -1,
		handler:          rh,
	}
}

func matchLocationConfig(lcs []*config.LocationConfig, path string) *config.LocationConfig {
	var pattern string
	var reg *regexp.Regexp
	var hitlc *config.LocationConfig
	for _, lc := range lcs {
		pattern = strings.TrimSpace(lc.Pattern)
		rt := RegexpCache.Get(pattern)
		if rt == nil {
			reg = regexp.MustCompile(pattern)
			RegexpCache.Put(pattern, reg)
		} else {
			reg = rt.(*regexp.Regexp)
		}
		if reg.MatchString(path) {
			hitlc = lc
			break
		}
	}
	return hitlc
}

func findUpstreamConfig(ucID string) *config.UpstreamConfig {
	if config.GlobalConfig.Upstreams != nil && len(config.GlobalConfig.Upstreams) > 0 {
		for _, v := range config.GlobalConfig.Upstreams {
			if strings.TrimSpace(ucID) == strings.TrimSpace(v.ID) {
				return &v
			}
		}
	}
	return nil
}

// NewRoutingHandler 创建反向代理处理器
func NewRoutingHandler(lc *config.LocationConfig, uc *config.UpstreamConfig) *RoutingHandler {
	ucID := strings.TrimSpace(uc.ID)
	if len(ucID) == 0 {
		panic("UpstreamConfig ID is empty")
	}

	clients, ok := clientsMap[ucID]
	if !ok {
		servers := uc.Servers
		if servers == nil || len(servers) <= 0 {
			panic("UpstreamConfig server list is empty")
		}

		clients = make([]*fasthttp.HostClient, len(servers))

		// 此处配置还需要细化
		for i, v := range servers {
			clients[i] = &fasthttp.HostClient{
				Addr:         v,
				Dial:         fasthttp.Dial,
				MaxConns:     config.DefaultClientMaxConnCount,
				ReadTimeout:  120 * time.Second,
				WriteTimeout: 5 * time.Second,
				// ReadBufferSize: *outMaxHeaderSize,
			}
		}
		clientsMap[ucID] = clients

	}

	balance := uc.Balance

	if len(strings.TrimSpace(balance)) == 0 {
		balance = "random"
	}

	timeout := time.Duration(config.DefaultRequestTimeout) * time.Millisecond

	if uc.Timeout > 0 {
		timeout = time.Duration(uc.Timeout) * time.Millisecond
	}

	return &RoutingHandler{
		Clients: clients,
		Balance: balance,
		Timeout: timeout,
		lc:      lc,
	}
}

// RoutingHandler 反向代理处理器
type RoutingHandler struct {
	Clients []*fasthttp.HostClient
	Balance string
	Timeout time.Duration
	lc      *config.LocationConfig
}

// Handle 反向代理处理器
func (rh *RoutingHandler) Handle(ctx *fasthttp.RequestCtx) {
	if rh.Clients == nil || len(rh.Clients) <= 0 {
		return
	}
	client := rh.Clients[0]
	balance := rh.Balance
	if balance == "random" {
		l := len(rh.Clients)
		index := rand.Intn(l)
		client = rh.Clients[index]
	}

	timeout := 30000 * time.Millisecond

	if rh.Timeout > 0 {
		timeout = rh.Timeout
	}

	ctx.Request.Header.ResetConnectionClose()

	if rh.lc != nil && rh.lc.Request != nil && len(rh.lc.Request) > 0 {
		for k, v := range rh.lc.Request {
			ctx.Request.Header.Set(k, v)
		}
	}

	e := client.DoTimeout(&ctx.Request, &ctx.Response, timeout*time.Millisecond)

	if rh.lc != nil && rh.lc.Response != nil && len(rh.lc.Response) > 0 {
		for k, v := range rh.lc.Response {
			ctx.Response.Header.Set(k, v)
		}
	}

	if e != nil {
		ctx.Response.SetStatusCode(config.HTTPStatusBadGateway)
		ctx.Response.SetBodyString("Bad Gateway")
		log.Println(e)
		return
	}
}

// StaticFileHandlerMapping 静态文件服务请求映射
type StaticFileHandlerMapping struct {
	LocConfig  map[string][]*config.LocationConfig
	handlerMap map[*config.LocationConfig]*DefaultFileHandler
}

// GetHandler 获取处理器
func (sfhm *StaticFileHandlerMapping) GetHandler(ctx *fasthttp.RequestCtx) *HandlerExecutionChain {
	httphost := string(ctx.Request.Host())
	path := string(ctx.Path())
	if sfhm.LocConfig == nil || len(sfhm.LocConfig) <= 0 {
		return nil
	}

	host := strings.Split(httphost, ":")[0]

	lcs, ok := sfhm.LocConfig[host]
	if !ok || lcs == nil || len(lcs) <= 0 {
		return nil
	}

	hitlc := matchLocationConfig(lcs, path)

	if hitlc == nil {
		return nil
	}

	if sfhm.handlerMap == nil {
		sfhm.handlerMap = make(map[*config.LocationConfig]*DefaultFileHandler, 32)
	}

	h, ok := sfhm.handlerMap[hitlc]

	if !ok {
		rootpath := hitlc.Root

		if len(strings.TrimSpace(rootpath)) == 0 {
			return nil
		}

		h = NewDefaultFileHandler(hitlc)
		if h == nil {
			return nil
		}
		sfhm.handlerMap[hitlc] = h

	}

	return &HandlerExecutionChain{
		interceptorIndex: -1,
		handler:          h,
	}
}

// NewDefaultFileHandler 创建文件处理器
func NewDefaultFileHandler(lc *config.LocationConfig) *DefaultFileHandler {

	fs := &fasthttp.FS{
		Root:               lc.Root,
		IndexNames:         strings.Split(lc.Index, ","),
		GenerateIndexPages: false,
		Compress:           false,
		AcceptByteRange:    false,
	}

	return &DefaultFileHandler{
		handler: fs.NewRequestHandler(),
		lc:      lc,
	}
}

// DefaultFileHandler 文件处理handler
type DefaultFileHandler struct {
	handler fasthttp.RequestHandler
	lc      *config.LocationConfig
}

// Handle 默认文件处理器
func (h *DefaultFileHandler) Handle(ctx *fasthttp.RequestCtx) {
	if h.lc != nil && h.lc.Request != nil && len(h.lc.Request) > 0 {
		for k, v := range h.lc.Request {
			ctx.Request.Header.Set(k, v)
		}
	}
	h.handler(ctx)
	if h.lc != nil && h.lc.Response != nil && len(h.lc.Response) > 0 {
		for k, v := range h.lc.Response {
			ctx.Response.Header.Set(k, v)
		}
	}
}
