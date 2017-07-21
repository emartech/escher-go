package mock_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/EscherAuth/escher/signer/mock"
	"github.com/stretchr/testify/assert"
)

func TestSignedURLBy_NothingAddedToTheMock_DefaultValueReturned(t *testing.T) {
	signer := mock.New()

	actually, err := signer.SignedURLBy("GET", "/test", 42)

	assert.Nil(t, err)
	assert.Equal(t, "/test", actually)

}

func TestSignedURLBy_OneElementAdded_ElementReturned(t *testing.T) {
	signer := mock.New()

	signer.AddSignURLResponse("/signed", nil)

	actually, err := signer.SignedURLBy("GET", "/unsigned", 42)

	assert.Nil(t, err)
	assert.Equal(t, "/signed", actually)

}

func TestSignedURLBy_MultipleElementAdded_ElementsReturnedAndAfterLast(t *testing.T) {

	signer := mock.New()
	times := rand.Intn(42)

	for i := 0; i < times; i++ {
		url := fmt.Sprintf("/signed/%v", i)
		signer.AddSignURLResponse(url, errors.New(url))
	}

	for i := 0; i < times; i++ {
		url := fmt.Sprintf("/signed/%v", i)
		actually, err := signer.SignedURLBy("GET", "/src", 42)
		assert.Error(t, err, url)
		assert.Equal(t, url, actually)
	}

	actually, err := signer.SignedURLBy("GET", "/src", 42)
	assert.Nil(t, err)
	assert.Equal(t, "/src", actually)

}
