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

	for _, check := range p.config.Checks {
		voucherResp, err := vClient.Check(check, canonicalRef)
		if nil != err {
			result.AddUnprocessible(image)
			return result, fmt.Errorf("running check %s failed: %s", check, err)
		}
		addVoucherResponseToResult(&voucherResp, canonicalRef, result)
	}

	return result, nil
}
