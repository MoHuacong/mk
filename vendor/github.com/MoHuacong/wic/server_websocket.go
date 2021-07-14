package wic

import (
	"net/http"
	"strings"
	"golang.org/x/net/websocket"
)

type WebSocket struct {
	ServerRealize
}

func init() {
	registerServerFactory(NewWebSocket)
}

func NewWebSocket() Server {
	ws := new(WebSocket)
	ws.SetName([]string{"ws", "websocket"})
	return ws
}

func (ws *WebSocket) Run(blocking bool) error {
	http.Handle("/", websocket.Handler(ws.upper))
	
	if blocking {
		return http.ListenAndServe(ws.addr, nil)
	}
	
	if !blocking {
		go func() error {
			return http.ListenAndServe(ws.addr, nil)
		}()
	}
	
	ws.callback.Init(ws)
	
	return nil
}

func (ws *WebSocket) Read(fd *Fd) string {
	return ""
}

func (ws *WebSocket) Write(fd *Fd, data string) bool {
	ws_conn := fd.Value.(*websocket.Conn)
	if err := websocket.Message.Send(ws_conn, strings.ToUpper(data)); err != nil {
			return false
		}

	return true
}

func (ws *WebSocket) Send(fd *Fd, data string) bool {
	return ws.Write(fd, data)
}

func (ws *WebSocket) upper(_ws *websocket.Conn) {
	var err error
	var fd *Fd = NewFd(ws, _ws)
	ws.callback.Connect(ws, fd)
	for {
		var data string
		if err = websocket.Message.Receive(_ws, &data); err != nil {
			ws.callback.Closes(ws, fd)
			break
		}
		ws.callback.Receive(ws, fd, data)
	}

}