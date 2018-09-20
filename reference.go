package main

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

func parseReference(image string) (reference.Named, error) {
	var namedRef reference.Named
	var ok bool

	ref, err := reference.Parse(image)
	if nil != err {
		return namedRef, fmt.Errorf("parsing image reference failed: %s", err)
	}

	if namedRef, ok = ref.(reference.Named); !ok {
		return namedRef, fmt.Errorf("couldn't get named version of reference: %s", err)
	}

	return namedRef, nil
}
