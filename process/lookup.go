package process

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Shopify/vouch4cluster/listers"
	"github.com/Shopify/voucher/auth/google"
)

// LookupAndAttest takes a listers.ImageLister and handles them, writing the output
// of the results to the passed io.Writer.
func LookupAndAttest(cfg *VoucherConfig, lister listers.ImageLister, output io.Writer) error {

	auth := google.NewAuth()

	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()

	images, err := lister.List()
	if nil != err {
		fmt.Printf("listing images failed: %s\n", err)
	}

	processor := Processor{
		ctx:    ctx,
		auth:   auth,
		config: cfg,
	}

	finalResults := new(Result)

	totalImages := len(images)

	inputChan := make(chan workerInput, totalImages)
	resultsChan := make(chan *Result, totalImages)

	workers := 100

	if workers > totalImages {
		workers = totalImages
	}

	for i := 1; i <= workers; i++ {
		go worker(i, &processor, inputChan, resultsChan)
	}

	// registering images to check
	for i, image := range images {
		fmt.Printf("- registering image (%d/%d)\n", i+1, totalImages)
		in := workerInput{
			id:    i + 1,
			image: image,
		}
		inputChan <- in
	}
	close(inputChan)

	for i := 1; i <= totalImages; i++ {
		finalResults.Combine(<-resultsChan)
		fmt.Printf("- got result of image (%d/%d)\n", i+1, totalImages)
	}

	finalResults.Write(output)
	return nil
}

func worker(id int, processor *Processor, inputChan <-chan workerInput, resultsChan chan<- *Result) {
	for input := range inputChan {
		fmt.Printf("- worker %d handling job %d\n", id, input.id)

		vClient, err := newVoucherClient(processor.ctx, processor.config)
		if nil != err {
			fmt.Printf("   - could not setup client: %s\n", err)
		}

		processResult, err := processor.Process(vClient, input.image)
		if nil != err {
			fmt.Printf("  - processing image \"%s\" failed: %s\n", input.image, err)
		}
		fmt.Printf("- worker %d completed job %d\n", id, input.id)

		resultsChan <- processResult
	}
}
