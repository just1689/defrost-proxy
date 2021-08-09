package k8s

import (
	"github.com/just1689/scale-aware-proxy-operator/model"
	"github.com/just1689/scale-aware-proxy-operator/util"
	"github.com/sirupsen/logrus"
)

//TODO: size to env to Helm
var scalerCacheSize = 256

type Scaler struct {
	c chan interface{}
}

func NewScaler() *Scaler {
	result := &Scaler{
		c: make(chan interface{}, scalerCacheSize),
	}
	result.run()
	return result
}

func (s *Scaler) Next() {
	s.c <- true
}

func (s *Scaler) run() {
	go util.ForEach(s.c, scaleUp)
}
func scaleUp() {
	logrus.Infoln("Scaler :: Going to scale", model.ThisTarget.Namespace, model.ThisTarget.Name)
	err, changed := KC.SwitchDeployment(model.ThisTarget, model.ThisTarget.Replicas)
	if err != nil {
		logrus.Errorln("could not warm up", model.ThisTarget.Name, model.ThisTarget.Name)
		return
	}
	if changed {
		logrus.Infoln("Scaler :: Going to scale", model.ThisTarget.Namespace, model.ThisTarget.Name, " :: OK")
		model.ScaledUp.Set(true)
	}
}
