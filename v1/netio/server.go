package netio

import (
	"github.com/sirupsen/logrus"
	"net"
	"os"
)

var (
	EnvLocalAddr  = "LOCAL_ADDR"
	EnvRemoteAddr = "REMOTE_ADDR"
	localAddr     = addressAsTCPAddr(EnvLocalAddr)
	remoteAddr    = addressAsTCPAddr(EnvRemoteAddr)
)

func newServer(listeners int, f func(*net.TCPConn, error)) (err error) {
	var listener *net.TCPListener
	listener, err = net.ListenTCP("tcp", localAddr)
	if err != nil {
		logrus.Panicln("Panic. Could not listen on address", localAddr.String())
		return
	}
	for i := 1; i <= listeners; i++ {
		logrus.Infoln("... starting listener", i)
		go func() {
			for {
				l, e := listener.AcceptTCP()
				if err != nil {
					logrus.Errorln(err)
					continue
				}
				go f(l, e)
			}
		}()
	}
	return
}

func addressAsTCPAddr(envVar string) *net.TCPAddr {
	addr := os.Getenv(envVar)
	if addr == "" {
		logrus.Panicln("Panic. no value found for environment variable", envVar, ". Exiting")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logrus.Panicln("Panic. Failed to resolve address: %s", err)
	}
	return tcpAddr

}
