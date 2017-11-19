package pool

import (
	"sync"
	"bufio"
	"net"
)

type PersistConn struct {
  	conn      net.Conn
  	br        *bufio.Reader       // from conn
	bw        *bufio.Writer       // to conn
	closed    bool
	mu        sync.Mutex
}

func (pc *PersistConn) Read(p []byte) (n int, err error) {
  	n, err = pc.conn.Read(p)
  	if err != nil {
		pc.closed = true
		pc.conn.Close()  
  	}
  	return
}
