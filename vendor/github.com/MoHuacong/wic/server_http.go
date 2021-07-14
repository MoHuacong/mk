package wic

import (
	//"strings"
	"net/http"
)

type Http struct {
	ServerRealize
	
	request *http.Request
	response http.ResponseWriter
}

func init() {
	registerServerFactory(NewHttp)
}

func NewHttp() Server {
	h := new(Http)
	h.SetName([]string{"http", "https"})
	return h
}

func (h *Http) Run(blocking bool) error {
	if blocking {
		return http.ListenAndServe(h.addr, h)
	}
	
	go func() error {
		return http.ListenAndServe(h.addr, h)
	}()
	
	return nil
}

func (h *Http) Read(fd *Fd) string {
	return ""
}

func (h *Http) Write(fd *Fd, data string) bool {
	fd.Value.(http.ResponseWriter).Write([]byte(data))
	return true
}

func (h *Http) Send(fd *Fd, data string) bool {
	return h.Write(fd, data)
}

func (h *Http) GetRaw(name string) interface{} {
	if name == "r" || name == "request" {
		return h.request
	} else if name == "w" || name == "response" {
		return h.response
	}
	
	return nil
}

func (h *Http) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	h.request = r
	h.response = w
	
	var ser Http = *h
	var fd *Fd = NewFd(&ser, w)
	
	h.callback.Connect(h, fd)
	
	h.callback.Receive(h, fd, r.URL.Path)
	
	h.callback.Closes(h, fd)
}