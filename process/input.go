package process

// workerInput is the input passed to a worker for processing. It contains the
// id of the image to check, and the image name to validate.
type workerInput struct {
	id    int
	image string
}
