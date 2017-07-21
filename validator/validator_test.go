package validator_test

import (
	"testing"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/signer"
	. "github.com/EscherAuth/escher/testing/cases"
	"github.com/EscherAuth/escher/validator"
	"github.com/stretchr/testify/assert"
)

// TODO presign url stuff to be tested

func isHappyPath(t testing.TB, validatorErr error, c config.Config, testConfig TestConfig) bool {

	if testConfig.Expected.Error != "" {
		return false
	}

	if validatorErr == nil {
		return true
	}

	t.Log("There shouldn't be any error but the following received: " + validatorErr.Error())
	t.Log("\n" + signer.New(c).canonicalizeRequest(&testConfig.Request, testConfig.HeadersToSign))
	t.Fail()

	return false
}

func TestValidateRequest(t *testing.T) {
	t.Log("Authenticate the incoming request")
	EachTestConfigFor(t, []string{"authenticate"}, []string{}, func(c config.Config, testConfig TestConfig) bool {

		apiKeyID, err := validator.New(c).Validate(&testConfig.Request, testConfig.KeyDB(), nil)

		if !isHappyPath(t, err, c, testConfig) {
			t.Log("not happy path case")
			return false
		}

		if testConfig.Expected.APIKeyID == "" {
			t.Log(testConfig.FilePath)
			t.Log(testConfig.FilePath)
			t.Log(testConfig.FilePath)
			t.Log(testConfig.FilePath)
			t.Log(testConfig.FilePath)
		}

		assert.Equal(t, testConfig.Expected.APIKeyID, apiKeyID)

		return true
	})
}

func TestValidateErrorCases(t *testing.T) {
	t.Log("Authenticate the incoming request")
	EachTestConfigFor(t, []string{"authenticate", "error"}, []string{}, func(c config.Config, testConfig TestConfig) bool {

		_, err := validator.New(c).Validate(&testConfig.Request, testConfig.KeyDB(), testConfig.MandatorySignedHeaders)

		expectedErrorMessage := testConfig.Expected.Error

		if expectedErrorMessage == "" {
			t.Log("no expectedErrorMessage found, skipping test")
			return false
		}

		if err == nil {
			t.Error("error object expected, but nothing was returned (" + expectedErrorMessage + ")")
			return false
		}

		return assert.Equal(t, expectedErrorMessage, err.Error(), expectedErrorMessage)

	})
}
