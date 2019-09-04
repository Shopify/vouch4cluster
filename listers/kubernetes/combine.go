package kubernetes

// uniqueImages returns a slice of the unique images in the passed list of
// images.
func uniqueImages(list []string) []string {

	imageExists := make(map[string]bool, len(list))
	uniqueImages := make([]string, 0, len(list))

	for _, item := range list {
		if _, value := imageExists[item]; !value {
			uniqueImages = append(uniqueImages, item)
		}
		imageExists[item] = true
	}

	return uniqueImages
}
