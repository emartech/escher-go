package mock_test

import (
	"testing"

	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/validator/mock"
	"github.com/stretchr/testify/assert"
)

func TestMockingReceivedRequestAvailable(t *testing.T) {

	mock := mock.New()

	req1 := requestBy(t, "GET", "/test", "Hello, World!")
	exampleKeyDB1 := keydb.NewByKeyValuePair("Foo", "Baz")
	mandatoryHeaders1 := []string{"X-Company-Stuff"}

	useValidator(t, mock, req1, exampleKeyDB1, mandatoryHeaders1)

	assert.Equal(t, 1, len(mock.ReceivedValidations))
	assert.Equal(t, req1, mock.ReceivedValidations[0].Request)
	assert.Equal(t, mandatoryHeaders1, mock.ReceivedValidations[0].MandatoryHeaders)
	assert.Equal(t, exampleKeyDB1, mock.ReceivedValidations[0].KeyDB)

	req2 := requestBy(t, "POST", "/monitoring/healthcheck", "Hello, You! :)")
	exampleKeyDB2 := keydb.NewByKeyValuePair("Baz", "Foo")
	mandatoryHeaders2 := []string{"X-Company-Stuff-Goes-Here"}

	useValidator(t, mock, req2, exampleKeyDB2, mandatoryHeaders2)

	assert.Equal(t, 2, len(mock.ReceivedValidations))
	assert.Equal(t, req2, mock.ReceivedValidations[1].Request)
	assert.Equal(t, mandatoryHeaders2, mock.ReceivedValidations[1].MandatoryHeaders)
	assert.Equal(t, exampleKeyDB2, mock.ReceivedValidations[1].KeyDB)

}
