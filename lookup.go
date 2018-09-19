package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/auth/google"
	"github.com/Shopify/voucher/client"
	"github.com/docker/distribution/reference"
)

func lookupAndAttest(images ...string) error {

	auth := google.NewAuth()

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	for _, image := range images {
		err := processImage(ctx, auth, image)
		if nil != err {
			return fmt.Errorf("processing image failed: %s", err)
		}

	}

	return nil
}

func processImage(ctx context.Context, auth voucher.Auth, image string) error {
	if !strings.HasPrefix(image, "gcr.io") {
		return fmt.Errorf("image is not in our registry: %s", image)
	}

	fmt.Printf("Attesting %s\n", image)

	namedRef, err := parseReference(image)
	if nil != err {
		return fmt.Errorf("getting reference failed: %s", err)
	}

	voucherClient, err := auth.ToClient(ctx, namedRef)
	if nil != err {
		return fmt.Errorf("creating authenticated client failed: %s", err)
	}

	canonicalRef, err := getCanonicalReference(voucherClient, namedRef)
	if nil != err {
		return fmt.Errorf("getting image digest failed: %s", err)
	}
	fmt.Printf(" - Canonical Image: %s\n", canonicalRef.String())

	voucherResp, err := client.SignImage("https://voucher-internal.shopifycloud.com", canonicalRef, "all")
	if nil != err {
		return fmt.Errorf("signing image failed: %s", err)
	}

	fmt.Println(formatResponse(&voucherResp))

	return nil
}

func parseReference(image string) (reference.Named, error) {
	var namedRef reference.Named
	var ok bool

	ref, err := reference.Parse(image)
	if nil != err {
		return namedRef, fmt.Errorf("parsing image reference failed: %s", err)
	}

	if namedRef, ok = ref.(reference.Named); !ok {
		return namedRef, fmt.Errorf("couldn't get named version of reference: %s", err)
	}

	return namedRef, nil
}

// formatResponse returns the response as a string.
func formatResponse(resp *voucher.Response) string {
	output := "checks status: \n"
	for _, result := range resp.Results {
		if result.Success && "" == result.Err {
			continue
		}
		if !result.Success {
			output += fmt.Sprintf(" - %s failed", result.Name)
		} else if !result.Attested {
			output += fmt.Sprintf(" - %s wasn't attested", result.Name)
		}

		if "" != result.Err {
			output += fmt.Sprintf(", err: %s", result.Err)
		}
		output += "\n"
	}

	return output
}
