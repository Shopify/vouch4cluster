package kubernetes

import (
	"github.com/Shopify/vouch4cluster/listers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// add GCP authentication support to the kubernetes code
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// k8sImageLister lists images running in Kubernetes.
type k8sImageLister struct {
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

// NewImageLister creates a new Kubernetes specific ImageLister
func NewImageLister(configFilename string) (listers.ImageLister, error) {
	k := new(k8sImageLister)

	kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configFilename},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}},
	).ClientConfig()
	if nil != err {
		return nil, err
	}

	k.client, err = kubernetes.NewForConfig(kubeconfig)
	return k, err
}
