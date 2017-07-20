package config_test

import (
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

func TestNewFromENV_ValidJSONIsPresentButWithMissingValues(t *testing.T) {
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
