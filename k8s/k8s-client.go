package k8s

import (
	"context"
	"github.com/just1689/scale-aware-proxy-operator/model"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var KC = connectToKubernetesAPI()

func updateDeploymentReplicas(d *v1.Deployment, target model.Target) (changed bool) {
	actual := 0
	if d.Spec.Replicas != nil {
		actual = int(*d.Spec.Replicas)
	}

	if actual != target.Replicas {
		logrus.Infoln(d.ObjectMeta.Namespace, d.ObjectMeta.Name, "changing replicas to", target.Replicas)
		d.Spec.Replicas = intToInt64P(target.Replicas)
		changed = true
		return
	} else {
		logrus.Infoln(d.ObjectMeta.Namespace, d.ObjectMeta.Name, "skipped as replicas already =", target.Replicas)
	}
	changed = false
	return
}

func intToInt64P(i int) *int32 {
	x := int32(i)
	return &x
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

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func (k *K8sClient) GetDeployment(ns, name string) *v1.Deployment {
	deployment, err := k.clientSet.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorln("could not get deployment", ns, name)
		logrus.Errorln(err)
	}
	return deployment
}

func (k *K8sClient) updateDeploymentInK8s(d *v1.Deployment) {
	if _, err := k.clientSet.AppsV1().Deployments(d.ObjectMeta.Namespace).Update(context.TODO(), d, metav1.UpdateOptions{}); err != nil {
		logrus.Errorln(err)
	}
}

func (k *K8sClient) SwitchOffDeployment() {
	k.GetDeployment(model.ThisTarget.Namespace, model.ThisTarget.Name)

}
