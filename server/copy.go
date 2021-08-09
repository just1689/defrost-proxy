package server

import (
	"github.com/just1689/scale-aware-proxy-operator/model"
	"io"
	"net"
)

func copier(l *net.TCPConn, r net.Conn) {
	model.Counter.Add()
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
	model.Counter.Sub()

}
