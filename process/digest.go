package process

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/docker"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
)

// toCanonicalImage returns the reference.Canonical version of the passed image
// and will look it up in GCR if required. Returns a reference.Canonical or an
// error.
//
// This is because Binary Authorization only supports canonical image
// references, as a non-canonical image reference could refer to multiple
// versions of the same image (with different contents).
func toCanonicalReference(ctx context.Context, auth voucher.Auth, image string) (reference.Canonical, error) {
	namedRef, err := parseReference(image)
	if nil != err {
		return nil, fmt.Errorf("getting reference failed: %s", err)
	}

	if canonicalRef, ok := namedRef.(reference.Canonical); ok {
		return canonicalRef, nil
	}

	authClient, err := voucher.AuthToClient(ctx, auth, namedRef)
	if nil != err {
		return nil, fmt.Errorf("creating authenticated client failed: %s", err)
	}

	return lookupCanonicalRef(authClient, namedRef)
}

// lookupCanonicalRef gets the canonical image reference for the passed
// image reference.
func lookupCanonicalRef(client *http.Client, ref reference.Reference) (reference.Canonical, error) {
	if taggedRef, ok := ref.(reference.NamedTagged); ok {
		imageDigest, err := docker.GetDigestFromTagged(client, taggedRef)
		if nil != err {
			return nil, fmt.Errorf("getting digest from tag failed: %s", err)
		}
		canonicalRef, err := reference.WithDigest(taggedRef, digest.Digest(imageDigest))
		if nil != err {
			return nil, fmt.Errorf("making canonical reference failed: %s", err)
		}
		return canonicalRef, nil
	}
	return nil, fmt.Errorf("reference cannot be converted to a canonical reference")
}
