package mock_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/EscherAuth/escher/signer/mock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSignature_NothingAddedToTheMock_DefaultValueReturned(t *testing.T) {
	signer := mock.New()

	expectedRequest := requestBy("GET", "/")
	actually := signer.GenerateSignature(expectedRequest, []string{})

	assert.Equal(t, "signature", actually)
}

func TestGenerateSignature_OneElementAdded_ElementReturned(t *testing.T) {
	signer := mock.New()

	expected := "/signedReq"
	signer.AddGeneratedSignature(expected)

	actually := signer.GenerateSignature(requestBy("GET", expected), []string{})
	assert.Equal(t, expected, actually)
}

func TestGenerateSignature_MultipleElementAdded_ElementsReturnedAndAfterLast(t *testing.T) {

	signer := mock.New()
	times := rand.Intn(42)

	urls := make([]string, 0, times)

	for i := 0; i < times; i++ {
		url := fmt.Sprintf("/signed/%v", i)
		signer.AddGeneratedSignature(url)
		urls = append(urls, url)
	}

	for _, expected := range urls {
		actually := signer.GenerateSignature(requestBy("GET", "/src"), []string{})
		assert.Equal(t, expected, actually)
	}

	actually := signer.GenerateSignature(requestBy("GET", "/original"), []string{})
	assert.Equal(t, "signature", actually)

}
