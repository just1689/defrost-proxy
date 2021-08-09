package k8s

import (
	"github.com/just1689/scale-aware-proxy-operator/model"
	"github.com/sirupsen/logrus"
)

func NewScaler() chan model.Target {
	//TODO: env var for chan size?
	scaler := make(chan model.Target, 256)
	startScaler(scaler)
	return scaler
}

func startScaler(in chan model.Target) {
	asyncForEachTarget(in, func(target model.Target) {
		logrus.Infoln("Scaler > Going to scale", target.Namespace, target.Name)
		deployment := KC.GetDeployment(target.Namespace, target.Name)
		if deployment == nil {
			logrus.Errorln("did not find object", target.Namespace, target.Name)
			return
		}
		changed := updateDeploymentReplicas(deployment, target)
		if changed {
			logrus.Infoln("Scaler >", target.Namespace, target.Name, "is being updated")
			KC.updateDeploymentInK8s(deployment)
			model.ScaledUp.Set(true)
		} else {
			logrus.Infoln("Scaler >", target.Namespace, target.Name, "does not need to change")
		}
	})
}

func asyncForEachTarget(c chan model.Target, f func(next model.Target)) {
	go func() {
		for target := range c {
			f(target)
		}
	}()
}
