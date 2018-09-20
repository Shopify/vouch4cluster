package listers

// ImageLister is a interface which has a List method.
type ImageLister interface {
	List() ([]string, error) // Return a slice of image names, or an error.
}
