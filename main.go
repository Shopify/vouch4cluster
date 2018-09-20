package main

import (
	"bufio"
	"fmt"
	"os"
)

func fromStdin() {
	images := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		image := scanner.Text()
		images = append(images, image)
	}
	if err := scanner.Err(); nil != err {
		fmt.Printf("scanning failed: %s", err)
	}

}

func main() {

	images, err := fromKubernetes(os.Getenv("HOME") + "/.kube/config")
	if nil != err {
		fmt.Printf("getting images failed: %s", err)
		os.Exit(1)
	}

	err = lookupAndAttest(images)
	if nil != err {
		fmt.Printf("attesting images failed: %s", err)
		os.Exit(1)
	}

}
