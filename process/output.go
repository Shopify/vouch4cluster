package process

import (
	"fmt"
	"io"

	"github.com/Shopify/voucher"
	"github.com/docker/distribution/reference"
)

// Result stores the result of a Processor.Process call.
type Result struct {
	Successes     []voucher.CheckResult
	Failures      []voucher.CheckResult
	Unprocessible []string
	ThirdParty    []string
}

// AddUnprocessible adds unprocessable images to our Result.
func (result *Result) AddUnprocessible(images ...string) {
	result.Unprocessible = append(result.Unprocessible, images...)
}

// AddThirdParty adds unprocessable images to our Result.
func (result *Result) AddThirdParty(images ...string) {
	result.ThirdParty = append(result.ThirdParty, images...)
}

// AddSuccesses adds successful images to our Result.
func (result *Result) AddSuccesses(checkResults ...voucher.CheckResult) {
	result.Successes = append(result.Successes, checkResults...)
}

// AddFailures adds failed images to our Result.
func (result *Result) AddFailures(checkResults ...voucher.CheckResult) {
	result.Failures = append(result.Failures, checkResults...)
}

// Combine adds the contents of the passed Result into the calling Result.
func (result *Result) Combine(newResult *Result) {
	result.AddSuccesses(newResult.Successes...)
	result.AddFailures(newResult.Failures...)
	result.AddUnprocessible(newResult.Unprocessible...)
	result.AddThirdParty(newResult.ThirdParty...)
}

// Write writes the Results to io.Writer.
func (result *Result) Write(output io.Writer) error {
	_, _ = fmt.Fprintln(output, "=== Processing Results ===")

	_, _ = fmt.Fprintln(output, "--- Successes ---")
	for _, success := range result.Successes {
		_, _ = fmt.Fprintf(output, "%s (passed %s)\n", success.ImageData, success.Name)
	}

	_, _ = fmt.Fprintln(output, "--- Failures ---")
	for _, failure := range result.Failures {
		_, _ = fmt.Fprint(output, failure.ImageData)
		if failure.Success {
			_, _ = fmt.Fprintf(output, "(passed %s,", failure.Name)
			if "" == failure.Err {
				_, _ = fmt.Fprint(output, " but wasn't attested)\n")
			} else {
				_, _ = fmt.Fprintf(output, " but attestation failed: %s)\n", failure.Err)
			}
		} else {
			_, _ = fmt.Fprintf(output, "(failed %s", failure.Name)
			if "" == failure.Err {
				_, _ = fmt.Fprint(output, " check)\n")
			} else {
				_, _ = fmt.Fprintf(output, ": %s)\n", failure.Err)
			}
		}
	}

	_, _ = fmt.Fprintln(output, "--- Unprocessible ---")
	for _, image := range result.Unprocessible {
		_, _ = fmt.Fprintln(output, image)
	}

	_, _ = fmt.Fprintln(output, "--- Third Party---")
	for _, image := range result.ThirdParty {
		_, _ = fmt.Fprintln(output, image)
	}

	return nil
}

// addVoucherResponseToResult adds the contents of a voucher.CheckResults to a
// Result.
func addVoucherResponseToResult(resp *voucher.Response, canonicalRef reference.Canonical, result *Result) {
	for _, checkResult := range resp.Results {
		checkResult.ImageData = voucher.ImageData(canonicalRef)
		if checkResult.Attested {
			result.AddSuccesses(checkResult)
		} else {
			result.AddFailures(checkResult)
		}
	}
}
