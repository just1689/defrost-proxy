package netio

import (
	"github.com/sirupsen/logrus"
	"net"
)

func NewOrchestration(synchronousScale func()) {
	newServer(10, func(conn *net.TCPConn, err error) {
		defer closer(conn)
		c, err := connectToRemote(1)
		defer closer(c)
		if err != nil {
			synchronousScale()
			c, err = connectToRemote(100)
		}
		if err != nil {
			logrus.Errorln(err)
		}
		newCopier(conn, c)
	})
}

func closer(c net.Conn) {
	if c != nil {
		c.Close()
	}
}
