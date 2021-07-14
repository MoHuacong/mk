package tools

import (
	"sync"
)

type Any interface{}

type SafeMap sync.Map

type AnyMap struct {
	Map map[interface{}]interface{}
	lock sync.RWMutex
}

func NewAnyMap() *AnyMap {
	return &AnyMap{Map: make(map[interface{}]interface{})}
}

func (m *AnyMap) Set(key, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Map[key] = value
}

func (m *AnyMap) Get(key interface{}) interface{} {
	return m.Map[key]
}

