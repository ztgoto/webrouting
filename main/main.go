package main

import (
	"github.com/ztgoto/webrouting/cmd"
)

// "log"
// "runtime"

// "github.com/ztgoto/webrouting/config"
// server "github.com/ztgoto/webrouting/http"

// func init() {
// 	cpuNum := runtime.NumCPU()
// 	log.Printf("CPU Num:%d\r\n", cpuNum)

// 	runtime.GOMAXPROCS(config.SC.MaxProcs)
// 	log.Printf("GOMAXPROCS:%d\r\n", config.SC.MaxProcs)
// }

func main() {
	// log.Println(config.ServerConfig)
	// server.StartServer(&config.ServerConfig)

	cmd.Execute()
}
