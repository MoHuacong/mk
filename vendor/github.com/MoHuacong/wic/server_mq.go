package wic

import (
	"fmt"
	"github.com/MoHuacong/wic/tools"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

type Mq MessageQueue

type Message interface{}
type Queue chan Message

/* 往队列添加消息 */
func (queue Queue) Add(msg Message) bool {
	select {
		case queue <- msg:
			return true
		default:
			return false
	}
	return false
}

/* 从队列读取消息 */
func (queue Queue) Read(len int) []Message {
	data := make([]Message, len)
	for k, _ := range data {
		select {
		case data[k] = <- queue:
			continue
		default:
			return data
		}
	}
	return data
}

func (queue Queue) Readk(len int) []Message {
	data := make([]Message, len)
	for k, _ := range data {
		v, ok := <- queue
		if ok {
			data[k] = v
		}
	}
	return data
}


type CallBackTopic func(id int, data interface{}) bool

/* 主题/订阅 */
type Topic struct {
	name string
	queue []Queue
	callBack []CallBackTopic
}

/* 获得主题名称 */
func (topic *Topic) GetName() string {
	return topic.name
}

/* 设置主题名称 */
func (topic *Topic) SetName(name string) {
	topic.name = name
}

func (topic *Topic) GetQueue(id int) Queue {
	return topic.queue[id]
}

/* 创建并初始化消息队列 */
func (topic *Topic) NewQueue(number uint, size uint) {
	topic.queue = make([]Queue, number)
	topic.callBack = make([]CallBackTopic, number)
	
	for k, _ := range topic.queue {
		topic.queue[k] = make(Queue, size)
	}
}

/* 添加队列 */
func (topic *Topic) AddQueue(number int, size uint) int {
	s := make([]Queue, len(topic.queue), (cap(topic.queue)) + number)
	copy(s, topic.queue)
	topic.queue = s
	length := len(topic.queue)
	topic.queue[length] = make(Queue, size)

	sc := make([]CallBackTopic, len(topic.callBack), cap(topic.callBack) + number)
	copy(sc, topic.callBack)
	topic.callBack = sc
	return length
}

/* 发送消息(生产者) */
func (topic *Topic) Send(msg Message) (int, int) {
	var ok, err int = 0, 0
	for _, queue := range topic.queue {
		//queue <- msg
		if queue.Add(msg) {
			ok++
			continue
		}
		err++
	}
	return ok, err
}

/* 读取消息(消费者) */
func (topic *Topic) Read(id, len int) []Message {
	return topic.queue[id].Readk(len)
}

func (topic *Topic) CallBack(callback CallBackTopic, args ...int) bool {
	var id, l int = 0, 0
	if topic.callBack[id] != nil { return true }
	topic.callBack[id] = callback
	if len(args) >= 0 { id = 0; l = 0  }
	if len(args) == 1 { id = args[0]; l = 1 }
	if len(args) >= 2 { id = args[0]; l = args[1] }
	go func() {
		for {
			if topic.callBack[id] == nil { return }
			for _, msg := range topic.Read(id, l) {
				if msg == nil { continue }
				go topic.callBack[id](id, msg)
			}
		}
	}()
	return true
}

func (topic *Topic) SetCallBack(id int, callbak CallBackTopic) bool {
	if len(topic.callBack) <= id { return false }
	topic.callBack[id] = callbak
	if topic.callBack[id] == nil { return false }
	return true
}

func (topic *Topic) SetPush(id int, _url string) {
	topic.callBack[id] = func(id int, data interface{}) bool {
		var err error
		var buf []byte
		value := url.Values{}

		typ, str := tools.InterfaceToStr(data)
		value.Add("type", typ)
		value.Add("data", str)

		resp, _ := http.PostForm(_url, value)
		defer resp.Body.Close()

		if buf, err = ioutil.ReadAll(resp.Body); err != nil {
			return true
		}

		ret := string(buf)
		if ret != "true" { return false }
		return true
	}

	go func() {
		for {
			if topic.callBack[id] == nil { return }
			for _, msg := range topic.Read(id, 1) {
				if msg == nil { continue }
				go topic.callBack[id](id, msg)
			}
		}
	}()
}




/* 消息队列 */
type MessageQueue struct {
	Tcp
	size uint
	list []*Topic
	number uint
}

func init() {
	registerServerFactory(NewMq)
}

func NewMq() Server {
	mq := new(Mq)
	mq.number = 3
	mq.size = 20
	mq.list = make([]*Topic, 20)
	mq.SetName([]string{"mq", "MessageQueue"})
	return mq
}

func (mq *Mq) Run(blocking bool) error {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", mq.addr)
	mq.server_tcp, _ = net.ListenTCP("tcp", tcpAddr)
	//defer tcp.server_tcp.Close()
	mq.callback.Init(mq)
	
	if blocking {
		mq.upper()
	}
	
	go mq.upper()
	
	return nil
}

func (mq *Mq) upper() {
	fmt.Println("moid")
}

/* 创建主题 */
func (mq *Mq) newTopic(name string, number uint, size uint) *Topic {
	topic := new(Topic)
	topic.SetName(name)
	topic.NewQueue(number, size)
	return topic
}

/* 判断主题是否存在 */
func (mq *Mq) isTopic(name string) (bool, *Topic) {
	for _, topic := range mq.list {
		if topic == nil { continue }
		if topic.GetName() == name {
			return true, topic
		}
	}
	return false, nil
}

/* 设置主题参数 */
func (mq *Mq) SetTopic(number uint, size uint) {
	mq.size = size
	mq.number = number
}

/* 获得主题，没有则新建 */
func (mq *Mq) Topic(name string) *Topic {
	ok, topic := mq.isTopic(name)
	if ok {
		return topic
	}
	
	topic = mq.newTopic(name, mq.number, mq.size)
	mq.list = append(mq.list, topic)
	return topic
}