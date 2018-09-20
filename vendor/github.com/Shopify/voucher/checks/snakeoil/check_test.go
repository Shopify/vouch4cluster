package snakeoil

import (
	"strings"
	"testing"

	"github.com/Shopify/voucher"
	vtesting "github.com/Shopify/voucher/testing"
)

func TestSnakeoilWithBadScanner(t *testing.T) {
	check := new(check)

	i, err := voucher.NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	if nil != err {
		t.Fatalf("failed to get ImageData: %s", err)
	}

	status, err := check.Check(i)

	if err != ErrNoScanner {
		t.Errorf("got wrong error for check: %s", err)
	}

	if status {
		t.Error("check passed when it was not technically possible")
	}
}

func TestSnakeoil(t *testing.T) {
	check := new(check)

	i, err := voucher.NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	if nil != err {
		t.Fatalf("failed to get ImageData: %s", err)
	}

	check.SetScanner(vtesting.NewScanner(t))

	status, err := check.Check(i)
	if err != nil {
		t.Errorf("check failed with error: %s", err)
	}

	if !status {
		t.Error("check failed when it should have passed")
	}
}

func TestSnakeoilWithVulnerabilities(t *testing.T) {
	check := new(check)

	i, err := voucher.NewImageData("gcr.io/path/to/image@sha256:97db2bc359ccc94d3b2d6f5daa4173e9e91c513b0dcd961408adbb95ec5e5ce5")
	if nil != err {
		t.Fatalf("failed to get ImageData: %s", err)
	}

	scanner := vtesting.NewScanner(t,
		voucher.Vulnerability{
			Name:        "cve-the-worst",
			Description: "it's really bad",
			Severity:    voucher.CriticalSeverity,
		},
		voucher.Vulnerability{
			Name:        "cve-this-is-fine",
			Description: "it's fine",
			Severity:    voucher.NegligibleSeverity,
		},
	)

	check.SetScanner(scanner)

	status, err := check.Check(i)
	if err == nil {
		t.Fatalf("check returned no errors, when it should have")
	}

	if !strings.HasPrefix(err.Error(), "vulnernable to 2 vulnerabilities:") {
		t.Errorf("error message is incorrectly formatted: %s", err)
	}
	if !strings.Contains(err.Error(), "cve-the-worst (critical)") || !strings.Contains(err.Error(), "cve-this-is-fine (negligible)") {
		t.Errorf("error message is incorrectly formatted: %s", err)

	}
	if status {
		t.Error("check passed when it should have failed")
	}
}
