package mock_test

import (
	"testing"

	"github.com/EscherAuth/escher/validator"
	"github.com/EscherAuth/escher/validator/mock"
	"github.com/stretchr/testify/assert"
)

func TestErrorRecording(t *testing.T) {

	mock := mock.New()

	testCases := []error{
		validator.MissingDateParam,
		validator.InvalidEscherKey,
		validator.AlgorithmNotAllowed,
		validator.RequestMethodIsInvalid,
		validator.POSTRequestBodyIsEmpty,
		validator.CredentialDateNotMatching,
		validator.RequestDateNotAcceptable,
		validator.AuthorizationDateIsInvalid,
		validator.InvalidCredentialScope,
		validator.HostHeaderNotSigned,
		validator.DateHeaderIsNotSigned,
		validator.SignatureDoesNotMatch,
		validator.MalformedAuthHeader,
		validator.SigningParamNotFound,
		validator.SchemaInURLNotAllowed,
		validator.HTTPSchemaFoundInTheURL,
		nil,
	}

	apiKey := "TEST"
	for _, expectedError := range testCases {
		t.Logf("when expected Error is %v", expectedError)
		mock.AddValidationResult(apiKey, expectedError)
		testWithValidator(t, mock, apiKey, expectedError)
	}

}

func testWithValidator(t testing.TB, v validator.Validator, expectedApiKey string, expectedError error) {

	req := requestBy(t, "GET", "/test", "Hello, World!")
	exampleKeyDB := keyDBBy("Foo", "Baz")
	mandatoryHeaders := []string{"X-Company-Stuff"}

	_, err := v.Validate(req, exampleKeyDB, mandatoryHeaders)

	if expectedError == nil {
		assert.Nil(t, err)
	} else {
		assert.Equal(t, expectedError, err)
	}

}
