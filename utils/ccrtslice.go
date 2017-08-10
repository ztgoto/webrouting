package utils

import "sync"

// ConcurrentSlice 线程安全slice
type ConcurrentSlice struct {
	lock sync.RWMutex
	data []interface{}
}
