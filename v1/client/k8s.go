package client

import (
	"context"
	"github.com/just1689/scale-aware-proxy-operator/v1/model"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

var Scaler = newScaler()

func newScaler() *scaler {
	r := &scaler{
		k: connectToK8s(),
	}
	return r
}

type scaler struct {
	l sync.Mutex
	k *kubernetes.Clientset
}

func (s *scaler) HandleInstruction(si model.ScaleInstruction) (changed bool, err error) {
	s.l.Lock()
	defer s.l.Unlock()
	var d *v1.Deployment
	d, err = getDeployment(s, si.Namespace, si.Name)
	if d == nil || err != nil {
		return
	}
	if changed = changeReplicaCount(d, si.Replicas); !changed {
		logrus.Infoln("... app already has", si.Replicas, "replicas")
		return
	}
	err = updateDeployment(s, d)
	getReplicas(d)
	changed = err == nil
	return
}

func connectToK8s() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Panicln(err)
	}
	return clientSet
}

func getDeployment(s *scaler, ns, name string) (d *v1.Deployment, err error) {
	d, err = s.k.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorln("could not get deployment", ns, name)
		logrus.Errorln(err)
	}
	return
}

func getReplicas(d *v1.Deployment) int {
	if d == nil {
		return -1
	}
	if d.Spec.Replicas == nil {
		return -1
	}
	result := int(*d.Spec.Replicas)
	model.LastKnownReplicas.Set(result)
	return result
}

func changeReplicaCount(d *v1.Deployment, targetReplicaCount int) (changed bool) {
	currentReplicas := getReplicas(d)
	if currentReplicas == targetReplicaCount {
		changed = false
		return
	}
	d.Spec.Replicas = intToInt64P(targetReplicaCount)
	changed = true
	return
}

func updateDeployment(s *scaler, d *v1.Deployment) (err error) {
	if _, err = s.k.AppsV1().Deployments(d.ObjectMeta.Namespace).Update(context.TODO(), d, metav1.UpdateOptions{}); err != nil {
		logrus.Errorln(err)
	}
	return
}

func intToInt64P(i int) *int32 {
	x := int32(i)
	return &x
}
