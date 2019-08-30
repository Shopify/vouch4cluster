package kubernetes

import (
	"github.com/Shopify/vouch4cluster/listers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// k8sImageLister lists images running in Kubernetes.
type k8sImageLister struct {
	client kubernetes.Interface
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

// NewImageLister creates a new Kubernetes specific ImageLister with the passed
// kubernetes.Interface.
func NewImageLister(clientSet kubernetes.Interface) listers.ImageLister {
	k := &k8sImageLister{
		client: clientSet,
	}

	return k
}
