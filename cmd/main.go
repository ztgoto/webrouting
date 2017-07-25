package main

import "fmt"

// "log"
// "runtime"

// "github.com/ztgoto/webrouting/config"
// server "github.com/ztgoto/webrouting/http"

const banner string = `

                  ___.                                   __  .__                
    __  _  __ ____\_ |__           _______  ____  __ ___/  |_|__| ____    ____  
    \ \/ \/ // __ \| __ \   ______ \_  __ \/  _ \|  |  \   __\  |/    \  / ___\ 
     \     /\  ___/| \_\ \ /_____/  |  | \(  <_> )  |  /|  | |  |   |  \/ /_/  >
      \/\_/  \___  >___  /          |__|   \____/|____/ |__| |__|___|  /\___  / 
                 \/    \/                                            \//_____/  

`

// func init() {
// 	cpuNum := runtime.NumCPU()
// 	log.Printf("CPU Num:%d\r\n", cpuNum)

// 	runtime.GOMAXPROCS(config.SC.MaxProcs)
// 	log.Printf("GOMAXPROCS:%d\r\n", config.SC.MaxProcs)
// }

func main() {
	// log.Println(config.ServerConfig)
	// server.StartServer(&config.ServerConfig)
	fmt.Println(banner)
}
