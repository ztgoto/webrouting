package config

import (
	"fmt"
	"testing"
)

var content = `
# 系统配置
application:
  processes: 1  # runtime.GOMAXPROCS(processes)

upstream:
  - id: server1
    balance: random
    server: ["127.0.0.1:8080","127.0.0.1:8081","127.0.0.1:8082"]
  - id: server2
    balance: random
    server: ["127.0.0.1:8083","127.0.0.1:8084"]

# http 服务配置
http:
  server:
    - listen: 80
      host: 127.0.0.1
      location:
        - pattern: "/*"
          proxy: server1
          request: {"head1": "m1"}
          response: {"Server": "webrouting"}
    - listen: 80
      host: localhost
      location:
        - pattern: "/*"
          proxy: server2
          request: {"head1": "m1"}
          response: {"Server": "webrouting"}
    - listen: 8080
      host: localhost
      ssl: true
      cert: "/aa/bb/cc/xx.cert"
      key: "/aa/bb/cc/xx.key"
      location:
        - pattern: "/aa/*"
          root: "/aa/bb/cc"
          index: "index.html,index.htm"
          request: {"head1": "m1"}
          response: {"Server": "webrouting"}
        - pattern: "/*"
          proxy: server2
          request: {"head1": "m1"}
          response: {"Server": "webrouting"}
`

/*
go test -v github.com\ztgoto\webrouting\config
*/

func TestParseConfig(t *testing.T) {
	app,e:=ParseConfig([]byte(content))
	if e!=nil {
		fmt.Println(e)
	}
	fmt.Printf("%+v\n",app)
}

/*
go test github.com\ztgoto\webrouting\config -bench=".*"
*/

func BenchmarkParseConfig(b *testing.B){
	for i := 0; i < b.N; i++ { //use b.N for looping 
		_,e:=ParseConfig([]byte(content))
		if (e!=nil) {
			panic(e)
		}
		
    }
	
}