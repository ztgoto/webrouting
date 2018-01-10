package client

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var errorChPool sync.Pool
var timerPool sync.Pool

// BaseClient 自定义http客户端 覆盖HostClient部分方法
type BaseClient struct {
	fasthttp.HostClient
}

// DoDeadline 覆盖fasthttp.HostClient.DoDeadline 方法
// 因原方法中(fasthttp.clientDoDeadline)是通过request response copy来实现，发现在做代理应用时处理文件上传，文件流会丢失掉无法转发到后台服务器
func (c *BaseClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	timeout := -time.Since(deadline)
	if timeout <= 0 {
		return fasthttp.ErrTimeout
	}

	var ch chan error
	chv := errorChPool.Get()
	if chv == nil {
		chv = make(chan error, 1)
	}
	ch = chv.(chan error)

	// Note that the request continues execution on ErrTimeout until
	// client-specific ReadTimeout exceeds. This helps limiting load
	// on slow hosts by MaxConns* concurrent requests.
	//
	// Without this 'hack' the load on slow host could exceed MaxConns*
	// concurrent requests, since timed out requests on client side
	// usually continue execution on the host.
	go func() {
		ch <- c.HostClient.Do(req, resp)
	}()

	tc := acquireTimer(timeout)
	var err error
	select {
	case err = <-ch:
		errorChPool.Put(chv)
	case <-tc.C:
		err = fasthttp.ErrTimeout
	}
	releaseTimer(tc)

	return err

}

func initTimer(t *time.Timer, timeout time.Duration) *time.Timer {
	if t == nil {
		return time.NewTimer(timeout)
	}
	if t.Reset(timeout) {
		panic("BUG: active timer trapped into initTimer()")
	}
	return t
}

func stopTimer(t *time.Timer) {
	if !t.Stop() {
		// Collect possibly added time from the channel
		// if timer has been stopped and nobody collected its' value.
		select {
		case <-t.C:
		default:
		}
	}
}

func acquireTimer(timeout time.Duration) *time.Timer {
	v := timerPool.Get()
	if v == nil {
		return time.NewTimer(timeout)
	}
	t := v.(*time.Timer)
	initTimer(t, timeout)
	return t
}

func releaseTimer(t *time.Timer) {
	stopTimer(t)
	timerPool.Put(t)
}
