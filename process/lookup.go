package process

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Shopify/vouch4cluster/listers"
	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/auth/google"
	"github.com/Shopify/voucher/client"
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

	processor := processor{
		ctx:  ctx,
		auth: auth,
	}

	totalImages := len(images)

	for i, image := range images {
		fmt.Printf("- handling image (%d/%d)\n", i+1, totalImages)
		vClient, err := newVoucherClient(ctx, cfg)
		if nil != err {
			fmt.Printf("   - could not setup client: %s\n", err)
		}

		err = processor.Process(vClient, image)
		if nil != err {
			fmt.Printf("  - processing image \"%s\" failed: %s\n", image, err)
		}

	}

	_ = writeProcessResult(&processor, output)
	return nil
}

// processor handles the images returned by the ImageLister.
type processor struct {
	ctx           context.Context
	auth          voucher.Auth
	successes     []voucher.CheckResult
	failures      []voucher.CheckResult
	unprocessible []string
	thirdParty    []string
}

// processImage processes the passed image.
func (p *processor) Process(vClient *client.VoucherClient, image string) error {
	if !strings.HasPrefix(image, "gcr.io/shopify-docker-images") {
		p.thirdParty = append(p.thirdParty, image)
		return nil
	}

	namedRef, err := parseReference(image)
	if nil != err {
		p.unprocessible = append(p.unprocessible, image)
		return fmt.Errorf("getting reference failed: %s", err)
	}

	authClient, err := voucher.AuthToClient(p.ctx, p.auth, namedRef)
	if nil != err {
		p.unprocessible = append(p.unprocessible, image)
		return fmt.Errorf("creating authenticated client failed: %s", err)
	}

	canonicalRef, err := getCanonicalReference(authClient, namedRef)
	if nil != err {
		p.unprocessible = append(p.unprocessible, image)
		return fmt.Errorf("getting image digest failed: %s", err)
	}

	voucherResp, err := vClient.Check("all", canonicalRef)
	if nil != err {
		p.unprocessible = append(p.unprocessible, image)
		return fmt.Errorf("signing image failed: %s", err)
	}

	for _, result := range voucherResp.Results {
		result.ImageData = voucher.ImageData(canonicalRef)
		if result.Attested {
			p.successes = append(p.successes, result)
		} else {
			p.failures = append(p.failures, result)
		}
	}

	return nil
}
