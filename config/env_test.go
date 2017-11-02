package config_test

import (
	"reflect"
	"testing"

	"github.com/EscherAuth/escher/config"
	. "github.com/EscherAuth/escher/testing/env"
	"github.com/stretchr/testify/assert"
)

var exampleEscherConfig = `
{
    "vendorKey": 		"ZZ",
    "algoPrefix": 		"VV",
    "hashAlgo": 		"SHA512",
	"credentialScope": 	"us-east-1/host/aws4_request",
	"authHeaderName": "X-Escher-Auth",
	"dateHeaderName": "X-Escher-Date"
}
`

func TestNewFromENV_ConfigJSONIsPresentInTheEnv(t *testing.T) {
	defer SetEnvForTheTest(t, "ESCHER_CONFIG", exampleEscherConfig)()

	config, err := config.NewFromENV()

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "ZZ", config.VendorKey)
	assert.Equal(t, "VV", config.AlgoPrefix)
	assert.Equal(t, "SHA512", config.HashAlgo)
	assert.Equal(t, "us-east-1/host/aws4_request", config.CredentialScope)
	assert.Equal(t, "", config.Date)

}

func TestNewFromENV_ValidJSONIsPresentButOnlyCredentialScopeIsProvided(t *testing.T) {
	defer SetEnvForTheTest(t, "ESCHER_CONFIG", `{"credentialScope": "a/b/c"}`)()

	config, err := config.NewFromENV()

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "ESR", config.AlgoPrefix)
	assert.Equal(t, "SHA256", config.HashAlgo)
	assert.Equal(t, "Escher", config.VendorKey)
	assert.Equal(t, "X-Escher-Auth", config.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", config.DateHeaderName)
	assert.Equal(t, "a/b/c", config.CredentialScope)

}

func TestNewFromENV_EveryValueIsProvidedInEnvVariables(t *testing.T) {
	defer UnsetEnvForTheTest(t, "ESCHER_CONFIG")()
	defer SetEnvForTheTest(t, "ESCHER_ALGO_PREFIX", "ALGO_PREFIX")()
	defer SetEnvForTheTest(t, "ESCHER_HASH_ALGO", "HASH_ALGO")()
	defer SetEnvForTheTest(t, "ESCHER_VENDOR_KEY", "VENDOR_KEY")()
	defer SetEnvForTheTest(t, "ESCHER_AUTH_HEADER_NAME", "AUTH_HEADER_NAME")()
	defer SetEnvForTheTest(t, "ESCHER_DATE_HEADER_NAME", "DATE_HEADER_NAME")()
	defer SetEnvForTheTest(t, "ESCHER_CREDENTIAL_SCOPE", "CREDENTIAL_SCOPE")()

	config, err := config.NewFromENV()

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "ALGO_PREFIX", config.AlgoPrefix)
	assert.Equal(t, "HASH_ALGO", config.HashAlgo)
	assert.Equal(t, "VENDOR_KEY", config.VendorKey)
	assert.Equal(t, "AUTH_HEADER_NAME", config.AuthHeaderName)
	assert.Equal(t, "DATE_HEADER_NAME", config.DateHeaderName)
	assert.Equal(t, "CREDENTIAL_SCOPE", config.CredentialScope)
}

func TestNewFromENV_OneValueAtLeastProvidedInTheENVWithExplicitValueSetting(t *testing.T) {
	defer UnsetEnvForTheTest(t, "ESCHER_CONFIG")()
	defer SetEnvForTheTest(t, "ESCHER_CREDENTIAL_SCOPE", "TEST")()

	cases := map[string]string{
		"ESCHER_ALGO_PREFIX":      "AlgoPrefix",
		"ESCHER_HASH_ALGO":        "HashAlgo",
		"ESCHER_VENDOR_KEY":       "VendorKey",
		"ESCHER_AUTH_HEADER_NAME": "AuthHeaderName",
		"ESCHER_DATE_HEADER_NAME": "DateHeaderName",
		"ESCHER_CREDENTIAL_SCOPE": "CredentialScope",
		"ESCHER_API_SECRET":       "ApiSecret",
		"ESCHER_ACCESS_KEY_ID":    "AccessKeyId",
	}

	for envKey, envValue := range cases {
		tearDown := SetEnvForTheTest(t, envKey, envValue)

		config, err := config.NewFromENV()

		if err != nil {
			t.Fatal(err)
		}

		r := reflect.ValueOf(config)
		actuallyValue := reflect.Indirect(r).FieldByName(envValue).String()
		assert.Equal(t, envValue, actuallyValue)

		tearDown()
	}

}

func TestNewFromENV_InvalidJSONConfig_ErrorReturned(t *testing.T) {
	defer SetEnvForTheTest(t, "ESCHER_CONFIG", `{credentialScope:"not/valid/json"}`)()

	_, err := config.NewFromENV()

	assert.NotNil(t, err)
}

func TestNewFromENV_CredentialScopeIsNotGiven_ErrorIsReturned(t *testing.T) {
	defer UnsetEnvForTheTest(t, "ESCHER_CREDENTIAL_SCOPE")()

	_, err := config.NewFromENV()

	assert.Error(t, err, "Credential Scope is missing")
}
