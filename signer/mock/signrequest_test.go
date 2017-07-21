package mock_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/signer/mock"
	"github.com/stretchr/testify/assert"
)

func TestSignRequest_NothingAddedToTheMock_DefaultValueReturned(t *testing.T) {
	signer := mock.New()

	expectedRequest := requestBy("GET", "/")

	actuallyRequest, err := signer.SignRequest(expectedRequest, []string{})
	assert.Nil(t, err)

	assert.Equal(t, expectedRequest, actuallyRequest)

}

func TestSignRequest_OneElementAdded_ElementReturned(t *testing.T) {
	signer := mock.New()

	expectedRequest := requestBy("GET", "/signedReq")
	signer.AddSignRequestResponse(expectedRequest, nil)

	actuallyRequest, err := signer.SignRequest(requestBy("GET", "/original"), []string{})
	assert.Nil(t, err)

	assert.Equal(t, expectedRequest, actuallyRequest)

}

func TestSignRequest_MultipleElementAdded_ElementsReturnedAndAfterLast(t *testing.T) {

	signer := mock.New()
	times := rand.Intn(42)

	expectedRequests := make([]*request.Request, 0, times)

	for i := 0; i < times; i++ {
		url := fmt.Sprintf("/signed/%v", i)
		expectedRequest := requestBy("GET", url)
		signer.AddSignRequestResponse(expectedRequest, errors.New(url))
		expectedRequests = append(expectedRequests, expectedRequest)
	}

	for i, expectedRequest := range expectedRequests {
		actuallyRequest, err := signer.SignRequest(requestBy("GET", "/src"), []string{})
		url := fmt.Sprintf("/signed/%v", i)
		assert.Error(t, err, url)
		assert.Equal(t, expectedRequest, actuallyRequest)
	}

	requestForLastCall := requestBy("GET", "/original")
	actuallyRequest, err := signer.SignRequest(requestForLastCall, []string{})

	assert.Nil(t, err)
	assert.Equal(t, requestForLastCall, actuallyRequest)

}
