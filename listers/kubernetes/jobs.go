package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listJobs(client kubernetes.Interface, namespace string) ([]string, error) {
	jobs, err := client.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if nil != err {
		return []string{}, nil
	}

	images := make([]string, 0, len(jobs.Items)*2)

	for _, job := range jobs.Items {
		images = append(images, GetImages(&job.Spec.Template.Spec)...)
	}

	return images, nil
}
