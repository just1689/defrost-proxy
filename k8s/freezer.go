package k8s

import (
	"github.com/just1689/scale-aware-proxy-operator/model"
	"github.com/sirupsen/logrus"
	"time"
)

var Freezer = newFreezer()

func newFreezer() *freezer {
	//TODO: include cooldown as env var
	f := &freezer{
		ping:    make(chan interface{}, 256),
		seconds: 60,
	}
	go f.run()
	return f
}

type freezer struct {
	ping    chan interface{}
	seconds int
}

func (f *freezer) Ping() {
	f.ping <- true
}

func (f *freezer) run() {
	for {
		select {
		case <-time.After(time.Second * time.Duration(f.seconds)):
			handleNextFreeze()
		case <-f.ping:
			logrus.Infoln("< ping >")
		}
	}
}

func handleNextFreeze() {
	if !model.ScaledUp.Get() {
		//Not running, can ignore
		logrus.Infoln("not scaled up, going back to sleep")
		return
	}
	if model.Counter.Get() > 0 {
		//Open long-running connection...
		logrus.Infoln("more than one open connection, going back to sleep")
		return
	}
	logrus.Infoln("Time to scale down....")
	model.ScaledUp.Set(false)
}
