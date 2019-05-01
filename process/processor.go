package process

import (
	"context"
	"fmt"
	"io"

	"github.com/Shopify/vouch4cluster/listers"
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/auth/google"
	"github.com/Shopify/voucher/client"
)

// Processor handles the images returned by the ImageLister.
type Processor struct {
	ctx     context.Context
	auth    voucher.Auth
	config  *VoucherConfig
	lister  listers.ImageLister
	output  io.Writer
	images  chan string
	results chan *Result
}

// LookupAndAttest takes a listers.ImageLister and handles them, writing the output
// of the results to the passed io.Writer.
func (p *Processor) LookupAndAttest() error {

	images, err := p.lister.List()
	if nil != err {
		fmt.Printf("listing images failed: %s\n", err)
	}

	finalResults := new(Result)

	totalImages := len(images)

	p.images = make(chan string, totalImages)
	p.results = make(chan *Result, totalImages)

	workers := p.config.Workers

	if workers > totalImages {
		workers = totalImages
	}

	for i := 1; i <= workers; i++ {
		go p.worker()
	}

	// registering images to check
	for i, image := range images {
		fmt.Printf("- registering image (%d/%d)\n", i+1, totalImages)
		p.images <- image
	}
	close(p.images)

	for i := 1; i <= totalImages; i++ {
		finalResults.Combine(<-p.results)
		fmt.Printf("- got result of image (%d/%d)\n", i+1, totalImages)
	}

	finalResults.Write(p.output)
	return nil
}

// newVoucherClient returns a new voucher client using the VoucherConfig and
// context.Context specific to the Processor.
func (p *Processor) newVoucherClient() (*client.VoucherClient, error) {
	return newVoucherClient(p.ctx, p.config)
}

// worker is the function that handles image processing.
func (p *Processor) worker() {
	for image := range p.images {
		result, err := checkImage(p, image)
		if nil != err {
			fmt.Printf("  - processing image \"%s\" failed: %s\n", image, err)
		}

		p.results <- result
	}
}

// NewProcessor creates a new Processor,
func NewProcessor(ctx context.Context, cfg *VoucherConfig, lister listers.ImageLister, output io.Writer) *Processor {
	auth := google.NewAuth()

	processor := Processor{
		ctx:    ctx,
		auth:   auth,
		config: cfg,
		lister: lister,
		output: output,
	}

	return &processor
}
