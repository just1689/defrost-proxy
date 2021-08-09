package k8s

import (
	"context"
	"github.com/just1689/scale-aware-proxy-operator/model"
	"github.com/just1689/scale-aware-proxy-operator/util"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sClient) SwitchDeployment(target model.Target, toReplicas int) (err error, changed bool) {
	var d *v1.Deployment
	err, d = k.GetDeployment(target.Namespace, target.Name)
	if err != nil || d == nil {
		return
	}
	err, changed = k.SetDeploymentReplica(d, toReplicas)
	if err != nil || !changed {
		return
	}
	err = k.SaveDeployment(d)
	changed = err == nil
	return
}

func (k *K8sClient) GetDeployment(ns, name string) (err error, d *v1.Deployment) {
	d, err = k.clientSet.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorln("could not get deployment", ns, name)
		logrus.Errorln(err)
	}
	return
}

func (k *K8sClient) SaveDeployment(d *v1.Deployment) (err error) {
	if _, err = k.clientSet.AppsV1().Deployments(d.ObjectMeta.Namespace).Update(context.TODO(), d, metav1.UpdateOptions{}); err != nil {
		logrus.Errorln(err)
	}
	return
}

func (k *K8sClient) SetDeploymentReplica(d *v1.Deployment, targetReplicas int) (err error, changed bool) {
	actual := 0
	if d.Spec.Replicas != nil {
		actual = int(*d.Spec.Replicas)
	}
	if actual != targetReplicas {
		logrus.Infoln(d.ObjectMeta.Namespace, d.ObjectMeta.Name, "changing replicas to", targetReplicas)
		d.Spec.Replicas = util.IntToInt64P(targetReplicas)
		changed = true
		return
	} else {
		logrus.Infoln(d.ObjectMeta.Namespace, d.ObjectMeta.Name, "skipped as replicas already =", targetReplicas)
	}
	changed = false
	return
}
