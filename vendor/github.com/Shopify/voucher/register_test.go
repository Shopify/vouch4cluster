package voucher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCheckFactory(t *testing.T) {
	factories := make(CheckFactories)

	factories.Register("goodcheck", func() Check {
		return newTestCheck(true)
	})
	factories.Register("badcheck", func() Check {
		return newTestCheck(false)
	})

	checks, err := factories.GetNewChecks("goodcheck", "badcheck")
	if nil != err {
		t.Fatalf("failed to get checks: %s", err)
	}

	assert.Len(t, checks, 2)

	i := newTestImageData(t)

	if assert.NotNil(t, checks["goodcheck"]) {
		ok, checkErr := checks["goodcheck"].Check(i)
		assert.Nil(t, checkErr)
		assert.True(t, ok)
	}

	if assert.NotNil(t, checks["badcheck"]) {
		ok, checkErr := checks["badcheck"].Check(i)
		assert.Nil(t, checkErr)
		assert.False(t, ok)
	}

}

func TestEmptyCheckFactory(t *testing.T) {
	factories := make(CheckFactories)
	_, err := factories.GetNewChecks("nilcheck")
	assert.Contains(t, "requested check \"nilcheck\" does not exist", err.Error())
}
