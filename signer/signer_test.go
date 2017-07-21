package signer_test

import (
	"testing"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/signer"
	. "github.com/EscherAuth/escher/testing/cases"
	"github.com/stretchr/testify/assert"
)

func TestSignRequest_RequestIsValid_SignedRequestReturned(t *testing.T) {
	t.Log("SignRequest should return with a properly signed request")
	EachTestConfigFor(t, []string{"signRequest"}, []string{"error"}, func(c config.Config, testConfig TestConfig) bool {
		signedRequest, err := signer.New(c).SignRequest(&testConfig.Request, testConfig.HeadersToSign)

		if !assert.NoError(t, err) {
			return false
		}

		if !assert.NotNil(t, signedRequest) {
			return false
		}

		return assert.Equal(t, testConfig.Expected.Request, *signedRequest)
	})
}

func TestSignRequest_ErrorOnSigning_ErrorReturnedThatTellsTheProblem(t *testing.T) {
	t.Log("SignRequest should return error about what was wrong with the given request to sign")
	EachTestConfigFor(t, []string{"signRequest", "error"}, []string{}, func(c config.Config, testConfig TestConfig) bool {
		_, err := signer.New(c).SignRequest(&testConfig.Request, testConfig.HeadersToSign)

		return assert.EqualError(t, err, testConfig.Expected.Error)
	})
}

func TestSignedURLBy(t *testing.T) {
	t.Log("SignRequest should return with a properly signed request")
	EachTestConfigFor(t, []string{"presignurl"}, []string{}, func(c config.Config, testConfig TestConfig) bool {
		r := &testConfig.Request

		signedURLStr, err := signer.New(c).SignedURLBy(r.Method(), r.RawURL(), r.Expires())

		if err != nil {
			t.Fatal(err)
		}

		return assert.Equal(t, testConfig.Expected.URL, signedURLStr)
	})
}
