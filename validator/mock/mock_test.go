package mock_test

import (
	"testing"

	"github.com/EscherAuth/escher/validator"
	"github.com/EscherAuth/escher/validator/mock"
	"github.com/stretchr/testify/assert"
)

func TestMockingValidatorIsABehavesToValidatorInterface(t *testing.T) {
	assert.Implements(t, (*validator.Validator)(nil), mock.New())
}
