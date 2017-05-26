package validator_test

import (
	"testing"

	escher "github.com/adamluzsi/escher-go"
	. "github.com/adamluzsi/escher-go/testing"
	"github.com/adamluzsi/escher-go/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	t.Log("Authenticate the incomming request")
	EachTestConfigFor(t, "authenticate", func(config escher.Config, testConfig TestConfig) bool {
		apiKeyID, err := validator.New(config).Validate(testConfig.Request, testConfig.KeyDB(), nil)

		if testConfig.Expected.Error != "" {
			return assert.Equal(t, testConfig.Expected.Error, err.Error(), testConfig.Expected.Error)
		}

		if testConfig.Expected.APIKeyID != "" {
			assert.Equal(t, testConfig.Expected.APIKeyID, apiKeyID)
		}

		return true
	})
}
