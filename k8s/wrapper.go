package k8s

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var KC = connectToKubernetesAPI()

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func connectToKubernetesAPI() *K8sClient {
	result := &K8sClient{}
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	result.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Panicln(err)
	}
	return result
}
