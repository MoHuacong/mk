package wic

/* socket io */
type IO interface {
	Close(fd *Fd) bool
	Read(fd *Fd) string
	Write(fd *Fd, data string) bool
	
	Send(fd *Fd, data string) bool
}