package main

import (
	"bufio"
	"fmt"
	"io"
)

// readerImageLister is an implementation of ImageLister which reads
// Images from a file.
type readerImageLister struct {
	reader io.Reader
}

// List returns a list of images from the passed io.Reader.
func (r *readerImageLister) List() ([]string, error) {
	images := make([]string, 0)

	scanner := bufio.NewScanner(r.reader)
	for scanner.Scan() {
		image := scanner.Text()
		images = append(images, image)
	}
	if err := scanner.Err(); nil != err {
		return images, fmt.Errorf("scanning failed: %s", err)
	}

	return images, nil
}

// NewReaderImageLister creates a new ImageLister that reads from an
// io.Reader.
func NewReaderImageLister(reader io.Reader) ImageLister {
	lister := new(readerImageLister)
	lister.reader = reader
	return lister
}
