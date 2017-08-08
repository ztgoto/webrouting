package http

import (
	"bufio"
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"regexp"
	"time"

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
	Data *ServerData
	Key  interface{}
	h    Handler
}

// RequestHandler 数据处理
func (s *ServerContext) RequestHandler(ctx *fasthttp.RequestCtx) {
	c := context.WithValue(context.Background(), s.Key, s.Data)
	e := s.h.Handle(c, ctx)
	if e != nil {
		log.Println(e)
	}
}

// RequestMapping http请求映射
type RequestMapping struct {
	Pattern string
	h       Handler
}

// DefaultFileHandler 文件处理handler
type DefaultFileHandler struct {
	handler fasthttp.RequestHandler
}

// Handle 默认文件处理器
func (h *DefaultFileHandler) Handle(c context.Context, ctx *fasthttp.RequestCtx) error {
	h.handler(ctx)
	return nil
}

// DefaultRoutingHandler 默认代理路由分发
type DefaultRoutingHandler struct {
	RoutingCnf      *config.UpStreamConfig
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
}

// Handle 默认路由处理器
func (h *DefaultRoutingHandler) Handle(c context.Context, ctx *fasthttp.RequestCtx) error {
	conn, e := h.randomConnection()
	if e != nil {
		ctx.Response.SetStatusCode(config.HTTPStatusBadGateway)
		ctx.Response.SetBodyString("Bad Gateway")
		return e
	}
	defer conn.Close()

	for k, v := range h.RequestHeaders {
		ctx.Request.Header.Set(k, v)
	}

	ctx.Request.WriteTo(conn)

	ctx.Response.Read(bufio.NewReader(conn))
	for k, v := range h.ResponseHeaders {
		ctx.Response.Header.Set(k, v)
	}

	return nil
}

func (h *DefaultRoutingHandler) randomConnection() (net.Conn, error) {
	algorithm := h.RoutingCnf.Algorithm
	servers := h.RoutingCnf.Servers
	timeout := time.Duration(config.DefaultTCPTimeout)
	retries := config.DefaultTCPRetries
	if h.RoutingCnf.Timeout > 0 {
		timeout = time.Duration(h.RoutingCnf.Timeout)
	}

	if h.RoutingCnf.Retries > 0 {
		retries = h.RoutingCnf.Retries
	}

	if len(servers) == 0 {
		return nil, errors.New("server list is empty")
	}
	if algorithm == "random" {
		index := rand.Intn(len(servers))
		uc := servers[index]
		if uc.Status == 1 {
			for i := 0; i < retries+1; i++ {
				log.Printf("random connection addr:%s\n", uc.Addr)
				conn, e := net.DialTimeout("tcp", uc.Addr, timeout*time.Millisecond)
				// 该处应该连接失败后对该服务器做故障排除处理(暂未实现)
				if e != nil {
					log.Println(e)
					continue
				}

				return conn, e
			}

		}

	}
	return nil, errors.New("server exception")
}

var (
	// RegexpCache 正则匹配缓存对象
	RegexpCache = make(map[string]*regexp.Regexp, 32)
)

// DefaultHandle 基础http请求处理器
func DefaultHandle(c context.Context, ctx *fasthttp.RequestCtx) error {

	cf := c.Value(Key).(*ServerData)
	log.Printf("%+v\n", cf)
	host, _, e := net.SplitHostPort(string(ctx.Request.Host()))

	if e != nil {
		return e
	}
	log.Printf("host:%s\n", host)

	rules := cf.Rules[host]

	if len(rules) == 0 {
		return nil
	}

	path := string(ctx.Path())

	log.Printf("request path:%q\n", path)

	for _, rule := range rules {
		pattern := rule.Pattern
		var reg *regexp.Regexp
		if v, ok := RegexpCache[pattern]; ok {
			reg = v
		} else {
			reg = regexp.MustCompile(pattern)
			RegexpCache[pattern] = reg
		}
		if reg.MatchString(path) {
			e := rule.h.Handle(nil, ctx)
			return e
		}
	}

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
