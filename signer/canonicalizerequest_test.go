package signer_test

import (
	"testing"

	"github.com/EscherAuth/escher/signer"

	"github.com/EscherAuth/escher/config"
	. "github.com/EscherAuth/escher/testing/cases"
	"github.com/stretchr/testify/assert"
)

func TestCanonicalizeRequest(t *testing.T) {
	t.Log("CanonicalizeRequest should return with a proper string")
	EachTestConfigFor(t, []string{"signRequest"}, []string{"error"}, func(t *testing.T, c config.Config, testConfig TestConfig) bool {
		canonicalizedRequest := signer.New(c).CanonicalizeRequest(&testConfig.Request, testConfig.HeadersToSign)

		return assert.Equal(t, testConfig.Expected.CanonicalizedRequest, canonicalizedRequest, "canonicalizedRequest should be eq")
	})
}
