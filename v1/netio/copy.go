package netio

import (
	"io"
	"net"
)

func newCopier(l, r net.Conn) (f func(*net.TCPConn, error), err error) {
	f = func(conn *net.TCPConn, err error) {
		done := make(chan bool)
		go func() {
			io.Copy(l, r)
			done <- true
		}()
		go func() {
			io.Copy(r, l)
			done <- true
		}()
		<-done

	}
	return
}
