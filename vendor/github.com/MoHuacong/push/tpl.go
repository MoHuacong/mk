package push

import (
	"github.com/MoHuacong/wic"
	"net/http"
)

var title string = "注册-Moid"

type Tpl struct {
	mk *Mk
}

func (tpl *Tpl) Index(w http.ResponseWriter, r *http.Request) wic.CallBackMap {
	tp := tpl.mk.web.Template(w)
	tp.Assign("title", title)

	if r.FormValue("ok") == "" {
		return tp.GetContext()
	}

	user := r.FormValue("user")
	pass := r.FormValue("pass")
	url := r.FormValue("url")

	tp.Assign("log", "true")

	if user == "" || pass == "" || url == "" {
		tp.Assign("logtitle", "注册失败")
		tp.Assign("logcontent", "user or pass or url")
		return tp.GetContext()
	}

	if tpl.mk.api.Register(user, pass, url) != 0 {
		tp.Assign("logtitle", "注册失败")
		tp.Assign("logcontent", tpl.mk.api.Error())
		return tp.GetContext()
	}

	tp.Assign("logtitle", "注册成功")
	tp.Assign("logcontent", tpl.mk.api.String())
	return tp.GetContext()
}

func (tpl *Tpl) Doc(w http.ResponseWriter, r *http.Request) wic.CallBackMap {
	tp := tpl.mk.web.Template(w)
	tp.Assign("title", "文档")
	tp.Assign("url", "http://"+r.Host+"/")
	return tp.GetContext()
}