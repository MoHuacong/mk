package main

import (
	"fmt"
)
import "github.com/MoHuacong/wic"

type A struct {
	wic.LogicRealize
}

func (a *A) Connect(ser wic.Server, fd *wic.Fd) bool {
	fmt.Printf("%d\n", fd)
	return true
}

func (a *A) Receive(ser wic.Server, fd *wic.Fd, data string) bool {
	fmt.Printf("%d\n", fd)
	fmt.Println(data)
	ser.Send(fd, "Moid-2333")
	return true
}

func (a *A) Close(ser wic.Server, fd *wic.Fd) bool {
	fmt.Printf("%d\n", fd)
	fmt.Println("close")
	return true
}

func main() {
	a := new(A)
	sf := wic.NewServerFactory()
	
	
	tcp, err := sf.On("tcp:0.0.0.0:8088", a)
	
	if err != 0 {
		return
	}
	
	tcp.Run(false)
	
	http, err := sf.On("http:0.0.0.0:8081", a)
	
	if err != 0 {
		return
	}
	
	http.Run(false)
	
	fmt.Println("start")
	
	ws, err := sf.On("websocket:0.0.0.0:8080", a)
	
	if err != 0 {
		return
	}
	
	ws.Run(true)
}