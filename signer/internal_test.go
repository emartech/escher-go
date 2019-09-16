package signer

import (
	"testing"

	"github.com/EscherAuth/escher/config"
	. "github.com/EscherAuth/escher/testing/cases"
	"github.com/stretchr/testify/assert"
)

func NewSubject(c config.Config) *signer {
	return &signer{c}
}

func TestGetStringToSign(t *testing.T) {
	t.Log("GetStringToSign should return with a proper string")
	EachTestConfigFor(t, []string{"signRequest"}, []string{"error"}, func(t *testing.T, c config.Config, testConfig TestConfig) bool {
		stringToSign := NewSubject(c).getStringToSign(&testConfig.Request, testConfig.HeadersToSign)

		return assert.Equal(t, stringToSign, testConfig.Expected.StringToSign, "stringToSign expected to eq with the test config expectation")
	})
}

func TestGenerateHeader(t *testing.T) {
	t.Log("GenerateHeader should return with a proper string")
	EachTestConfigFor(t, []string{"signRequest"}, []string{}, func(t *testing.T, c config.Config, testConfig TestConfig) bool {
		if testConfig.Expected.AuthHeader == "" {
			return false
		}

		authHeader := NewSubject(c).generateHeader(&testConfig.Request, testConfig.HeadersToSign)
		return assert.Equal(t, testConfig.Expected.AuthHeader, authHeader, "authHeader generation failed")
	})
}
