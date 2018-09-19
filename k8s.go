package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// k8sImageLister lists images running in Kubernetes.
type k8sImageLister struct {
	path   string
	client *kubernetes.Clientset
}

// List returns a list of all of the images available to the current Kubernetes context,
// or an error.
func (k *k8sImageLister) List() ([]string, error) {
	pods, err := k.client.CoreV1().Pods("").List(metav1.ListOptions{})
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

// NewK8SImageLister creates a new Kubernetes specific ImageLister
func NewK8SImageLister(path string) (ImageLister, error) {
	k := new(k8sImageLister)

	kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: k.path},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}},
	).ClientConfig()
	if nil != err {
		return nil, err
	}

	k.client, err = kubernetes.NewForConfig(kubeconfig)
	return k, err
}
