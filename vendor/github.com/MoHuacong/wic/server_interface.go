package wic

/* 服务解析处理接口 */
type Server interface {
	IO
	Init()
	GetPort() uint
	
	GetName() []string
	IsName(name string) bool
	SetName(name []string) bool
	
	Run(blocking bool) error
	SetIpPort(ip string, port int)
	SetLogic(logic Logic) bool
	
	GetRaw(name string) interface{}
	
	callBack(typ string, ser Server, value interface{}, data string)
}

//type Greeting func(name string) string
type ServerFunc func() Server