package kubernetes

import (
	"github.com/Shopify/vouch4cluster/listers"
	"k8s.io/client-go/kubernetes"
)

// k8sImageLister lists images running in Kubernetes.
type k8sImageLister struct {
	client kubernetes.Interface
}

// k8sResourceLister is a function that lists images used by a resource.
type k8sResourceLister func(kubernetes.Interface, string) ([]string, error)

// List returns a list of all of the images available to the current Kubernetes context,
// or an error.
func (k *k8sImageLister) List() ([]string, error) {

	allImages := make([]string, 0, 100)

	for _, resourceLister := range []k8sResourceLister{
		listPods,
		listJobs,
	} {

		resourceImages, err := resourceLister(k.client, "")
		if nil != err {
			return []string{}, err
		}

		allImages = append(allImages, resourceImages...)

	}

	images := combineAndReduce(allImages)
	return images, nil
}

// NewImageLister creates a new Kubernetes specific ImageLister with the passed
// kubernetes.Interface.
func NewImageLister(clientSet kubernetes.Interface) listers.ImageLister {
	k := &k8sImageLister{
		client: clientSet,
	}

	return k
}
