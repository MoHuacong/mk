package push

import (
	"fmt"
	"net/http"
	"strconv"
)

type Web struct {
	mk *Mk
}

func (web *Web) Index(w http.ResponseWriter, r *http.Request) interface{} {
	fmt.Println("index---")
	tpl := web.mk.web.Template(w)
	tpl.Assign("title", "Moid")
	if !tpl.Display("index.html") {
		return "Moid2333"
	}
	return nil
}

func (web *Web) Register(w http.ResponseWriter, r *http.Request) string {
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	url := r.FormValue("url")
	ret := web.mk.api.Register(user, pass, url)
	return strconv.FormatUint(uint64(ret), 10) + "_" + web.mk.api.String()
}

func (web *Web) Login(w http.ResponseWriter, r *http.Request) string {
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	ret := web.mk.api.Login(user, pass)
	return strconv.FormatUint(uint64(ret), 10) + "_" + web.mk.api.String()
}

func (web *Web) Api(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	web.mk.JsonApi(data)
	w.Write([]byte(web.mk.api.String()))
}

// 127.0.0.1:8032/api.go?data={"type":"send","user":"demo","pass":"demo","url":"http://demo.moid.red:81","id":{"list":[1,2,3],"data":"2333"}}
/*
{
	"type":"send",
	"user":"abc",
	"pass":"abc",
	"url":"http://www.moid.red:81",
	"id":{
		"list":[1,2,3],
		"data":"2333"
	}
}
*/