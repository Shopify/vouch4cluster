package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listCronJobs(client kubernetes.Interface, namespace string) ([]string, error) {
	cronJobs, err := client.BatchV1beta1().
		CronJobs(namespace).
		List(metav1.ListOptions{})

	if nil != err {
		return []string{}, nil
	}

	images := make([]string, 0, len(cronJobs.Items)*2)

	for _, cronJob := range cronJobs.Items {
		images = append(images, GetImages(&cronJob.Spec.JobTemplate.Spec.Template.Spec)...)
	}

	return images, nil
}
