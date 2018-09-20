package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func fromKubernetes(path string) ([]string, error) {
	kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}},
	).ClientConfig()
	if nil != err {
		return []string{}, nil
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if nil != err {
		return []string{}, nil
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if nil != err {
		return []string{}, nil
	}

	allImages := make(map[string]struct{})

	for _, pod := range pods.Items {
		for _, status := range pod.Status.ContainerStatuses {
			if _, ok := allImages[status.Image]; !ok {
				allImages[status.Image] = struct{}{}
			}
		}
	}

	uniqueImages := make([]string, 0, len(pods.Items))

	for image := range allImages {
		uniqueImages = append(uniqueImages, image)
	}

	return uniqueImages, err
}
