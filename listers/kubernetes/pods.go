package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listPods(client kubernetes.Interface, namespace string) ([]string, error) {
	pods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if nil != err {
		return []string{}, nil
	}

	images := make([]string, 0, len(pods.Items)*2)

	for _, pod := range pods.Items {
		images = append(images, GetImages(&pod.Spec)...)
	}

	return images, nil
}
