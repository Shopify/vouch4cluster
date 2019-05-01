package process

import (
	"fmt"
	"strings"
)

// checkImage gets the Canonical image reference based on the passed string,
// converting tags to digests, submits the reference to Voucher, and creates a
// Result for the call.
func checkImage(p *Processor, image string) (*Result, error) {
	result := new(Result)

	if !strings.HasPrefix(image, "gcr.io/") {
		result.AddThirdParty(image)
		return result, nil
	}

	vClient, err := p.newVoucherClient()
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("setting up client failed: %s", err)
	}

	canonicalRef, err := toCanonicalReference(p.ctx, p.auth, image)
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("getting image reference failed: %s", err)
	}

	voucherResp, err := vClient.Check("all", canonicalRef)
	if nil != err {
		result.AddUnprocessible(image)
		return result, fmt.Errorf("signing image failed: %s", err)
	}

	addVoucherResponseToResult(&voucherResp, canonicalRef, result)

	return result, nil
}
