package wic

import (
	"strconv"
)

type ServerRealize struct {
	port uint
	addr string
	name []string
	callback Logic
}

func (sr *ServerRealize) GetName() []string {
	return sr.name
}

func (sr *ServerRealize) IsName(name string) bool {
	for _, _name := range sr.name {
		if name == _name {
			return true
		}
	}
	return false
}

func (sr *ServerRealize) SetName(name []string) bool {
	sr.name = name
	return true
}

func (sr *ServerRealize) Init() {
	
}

func (sr *ServerRealize) SetLogic(logic Logic) bool {
	sr.callback = logic
	if sr.callback == nil {
		return false
	}
	return true
}

func (sr *ServerRealize) GetRaw(name string) interface{} {
	return sr
}

func (sr *ServerRealize) GetPort() uint {
	return sr.port
}

func (sr *ServerRealize) SetIpPort(ip string, port int) {
	sr.port = uint(port)
	sr.addr = ip + ":" + strconv.Itoa(port)
}

func (sr *ServerRealize) Close(fd *Fd) bool {
	return true
}

func (sr *ServerRealize) callBack(typ string, ser Server, value interface{}, data string) {
	if typ == "init" {
		sr.callback.Init(ser)
	}
	
	if typ == "end" {
		sr.callback.End(ser)
	}
	
	var fd *Fd = &Fd{Ser: ser, Value: value}
	
	if typ == "connect" {
		sr.callback.Connect(ser, fd)
	}
	
	if typ == "receive" {
		sr.callback.Receive(ser, fd, data)
	}
	
	if typ == "close" {
		sr.callback.Closes(ser, fd)
	}
}