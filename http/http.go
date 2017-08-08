package http

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/ztgoto/webrouting/config"
	"github.com/ztgoto/webrouting/utils"
)

// ServerData 用于存储解析后的数据信息
type ServerData struct {
	Listen   string
	SSL      bool
	CertFile string
	KeyFile  string
	Rules    map[string][]*RequestMapping
}

var (
	w sync.WaitGroup
	// Key 上下文存储的key名称
	Key = "root"
	// ListenServerList 解析处理后的HTTP服务配置
	ListenServerList map[string]*ServerData
)

// StartServer 启动http服务
func StartServer() {
	w.Add(1)
	for k, v := range ListenServerList {
		log.Printf("Listing addr:%s\n", k)
		go func(sf *ServerData) {
			w.Add(1)
			defer w.Done()
			e := AddServer(sf)
			if e != nil {
				log.Println(e)
			}
		}(v)
	}
	w.Wait()
}

// AddServer 添加启动http服务
func AddServer(sf *ServerData) (e error) {
	c := &ServerContext{
		Data: sf,
		Key:  Key,
		h:    HandlerFun(DefaultHandle),
	}
	if sf.SSL {
		e = fasthttp.ListenAndServeTLS(sf.Listen, sf.CertFile, sf.KeyFile, c.RequestHandler)
	} else {
		e = fasthttp.ListenAndServe(sf.Listen, c.RequestHandler)
	}
	return
}

// PrepareSetting 初始设置
func PrepareSetting() error {

	err := resolverHTTPServer()
	if err != nil {
		return err
	}
	return nil
}

// resolverHTTPServer 处理解析HTTP服务信息
func resolverHTTPServer() error {
	httpServers := config.AppConf.HTTP.Servers
	if ListenServerList == nil {
		ListenServerList = make(map[string]*ServerData, len(httpServers))
	}
	for _, server := range httpServers {
		log.Println(server)
		listen := utils.SpaceRegexp.ReplaceAllString(server.Listen, "")
		if len(listen) == 0 {
			return errors.New("http.servers.listen is empty")
		}

		sname := utils.SpaceRegexp.ReplaceAllString(server.ServerName, "")
		names := strings.Split(sname, ",")

		if v, ok := ListenServerList[listen]; ok {
			v.SSL = server.SSL
			v.CertFile = server.CertFile
			v.KeyFile = server.KeyFile
			rules := v.Rules
			sloc := server.Locations
			if rules == nil {
				rules = make(map[string][]*RequestMapping, len(names))
			}
			for _, name := range names {
				nname := utils.SpaceRegexp.ReplaceAllString(name, "")
				if len(nname) == 0 {
					continue
				}
				if _, sok := rules[nname]; sok {
					continue
				} else {
					rms := make([]*RequestMapping, 0, len(sloc))
					for i := 0; i < len(sloc); i++ {
						rm, e := createRequestMapping(&sloc[i])
						if e != nil {
							return e
						}
						rms = append(rms, rm)
					}
					rules[nname] = rms
				}
			}

		} else {
			s := &ServerData{
				Listen:   listen,
				SSL:      server.SSL,
				CertFile: server.CertFile,
				KeyFile:  server.KeyFile,
			}

			rules := make(map[string][]*RequestMapping, len(names))
			sloc := server.Locations
			rms := make([]*RequestMapping, 0, len(sloc))
			for i := 0; i < len(sloc); i++ {
				rm, e := createRequestMapping(&sloc[i])
				if e != nil {
					return e
				}
				rms = append(rms, rm)
			}
			for _, name := range names {
				nname := utils.SpaceRegexp.ReplaceAllString(name, "")
				if len(nname) == 0 {
					continue
				}
				if _, ok := rules[nname]; ok {
					continue
				} else {
					rules[nname] = rms
				}
			}
			s.Rules = rules
			ListenServerList[listen] = s
		}

	}

	log.Printf("%+v\n", ListenServerList)
	return nil
}

// createRequestMapping 将配置转成路由映射
func createRequestMapping(lc *config.LocationConfig) (*RequestMapping, error) {
	if lc == nil {
		return nil, errors.New("param error")
	}
	pattern := lc.Pattern
	if len(utils.SpaceRegexp.ReplaceAllString(pattern, "")) == 0 {
		return nil, errors.New("pattern error")
	}

	proxyPass := lc.ProxyPass
	root := utils.SpaceRegexp.ReplaceAllString(lc.Root, "")
	index := utils.SpaceRegexp.ReplaceAllString(lc.Index, "")
	r := &RequestMapping{
		Pattern: pattern,
	}
	if len(proxyPass) > 0 {

		if v, ok := config.AppConf.UpStreams[proxyPass]; ok {
			h := &DefaultRoutingHandler{
				RoutingCnf:      &v,
				RequestHeaders:  lc.RequestHeaders,
				ResponseHeaders: lc.ResponseHeaders,
			}
			r.h = h
		} else {
			return nil, fmt.Errorf("upstreams key[%s] not exist", proxyPass)
		}

	} else if len(root) > 0 {
		fs := &fasthttp.FS{
			Root:               root,
			IndexNames:         strings.Split(index, ","),
			GenerateIndexPages: false,
			Compress:           false,
			AcceptByteRange:    false,
		}

		h := &DefaultFileHandler{
			handler: fs.NewRequestHandler(),
		}
		r.h = h
	}
	return r, nil
}

// mergeLocations 合并Location 信息
// 将s2合并到s1
// func mergeLocations(s1, s2 []LocationConfig) []LocationConfig {
// 	cap := len(s1) + len(s2)
// 	if cap == 0 {
// 		return nil
// 	}
// 	result := make([]LocationConfig, 0, cap)
// 	ut := make(map[string]int, cap)

// 	for _, s := range s1 {
// 		pattern := strings.TrimSpace(s.Pattern)
// 		if len(pattern) == 0 {
// 			panic("pattern is empty")
// 		}
// 		if index, ok := ut[pattern]; ok {
// 			lc := result[index]
// 			mapCopy(lc.RequestHeaders, s.RequestHeaders)
// 			mapCopy(lc.ResponseHeaders, s.ResponseHeaders)
// 		} else {
// 			result = append(result, s)
// 			ut[pattern] = len(result) - 1
// 		}
// 	}

// 	for _, s := range s2 {
// 		pattern := strings.TrimSpace(s.Pattern)
// 		if len(pattern) == 0 {
// 			panic("pattern is empty")
// 		}
// 		if index, ok := ut[pattern]; ok {
// 			lc := result[index]
// 			mapCopy(lc.RequestHeaders, s.RequestHeaders)
// 			mapCopy(lc.ResponseHeaders, s.ResponseHeaders)
// 		} else {
// 			result = append(result, s)
// 			ut[pattern] = len(result) - 1
// 		}
// 	}

// 	return result
// }

// mapCopy map拷贝
func mapCopy(target, source map[string]string) {
	if target != nil && source != nil {
		for k, v := range source {
			target[k] = v
		}
	}
}
