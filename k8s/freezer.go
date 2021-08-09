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
			if !model.ScaledUp.Get() {
				//Not running, can ignore
				logrus.Infoln("not scaled up, going back to sleep")
				continue
			}
			if model.Counter.Get() > 0 {
				//Open long-running connection...
				logrus.Infoln("more than one open connection, going back to sleep")
				continue

			}

			//Cooldown
			logrus.Infoln("Time to scale down....")
			d := KC.GetDeployment(model.ThisTarget.Namespace, model.ThisTarget.Name)
			if d == nil {
				logrus.Errorln("could not get target", model.ThisTarget)
				continue
			}
			t := model.ThisTarget
			t.Replicas = 0
			updateDeploymentReplicas(d, t)
			KC.updateDeploymentInK8s(d) //TODO: make single threaded,
			//TODO: get result of `updateDeploymentInK8s`
			t.Replicas = 0
			model.ScaledUp.Set(false)
		case <-f.ping:
			logrus.Infoln("< ping >")

		}
		//TODO: impl
	}
}
