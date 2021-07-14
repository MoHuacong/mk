package wic

import (
	"net"
	"github.com/MoHuacong/wic/tools"
)

type Tcp struct {
	ServerRealize
	server_tcp *net.TCPListener
}

func init() {
	/* 注册服务 */
	registerServerFactory(NewTcp)
}

func NewTcp() Server {
	tcp := new(Tcp)
	tcp.SetName([]string{"tcp", "TCP"})
	return tcp
}

func (tcp *Tcp) Run(blocking bool) error {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", tcp.addr)
	tcp.server_tcp, _ = net.ListenTCP("tcp", tcpAddr)
	//defer tcp.server_tcp.Close()
	//tcp.callback.Init(tcp)
	tcp.callBack("int", tcp, nil, "")
	
	if blocking {
		tcp.upper()
	}
	
	go tcp.upper()
	
	return nil
}

func (tcp *Tcp) Read(fd *Fd) string {
	return ""
}

func (tcp *Tcp) Write(fd *Fd, data string) bool {
	fd.Value.(*net.TCPConn).Write([]byte(data))
	return true
}

func (tcp *Tcp) Send(fd *Fd, data string) bool {
	return tcp.Write(fd, data)
}

func (tcp *Tcp) Close(fd *Fd) bool {
	fd.Value.(*net.TCPConn).Close()
	return true
}

func (tcp *Tcp) upper() {
	for {
		conn, err := tcp.server_tcp.AcceptTCP()
		if err != nil {
			continue
		}
		
		go func() {
			var fd *Fd = NewFd(tcp, conn)
			tcp.callback.Connect(tcp, fd)
			
			for {
				buf := make([]byte, 1024 * 2)
				
				if _, err := conn.Read(buf); err != nil {
					tcp.callback.Closes(tcp, fd)
					conn.Close()
					return
				}
				tcp.callback.Receive(tcp, fd, string(tools.Lengths(buf)))
			}
		}()
	}
}