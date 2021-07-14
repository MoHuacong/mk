package wic

import (
	"github.com/MoHuacong/wic/tools"
	"strconv"
)

type CallBack func() interface{}

var _id *tools.Count = &tools.Count{}
var _list map[uint64]*Fd = make(map[uint64]*Fd)

func GetFd(id uint64) *Fd {
	return _list[id]
}

func NewFd(ser Server, v interface{}) *Fd {
	id := _id.Add()
	fd := &Fd{Id: id, Ser: ser, Validate: false, Value: v}
	_list[id] = fd
	return fd
}

type Fd struct {
	Id uint64
	Ser Server
	Validate bool
	Value interface{}
}

func (fd *Fd) SetValidate(v bool)  {
	fd.Validate = v
}

func (fd *Fd) IdString() string {
	return strconv.FormatUint(fd.Id, 10)
}