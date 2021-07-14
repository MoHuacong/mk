package push

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	/*
		"io/ioutil"
		"net/http"
	*/
	"encoding/json"
	"github.com/MoHuacong/wic"
	"github.com/MoHuacong/wic/tools"
	"github.com/go-redis/redis"
)

type Id map[int]wic.Fd

/* 包容任何类型的Map */
type Any map[interface{}]interface{}

type Mk struct {
	wic.LogicRealize
	
	/* service */
	mq *wic.Mq
	web *wic.Web
	port map[string]int
	
	/* data */
	api *Api
	id *tools.ListMap
	redis *redis.Client
}

func OpenMk(args ...int) *Mk {
	var ws, tcp, web int
	if len(args) >= 3 {
		ws = args[0]
		tcp = args[1]
		web = args[2]
	}

	if len(args) == 1 {
		ws = args[0]
		tcp = ws + 1
		web = tcp + 1
	}

	mk := new(Mk)
	mk.port = make(map[string]int)
	mk.port["ws"] = ws
	mk.port["tcp"] = tcp
	mk.port["web"] = web
	
	var err int
	var ser wic.Server
	sf := wic.GetServerFactory()
	
	if ser, err = sf.On("ws:0.0.0.0:" + strconv.Itoa(ws), mk); err != 0 {
		return nil
	}
	ser.Run(false)
	
	if ser, err = sf.On("tcp:0.0.0.0:" + strconv.Itoa(tcp), mk); err != 0 {
		return nil
	}
	ser.Run(false)
	
	if ser, err = sf.On("mq:0.0.0.0:" + strconv.Itoa(web + 1), mk); err != 0 {
		return nil
	}
	ser.Run(false)
	
	if ser, err = sf.On("web:0.0.0.0:" + strconv.Itoa(web), mk); err != 0 {
		return nil
	}
	ser.Run(true)
	
	return mk
}

/* 服务初始化 */
func (mk *Mk) Init(ser wic.Server) bool {
	var web *wic.Web
	
	if ser.IsName("mq") {
		mk.mq = ser.(*wic.Mq)
		mk.mq.SetTopic(3, 240)
	}
	
	if !ser.IsName("web") { return false }
	
	web = ser.(*wic.Web)
	
	web.SetRouterAutomatic(true)
	
	webObject := new(Web)
	webObject.mk = mk

	web.AddRouter("/(.*).go$", webObject)
	web.AddRouter("^/$", webObject.Index)

	tpl := new(Tpl)
	tpl.mk = mk

	web.AddTemplateRouter("/(.*).html$", tpl)
	web.AddTemplateRouter("^/$", tpl.Index)

	web.SetTemplateDir("/assets/html")
	
	mk.web = web
	
	mk.api = NewApi()
	mk.api.Open("mk", "123456", "localhost", "3306", "mk", "utf8")

	mk.redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})


	mk.id = tools.NewListMap()

	mk.mq.Topic("连接成功").CallBack(mk.mq_Connect, 0)
	mk.mq.Topic("发送数据").CallBack(mk.mq_Receive, 0, 10)
	mk.mq.Topic("断开连接").CallBack(mk.mq_Closes, 0)

	return true
}

/* 新连接处理(定时器) */
func (mk *Mk) Connect(ser wic.Server, fd *wic.Fd) bool {
	if !ser.IsName("ws") && !ser.IsName("tcp") {
		return true
	}
	
	ticker := time.NewTicker(time.Second * 5)
	
	go func() {
		for range ticker.C {
			ticker.Stop()
			if fd != nil && !fd.Validate {
				mk.Close(fd)
			}
		}
	}()
	return true
}

/* 接收数据处理 */
func (mk *Mk) Receive(ser wic.Server, fd *wic.Fd, data string) bool {
	if !ser.IsName("ws") && !ser.IsName("tcp") {
		return true
	}
	/* 如果fd已验证过 */
	if fd != nil && fd.Validate {
		mk.mq.Topic("发送数据").Send(mk.Map("fd", fd, "data", data))
		return true
	}
	
	// json
	jc := &JsonClient{}
	if err := json.Unmarshal([]byte(data), jc); err != nil {
		return false
	}
	
	/* 判断地址存在否 */
	user := new(User)
	if mk.id.IdKey(fd.Id) != nil && mk.api.IsUrl(user, jc.Url) != 0 {
		mk.Close(fd)
		return false
	}
	
	fd.Validate = true
	mk.AddId(jc.Url, fd.Id)
	mk.mq.Topic("连接成功").Send(mk.Map("url", jc.Url, "fd", fd))
	return true
}

func (mk *Mk) Closes(ser wic.Server, fd *wic.Fd) bool {
	if !ser.IsName("ws") && !ser.IsName("tcp") {
		return true
	}
	mk.mq.Topic("断开连接").Send(mk.Map("fd", fd))
	return true
}

func (mk *Mk) AddId(url string, id uint64) bool {
	return mk.id.Set(id, url)
}

/* 格式并初始化 */
func (mk *Mk) Map(args ...interface{}) Any {
	m := make(Any)
	for i := 0; i < len(args); i+=2 {
		m[args[i]] = args[i+1]
	}
	return m
}

func (mk *Mk) Format(typ string, fd *wic.Fd, client, ata string) url.Values {
	urlValues := url.Values{}
	urlValues.Add("type", typ)
	urlValues.Add("client", client)
	urlValues.Add("id", fd.IdString())
	if ata != "" {
		urlValues.Add("data", ata)
	}
	return urlValues
}

func (mk *Mk) Http(url string, value url.Values) string {
	resp, err := http.PostForm(url, value)
	if err != nil { return "" }
	defer resp.	Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { return "" }
	return string(body)
}

func (mk *Mk) Id_Send(id uint64, data string) bool {
	fd := wic.GetFd(id)
	if fd == nil { return false }
	return fd.Ser.Send(fd, data)
}

func (mk *Mk) Mk_Send(data *JsonResponse) (int, int) {
	var ok, err int
	for _, id := range data.Id.List {
		tmp := mk.id.IdKey(id)
		if tmp == nil { err++; continue }
		url := tmp.(string)
		if url == data.Url {
			if mk.Id_Send(id, data.Id.Data) {
				ok++
				continue
			}
		}
		err++
	}
	return ok, err
}

func (mk *Mk) Mk_SendAll(data *JsonResponse) (int, int) {
	var ok, err int
	for _id, _ := range mk.id.UrlKeyList(data.Url) {
		id := _id.(uint64)
		tmp := mk.id.IdKey(id)
		if tmp == nil { err++; continue }
		url := tmp.(string)
		if url == data.	Url {
			if mk.Id_Send(id, data.Id.Data) {
				ok++
				continue
			}
		}
		err++
	}
	return ok, err
}

func (mk *Mk) JsonApi(data string) bool {
	jr := new(JsonResponse)
	if json.Unmarshal([]byte(data), jr) != nil {
		mk.api.formatExit(1, "JSON_NO")
		return false
	}
	if mk.api.Is(jr.User, jr.Pass, jr.Url) != 0 { return false }
	if jr.Type == "send" {
		ok, err := mk.Mk_Send(jr)
		str := "{\"ok\":" + strconv.Itoa(ok) + ",\"err\":" + strconv.Itoa(err) + "}"
		mk.api.formatExit(0, str)
	} else if jr.Type == "sendAll" {
		ok, err := mk.Mk_SendAll(jr)
		str := "{\"ok\":" + strconv.Itoa(ok) + ",\"err\":" + strconv.Itoa(err) + "}"
		mk.api.formatExit(0, str)
	}
	return true
}

func (mk *Mk) mq_Connect(_id int, i interface{}) bool {
	data := i.(Any)
	fd := data["fd"].(*wic.Fd)
	url := data["url"].(string)
	ret := mk.Http(url, mk.Format("connect", fd, fd.Ser.GetName()[0], ""))
	return mk.JsonApi(ret)
}

func (mk *Mk) mq_Receive(_id int, i interface{}) bool {
	data := i.(Any)
	fd := data["fd"].(*wic.Fd)
	datas := data["data"].(string)
	tmp := mk.id.IdKey(fd.Id)
	if tmp == nil { return false }
	url := tmp.(string)
	ret := mk.Http(url, mk.Format("receive", fd, fd.Ser.GetName()[0], datas))
	return mk.JsonApi(ret)
}

func (mk *Mk) mq_Closes(_id int, i interface{}) bool {
	data := i.(Any)
	fd := data["fd"].(*wic.Fd)
	tmp := mk.id.IdKey(fd.Id)
	if tmp == nil { return false }
	url := tmp.(string)
	mk.id.Delete(fd.Id)
	ret := mk.Http(url, mk.Format("close", fd, fd.Ser.GetName()[0], ""))
	return mk.JsonApi(ret)
}