package kedahttp

import (
	"context"
	"fmt"

	naml "github.com/kris-nova/naml/pkg"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Installation struct {
	metav1.ObjectMeta
}

var _ naml.Deployable = Installation{}

func New(namespace, name string) Installation {
	return Installation{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			ResourceVersion: "v1.0.0",
			Labels: map[string]string{
				"k8s-app": "keda-http-add-on",
				"app":     "keda-http-add-on",
			},
			Annotations: map[string]string{
				"installed-by": "http-add-on-installer",
			},
		},
	}
}

func (i Installation) Install(client *kubernetes.Clientset) error {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: i.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: naml.I32p(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: i.Labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: i.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  i.Name,
							Image: "busybox",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := client.AppsV1().Deployments(i.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to install deployment in Kubernetes: %v", err)
	}
	return nil
}

// Uninstall will attempt to uninstall in Kubernetes
func (i Installation) Uninstall(client *kubernetes.Clientset) error {
	return client.AppsV1().Deployments(i.Namespace).Delete(context.TODO(), i.Name, metav1.DeleteOptions{})

}

// Meta returns the Kubernetes native ObjectMeta which is used to manage applications with naml.
func (i Installation) Meta() *metav1.ObjectMeta {
	return &i.ObjectMeta
}
