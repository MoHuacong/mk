package wic

type LogicRealize struct {
	current Server
	server map[string]map[uint]Server
}

func (lr *LogicRealize) News() {
	if lr.server == nil {
		lr.server = make(map[string]map[uint]Server)
	}
}

func (lr *LogicRealize) Init(ser Server) bool {
	return true
}

func (lr *LogicRealize) End(ser Server) bool {
	return true
}

func (lr *LogicRealize) SetServer(ser Server) bool {
	lr.current = ser
	return true
}

func (lr *LogicRealize) GetServerList() ServerMap {
	return lr.server
}

func (lr *LogicRealize) GetServer(name string) map[uint]Server {
	return lr.server[name]
}

func (lr *LogicRealize) AddServer(ser Server) bool {
	if lr.server[ser.GetName()[0]] == nil {
		lr.server[ser.GetName()[0]] = make(map[uint]Server)
	}
	lr.server[ser.GetName()[0]][ser.GetPort()] = ser
	return true
}

func (lr *LogicRealize) Read(fd *Fd) string {
	return lr.current.Read(fd)
}

func (lr *LogicRealize) Write(fd *Fd, data string) bool {
	//return lr.current.Write(fd, data)
	return fd.Ser.Write(fd, data)
}

func (lr *LogicRealize) Send(fd *Fd, data string) bool {
	//return lr.current.Send(fd, data)
	return fd.Ser.Send(fd, data)
}


func (lr *LogicRealize) Close(fd *Fd) bool {
	//return lr.current.Send(fd, data)
	return fd.Ser.Close(fd)
}


func (lr *LogicRealize) Connect(ser Server, fd *Fd) bool {
	return true
}

func (lr *LogicRealize) Receive(ser Server, fd *Fd, data string) bool {
	return true
}

func (lr *LogicRealize) Closes(ser Server, fd *Fd) bool {
	return true
}