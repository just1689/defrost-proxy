package netio

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

func connectToRemote(tryLimit int) (c net.Conn, err error) {
	success := false
	tryNo := 0
	for !success && tryNo < tryLimit {
		c, err = net.DialTimeout("tcp", remoteAddr.String(), 1*time.Second)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				//try again
				tryNo++
				continue
			}
			logrus.Errorln(err)
			return
		}
		success = true
	}
	return
}
