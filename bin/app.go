package main

import (
	"github.com/just1689/scale-aware-proxy-operator/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("starting up...")
	server.StartServer()
}
