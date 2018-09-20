package nobody

import (
	"testing"

	"github.com/Shopify/voucher"
	vtesting "github.com/Shopify/voucher/testing"
)

func TestNobodyCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	nobodyCheck := new(check)
	nobodyCheck.SetAuth(auth)

	i := vtesting.NewTestReference(t)

	pass, err := nobodyCheck.Check(i)

	if nil != err {
		t.Errorf("check failed with error: %s", err)
	}

	if pass {
		t.Error("check passed when it should have failed")
	}
}

func TestNobodyCheckWithNoAuth(t *testing.T) {
	i := vtesting.NewTestReference(t)

	nobodyCheck := new(check)

	// run check without setting up Auth.
	pass, err := nobodyCheck.Check(i)
	if err != voucher.ErrNoAuth {
		t.Fatal("check should have failed due to lack of Auth, but didn't")
	}
	if pass {
		t.Error("check passed when it should have failed due to no Auth")
	}
}
