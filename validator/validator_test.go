package validator_test

import (
	"testing"

	escher "github.com/adamluzsi/escher-go"
	. "github.com/adamluzsi/escher-go/testing"
	"github.com/adamluzsi/escher-go/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	t.Log("Authenticate the incoming request")
	EachTestConfigFor(t, "authenticate", func(config escher.Config, testConfig TestConfig) bool {
		apiKeyID, err := validator.New(config).Validate(testConfig.Request, testConfig.KeyDB(), nil)

		expectedErrorMessage := testConfig.Expected.Error

		if expectedErrorMessage != "" {
			if err == nil {
				t.Fatal("error object expected, but nothing was returned (" + expectedErrorMessage + ")")
			}

			return assert.Equal(t, expectedErrorMessage, err.Error(), expectedErrorMessage)
		}

		if err != nil {
			t.Fatal("There shouldn't be any error but the following received: " + err.Error())
		}

		if testConfig.Expected.APIKeyID != "" {
			assert.Equal(t, testConfig.Expected.APIKeyID, apiKeyID)
		}

		return true
	})
}

// func TestSignThisWeirdStuff(t *testing.T) {
// 	EachTestConfigFor(t, "authenticate", func(config escher.Config, testConfig TestConfig) bool {
// 		if testConfig.Title == "should check if the date is in the allowed range" {

// 			s := signer.New(config)

// 			signature := s.GenerateSignature(testConfig.Request, testConfig.HeadersToSign)

// 			if "06fc6d7f2ff5587b8a7dd9411481b4901aba6cf28387efc1bc8cc3c13d543a30" != signature {
// 				t.Fatal("ohh...")
// 			}

// 		}
// 		return true
// 	})
// }
