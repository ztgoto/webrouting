package utils

import "sync"

// ConcurrentMap 线程安全map
type ConcurrentMap struct {
	lock sync.RWMutex
	data map[interface{}]interface{}
}

// NewConcurrentMap 创建一个并发map
func NewConcurrentMap(size int) *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[interface{}]interface{}, size),
	}
}

// Put 添加键值
func (m *ConcurrentMap) Put(key, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]interface{})
	}
	m.data[key] = value
}

// Get 获取指定值
func (m *ConcurrentMap) Get(key interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.data == nil {
		return nil
	}
	v := m.data[key]
	return v
}

// Remove 移除指定值
func (m *ConcurrentMap) Remove(key interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.data != nil {
		delete(m.data, key)
	}

}

// Size 获取大小
func (m *ConcurrentMap) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.data == nil {
		return 0
	}
	return len(m.data)
}
