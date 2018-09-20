package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	//	images, err := fromKubernetes(os.Getenv("HOME") + "/.kube/config")
	//	if nil != err {
	//		fmt.Printf("getting images failed: %s", err)
	//		os.Exit(1)
	//	}

	lister := NewReaderImageLister(os.Stdin)

	err = lookupAndAttest(lister, os.Stdout)
	if nil != err {
		fmt.Printf("attesting images failed: %s", err)
		os.Exit(1)
	}

}
