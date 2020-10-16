package kubernetes

import (
	"context"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type StatusClient struct {
	k8sClient *kubernetes.Clientset
	logger    *zap.SugaredLogger
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
		logger:    logger,
	}, nil
}

func (sc *StatusClient) IsDeploymentReady(deploymentName, namespace string, minReadyReplicas int32) (bool, error) {
	sc.logger.Debugf("checking deployment status for deployment '%s' in namespace '%s'", deploymentName, namespace)
	deployment, err := sc.k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	readyReplicas := deployment.Status.ReadyReplicas
	if readyReplicas >= minReadyReplicas {
		sc.logger.Debugf("deployment '%s' in namespace '%s' has sufficient ready replicas (%d/%d)", deploymentName, namespace, readyReplicas, minReadyReplicas)
		return true, nil
	} else {
		sc.logger.Debugf("deployment '%s' in namespace '%s' has insufficient ready replicas (%d/%d)", deploymentName, namespace, readyReplicas, minReadyReplicas)
		return false, nil
	}
}
