package validator_test

import (
	"testing"

	escher "github.com/EscherAuth/escher"
	"github.com/EscherAuth/escher/signer"
	. "github.com/EscherAuth/escher/testing"
	"github.com/EscherAuth/escher/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	t.Log("Authenticate the incoming request")
	EachTestConfigFor(t, "authenticate", func(config escher.Config, testConfig TestConfig) bool {

		escherValidator := validator.New(config)
		apiKeyID, err := escherValidator.Validate(testConfig.Request, testConfig.KeyDB(), nil)
		expectedErrorMessage := testConfig.Expected.Error

		if expectedErrorMessage != "" {
			if err == nil {
				t.Fatal("error object expected, but nothing was returned (" + expectedErrorMessage + ")")
			}

			return assert.Equal(t, expectedErrorMessage, err.Error(), expectedErrorMessage)
		}

		if err != nil {
			t.Log("There shouldn't be any error but the following received: " + err.Error())
			escherSigner := signer.New(config)
			canonizedRequest := escherSigner.CanonicalizeRequest(testConfig.Request, testConfig.HeadersToSign)
			t.Log("\n" + canonizedRequest)
			t.FailNow()
		}

		if testConfig.Expected.APIKeyID != "" {
			assert.Equal(t, testConfig.Expected.APIKeyID, apiKeyID)
		}

		return true
	})
}
