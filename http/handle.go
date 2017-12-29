package http

import (
	// "bufio"
	"errors"
	"log"
	"math/rand"
	"net"
	"regexp"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
	"github.com/ztgoto/webrouting/utils"
)

// "net/http" type Transport

// Handler 处理器
type Handler interface {
	Handle(c interface{}, ctx *fasthttp.RequestCtx) error
}

// HandlerFun 处理函数
type HandlerFun func(c interface{}, ctx *fasthttp.RequestCtx) error

// Handle 处理器
func (f HandlerFun) Handle(c interface{}, ctx *fasthttp.RequestCtx) error {
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
	// c := context.WithValue(context.Background(), s.Key, s.Data)
	e := s.h.Handle(s.Data, ctx)
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
func (h *DefaultFileHandler) Handle(c interface{}, ctx *fasthttp.RequestCtx) error {
	h.handler(ctx)
	return nil
}

// Stream 服务
type Stream struct {
	Addr   string
	Weight int
	Status int
}

// Streams 服务列表
type Streams struct {
	Algorithm string
	Timeout   int64
	Retries   int
	Up        []*Stream
	Clints    []*fasthttp.HostClient
	Lock      sync.RWMutex
}

// DefaultRoutingHandler 默认代理路由分发
type DefaultRoutingHandler struct {
	Routing         *Streams
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
}

// Handle 默认路由处理器
func (h *DefaultRoutingHandler) Handle(c interface{}, ctx *fasthttp.RequestCtx) error {
	// conn, e := h.getConnection()
	// if e != nil {
	// 	ctx.Response.SetStatusCode(config.HTTPStatusBadGateway)
	// 	ctx.Response.SetBodyString("Bad Gateway")
	// 	return e
	// }
	// defer conn.Close()

	// for k, v := range h.RequestHeaders {
	// 	ctx.Request.Header.Set(k, v)
	// }

	// n, err := ctx.Request.WriteTo(conn)
	// // 测试
	// log.Printf("WriteTo Length:%d\n", n)
	// if err != nil {
	// 	log.Println(err)
	// }

	// f, err := os.Create("C:\\Users\\rax\\Desktop\\testhttp.txt")
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer f.Close()

	// n1, err1 := ctx.Request.WriteTo(f)
	// if err1 != nil {
	// 	log.Println(err1)
	// }
	// log.Printf("WriteTo FileLength:%d\n", n1)

	// err = ctx.Response.Read(bufio.NewReader(conn))
	// if err != nil {
	// 	log.Println(err)
	// }
	// for k, v := range h.ResponseHeaders {
	// 	ctx.Response.Header.Set(k, v)
	// }

	algorithm := h.Routing.Algorithm
	clients := h.Routing.Clints
	timeout := time.Duration(config.DefaultTCPTimeout)

	if h.Routing.Timeout > 0 {
		timeout = time.Duration(h.Routing.Timeout)
	}

	if len(clients) == 0 {
		return errors.New("server list is empty")
	}

	client := clients[0]
	if algorithm == "random" {
		l := len(clients)
		index := rand.Intn(l)
		client = clients[index]
	}

	// Reset 'Connection: close' request header in order to prevent
	// from closing keep-alive connections to -out servers.
	ctx.Request.Header.ResetConnectionClose()

	for k, v := range h.RequestHeaders {
		ctx.Request.Header.Set(k, v)
	}

	e := client.DoTimeout(&ctx.Request, &ctx.Response, timeout*time.Millisecond)

	if e != nil {
		ctx.Response.SetStatusCode(config.HTTPStatusBadGateway)
		ctx.Response.SetBodyString("Bad Gateway")
		log.Println(e)
		return e
	}
	for k, v := range h.ResponseHeaders {
		ctx.Response.Header.Set(k, v)
	}

	return nil
}

func (h *DefaultRoutingHandler) getConnection() (net.Conn, error) {
	algorithm := h.Routing.Algorithm
	servers := h.Routing.Up
	timeout := time.Duration(config.DefaultTCPTimeout)
	retries := config.DefaultTCPRetries
	if h.Routing.Timeout > 0 {
		timeout = time.Duration(h.Routing.Timeout)
	}

	if h.Routing.Retries > 0 {
		retries = h.Routing.Retries
	}

	if len(servers) == 0 {
		return nil, errors.New("server list is empty")
	}
	if algorithm == "random" {
		// 存在并发问题
		for true {
			list := getOkAddr(servers)
			l := len(list)
			if l == 0 {
				return nil, errors.New("server exception")
			}
			index := rand.Intn(l)
			key := list[index]
			server := servers[key]
			for i := 0; i < retries+1; i++ {
				log.Printf("random connection addr:%s\n", server.Addr)
				conn, e := net.DialTimeout("tcp", server.Addr, timeout*time.Millisecond)

				if e != nil {
					log.Println(e)
					continue
				}

				return conn, e
			}
			// servers[key].Status = 0

		}

	}

	return nil, errors.New("server exception")
}

func getOkAddr(list []*Stream) []int {
	oks := make([]int, 0, len(list))
	for i, v := range list {
		if v.Status == 1 {
			oks = append(oks, i)
		}
	}
	return oks
}

var (
	// RegexpCache 正则匹配缓存对象
	RegexpCache = utils.NewConcurrentMap(32)
)

// DefaultHandle 基础http请求处理器
func DefaultHandle(c interface{}, ctx *fasthttp.RequestCtx) error {

	cf := c.(*ServerData)
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
		if v := RegexpCache.Get(pattern); v != nil {
			reg = v.(*regexp.Regexp)
		} else {
			reg = regexp.MustCompile(pattern)
			RegexpCache.Put(pattern, reg)
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
