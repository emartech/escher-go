package mock_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/validator"
	"github.com/EscherAuth/escher/validator/mock"
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
	}

	for _, expectedError := range testCases {
		t.Run(fmt.Sprintf(`when the expected error is %q`, expectedError.Error()), func(t *testing.T) {
			mock.AddValidationResult(`TEST`, expectedError)
			testWithValidator(t, mock, expectedError)
		})
	}

	t.Run(`when no error expected from the mock`, func(t *testing.T) {
		mock.AddValidationResult(`TEST`, nil)

		_, err := mock.Validate(
			requestBy(t, "GET", "/test", "Hello, World!"),
			keydb.NewByKeyValuePair("Foo", "Baz"),
			[]string{"X-Company-Stuff"})

		assert.Nil(t, err)
	})

}

func testWithValidator(t testing.TB, v validator.Validator, expectedError error) {
	req := requestBy(t, "GET", "/test", "Hello, World!")
	exampleKeyDB := keydb.NewByKeyValuePair("Foo", "Baz")
	mandatoryHeaders := []string{"X-Company-Stuff"}

	_, err := v.Validate(req, exampleKeyDB, mandatoryHeaders)
	assert.Equal(t, expectedError, err)
}
