package diy

import (
	"strings"
	"testing"

	"github.com/Shopify/voucher"
	vtesting "github.com/Shopify/voucher/testing"
)

func TestDIYCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	i := vtesting.NewTestReference(t)

	diyCheck := new(check)
	diyCheck.SetAuth(vtesting.NewAuth(server))

	pass, err := diyCheck.Check(i)

	if nil != err {
		t.Errorf("check failed with error: %s", err)
	}

	if !pass {
		t.Error("check failed when it should have passed")
	}
}

func TestDIYCheckWithNoAuth(t *testing.T) {
	i := vtesting.NewTestReference(t)

	diyCheck := new(check)

	// run check without setting up Auth.
	pass, err := diyCheck.Check(i)
	if err != voucher.ErrNoAuth {
		t.Fatal("check should have failed due to lack of Auth, but didn't")
	}
	if pass {
		t.Error("check passed when it should have failed due to no Auth")
	}
}

func TestFailingDIYCheck(t *testing.T) {
	server := vtesting.NewTestDockerServer(t)

	auth := vtesting.NewAuth(server)

	diyCheck := new(check)
	diyCheck.SetAuth(auth)

	i := vtesting.NewBadTestReference(t)

	pass, err := diyCheck.Check(i)

	if nil == err {
		t.Fatal("check should have failed with error, but didn't")
	}

	if !strings.Contains(err.Error(), "image doesn't exist") {
		t.Errorf("check error format is incorrect: \"%s\"", err)
	}

	if pass {
		t.Error("check passed when it should have failed")
	}
}
