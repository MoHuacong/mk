package wic

import (
	"encoding/json"
	"errors"
	"github.com/MoHuacong/wic/tools"
	"github.com/flosch/pongo2"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type CallBackMap map[string]interface{}
type RouterCallBack func(http.ResponseWriter, *http.Request)

/* 模板 */
type Template struct {
	dir string
	web *Web
	data pongo2.Context
	w http.ResponseWriter
}

/* 赋值 */
func (tpl *Template) Assign(name string, v interface{}) bool {
	tpl.data[name] = v
	return (tpl.data[name] == v)
}

/* 显示/渲染 */
func (tpl *Template) Display(name string) bool {
	path := tpl.dir + "/" + name
	t, err := pongo2.FromFile(path)
	
	if t == nil || err != nil { return tpl.Error(err) }
	
	err = t.ExecuteWriter(tpl.data, tpl.w)
	
	_, types := tpl.web.mime(name)
	tpl.w.Header().Set("Content-Type", types)
	
	if err != nil { return tpl.Error(err)}
	return true
}

/* 设置模板路径 */
func (tpl *Template) SetTemplateDir(dir string) {
	tpl.dir = dir
}

func (tpl *Template) GetContext() CallBackMap {
	data := make(CallBackMap)
	for k, v := range tpl.data {
		data[k] = v
	}
	return data
}

/* 错误输出 */
func (tpl *Template) Error(err error) bool {
	http.Error(tpl.w, err.Error(), http.StatusInternalServerError)
	return true
}

type Web struct {
	Http
	ait bool
	dir string
	tpl_dir string
	Router map[string]interface{}
	TplRouter map[string]interface{}
}

func init() {
	/* 注册服务 */
	registerServerFactory(NewWeb)
}

func NewWeb() Server {
	web := new(Web)
	web.SetName([]string{"goweb", "web"})
	dir, _ := os.Getwd()
	web.dir = dir
	web.Router = make(map[string]interface{})
	web.TplRouter = make(map[string]interface{})
	return web
}

/* 开始运行 */
func (web *Web) Run(blocking bool) error {
	web.callback.Init(web)
	if blocking {
		return http.ListenAndServe(web.addr, web)
	}
	
	go func() error {
		return http.ListenAndServe(web.addr, web)
	}()
	
	return nil
}

/* net/http包的接口实现 */
func (web *Web) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	web.request = r
	web.response = w

	if web.templateRouters(w, r) {
		return
	}

	if web.routers(w, r) {
		return
	}
	
	if web.header(w, r) {
		return
	}
	
	var ser Web = *web
	var fd *Fd = NewFd(&ser, w)
	
	web.callback.Connect(web, fd)
	
	web.callback.Receive(web, fd, r.URL.Path)
	
	web.callback.Closes(web, fd)
}

func (web *Web) GetDir() string {
	return web.dir
}

func (web *Web) SetDir(dir string) {
	web.dir = dir
}

/* 获取模板 */
func (web *Web) Template(w http.ResponseWriter) *Template {
	return &Template{web.tpl_dir, web, make(pongo2.Context), w}
}

/* 设置模板路径 */
func (web *Web) SetTemplateDir(dir string) {
	web.tpl_dir = web.dir + "/" + dir
}

/* 设置路由器状态 */
func (web *Web) SetRouterAutomatic(ok bool) {
	web.ait = ok
}

/* 添加路由处理函数 */
func (web *Web) AddRouter(reg string, lsp interface{}) bool {
	web.Router[reg] = lsp
	return true
}

/* 添加模板路由处理函数 */
func (web *Web) AddTemplateRouter(reg string, lsp interface{}) bool {
	web.TplRouter[reg] = lsp
	return true
}

/* 调用并运行路由处理函数 */
func (web *Web) routers(w http.ResponseWriter, r *http.Request) bool {
	for reg, lsp := range web.Router {
		typ := reflect.TypeOf(lsp).Kind().String()
		
		// 函数路由器
		if typ == "func" {
			if ok, _ := regexp.MatchString(reg, r.URL.Path); ok {
				fu := reflect.ValueOf(lsp)
				ret := fu.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})
				if len(ret) > 0 {
					s := tools.ValueToSliceInterface(ret)
					if len(s) == 0 || s[0] == nil { return true }
					js, err := json.Marshal(s)
					if err != nil { continue }
					w.Write(js)
					return true
				}
				continue
			}
		}
		// 结构体路由器
		if typ == "struct" || typ == "ptr" {
			var url string = ""
			regx, _ := regexp.Compile(reg)
			params := regx.FindStringSubmatch(r.URL.Path)
			
			if params != nil && len(params) != 0 {
				url = params[len(params) - 1]
			}
			
			if len(params) >= 3 {
				url = ""
				for i := 0; i < len(params); i++ {
					if i == len(params) - 1 {
						url = url + params[i]
					}
					
					if i != 0 && i != len(params) - 1{
						url = url + params[i] + "/"
					}
				}
			}
			url = tools.StrToUpper(url)
			ret, err := web.automatic(url, "", lsp, []interface{}{w, r})
			if err != nil { return false }
			if ret == nil || len(ret) == 0 { return true }

			s := tools.ValueToSliceInterface(ret)
			if len(s) == 0 || s[0] == nil { return true }
			js, err := json.Marshal(s)
			if err != nil { continue }
			w.Write(js)
			return true
		}
	}
	return false
}

func (web *Web) templateRouters(w http.ResponseWriter, r *http.Request) bool {
	var path string = web.tpl_dir + r.URL.Path
	for reg, lsp := range web.TplRouter {
		typ := reflect.TypeOf(lsp).Kind().String()

		// 函数路由器
		if typ == "func" {
			if ok, _ := regexp.MatchString(reg, r.URL.Path); ok {
				fu := reflect.ValueOf(lsp)
				ret := fu.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)})

				if ret == nil || len(ret) == 0 { return true }
				s := tools.ValueToSliceInterface(ret)
				if len(s) == 0 || s[0] == nil { return true }

				m := s[0].(CallBackMap)
				if m == nil { return true }

				if name, types := web.mime(path); name == "" || types == "" {
					path += "index.html"
				}

				web.TplOut(w, path, m)
				return true
			}
		}
		// 结构体路由器
		if typ == "struct" || typ == "ptr" {
			var url string = ""
			regx, _ := regexp.Compile(reg)
			params := regx.FindStringSubmatch(r.URL.Path)

			if params != nil && len(params) != 0 {
				url = params[len(params) - 1]
			}

			if len(params) >= 3 {
				url = ""
				for i := 0; i < len(params); i++ {
					if i == len(params) - 1 {
						url = url + params[i]
					}

					if i != 0 && i != len(params) - 1{
						url = url + params[i] + "/"
					}
				}
			}
			url = tools.StrToUpper(url)

			ret, err := web.automatic(url, "", lsp, []interface{}{w, r})
			if err != nil { return false }
			if ret == nil || len(ret) == 0 { return true }

			s := tools.ValueToSliceInterface(ret)
			if len(s) == 0 || s[0] == nil { return true }

			m := s[0].(CallBackMap)
			if m == nil { return true }
			web.TplOut(w, path, m)
			return true
		}
	}
	return false
}

func (web *Web) TplOut(w http.ResponseWriter, path string, data map[string]interface{}) bool {
	_ , types := web.mime(path)

	tpl, err := pongo2.FromFile(path)

	if tpl == nil || err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	err = tpl.ExecuteWriter(data, w)

	w.Header().Set("Content-Type", types)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

/* 反射调用 */
func (web *Web) reflectCall(i interface{}, method string, args []reflect.Value) ([]reflect.Value, error) {
	var zero reflect.Value
	
	if i == nil || method == "" {
		return nil, errors.New("object in method no")
	}
	
	v := reflect.ValueOf(i)
	mv := v.MethodByName(method)
	
	if mv == zero {
		return nil, errors.New("method nonono")
	}
	
	if mv.Type().NumIn() < len(args) {
		return nil, errors.New("args no")
	}
	
	return mv.Call(args), nil
}

/* 路由url处理并反射调用 */
func (web *Web) automatic(url, head string, object interface{}, args []interface{}) ([]reflect.Value, error) {

	method := strings.Replace(url, "/", "_", -1)
	method = head + method
	
	i := strings.Index(method, ".")
	
	if i != -1 {
		method = method[:i]
	}
	
	var list []reflect.Value = make([]reflect.Value, len(args))
	
	for k, v := range args {
		list[k] = reflect.ValueOf(v)
	}
	
	return web.reflectCall(object, method, list)
}

/* mime类型获得 */
func (web *Web) mime(url string) (string, string) {
	arr := strings.Split(url, ".")
	mime.AddExtensionType(".txt", "text/plain")
	types := mime.TypeByExtension("." + arr[len(arr) - 1])
	return arr[len(arr) - 1], types
}

/* 预先处理请求头 */
func (web *Web) header(w http.ResponseWriter, r *http.Request) bool {
	name , types := web.mime(r.URL.Path)
	if types == "" && name != "go" {
		return false
	}
	
	var path string = web.dir + r.URL.Path

	_, err := os.Lstat(path)
	if os.IsNotExist(err) {
		w.Header().Set("content-type", "text/html")
		w.Write([]byte("<h1>404-Moid</h1>"))
		return true
	}

	f, err := ioutil.ReadFile(path)
	w.Write(f)
	return true
}