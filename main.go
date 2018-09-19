package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const timeout = 120 * time.Second

func main() {
	images := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		image := scanner.Text()
		images = append(images, image)
	}
	if err := scanner.Err(); nil != err {
		fmt.Printf("scanning failed: %s", err)
	}

	lookupAndAttest(images...)
}
