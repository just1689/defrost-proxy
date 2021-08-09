package server

import (
	"fmt"
	"github.com/just1689/scale-aware-proxy-operator/k8s"
	"github.com/just1689/scale-aware-proxy-operator/model"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"time"
)

var (
	EnvLocalAddr  = "LOCAL_ADDR"
	EnvRemoteAddr = "REMOTE_ADDR"
	scaler        = k8s.NewScaler()
	localAddr     = addressAsTCPAddr(EnvLocalAddr)
	remoteAddr    = addressAsTCPAddr(EnvRemoteAddr)
)

func StartServer() {
	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		logrus.Panicln("Failed to open local port to listen: %s", err)
	}
	for i := 1; i <= 9; i++ {
		logrus.Infoln("starting worker", i)
		go handleIncomingConnections(listener)
	}
	logrus.Infoln("starting worker", 10)
	handleIncomingConnections(listener)

}

func handleIncomingConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logrus.Errorln("Failed to accept TCP conn on local port: %s", err)
			continue
		}

		if !model.ScaledUp.Get() {
			logrus.Infoln("new connection < scaler getTarget()")
			scaler <- model.ThisTarget
		}
		go handleIncomingConnection(conn)
	}
}

func handleIncomingConnection(conn *net.TCPConn) {
	ok := false
	tries := 0
	start := time.Now()
	requestedBoot := false
	//TODO: env vars for both?
	for !ok && tries < 100 {
		//TODO: env var for TCP dial timeout
		remoteConn, err := net.DialTimeout("tcp", remoteAddr.String(), 1*time.Second)
		if err != nil {
			tries++
			logrus.Errorln(err)
			logrus.Errorln("tries:", tries)
			if time.Since(start) > time.Duration(1) && !requestedBoot {
				requestedBoot = true
				scaler <- model.ThisTarget
			}
			continue
		}
		k8s.Freezer.Ping()
		copier(conn, remoteConn)
		closeNamedConnection("remoteConn", remoteConn)
		closeNamedConnection("conn", conn)
		ok = true
		return
	}
	closeNamedConnection("conn", conn)

}

func closeNamedConnection(name string, conn net.Conn) {
	if err := conn.Close(); err != nil {
		logrus.Errorln(fmt.Sprintf("%s.Close()", name), err)
	}
}

func addressAsTCPAddr(envVar string) *net.TCPAddr {
	addr := os.Getenv(envVar)
	if addr == "" {
		logrus.Panicln("no value found for environment variable", envVar, ". Exiting")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logrus.Panicln("Failed to resolve address: %s", err)
	}
	return tcpAddr

}
