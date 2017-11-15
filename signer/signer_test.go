package signer_test

import (
	"testing"

	"github.com/EscherAuth/escher/request"

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

func TestSignRequest_HasConfiguredDate_DateHeaderIsInUTCFormat(t *testing.T) {
	t.Log("SignRequest with configured date should add date header in UTC format")
	c := config.Config{
		Date:        "20171114T212223+01",
		AccessKeyId: "dummy key",
		ApiSecret:   "dummy secret",
	}
	config.SetDefaults(&c)

	req := request.New(
		"GET",
		"/",
		[][2]string{{"Host", "example.com"}},
		"",
		300)

	signedRequest, _ := signer.New(c).SignRequest(req, []string{})
	actualDateHeader, _ := signedRequest.Headers().Get("X-Escher-Date")

	assert.Equal(t, "20171114T202223Z", actualDateHeader)
}

func TestSignRequest_WithoutConfiguredDate_DateHeaderIsInUTCFormat(t *testing.T) {
	t.Log("SignRequest without configured date should add date header in UTC format")
	c := config.Config{
		AccessKeyId: "dummy key",
		ApiSecret:   "dummy secret",
	}
	config.SetDefaults(&c)

	assertSignedRequestDateHeaderIsUTC(t, c)
}

func TestSignRequest_WithInvalidConfiguredDate_DateHeaderIsInUTCFormat(t *testing.T) {
	t.Log("SignRequest with malformed configured date should add date header in UTC format")
	c := config.Config{
		Date:        "invalid date",
		AccessKeyId: "dummy key",
		ApiSecret:   "dummy secret",
	}
	config.SetDefaults(&c)

	assertSignedRequestDateHeaderIsUTC(t, c)
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

func assertSignedRequestDateHeaderIsUTC(t *testing.T, c config.Config) {
	req := request.New(
		"GET",
		"/",
		[][2]string{{"Host", "example.com"}},
		"",
		300)

	signedRequest, _ := signer.New(c).SignRequest(req, []string{})
	actualDateHeader, _ := signedRequest.Headers().Get("X-Escher-Date")

	assert.Regexp(t, "^\\d{8}T\\d{6}Z$", actualDateHeader)

}
