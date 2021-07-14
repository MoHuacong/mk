package wic

type ServerMap map[string]map[uint]Server

/* 逻辑业务接口 */
type Logic interface {
	IO
	News()
	Init(ser Server) bool
	End(ser Server) bool
	
	SetServer(ser Server) bool
	AddServer(ser Server) bool
	GetServer(name string) map[uint]Server
	
	GetServerList() ServerMap
	
	Connect(ser Server, fd *Fd) bool
	Receive(ser Server, fd *Fd, data string) bool
	Closes(ser Server, fd *Fd) bool
}