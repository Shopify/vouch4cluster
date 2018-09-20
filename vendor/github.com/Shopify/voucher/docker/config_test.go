package docker

import (
	"testing"

	vtesting "github.com/Shopify/voucher/testing"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestRequestConfig(t *testing.T) {
	ref := vtesting.NewTestReference(t)

	client, server := PrepareDockerTest(t, ref)
	defer server.Close()

	config, err := RequestImageConfig(client, ref)
	if nil != err {
		t.Fatalf("failed to get config: %s", err)
	}

	expectedConfig := ImageConfig{
		ContainerConfig: dockerTypes.ExecConfig{
			User: "root",
		},
	}

	assert.Equal(t, expectedConfig, config)

	assert.True(t, config.RunsAsRoot())
}
