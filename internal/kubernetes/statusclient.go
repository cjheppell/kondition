package kubernetes

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type statusClient struct {
	k8sClient *kubernetes.Clientset
}

func NewStatusClient(kubeConfigPath string) (*statusClient, error) {
	var k8sConfig *rest.Config
	if kubeConfigPath == "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		k8sConfig = config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
		k8sConfig = config
	}

	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return &statusClient{
		k8sClient: k8sClient,
	}, nil
}

func (sc *statusClient) IsDeploymentReady(deploymentName, namespace string) (bool, error) {
	deployment, err := sc.k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable {
			if condition.Status == corev1.ConditionTrue {
				return true, nil
			}
			return false, nil
		}
		continue
	}

	return false, fmt.Errorf("could not find DeploymentAvailable condition on the deployment")
}