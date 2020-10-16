package kubernetes

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type StatusClient struct {
	k8sClient *kubernetes.Clientset
	logger *zap.SugaredLogger
}

func NewStatusClient(kubeConfigPath string, logger *zap.SugaredLogger) (*StatusClient, error) {
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

	return &StatusClient{
		k8sClient: k8sClient,
		logger: logger,
	}, nil
}

func (sc *StatusClient) IsDeploymentReady(deploymentName, namespace string) (bool, error) {
	sc.logger.Debugf("checking deployment status for deployment '%s' in namespace '%s'", deploymentName, namespace)
	deployment, err := sc.k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable {
			if condition.Status == corev1.ConditionTrue {
				sc.logger.Debugf("deployment '%s' in namespace '%s' was available", deploymentName, namespace)
				return true, nil
			}
			sc.logger.Debugf("deployment '%s' in namespace '%s' was not available", deploymentName, namespace)
			return false, nil
		}
		continue
	}

	return false, fmt.Errorf("could not find DeploymentAvailable condition on the deployment")
}