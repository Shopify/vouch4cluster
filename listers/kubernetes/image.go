package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
)

// GetImages returns a slice of strings containing all of the images in the
// passed pod.
func GetImages(podSpec *corev1.PodSpec) []string {
	images := make([]string, 0, len(podSpec.Containers)+len(podSpec.InitContainers))

	for _, list := range [][]corev1.Container{
		podSpec.InitContainers,
		podSpec.Containers,
	} {
		for _, container := range list {
			images = append(images, container.Image)
		}
	}

	return images
}
