package model

import (
	"github.com/sirupsen/logrus"
	"os"
)

var ThisTarget = getTarget()

func getTarget() Target {
	return Target{
		Namespace: getEnvOrDie("TARGET_NAMESPACE"),
		Name:      getEnvOrDie("TARGET_NAME"),
		Replicas:  1,
	}
}

type Target struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Replicas  int    `json:"replicas"`
}

func getEnvOrDie(name string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		logrus.Panicln("could not get ENV", name, ". Exiting")
	}
	return
}
