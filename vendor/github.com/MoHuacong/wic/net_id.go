package wic

import (
	"sync"
)

type Data struct {
	Fd int64
	Ser Server
}

type Id struct {
	m sync.Mutex
	index map[int64]Server
}

var addr *Id

func NewId() *Id {
	if addr != nil {
		return addr
	}
	
	addr = &Id{index: make(map[int64]Server)}
	return addr
}

func (id *Id) Add(fd int64, ser Server) bool {
	id.m.Lock()
	defer id.m.Unlock()
	
	delete(id.index, fd)
	id.index[fd] = ser
	
	_, ok := id.index[fd]
	
	if !ok {
		return false
	}
	return true
}

func (id *Id) Check(fd int64) Server {
	id.m.Lock()
	defer id.m.Unlock()
	
	ser, ok := id.index[fd]
	
	if !ok {
		return nil
	}
	return ser
}