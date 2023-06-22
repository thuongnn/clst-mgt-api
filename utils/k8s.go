package utils

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

func GetKubeConfig(isInCluster bool) (kubeConfig *rest.Config, err error) {
	if isInCluster {
		kubeConfig, err = rest.InClusterConfig()
	} else {
		//currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		//if err != nil {
		//	return nil, err
		//}

		currentDir := "/Users/thuongnn/Desktop/thuongnn/projects/golang-mongodb-api"

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
