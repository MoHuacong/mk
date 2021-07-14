package tools

import (
	"sync"
)

type Count struct {
	v uint64
	lock sync.RWMutex
}

func (c *Count) Add() uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.v++
	return c.v
}

func (c *Count) Value() uint64 {
	return c.v
}

type IdCount struct {
	any *AnyMap
	count *Count
}

func NewIdCount() *IdCount {
	return &IdCount{any: NewAnyMap(), count: &Count{v:0}}
}

func (i *IdCount) Add(v interface{}) uint64 {
	id := i.count.Add()
	i.any.Set(id, v)
	return id
}

func (i *IdCount) Get(id uint64) interface{} {
	return i.any.Get(id)
}

func (i *IdCount) Check(k interface{}) (interface{}, bool) {
	v := i.any.Get(k)
	if v == nil {
		return nil, false
	}
	return v, true
}