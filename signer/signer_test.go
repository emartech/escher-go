package signer_test

import (
	"encoding/json"
	"testing"

	"github.com/EscherAuth/escher"
	"github.com/EscherAuth/escher/signer"
	. "github.com/EscherAuth/escher/testing"
	"github.com/stretchr/testify/assert"
)

func TestCanonicalizeRequest(t *testing.T) {
	t.Log("CanonicalizeRequest should return with a proper string")
	EachTestConfigFor(t, []string{"signRequest"}, []string{"error"}, func(config escher.Config, testConfig TestConfig) bool {
		canonicalizedRequest := signer.New(config).CanonicalizeRequest(&testConfig.Request, testConfig.HeadersToSign)

		return assert.Equal(t, testConfig.Expected.CanonicalizedRequest, canonicalizedRequest, "canonicalizedRequest should be eq")
	})
}

func TestGetStringToSign(t *testing.T) {
	t.Log("GetStringToSign should return with a proper string")
	EachTestConfigFor(t, []string{"signRequest"}, []string{"error"}, func(config escher.Config, testConfig TestConfig) bool {
		stringToSign := signer.New(config).GetStringToSign(&testConfig.Request, testConfig.HeadersToSign)

		return assert.Equal(t, stringToSign, testConfig.Expected.StringToSign, "stringToSign expected to eq with the test config expectation")
	})
}

func TestGenerateHeader(t *testing.T) {
	t.Log("GenerateHeader should return with a proper string")
	EachTestConfigFor(t, []string{"signRequest"}, []string{}, func(config escher.Config, testConfig TestConfig) bool {
		if testConfig.Expected.AuthHeader == "" {
			return false
		}

		authHeader := signer.New(config).GenerateHeader(&testConfig.Request, testConfig.HeadersToSign)
		return assert.Equal(t, testConfig.Expected.AuthHeader, authHeader, "authHeader generation failed")
	})
}

func TestSignRequestHappyPath(t *testing.T) {
	t.Log("SignRequest should return with a properly signed request")
	EachTestConfigFor(t, []string{"signRequest"}, []string{"error"}, func(config escher.Config, testConfig TestConfig) bool {
		signedRequest, err := signer.New(config).SignRequest(&testConfig.Request, testConfig.HeadersToSign)

		if !assert.NoError(t, err) {
			return false
		}

		if !assert.NotNil(t, signedRequest) {
			return false
		}

		return assert.Equal(t, testConfig.Expected.Request, *signedRequest)
	})
}

func TestSignRequestError(t *testing.T) {
	t.Log("SignRequest should return error about what was wrong with the given request to sign")
	EachTestConfigFor(t, []string{"signRequest", "error"}, []string{}, func(config escher.Config, testConfig TestConfig) bool {
		_, err := signer.New(config).SignRequest(&testConfig.Request, testConfig.HeadersToSign)

		return assert.EqualError(t, err, testConfig.Expected.Error)
	})
}

func TestSignedURLBy(t *testing.T) {
	t.Log("SignRequest should return with a properly signed request")
	EachTestConfigFor(t, []string{"presignurl"}, []string{}, func(config escher.Config, testConfig TestConfig) bool {
		r := &testConfig.Request

		signedURLStr, err := signer.New(config).SignedURLBy(r.Method(), r.RawURL(), r.Expires())

		if err != nil {
			t.Fatal(err)
		}

		return assert.Equal(t, testConfig.Expected.URL, signedURLStr)
	})
}

func escapeToJSONStringFormat(s string) string {
	bs, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
