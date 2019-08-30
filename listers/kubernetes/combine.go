package kubernetes

func combineAndReduce(list []string) []string {

	unique := make(map[string]bool, len(list))

	for _, item := range list {
		unique[item] = true
	}

	images := make([]string, 0, len(unique))
	for item := range unique {
		images = append(images, item)
	}

	return images
}
