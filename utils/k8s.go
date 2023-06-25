package utils

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func GetKubeConfig(isInCluster bool) (kubeConfig *rest.Config, err error) {
	if isInCluster {
		kubeConfig, err = rest.InClusterConfig()
	} else {
		currentDir, errWd := os.Getwd()
		if errWd != nil {
			return nil, fmt.Errorf("Failed to get current working directory: %v ", err)
		}

		// use the current context in kubeConfig
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(currentDir, "kubeConfig"))
	}

	return kubeConfig, err
}

func K8SHealth(k8sClient *kubernetes.Clientset, ctx context.Context) error {
	path := "/healthz"

	content, err := k8sClient.WithLegacy().RESTClient().Get().AbsPath(path).DoRaw(ctx)
	if err != nil {
		return fmt.Errorf("ErrorBadRequst : %s\n", err.Error())
	}

	contentStr := string(content)
	if contentStr != "ok" {
		return fmt.Errorf("ErrorNotOk : response != 'ok' : %s\n", contentStr)
	}

	return nil
}

// GetCurrentNodeId Get node id of current Pod (Runas DaemonSet)
func GetCurrentNodeId(k8sClient *kubernetes.Clientset, ctx context.Context) (string, error) {
	podName := os.Getenv("HOSTNAME")
	//podNamespace := os.Getenv("POD_NAMESPACE")
	pod, err := k8sClient.CoreV1().Pods("default").Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	node, err := k8sClient.CoreV1().Nodes().Get(ctx, pod.Spec.NodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(node.UID), nil
}
