package docker

import (
	"testing"

	vtesting "github.com/Shopify/voucher/testing"
	"github.com/stretchr/testify/assert"
)

func TestRequestManifest(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	manifest, err := RequestManifest(client, ref)
	if nil != err {
		t.Fatalf("failed to get manifest: %s", err)
	}

	assert.Equal(t, vtesting.NewTestManifest(), manifest)
}

func TestRequestBadManifest(t *testing.T) {
	ref := vtesting.NewBadTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	_, err := RequestManifest(client, ref)
	assert.NotNilf(t, err, "should have failed to get manifest, but didn't")
	assert.Contains(t, err.Error(), "failed to load resource with status \"404 Not Found\":")
}
