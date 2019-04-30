package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/client"
)

// Processor handles the images returned by the ImageLister.
type Processor struct {
	ctx    context.Context
	auth   voucher.Auth
	config *VoucherConfig
}

// Process processes the passed image using the passed VoucherClient. Returns
// a *Result or an error if something goes wrong.
func (p *Processor) Process(vClient *client.VoucherClient, image string) (*Result, error) {
	result := new(Result)

	if !strings.HasPrefix(image, "gcr.io/") {
		result.AddThirdParty(image)
		return result, nil
	}

	namedRef, err := parseReference(image)
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("getting reference failed: %s", err)
	}

	authClient, err := voucher.AuthToClient(p.ctx, p.auth, namedRef)
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("creating authenticated client failed: %s", err)
	}

	canonicalRef, err := getCanonicalReference(authClient, namedRef)
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("getting image digest failed: %s", err)
	}

	voucherResp, err := vClient.Check("all", canonicalRef)
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("signing image failed: %s", err)
	}

	addVoucherResponseToResult(&voucherResp, canonicalRef, result)

	return result, nil
}
