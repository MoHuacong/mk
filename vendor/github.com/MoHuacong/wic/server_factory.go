package wic

import (
	"strings"
	"strconv"
)

var sf *ServerFactory

type ServerFactory struct {
	fun []ServerFunc
}

func NewServerFactory() *ServerFactory {
	//return new(ServerFactory)
	return &ServerFactory{fun: make([]ServerFunc, 1)}
}

func (sf *ServerFactory) Create(name string) Server {
	
	for _, f := range sf.fun {
		if f == nil { continue }
		
		ser := f()
		for _, _name := range ser.GetName() {
			if name == _name {
				return ser
			}
		}
	}
	
	return nil
}

func (sf *ServerFactory) Router(url string) (Server, int) {
	var port int
	var ip string
	var err error
	var ser Server
	arr := strings.Split(url, ":")
	
	if arr == nil || len(arr) == 0 {
		return nil, 1
	}
	
	if ser = sf.Create(arr[0]); ser == nil {
		return nil, 2
	}
	
	ip = arr[1]
	if ip == "" {
		ip = "0.0.0.0"
	}
	
	port, err = strconv.Atoi(arr[2])
	if err != nil || port == 0 {
		port = 8080
	}
	
	ser.SetIpPort(ip, port)
	return ser, 0
}

func (sf *ServerFactory) On(url string, logic Logic) (Server, int) {
	ser, err := sf.Router(url)
	if err != 0 {
		return ser, err
	}
	
	logic.News()
	logic.AddServer(ser)
	ser.SetLogic(logic)
	return ser, err
}


func GetServerFactory() *ServerFactory {
	if sf == nil {
		sf = NewServerFactory()
	}
	
	return sf
}

func AddServerFactory(fun ServerFunc) bool {
	if sf == nil {
		sf = NewServerFactory()
	}
	
	if len(sf.fun) >= cap(sf.fun) - 1 {
		tmp := make([]ServerFunc, len(sf.fun), (cap(sf.fun))*2)
		copy(tmp, sf.fun)
		sf.fun = tmp
	}
	
	sf.fun = append(sf.fun, fun)
	return true
}

var registerServerFactory func (ServerFunc) bool = AddServerFactory

func init() {
	/*
	AddServerFactory(NewMq)
	AddServerFactory(NewTcp)
	AddServerFactory(NewWeb)
	AddServerFactory(NewHttp)
	AddServerFactory(NewWebSocket)
	*/
}