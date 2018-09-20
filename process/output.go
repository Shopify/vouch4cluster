package process

import (
	"fmt"
	"io"
)

// writeProcessResult writes the processor's output to the passed io.writer.
func writeProcessResult(processor *processor, output io.Writer) error {
	_, _ = fmt.Fprintln(output, "=== Processing Results ===")

	_, _ = fmt.Fprintln(output, "--- Successes ---")
	for _, success := range processor.successes {
		_, _ = fmt.Fprintf(output, "%s (passed %s)\n", success.ImageData, success.Name)
	}

	_, _ = fmt.Fprintln(output, "--- Failures ---")
	for _, failure := range processor.failures {
		if failure.Success {
			_, _ = fmt.Fprintf(output, "%s (passed %s, but attestation failed: %s)\n", failure.ImageData, failure.Name, failure.Err)
		} else {
			if "" == failure.Err {
				_, _ = fmt.Fprintf(output, "%s (failed %s check)\n", failure.ImageData, failure.Name)
			} else {
				_, _ = fmt.Fprintf(output, "%s (failed %s: %s)\n", failure.ImageData, failure.Name, failure.Err)
			}
		}
	}

	_, _ = fmt.Fprintln(output, "--- Unprocessable ---")
	for _, image := range processor.unprocessible {
		_, _ = fmt.Fprintln(output, image)
	}

	_, _ = fmt.Fprintln(output, "--- Third Party---")
	for _, image := range processor.thirdParty {
		_, _ = fmt.Fprintln(output, image)
	}

	return nil
}
